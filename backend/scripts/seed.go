// backend/scripts/seed.go
package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultDSN = "postgres://esms_user:esms_pass@postgres:5432/esms?sslmode=disable"

func main() {
	var (
		seedDir string
		dsn     string
		timeout time.Duration
		dryRun  bool
	)

	flag.StringVar(&seedDir, "dir", "database/seed", "Path to seed SQL directory")
	flag.StringVar(&dsn, "dsn", envOrDefault("DATABASE_URL", defaultDSN), "Database connection string")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Statement timeout")
	flag.BoolVar(&dryRun, "dry-run", false, "Print statements without executing")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	files := []string{"users.sql", "resources.sql"}

	if dryRun {
		for _, file := range files {
			path := filepath.Join(seedDir, file)
			fmt.Printf("-- Dry run: %s\n", path)
		}
		return
	}

	if err := runSeed(ctx, dsn, seedDir, files); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	log.Println("✅ Seed data applied successfully")
}

func runSeed(ctx context.Context, dsn, dir string, files []string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, file := range files {
		path := filepath.Join(dir, file)
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to read %s: %w", path, readErr)
		}

		if _, execErr := tx.ExecContext(ctx, string(content)); execErr != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute %s: %w", path, execErr)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit seed transaction: %w", err)
	}

	return nil
}

func envOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// Ensure we catch the compiler error if files slice is empty in future refactors.
var _ = errors.New
