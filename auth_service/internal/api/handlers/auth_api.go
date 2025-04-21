package handlers

import (
	"context"

	pb "auth.service/api/proto"
	"auth.service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthServiceHandler(authService service.AuthService) *AuthServiceHandler {
	return &AuthServiceHandler{
		authService: authService,
	}
}

func (h *AuthServiceHandler) GetAccessToken(
	ctx context.Context,
	req *pb.RefreshTokenRequest,
) (*pb.AccessTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"refresh token is required",
		)
	}

	tokenPair, err := h.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(
			codes.Internal, "internal server error",
		)
	}

	return &pb.AccessTokenResponse{
		AccessToken: tokenPair.AccessToken,
	}, nil
}

func (h *AuthServiceHandler) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"username and password are required",
		)
	}

	tokenPair, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			return nil, status.Error(
				codes.Unauthenticated,
				"invalid credentials",
			)
		default:
			return nil, status.Error(
				codes.Internal, "internal server error",
			)
		}
	}

	return &pb.LoginResponse{
		UserId:       tokenPair.UserID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
