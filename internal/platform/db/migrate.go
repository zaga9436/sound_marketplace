package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func ApplyMigrations(db *sql.DB, migrationsDir string) error {
	log.Printf("migration runner: using migrations dir=%s", migrationsDir)

	if err := ensureSchemaMigrationsTable(db); err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("scan migrations: %w", err)
	}
	log.Printf("migration runner: found %d migration file(s)", len(files))
	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", migrationsDir)
	}
	sort.Strings(files)

	for _, file := range files {
		log.Printf("migration runner: checking %s", filepath.Base(file))
		applied, err := isMigrationApplied(db, filepath.Base(file))
		if err != nil {
			return err
		}
		if applied {
			log.Printf("migration runner: skipping already applied %s", filepath.Base(file))
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", file, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (name) VALUES ($1)`, filepath.Base(file)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", file, err)
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		log.Printf("migration runner: applied %s", filepath.Base(file))
	}

	return nil
}

func ensureSchemaMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func isMigrationApplied(db *sql.DB, name string) (bool, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name = $1)`, name).Scan(&exists)
	return exists, err
}
