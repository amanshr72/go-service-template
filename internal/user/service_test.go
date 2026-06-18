package user

import (
	"errors"
	"go-crud2/internal/notification"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockNotifier struct {
	SendCalled bool
	ShouldFail bool
}

func (m *MockNotifier) SendEmail(req notification.SendEmailRequest) (*notification.SendEmailResponse, error) {
	m.SendCalled = true
	if m.ShouldFail {
		return nil, errors.New("mock vendor failure")
	}
	return &notification.SendEmailResponse{MessageID: "mock_123", Status: "queued"}, nil
}
func TestService_Create_NotifierFailureDoesNotBlockCreate(t *testing.T) {
	notifier := &MockNotifier{ShouldFail: true}
	svc := NewService(NewMockRepository(), notifier)

	u, err := svc.Create(CreateUserInput{Name: "Aman", Email: "a@t.com"})

	assert.NoError(t, err)
	assert.NotZero(t, u.ID)
	assert.True(t, notifier.SendCalled)
}

func TestService_GetByActive(t *testing.T) {
	svc := NewService(NewMockRepository(), &MockNotifier{})
	_, _ = svc.Create(CreateUserInput{Name: "A", Email: "a@t.com"}) // IsActive=true by default

	result, err := svc.GetByActive(true)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestService_GetCount(t *testing.T) {
	svc := NewService(NewMockRepository(), &MockNotifier{})
	_, _ = svc.Create(CreateUserInput{Name: "A", Email: "a@t.com"})
	_, _ = svc.Create(CreateUserInput{Name: "B", Email: "b@t.com"})

	count, err := svc.GetCount()
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestService_Update_IsActive(t *testing.T) {
	svc := NewService(NewMockRepository(), &MockNotifier{})
	u, _ := svc.Create(CreateUserInput{Name: "A", Email: "a@t.com"})

	f := false
	updated, err := svc.Update(u.ID, UpdateUserInput{IsActive: &f})
	assert.NoError(t, err)
	assert.False(t, updated.IsActive)
}
