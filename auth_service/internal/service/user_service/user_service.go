package user_service

import (
	"context"
	"errors"
	"fmt"

	"auth.service/internal/repository"
	"auth.service/internal/service"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
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
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if existingUser != nil {
		return "", ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(password)
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

func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID string) (*service.User, error) {
	op := "UserService.UserByID"

	user, err := s.userRepo.UserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &service.User{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
func (s *UserServiceImpl) UpdateUser(ctx context.Context, userID, username, password string) error {
	op := "UserService.UpdateUser"

	user, err := s.userRepo.UserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	if username != "" && user.Username != username {
		existingUser, err := s.userRepo.UserByUsername(ctx, username)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if existingUser == nil {
			return ErrUserAlreadyExists
		}

		user.Username = username
	}

	if password != "" {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hashedPassword
	}

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
