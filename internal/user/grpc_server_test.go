package user

import (
	"context"
	"go-crud2/internal/user/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPC_CreateUser(t *testing.T) {
	svc := NewService(NewMockRepository(), &MockNotifier{})
	s := NewGRPCServer(svc)

	resp, err := s.CreateUser(context.Background(), &pb.CreateUserRequest{Name: "Aman", Email: "a@t.com"})
	assert.NoError(t, err)
	assert.Equal(t, "Aman", resp.Name)
}

func TestGRPC_GetUser_NotFound(t *testing.T) {
	svc := NewService(NewMockRepository(), &MockNotifier{})
	s := NewGRPCServer(svc)

	_, err := s.GetUser(context.Background(), &pb.GetUserRequest{Id: 999})
	assert.Error(t, err)

	// gRPC errors carry status codes — verify the right one came back
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGRPC_ListUsers(t *testing.T) {
	svc := NewService(NewMockRepository(), &MockNotifier{})
	_, _ = svc.Create(CreateUserInput{Name: "A", Email: "a@t.com"})
	s := NewGRPCServer(svc)

	resp, err := s.ListUsers(context.Background(), &pb.ListUsersRequest{})
	assert.NoError(t, err)
	assert.Len(t, resp.Users, 1)
}

func TestGRPC_DeleteUser(t *testing.T) {
	svc := NewService(NewMockRepository(), &MockNotifier{})
	u, _ := svc.Create(CreateUserInput{Name: "Del", Email: "d@t.com"})
	s := NewGRPCServer(svc)

	resp, err := s.DeleteUser(context.Background(), &pb.DeleteUserRequest{Id: int32(u.ID)})
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}
