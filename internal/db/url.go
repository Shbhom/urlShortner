package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *DB) GetUrlByCode(short_code string) (string, error) {
	var url string
	if err := db.Client.QueryRow(`SELECT url from shornted_url where shortnedkey = $1`, short_code).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("No record found for short key %s", short_code)
		}
		return "", fmt.Errorf("Error while fetching url from shortcode: %w", err)
	}
	return url, nil
}

func (db *DB) AddUrl(url, short_code string) error {
	if _, err := db.Client.Exec(`INSERT into shornted_url (shortnedkey,url) VALUES ($1,$2)`, short_code, url); err != nil {
		return fmt.Errorf("Erorr while inserting short Url to db: %w", err)
	}
	return nil
}
