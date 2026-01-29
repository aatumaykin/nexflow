package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// UserUseCase handles user-related business logic
type UserUseCase struct {
	userRepo repository.UserRepository
	logger   logging.Logger
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(
	userRepo repository.UserRepository,
	logger logging.Logger,
) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

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

// DeleteUser deletes a user by ID
func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	// Check if user exists
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return dto.ErrorUserResponse(fmt.Errorf("user not found: %w", err)), fmt.Errorf("user not found: %w", err)
	}

	// Delete user
	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return dto.ErrorUserResponse(fmt.Errorf("failed to delete user: %w", err)), fmt.Errorf("failed to delete user: %w", err)
	}

	uc.logger.Info("user deleted", "user_id", user.ID, "channel", user.Channel)

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
