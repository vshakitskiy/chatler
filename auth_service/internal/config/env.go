package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	SqlitePath      string
	GRPCPort        string
	JWTSecret       string
	AccessTokenTTL  string
	RefreshTokenTTL string
}

var Env *env

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}

	sqdsn := os.Getenv("SQLITE_PATH")
	if sqdsn == "" {
		return errors.New("SQLITE_PATH is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return errors.New("JWT_SECRET_KEY is not set")
	}

	port := getEnv("GRPC_PORT", "50051")

	accessTokenTTL := getEnv("ACCESS_TOKEN_TTL", "15m")
	refreshTokenTTL := getEnv("REFRESH_TOKEN_TTL", "24h")

	env := &env{
		SqlitePath:      sqdsn,
		GRPCPort:        port,
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
	}

	Env = env
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	log.Printf("Using default %s: %s\n", key, defaultValue)
	return defaultValue
}
