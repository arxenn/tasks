package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	DefaulListCountNumber = 3

	ListTimeDisplayFormat = "01/02 15:04"
)

type TaskPriority string

const (
	BlockTaskPriority  TaskPriority = "block"
	HighTaskPriority   TaskPriority = "high"
	MediumTaskPriority TaskPriority = "medium"
	LowTaskPriority    TaskPriority = "low"
)

var (
	UndefinedPriorityErr    error = errors.New("undefined priority")
	UndefinedStatusErr      error = errors.New("undefined status")
	ContentCannotBeEmptyErr error = errors.New("content cannot be empty")
)

// stores numerical value of priority for sorting purposes
var PriorityToIntMap = map[TaskPriority]int{
	BlockTaskPriority:  4,
	HighTaskPriority:   3,
	MediumTaskPriority: 2,
	LowTaskPriority:    1,
}

func StrToTaskPriorityDomain(p string) (TaskPriority, error) {
	switch strings.ToLower(p) {
	case "block":
		return BlockTaskPriority, nil
	case "high":
		return HighTaskPriority, nil
	case "medium":
		return MediumTaskPriority, nil
	case "low":
		return LowTaskPriority, nil
	default:
		return "", UndefinedPriorityErr
	}
}

type TaskStatus string

const (
	DoneTaskStatus TaskStatus = "done"
	TodoTaskStatus TaskStatus = "todo"
)

func StrToTaskStatusDomain(s string) (TaskStatus, error) {
	switch strings.ToLower(s) {
	case "todo":
		return TodoTaskStatus, nil
	case "done":
		return DoneTaskStatus, nil
	default:
		return "", UndefinedStatusErr
	}
}

type Task struct {
	ID        int
	Content   string
	Priority  TaskPriority
	Status    TaskStatus
	CreatedAt time.Time
	DoneAt    *time.Time
}

type TaskFilters struct {
	Priority TaskPriority
	Done     bool
	Limit    int
}
