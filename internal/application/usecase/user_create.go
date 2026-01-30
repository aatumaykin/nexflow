package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	existingUser, err := uc.userRepo.FindByChannel(ctx, req.Channel, req.ChannelID)
	if err == nil && existingUser != nil {
		return handleUserError(fmt.Errorf("user already exists: channel=%s, channelID=%s", req.Channel, req.ChannelID), "user already exists")
	}

	user := entity.NewUser(req.Channel, req.ChannelID)
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return handleUserError(err, "failed to create user")
	}

	uc.logger.Info("user created", "user_id", user.ID, "channel", user.Channel)

	return dto.SuccessUserResponse(dto.UserDTOFromEntity(user)), nil
}

// GetOrCreateUser gets an existing user or creates a new one
func (uc *UserUseCase) GetOrCreateUser(ctx context.Context, channel, channelID string) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByChannel(ctx, channel, channelID)
	if err == nil && user != nil {
		return dto.SuccessUserResponse(dto.UserDTOFromEntity(user)), nil
	}

	newUser := entity.NewUser(channel, channelID)
	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return handleUserError(err, "failed to create user")
	}

	uc.logger.Info("user created", "user_id", newUser.ID, "channel", newUser.Channel)

	return dto.SuccessUserResponse(dto.UserDTOFromEntity(newUser)), nil
}
