package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

type DB struct {
	Client *sql.DB
}

func NewPostgres(dbUrl string) (*DB, error) {
	c, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to db: %w", err)
	}
	if err := c.Ping(); err != nil {
		return nil, fmt.Errorf("error making ping request to db: %w", err)
	}
	slog.Info("Connected to Postgres Successfully", "DbUrl", dbUrl)
	c.SetMaxOpenConns(100)
	c.SetMaxIdleConns(30)
	c.SetConnMaxLifetime(5 * time.Minute)

	return &DB{
		Client: c,
	}, err
}
