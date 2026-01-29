# Обзор архитектуры Nexflow

Документ описывает общую архитектуру проекта Nexflow, ключевые компоненты и их взаимодействие.

## Введение

Nexflow - self-hosted ИИ-агент на Go, который управляет цифровыми потоками задач через multiple channels (Telegram, Discord, Web UI) с использованием LLM-провайдеров и системы навыков (skills).

## Архитектурные принципы

Проект следует следующим принципам:

1. **Clean Layered Architecture** - четкое разделение на слои
2. **Domain-Driven Design** - бизнес-логика в domain слое
3. **Dependency Inversion** - зависимости только от интерфейсов
4. **Single Responsibility** - каждый компонент имеет одну ответственность
5. **Testability** - все слои легко тестируемы

## Высокоуровневая архитектура

```
┌─────────────────────────────────────────────────────────┐
│                  Presentation Layer                   │
│  API Endpoints | Telegram Bot | Discord Bot | Web UI  │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP/WebSocket
┌──────────────────────▼──────────────────────────────┐
│                Application Layer                  │
│        Use Cases (Orchestration)               │
│    UserUseCase | ChatUseCase | ...             │
└──────────────────────┬──────────────────────────────┘
                       │ Domain entities
┌──────────────────────▼──────────────────────────────┐
│                   Domain Layer                   │
│        Entities: User, Session, Message...          │
│     Repository Interfaces (Ports)                  │
└──────────────────────┬──────────────────────────────┘
                       │ Repository implementations
┌──────────────────────▼──────────────────────────────┐
│              Infrastructure Layer                 │
│         SQLite/PostgreSQL | LLM Providers        │
│         Channels | Skill Runtime                   │
└──────────────────────────────────────────────────────┘
```

## Ключевые компоненты

### Presentation Layer

**Расположение:** `cmd/`

**Ответственность:**
- Entry points приложения
- HTTP API endpoints
- WebSocket connections
- Telegram/Discord bots

**Компоненты:**
- `cmd/server/` - HTTP сервер
- `cmd/telegram-bot/` - Telegram bot
- `cmd/discord-bot/` - Discord bot (TBD)

### Application Layer

**Расположение:** `internal/application/`

**Ответственность:**
- Orchestration бизнес-процессов
- Use case implementation
- DTO преобразование

**Компоненты:**
- `usecase/user_usecase.go` - управление пользователями
- `usecase/chat_usecase.go` - чат-функциональность
- `dto/` - Data Transfer Objects
- `ports/` - интерфейсы внешних зависимостей

### Domain Layer

**Расположение:** `internal/domain/`

**Ответственность:**
- Бизнес-сущности
- Бизнес-логика
- Repository interfaces

**Компоненты:**
- `entity/user.go` - пользователь
- `entity/session.go` - сессия
- `entity/message.go` - сообщение
- `entity/task.go` - задача
- `entity/skill.go` - навык
- `entity/schedule.go` - расписание
- `entity/log.go` - лог
- `repository/*` - репозитории (интерфейсы)

### Infrastructure Layer

**Расположение:** `internal/infrastructure/`

**Ответственность:**
- Реализация репозиториев
- Интеграции с внешними сервисами
- База данных

**Компоненты:**
- `persistence/database/sqlite/` - SQLite репозитории
- `channels/` - интеграции с Telegram, Discord
- `llm/` - LLM провайдеры (OpenAI, Anthropic, Ollama)
- `skills/` - runtime для навыков

## Потоки данных

### User Creation Flow

```
User Request (Telegram/Web)
    ↓
UserUseCase.CreateUser()
    ↓
UserRepository.Create()
    ↓
SQLite Database
    ↓
User Response
```

### Chat Flow

```
User Message (Telegram/Web)
    ↓
ChatUseCase.SendMessage()
    ↓
1. FindOrCreateUser()
2. CreateSession()
3. SaveMessage() [user]
4. GetConversationHistory()
    ↓
LLMProvider.Generate()
    ↓
5. SaveMessage() [assistant]
6. UpdateSession()
    ↓
AI Response
```

### Skill Execution Flow

```
AI Request → Execute Skill
    ↓
ChatUseCase.ExecuteSkill()
    ↓
TaskRepository.Create()
    ↓
Task.SetRunning()
    ↓
SkillRuntime.Execute()
    ↓
Task.SetCompleted()
    ↓
TaskRepository.Update()
    ↓
Result to AI
```

## Технологический стек

### Backend
- **Язык:** Go 1.25.5
- **База данных:** SQLite 3 / PostgreSQL
- **ORM:** SQLC (type-safe SQL)
- **Логирование:** slog (structured logging)

### LLM Integration
- **Провайдеры:** OpenAI, Anthropic, Ollama, Custom
- **Маршрутизация:** динамическая по названию провайдера
- **Формат:** Unified Completion API

### Channels
- **Telegram:** Bot API
- **Discord:** Bot API (TBD)
- **Web:** WebSocket + REST API

## Безопасность

1. **Секреты:** через ENV переменные
2. **Маскирование:** секретов в логах
3. **Валидация:** всех входных данных
4. **Sandbox:** для опасных навыков
5. **Whitelist:** разрешенных ресурсов

## Масштабирование

### Горизонтальное масштабирование
- Несколько инстансов Nexflow
- Load balancing через nginx/haproxy
- Shared database (PostgreSQL)

### Вертикальное масштабирование
- Оптимизация БД queries
- Кэширование (Redis TBD)
- Асинхронная обработка задач (Worker pool TBD)

## Мониторинг

- **Логирование:** структурированные логи
- **База данных:** query logs
- **Приложение:** health checks
- **Tasks:** статусы выполнения

## Документация

- [README.md](../README.md) - быстрый старт
- [Layer Architecture](./layer-architecture.md) - детальная слоистая архитектура
- [Migration Guide](./MIGRATION.md) - руководство по миграции
- [Development Guide](./development-guide.md) - разработка
- [Testing Guide](./testing-guide.md) - тестирование
- [API Reference](./api-reference.md) - API endpoints

## Заключение

Архитектура Nexflow спроектирована для:
- **Масштабируемости:** горизонтальное и вертикальное
- **Тестируемости:** каждый слой легко тестируется
- **Поддерживаемости:** простое добавление новых каналов и LLM провайдеров
- **Безопасности:** встроенные механизмы безопасности
- **Производительности:** эффективная работа с БД и кэширование

---

**Последнее обновление:** Январь 2026
