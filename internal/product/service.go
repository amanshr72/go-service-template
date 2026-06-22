package product

import (
	"errors"
	"fmt"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func validate(req CreateProductRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Price <= 0 {
		return errors.New("price must be positive")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	if req.Category == "" {
		return errors.New("category is required")
	}
	if req.Stock < 0 {
		return errors.New("stock must be >= 0")
	}
	return nil
}

func validateFields(req CreateProductRequest) map[string]interface{} {
	errs := map[string]interface{}{}
	if req.Name == "" {
		errs["name"] = "required"
	}
	if req.Price <= 0 {
		errs["price"] = "must be > 0"
	}
	if req.Description == "" {
		errs["description"] = "required"
	}
	if req.Category == "" {
		errs["category"] = "required"
	}
	if req.Stock < 0 {
		errs["stock"] = "must be >= 0"
	}
	return errs
}

func (s *service) Create(req CreateProductRequest) (*Product, error) {
	if errs := validateFields(req); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}
	p := &Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Category:    req.Category,
		Stock:       req.Stock,
	}
	if err := s.repo.Create(p); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	return p, nil
}

func (s *service) GetAll() ([]Product, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id int64) (*Product, error) {
	return s.repo.GetByID(id)
}

type ValidationError struct {
	Fields map[string]interface{}
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
