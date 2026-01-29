# Архитектура проекта Nexflow

Этот модуль описывает архитектурную модель, слои и зависимости.

## Архитектурный паттерн

**Layered Architecture** с элементами **Clean Architecture**.

### High-level архитектура

```
┌─────────────────────────────────────────────────────────────┐
│                      Channels Layer                         │
│  Telegram  │  Discord  │  Web UI  │  Email  │  Webhooks     │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                    Core Gateway (Go)                        │
│  Message Router │ Orchestrator │ LLM Router │ Skills         │
└─────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│                   Storage & Execution                       │
│  SQLite/Postgres │  FS (Markdown) │  Sandbox Containers   │
└─────────────────────────────────────────────────────────────┘
```

## Основные слои

### 1. Channels Layer (`internal/channels/`)

**Интерфейс:**
```go
type Connector interface {
    Start(ctx context.Context) error
    Stop() error
    Events() <-chan Event
    SendMessage(ctx context.Context, userID, message string) error
}
```

**Правила:**
- Все коннекторы реализуют единый интерфейс
- Входящие события → `Events()` channel
- Без бизнес-логики — только адаптация
- Не используют БД и LLM напрямую

### 2. Core Gateway (`cmd/server/`)

**Компоненты:**
- **Message Router:** события от Connectors → Orchestrator
- **Orchestrator:** управление потоком, выбор LLM, выполнение навыков
- **LLM Router:** выбор модели по политикам, балансировка

**Правила:**
- Минимальная логика в main.go — только инициализация
- Использует интерфейсы для зависимостей

### 3. Skills Layer (`internal/skills/` + `skills/`)

**Интерфейс:**
```go
type Skill interface {
    Name() string
    Execute(ctx context.Context, input string) (string, error)
    RequiresSandbox() bool
    Permissions() []string
}
```

**Формат:** директория с `SKILL.md` + скрипты

**Правила:**
- Навыки изолированы
- Таймауты для выполнения
- Sandbox для опасных permissions

### 4. Storage Layer (`internal/database/`)

**Интерфейс:**
```go
type Database interface {
    CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
    GetUserByID(ctx context.Context, id string) (User, error)
    // ... другие методы
    Migrate(ctx context.Context) error
    Close() error
}
```

**Правила:**
- SQLC для генерации типов
- Интерфейс для мокинга
- Connection pool: 25 max open/idle, 5min lifetime
- Не зависит от Channels и LLM

### 5. LLM Layer (`internal/llm/`)

**Интерфейс:**
```go
type LLMProvider interface {
    CreateCompletion(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
    CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)
}
```

**Провайдеры:** Anthropic, OpenAI, Google Gemini, z.ai, OpenRouter, Ollama, Custom

**Правила:**
- Единый интерфейс для всех провайдеров
- Routing по политикам
- Retry логика

### 6. Config Layer (`internal/config/`)

```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    LLM      LLMConfig
    Channels ChannelsConfig
    Skills   SkillsConfig
    Logging  LoggingConfig
}
```

**Правила:**
- YAML/JSON с `${VAR_NAME}` подстановкой
- Валидация после загрузки
- Секреты через ENV

### 7. Logging Layer (`internal/logging/`)

```go
type Logger interface {
    Info(msg string, args ...any)
    Error(msg string, args ...any)
    With(args ...any) Logger
    WithContext(ctx context.Context) Logger
}
```

**Правила:**
- slog (structured logging)
- JSON/text формат
- Маскирование секретов

## Направление зависимостей

### Допустимые зависимости

```
Channels → Config, Logging
    ↓
Core Gateway → Channels, LLM, Skills, Database
    ↓
Skills → Config, Logging
    ↓
LLM → Config, Logging
    ↓
Database → Config, Logging
```

### Принципы

1. **Верхние → нижние:** Connector → Config
2. **Нижние ⊄ верхние:** Database ⊄ Channels
3. **Пакеты в internal/** изолированы

### Правильные зависимости

```go
// ✅ ПРАВИЛЬНО: Connector может зависеть от Config
type Connector struct {
    config  *config.TelegramConfig
    logger  logging.Logger
}

// ✅ ПРАВИЛЬНО: main может координировать все слои
import (
    "github.com/atumaikin/nexflow/internal/channels"
    "github.com/atumaikin/nexflow/internal/database"
    "github.com/atumaikin/nexflow/internal/llm"
)
```

### Запрещённые зависимости

```go
// ❌ НЕПРАВИЛЬНО: Database ⊄ Channels
import "github.com/atumaikin/nexflow/internal/channels"

// ❌ НЕПРАВИЛЬНО: LLM ⊄ Database
import "github.com/atumaikin/nexflow/internal/database"
```

## Пакеты в internal/

| Пакет | Назначение | Зависимости | Используется |
|-------|-----------|-------------|--------------|
| config | Конфигурация | нет | всеми |
| database | БД | config, logging | Core Gateway |
| logging | Логирование | нет | всеми |
| channels | Connectors | config, logging | Core Gateway |
| llm | LLM провайдеры | config, logging | Core Gateway |
| skills | Skills execution | config, logging | Core Gateway |

## Пакеты в pkg/

Публичные библиотеки, API стабильно, документировано, не зависит от `internal/`.

## Циклические зависимости

**НИКОГДА не создавайте циклические зависимости.**

Для проверки:
```bash
go mod graph | grep -v std | grep -v "nexflow"
```

Если обнаружили:
1. Выделите общую логику в отдельный пакет
2. Используйте интерфейсы для разрыва связей
3. Переместите зависимость в более верхний слой

## SOLID в архитектуре

### Single Responsibility
Каждый пакет — одна ответственность:
- `internal/config` — только конфигурация
- `internal/database` — только БД

### Open/Closed
Новые каналы через реализацию интерфейса `Connector`.

### Liskov Substitution
Все реализации интерфейсов взаимозаменяемы.

### Interface Segregation
Интерфейсы минимальны и специфичны.

### Dependency Inversion
Высокоуровневые модули зависят от абстракций:
```go
type Orchestrator struct {
    llm     llm.LLMRouter
    skills  skills.SkillRegistry
    db      database.Database
}
```

## Создание новых пакетов

### В internal/

Новый пакет → `internal/package-name/`

### В pkg/

Новый пакет → `pkg/package-name/`
- API стабильно
- Не зависит от `internal/`
- Тестировано и документировано

## Критические архитектурные правила

1. НИКОГДА не создавайте циклические зависимости
2. ВСЕГДА используйте интерфейсы для тестирования
3. НИКОГДА не позволяйте нижним слоям зависеть от верхних
4. ВСЕГДА разделяйте логику по слоям
5. НИКОГДА не смешивайте ответственность слоёв
6. ВСЕГДА следуйте принципу зависимости вниз
7. НИКОГДА не импортируйте `internal/` в `pkg/`
8. ВСЕГДА используйте `context.Context` для cross-layer операций
9. НИКОГДА не нарушайте SOLID
10. ВСЕГДА минимизируйте зависимости

---

**Памятка:** Архитектура — это руководство, не догма. Изменения требуют обсуждения.
