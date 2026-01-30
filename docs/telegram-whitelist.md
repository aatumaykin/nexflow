# Telegram Whitelist Configuration

This document explains how to configure whitelist security for Telegram connector.

## Overview

Nexflow implements a whitelist-based security model for Telegram:
- Only authorized users/chats can interact with the bot
- At least one of `allowed_users` or `allowed_chats` must be specified
- Unauthorized access attempts are logged and blocked

## Configuration

### Required Fields

```yaml
channels:
  telegram:
    enabled: true
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users: ["123456789"]  # Required if allowed_chats is empty
    allowed_chats: []             # Required if allowed_users is empty
```

### Example Configurations

#### Single User (Private Chat)

```yaml
channels:
  telegram:
    enabled: true
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users: ["123456789"]  # Your Telegram user ID
    allowed_chats: []
```

#### Multiple Users

```yaml
channels:
  telegram:
    enabled: true
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users:
      - "123456789"
      - "987654321"
      - "555555555"
    allowed_chats: []
```

#### Group Chat Only

```yaml
channels:
  telegram:
    enabled: true
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users: []
    allowed_chats: ["-1001234567890"]  # Group chat ID
```

#### Private Chat + Group Chat

```yaml
channels:
  telegram:
    enabled: true
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users: ["123456789"]
    allowed_chats: ["-1001234567890"]
```

## Getting Your Telegram User/Chat IDs

### Method 1: Using @userinfobot

1. Open Telegram
2. Search for `@userinfobot`
3. Send any message
4. The bot will reply with your User ID

### Method 2: Using the Bot

1. Temporarily set `allowed_users: []` and `allowed_chats: []`
2. Start the bot
3. Send a message to your bot
4. Check the logs for the warning message:
   ```
   Message from unauthorized user/chat ignored
       chat_id=123456789 user_id=987654321 username=yourname
   ```
5. Add the `user_id` or `chat_id` to your config

### Method 3: Using a Forwarded Message

1. Forward any message from your target chat to the bot
2. Check the logs for `forward_from_user_id` or chat information

## Security Best Practices

1. **Always use whitelist**: Never leave both `allowed_users` and `allowed_chats` empty
2. **Use environment variables**: Store bot token in `TELEGRAM_BOT_TOKEN` environment variable
3. **Monitor logs**: Check for unauthorized access attempts
4. **Regular updates**: Review and update allowed users/chats periodically
5. **Use unique IDs**: Telegram user/chat IDs are unique and permanent

## Troubleshooting

### Bot not responding

1. Check that `enabled: true` is set
2. Verify `bot_token` is correct
3. Ensure your user/chat ID is in the whitelist
4. Check logs for error messages

### "Unauthorized user" warning

This means your user/chat ID is not in the whitelist:
```
Message from unauthorized user/chat ignored
    chat_id=123456789 user_id=987654321 username=yourname
```

Add the `user_id` (for private chats) or `chat_id` (for groups) to your config.

### Validation error

If you get this error:
```
telegram: at least one of allowed_users or allowed_chats must be specified for security when telegram is enabled
```

Add at least one user or chat to the whitelist.

## Testing

You can test the whitelist with unit tests:

```bash
go test ./internal/infrastructure/channels/telegram/... -v
```

## Related Documentation

- [Channels Documentation](../channels.md)
- [Security Rules](../rules/security.md)
- [Getting Telegram Bot Token](../channels.md#getting-your-telegram-bot-token)
