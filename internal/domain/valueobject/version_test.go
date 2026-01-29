package valueobject

import (
	"encoding/json"
	"testing"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    Version
		wantErr bool
	}{
		{
			name:    "valid version 1.0.0",
			version: "1.0.0",
			want:    Version("1.0.0"),
			wantErr: false,
		},
		{
			name:    "valid version 2.3.1",
			version: "2.3.1",
			want:    Version("2.3.1"),
			wantErr: false,
		},
		{
			name:    "valid version 10.20.30",
			version: "10.20.30",
			want:    Version("10.20.30"),
			wantErr: false,
		},
		{
			name:    "empty version",
			version: "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid version without patch",
			version: "1.0",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid version without minor and patch",
			version: "1",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid version format",
			version: "v1.0.0",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid version with letters",
			version: "1.0.0a",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustNewVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		wantPanic bool
	}{
		{
			name:      "valid version",
			version:   "1.0.0",
			wantPanic: false,
		},
		{
			name:      "invalid version",
			version:   "1.0",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("MustNewVersion() did not panic")
					}
				}()
				MustNewVersion(tt.version)
			} else {
				MustNewVersion(tt.version)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name string
		v    Version
		want string
	}{
		{
			name: "1.0.0",
			v:    Version("1.0.0"),
			want: "1.0.0",
		},
		{
			name: "2.3.1",
			v:    Version("2.3.1"),
			want: "2.3.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_IsEmpty(t *testing.T) {
	if Version("1.0.0").IsEmpty() {
		t.Error("Version(1.0.0).IsEmpty() returned true")
	}
	if !Version("").IsEmpty() {
		t.Error("Version(\"\").IsEmpty() returned false")
	}
}

func TestVersion_IsValid(t *testing.T) {
	tests := []struct {
		name string
		v    Version
		want bool
	}{
		{
			name: "valid 1.0.0",
			v:    Version("1.0.0"),
			want: true,
		},
		{
			name: "valid 10.20.30",
			v:    Version("10.20.30"),
			want: true,
		},
		{
			name: "invalid 1.0",
			v:    Version("1.0"),
			want: false,
		},
		{
			name: "empty",
			v:    Version(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.IsValid(); got != tt.want {
				t.Errorf("Version.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		v    Version
		want string
	}{
		{
			name: "1.0.0",
			v:    Version("1.0.0"),
			want: `"1.0.0"`,
		},
		{
			name: "2.3.1",
			v:    Version("2.3.1"),
			want: `"2.3.1"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.v)
			if err != nil {
				t.Errorf("Version.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Version.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestVersion_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    Version
		wantErr bool
	}{
		{
			name:    "valid 1.0.0",
			data:    `"1.0.0"`,
			want:    Version("1.0.0"),
			wantErr: false,
		},
		{
			name:    "valid 2.3.1",
			data:    `"2.3.1"`,
			want:    Version("2.3.1"),
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
			data:    `"1.0"`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Version
			err := json.Unmarshal([]byte(tt.data), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Version.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Version.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
