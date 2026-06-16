package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_GetByActive(t *testing.T) {
	svc := NewService(NewMockRepository())
	_, _ = svc.Create(CreateUserInput{Name: "A", Email: "a@t.com"}) // IsActive=true by default

	result, err := svc.GetByActive(true)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestService_GetCount(t *testing.T) {
	svc := NewService(NewMockRepository())
	_, _ = svc.Create(CreateUserInput{Name: "A", Email: "a@t.com"})
	_, _ = svc.Create(CreateUserInput{Name: "B", Email: "b@t.com"})

	count, err := svc.GetCount()
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestService_Update_IsActive(t *testing.T) {
	svc := NewService(NewMockRepository())
	u, _ := svc.Create(CreateUserInput{Name: "A", Email: "a@t.com"})

	f := false
	updated, err := svc.Update(u.ID, UpdateUserInput{IsActive: &f})
	assert.NoError(t, err)
	assert.False(t, updated.IsActive)
}
