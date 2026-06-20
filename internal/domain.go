package internal

import "time"

type TaskPriority string

const (
	BlockTaskPriority  TaskPriority = "block"
	HighTaskPriority   TaskPriority = "high"
	MediumTaskPriority TaskPriority = "med"
	LowTaskPriority    TaskPriority = "low"
)

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

type TaskFilters struct{}
