# Database Package

Пакет `database` обеспечивает доступ к базе данных для проекта Nexflow. Он поддерживает SQLite (для разработки) и PostgreSQL (для продакшена).

## Установка

```bash
go get github.com/atumaikin/nexflow/internal/database
```

## Использование

### Инициализация

```go
import (
    "context"
    "github.com/atumaikin/nexflow/internal/config"
    "github.com/atumaikin/nexflow/internal/database"
)

// Загрузка конфигурации
cfg, err := config.Load("config.yaml")
if err != nil {
    log.Fatal(err)
}

// Создание подключения к базе данных
// По умолчанию используется NoopLogger (не логирует)
db, err := database.NewDatabase(&cfg.Database)
if err != nil {
    log.Fatal(err)
}

// Или с кастомным логгером
// logger, err := logging.New("info", "json")
// if err != nil {
//     log.Fatal(err)
// }
// db, err := database.NewDatabase(&cfg.Database, database.WithLogger(logger))
// if err != nil {
//     log.Fatal(err)
// }
defer db.Close()

// Запуск миграций
if err := db.Migrate(context.Background()); err != nil {
    log.Fatal(err)
}
```

### Примеры использования

#### Создание пользователя

```go
userID := uuid.New().String()
params := database.CreateUserParams{
    ID:            userID,
    Channel:       "telegram",
    ChannelUserID: "123456",
    CreatedAt:     time.Now().Format(time.RFC3339),
}

user, err := db.CreateUser(ctx, params)
if err != nil {
    log.Fatal(err)
}
```

#### Получение пользователя

```go
// По ID
user, err := db.GetUserByID(ctx, userID)
if err != nil {
    log.Fatal(err)
}

// По каналу
params := database.GetUserByChannelParams{
    Channel:       "telegram",
    ChannelUserID: "123456",
}
user, err := db.GetUserByChannel(ctx, params)
if err != nil {
    log.Fatal(err)
}
```

#### Создание сессии

```go
sessionID := uuid.New().String()
now := time.Now().Format(time.RFC3339)
params := database.CreateSessionParams{
    ID:        sessionID,
    UserID:    userID,
    CreatedAt: now,
    UpdatedAt: now,
}

session, err := db.CreateSession(ctx, params)
if err != nil {
    log.Fatal(err)
}
```

#### Создание сообщения

```go
messageID := uuid.New().String()
params := database.CreateMessageParams{
    ID:        messageID,
    SessionID: sessionID,
    Role:      "user",
    Content:   "Hello, world!",
    CreatedAt: now,
}

message, err := db.CreateMessage(ctx, params)
if err != nil {
    log.Fatal(err)
}
```

#### Создание задачи

```go
taskID := uuid.New().String()
params := database.CreateTaskParams{
    ID:        taskID,
    SessionID: sessionID,
    Skill:     "my-skill",
    Input:     "input data",
    Status:    "pending",
    CreatedAt: now,
    UpdatedAt: now,
}

task, err := db.CreateTask(ctx, params)
if err != nil {
    log.Fatal(err)
}
```

#### Обновление задачи

```go
params := database.UpdateTaskParams{
    Output:    sql.NullString{String: "result", Valid: true},
    Status:    "completed",
    Error:     sql.NullString{Valid: false},
    UpdatedAt: now,
    ID:        taskID,
}

task, err := db.UpdateTask(ctx, params)
if err != nil {
    log.Fatal(err)
}
```

## Миграции

Миграции находятся в директории `migrations/`:

- `001_create_schema.up.sql` - создание схемы для SQLite
- `001_create_schema.down.sql` - откат схемы для SQLite
- `001_create_schema.postgres.up.sql` - создание схемы для PostgreSQL
- `001_create_schema.postgres.down.sql` - откат схемы для PostgreSQL

### Ручной запуск миграций

```go
// Применить миграции
if err := db.Migrate(ctx); err != nil {
    log.Fatal(err)
}

// Откатить последнюю миграцию
if err := db.Rollback(ctx); err != nil {
    log.Fatal(err)
}
```

## Структура базы данных

### Users
- `id` - уникальный идентификатор пользователя
- `channel` - канал связи (telegram, web, etc.)
- `channel_user_id` - идентификатор пользователя в канале
- `created_at` - время создания

### Sessions
- `id` - уникальный идентификатор сессии
- `user_id` - ссылка на пользователя
- `created_at` - время создания
- `updated_at` - время последнего обновления

### Messages
- `id` - уникальный идентификатор сообщения
- `session_id` - ссылка на сессию
- `role` - роль отправителя (user, assistant, system)
- `content` - содержимое сообщения
- `created_at` - время создания

### Tasks
- `id` - уникальный идентификатор задачи
- `session_id` - ссылка на сессию
- `skill` - название скилла для выполнения
- `input` - входные данные
- `output` - результат выполнения
- `status` - статус (pending, in_progress, completed, failed)
- `error` - ошибка (если есть)
- `created_at` - время создания
- `updated_at` - время последнего обновления

### Skills
- `id` - уникальный идентификатор скилла
- `name` - название скилла (уникальное)
- `version` - версия
- `location` - путь к файлу скилла
- `permissions` - необходимые разрешения (JSON)
- `metadata` - метаданные (JSON)
- `created_at` - время создания

### Schedules
- `id` - уникальный идентификатор расписания
- `skill` - ссылка на скилл
- `cron_expression` - cron выражение
- `input` - входные данные
- `enabled` - включено ли расписание
- `created_at` - время создания

### Logs
- `id` - уникальный идентификатор лога
- `level` - уровень (debug, info, warn, error)
- `source` - источник лога
- `message` - сообщение
- `metadata` - дополнительные данные (JSON)
- `created_at` - время создания

## Генерация кода

Код генерируется с помощью [sqlc](https://sqlc.dev/):

```bash
# Установка sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Генерация кода
cd internal/database
sqlc generate
```

После изменения SQL-файлов (schema.sql или query.sql) необходимо перегенерировать код.

## Тестирование

```bash
# Запуск всех тестов
go test ./internal/database/...

# Запуск с покрытием
go test -cover ./internal/database/...
```

## Поддерживаемые базы данных

- **SQLite** - для локальной разработки
- **PostgreSQL** - для продакшена
