package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/atumaikin/nexflow/internal/config"
	"github.com/atumaikin/nexflow/internal/logging"
)

// DB is main database implementation
type DB struct {
	*Queries
	db     *sql.DB
	config *DBConfig
	logger logging.Logger
}

// NewDatabase creates a new database instance.
// By default, it uses a NoopLogger that does nothing. Use WithLogger option to provide a custom logger.
func NewDatabase(cfg *config.DatabaseConfig, opts ...Option) (Database, error) {
	dbConfig := &DBConfig{
		Type:            cfg.Type,
		Path:            cfg.Path,
		MigrationsPath:  cfg.MigrationsPath,
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}

	// Validate configuration
	if err := dbConfig.Validate(); err != nil {
		return nil, err
	}

	// Set default connection pool settings if not provided
	if dbConfig.MaxOpenConns == 0 {
		dbConfig.MaxOpenConns = 25
	}
	if dbConfig.MaxIdleConns == 0 {
		dbConfig.MaxIdleConns = 25
	}
	if dbConfig.ConnMaxLifetime == 0 {
		dbConfig.ConnMaxLifetime = 5 * time.Minute
	}

	var db *sql.DB
	var err error

	switch dbConfig.Type {
	case "sqlite":
		db, err = openSQLite(dbConfig.Path)
	case "postgres":
		db, err = openPostgres(dbConfig.Path)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings from config
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	queries := New(db)

	// Create database instance with default NoopLogger
	dbInstance := &DB{
		Queries: queries,
		db:      db,
		config:  dbConfig,
		logger:  logging.NewNoopLogger(), // Default to NoopLogger
	}

	// Apply options
	for _, opt := range opts {
		opt(dbInstance)
	}

	return dbInstance, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	d.logger.Info("Closing database connection", "type", d.config.Type)
	if err := d.db.Close(); err != nil {
		d.logger.Error("Failed to close database connection", "error", err)
		return err
	}
	d.logger.Info("Database connection closed successfully")
	return nil
}

// GetDB returns the underlying *sql.DB instance
func (d *DB) GetDB() *sql.DB {
	return d.db
}
