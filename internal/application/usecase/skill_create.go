package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// CreateSkill creates a new skill
func (uc *SkillUseCase) CreateSkill(ctx context.Context, req dto.CreateSkillRequest) (*dto.SkillResponse, error) {
	// Create new skill using entity constructor
	newSkill := entity.NewSkill(req.Name, req.Version, req.Location, req.Permissions, req.Metadata)

	if err := uc.skillRepo.Create(ctx, newSkill); err != nil {
		return &dto.SkillResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create skill: %v", err),
		}, fmt.Errorf("failed to create skill: %w", err)
	}

	uc.logger.Info("skill created", "skill_id", newSkill.ID, "name", newSkill.Name)

	return &dto.SkillResponse{
		Success: true,
		Skill:   dto.SkillDTOFromEntity(newSkill),
	}, nil
}
