package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockRepository defined in repository_postgres_test.go (same package)
// No need to redefine — Go test files in same package share declarations

func TestMongo_CreateAndFindByID(t *testing.T) {
	repo := NewMockRepository()
	u := &User{Name: "Mongo User", Email: "mongo@t.com", IsActive: true}

	err := repo.Create(u)
	assert.NoError(t, err)

	found, err := repo.FindByID(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Mongo User", found.Name)
}

func TestMongo_FindByActive(t *testing.T) {
	repo := NewMockRepository()
	_ = repo.Create(&User{Name: "Active", Email: "active@t.com", IsActive: true})
	_ = repo.Create(&User{Name: "Inactive", Email: "inactive@t.com", IsActive: false})

	result, err := repo.FindByActive(false)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Inactive", result[0].Name)
}

func TestMongo_Count(t *testing.T) {
	repo := NewMockRepository()
	_ = repo.Create(&User{Name: "A", Email: "a@mongo.com"})

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestMongo_Delete_NotFound(t *testing.T) {
	repo := NewMockRepository()
	err := repo.Delete(999)
	assert.EqualError(t, err, "user not found")
}
