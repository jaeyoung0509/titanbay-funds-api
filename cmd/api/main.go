package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	postgresadapter "github.com/jaeyoung0509/titanbay-funds-api/internal/adapter/postgres"
	"github.com/jaeyoung0509/titanbay-funds-api/internal/app"
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

	startupCtx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	if err := platform.WaitForDB(startupCtx, database, 30, time.Second); err != nil {
		logger.Fatal().Err(err).Msg("wait for database")
	}

	if err := platform.Migrate(database, cfg.MigrationsDir); err != nil {
		logger.Fatal().Err(err).Msg("run migrations")
	}

	server := app.New(app.Dependencies{
		Repo:   postgresadapter.New(database),
		Logger: &logger,
	})

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	listenErr := make(chan error, 1)
	logger.Info().Str("port", cfg.Port).Msg("starting api")
	go func() {
		listenErr <- server.Listen(fmt.Sprintf(":%s", cfg.Port))
	}()

	select {
	case err := <-listenErr:
		if err != nil {
			logger.Fatal().Err(err).Msg("listen")
		}
		logger.Info().Msg("api stopped")
		return
	case <-signalCtx.Done():
		logger.Info().Msg("shutdown signal received")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Fatal().Err(err).Msg("shutdown api")
	}

	if err := <-listenErr; err != nil {
		logger.Fatal().Err(err).Msg("listen stopped")
	}

	logger.Info().Msg("api stopped")
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
