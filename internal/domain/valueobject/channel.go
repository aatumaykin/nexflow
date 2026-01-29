package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrInvalidChannel is returned when an invalid channel type is provided.
	ErrInvalidChannel = errors.New("invalid channel")
	// ErrEmptyChannel is returned when an empty channel is provided.
	ErrEmptyChannel = errors.New("channel cannot be empty")
)

// Channel represents the type of communication channel.
// It's a value object that ensures type safety for channel types.
type Channel string

const (
	// ChannelTelegram represents the Telegram messaging platform.
	ChannelTelegram Channel = "telegram"
	// ChannelDiscord represents the Discord messaging platform.
	ChannelDiscord Channel = "discord"
	// ChannelWeb represents the web interface.
	ChannelWeb Channel = "web"
)

// String returns the string representation of the channel.
func (c Channel) String() string {
	return string(c)
}

// IsValid checks if the channel is valid.
func (c Channel) IsValid() bool {
	switch c {
	case ChannelTelegram, ChannelDiscord, ChannelWeb:
		return true
	default:
		return false
	}
}

// IsTelegram returns true if the channel is Telegram.
func (c Channel) IsTelegram() bool {
	return c == ChannelTelegram
}

// IsDiscord returns true if the channel is Discord.
func (c Channel) IsDiscord() bool {
	return c == ChannelDiscord
}

// IsWeb returns true if the channel is Web.
func (c Channel) IsWeb() bool {
	return c == ChannelWeb
}

// Equals checks if the channel equals another channel.
func (c Channel) Equals(other Channel) bool {
	return c == other
}

// MarshalJSON implements json.Marshaler interface.
func (c Channel) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (c *Channel) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyChannel
	}
	ch := Channel(str)
	if !ch.IsValid() {
		*c = "" // Reset to empty on error
		return fmt.Errorf("%w: %s", ErrInvalidChannel, str)
	}
	*c = ch
	return nil
}

// NewChannel creates a new Channel from a string.
// Returns an error if the string is not a valid channel.
func NewChannel(channel string) (Channel, error) {
	if channel == "" {
		return "", ErrEmptyChannel
	}
	c := Channel(channel)
	if !c.IsValid() {
		return "", ErrInvalidChannel
	}
	return c, nil
}

// MustNewChannel creates a new Channel from a string.
// Panics if the string is not a valid channel.
func MustNewChannel(channel string) Channel {
	c, err := NewChannel(channel)
	if err != nil {
		panic(err)
	}
	return c
}
