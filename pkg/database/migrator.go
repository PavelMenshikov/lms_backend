package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	slog.Info("goose: starting migrations")

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	slog.Info("goose: migrations finished")
	return nil
}
