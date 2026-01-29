package config

// Default configuration constants

// Timeouts
const (
	// Default timeout in seconds for skill execution
	DefaultSkillTimeoutSec = 30

	// HTTP server timeouts
	DefaultIdleTimeoutSec = 60
	DefaultReadTimeoutSec = 60
	DefaultWriteTimeoutSec = 60
)

// Limits
const (
	// Default channel buffer size for incoming messages
	DefaultChannelBufferSize = 100

	// Default max tokens for LLM requests
	DefaultMaxTokens = 100

	// Default connection pool size
	DefaultMaxOpenConnections = 25
	DefaultMaxIdleConnections = 25
)

// Costs
const (
	// Default cost per 1000 tokens (USD)
	DefaultCostPer1000Tokens = 0.02
)

// Version
const (
	// Current SQLC version (from generate output)
	SQLCVersion = "1.30.0"
)
