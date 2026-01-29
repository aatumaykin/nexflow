package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.LogRepository = (*LogRepository)(nil)

// LogRepository implements repository.LogRepository using SQLC-generated queries
type LogRepository struct {
	db *sql.DB
}

// NewLogRepository creates a new LogRepository instance
func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{db: db}
}

// Create saves a new log entry
func (r *LogRepository) Create(ctx context.Context, log *entity.Log) error {
	dbLog := mappers.LogToDB(log)
	if dbLog == nil {
		return fmt.Errorf("failed to convert log to db model")
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO logs (id, level, source, message, metadata, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		dbLog.ID, dbLog.Level, dbLog.Source, dbLog.Message, dbLog.Metadata, dbLog.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	return nil
}

// FindByID retrieves a log entry by ID
func (r *LogRepository) FindByID(ctx context.Context, id string) (*entity.Log, error) {
	var dbLog dbmodel.Log

	err := r.db.QueryRowContext(ctx,
		`SELECT id, level, source, message, metadata, created_at FROM logs WHERE id = ? LIMIT 1`,
		id,
	).Scan(&dbLog.ID, &dbLog.Level, &dbLog.Source, &dbLog.Message, &dbLog.Metadata, &dbLog.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("log not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find log by id: %w", err)
	}

	return mappers.LogToDomain(&dbLog), nil
}

// FindByLevel retrieves logs by level
func (r *LogRepository) FindByLevel(ctx context.Context, level string, limit int) ([]*entity.Log, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, level, source, message, metadata, created_at FROM logs WHERE level = ? ORDER BY created_at DESC LIMIT ?`,
		level, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find logs by level: %w", err)
	}
	defer rows.Close()

	var dbLogs []dbmodel.Log
	for rows.Next() {
		var dbLog dbmodel.Log

		if err := rows.Scan(&dbLog.ID, &dbLog.Level, &dbLog.Source, &dbLog.Message, &dbLog.Metadata, &dbLog.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}

		dbLogs = append(dbLogs, dbLog)
	}

	return mappers.LogsToDomain(dbLogs), nil
}

// FindBySource retrieves logs by source
func (r *LogRepository) FindBySource(ctx context.Context, source string, limit int) ([]*entity.Log, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, level, source, message, metadata, created_at FROM logs WHERE source = ? ORDER BY created_at DESC LIMIT ?`,
		source, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find logs by source: %w", err)
	}
	defer rows.Close()

	var dbLogs []dbmodel.Log
	for rows.Next() {
		var dbLog dbmodel.Log

		if err := rows.Scan(&dbLog.ID, &dbLog.Level, &dbLog.Source, &dbLog.Message, &dbLog.Metadata, &dbLog.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}

		dbLogs = append(dbLogs, dbLog)
	}

	return mappers.LogsToDomain(dbLogs), nil
}

// FindByDateRange retrieves logs within a date range
func (r *LogRepository) FindByDateRange(ctx context.Context, startDate, endDate string, limit int) ([]*entity.Log, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, level, source, message, metadata, created_at FROM logs WHERE created_at >= ? AND created_at <= ? ORDER BY created_at DESC LIMIT ?`,
		startDate, endDate, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find logs by date range: %w", err)
	}
	defer rows.Close()

	var dbLogs []dbmodel.Log
	for rows.Next() {
		var dbLog dbmodel.Log

		if err := rows.Scan(&dbLog.ID, &dbLog.Level, &dbLog.Source, &dbLog.Message, &dbLog.Metadata, &dbLog.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}

		dbLogs = append(dbLogs, dbLog)
	}

	return mappers.LogsToDomain(dbLogs), nil
}

// Delete removes a log entry
func (r *LogRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM logs WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("log not found: %s", id)
	}

	return nil
}

// DeleteOlderThan removes logs older than a specific date
func (r *LogRepository) DeleteOlderThan(ctx context.Context, date string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM logs WHERE created_at < ?`, date)
	if err != nil {
		return fmt.Errorf("failed to delete logs older than: %w", err)
	}

	return nil
}

// CountByLevel counts logs by level
func (r *LogRepository) CountByLevel(ctx context.Context, level string) (int, error) {
	var count int

	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM logs WHERE level = ?`,
		level,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count logs by level: %w", err)
	}

	return count, nil
}
