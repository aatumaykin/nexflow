package database

import (
	"context"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/shared/config"
	"github.com/atumaikin/nexflow/internal/shared/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatabase(t *testing.T) {
	tests := []struct {
		name          string
		config        *config.DatabaseConfig
		opts          []Option
		wantErr       bool
		errContains   string
		setupDBConfig func() *DBConfig
	}{
		{
			name: "Valid SQLite config with defaults",
			config: &config.DatabaseConfig{
				Type:           "sqlite",
				Path:           ":memory:",
				MigrationsPath: "migrations",
			},
			wantErr: false,
			setupDBConfig: func() *DBConfig {
				return &DBConfig{
					Type:            "sqlite",
					Path:            ":memory:",
					MigrationsPath:  "migrations",
					MaxOpenConns:    25,
					MaxIdleConns:    25,
					ConnMaxLifetime: 5 * time.Minute,
				}
			},
		},
		{
			name: "Valid PostgreSQL config",
			config: &config.DatabaseConfig{
				Type:            "postgres",
				Path:            "postgres://localhost/test?sslmode=disable",
				MigrationsPath:  "migrations",
				MaxOpenConns:    10,
				MaxIdleConns:    5,
				ConnMaxLifetime: 10 * time.Minute,
			},
			wantErr: false,
			setupDBConfig: func() *DBConfig {
				return &DBConfig{
					Type:            "postgres",
					Path:            "postgres://localhost/test?sslmode=disable",
					MigrationsPath:  "migrations",
					MaxOpenConns:    10,
					MaxIdleConns:    5,
					ConnMaxLifetime: 10 * time.Minute,
				}
			},
		},
		{
			name: "Missing database type",
			config: &config.DatabaseConfig{
				Path: "./test.db",
			},
			wantErr:     true,
			errContains: "database type is required",
		},
		{
			name: "Missing database path",
			config: &config.DatabaseConfig{
				Type: "sqlite",
			},
			wantErr:     true,
			errContains: "database path is required",
		},
		{
			name: "Unsupported database type",
			config: &config.DatabaseConfig{
				Type: "mysql",
				Path: "./test.db",
			},
			wantErr:     true,
			errContains: "unsupported database type",
		},
		{
			name: "With custom logger option",
			config: &config.DatabaseConfig{
				Type:           "sqlite",
				Path:           ":memory:",
				MigrationsPath: "migrations",
			},
			opts:    []Option{WithLogger(&TestLogger{})},
			wantErr: false,
			setupDBConfig: func() *DBConfig {
				return &DBConfig{
					Type:            "sqlite",
					Path:            ":memory:",
					MigrationsPath:  "migrations",
					MaxOpenConns:    25,
					MaxIdleConns:    25,
					ConnMaxLifetime: 5 * time.Minute,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip PostgreSQL test if server is not available
			if tt.name == "Valid PostgreSQL config" {
				_, err := NewDatabase(tt.config, tt.opts...)
				if err != nil {
					t.Skipf("PostgreSQL server not available: %v", err)
				}
			}

			db, err := NewDatabase(tt.config, tt.opts...)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, db)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, db)
			require.NotNil(t, db.(*DB).db)
			assert.Equal(t, tt.setupDBConfig().Type, db.(*DB).config.Type)
			assert.Equal(t, tt.setupDBConfig().MaxOpenConns, db.(*DB).config.MaxOpenConns)
			assert.Equal(t, tt.setupDBConfig().MaxIdleConns, db.(*DB).config.MaxIdleConns)
			assert.Equal(t, tt.setupDBConfig().ConnMaxLifetime, db.(*DB).config.ConnMaxLifetime)
		})
	}
}

func TestNewDatabase_Close(t *testing.T) {
	// Setup test config
	config := &config.DatabaseConfig{
		Type:           "sqlite",
		Path:           ":memory:",
		MigrationsPath: "migrations",
	}

	db, err := NewDatabase(config)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Verify connection is alive before close
	dbImpl := db.(*DB)
	err = dbImpl.GetDB().Ping()
	require.NoError(t, err, "Database should be accessible before close")

	// Close database
	err = db.Close()
	require.NoError(t, err)

	// Verify connection is closed (Ping should fail)
	err = dbImpl.GetDB().Ping()
	assert.Error(t, err, "Database should not be accessible after close")
}

func TestDB_GetDB(t *testing.T) {
	config := &config.DatabaseConfig{
		Type:           "sqlite",
		Path:           ":memory:",
		MigrationsPath: "migrations",
	}

	db, err := NewDatabase(config)
	require.NoError(t, err)
	require.NotNil(t, db)

	dbImpl := db.(*DB)
	sqlDB := dbImpl.GetDB()
	require.NotNil(t, sqlDB)
	assert.NoError(t, sqlDB.PingContext(context.Background()))

	err = db.Close()
	require.NoError(t, err)
}

type TestLogger struct {
	messages []string
}

func (l *TestLogger) Info(msg string, args ...any) {
	l.messages = append(l.messages, msg)
}

func (l *TestLogger) Error(msg string, args ...any) {
	l.messages = append(l.messages, "ERROR: "+msg)
}

func (l *TestLogger) Debug(msg string, args ...any) {
	l.messages = append(l.messages, "DEBUG: "+msg)
}

func (l *TestLogger) Warn(msg string, args ...any) {
	l.messages = append(l.messages, "WARN: "+msg)
}

func (l *TestLogger) With(args ...any) logging.Logger {
	return l
}

func (l *TestLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.Info(msg, args...)
}

func (l *TestLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.Error(msg, args...)
}

func (l *TestLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.Debug(msg, args...)
}

func (l *TestLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.Warn(msg, args...)
}

func (l *TestLogger) WithContext(ctx context.Context) logging.Logger {
	return l
}
