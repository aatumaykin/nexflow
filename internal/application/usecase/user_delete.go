package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// DeleteUser deletes a user by ID
func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return handleUserError(err, "user not found")
	}

	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return handleUserError(err, "failed to delete user")
	}

	uc.logger.Info("user deleted", "user_id", user.ID, "channel", user.Channel)

	return dto.SuccessUserResponse(dto.UserDTOFromEntity(user)), nil
}
