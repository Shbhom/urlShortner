package main

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/lib/pq"
	"github.com/shbhom/urlShortner/internal/config"
	"github.com/shbhom/urlShortner/internal/db"
)

func main() {
	slog.Info("Loading config for seeder...")
	cfg := config.LoadConfig("local")

	slog.Info("Connecting to Postgres...")
	// IMPORTANT NOTE: db.NewPostgres returns a struct that contains `Client *sql.DB`. 
	// You MUST extract the Client to use it.
	database := db.NewPostgres(cfg.DB_URL)
	defer database.Client.Close()

	slog.Info("Truncating table shornted_url...")
	_, err := database.Client.Exec("TRUNCATE TABLE shornted_url;")
	if err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}

	slog.Info("Starting COPY IN protocol for 1,000,000 rows...")
	start := time.Now()

	txn, err := database.Client.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	// Prepare COPY statement for shortnedKey and url
	stmt, err := txn.Prepare(pq.CopyIn("shornted_url", "shortnedkey", "url"))
	if err != nil {
		log.Fatalf("Failed to prepare COPY IN: %v", err)
	}

	for i := 0; i < 1000000; i++ {
		code := fmt.Sprintf("seed_%06d", i)
		targetURL := fmt.Sprintf("https://example.com/target/%06d", i)
		_, err = stmt.Exec(code, targetURL)
		if err != nil {
			log.Fatalf("Failed to execute copy inside loop at index %d: %v", i, err)
		}
	}

	slog.Info("Flushing COPY IN data...")
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("Failed to execute final flush: %v", err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatalf("Failed to close statement: %v", err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	slog.Info("Seeding complete!", "duration", time.Since(start).String(), "rows", 1000000)
}
