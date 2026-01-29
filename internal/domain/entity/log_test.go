package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLog(t *testing.T) {
	// Arrange & Act
	log := NewLog(LogLevelInfo, "module1", "test message", map[string]interface{}{"key": "value"})

	// Assert
	require.NotEmpty(t, log.ID)
	assert.Equal(t, string(LogLevelInfo), log.Level)
	assert.Equal(t, "module1", log.Source)
	assert.Equal(t, "test message", log.Message)
	assert.Equal(t, `{"key":"value"}`, log.Metadata) // marshaled
	assert.WithinDuration(t, time.Now(), log.CreatedAt, time.Second)
}

func TestNewLog_EmptyMetadata(t *testing.T) {
	// Arrange & Act
	log := NewLog(LogLevelInfo, "module1", "test message", map[string]interface{}{})

	// Assert
	assert.Equal(t, "{}", log.Metadata)
}

func TestLog_IsDebug(t *testing.T) {
	// Arrange
	log := NewLog(LogLevelDebug, "module1", "debug message", nil)

	// Act & Assert
	assert.True(t, log.IsDebug())
	assert.False(t, log.IsInfo())
	assert.False(t, log.IsWarn())
	assert.False(t, log.IsError())
}

func TestLog_IsInfo(t *testing.T) {
	// Arrange
	log := NewLog(LogLevelInfo, "module1", "info message", nil)

	// Act & Assert
	assert.False(t, log.IsDebug())
	assert.True(t, log.IsInfo())
	assert.False(t, log.IsWarn())
	assert.False(t, log.IsError())
}

func TestLog_IsWarn(t *testing.T) {
	// Arrange
	log := NewLog(LogLevelWarn, "module1", "warn message", nil)

	// Act & Assert
	assert.False(t, log.IsDebug())
	assert.False(t, log.IsInfo())
	assert.True(t, log.IsWarn())
	assert.False(t, log.IsError())
}

func TestLog_IsError(t *testing.T) {
	// Arrange
	log := NewLog(LogLevelError, "module1", "error message", nil)

	// Act & Assert
	assert.False(t, log.IsDebug())
	assert.False(t, log.IsInfo())
	assert.False(t, log.IsWarn())
	assert.True(t, log.IsError())
}

func TestLog_IsFromSource(t *testing.T) {
	// Arrange
	log := NewLog(LogLevelInfo, "module1", "message", nil)

	// Act & Assert
	assert.True(t, log.IsFromSource("module1"))
	assert.False(t, log.IsFromSource("module2"))
}

func TestLog_DifferentLevels(t *testing.T) {
	// Arrange
	debugLog := NewLog(LogLevelDebug, "module", "msg", nil)
	infoLog := NewLog(LogLevelInfo, "module", "msg", nil)
	warnLog := NewLog(LogLevelWarn, "module", "msg", nil)
	errorLog := NewLog(LogLevelError, "module", "msg", nil)

	// Act & Assert
	assert.True(t, debugLog.IsDebug())
	assert.True(t, infoLog.IsInfo())
	assert.True(t, warnLog.IsWarn())
	assert.True(t, errorLog.IsError())

	assert.False(t, debugLog.IsError())
	assert.False(t, errorLog.IsDebug())
}

func TestLog_DifferentSources(t *testing.T) {
	// Arrange
	log1 := NewLog(LogLevelInfo, "module1", "msg", nil)
	log2 := NewLog(LogLevelInfo, "module2", "msg", nil)
	log3 := NewLog(LogLevelInfo, "module1", "msg", nil)

	// Act & Assert
	assert.True(t, log1.IsFromSource("module1"))
	assert.True(t, log2.IsFromSource("module2"))
	assert.True(t, log3.IsFromSource("module1"))

	assert.False(t, log1.IsFromSource("module2"))
	assert.False(t, log2.IsFromSource("module1"))
}

func TestLog_UniqueIDs(t *testing.T) {
	// Arrange
	log1 := NewLog(LogLevelInfo, "module", "msg", nil)
	log2 := NewLog(LogLevelInfo, "module", "msg", nil)

	// Act & Assert
	assert.NotEqual(t, log1.ID, log2.ID)
}

func TestLog_MultipleLogs(t *testing.T) {
	// Arrange
	logs := []*Log{
		NewLog(LogLevelDebug, "module1", "debug msg", nil),
		NewLog(LogLevelInfo, "module1", "info msg", nil),
		NewLog(LogLevelWarn, "module2", "warn msg", nil),
		NewLog(LogLevelError, "module2", "error msg", nil),
	}

	// Act & Assert
	assert.Equal(t, 4, len(logs))
	assert.True(t, logs[0].IsDebug())
	assert.True(t, logs[1].IsInfo())
	assert.True(t, logs[2].IsWarn())
	assert.True(t, logs[3].IsError())

	assert.True(t, logs[0].IsFromSource("module1"))
	assert.True(t, logs[1].IsFromSource("module1"))
	assert.True(t, logs[2].IsFromSource("module2"))
	assert.True(t, logs[3].IsFromSource("module2"))
}
