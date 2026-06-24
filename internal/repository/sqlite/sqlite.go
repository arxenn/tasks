package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arxenn/tasks/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidID    = errors.New("invalid task ID")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository() (*SQLiteRepository, error) {
	appName := "tasks"
	appDataDir, err := getAppDataDir(appName)
	if err != nil {
		return nil, fmt.Errorf("failed to get app data dir: %w", err)
	}
	dbPath := fmt.Sprintf("%s/%s.db", appDataDir, appName)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteRepository{db: db}, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		status TEXT NOT NULL,
		priority TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		done_at TIMESTAMP NULL
	);
	`
	_, err := db.Exec(query)
	return err
}

func (r *SQLiteRepository) Add(t domain.Task) (int, error) {
	query := `INSERT INTO tasks (content, status, priority, created_at) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, t.Content, t.Status, t.Priority, t.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(id), nil
}

func (r *SQLiteRepository) List(filters domain.TaskFilters) ([]domain.Task, error) {
	query := `SELECT id, content, status, priority, created_at, done_at FROM tasks WHERE 1=1`
	args := []any{}

	if filters.Priority != "" {
		query += " AND priority = ?"
		args = append(args, filters.Priority)
	}

	if filters.Done {
		query += " AND status = ?"
		args = append(args, domain.DoneTaskStatus)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(&task.ID, &task.Content, &task.Status, &task.Priority, &task.CreatedAt, &task.DoneAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tasks, nil
}

func (r *SQLiteRepository) Update(id int, t domain.Task) error {
	if id <= 0 {
		return ErrInvalidID
	}

	query := strings.Builder{}
	query.WriteString("UPDATE tasks SET ")
	args := []any{}
	updates := []string{}

	var currentTask domain.Task
	err := r.db.QueryRow("SELECT content, status, priority FROM tasks WHERE id = ?", id).
		Scan(&currentTask.Content, &currentTask.Status, &currentTask.Priority)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrTaskNotFound
		}
		return fmt.Errorf("failed to get current task: %w", err)
	}

	if t.Content != "" && t.Content != currentTask.Content {
		updates = append(updates, "content = ?")
		args = append(args, t.Content)
	}

	if t.Status != "" && t.Status != currentTask.Status {
		updates = append(updates, "status = ?")
		args = append(args, t.Status)
		if t.Status == domain.DoneTaskStatus {
			updates = append(updates, "done_at = ?")
			args = append(args, time.Now())
		}
	}

	if t.Priority != "" && t.Priority != currentTask.Priority {
		updates = append(updates, "priority = ?")
		args = append(args, t.Priority)
	}

	if len(updates) == 0 {
		return nil
	}

	query.WriteString(updates[0])
	for i := 1; i < len(updates); i++ {
		query.WriteString(", " + updates[i])
	}
	query.WriteString(" WHERE id = ?")
	args = append(args, id)

	_, err = r.db.Exec(query.String(), args...)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (r *SQLiteRepository) Delete(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	result, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *SQLiteRepository) Clear(all bool) error {
	var query string
	var args []interface{}

	if all {
		query = "DELETE FROM tasks"
	} else {
		query = "DELETE FROM tasks WHERE status = ?"
		args = append(args, domain.DoneTaskStatus)
	}

	if _, err := r.db.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to clear tasks: %w", err)
	}

	return nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
