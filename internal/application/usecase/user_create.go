package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.FindByChannel(ctx, req.Channel, req.ChannelID)
	if err == nil && existingUser != nil {
		return dto.ErrorUserResponse(fmt.Errorf("user already exists: channel=%s, channelID=%s", req.Channel, req.ChannelID)), fmt.Errorf("user already exists")
	}

	// Create new user
	user := entity.NewUser(req.Channel, req.ChannelID)
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return dto.ErrorUserResponse(fmt.Errorf("failed to create user: %w", err)), fmt.Errorf("failed to create user: %w", err)
	}

	uc.logger.Info("user created", "user_id", user.ID, "channel", user.Channel)

	return dto.SuccessUserResponse(dto.UserDTOFromEntity(user)), nil
}

// GetOrCreateUser gets an existing user or creates a new one
func (uc *UserUseCase) GetOrCreateUser(ctx context.Context, channel, channelID string) (*dto.UserResponse, error) {
	// Try to find existing user
	user, err := uc.userRepo.FindByChannel(ctx, channel, channelID)
	if err == nil && user != nil {
		return dto.SuccessUserResponse(dto.UserDTOFromEntity(user)), nil
	}

	// Create new user
	newUser := entity.NewUser(channel, channelID)
	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return dto.ErrorUserResponse(fmt.Errorf("failed to create user: %w", err)), fmt.Errorf("failed to create user: %w", err)
	}

	uc.logger.Info("user created", "user_id", newUser.ID, "channel", newUser.Channel)

	return dto.SuccessUserResponse(dto.UserDTOFromEntity(newUser)), nil
}
