package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.LogRepository = (*LogRepository)(nil)

type LogRepository struct {
	queries *database.Queries
}

func NewLogRepository(queries *database.Queries) *LogRepository {
	return &LogRepository{queries: queries}
}

func (r *LogRepository) Create(ctx context.Context, log *entity.Log) error {
	dbLog := mappers.LogToDB(log)
	if dbLog == nil {
		return fmt.Errorf("failed to convert log to db model")
	}

	var metadata sql.NullString
	if log.Metadata != "" {
		metadata.Valid = true
		metadata.String = log.Metadata
	}

	_, err := r.queries.CreateLog(ctx, database.CreateLogParams{
		ID:        dbLog.ID,
		Level:     dbLog.Level,
		Source:    dbLog.Source,
		Message:   dbLog.Message,
		Metadata:  metadata,
		CreatedAt: dbLog.CreatedAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create log: %w", err)
	}

	return nil
}

func (r *LogRepository) FindByID(ctx context.Context, id string) (*entity.Log, error) {
	dbLog, err := r.queries.GetLogByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("log not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find log by id: %w", err)
	}

	return mappers.LogToDomain(&dbLog), nil
}

func (r *LogRepository) FindByLevel(ctx context.Context, level string, limit int) ([]*entity.Log, error) {
	dbLogs, err := r.queries.GetLogsByLevel(ctx, database.GetLogsByLevelParams{
		Level: level,
		Limit: int64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find logs by level: %w", err)
	}

	return mappers.LogsToDomain(dbLogs), nil
}

func (r *LogRepository) FindBySource(ctx context.Context, source string, limit int) ([]*entity.Log, error) {
	dbLogs, err := r.queries.GetLogsBySource(ctx, database.GetLogsBySourceParams{
		Source: source,
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find logs by source: %w", err)
	}

	return mappers.LogsToDomain(dbLogs), nil
}

func (r *LogRepository) FindByDateRange(ctx context.Context, startDate, endDate string, limit int) ([]*entity.Log, error) {
	dbLogs, err := r.queries.GetLogsByDateRange(ctx, database.GetLogsByDateRangeParams{
		CreatedAt:   startDate,
		CreatedAt_2: endDate,
		Limit:       int64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find logs by date range: %w", err)
	}

	return mappers.LogsToDomain(dbLogs), nil
}

func (r *LogRepository) Delete(ctx context.Context, id string) error {
	_, err := r.queries.GetLogByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("log not found: %s", id)
		}
		return fmt.Errorf("failed to check log existence: %w", err)
	}

	err = r.queries.DeleteLog(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}

	return nil
}

func (r *LogRepository) DeleteOlderThan(ctx context.Context, date string) error {
	err := r.queries.DeleteLogsOlderThan(ctx, date)
	if err != nil {
		return fmt.Errorf("failed to delete logs older than: %w", err)
	}

	return nil
}

func (r *LogRepository) CountByLevel(ctx context.Context, level string) (int, error) {
	dbLogs, err := r.queries.GetLogsByLevel(ctx, database.GetLogsByLevelParams{
		Level: level,
		Limit: 1000000,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count logs by level: %w", err)
	}

	return len(dbLogs), nil
}
