package database

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/atumaikin/nexflow/internal/logging"
)

func TestMigrations(t *testing.T) {
	// Create a temporary database file
	tmpDir, err := os.MkdirTemp("", "nexflow-migrate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	// Create DB instance
	dbConfig := &DBConfig{
		Type:           "sqlite",
		Path:           dbPath,
		MigrationsPath: "../../migrations",
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	queries := New(db)

	testDB := &DB{
		Queries: queries,
		db:      db,
		config:  dbConfig,
		logger:  logging.NewNoopLogger(), // use NoopLogger for tests
	}

	// Run migrations
	ctx := context.Background()
	if err := testDB.Migrate(ctx); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Verify tables exist before DB is closed by migrations
	tables := []string{"users", "sessions", "messages", "tasks", "skills", "schedules", "logs"}
	for _, table := range tables {
		var exists int
		err := db.QueryRow(`
			SELECT COUNT(*) FROM sqlite_master
			WHERE type='table' AND name=?
		`, table).Scan(&exists)
		if err != nil {
			t.Fatalf("failed to check table %s: %v", table, err)
		}
		if exists != 1 {
			t.Errorf("table %s does not exist after migration", table)
		}
	}
}
