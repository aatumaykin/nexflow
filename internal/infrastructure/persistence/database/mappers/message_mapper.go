package mappers

import (
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// MessageToDomain converts SQLC Message model to domain Message entity.
func MessageToDomain(dbMessage *dbmodel.Message) *entity.Message {
	if dbMessage == nil {
		return nil
	}

	return &entity.Message{
		ID:        valueobject.MessageID(dbMessage.ID),
		SessionID: valueobject.MustNewSessionID(dbMessage.SessionID),
		Role:      valueobject.MustNewMessageRole(dbMessage.Role),
		Content:   dbMessage.Content,
		CreatedAt: utils.ParseTimeRFC3339(dbMessage.CreatedAt),
	}
}

// MessageToDB converts domain Message entity to SQLC Message model.
func MessageToDB(message *entity.Message) *dbmodel.Message {
	if message == nil {
		return nil
	}

	return &dbmodel.Message{
		ID:        string(message.ID),
		SessionID: string(message.SessionID),
		Role:      string(message.Role),
		Content:   message.Content,
		CreatedAt: utils.FormatTimeRFC3339(message.CreatedAt),
	}
}

// MessagesToDomain converts slice of SQLC Message models to domain Message entities.
func MessagesToDomain(dbMessages []dbmodel.Message) []*entity.Message {
	messages := make([]*entity.Message, 0, len(dbMessages))
	for i := range dbMessages {
		messages = append(messages, MessageToDomain(&dbMessages[i]))
	}
	return messages
}
