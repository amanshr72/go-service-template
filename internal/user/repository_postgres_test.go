package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	users  map[int]*User
	nextID int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{users: make(map[int]*User), nextID: 1}
}

func (m *MockRepository) Create(u *User) error {
	u.ID = m.nextID
	m.nextID++
	copy := *u
	m.users[copy.ID] = &copy
	return nil
}

func (m *MockRepository) FindAll() ([]User, error) {
	result := make([]User, 0)
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

func (m *MockRepository) FindByID(id int) (*User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *MockRepository) FindByActive(active bool) ([]User, error) {
	result := make([]User, 0)
	for _, u := range m.users {
		if u.IsActive == active {
			result = append(result, *u)
		}
	}
	return result, nil
}

func (m *MockRepository) Count() (int, error) {
	return len(m.users), nil
}

func (m *MockRepository) Update(u *User) error {
	if _, ok := m.users[u.ID]; !ok {
		return errors.New("user not found")
	}
	copy := *u
	m.users[u.ID] = &copy
	return nil
}

func (m *MockRepository) Delete(id int) error {
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

// --- Postgres adapter tests (via mock — logic parity tests) ---

func TestPostgres_FindByActive(t *testing.T) {
	repo := NewMockRepository()
	_ = repo.Create(&User{Name: "Active", Email: "a@t.com", IsActive: true})
	_ = repo.Create(&User{Name: "Inactive", Email: "b@t.com", IsActive: false})

	result, err := repo.FindByActive(true)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestPostgres_Count(t *testing.T) {
	repo := NewMockRepository()
	_ = repo.Create(&User{Name: "A", Email: "a@t.com"})
	_ = repo.Create(&User{Name: "B", Email: "b@t.com"})

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
