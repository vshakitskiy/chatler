package service

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found 404")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	ID       string
	Username string
}

type TokenPair struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	UserID    string
	Username  string
	ExpiresAt time.Time
}

type UserService interface {
	CreateUser(ctx context.Context, username, password string) (string, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	UpdateUser(ctx context.Context, userID, username, password string) error
	DeleteUser(ctx context.Context, userID string) error
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (*TokenPair, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error)
	ValidateToken(ctx context.Context, accessToken string) (*TokenClaims, error)
}

type AcccessService interface {
	Check(ctx context.Context, accessToken string) (bool, string, error)
}
