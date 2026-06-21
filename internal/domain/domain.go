package domain

import (
	"strings"
	"time"
)

type TaskPriority string

const (
	BlockTaskPriority  TaskPriority = "block"
	HighTaskPriority   TaskPriority = "high"
	MediumTaskPriority TaskPriority = "med"
	LowTaskPriority    TaskPriority = "low"
)

func StrToTaskPriorityDomain(p string) TaskPriority {
	switch strings.ToLower(p) {
	case "block":
		return BlockTaskPriority
	case "high":
		return HighTaskPriority
	case "med":
		return MediumTaskPriority
	case "low":
		return LowTaskPriority
	default:
		return ""
	}
}

type TaskStatus string

const (
	DoneTaskStatus TaskStatus = "done"
	TodoTaskStatus TaskStatus = "todo"
)

type Task struct {
	ID       int
	Content  string
	Priority TaskPriority
	Status   TaskStatus
	Time     time.Time
}

type TaskFilters struct {
	Status   TaskStatus
	Priority TaskPriority
}
