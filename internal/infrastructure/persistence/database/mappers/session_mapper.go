package mappers

import (
	"github.com/atumaikin/nexflow/internal/domain/entity"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// SessionToDomain converts SQLC Session model to domain Session entity.
func SessionToDomain(dbSession *dbmodel.Session) *entity.Session {
	if dbSession == nil {
		return nil
	}

	return &entity.Session{
		ID:        dbSession.ID,
		UserID:    dbSession.UserID,
		CreatedAt: utils.ParseTimeRFC3339(dbSession.CreatedAt),
		UpdatedAt: utils.ParseTimeRFC3339(dbSession.UpdatedAt),
	}
}

// SessionToDB converts domain Session entity to SQLC Session model.
func SessionToDB(session *entity.Session) *dbmodel.Session {
	if session == nil {
		return nil
	}

	return &dbmodel.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		CreatedAt: utils.FormatTimeRFC3339(session.CreatedAt),
		UpdatedAt: utils.FormatTimeRFC3339(session.UpdatedAt),
	}
}

// SessionsToDomain converts slice of SQLC Session models to domain Session entities.
func SessionsToDomain(dbSessions []dbmodel.Session) []*entity.Session {
	sessions := make([]*entity.Session, 0, len(dbSessions))
	for i := range dbSessions {
		sessions = append(sessions, SessionToDomain(&dbSessions[i]))
	}
	return sessions
}
