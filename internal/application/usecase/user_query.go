package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return handleUserError(err, "failed to find user")
	}

	return dto.SuccessUserResponse(dto.UserDTOFromEntity(user)), nil
}

// GetUserByChannel retrieves a user by channel and channel ID
func (uc *UserUseCase) GetUserByChannel(ctx context.Context, channel, channelID string) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByChannel(ctx, channel, channelID)
	if err != nil {
		return handleUserError(err, "failed to find user by channel")
	}

	return dto.SuccessUserResponse(dto.UserDTOFromEntity(user)), nil
}

// ListUsers retrieves all users
func (uc *UserUseCase) ListUsers(ctx context.Context) (*dto.UsersResponse, error) {
	users, err := uc.userRepo.List(ctx)
	if err != nil {
		return dto.ErrorUsersResponse(err), err
	}

	userDTOs := make([]*dto.UserDTO, 0, len(users))
	for _, user := range users {
		userDTOs = append(userDTOs, dto.UserDTOFromEntity(user))
	}

	return dto.SuccessUsersResponse(userDTOs), nil
}
