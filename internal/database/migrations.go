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
	d.logger.Info("Running SQLite migrations", "path", d.config.MigrationsPath)

	driver, err := sqlite3.WithInstance(d.db, &sqlite3.Config{})
	if err != nil {
		d.logger.Error("Failed to create sqlite migration driver", "error", err)
		return fmt.Errorf("failed to create sqlite migration driver: %w", err)
	}

	// Create migration instance from file system
	migrationPath := d.config.MigrationsPath + "/sqlite"
	d.logger.Debug("Migration path", "path", migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"sqlite3", driver,
	)
	if err != nil {
		d.logger.Error("Failed to create migration instance", "error", err)
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to run migrations", "error", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	d.logger.Info("SQLite migrations completed successfully")
	return nil
}

// migratePostgres runs migrations for PostgreSQL
func (d *DB) migratePostgres(ctx context.Context) error {
	d.logger.Info("Running PostgreSQL migrations", "path", d.config.MigrationsPath)

	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		d.logger.Error("Failed to create postgres migration driver", "error", err)
		return fmt.Errorf("failed to create postgres migration driver: %w", err)
	}

	// Create migration instance from file system
	migrationPath := d.config.MigrationsPath + "/postgres"
	d.logger.Debug("Migration path", "path", migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres", driver,
	)
	if err != nil {
		d.logger.Error("Failed to create migration instance", "error", err)
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to run migrations", "error", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	d.logger.Info("PostgreSQL migrations completed successfully")
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
	d.logger.Info("Rolling back SQLite migrations", "path", d.config.MigrationsPath)

	driver, err := sqlite3.WithInstance(d.db, &sqlite3.Config{})
	if err != nil {
		d.logger.Error("Failed to create sqlite migration driver", "error", err)
		return fmt.Errorf("failed to create sqlite migration driver: %w", err)
	}

	migrationPath := d.config.MigrationsPath + "/sqlite"
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"sqlite3", driver,
	)
	if err != nil {
		d.logger.Error("Failed to create migration instance", "error", err)
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	err = m.Steps(-1)
	if err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to rollback migration", "error", err)
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	d.logger.Info("SQLite migration rollback completed successfully")
	return nil
}

// rollbackPostgres rolls back the last migration for PostgreSQL
func (d *DB) rollbackPostgres(ctx context.Context) error {
	d.logger.Info("Rolling back PostgreSQL migrations", "path", d.config.MigrationsPath)

	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		d.logger.Error("Failed to create postgres migration driver", "error", err)
		return fmt.Errorf("failed to create postgres migration driver: %w", err)
	}

	migrationPath := d.config.MigrationsPath + "/postgres"
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres", driver,
	)
	if err != nil {
		d.logger.Error("Failed to create migration instance", "error", err)
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	err = m.Steps(-1)
	if err != nil && err != migrate.ErrNoChange {
		d.logger.Error("Failed to rollback migration", "error", err)
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	d.logger.Info("PostgreSQL migration rollback completed successfully")
	return nil
}
