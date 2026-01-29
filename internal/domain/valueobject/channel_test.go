package valueobject

import (
	"encoding/json"
	"testing"
)

func TestNewChannel(t *testing.T) {
	tests := []struct {
		name    string
		channel string
		want    Channel
		wantErr bool
	}{
		{
			name:    "valid telegram channel",
			channel: "telegram",
			want:    ChannelTelegram,
			wantErr: false,
		},
		{
			name:    "valid discord channel",
			channel: "discord",
			want:    ChannelDiscord,
			wantErr: false,
		},
		{
			name:    "valid web channel",
			channel: "web",
			want:    ChannelWeb,
			wantErr: false,
		},
		{
			name:    "empty channel",
			channel: "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid channel",
			channel: "invalid",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChannel(tt.channel)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustNewChannel(t *testing.T) {
	tests := []struct {
		name      string
		channel   string
		wantPanic bool
	}{
		{
			name:      "valid channel",
			channel:   "telegram",
			wantPanic: false,
		},
		{
			name:      "invalid channel",
			channel:   "invalid",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("MustNewChannel() did not panic")
					}
				}()
				MustNewChannel(tt.channel)
			} else {
				MustNewChannel(tt.channel)
			}
		})
	}
}

func TestChannel_String(t *testing.T) {
	tests := []struct {
		name string
		c    Channel
		want string
	}{
		{
			name: "telegram",
			c:    ChannelTelegram,
			want: "telegram",
		},
		{
			name: "discord",
			c:    ChannelDiscord,
			want: "discord",
		},
		{
			name: "web",
			c:    ChannelWeb,
			want: "web",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Channel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    Channel
		want bool
	}{
		{
			name: "valid telegram",
			c:    ChannelTelegram,
			want: true,
		},
		{
			name: "valid discord",
			c:    ChannelDiscord,
			want: true,
		},
		{
			name: "valid web",
			c:    ChannelWeb,
			want: true,
		},
		{
			name: "invalid",
			c:    Channel("invalid"),
			want: false,
		},
		{
			name: "empty",
			c:    Channel(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("Channel.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_IsTelegram(t *testing.T) {
	if !ChannelTelegram.IsTelegram() {
		t.Error("ChannelTelegram.IsTelegram() returned false")
	}
	if ChannelDiscord.IsTelegram() {
		t.Error("ChannelDiscord.IsTelegram() returned true")
	}
}

func TestChannel_IsDiscord(t *testing.T) {
	if !ChannelDiscord.IsDiscord() {
		t.Error("ChannelDiscord.IsDiscord() returned false")
	}
	if ChannelWeb.IsDiscord() {
		t.Error("ChannelWeb.IsDiscord() returned true")
	}
}

func TestChannel_IsWeb(t *testing.T) {
	if !ChannelWeb.IsWeb() {
		t.Error("ChannelWeb.IsWeb() returned false")
	}
	if ChannelTelegram.IsWeb() {
		t.Error("ChannelTelegram.IsWeb() returned true")
	}
}

func TestChannel_Equals(t *testing.T) {
	if !ChannelTelegram.Equals(ChannelTelegram) {
		t.Error("ChannelTelegram.Equals(ChannelTelegram) returned false")
	}
	if ChannelTelegram.Equals(ChannelDiscord) {
		t.Error("ChannelTelegram.Equals(ChannelDiscord) returned true")
	}
}

func TestChannel_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		c    Channel
		want string
	}{
		{
			name: "telegram",
			c:    ChannelTelegram,
			want: `"telegram"`,
		},
		{
			name: "discord",
			c:    ChannelDiscord,
			want: `"discord"`,
		},
		{
			name: "web",
			c:    ChannelWeb,
			want: `"web"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.c)
			if err != nil {
				t.Errorf("Channel.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Channel.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestChannel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    Channel
		wantErr bool
	}{
		{
			name:    "valid telegram",
			data:    `"telegram"`,
			want:    ChannelTelegram,
			wantErr: false,
		},
		{
			name:    "valid discord",
			data:    `"discord"`,
			want:    ChannelDiscord,
			wantErr: false,
		},
		{
			name:    "valid web",
			data:    `"web"`,
			want:    ChannelWeb,
			wantErr: false,
		},
		{
			name:    "empty",
			data:    `""`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid",
			data:    `"invalid"`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Channel
			err := json.Unmarshal([]byte(tt.data), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Channel.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Channel.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
