# Работа с базой данных для Nexflow

Этот модуль описывает правила работы с БД, включая SQLC, миграции и паттерны доступа.

## Технологический стек

- SQLite (по умолчанию) или PostgreSQL
- SQLC — генерация типобезопасного SQL кода
- golang-migrate/migrate — миграции
- sql.DB — стандартный драйвер Go

## SQLC

### Основные принципы

SQLC генерирует типобезопасный Go код из SQL запросов.

**Преимущества:**
- Типобезопасность (compile-time проверки)
- SQL остается SQL
- Нет runtime SQL injection
- Хорошая интеграция с Go

### Файл query.sql

```sql
-- name: CreateUser :one
INSERT INTO users (id, name, email, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users SET name = $2, email = $3, updated_at = $4
WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
```

### SQLC синтаксис

**Комментарии с метаданными:**
```sql
-- name: FunctionName :result_type
-- sqlc.slice_param: param_name
-- sqlc.noreturn: true
```

**Типы результатов:**
- `:one` — одна запись (struct)
- `:many` — много записей ([]struct)
- `:exec` — без результата (INSERT, UPDATE, DELETE)
- `:execresult` — с результатом (RowsAffected)

### Конфигурация SQLC

**sqlc.yaml:**
```yaml
version: "2"
sql:
  - engine: "postgresql" # или "sqlite"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "database"
        out: "internal/database"
        sql_package: "database/sql"
        emit_json_tags: true
        emit_interface: true
```

## Интерфейс Database

```go
type Database interface {
    // Users
    CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
    GetUserByID(ctx context.Context, id string) (User, error)
    GetUserByEmail(ctx context.Context, email string) (User, error)
    ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
    UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
    DeleteUser(ctx context.Context, id string) error

    // Sessions, Messages, Tasks, Skills, Schedules, Logs
    // ... другие методы

    Migrate(ctx context.Context) error
    Close() error
}
```

## Создание базы данных

```go
func NewDatabase(cfg *config.DatabaseConfig, logger logging.Logger) (Database, error) {
    var db *sql.DB
    var err error

    switch cfg.Type {
    case "sqlite":
        db, err = sql.Open("sqlite3", cfg.Path)
    case "postgres":
        db, err = sql.Open("postgres", cfg.Path)
    }

    // Connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    queries := New(db)
    return &DB{Queries: queries, db: db, config: cfg, logger: logger}, nil
}
```

## Connection Pooling

### Настройки

```go
db.SetMaxOpenConns(25)   // Максимальное количество открытых соединений
db.SetMaxIdleConns(25)   // Максимальное количество idle соединений
db.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни соединения
db.SetConnMaxIdleTime(10 * time.Minute) // Максимальное время простоя
```

### Рекомендации

- SQLite: MaxOpen=1 (SQLite не поддерживает concurrent writes), но SQLC генерирует safe код
- PostgreSQL: MaxOpen=25, MaxIdle=25, Lifetime=5min

### Мониторинг

```go
stats := db.Stats()
logger.Info("Connection pool stats",
    "open_connections", stats.OpenConnections,
    "in_use", stats.InUse,
    "idle", stats.Idle)
```

## Миграции

### Структура

```
migrations/
├── postgres/
│   ├── 000001_init_schema.up.sql
│   ├── 000001_init_schema.down.sql
└── sqlite/
    ├── 000001_init_schema.up.sql
    └── 000001_init_schema.down.sql
```

### Именование

Формат: `NNNNNN_description.{up,down}.sql`

Примеры:
- `000001_init_schema.up.sql`
- `000002_add_users_table.up.sql`
- `000003_add_foreign_keys.up.sql`

### Пример миграции (SQLite)

**000001_init_schema.up.sql:**
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    channel TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

**000001_init_schema.down.sql:**
```sql
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
```

### Выполнение миграций

```go
func (d *DB) Migrate(ctx context.Context) error {
    m, err := migrate.New(
        "file://migrations/"+d.config.Type,
        d.config.Path,
    )
    if err != nil {
        return err
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    return nil
}
```

## Паттерны запросов

### Использование context

**✅ ХОРОШО:**
```go
user, err := db.GetUserByID(ctx, userID)
```

**❌ ПЛОХО:**
```go
user, err := db.GetUserByID(userID) // Нет context!
```

### Транзакции

```go
func (d *DB) CreateUserWithMessages(ctx context.Context, user User, messages []Message) error {
    tx, err := d.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    qtx := d.WithTx(tx)

    user, err := qtx.CreateUser(ctx, CreateUserParams{...})
    if err != nil {
        return err
    }

    for _, msg := range messages {
        msg.UserID = user.ID
        _, err := qtx.CreateMessage(ctx, CreateMessageParams{...})
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}
```

### Пагинация

```sql
-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

// Использование
users, err := db.ListUsers(ctx, ListUsersParams{Limit: 50, Offset: 0})
```

### Фильтрация

```sql
-- name: GetLogsByLevel :many
SELECT * FROM logs WHERE level = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetLogsByDateRange :many
SELECT * FROM logs WHERE created_at >= $1 AND created_at < $2 ORDER BY created_at DESC;
```

## Тестирование

### Unit тесты с mock

```go
func TestUserService_CreateUser(t *testing.T) {
    mockDB := &MockDatabase{users: make(map[string]User)}
    service := NewUserService(mockDB, logger)

    user, err := service.CreateUser(ctx, CreateUserRequest{Name: "John"})
    require.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}
```

### Интеграционные тесты

```go
func TestDatabase_Integration(t *testing.T) {
    tmpDir, _ := os.MkdirTemp("", "nexflow-test-*")
    defer os.RemoveAll(tmpDir)

    db, err := NewDatabase(&config.DatabaseConfig{Type: "sqlite", Path: tmpDir + "/test.db"}, logger)
    require.NoError(t, err)
    defer db.Close()

    ctx := context.Background()
    err = db.Migrate(ctx)
    require.NoError(t, err)

    user, err := db.CreateUser(ctx, CreateUserParams{ID: "user-1", Name: "Test", ...})
    require.NoError(t, err)
    assert.Equal(t, "user-1", user.ID)
}
```

## Критические правила

1. ВСЕГДА используйте SQLC для типобезопасных запросов
2. НИКОГДА не используйте string formatting для SQL (SQL injection!)
3. ВСЕГДА используйте `context.Context` для всех запросов
4. НИКОГДА не игнорируйте ошибки БД
5. ВСЕГДА оборачивайте транзакции в Begin/Commit/Rollback
6. НИКОГДА не забывайте закрывать соединения (defer db.Close())
7. ВСЕГДА настройте connection pool правильно
8. НИКОГДА не создавайте connection leaks (defer rows.Close())
9. ВСЕГДА пишите миграции в up/down формате
10. НИКОГДА не коммите миграции без review

---

**Памятка:** База данных — это критический компонент. Будьте осторожны и тестируйте всё.
