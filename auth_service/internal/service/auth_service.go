package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"auth.service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthServiceImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtSecret   []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	jwtSecret []byte,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		// TODO: DEFAULT SECRET
		jwtSecret: []byte(getEnv("JWT_SECRET_KEY", "asdasdasdasdasdas")),
		accessTTL: parseDuration(getEnv(
			"ACCESS_TOKEN_TTL",
			"15m",
		)),
		refreshTTL: parseDuration(getEnv(
			"REFRESH_TOKEN_TTL",
			"24h",
		)),
	}
}

func (s *AuthServiceImpl) Login(
	ctx context.Context,
	username, password string,
) (*TokenPair, error) {
	op := "AuthService.Login"

	user, err := s.userRepo.UserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if user == nil {
		return nil, ErrInvalidToken
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	u := &User{
		ID:       user.ID,
		Username: user.Username,
	}

	return s.createTokens(ctx, u)
}

func (s *AuthServiceImpl) createTokens(
	ctx context.Context,
	user *User,
) (*TokenPair, error) {
	op := "AuthService.createTokens"

	now := time.Now()

	accessToken, err := s.generateAccessToken(user, now)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken := uuid.New().String()

	session := &repository.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(s.refreshTTL),
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
	}, nil
}

func (s *AuthServiceImpl) generateAccessToken(
	user *User,
	now time.Time,
) (string, error) {
	claims := CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "auth.service",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("Unable to create token")
	}

	return signedToken, nil
}

func (s *AuthServiceImpl) RefreshTokens(
	ctx context.Context,
	refreshToken string,
) (*TokenPair, error) {
	op := "AuthService.RefreshTokens"

	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.sessionRepo.DeleteSession(ctx, session.ID)
	}

	user, err := s.userRepo.UserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_ = s.sessionRepo.DeleteSession(ctx, session.ID)

	u := &User{
		ID:       user.ID,
		Username: user.Username,
	}

	return s.createTokens(ctx, u)
}

func (s *AuthServiceImpl) ValidateToken(
	ctx context.Context,
	accessToken string,
) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}

			return s.jwtSecret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse token %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if time.Now().Unix() > claims.ExpiresAt.Time.Unix() {
			return nil, ErrExpiredToken
		}

		return &TokenClaims{
			UserID:    claims.UserID,
			Username:  claims.Username,
			ExpiresAt: claims.ExpiresAt.Time,
		}, nil
	}

	return nil, ErrInvalidToken
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		if hours, err := strconv.Atoi(value); err != nil {
			return time.Duration(hours) * time.Hour
		}
		return 15 * time.Minute
	}

	return duration
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
