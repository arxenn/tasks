package service

import (
	"fmt"
	"slices"
	"time"

	"github.com/arxenn/tasks/internal/domain"
	"github.com/arxenn/tasks/internal/repository"
)

type Service interface {
	Add(content, priority string) (int, error)
	Done(id int) error
	List(priority string, done bool, num int) ([]domain.Task, error)
	Delete(id int) error
	Clear(all bool) error
}

type service struct {
	repo repository.Repository
}

func (s *service) Add(content, priorityStr string) (int, error) {
	priority, err := domain.StrToTaskPriorityDomain(priorityStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert priority: %w", err)
	}

	id, err := s.repo.Add(domain.Task{
		Content:   content,
		Priority:  priority,
		Status:    domain.TodoTaskStatus,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return 0, fmt.Errorf("adding task to repo failed: %w", err)
	}

	return id, nil
}

func (s *service) Done(id int) error {
	return s.repo.Update(id, domain.Task{
		Status: domain.DoneTaskStatus,
	})
}

func (s *service) List(priority string, done bool, num int) ([]domain.Task, error) {
	var filter domain.TaskFilters

	p, err := domain.StrToTaskPriorityDomain(priority)
	if err == nil {
		filter.Priority = p
	}

	filter.Done = done
	filter.Limit = num

	tasks, err := s.repo.List(filter)
	if err != nil {
		return nil, fmt.Errorf("get tasks list from repo failed: %w", err)
	}

	slices.SortFunc(tasks, func(a, b domain.Task) int {
		return domain.PriorityToIntMap[b.Priority] - domain.PriorityToIntMap[a.Priority]
	})

	// we don't use db's LIMIT functionality because we cannot sort
	// tasks based on priority in db (we store priority as string)
	// so we slice the final list after sort
	if len(tasks) > filter.Limit {
		tasks = tasks[:filter.Limit]
	}

	return tasks, nil
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *service) Clear(all bool) error {
	return s.repo.Clear(all)
}

func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}
