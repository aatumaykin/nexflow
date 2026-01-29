package usecase

import (
	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// SkillUseCase handles skill-related business logic
type SkillUseCase struct {
	skillRepo    repository.SkillRepository
	skillRuntime ports.SkillRuntime
	logger       logging.Logger
}

// NewSkillUseCase creates a new SkillUseCase
func NewSkillUseCase(
	skillRepo repository.SkillRepository,
	skillRuntime ports.SkillRuntime,
	logger logging.Logger,
) *SkillUseCase {
	return &SkillUseCase{
		skillRepo:    skillRepo,
		skillRuntime: skillRuntime,
		logger:       logger,
	}
}
