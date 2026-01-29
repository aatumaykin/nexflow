package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return dto.ErrorUserResponse(fmt.Errorf("failed to find user: %w", err)), fmt.Errorf("failed to find user: %w", err)
	}

	return dto.SuccessUserResponse(dto.UserDTOFromEntity(user)), nil
}

// GetUserByChannel retrieves a user by channel and channel ID
func (uc *UserUseCase) GetUserByChannel(ctx context.Context, channel, channelID string) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByChannel(ctx, channel, channelID)
	if err != nil {
		return dto.ErrorUserResponse(fmt.Errorf("failed to find user by channel: %w", err)), fmt.Errorf("failed to find user by channel: %w", err)
	}

	return dto.SuccessUserResponse(dto.UserDTOFromEntity(user)), nil
}

// ListUsers retrieves all users
func (uc *UserUseCase) ListUsers(ctx context.Context) (*dto.UsersResponse, error) {
	users, err := uc.userRepo.List(ctx)
	if err != nil {
		return dto.ErrorUsersResponse(fmt.Errorf("failed to list users: %w", err)), fmt.Errorf("failed to list users: %w", err)
	}

	userDTOs := make([]*dto.UserDTO, 0, len(users))
	for _, user := range users {
		userDTOs = append(userDTOs, dto.UserDTOFromEntity(user))
	}

	return dto.SuccessUsersResponse(userDTOs), nil
}
