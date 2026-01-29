# Logging Package

Структурированная система логирования для проекта Nexflow с поддержкой JSON формата и автоматическим маскированием секретов.

## Возможности

- **Структурированное логирование** на базе `log/slog` (Go 1.21+)
- **JSON формат** для удобства парсинга логов
- **Автоматическое маскирование секретов** (API ключи, токены, пароли и т.д.)
- **Уровни логирования**: DEBUG, INFO, WARN, ERROR, FATAL
- **Контекстное логирование** с поддержкой request_id и других полей
- **HTTP middleware** для автоматического логирования запросов

## Использование

### Базовое использование

```go
import "github.com/atumaikin/nexflow/internal/logging"

// Создание logger с уровнем info и JSON форматом
logger, err := logging.New("info", "json")
if err != nil {
    log.Fatal(err)
}

// Логирование сообщений
logger.Info("Application started", "version", "1.0.0")
logger.Debug("Database connection", "host", "localhost", "port", 5432)
logger.Warn("Rate limit approaching", "current", 90, "max", 100)
logger.Error("Failed to connect", "error", err)
```

### Логирование с секретами (автоматическое маскирование)

```go
// Секреты будут автоматически замаскированы
logger.Info("API request",
    "api_key", "sk-1234567890abcdef",
    "endpoint", "/users",
    "user_id", "12345")
```

Результат в JSON:
```json
{
  "level": "INFO",
  "msg": "API request",
  "api_key": "sk**************ef",
  "endpoint": "/users",
  "user_id": "12345",
  "source": "nexflow",
  "time": "2024-01-29T21:00:00.000Z"
}
```

### Логирование с контекстом

```go
import "context"

// Создание logger с контекстом
ctx := context.Background()
loggerWithContext := logger.WithContext(ctx)

// Логирование с контекстом
loggerWithContext.InfoContext(ctx, "Request processed", "user_id", "12345")
```

### Добавление полей к logger

```go
// Создание logger с постоянными полями
serviceLogger := logger.With("service", "auth", "version", "1.0.0")

// Все сообщения будут включать эти поля
serviceLogger.Info("User logged in", "user_id", "12345")
```

### HTTP Middleware

```go
import (
    "net/http"
    "github.com/atumaikin/nexflow/internal/logging"
)

func main() {
    logger, _ := logging.New("info", "json")

    // Создание middleware
    loggingMiddleware := logging.Middleware(logger)

    // Использование в HTTP server
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    http.ListenAndServe(":8080", loggingMiddleware(mux))
}
```

## Маскирование секретов

Следующие поля автоматически маскируются в логах:

- `api_key`, `apikey`, `apiKey`
- `token`, `access_token`, `accessToken`
- `refresh_token`, `refreshtoken`, `refreshToken`
- `password`, `pass`
- `secret`
- `private_key`, `privatekey`, `privateKey`
- `bot_token`, `bottoken`, `botToken`
- `auth_token`, `authtoken`, `authToken`
- `bearer_token`, `bearertoken`, `bearerToken`
- `client_secret`, `clientsecret`, `clientSecret`
- `api_secret`, `apisecret`, `apiSecret`
- `session_token`, `sessiontoken`, `sessionToken`
- `csrf_token`, `csrftoken`, `csrfToken`
- `authorization`
- `credentials`, `credential`

Маскирование показывает только первые 2 и последние 2 символа значения.

## Уровни логирования

- **DEBUG** - Детальная информация для отладки
- **INFO** - Информационные сообщения о нормальной работе
- **WARN** - Предупреждения о потенциальных проблемах
- **ERROR** - Ошибки, которые не являются критическими
- **FATAL** - Критические ошибки (приводят к завершению приложения)

Уровень логирования настраивается через конфигурационный файл:

```yaml
logging:
  level: "info"  # debug, info, warn, error, fatal
  format: "json" # json, text
```

## Конфигурация через переменные окружения

```bash
export LOGGING_LEVEL=debug
export LOGGING_FORMAT=json
```

## Интеграция с другими пакетами

### В database пакете

```go
import "github.com/atumaikin/nexflow/internal/logging"

type Database struct {
    db     *sql.DB
    logger logging.Logger
}

func NewDatabase(dbPath string, logger logging.Logger) (*Database, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        logger.Error("Failed to open database", "error", err, "path", dbPath)
        return nil, err
    }

    logger.Info("Database connected", "path", dbPath)
    return &Database{db: db, logger: logger}, nil
}
```

### В handlers пакете

```go
type Handler struct {
    logger logging.Logger
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("Processing request", "path", r.URL.Path, "method", r.Method)

    // ... обработка запроса

    h.logger.Info("Request completed", "status", "200")
}
```

## Тестирование

Запуск тестов:

```bash
go test ./internal/logging/...
```

Запуск тестов с покрытием:

```bash
go test -cover ./internal/logging/...
```

## Производительность

- Использование `log/slog` обеспечивает высокую производительность
- Логирование асинхронно через стандартный output stream
- Маскирование секретов выполняется при записи, без влияния на бизнес-логику

## Лучшие практики

1. **Используйте правильные уровни логирования**:
   - DEBUG: Для детальной отладки
   - INFO: Для нормальной работы приложения
   - WARN: Для потенциальных проблем
   - ERROR: Для ошибок, которые можно восстановить
   - FATAL: Для критических ошибок

2. **Добавляйте контекстную информацию**:
   ```go
   logger.Info("User action", "user_id", "12345", "action", "login", "ip", "192.168.1.1")
   ```

3. **Не логируйте секреты напрямую**:
   ```go
   // Плохо
   logger.Info("API call", "credentials", "user:password")

   // Хорошо (автоматически маскируется)
   logger.Info("API call", "api_key", "sk-1234567890")
   ```

4. **Используйте структурированные данные**:
   ```go
   // Хорошо
   logger.Info("User created", "user_id", "12345", "email", "user@example.com", "role", "admin")
   ```

## Требования

- Go 1.21 или выше (для `log/slog`)
