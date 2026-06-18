package product

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Create(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	p, err := svc.Create(
		"MacBook",
		120000,
	)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "MacBook", p.Name)
}

func TestService_Create_EmptyName(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	_, err := svc.Create(
		"",
		120000,
	)

	assert.Error(t, err)
	assert.Equal(t, "name is required", err.Error())
}

func TestService_Create_InvalidPrice(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	_, err := svc.Create(
		"Laptop",
		0,
	)

	assert.Error(t, err)
	assert.Equal(t, "price must be positive", err.Error())
}

func TestService_GetByID(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	created, _ := svc.Create(
		"Keyboard",
		3000,
	)

	found, err := svc.GetByID(created.ID)

	assert.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}

func TestService_GetAll(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	_, _ = svc.Create("A", 10)
	_, _ = svc.Create("B", 20)

	products, err := svc.GetAll()

	assert.NoError(t, err)
	assert.Len(t, products, 2)
}
