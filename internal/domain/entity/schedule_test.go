package entity

import (
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSchedule(t *testing.T) {
	// Arrange & Act
	schedule := NewSchedule("my-skill", "0 * * * *", `{"param": "value"}`)

	// Assert
	require.NotEmpty(t, schedule.ID)
	assert.Equal(t, "my-skill", schedule.Skill)
	assert.Equal(t, valueobject.CronExpression("0 * * * *"), schedule.CronExpression)
	assert.Equal(t, `{"param": "value"}`, schedule.Input)
	assert.True(t, schedule.Enabled)
	assert.WithinDuration(t, time.Now(), schedule.CreatedAt, time.Second)
}

func TestSchedule_Enable(t *testing.T) {
	// Arrange
	schedule := NewSchedule("skill", "0 * * * *", "{}")
	schedule.Disable()
	assert.False(t, schedule.IsEnabled())

	// Act
	schedule.Enable()

	// Assert
	assert.True(t, schedule.IsEnabled())
}

func TestSchedule_Disable(t *testing.T) {
	// Arrange
	schedule := NewSchedule("skill", "0 * * * *", "{}")
	assert.True(t, schedule.IsEnabled())

	// Act
	schedule.Disable()

	// Assert
	assert.False(t, schedule.IsEnabled())
}

func TestSchedule_IsEnabled(t *testing.T) {
	// Arrange
	schedule := NewSchedule("skill", "0 * * * *", "{}")

	// Act & Assert
	assert.True(t, schedule.IsEnabled())

	schedule.Disable()
	assert.False(t, schedule.IsEnabled())

	schedule.Enable()
	assert.True(t, schedule.IsEnabled())
}

func TestSchedule_BelongsToSkill(t *testing.T) {
	// Arrange
	schedule := NewSchedule("my-skill", "0 * * * *", "{}")

	// Act & Assert
	assert.True(t, schedule.BelongsToSkill("my-skill"))
	assert.False(t, schedule.BelongsToSkill("other-skill"))
}

func TestSchedule_CronExpressions(t *testing.T) {
	// Arrange
	schedules := []*Schedule{
		NewSchedule("skill1", "0 * * * *", "{}"),   // Every hour
		NewSchedule("skill2", "0 0 * * *", "{}"),   // Every day at midnight
		NewSchedule("skill3", "*/5 * * * *", "{}"), // Every 5 minutes
		NewSchedule("skill4", "0 9 * * 1-5", "{}"), // Weekdays at 9 AM
	}

	// Act & Assert
	assert.Equal(t, valueobject.CronExpression("0 * * * *"), schedules[0].CronExpression)
	assert.Equal(t, valueobject.CronExpression("0 0 * * *"), schedules[1].CronExpression)
	assert.Equal(t, valueobject.CronExpression("*/5 * * * *"), schedules[2].CronExpression)
	assert.Equal(t, valueobject.CronExpression("0 9 * * 1-5"), schedules[3].CronExpression)
}

func TestSchedule_MultipleSchedulesForSameSkill(t *testing.T) {
	// Arrange
	skill := "my-skill"
	schedule1 := NewSchedule(skill, "0 * * * *", "{}")
	schedule2 := NewSchedule(skill, "0 0 * * *", "{}")

	// Act & Assert
	assert.NotEqual(t, schedule1.ID, schedule2.ID)
	assert.True(t, schedule1.BelongsToSkill(skill))
	assert.True(t, schedule2.BelongsToSkill(skill))
	assert.NotEqual(t, schedule1.CronExpression, schedule2.CronExpression)
}

func TestSchedule_EnableDisableCycle(t *testing.T) {
	// Arrange
	schedule := NewSchedule("skill", "0 * * * *", "{}")

	// Act & Assert
	assert.True(t, schedule.IsEnabled())
	schedule.Disable()
	assert.False(t, schedule.IsEnabled())
	schedule.Enable()
	assert.True(t, schedule.IsEnabled())
	schedule.Disable()
	assert.False(t, schedule.IsEnabled())
	schedule.Enable()
	assert.True(t, schedule.IsEnabled())
}

func TestSchedule_InputJSON(t *testing.T) {
	// Arrange
	input1 := `{"key": "value", "number": 123}`
	input2 := `[]`
	input3 := `null`

	schedule1 := NewSchedule("skill", "0 * * * *", input1)
	schedule2 := NewSchedule("skill", "0 * * * *", input2)
	schedule3 := NewSchedule("skill", "0 * * * *", input3)

	// Act & Assert
	assert.Equal(t, input1, schedule1.Input)
	assert.Equal(t, input2, schedule2.Input)
	assert.Equal(t, input3, schedule3.Input)
}

func TestSchedule_UniqueIDs(t *testing.T) {
	// Arrange
	schedule1 := NewSchedule("skill", "0 * * * *", "{}")
	schedule2 := NewSchedule("skill", "0 * * * *", "{}")

	// Act & Assert
	assert.NotEqual(t, schedule1.ID, schedule2.ID)
}
