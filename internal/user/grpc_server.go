package user

import (
	"context"
	"go-crud2/internal/user/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer implements the generated UserServiceServer interface.
// Embeds UnimplementedUserServiceServer for forward-compat (future proto methods won't break build).
type GRPCServer struct {
	pb.UnimplementedUserServiceServer
	svc Service // same port REST/GraphQL use
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func toProtoUser(u *User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:        int32(u.ID),
		Name:      u.Name,
		Email:     u.Email,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
	}
}

func (s *GRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	u, err := s.svc.Create(CreateUserInput{Name: req.Name, Email: req.Email})
	if err != nil {
		// gRPC has its own status codes — InvalidArgument maps to your validation errors
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return toProtoUser(u), nil
}

func (s *GRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	u, err := s.svc.GetByID(int(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoUser(u), nil
}

func (s *GRPCServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := s.svc.GetAll()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &pb.ListUsersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, toProtoUser(&u))
	}
	return resp, nil
}

func (s *GRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	u, err := s.svc.Update(int(req.Id), UpdateUserInput{Name: req.Name, Email: req.Email})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoUser(u), nil
}

func (s *GRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.svc.Delete(int(req.Id)); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}
