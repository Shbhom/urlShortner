package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	Client *sql.DB
}

func NewPostgres(dbUrl string) *DB {
	c, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error connecting to db: ", err)
	}
	if err := c.Ping(); err != nil {
		log.Fatal("error making ping request to db: ", err)
	}
	c.SetMaxOpenConns(100)
	c.SetMaxIdleConns(30)
	c.SetConnMaxLifetime(5 * time.Minute)
	
	return &DB{
		Client: c,
	}
}
