# LLM and Skill Runtime Providers Implementation

## Overview

This document describes the implementation of real LLM and Skill Runtime providers in Nexflow, replacing mock implementations with production-ready providers.

## LLM Providers

### OpenAI Provider

**Location:** `internal/infrastructure/llm/openai/provider.go`

Features:
- Full OpenAI Chat Completions API support
- Temperature and max tokens configuration
- Token usage tracking
- Error handling with proper error messages
- Structured logging
- Availability checking

Configuration:
```yaml
llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      base_url: "https://api.openai.com/v1"  # Optional
      model: "gpt-4"
```

### Anthropic Provider

**Location:** `internal/infrastructure/llm/anthropic/provider.go`

Features:
- Anthropic Messages API support
- System message handling
- Temperature and max tokens configuration
- Token usage tracking (input/output tokens)
- Stop reason tracking
- Error handling with proper error messages
- Structured logging
- Availability checking

Configuration:
```yaml
llm:
  default_provider: "anthropic"
  providers:
    anthropic:
      api_key: "${ANTHROPIC_API_KEY}"
      base_url: "https://api.anthropic.com/v1"  # Optional
      model: "claude-3-sonnet-20240229"
```

### Ollama Provider

**Location:** `internal/infrastructure/llm/ollama/provider.go`

Features:
- Ollama Chat API support
- Local LLM support (self-hosted)
- Temperature and max tokens configuration
- Token estimation
- Error handling
- Structured logging
- Availability checking

Configuration:
```yaml
llm:
  default_provider: "ollama"
  providers:
    ollama:
      base_url: "http://localhost:11434"  # Optional
      model: "llama2"
```

## Skill Runtime

### Local Runtime

**Location:** `internal/infrastructure/skills/local_runtime.go`

Features:
- Execute skills as local executables
- Environment variable input passing
- JSON output parsing
- Execution timeout handling
- Skill directory scanning
- Metadata retrieval
- Availability checking
- Structured logging
- Error handling

Configuration:
```yaml
skills:
  directory: "./skills"
  timeout_sec: 60
  sandbox_enabled: false  # Not yet implemented
```

Skill Execution:
- Skills receive input via environment variables (prefixed with `NEXFLOW_`)
- Output should be JSON for proper parsing
- Skills must be executable files (e.g., shell scripts, Python scripts with shebang)

## Provider Adapter Pattern

### LLM Provider Adapter

**Location:** `internal/infrastructure/llm/adapter.go`

Adapts infrastructure `Provider` interface to application `ports.LLMProvider` interface:
- Converts request/response types
- Implements streaming simulation
- Provides cost estimation
- Tool calling (currently delegated to Generate)

### Skill Runtime Adapter

**Location:** `internal/infrastructure/skills/runtime_adapter.go`

Adapts infrastructure `Executor` interface to application `ports.SkillRuntime` interface:
- Converts execution results to JSON strings
- Validates skill availability
- Lists available skills
- Retrieves skill metadata

## Dependency Injection

**Location:** `cmd/server/di.go`

The DI container now:
1. Reads provider configuration from config
2. Creates appropriate provider based on `llm.default_provider`
3. Wraps providers with adapters
4. Falls back to mock implementations if config is invalid or logger is not available

Configuration Flow:
1. Check if `llm.default_provider` is set
2. Find provider config in `llm.providers[default_provider]`
3. Create provider based on provider name (openai/anthropic/ollama)
4. Wrap with adapter and inject into use cases

## Configuration

See `config.example.yml` for a complete example configuration.

### Environment Variables

Security best practices:
- API keys should be stored in environment variables
- Use `${VAR_NAME}` syntax in YAML config
- Never commit API keys to version control

Example:
```bash
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
```

## Error Handling

All providers implement proper error handling:
- Invalid configuration errors
- Network errors
- API errors with descriptive messages
- Timeout errors

## Logging

Structured logging with slog:
- Debug: Request details, responses
- Info: Initialization, successful operations
- Warn: Fallbacks, unknown configurations
- Error: Failures, API errors

## Testing

- Mock implementations are retained for testing
- All providers follow the same interface
- Can be easily swapped in tests

## Future Enhancements

1. Tool calling support in LLM providers
2. Streaming support in LLM providers
3. Sandbox mode in LocalRuntime
4. Additional LLM providers (Google, Mistral, etc.)
5. Remote skill runtime (gRPC, HTTP)
6. Skill caching and optimization
