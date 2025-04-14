package repository

import (
	"context"

	m "auth.service/internal/models"
)

type Context = context.Context

type UserRepository interface {
	CreateUser(ctx Context, user *m.User) error
	UserByID(ctx Context, id string) (*m.User, error)
	UserByUsername(ctx Context, username string) (*m.User, error)
	UpdateUser(ctx Context, user *m.User) error
	DeleteUser(ctx Context, id string) error
}

type SessionRepository interface {
	CreateSession(ctx Context, session *m.Session) error
	GetByRefreshToken(ctx Context, refreshToken string) (*m.Session, error)
	DeleteSession(ctx Context, id string) error
	DeleteByUserID(ctx Context, userID string) error
}
