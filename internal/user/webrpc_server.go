package user

import (
	"context"
	webrpcgen "go-crud2/internal/user/webrpc"
)

type WebRPCServer struct {
	svc Service
}

func NewWebRPCServer(svc Service) *WebRPCServer {
	return &WebRPCServer{svc: svc}
}

func toWebRPCUser(u *User) *webrpcgen.User {
	return &webrpcgen.User{
		Id:       uint32(u.ID),
		Name:     u.Name,
		Email:    u.Email,
		IsActive: u.IsActive,
	}
}

func (s *WebRPCServer) CreateUser(ctx context.Context, input *webrpcgen.CreateUserInput) (*webrpcgen.User, error) {
	u, err := s.svc.Create(CreateUserInput{Name: input.Name, Email: input.Email})
	if err != nil {
		return nil, err
	}
	return toWebRPCUser(u), nil
}

func (s *WebRPCServer) GetUser(ctx context.Context, id uint32) (*webrpcgen.User, error) {
	u, err := s.svc.GetByID(int(id))
	if err != nil {
		return nil, err
	}
	return toWebRPCUser(u), nil
}

func (s *WebRPCServer) ListUsers(ctx context.Context) ([]*webrpcgen.User, error) {
	users, err := s.svc.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]*webrpcgen.User, 0, len(users))
	for _, u := range users {
		u := u
		result = append(result, toWebRPCUser(&u))
	}
	return result, nil
}

func (s *WebRPCServer) UpdateUser(ctx context.Context, id uint32, input *webrpcgen.UpdateUserInput) (*webrpcgen.User, error) {
	u, err := s.svc.Update(int(id), UpdateUserInput{Name: input.Name, Email: input.Email})
	if err != nil {
		return nil, err
	}
	return toWebRPCUser(u), nil
}

func (s *WebRPCServer) DeleteUser(ctx context.Context, id uint32) (bool, error) {
	if err := s.svc.Delete(int(id)); err != nil {
		return false, err
	}
	return true, nil
}

func (s *WebRPCServer) GetCount(ctx context.Context) (int32, error) {
	count, err := s.svc.GetCount()
	return int32(count), err
}

func (s *WebRPCServer) GetActiveUsers(ctx context.Context, active bool) ([]*webrpcgen.User, error) {
	users, err := s.svc.GetByActive(active)
	if err != nil {
		return nil, err
	}
	result := make([]*webrpcgen.User, 0, len(users))
	for _, u := range users {
		u := u
		result = append(result, toWebRPCUser(&u))
	}
	return result, nil
}
