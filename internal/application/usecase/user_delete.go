package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

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
