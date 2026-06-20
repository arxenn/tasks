package internal

type Repository interface {
	Add(t Task) (int, error)
	Get(id int) (Task, error)
	List(filters TaskFilters) ([]Task, error)
	Update(id int, t Task) error
	Delete(id int) error
}
