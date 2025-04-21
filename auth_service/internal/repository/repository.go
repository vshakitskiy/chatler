package repository

import (
	"context"
	"time"
)

type User struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Session struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	UserByID(ctx context.Context, id string) (*User, error)
	UserByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (*Session, error)
	DeleteSession(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
