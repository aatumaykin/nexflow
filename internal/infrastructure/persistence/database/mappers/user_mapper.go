package mappers

import (
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// UserToDomain converts SQLC User model to domain User entity.
func UserToDomain(dbUser *dbmodel.User) *entity.User {
	if dbUser == nil {
		return nil
	}

	return &entity.User{
		ID:        valueobject.UserID(dbUser.ID),
		Channel:   valueobject.MustNewChannel(dbUser.Channel),
		ChannelID: dbUser.ChannelUserID,
		CreatedAt: utils.ParseTimeRFC3339(dbUser.CreatedAt),
	}
}

// UserToDB converts domain User entity to SQLC User model.
func UserToDB(user *entity.User) *dbmodel.User {
	if user == nil {
		return nil
	}

	return &dbmodel.User{
		ID:            string(user.ID),
		Channel:       string(user.Channel),
		ChannelUserID: user.ChannelID,
		CreatedAt:     utils.FormatTimeRFC3339(user.CreatedAt),
	}
}

// UsersToDomain converts slice of SQLC User models to domain User entities.
func UsersToDomain(dbUsers []dbmodel.User) []*entity.User {
	users := make([]*entity.User, 0, len(dbUsers))
	for i := range dbUsers {
		users = append(users, UserToDomain(&dbUsers[i]))
	}
	return users
}
