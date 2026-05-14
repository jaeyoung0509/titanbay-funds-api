package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jaeyoung0509/titanbay-funds-api/internal/app"
	postgresadapter "github.com/jaeyoung0509/titanbay-funds-api/internal/adapter/postgres"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/platform"
	"github.com/rs/zerolog"
)

type Config struct {
	DatabaseURL   string
	Port          string
	MigrationsDir string
}

func loadConfig() Config {
	cfg := Config{
		DatabaseURL:   getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/titanbay?sslmode=disable"),
		Port:          getenv("PORT", "8080"),
		MigrationsDir: getenv("MIGRATIONS_DIR", "migrations"),
	}
	return cfg
}

func main() {
	flag.Parse()

	cfg := loadConfig()
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("service", "titanbay-funds-api").Logger()

	database, err := platform.Open(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("open database")
	}
	defer database.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	if err := platform.WaitForDB(ctx, database, 30, time.Second); err != nil {
		logger.Fatal().Err(err).Msg("wait for database")
	}

	if err := platform.Migrate(database, cfg.MigrationsDir); err != nil {
		logger.Fatal().Err(err).Msg("run migrations")
	}

	server := app.New(app.Dependencies{
		Repo: postgresadapter.New(database),
		Logger: &logger,
	})

	logger.Info().Str("port", cfg.Port).Msg("starting api")
	if err := server.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		logger.Fatal().Err(err).Msg("listen")
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
