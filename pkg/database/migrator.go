package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	log.Println("--- GOOSE: STARTING MIGRATIONS ---")

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	log.Println("--- GOOSE: MIGRATIONS FINISHED ---")
	return nil
}
