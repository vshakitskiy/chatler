package service

import (
	"context"
	"errors"
	"fmt"

	"auth.service/internal/repository"
	"auth.service/pkg"
)

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{userRepo: userRepo}
}

func (s *UserServiceImpl) CreateUser(
	ctx context.Context,
	username, password string,
) (string, error) {
	op := "UserService.CreateUser"

	existingUser, err := s.userRepo.UserByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			break
		default:
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}
	if existingUser != nil {
		return "", ErrUserAlreadyExists
	}

	hashedPassword, err := pkg.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &repository.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return user.ID, nil
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID string) (*User, error) {
	op := "UserService.UserByID"

	user, err := s.userRepo.UserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &User{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *UserServiceImpl) UpdateUser(
	ctx context.Context,
	userID, username, password string,
) error {
	op := "UserService.UpdateUser"

	user, err := s.userRepo.UserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	isChanged := false

	if username != "" && user.Username != username {
		existingUser, err := s.userRepo.UserByUsername(ctx, username)

		if err != nil {
			switch err {
			case repository.ErrUserNotFound:
				break
			default:
				return fmt.Errorf("%s: %w", op, err)
			}
		}
		if existingUser != nil {
			return ErrUserAlreadyExists
		}

		user.Username = username
		isChanged = true
	}

	if password != "" {
		hashedPassword, err := pkg.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hashedPassword
		isChanged = true
	}

	if isChanged {
		if err := s.userRepo.UpdateUser(ctx, user); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID string) error {
	op := "UserService.DeleteUser"

	if err := s.userRepo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
