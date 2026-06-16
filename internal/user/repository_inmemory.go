package user

import (
	"errors"
	"sync"
	"time"
)

type InMemoryRepository struct {
	mu     sync.RWMutex
	users  map[int]*User
	nextID int
}

func NewInMemoryRepository() Repository {
	return &InMemoryRepository{users: make(map[int]*User), nextID: 1}
}

func (r *InMemoryRepository) Create(u *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = r.nextID
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	r.nextID++
	copy := *u
	r.users[copy.ID] = &copy
	return nil
}

func (r *InMemoryRepository) FindAll() ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, *u)
	}
	return result, nil
}

func (r *InMemoryRepository) FindByID(id int) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *InMemoryRepository) FindByActive(active bool) ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]User, 0)
	for _, u := range r.users {
		if u.IsActive == active {
			result = append(result, *u)
		}
	}
	return result, nil
}

func (r *InMemoryRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.users), nil
}

func (r *InMemoryRepository) Update(u *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[u.ID]; !ok {
		return errors.New("user not found")
	}
	u.UpdatedAt = time.Now()
	copy := *u
	r.users[u.ID] = &copy
	return nil
}

func (r *InMemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}
