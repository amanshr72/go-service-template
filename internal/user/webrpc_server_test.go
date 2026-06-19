package user

import (
	"context"
	webrpcuser "go-crud2/internal/user/webrpc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupWebRPCServer() *WebRPCServer {
	repo := NewMockRepository()
	notifier := &MockNotifier{}

	svc := NewService(repo, notifier)

	return NewWebRPCServer(svc)
}

func TestWebRPC_CreateUser(t *testing.T) {
	s := setupWebRPCServer()

	u, err := s.CreateUser(context.Background(), &webrpcuser.CreateUserInput{
		Name:  "Aman",
		Email: "aman@t.com",
	})

	assert.NoError(t, err)
	assert.Equal(t, "Aman", u.Name)
	assert.NotZero(t, u.Id)
}

func TestWebRPC_GetUser_NotFound(t *testing.T) {
	s := setupWebRPCServer()

	_, err := s.GetUser(context.Background(), 999)
	assert.Error(t, err)
}

func TestWebRPC_ListUsers(t *testing.T) {
	s := setupWebRPCServer()
	_, _ = s.CreateUser(context.Background(), &webrpcuser.CreateUserInput{Name: "A", Email: "a@t.com"})
	_, _ = s.CreateUser(context.Background(), &webrpcuser.CreateUserInput{Name: "B", Email: "b@t.com"})

	users, err := s.ListUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestWebRPC_DeleteUser(t *testing.T) {
	s := setupWebRPCServer()
	u, _ := s.CreateUser(context.Background(), &webrpcuser.CreateUserInput{Name: "Del", Email: "d@t.com"})

	ok, err := s.DeleteUser(context.Background(), u.Id)
	assert.NoError(t, err)
	assert.True(t, ok)

	// confirm actually gone
	_, err = s.GetUser(context.Background(), u.Id)
	assert.Error(t, err)
}

func TestWebRPC_GetCount(t *testing.T) {
	s := setupWebRPCServer()
	_, _ = s.CreateUser(context.Background(), &webrpcuser.CreateUserInput{Name: "A", Email: "a@t.com"})

	count, err := s.GetCount(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int32(1), count)
}

func TestWebRPC_GetActiveUsers(t *testing.T) {
	s := setupWebRPCServer()
	_, _ = s.CreateUser(context.Background(), &webrpcuser.CreateUserInput{Name: "A", Email: "a@t.com"})

	users, err := s.GetActiveUsers(context.Background(), true)
	assert.NoError(t, err)
	assert.Len(t, users, 1) // newly created users are active by default
}
