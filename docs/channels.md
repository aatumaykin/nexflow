# Channel Connectors Implementation

## Overview

This document describes the implementation of channel connectors in Nexflow. Channel connectors enable Nexflow to communicate with users through different platforms such as Telegram, Discord, and Web.

## Architecture

All channel connectors implement the `Connector` interface defined in `internal/infrastructure/channels/connector.go`:

```go
type Connector interface {
    Name() string
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    SendResponse(ctx context.Context, userID string, response *Response) error
    Incoming() <-chan *Message
    IsRunning() bool
    GetUser(ctx context.Context, channelUserID string) (*entity.User, error)
    CreateUser(ctx context.Context, channelUserID string) (*entity.User, error)
}
```

## Telegram Connector

**Location:** `internal/infrastructure/channels/telegram/telegram.go`

### Features

- Full Telegram Bot API support using [go-telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api)
- Support for both webhook and long polling modes
- Whitelist-based security (allowed_users and allowed_chats)
- Automatic user creation and retrieval
- Rich metadata in messages (user info, chat type, message details)
- Graceful shutdown
- Structured logging

### Configuration

```yaml
channels:
  telegram:
    enabled: true
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users: ["123456789"]  # List of allowed Telegram user IDs
    allowed_chats: []             # List of allowed Telegram chat IDs
    webhook_url: ""               # Optional: use webhook instead of long polling
```

### Security

The Telegram connector implements a whitelist-based security model:

- **allowed_users**: Only specific Telegram user IDs can interact with the bot
- **allowed_chats**: Only specific chat IDs can interact with the bot
- At least one of these lists must be specified when Telegram is enabled

### Message Format

Incoming messages from Telegram are normalized into the `Message` struct:

```go
type Message struct {
    UserID    string                 // Format: "telegram_user_id:chat_id"
    ChannelID string                 // Chat ID as string
    Content   string                 // Message text
    Metadata  map[string]interface{} // Additional info (username, chat_type, etc.)
}
```

Metadata includes:
- `chat_id`: Telegram chat ID
- `user_id`: Telegram user ID
- `username`: Telegram username
- `first_name`: User's first name
- `last_name`: User's last name
- `message_id`: Telegram message ID
- `is_command`: Whether message is a bot command
- `chat_type`: Chat type (private, group, supergroup, channel)

### Getting Your Telegram Bot Token

1. Open [BotFather](https://t.me/botfather) in Telegram
2. Send `/newbot` command
3. Follow instructions to create a new bot
4. BotFather will provide you with an API token

### Getting Your Telegram User/Chat IDs

You can get your Telegram user ID by:

1. Using the `@userinfobot` bot in Telegram
2. Or send a message to your bot and check the logs for the user ID

### Usage Example

The Telegram connector is automatically initialized and started by the DI container when `enabled: true` is set in the configuration.

```go
// In DI container (cmd/server/di.go):
c.telegramConnector = telegramconn.NewConnector(
    c.config.Channels.Telegram,
    c.userRepo,
    slogLogger.GetSlogLogger(),
)
```

### Testing

The Telegram connector has comprehensive unit tests covering:

- Connector initialization
- Start/Stop lifecycle
- Message sending
- User creation and retrieval
- Whitelist validation
- Webhook and polling modes

Run tests with:
```bash
go test ./internal/infrastructure/channels/telegram/... -v
```

### Error Handling

The connector handles various error scenarios:

- Invalid bot token
- Network failures
- Unauthorized users/chats
- Malformed user IDs
- Channel buffer full

All errors are logged with appropriate context for debugging.

## Discord Connector

**Location:** `internal/infrastructure/channels/mock/discord.go` (mock implementation)

Currently, a mock implementation is available. Real Discord connector implementation is a future task.

## Web Connector

**Location:** `internal/infrastructure/channels/mock/web.go` (mock implementation)

Currently, a mock implementation is available. Real Web connector implementation is a future task.

## Mock Connectors

For testing and development purposes, mock implementations are available:

```go
import "github.com/atumaikin/nexflow/internal/infrastructure/channels/mock"

mockConnector := mock.NewTelegramConnector()
```

Mock connectors provide:
- In-memory user storage
- Test message injection
- Response tracking
- No external dependencies

## Future Enhancements

Potential improvements to channel connectors:

- Support for inline keyboards (Telegram)
- File uploads/downloads
- Voice messages
- Rich message formatting
- Threaded conversations
- Rate limiting
- Message persistence
- Retry logic for failed sends
- Health checks
- Metrics and monitoring

## Troubleshooting

### Telegram Bot Not Responding

1. Check that `bot_token` is correct and valid
2. Verify that the bot is started and running
3. Check logs for error messages
4. Ensure your user/chat ID is in the whitelist
5. Verify network connectivity to Telegram servers

### Webhook Not Receiving Updates

1. Ensure your webhook URL is publicly accessible
2. Verify SSL certificate is valid (Telegram requires HTTPS)
3. Check that your server's firewall allows incoming connections
4. Test webhook URL using Telegram's API

### Polling Mode Issues

1. Check for network connectivity
2. Verify bot token permissions
3. Check logs for polling errors
4. Consider using webhook mode for better reliability

## References

- [Telegram Bot API Documentation](https://core.telegram.org/bots/api)
- [go-telegram-bot-api Library](https://github.com/go-telegram-bot-api/telegram-bot-api)
- [BotFather](https://t.me/botfather)
- [Discord API Documentation](https://discord.com/developers/docs/intro)
