package database

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrate runs database migrations
func (d *DB) Migrate(ctx context.Context) error {
	switch d.config.Type {
	case "sqlite":
		return d.migrateSQLite(ctx)
	case "postgres":
		return d.migratePostgres(ctx)
	default:
		return fmt.Errorf("unsupported database type for migration: %s", d.config.Type)
	}
}

// migrateSQLite runs migrations for SQLite
func (d *DB) migrateSQLite(ctx context.Context) error {
	driver, err := sqlite3.WithInstance(d.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite migration driver: %w", err)
	}

	// Create migration instance from file system
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations/sqlite",
		"sqlite3", driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// migratePostgres runs migrations for PostgreSQL
func (d *DB) migratePostgres(ctx context.Context) error {
	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres migration driver: %w", err)
	}

	// Create migration instance from file system
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations/postgres",
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Rollback rolls back the last migration
func (d *DB) Rollback(ctx context.Context) error {
	switch d.config.Type {
	case "sqlite":
		return d.rollbackSQLite(ctx)
	case "postgres":
		return d.rollbackPostgres(ctx)
	default:
		return fmt.Errorf("unsupported database type for rollback: %s", d.config.Type)
	}
}

// rollbackSQLite rolls back the last migration for SQLite
func (d *DB) rollbackSQLite(ctx context.Context) error {
	driver, err := sqlite3.WithInstance(d.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations/sqlite",
		"sqlite3", driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	err = m.Steps(-1)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// rollbackPostgres rolls back the last migration for PostgreSQL
func (d *DB) rollbackPostgres(ctx context.Context) error {
	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations/postgres",
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	err = m.Steps(-1)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}
