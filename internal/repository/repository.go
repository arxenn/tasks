package repository

import "github.com/arxenn/tasks/internal/domain"

type Repository interface {
	Add(t domain.Task) (int, error)
	Get(id int) (domain.Task, error)
	List(filters domain.TaskFilters) ([]domain.Task, error)
	Update(id int, t domain.Task) error
	Delete(id int) error
}
