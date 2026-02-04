package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB) error {
	log.Println("--- STARTING MIGRATIONS ---")

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		filename VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT NOW()
	);`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	files, err := ioutil.ReadDir("./migrations")
	if err != nil {
		files, err = ioutil.ReadDir("../../migrations")
		if err != nil {
			return fmt.Errorf("failed to read migrations directory: %w", err)
		}
	}

	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, filename := range sqlFiles {
		if isApplied(db, filename) {
			continue
		}

		log.Printf("Applying migration: %s...", filename)

		content, err := readFile(filename)
		if err != nil {
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing migration %s: %w", filename, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (filename) VALUES ($1)", filename); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		log.Printf("Migration %s applied successfully.", filename)
	}

	log.Println("--- MIGRATIONS FINISHED ---")
	return nil
}

func isApplied(db *sql.DB, filename string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)", filename).Scan(&exists)
	if err != nil {
		log.Printf("Error checking migration status: %v", err)
		return false
	}
	return exists
}

func readFile(filename string) ([]byte, error) {
	path := filepath.Join("migrations", filename)
	content, err := os.ReadFile(path)
	if err != nil {
		path = filepath.Join("../../migrations", filename)
		content, err = os.ReadFile(path)
	}
	return content, err
}
