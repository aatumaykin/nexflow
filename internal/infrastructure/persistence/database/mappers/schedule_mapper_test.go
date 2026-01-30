package mappers

import (
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduleToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbSchedule  *dbmodel.Schedule
		expected    *entity.Schedule
		expectedNil bool
	}{
		{
			name: "Valid schedule with enabled=true",
			dbSchedule: &dbmodel.Schedule{
				ID:             "test-id",
				Skill:          "test-skill",
				CronExpression: "0 0 * * *",
				Input:          "test input",
				Enabled:        1,
				CreatedAt:      time.Now().Format(time.RFC3339),
			},
			expected: &entity.Schedule{
				ID:             valueobject.ScheduleID("test-id"),
				Skill:          "test-skill",
				CronExpression: valueobject.MustNewCronExpression("0 0 * * *"),
				Input:          "test input",
				Enabled:        true,
				CreatedAt:      time.Now(),
			},
			expectedNil: false,
		},
		{
			name: "Valid schedule with enabled=false",
			dbSchedule: &dbmodel.Schedule{
				ID:             "test-id-2",
				Skill:          "test-skill-2",
				CronExpression: "0 1 * * *",
				Input:          "test input 2",
				Enabled:        0,
				CreatedAt:      time.Now().Format(time.RFC3339),
			},
			expected: &entity.Schedule{
				ID:             valueobject.ScheduleID("test-id-2"),
				Skill:          "test-skill-2",
				CronExpression: valueobject.MustNewCronExpression("0 1 * * *"),
				Input:          "test input 2",
				Enabled:        false,
				CreatedAt:      time.Now(),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			dbSchedule:  nil,
			expected:    nil,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScheduleToDomain(tt.dbSchedule)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Skill, result.Skill)
			assert.Equal(t, tt.expected.CronExpression, result.CronExpression)
			assert.Equal(t, tt.expected.Input, result.Input)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.WithinDuration(t, tt.expected.CreatedAt, result.CreatedAt, time.Second)
		})
	}
}

func TestScheduleToDB(t *testing.T) {
	tests := []struct {
		name        string
		schedule    *entity.Schedule
		expected    *dbmodel.Schedule
		expectedNil bool
	}{
		{
			name: "Valid schedule with enabled=true",
			schedule: &entity.Schedule{
				ID:             valueobject.ScheduleID("test-id"),
				Skill:          "test-skill",
				CronExpression: valueobject.MustNewCronExpression("0 0 * * *"),
				Input:          "test input",
				Enabled:        true,
				CreatedAt:      time.Now(),
			},
			expected: &dbmodel.Schedule{
				ID:             "test-id",
				Skill:          "test-skill",
				CronExpression: "0 0 * * *",
				Input:          "test input",
				Enabled:        1,
				CreatedAt:      time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name: "Valid schedule with enabled=false",
			schedule: &entity.Schedule{
				ID:             valueobject.ScheduleID("test-id-2"),
				Skill:          "test-skill-2",
				CronExpression: valueobject.MustNewCronExpression("0 1 * * *"),
				Input:          "test input 2",
				Enabled:        false,
				CreatedAt:      time.Now(),
			},
			expected: &dbmodel.Schedule{
				ID:             "test-id-2",
				Skill:          "test-skill-2",
				CronExpression: "0 1 * * *",
				Input:          "test input 2",
				Enabled:        0,
				CreatedAt:      time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			schedule:    nil,
			expected:    nil,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScheduleToDB(tt.schedule)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Skill, result.Skill)
			assert.Equal(t, tt.expected.CronExpression, result.CronExpression)
			assert.Equal(t, tt.expected.Input, result.Input)
			assert.Equal(t, tt.expected.Enabled, result.Enabled)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
		})
	}
}

func TestSchedulesToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbSchedules []dbmodel.Schedule
		expectedLen int
	}{
		{
			name: "Multiple schedules",
			dbSchedules: []dbmodel.Schedule{
				{
					ID:             "sched-1",
					Skill:          "skill-1",
					CronExpression: "0 0 * * *",
					Input:          "input-1",
					Enabled:        1,
					CreatedAt:      time.Now().Format(time.RFC3339),
				},
				{
					ID:             "sched-2",
					Skill:          "skill-2",
					CronExpression: "0 1 * * *",
					Input:          "input-2",
					Enabled:        0,
					CreatedAt:      time.Now().Format(time.RFC3339),
				},
				{
					ID:             "sched-3",
					Skill:          "skill-3",
					CronExpression: "0 2 * * *",
					Input:          "input-3",
					Enabled:        1,
					CreatedAt:      time.Now().Format(time.RFC3339),
				},
			},
			expectedLen: 3,
		},
		{
			name:        "Empty slice",
			dbSchedules: []dbmodel.Schedule{},
			expectedLen: 0,
		},
		{
			name:        "Nil input",
			dbSchedules: nil,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SchedulesToDomain(tt.dbSchedules)

			if tt.dbSchedules == nil {
				assert.Empty(t, result)
				return
			}

			assert.Len(t, result, tt.expectedLen)
			for i, schedule := range result {
				assert.Equal(t, tt.dbSchedules[i].ID, string(schedule.ID))
				assert.Equal(t, tt.dbSchedules[i].Skill, schedule.Skill)
				assert.Equal(t, tt.dbSchedules[i].Enabled == 1, schedule.Enabled)
			}
		})
	}
}
