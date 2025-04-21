package api

import (
	"context"
	"log"

	pb "auth.service/api/proto"
	"auth.service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserServiceHandler(userService service.UserService) *UserServiceHandler {
	return &UserServiceHandler{
		userService: userService,
	}
}

func (h *UserServiceHandler) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.UserResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"username and password are required",
		)
	}

	userID, err := h.userService.CreateUser(ctx, req.Username, req.Password)
	if err != nil {
		log.Printf("failed to create user: %v", err)
		switch err {
		case service.ErrUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pb.UserResponse{
		UserId:   userID,
		Username: req.Username,
	}, nil
}

func (h *UserServiceHandler) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.UserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"user ID is required",
		)
	}

	user, err := h.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		log.Printf("failed to create user: %v", err)
		switch err {
		case service.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pb.UserResponse{
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

func (h *UserServiceHandler) UpdateUser(
	ctx context.Context,
	req *pb.UpdateUserRequest,
) (*pb.UserResponse, error) {
	if req.UserId == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"user ID is required",
		)
	}

	userID := req.UserId.Value
	var username, password string

	if req.Username != nil {
		username = req.Username.Value
	}
	if req.Password != nil {
		password = req.Password.Value
	}

	if username == "" && password == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"username and password are required",
		)
	}

	err := h.userService.UpdateUser(ctx, userID, username, password)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		switch err {
		case service.ErrUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		case service.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.UserResponse{
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

func (h *UserServiceHandler) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (*pb.UserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"user ID is required",
		)
	}

	user, err := h.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		log.Printf("failed to create user: %v", err)
		switch err {
		case service.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	err = h.userService.DeleteUser(ctx, req.UserId)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.UserResponse{
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}
