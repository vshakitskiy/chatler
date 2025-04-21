package service

import (
	"context"
)

type AccessServiceImpl struct {
	authService AuthService
}

func NewAccessService(authService AuthService) *AccessServiceImpl {
	return &AccessServiceImpl{
		authService: authService,
	}
}

func (s *AccessServiceImpl) Check(
	ctx context.Context,
	accessToken string,
) (bool, string, error) {
	claims, err := s.authService.ValidateToken(ctx, accessToken)
	if err != nil {
		if err == ErrExpiredToken {
			return false, "", ErrExpiredToken
		}
		return false, "", ErrInvalidToken
	}

	return true, claims.UserID, nil
}
