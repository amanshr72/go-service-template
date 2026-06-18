package product

import (
	"errors"
	"sync"
)

type InMemoryRepository struct {
	mu     sync.RWMutex
	items  map[int]*Product
	nextID int
}

func NewInMemoryRepository() Repository {
	return &InMemoryRepository{items: make(map[int]*Product), nextID: 1}
}

func (r *InMemoryRepository) Create(p *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = r.nextID
	r.nextID++
	copy := *p
	r.items[copy.ID] = &copy
	return nil
}

func (r *InMemoryRepository) GetAll() ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Product, 0, len(r.items))
	for _, p := range r.items {
		result = append(result, *p)
	}
	return result, nil
}

func (r *InMemoryRepository) GetByID(id int) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.items[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return p, nil
}
