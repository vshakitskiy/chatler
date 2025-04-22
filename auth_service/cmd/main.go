package main

import (
	"context"
	"log"

	a "auth.service/internal/app"
	"auth.service/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.ConnectSqlite()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	app, err := a.NewApp(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
