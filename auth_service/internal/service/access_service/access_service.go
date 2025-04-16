package access_service

import (
	"context"
	"errors"

	"auth.service/internal/service"
)

var (
	ErrInvalidToken       = errors.New("Invalid token")
	ErrExpiredToken       = errors.New("Expired token")
	ErrTokenNotFound      = errors.New("Token not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AccessServiceImpl struct {
	authService service.AuthService
}

func NewAccessService(authService service.AuthService) *AccessServiceImpl {
	return &AccessServiceImpl{
		authService: authService,
	}
}

func (s *AccessServiceImpl) Check(
	ctx context.Context,
	accessToken string,
) (bool, string, error) {
	// op := "AccessService.CheckAccess"

	claims, err := s.authService.ValidateToken(ctx, accessToken)
	if err != nil {
		if err == ErrExpiredToken {
			return false, "", ErrExpiredToken
		}
		return false, "", ErrInvalidToken
	}

	return true, claims.UserID, nil
}
