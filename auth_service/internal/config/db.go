package config

import (
	"github.com/jmoiron/sqlx"
)

func ConnectSqlite() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", Env.SqlitePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
