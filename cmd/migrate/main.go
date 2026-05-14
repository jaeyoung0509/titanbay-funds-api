package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/platform"
)

func main() {
	databaseURL := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/titanbay?sslmode=disable")
	migrationsDir := getenv("MIGRATIONS_DIR", "migrations")

	db, err := platform.Open(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	if err := platform.WaitForDB(ctx, db, 30, time.Second); err != nil {
		log.Fatal(err)
	}

	if err := platform.Migrate(db, migrationsDir); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
