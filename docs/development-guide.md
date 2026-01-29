# Руководство разработчика Nexflow

Это руководство предназначено для разработчиков, которые хотят внести вклад в проект Nexflow или понять его внутреннюю структуру.

## Содержание

- [Обзор проекта](#обзор-проекта)
- [Структура проекта](#структура-проекта)
- [Начало работы](#начало-работы)
- [Архитектура](#архитектура)
- [Разработка новых компонентов](#разработка-новых-компонентов)
- [Тестирование](#тестирование)
- [Логирование](#логирование)
- [Конфигурация](#конфигурация)
- [База данных](#база-данных)
- [Деплой](#деплой)

## Обзор проекта

**Nexflow** — self-hosted ИИ-агент на Go, управляющий цифровыми потоками задач через multiple channels (Telegram, Discord, Web UI) с LLM-провайдерами (Anthropic, OpenAI, Ollama и др.) и навыками (skills).

### Технологический стек

- **Язык:** Go 1.25.5+
- **БД:** SQLite/Postgres
- **ORM:** SQLC
- **Логирование:** slog
- **Конфигурация:** YAML/JSON + ENV
- **CI/CD:** GitHub Actions

### Ключевые принципы

1. **Чистая слоистая архитектура** — разделение на Domain, Application, Infrastructure слои
2. **Безопасность** — приоритет выше всего
3. **Тестирование** — unit тесты для каждого модуля
4. **Документация** — godoc comments для всех экспортируемых типов

## Структура проекта

```
nexflow/
├── cmd/                      # Entry points
│   └── server/
│       └── main.go           # Главный файл приложения
├── internal/                 # Приватные пакеты
│   ├── application/          # Application слой
│   │   ├── dto/             # Data Transfer Objects
│   │   ├── ports/           # Ports (interfaces)
│   │   └── usecase/         # Use cases
│   ├── domain/              # Domain слой
│   │   ├── entity/          # Сущности
│   │   └── repository/      # Repository интерфейсы
│   ├── infrastructure/      # Infrastructure слой
│   │   ├── channels/        # Connectors (Telegram, Discord, Web)
│   │   ├── llm/            # LLM провайдеры
│   │   └── persistence/     # Persistence (БД)
│   └── shared/             # Shared utilities
│       ├── config/          # Конфигурация
│       ├── logging/         # Логирование
│       └── utils/           # Helper функции
├── pkg/                     # Публичные библиотеки
├── docs/                    # Документация
│   └── rules/               # Правила проекта
├── migrations/              # SQL миграции
├── skills/                  # Навыки (SKILL.md)
└── .github/workflows/      # CI/CD конфигурации
```

## Начало работы

### Требования

- Go 1.25.5 или новее
- Git
- Docker (опционально, для запуска БД)
- Make (опционально)

### Установка и запуск

```bash
# Клонирование репозитория
git clone https://github.com/atumaikin/nexflow.git
cd nexflow

# Установка зависимостей
go mod download

# Запуск тестов
go test ./...

# Запуск сервера
go run cmd/server/main.go
```

### Настройка окружения

Создайте файл `config.yml` или используйте переменные окружения:

```yaml
server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "./data/nexflow.db"
  migrations_path: "./migrations/sqlite"

llm:
  default_provider: "anthropic"
  providers:
    anthropic:
      api_key: "${ANTHROPIC_API_KEY}"
      model: "claude-sonnet-4"

logging:
  level: "info"
  format: "json"
```

## Архитектура

### Слои архитектуры

#### Domain Layer (`internal/domain/`)

Содержит бизнес-логику и сущности, не зависящие от внешней инфраструктуры.

```go
// internal/domain/entity/user.go
type User struct {
    ID        string
    Channel   string
    ChannelID string
    CreatedAt time.Time
}

func NewUser(channel, channelID string) *User {
    return &User{
        ID:        utils.GenerateID(),
        Channel:   channel,
        ChannelID: channelID,
        CreatedAt: utils.Now(),
    }
}
```

#### Application Layer (`internal/application/`)

Содержит use cases и DTO для взаимодействия с внешним миром.

```go
// internal/application/ports/connector.go
type Connector interface {
    Start(ctx context.Context) error
    Stop() error
    Events() <-chan Event
    SendMessage(ctx context.Context, userID, message string) error
    Name() string
}
```

#### Infrastructure Layer (`internal/infrastructure/`)

Содержит реализацию внешних зависимостей: БД, LLM, channels.

```go
// internal/infrastructure/persistence/database/mappers/user_mapper.go
func UserToDomain(dbUser *dbmodel.User) *entity.User {
    if dbUser == nil {
        return nil
    }
    return &entity.User{
        ID:        dbUser.ID,
        Channel:   dbUser.Channel,
        ChannelID: dbUser.ChannelUserID,
        CreatedAt: utils.ParseTimeRFC3339(dbUser.CreatedAt),
    }
}
```

### Направление зависимостей

```
Application → Domain
Infrastructure → Domain
Infrastructure → Application (через ports)
```

## Разработка новых компонентов

### Добавление новой сущности

1. Создайте файл в `internal/domain/entity/`
2. Определите структуру с godoc comments
3. Реализуйте конструктор с использованием `utils.GenerateID()` и `utils.Now()`
4. Добавьте методы валидации
5. Создайте unit тесты

```go
package entity

import (
    "time"

    "github.com/atumaikin/nexflow/internal/shared/utils"
)

// MyEntity represents a new domain entity.
type MyEntity struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

// NewMyEntity creates a new instance of MyEntity.
func NewMyEntity(name string) *MyEntity {
    return &MyEntity{
        ID:        utils.GenerateID(),
        Name:      name,
        CreatedAt: utils.Now(),
    }
}
```

### Добавление нового Repository

1. Создайте интерфейс в `internal/domain/repository/`
2. Реализуйте его в `internal/infrastructure/persistence/database/`
3. Создайте mapper в `internal/infrastructure/persistence/database/mappers/`
4. Добавьте unit тесты с использованием mock

```go
// internal/domain/repository/myentity_repository.go
package repository

import (
    "context"
    "github.com/atumaikin/nexflow/internal/domain/entity"
)

type MyEntityRepository interface {
    Create(ctx context.Context, entity *entity.MyEntity) error
    GetByID(ctx context.Context, id string) (*entity.MyEntity, error)
    Update(ctx context.Context, entity *entity.MyEntity) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context) ([]*entity.MyEntity, error)
}
```

### Добавление нового Use Case

1. Создайте файл в `internal/application/usecase/`
2. Определите структуру с зависимостями
3. Реализуйте бизнес-логику
4. Создайте unit тесты

```go
package usecase

import (
    "context"
    "github.com/atumaikin/nexflow/internal/domain/entity"
    "github.com/atumaikin/nexflow/internal/domain/repository"
)

type MyEntityUseCase struct {
    repo repository.MyEntityRepository
}

func NewMyEntityUseCase(repo repository.MyEntityRepository) *MyEntityUseCase {
    return &MyEntityUseCase{
        repo: repo,
    }
}

func (uc *MyEntityUseCase) Create(ctx context.Context, name string) (*entity.MyEntity, error) {
    ent := entity.NewMyEntity(name)
    if err := uc.repo.Create(ctx, ent); err != nil {
        return nil, fmt.Errorf("failed to create entity: %w", err)
    }
    return ent, nil
}
```

## Тестирование

### Unit тесты

Каждый модуль должен иметь unit тесты в файле `*_test.go`.

```go
package entity

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNewMyEntity(t *testing.T) {
    ent := NewMyEntity("test")

    assert.NotEmpty(t, ent.ID)
    assert.Equal(t, "test", ent.Name)
    assert.NotZero(t, ent.CreatedAt)
}
```

### Таблицарные тесты

Используйте таблицарные тесты для множественных сценариев:

```go
func TestMyEntity_Validate(t *testing.T) {
    tests := []struct {
        name    string
        entity  *MyEntity
        wantErr bool
    }{
        {
            name:    "valid entity",
            entity:  NewMyEntity("valid"),
            wantErr: false,
        },
        {
            name:    "empty name",
            entity:  NewMyEntity(""),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.entity.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Запуск тестов

```bash
# Все тесты
go test ./...

# Тесты с покрытием
go test -cover ./...

# Тесты с race detector
go test -race ./...

# Детальный вывод
go test -v ./...
```

## Логирование

### Использование Logger

```go
import (
    "github.com/atumaikin/nexflow/internal/shared/logging"
)

type MyService struct {
    logger logging.Logger
}

func NewMyService(logger logging.Logger) *MyService {
    return &MyService{
        logger: logger,
    }
}

func (s *MyService) DoSomething(ctx context.Context) {
    s.logger.Info("Starting operation", "component", "MyService")
    // ...
}
```

### Маскирование секретов

Logger автоматически маскирует поля с ключами: `token`, `key`, `password`, `secret`.

```go
s.logger.Info("Connecting to database",
    "host", "localhost",
    "password", "secret123", // Будет замаскировано в логах
)
// Output: {"level":"info","msg":"Connecting to database","host":"localhost","password":"***"}
```

## Конфигурация

### Загрузка конфигурации

```go
import (
    "github.com/atumaikin/nexflow/internal/shared/config"
)

cfg, err := config.Load("config.yml")
if err != nil {
    log.Fatal("Failed to load config:", err)
}

if err := cfg.Validate(); err != nil {
    log.Fatal("Invalid config:", err)
}
```

### Использование переменных окружения

Конфигурация поддерживает подстановку переменных окружения:

```yaml
database:
  password: "${DB_PASSWORD}"  # Будет заменено значением из ENV
```

## База данных

### Использование SQLC

SQLC генерирует типобезопасные Go код из SQL запросов.

```sql
-- migrations/sqlite/0001_init.sql
-- name: CreateUser :one
INSERT INTO users (id, channel, channel_id, created_at)
VALUES (?, ?, ?, ?)
RETURNING *;
```

```go
// Запуск SQLC
sqlc generate

// Использование сгенерированного кода
user, err := db.Queries().CreateUser(ctx, db.CreateUserParams{
    ID:        userID,
    Channel:   channel,
    ChannelID: channelID,
    CreatedAt: time.Now().Format(time.RFC3339),
})
```

### Запуск миграций

```go
import (
    "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
)

db := database.New(cfg.Database)
if err := db.Migrate(ctx); err != nil {
    log.Fatal("Failed to migrate database:", err)
}
```

## Деплой

### Docker

```dockerfile
# Dockerfile
FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o nexflow cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/nexflow .
COPY config.yml .
CMD ["./nexflow"]
```

### Docker Compose

```yaml
version: '3.8'
services:
  nexflow:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
    volumes:
      - ./data:/app/data
```

## Полезные команды

```bash
# Форматирование кода
go fmt ./...

# Линтинг
golangci-lint run

# Генерация зависимостей
go mod tidy

# Генерация документации
go doc ./...

# Сборка для всех платформ
GOOS=linux GOARCH=amd64 go build -o nexflow-linux cmd/server/main.go
GOOS=darwin GOARCH=amd64 go build -o nexflow-macos cmd/server/main.go
GOOS=windows GOARCH=amd64 go build -o nexflow.exe cmd/server/main.go
```

## Ресурсы

- [Godoc](https://pkg.go.dev/github.com/atumaikin/nexflow)
- [Архитектура](./rules/architecture.md)
- [Тестирование](./rules/testing.md)
- [Качество кода](./rules/codequality.md)
- [Безопасность](./rules/security.md)
