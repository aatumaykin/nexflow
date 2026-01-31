# Nexflow

![codecov](https://codecov.io/gh/aatumaykin/nexflow/branch/main/graph/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/aatumaykin/nexflow)](https://goreportcard.com/report/github.com/aatumaykin/nexflow)

Self-hosted ИИ-агент на Go, управляющий цифровыми потоками задач через multiple channels (Telegram, Discord, Web UI) с LLM-провайдерами (Anthropic, OpenAI, Ollama и др.) и навыками (skills).

## 🚀 Быстрый старт

```bash
# Клонировать репозиторий
git clone https://github.com/your-repo/nexflow.git
cd nexflow

# Установить зависимости
go mod download

# Запустить сервер
go run cmd/server/main.go
```

## 📋 Требования

- Go 1.25.5 или выше
- SQLite 3 (по умолчанию) или PostgreSQL
- Доступ к LLM провайдеру (API ключ)

## 🏗️ Архитектура

Nexflow следует принципам чистой слоистой архитектуры (Clean Layered Architecture):

```
┌─────────────────────────────────────────────────────────┐
│                   Presentation Layer                  │
│  (API, Telegram Bot, Discord Bot, Web UI)      │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                Application Layer                  │
│          (Use Cases, DTOs, Ports)               │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                   Domain Layer                   │
│      (Entities, Value Objects, Repositories)    │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│              Infrastructure Layer                 │
│    (Database, LLM Providers, Channels, Skills)  │
└──────────────────────────────────────────────────────┘
```

### Слои

- **Domain Layer** (`internal/domain/`) - Бизнес-логика и сущности
- **Application Layer** (`internal/application/`) - Use cases и orchestration
- **Infrastructure Layer** (`internal/infrastructure/`) - Внешние зависимости
- **Presentation Layer** (`cmd/`) - Entry points и API

## 📁 Структура проекта

```
nexflow/
├── cmd/                    # Entry points (main.go)
│   └── server/
├── internal/
│   ├── application/         # Application layer
│   │   ├── dto/          # Data Transfer Objects
│   │   ├── ports/        # Interfaces for external dependencies
│   │   └── usecase/      # Business logic (use cases)
│   ├── domain/            # Domain layer
│   │   ├── entity/       # Business entities
│   │   └── repository/   # Repository interfaces
│   ├── infrastructure/    # Infrastructure layer
│   │   └── persistence/   # Database implementations
│   └── shared/           # Shared utilities
├── docs/                 # Documentation
├── skills/               # Skill definitions
└── migrations/           # Database migrations
```

## 🧪 Тестирование

```bash
# Запустить все тесты
go test ./...

# Запустить с покрытием
go test -cover ./...

# Запустить с race detection
go test -race ./...

# Запустить только unit тесты
go test ./internal/domain/... ./internal/application/...

# Запустить только integration тесты
go test ./internal/infrastructure/...

# Сгенерировать HTML отчет покрытия
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out -o coverage.html

# Проверить покрытие по функциям
go tool cover -func=coverage.out

# Показать общее покрытие
go tool cover -func=coverage.out | tail -1
```

### 📊 Покрытие кода

- **Целевое покрытие:** 60%
- **Текущее покрытие:** обновляется на каждом PR
- **Тренд:** отслеживается через [Codecov](https://codecov.io/gh/aatumaykin/nexflow)

Детальные отчеты доступны в CI/CD artifacts и на странице [Codecov](https://codecov.io/gh/aatumaykin/nexflow).

## 🔧 Конфигурация

Конфигурация загружается из YAML/JSON файла или переменных окружения:

```yaml
server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "./data/nexflow.db"

llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
      model: "gpt-4"

channels:
  telegram:
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users: []
```

## 📚 Документация

- [Development Guide](docs/development-guide.md)
- [Testing Guide](docs/testing-guide.md)
- [Refactoring Guide](docs/refactoring-guide.md)
- [API Reference](docs/api-reference.md)
- [Database Configuration](docs/database-config.md)
- [PRD v2](docs/nexflow-prd-v2.md)
- [Implementation Plan](docs/nexflow-implementation-plan.md)

## 🤝 Contributing

1. Fork проекта
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📝 Лицензия

MIT License - см. файл LICENSE для деталей

## 👥 Поддержка

- GitHub Issues: https://github.com/your-repo/nexflow/issues
- Документация: https://github.com/your-repo/nexflow/tree/main/docs
