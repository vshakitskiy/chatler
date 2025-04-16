package repository

import (
	"context"
	"time"
)

type Context = context.Context

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
	CreateUser(ctx Context, user *User) error
	UserByID(ctx Context, id string) (*User, error)
	UserByUsername(ctx Context, username string) (*User, error)
	UpdateUser(ctx Context, user *User) error
	DeleteUser(ctx Context, id string) error
}

type SessionRepository interface {
	CreateSession(ctx Context, session *Session) error
	GetByRefreshToken(ctx Context, refreshToken string) (*Session, error)
	DeleteSession(ctx Context, id string) error
	DeleteByUserID(ctx Context, userID string) error
}
