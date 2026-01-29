package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// ExecuteSkill executes a skill with given input parameters
func (uc *SkillUseCase) ExecuteSkill(ctx context.Context, req dto.SkillExecutionRequest) (*dto.SkillExecutionResponse, error) {
	// Validate skill exists
	_, err := uc.skillRepo.FindByName(ctx, req.Skill)
	if err != nil {
		return &dto.SkillExecutionResponse{
			Success: false,
			Error:   fmt.Sprintf("skill not found: %v", err),
		}, fmt.Errorf("skill not found: %w", err)
	}

	// Execute skill through runtime
	result, err := uc.skillRuntime.Execute(ctx, req.Skill, req.Input)
	if err != nil {
		return &dto.SkillExecutionResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to execute skill: %v", err),
		}, fmt.Errorf("failed to execute skill: %w", err)
	}

	uc.logger.Info("skill executed", "skill", req.Skill, "success", result.Success)

	return &dto.SkillExecutionResponse{
		Success: result.Success,
		Output:  result.Output,
		Error:   result.Error,
	}, nil
}

// ValidateSkill validates a skill
func (uc *SkillUseCase) ValidateSkill(ctx context.Context, skillName string) error {
	return uc.skillRuntime.Validate(skillName)
}

// ListAvailableSkills returns list of available skill names
func (uc *SkillUseCase) ListAvailableSkills(ctx context.Context) ([]string, error) {
	return uc.skillRuntime.List()
}

// GetSkillDetails returns detailed skill information
func (uc *SkillUseCase) GetSkillDetails(ctx context.Context, skillName string) (map[string]interface{}, error) {
	return uc.skillRuntime.GetSkill(skillName)
}
