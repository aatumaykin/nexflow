package repository

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// LogRepository defines the interface for log data operations
type LogRepository interface {
	// Create saves a new log entry
	Create(ctx context.Context, log *entity.Log) error

	// FindByID retrieves a log entry by ID
	FindByID(ctx context.Context, id string) (*entity.Log, error)

	// FindByLevel retrieves logs by level
	FindByLevel(ctx context.Context, level string, limit int) ([]*entity.Log, error)

	// FindBySource retrieves logs by source
	FindBySource(ctx context.Context, source string, limit int) ([]*entity.Log, error)

	// FindByDateRange retrieves logs within a date range
	FindByDateRange(ctx context.Context, startDate, endDate string, limit int) ([]*entity.Log, error)

	// Delete removes a log entry
	Delete(ctx context.Context, id string) error

	// DeleteOlderThan removes logs older than a specific date
	DeleteOlderThan(ctx context.Context, date string) error

	// CountByLevel counts logs by level
	CountByLevel(ctx context.Context, level string) (int, error)
}
