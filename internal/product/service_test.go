package product

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func req(name string, price float64) CreateProductRequest {
	return CreateProductRequest{Name: name, Price: price, Description: "desc", Category: "cat", Stock: 1}
}

func TestService_Create(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	p, err := svc.Create(req("MacBook", 120000))
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "MacBook", p.Name)
}

func TestService_Create_EmptyName(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	_, err := svc.Create(req("", 120000))
	assert.Error(t, err)
}

func TestService_Create_InvalidPrice(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	_, err := svc.Create(req("Laptop", 0))
	assert.Error(t, err)
}

func TestService_GetByID(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	created, _ := svc.Create(req("Keyboard", 3000))
	found, err := svc.GetByID(created.Id)
	assert.NoError(t, err)
	assert.Equal(t, created.Id, found.Id)
}

func TestService_GetAll(t *testing.T) {
	svc := NewService(NewInMemoryRepository())
	_, _ = svc.Create(req("A", 10))
	_, _ = svc.Create(req("B", 20))
	products, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, products, 2)
}
