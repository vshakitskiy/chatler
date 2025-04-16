package sqlite

import (
	"context"
	"fmt"
	"time"

	"auth.service/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Context = context.Context

type SqliteUserRepository struct {
	db *sqlx.DB
}

type SqliteSessionRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *SqliteUserRepository {
	return &SqliteUserRepository{db: db}
}

func NewSessionRepository(db *sqlx.DB) *SqliteSessionRepository {
	return &SqliteSessionRepository{db: db}
}

func (r *SqliteUserRepository) CreateUser(ctx Context, user *repository.User) error {
	op := "repository.UserRepository.CreateUser"

	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, username, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *SqliteUserRepository) UserByID(
	ctx Context,
	id string,
) (*repository.User, error) {
	op := "repository.UserRepository.UserByID"
	user := new(repository.User)

	query := `
		SELECT id, username, password_hash, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *SqliteUserRepository) UserByUsername(
	ctx Context,
	username string,
) (*repository.User, error) {
	op := "repository.UserRepository.UserByUsername"
	user := new(repository.User)

	query := `
		SELECT id, username, password_hash, created_at, updated_at
		FROM users
		WHERE username = ?
	`

	err := r.db.GetContext(ctx, user, query, username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *SqliteUserRepository) UpdateUser(
	ctx Context,
	user *repository.User,
) error {
	op := "repository.UserRepository.UpdateUser"

	user.UpdatedAt = time.Now()
	query := `
		UPDATE users
		SET username = ?, password_hash = ?, updated_at = ?
		WHERE id = ?
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %s", op, "user not found")
	}

	return nil
}

func (r *SqliteUserRepository) DeleteUser(
	ctx Context,
	id string,
) error {
	op := "repository.UserRepository.DeleteUser"

	query := `
		DELETE FROM users
		WHERE id = ?
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %s", op, "user not found")
	}

	return nil
}

func (r *SqliteSessionRepository) CreateSession(
	ctx Context,
	session *repository.Session,
) error {
	op := "repository.SessionRepository.CreateSession"

	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	now := time.Now()
	session.CreatedAt = now
	session.ExpiresAt = now.Add(24 * time.Hour)

	query := `
		INSERT INTO sessions (id, user_id, refresh_token, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.ExpiresAt,
		session.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *SqliteSessionRepository) GetByRefreshToken(
	ctx Context,
	refreshToken string,
) (*repository.Session, error) {
	op := "repository.SessionRepository.SessionByID"
	session := new(repository.Session)

	query := `
		SELECT id, user_id, refresh_token, expires_at, created_at
		FROM sessions
		WHERE refresh_token = ?
	`
	err := r.db.GetContext(ctx, session, query, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (r *SqliteSessionRepository) DeleteSession(
	ctx Context,
	id string,
) error {
	op := "repository.SessionRepository.DeleteSession"

	query := `
		DELETE FROM sessions
		WHERE id = ?
	`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %s", op, "session not found")
	}
	return nil
}

func (r *SqliteSessionRepository) DeleteByUserID(
	ctx Context,
	userID string,
) error {
	op := "repository.SessionRepository.DeleteByUserID"
	query := `
		DELETE FROM sessions
		WHERE user_id = ?
	`

	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %s", op, "sessions not found")
	}

	return nil
}
