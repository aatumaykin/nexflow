package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// GetSkillByID retrieves a skill by ID
func (uc *SkillUseCase) GetSkillByID(ctx context.Context, id string) (*dto.SkillResponse, error) {
	skill, err := uc.skillRepo.FindByID(ctx, id)
	if err != nil {
		return handleSkillError(err, "failed to find skill")
	}

	return dto.SuccessSkillResponse(dto.SkillDTOFromEntity(skill)), nil
}

// GetSkillByName retrieves a skill by name
func (uc *SkillUseCase) GetSkillByName(ctx context.Context, name string) (*dto.SkillResponse, error) {
	skill, err := uc.skillRepo.FindByName(ctx, name)
	if err != nil {
		return handleSkillError(err, "failed to find skill by name")
	}

	return dto.SuccessSkillResponse(dto.SkillDTOFromEntity(skill)), nil
}

// ListSkills retrieves all skills
func (uc *SkillUseCase) ListSkills(ctx context.Context) (*dto.SkillsResponse, error) {
	skills, err := uc.skillRepo.List(ctx)
	if err != nil {
		return dto.ErrorSkillsResponse(err), err
	}

	skillDTOs := make([]*dto.SkillDTO, 0, len(skills))
	for _, skill := range skills {
		skillDTOs = append(skillDTOs, dto.SkillDTOFromEntity(skill))
	}

	return dto.SuccessSkillsResponse(skillDTOs), nil
}
