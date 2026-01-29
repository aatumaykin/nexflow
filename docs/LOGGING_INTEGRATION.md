# Интеграция логирования в проект

Этот документ описывает, как использовать систему логирования в проекте Nexflow.

## Обзор

Проект использует структурированную систему логирования на базе `log/slog` (Go 1.21+) с поддержкой:
- JSON формата для логов
- Автоматического маскирования секретов
- Контекстного логирования
- Уровней логирования (DEBUG, INFO, WARN, ERROR)
- **NoopLogger** - логгер-пустышка для тестов и случаев, когда логирование не требуется

## Структура

```
internal/logging/
├── logger.go           # Основной интерфейс Logger и реализация SlogLogger
├── secret_mask.go      # Маскирование секретов (API ключи, токены, пароли)
├── middleware.go       # HTTP middleware для логирования запросов
├── logger_test.go      # Unit-тесты
└── README.md           # Документация пакета
```

## NoopLogger (Logger-пустышка)

`NoopLogger` - это реализация интерфейса `Logger`, которая ничего не делает. Полезна для:
- **Тестов** - чтобы не захламлять вывод тестов логами
- **Компонентов, где логирование не нужно** - по умолчанию в конструкторах
- **Снижения накладных расходов** - когда логирование отключено

```go
import "github.com/atumaikin/nexflow/internal/logging"

// Создать NoopLogger
noopLogger := logging.NewNoopLogger()

// Использование такое же, как и обычного logger
noopLogger.Info("Это сообщение не будет выведено")
noopLogger.Error("Это тоже не будет выведено")

// Методы возвращают сам себя (chainable)
logger := noopLogger.With("component", "test")
logger.InfoContext(ctx, "С сообщением")
```

### Использование в конструкторах

Большинство конструкторов в проекте по умолчанию используют `NoopLogger`. Для включения логирования передайте опцию `WithLogger`:

```go
// Без логирования (по умолчанию)
db, err := database.NewDatabase(&cfg.Database)

// С логированием
logger, _ := logging.New("info", "json")
db, err := database.NewDatabase(&cfg.Database, database.WithLogger(logger))
```

## Использование в новых компонентах

### 1. Добавление logger в структуру

```go
import "github.com/atumaikin/nexflow/internal/logging"

type MyService struct {
    logger logging.Logger
    // ... другие поля
}

func NewService(logger logging.Logger) *MyService {
    return &MyService{
        logger: logger,
    }
}
```

### 2. Логирование операций

```go
// Информационное сообщение
service.logger.Info("Service started", "config", cfg)

// Отладочное сообщение
service.logger.Debug("Processing item", "item_id", "12345")

// Предупреждение
service.logger.Warn("Rate limit approaching", "current", 90, "max", 100)

// Ошибка
service.logger.Error("Failed to process", "error", err, "item_id", "12345")
```

### 3. Логирование с контекстом

```go
import "context"

func (s *MyService) Process(ctx context.Context, data string) error {
    s.logger.InfoContext(ctx, "Processing request", "data", data)
    // ... логика
    return nil
}
```

### 4. Добавление постоянных полей

```go
// Создать logger с постоянными полями
serviceLogger := service.logger.With(
    "service", "my-service",
    "version", "1.0.0",
)

// Все сообщения будут включать эти поля
serviceLogger.Info("Processing started")
```

## Маскирование секретов

Секреты маскируются автоматически при логировании. Следующие поля будут замаскированы:

- `api_key`, `apikey`, `apiKey`
- `token`, `access_token`, `accessToken`
- `password`, `pass`
- `secret`
- `bot_token`, `botToken`
- И другие (см. `internal/logging/secret_mask.go`)

Пример:
```go
// Секрет будет автоматически замаскирован
logger.Info("API request",
    "api_key", "sk-1234567890abcdef",  // -> sk***************ef
    "endpoint", "/users",
    "user_id", "12345")
```

## Интеграция с существующими пакетами

### Database пакет

Logger интегрирован в `internal/database/database.go`:

```go
import (
    "github.com/atumaikin/nexflow/internal/database"
    "github.com/atumaikin/nexflow/internal/logging"
)

logger, _ := logging.New("info", "json")
// По умолчанию используется NoopLogger, но можно передать свой
db, err := database.NewDatabase(&cfg.Database, database.WithLogger(logger))
```

Database пакет логирует:
- Подключение к базе данных
- Ошибки подключения
- Закрытие соединения

### HTTP Middleware

Для HTTP серверов используйте middleware из `internal/logging/middleware.go`:

```go
import (
    "github.com/atumaikin/nexflow/internal/logging"
    "net/http"
)

func main() {
    logger, _ := logging.New("info", "json")
    // По умолчанию использует NoopLogger, но можно передать свой
    loggingMiddleware := logging.Middleware(logging.WithMiddlewareLogger(logger))

    mux := http.NewServeMux()
    // ... добавление handlers

    http.ListenAndServe(":8080", loggingMiddleware(mux))
}
```

Middleware логирует:
- Начало запроса (method, path, remote_addr)
- Завершение запроса (status, duration_ms)

## Конфигурация

Логирование настраивается через конфигурационный файл:

```yaml
logging:
  level: "info"  # debug, info, warn, error, fatal
  format: "json" # json, text
```

Или через переменные окружения:

```bash
export LOGGING_LEVEL=debug
export LOGGING_FORMAT=json
```

## Best Practices

### 1. Выбор уровня логирования

- **DEBUG**: Детальная информация для отладки (только в dev)
- **INFO**: Нормальная работа приложения (по умолчанию)
- **WARN**: Потенциальные проблемы, не критичные
- **ERROR**: Ошибки, которые можно восстановить
- **FATAL**: Критические ошибки (приводят к завершению)

### 2. Добавление контекста

```go
// Хорошо: добавить контекстную информацию
logger.Info("User logged in",
    "user_id", "12345",
    "ip", "192.168.1.1",
    "method", "oauth")

// Плохо: недостаточно контекста
logger.Info("User logged in")
```

### 3. Структурированные данные

```go
// Хорошо: использовать структурированные данные
logger.Info("Order created",
    "order_id", "ORD-001",
    "user_id", "12345",
    "total", 99.99,
    "items", 3,
    "status", "pending")

// Плохо: всё в одном сообщении
logger.Info("Order created ORD-001 by user 12345 for 99.99 with 3 items pending")
```

### 4. Обработка ошибок

```go
// Хорошо: передавать ошибку как отдельный параметр
if err != nil {
    logger.Error("Failed to process order",
        "error", err,
        "order_id", "ORD-001",
        "user_id", "12345")
    return err
}

// Плохо: форматировать ошибку в сообщении
if err != nil {
    logger.Error("Failed to process order ORD-001: " + err.Error())
    return err
}
```

## Примеры

### Пример 1: Новый сервис

```go
package service

import (
    "context"
    "github.com/atumaikin/nexflow/internal/logging"
)

type UserService struct {
    logger logging.Logger
    // ... другие зависимости
}

func NewUserService(logger logging.Logger) *UserService {
    return &UserService{
        logger: logger.With("service", "user-service"),
    }
}

func (s *UserService) CreateUser(ctx context.Context, email string) (string, error) {
    s.logger.Info("Creating user", "email", email)

    // ... логика создания пользователя

    s.logger.Info("User created",
        "user_id", userID,
        "email", email)

    return userID, nil
}
```

### Пример 2: HTTP Handler

```go
package handlers

import (
    "net/http"
    "github.com/atumaikin/nexflow/internal/logging"
)

type Handler struct {
    logger logging.Logger
}

func NewHandler(logger logging.Logger) *Handler {
    return &Handler{
        logger: logger.With("component", "handler"),
    }
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")

    h.logger.Info("Fetching user", "user_id", userID)

    // ... логика получения пользователя

    h.logger.Info("User fetched successfully", "user_id", userID)
}
```

### Пример 3: Background Worker

```go
package worker

import (
    "context"
    "time"
    "github.com/atumaikin/nexflow/internal/logging"
)

type Worker struct {
    logger logging.Logger
}

func NewWorker(logger logging.Logger) *Worker {
    return &Worker{
        logger: logger.With("component", "worker"),
    }
}

func (w *Worker) Start(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    w.logger.Info("Worker started")

    for {
        select {
        case <-ctx.Done():
            w.logger.Info("Worker stopped")
            return
        case <-ticker.C:
            w.logger.Debug("Processing batch", "time", time.Now())
            // ... обработка
        }
    }
}
```

## Тестирование

Для тестирования с логированием используйте `NoopLogger`, который ничего не выводит:

```go
import (
    "testing"
    "github.com/atumaikin/nexflow/internal/logging"
)

func TestMyFunction(t *testing.T) {
    logger := logging.NewNoopLogger()

    result, err := MyFunction(logger)
    if err != nil {
        t.Fatal(err)
    }

    // ... проверки
}
```

Если нужно видеть логи при отладке тестов, можно использовать обычный logger:

```go
func TestMyFunctionWithLogging(t *testing.T) {
    logger, _ := logging.New("debug", "text")

    result, err := MyFunction(logger)
    if err != nil {
        t.Fatal(err)
    }

    // ... проверки
}
```

## Мониторинг

Логи в JSON формате можно легко парсить и отправлять в системы мониторинга:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Grafana Loki
- CloudWatch Logs
- Datadog
- И другие системы агрегации логов

Пример запроса для ELK:

```json
{
  "query": {
    "bool": {
      "must": [
        { "match": { "level": "ERROR" } },
        { "match": { "source": "nexflow" } }
      ]
    }
  }
}
```

## Дополнительные ресурсы

- [log/slog документация](https://pkg.go.dev/log/slog)
- [Структурированное логирование в Go](https://go.dev/blog/slog)
- [Логирование best practices](https://github.com/gorilla/websocket/issues/747)
