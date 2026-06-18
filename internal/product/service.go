package product

import "errors"

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(name string, price float64) (*Product, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be positive")
	}
	p := &Product{Name: name, Price: price}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) GetAll() ([]Product, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id int) (*Product, error) {
	return s.repo.GetByID(id)
}
