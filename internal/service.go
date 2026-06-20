package internal

type Service interface {
	Add(content, priority string) (int, error)
	Get(id int) (Task, error)
	List(filters TaskFilters) ([]Task, error)
	Update(id int, t Task) error
	Delete(id int) error
}

type service struct {
	repo Repository
}

func (s *service) Add(content, priority string) (int, error) {

	return 0, nil
}

func (s *service) Get(id int) (Task, error) {
	return s.repo.Get(id)
}

func (s *service) List(name, priority, status string) ([]Task, error) {
	return nil, nil
}

func (s *service) Update(id int, task Task) error {
	return nil
}

func (s *service) Delete(id int) error {
	return s.Delete(id)
}
