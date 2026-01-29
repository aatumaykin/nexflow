package usecase

import (
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
