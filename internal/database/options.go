package database

import (
	"github.com/atumaikin/nexflow/internal/logging"
)

// Option is a function that configures database.
type Option func(*DB)

// WithLogger sets logger for database.
func WithLogger(logger logging.Logger) Option {
	return func(db *DB) {
		db.logger = logger
	}
}
