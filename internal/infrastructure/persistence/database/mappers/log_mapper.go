package mappers

import (
	"database/sql"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// LogToDomain converts SQLC Log model to domain Log entity.
func LogToDomain(dbLog *dbmodel.Log) *entity.Log {
	if dbLog == nil {
		return nil
	}

	metadata := ""
	if dbLog.Metadata.Valid {
		metadata = dbLog.Metadata.String
	}

	return &entity.Log{
		ID:        valueobject.LogID(dbLog.ID),
		Level:     valueobject.MustNewLogLevel(dbLog.Level),
		Source:    dbLog.Source,
		Message:   dbLog.Message,
		Metadata:  metadata,
		CreatedAt: utils.ParseTimeRFC3339(dbLog.CreatedAt),
	}
}

// LogToDB converts domain Log entity to SQLC Log model.
func LogToDB(log *entity.Log) *dbmodel.Log {
	if log == nil {
		return nil
	}

	var metadata sql.NullString
	if log.Metadata != "" {
		metadata.Valid = true
		metadata.String = log.Metadata
	}

	return &dbmodel.Log{
		ID:        string(log.ID),
		Level:     string(log.Level),
		Source:    log.Source,
		Message:   log.Message,
		Metadata:  metadata,
		CreatedAt: utils.FormatTimeRFC3339(log.CreatedAt),
	}
}

// LogsToDomain converts slice of SQLC Log models to domain Log entities.
func LogsToDomain(dbLogs []dbmodel.Log) []*entity.Log {
	logs := make([]*entity.Log, 0, len(dbLogs))
	for i := range dbLogs {
		logs = append(logs, LogToDomain(&dbLogs[i]))
	}
	return logs
}
