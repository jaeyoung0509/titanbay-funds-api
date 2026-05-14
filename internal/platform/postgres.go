package platform

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

func WaitForDB(ctx context.Context, db *sql.DB, attempts int, delay time.Duration) error {
	if attempts < 1 {
		attempts = 1
	}
	if delay <= 0 {
		delay = time.Second
	}

	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := db.PingContext(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			if lastErr != nil {
				return lastErr
			}
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	if lastErr == nil {
		return errors.New("database not ready")
	}
	return lastErr
}

func Migrate(database *sql.DB, migrationsDir string) error {
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(database, migrationsDir)
}

func SeedFromFile(ctx context.Context, database *sql.DB, seedPath string) error {
	if seedPath == "" {
		return nil
	}

	content, err := os.ReadFile(seedPath)
	if err != nil {
		return err
	}

	statements := splitSQLStatements(string(content))
	if len(statements) == 0 {
		return nil
	}

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func splitSQLStatements(raw string) []string {
	lines := strings.Split(raw, "\n")
	var builder strings.Builder
	var statements []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		builder.WriteString(line)
		builder.WriteByte('\n')
		if strings.HasSuffix(trimmed, ";") {
			statement := strings.TrimSpace(builder.String())
			statement = strings.TrimSuffix(statement, ";")
			statement = strings.TrimSpace(statement)
			if statement != "" {
				statements = append(statements, statement)
			}
			builder.Reset()
		}
	}

	remaining := strings.TrimSpace(builder.String())
	if remaining != "" {
		statements = append(statements, strings.TrimSuffix(remaining, ";"))
	}
	return statements
}
