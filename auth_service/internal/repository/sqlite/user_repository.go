package sqlite

import (
	"context"
	"fmt"
	"time"

	"auth.service/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SqliteUserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *SqliteUserRepository {
	return &SqliteUserRepository{db: db}
}

func (r *SqliteUserRepository) CreateUser(ctx context.Context, user *repository.User) error {
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
	ctx context.Context,
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
	ctx context.Context,
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
		switch {
		case err.Error() == "sql: no rows in result set":
			return nil, repository.ErrUserNotFound
		default:
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return user, nil
}

func (r *SqliteUserRepository) UpdateUser(
	ctx context.Context,
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
	ctx context.Context,
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
