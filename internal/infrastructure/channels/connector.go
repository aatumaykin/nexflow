package channels

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// Message represents a message from a channel
type Message struct {
	UserID    string // Channel-specific user ID
	ChannelID string // Channel-specific channel or chat ID
	Content   string // Message content
	Metadata  map[string]interface{}
}

// Response represents a response to send back to a channel
type Response struct {
	Type      ResponseType           // Type of response
	Content   string                 // Text content (for text, photo, document, audio, video messages)
	Caption   string                 // Caption for media messages
	Media     *MediaContent          // Media content (for photo, document, audio, video)
	Buttons   []InlineButton         // Inline buttons
	Markup    interface{}            // Custom reply markup (channel-specific)
	MessageID string                 // For editing existing messages
	Metadata  map[string]interface{} // Additional metadata
}

// ResponseType represents the type of response message
type ResponseType string

const (
	ResponseTypeText     ResponseType = "text"     // Plain text message
	ResponseTypePhoto    ResponseType = "photo"    // Photo with optional caption
	ResponseTypeDocument ResponseType = "document" // Document/file with optional caption
	ResponseTypeAudio    ResponseType = "audio"    // Audio file with optional caption
	ResponseTypeVideo    ResponseType = "video"    // Video file with optional caption
	ResponseTypeSticker  ResponseType = "sticker"  // Sticker
)

// MediaContent represents media content for a response
type MediaContent struct {
	URL      string // URL to media file
	FileID   string // Telegram file ID (for reusing existing files)
	FileData []byte // Raw file data (for uploading new files)
	FileName string // File name (for documents)
}

// InlineButton represents an inline keyboard button
type InlineButton struct {
	Text         string // Button text
	Data         string // Callback data (up to 64 bytes)
	URL          string // URL to open (for link buttons)
	InlineData   string // Inline query data (for inline mode)
	SwitchInline string // Switch to inline query data
	CallbackData string // Callback data (alternative field)
}

// InlineKeyboard represents a row of inline buttons
type InlineKeyboard struct {
	Buttons []InlineButton
}

// InlineMarkup represents a complete inline keyboard markup
type InlineMarkup struct {
	Keyboard [][]InlineButton
}

// Connector defines the interface for all channel connectors
type Connector interface {
	// Name returns the name of the channel (telegram, discord, web, etc.)
	Name() string

	// Start initializes and starts the connector
	Start(ctx context.Context) error

	// Stop gracefully stops the connector
	Stop(ctx context.Context) error

	// SendResponse sends a response to a user
	SendResponse(ctx context.Context, userID string, response *Response) error

	// Incoming returns a channel for incoming messages
	Incoming() <-chan *Message

	// IsRunning returns whether the connector is currently running
	IsRunning() bool

	// GetUser retrieves a user by channel-specific ID
	GetUser(ctx context.Context, channelUserID string) (*entity.User, error)

	// CreateUser creates a new user in the system
	CreateUser(ctx context.Context, channelUserID string) (*entity.User, error)
}
