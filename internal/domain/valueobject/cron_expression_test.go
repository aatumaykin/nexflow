package valueobject

import (
	"encoding/json"
	"testing"
)

func TestNewCronExpression(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    CronExpression
		wantErr bool
	}{
		{
			name:    "valid every hour",
			expr:    "0 * * * *",
			want:    CronExpression("0 * * * *"),
			wantErr: false,
		},
		{
			name:    "valid every 5 minutes",
			expr:    "*/5 * * * *",
			want:    CronExpression("*/5 * * * *"),
			wantErr: false,
		},
		{
			name:    "valid every day at midnight",
			expr:    "0 0 * * *",
			want:    CronExpression("0 0 * * *"),
			wantErr: false,
		},
		{
			name:    "valid every Monday",
			expr:    "0 0 * * 1",
			want:    CronExpression("0 0 * * 1"),
			wantErr: false,
		},
		{
			name:    "valid specific time",
			expr:    "30 14 * * *",
			want:    CronExpression("30 14 * * *"),
			wantErr: false,
		},
		{
			name:    "empty expression",
			expr:    "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid - too few parts",
			expr:    "0 * * *",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid - too many parts",
			expr:    "0 * * * * *",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid - invalid minute",
			expr:    "60 * * * *",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid - invalid hour",
			expr:    "* 24 * * *",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCronExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCronExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewCronExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustNewCronExpression(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		wantPanic bool
	}{
		{
			name:      "valid expression",
			expr:      "0 * * * *",
			wantPanic: false,
		},
		{
			name:      "invalid expression",
			expr:      "invalid",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("MustNewCronExpression() did not panic")
					}
				}()
				MustNewCronExpression(tt.expr)
			} else {
				MustNewCronExpression(tt.expr)
			}
		})
	}
}

func TestCronExpression_String(t *testing.T) {
	tests := []struct {
		name string
		c    CronExpression
		want string
	}{
		{
			name: "every hour",
			c:    CronExpression("0 * * * *"),
			want: "0 * * * *",
		},
		{
			name: "every 5 minutes",
			c:    CronExpression("*/5 * * * *"),
			want: "*/5 * * * *",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("CronExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronExpression_IsEmpty(t *testing.T) {
	if CronExpression("0 * * * *").IsEmpty() {
		t.Error("CronExpression(0 * * * *).IsEmpty() returned true")
	}
	if !CronExpression("").IsEmpty() {
		t.Error("CronExpression(\"\").IsEmpty() returned false")
	}
}

func TestCronExpression_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    CronExpression
		want bool
	}{
		{
			name: "valid every hour",
			c:    CronExpression("0 * * * *"),
			want: true,
		},
		{
			name: "valid every 5 minutes",
			c:    CronExpression("*/5 * * * *"),
			want: true,
		},
		{
			name: "invalid - too few parts",
			c:    CronExpression("0 * * *"),
			want: false,
		},
		{
			name: "empty",
			c:    CronExpression(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("CronExpression.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCronExpression_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		c    CronExpression
		want string
	}{
		{
			name: "every hour",
			c:    CronExpression("0 * * * *"),
			want: `"0 * * * *"`,
		},
		{
			name: "every 5 minutes",
			c:    CronExpression("*/5 * * * *"),
			want: `"*/5 * * * *"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.c)
			if err != nil {
				t.Errorf("CronExpression.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("CronExpression.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestCronExpression_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    CronExpression
		wantErr bool
	}{
		{
			name:    "valid every hour",
			data:    `"0 * * * *"`,
			want:    CronExpression("0 * * * *"),
			wantErr: false,
		},
		{
			name:    "valid every 5 minutes",
			data:    `"*/5 * * * *"`,
			want:    CronExpression("*/5 * * * *"),
			wantErr: false,
		},
		{
			name:    "empty",
			data:    `""`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			data:    `"invalid"`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got CronExpression
			err := json.Unmarshal([]byte(tt.data), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("CronExpression.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CronExpression.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
