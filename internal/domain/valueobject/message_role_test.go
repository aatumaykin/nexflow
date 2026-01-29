package valueobject

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageRole_String(t *testing.T) {
	tests := []struct {
		name     string
		role     MessageRole
		expected string
	}{
		{"user", RoleUser, "user"},
		{"assistant", RoleAssistant, "assistant"},
		{"system", RoleSystem, "system"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.String())
		})
	}
}

func TestMessageRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     MessageRole
		expected bool
	}{
		{"user", RoleUser, true},
		{"assistant", RoleAssistant, true},
		{"system", RoleSystem, true},
		{"invalid", MessageRole("invalid"), false},
		{"empty", MessageRole(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.IsValid())
		})
	}
}

func TestMessageRole_IsUser(t *testing.T) {
	tests := []struct {
		name     string
		role     MessageRole
		expected bool
	}{
		{"user", RoleUser, true},
		{"assistant", RoleAssistant, false},
		{"system", RoleSystem, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.IsUser())
		})
	}
}

func TestMessageRole_IsAssistant(t *testing.T) {
	tests := []struct {
		name     string
		role     MessageRole
		expected bool
	}{
		{"user", RoleUser, false},
		{"assistant", RoleAssistant, true},
		{"system", RoleSystem, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.IsAssistant())
		})
	}
}

func TestMessageRole_IsSystem(t *testing.T) {
	tests := []struct {
		name     string
		role     MessageRole
		expected bool
	}{
		{"user", RoleUser, false},
		{"assistant", RoleAssistant, false},
		{"system", RoleSystem, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.role.IsSystem())
		})
	}
}

func TestMessageRole_MarshalJSON(t *testing.T) {
	role := RoleUser
	data, err := json.Marshal(role)

	require.NoError(t, err)
	assert.JSONEq(t, `"user"`, string(data))
}

func TestMessageRole_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		expected  MessageRole
		expectErr bool
	}{
		{"user", `"user"`, RoleUser, false},
		{"assistant", `"assistant"`, RoleAssistant, false},
		{"system", `"system"`, RoleSystem, false},
		{"invalid", `"invalid"`, MessageRole(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var role MessageRole
			err := json.Unmarshal([]byte(tt.data), &role)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, role)
			}
		})
	}
}

func TestNewMessageRole(t *testing.T) {
	tests := []struct {
		name      string
		roleStr   string
		expected  MessageRole
		expectErr bool
	}{
		{"user", "user", RoleUser, false},
		{"assistant", "assistant", RoleAssistant, false},
		{"system", "system", RoleSystem, false},
		{"invalid", "invalid", MessageRole(""), true},
		{"empty", "", MessageRole(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := NewMessageRole(tt.roleStr)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, role)
			}
		})
	}
}

func TestMustNewMessageRole(t *testing.T) {
	assert.NotPanics(t, func() {
		MustNewMessageRole("user")
	})

	assert.Panics(t, func() {
		MustNewMessageRole("invalid")
	})
}
