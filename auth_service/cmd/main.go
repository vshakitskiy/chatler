package main

import (
	"context"
	"log"
	"os"

	a "auth.service/internal/app"
	"auth.service/internal/repository/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	ctx := context.Background()

	sqdsn := os.Getenv("SQLITE_PATH")
	db, err := sqlx.Connect("sqlite3", sqdsn)
	if err != nil {
		panic(err)
	}

	userRepo := sqlite.NewUserRepository(db)
	sessionRepo := sqlite.NewSessionRepository(db)

	app := a.NewApp(ctx, userRepo, sessionRepo)
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
