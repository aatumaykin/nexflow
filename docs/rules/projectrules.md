# Общие правила проекта Nexflow

Этот модуль описывает общие правила работы с проектом, включая структуру, соглашения и репозиторий.

## Структура проекта

### Основные директории

```
nexflow/
├── cmd/              # Entry points (main.go)
├── internal/         # Private packages (не импортируется извне)
│   ├── config/       # Конфигурация
│   ├── database/     # База данных, SQLC
│   ├── logging/      # Структурированное логирование
│   ├── channels/     # Connectors (Telegram, Discord, Web)
│   ├── llm/         # LLM провайдеры и роутинг
│   └── skills/      # Skills execution
├── pkg/              # Public libraries
├── docs/
│   └── rules/        # Модули правил
├── migrations/       # SQL миграции
├── skills/           # Навыки (SKILL.md)
└── .cursor/rules/    # agent_behavior.mdc
```

### Правила организации

**cmd/** — только `main.go`, бизнес-логики нет

**internal/** — внутренние пакеты, бизнес-логика, интерфейсы для тестируемости

**pkg/** — публичные библиотеки, стабильное API, документация

**docs/** — документация, `docs/rules/` — правила для ИИ-агентов

## Соглашения по именованию

### Пакеты

- lowercase, короткие, понятные
- Избегайте сокращений (кроме общепринятых: config, db, http)
- Имя пакета = имя директории

```go
package config  // ✅
package database // ✅

package conf   // ❌
package dbase  // ❌
```

### Файлы

- lowercase с подчеркиваниями
- Имя файла отражает основной тип/функцию
- Тесты: `file_test.go`

```go
config.go       // ✅
database.go     // ✅
config_test.go  // ✅
Config.go       // ❌
```

### Переменные и функции

- Локальные переменные: `camelCase`
- Экспортируемые функции: `PascalCase`
- Receiver: 1-2 символа

```go
userName := "john"
func NewDatabase() (*DB, error) { ... }
func (d *DB) Close() error { ... }
```

### Типы и интерфейсы

- `PascalCase`
- Интерфейсы: `-er` суффикс (`Reader`, `Writer`, `Database`)

```go
type User struct { ... }
type Logger interface { ... }
type Database interface { ... }
```

## Структура файлов

```go
// Пакетный комментарий
package packagename

import "fmt"

// Константы
const MaxRetries = 3

// Переменные
var defaultTimeout = 30 * time.Second

// Типы
type User struct { ... }

// Интерфейсы
type Database interface { ... }

// Конструкторы
func NewUser(name, email string) *User { ... }

// Методы
func (u *User) Validate() error { ... }

// Функции пакета
func CreateUser(name, email string) (User, error) { ... }
```

### Форматирование

- `gofmt` или `go fmt`
- Максимальная длина строки: 100-120 символов
- Отступы: TAB

### Структурные теги

```go
type User struct {
    ID        string    `json:"id" yaml:"id" db:"id"`
    Name      string    `json:"name" yaml:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" yaml:"created_at" db:"created_at"`
}
```

## Работа с репозиторием

### Ветвление

**Основная:** `main`

**Типы веток:**
- `feature/feature-name` — новая функциональность
- `fix/bug-name` — исправление бага
- `refactor/refactor-name` — рефакторинг
- `docs/doc-name` — документация
- `test/test-name` — тесты

### Commit сообщения (Conventional Commits)

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

**Типы:** feat, fix, refactor, docs, test, chore, style

**Примеры:**
```
feat(channels): add Telegram connector implementation

- Add TelegramBot connector
- Support user authentication

Closes #123

fix(database): prevent connection pool exhaustion

- Set max open connections (25)
- Add connection lifetime (5 min)

fixes #456
```

### Pull Requests

**Название:** следует формату Conventional Commits

**Описание:**
```markdown
## Summary
Краткое описание (2-3 предложения)

## Changes
- Change 1
- Change 2

## Testing
Как тестировали

## Checklist
- [ ] Tests pass
- [ ] Code follows project rules
```

## CI/CD

### GitHub Actions

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25.5'
      - name: Run tests
        run: go test ./...
      - name: Build
        run: go build ./cmd/server
```

### Локальная проверка

```bash
go fmt ./...
go test ./...
go build ./cmd/server
```

## Конфигурация

### ENV переменные

```bash
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export TELEGRAM_BOT_TOKEN="123456:ABC-..."
```

### Конфигурационные файлы

```yaml
server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "${DATABASE_PATH}"

llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      base_url: "https://api.openai.com/v1"
      model: "gpt-4"

channels:
  telegram:
    bot_token: "${TELEGRAM_BOT_TOKEN}"

logging:
  level: "info"
  format: "json"
```

## Документация

### Кодовая документация

```go
// Package database provides database access and SQLC-generated queries.
package database

// NewDatabase creates a new database connection.
// It supports both SQLite and PostgreSQL databases.
func NewDatabase(cfg *config.DatabaseConfig, logger logging.Logger) (Database, error) { ... }

// Database defines interface for all database operations.
// This interface allows for easy mocking in tests.
type Database interface {
    // CreateUser creates a new user in database.
    CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
}
```

## Go-специфичные соглашения

### Context

Всегда используйте `context.Context` для БД, HTTP, LLM, долгих операций.

```go
func (d *DB) GetUserByID(ctx context.Context, id string) (User, error) { ... }
```

### Обработка ошибок

Оборачивайте ошибки с `%w`:

```go
return fmt.Errorf("failed to open database: %w", err)
```

### Инициализация

Используйте конструкторы:

```go
// По умолчанию использует NoopLogger
db, err := database.NewDatabase(&cfg.Database)

// Или с кастомным логгером
db, err := database.NewDatabase(&cfg.Database, database.WithLogger(logger))
```

### Интерфейсы

Определяйте интерфейсы там, где они нужны (для тестирования):

```go
type Logger interface {
    Info(msg string, args ...any)
    Error(msg string, args ...any)
}
```

## Добавление новых модулей

1. Создайте пакет в `internal/` или `pkg/`
2. Определите интерфейс (если нужен)
3. Реализуйте функции
4. Валидация (если применимо)
5. Напишите тесты
6. Добавьте логирование
7. Обновите конфиг (если нужно)

## Критические правила

1. ВСЕГДА используйте `gofmt`
2. НИКОГДА не коммитьте секреты
3. ВСЕГДА пишите тесты
4. НИКОГДА не игнорируйте ошибки
5. ВСЕГДА используйте `context.Context`
6. НИКОГДА не создавайте циклические зависимости
7. ВСЕГДА оборачивайте ошибки с `%w`
8. НИКОГДА не хардкодите секреты
9. ВСЕГДА документируйте экспортируемые функции
10. НИКОГДА не прерывайте проверки CI/CD

---

**Памятка:** Эти правила обеспечивают качество и согласованность кода.
