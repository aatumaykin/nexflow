package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Logger represents a structured logger interface
type Logger interface {
	// Debug logs a message at debug level
	Debug(msg string, args ...any)
	// Info logs a message at info level
	Info(msg string, args ...any)
	// Warn logs a message at warn level
	Warn(msg string, args ...any)
	// Error logs a message at error level
	Error(msg string, args ...any)
	// With returns a new logger with additional fields
	With(args ...any) Logger
	// WithContext returns a new logger with context
	WithContext(ctx context.Context) Logger
	// DebugContext logs a message at debug level with context
	DebugContext(ctx context.Context, msg string, args ...any)
	// InfoContext logs a message at info level with context
	InfoContext(ctx context.Context, msg string, args ...any)
	// WarnContext logs a message at warn level with context
	WarnContext(ctx context.Context, msg string, args ...any)
	// ErrorContext logs a message at error level with context
	ErrorContext(ctx context.Context, msg string, args ...any)
}

// SlogLogger is a slog-based implementation of Logger
type SlogLogger struct {
	logger *slog.Logger
	ctx    context.Context
}

// New creates a new logger with the specified level and format
func New(level string, format string) (Logger, error) {
	// Parse log level
	logLevel, err := parseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Mask secret fields
			if shouldMask(a.Key) {
				return slog.String(a.Key, maskValue(a.Value.String()))
			}
			return a
		},
	}

	// Create handler based on format
	var handler slog.Handler
	if strings.ToLower(format) == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	// Add default attributes
	handler = handler.WithAttrs([]slog.Attr{
		slog.String("source", "nexflow"),
	})

	return &SlogLogger{
		logger: slog.New(handler),
		ctx:    context.Background(),
	}, nil
}

// parseLevel converts a string level to slog.Level
func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	case "fatal":
		return slog.LevelError, nil // Use error level for fatal
	default:
		return slog.LevelInfo, fmt.Errorf("unknown level: %s", level)
	}
}

// Debug logs a message at debug level
func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.DebugContext(l.ctx, msg, args...)
}

// Info logs a message at info level
func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.InfoContext(l.ctx, msg, args...)
}

// Warn logs a message at warn level
func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.WarnContext(l.ctx, msg, args...)
}

// Error logs a message at error level
func (l *SlogLogger) Error(msg string, args ...any) {
	l.logger.ErrorContext(l.ctx, msg, args...)
}

// With returns a new logger with additional fields
func (l *SlogLogger) With(args ...any) Logger {
	return &SlogLogger{
		logger: l.logger.With(args...),
		ctx:    l.ctx,
	}
}

// WithContext returns a new logger with context
func (l *SlogLogger) WithContext(ctx context.Context) Logger {
	return &SlogLogger{
		logger: l.logger,
		ctx:    ctx,
	}
}

// DebugContext logs a message at debug level with context
func (l *SlogLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

// InfoContext logs a message at info level with context
func (l *SlogLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

// WarnContext logs a message at warn level with context
func (l *SlogLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

// ErrorContext logs a message at error level with context
func (l *SlogLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}
