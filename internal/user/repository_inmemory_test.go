package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemory_CreateAndFind(t *testing.T) {
	repo := NewInMemoryRepository()
	u := &User{Name: "Aman", Email: "aman@t.com", IsActive: true}

	err := repo.Create(u)
	assert.NoError(t, err)
	assert.NotZero(t, u.ID)

	found, err := repo.FindByID(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Aman", found.Name)
}

func TestInMemory_FindByActive(t *testing.T) {
	repo := NewInMemoryRepository()
	_ = repo.Create(&User{Name: "Active", Email: "a@t.com", IsActive: true})
	_ = repo.Create(&User{Name: "Inactive", Email: "b@t.com", IsActive: false})

	active, err := repo.FindByActive(true)
	assert.NoError(t, err)
	assert.Len(t, active, 1)
	assert.Equal(t, "Active", active[0].Name)
}

func TestInMemory_Count(t *testing.T) {
	repo := NewInMemoryRepository()
	_ = repo.Create(&User{Name: "A", Email: "a@t.com"})
	_ = repo.Create(&User{Name: "B", Email: "b@t.com"})

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestInMemory_Update(t *testing.T) {
	repo := NewInMemoryRepository()
	u := &User{Name: "Old", Email: "old@t.com"}
	_ = repo.Create(u)

	u.Name = "New"
	err := repo.Update(u)
	assert.NoError(t, err)

	found, _ := repo.FindByID(u.ID)
	assert.Equal(t, "New", found.Name)
}

func TestInMemory_Delete(t *testing.T) {
	repo := NewInMemoryRepository()
	u := &User{Name: "Del", Email: "del@t.com"}
	_ = repo.Create(u)

	err := repo.Delete(u.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(u.ID)
	assert.Error(t, err)
}
