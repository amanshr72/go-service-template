package product

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	repo := NewInMemoryRepository()

	p := &Product{
		Name:  "Laptop",
		Price: 1000,
	}

	err := repo.Create(p)

	assert.NoError(t, err)
	assert.Equal(t, 1, p.ID)
}

func TestRepository_GetByID(t *testing.T) {
	repo := NewInMemoryRepository()

	p := &Product{
		Name:  "Laptop",
		Price: 1000,
	}

	_ = repo.Create(p)

	found, err := repo.GetByID(p.ID)

	assert.NoError(t, err)
	assert.Equal(t, p.Name, found.Name)
}

func TestRepository_GetByID_NotFound(t *testing.T) {
	repo := NewInMemoryRepository()

	_, err := repo.GetByID(999)

	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
}

func TestRepository_GetAll(t *testing.T) {
	repo := NewInMemoryRepository()

	_ = repo.Create(&Product{Name: "A", Price: 10})
	_ = repo.Create(&Product{Name: "B", Price: 20})

	products, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, products, 2)
}
