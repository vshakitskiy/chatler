package sqlite

import (
	"context"
	"fmt"
	"time"

	"auth.service/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SqliteSessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SqliteSessionRepository {
	return &SqliteSessionRepository{db: db}
}

func (r *SqliteSessionRepository) CreateSession(
	ctx context.Context,
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
	ctx context.Context,
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
	ctx context.Context,
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
	ctx context.Context,
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
