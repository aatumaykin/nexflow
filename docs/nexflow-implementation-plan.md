# Nexflow Implementation Plan

## Обзор

Документ описывает план реализации Nexflow — self-hosted персонального ИИ-агента. План разделен на фазы: MVP (4-6 недель) и v1.0 (8-12 недель).

## Стек технологий

- **Ядро:** Go 1.22+
- **БД:** SQLite (по умолчанию), Postgres (опционально)
- **Frontend:** Svelte
- **LLM:** Anthropic (Claude), OpenAI, Ollama, Google Gemini, z.ai, OpenRouter + кастомный провайдер
- **Конфигурация:** YAML + JSON
- **Навыки:** Bash, Python, Node.js
- **Деплой:** Docker/Docker Compose

---

## Фаза 1: Базовая инфраструктура (1 неделя)

### Задачи

#### 1.1 Проектная структура
- [ ] Создать Go проект с модулями
- [ ] Настроить структуру директорий: `cmd/`, `internal/`, `pkg/`, `skills/`, `docs/`
- [ ] Инициализировать Go modules
- [ ] Настроить CI/CD (GitHub Actions)

#### 1.2 Конфигурация
- [ ] Создать структуру конфигурации (struct)
- [ ] Реализовать парсер YAML
- [ ] Реализовать парсер JSON
- [ ] Создать пример `config.yml` и `config.json`
- [ ] Реализовать загрузку из ENV переменных

**Пример config.yml:**
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
    anthropic:
      api_key: "${ANTHROPIC_API_KEY}"
      model: "claude-opus-4"
    openai:
      api_key: "${OPENAI_API_KEY}"
      model: "gpt-4"
    ollama:
      base_url: "http://localhost:11434"
      model: "llama3"
    gemini:
      api_key: "${GEMINI_API_KEY}"
      model: "gemini-pro"
    zai:
      api_key: "${ZAI_API_KEY}"
      model: "glm-4"
    openrouter:
      api_key: "${OPENROUTER_API_KEY}"
      model: "anthropic/claude-3-opus"
    custom:
      base_url: "${CUSTOM_LLM_URL}"
      api_key: "${CUSTOM_LLM_KEY}"
      model: "custom-model"

channels:
  telegram:
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users: []
  web:
    enabled: true

skills:
  directory: "./skills"
  timeout_sec: 30
  sandbox_enabled: true

logging:
  level: "info"
  format: "json"
```

#### 1.3 База данных
- [ ] Создать SQLite схему (таблицы: users, sessions, messages, tasks, skills, schedules, logs)
- [ ] Создать Postgres схему (для прод)
- [ ] Реализовать Go ORM (ent/sqlc)
- [ ] Создать миграции

**Схема SQLite:**
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    channel TEXT NOT NULL,
    channel_user_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(channel, channel_user_id)
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);

CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    skill TEXT NOT NULL,
    input TEXT NOT NULL,
    output TEXT,
    status TEXT NOT NULL,
    error TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);

CREATE TABLE skills (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    version TEXT NOT NULL,
    location TEXT NOT NULL,
    permissions TEXT NOT NULL,
    metadata TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE schedules (
    id TEXT PRIMARY KEY,
    skill TEXT NOT NULL,
    cron_expression TEXT NOT NULL,
    input TEXT NOT NULL,
    enabled BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (skill) REFERENCES skills(id)
);

CREATE TABLE logs (
    id TEXT PRIMARY KEY,
    level TEXT NOT NULL,
    source TEXT NOT NULL,
    message TEXT NOT NULL,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### 1.4 Логирование
- [ ] Настроить структурированный logger (slog/logrus)
- [ ] Реализовать JSON формат
- [ ] Реализовать маскирование секретов
- [ ] Создать логирование по уровням (DEBUG, INFO, WARN, ERROR)

---

## Фаза 2: Core Gateway и API (1-2 недели)

### Задачи

#### 2.1 Message Router
- [ ] Определить интерфейс `Event`
- [ ] Реализовать router для обработки событий
- [ ] Реализовать диспетчеризацию по каналам
- [ ] Создать event bus для pub/sub

**Интерфейс Event:**
```go
type Event struct {
    ID        string                 `json:"id"`
    Channel   string                 `json:"channel"` // "telegram", "discord", "web"
    UserID    string                 `json:"user_id"`
    Message   string                 `json:"message"`
    Metadata  map[string]string      `json:"metadata"`
    Timestamp time.Time              `json:"timestamp"`
}

type EventHandler interface {
    Handle(ctx context.Context, event Event) error
}
```

#### 2.2 Orchestrator
- [ ] Реализовать orchestrator для выбора модели и цепочек навыков
- [ ] Создать промпт-темплейты (system, user, assistant)
- [ ] Реализовать управление контекстом (context window)
- [ ] Создать политику выбора модели

**Интерфейс Orchestrator:**
```go
type Orchestrator interface {
    ProcessMessage(ctx context.Context, event Event) (string, error)
    SelectLLM(taskType string) LLMProvider
    SelectSkills(ctx context.Context, task string) ([]Skill, error)
}
```

#### 2.3 HTTP API
- [ ] Создать HTTP сервер (chi/gin)
- [ ] Реализовать endpoints:
  - `POST /api/v1/chat` - отправить сообщение
  - `GET /api/v1/sessions` - список сессий
  - `GET /api/v1/sessions/{id}` - детали сессии
  - `POST /api/v1/skills/{name}` - выполнить навык
  - `GET /api/v1/skills` - список навыков
  - `GET /api/v1/metrics` - метрики
  - `GET /health` - health check
- [ ] Реализовать middleware (auth, logging, cors)

#### 2.4 WebSocket API
- [ ] Создать WebSocket сервер (gorilla/websocket)
- [ ] Реализовать endpoints:
  - `ws://host/ws/chat/{session}` - чат в реальном времени
  - `ws://host/ws/logs` - логи в реальном времени
- [ ] Создать manager для управления соединениями

---

## Фаза 3: Connectors (2-3 недели)

### Задачи

#### 3.1 Общий интерфейс
- [ ] Определить интерфейс `Connector`
- [ ] Создать registry для коннекторов
- [ ] Реализовать lifecycle management (start/stop)

**Интерфейс Connector:**
```go
type Connector interface {
    // Запуск коннектора
    Start(ctx context.Context) error

    // Остановка коннектора
    Stop() error

    // Канал событий
    Events() <-chan Event

    // Отправка ответа
    SendMessage(ctx context.Context, userID, message string) error
}
```

#### 3.2 Telegram Connector
- [ ] Интегрировать Telegram Bot API (go-telegram-bot-api)
- [ ] Реализовать обработку сообщений
- [ ] Реализовать отправку ответов
- [ ] Добавить поддержку файлов/изображений
- [ ] Создать Whitelist пользователей и чатов

#### 3.3 Discord Connector (v1.0)
- [ ] Интегрировать Discord Bot API (discordgo)
- [ ] Реализовать обработку сообщений
- [ ] Реализовать поддержку embed сообщений
- [ ] Создать Whitelist ролей и каналов

#### 3.4 Web UI Connector
- [ ] Реализовать HTTP endpoint для чата
- [ ] Создать WebSocket для реального времени
- [ ] Интегрировать с message router

---

## Фаза 4: Skills Layer (2-3 недели)

### Задачи

#### 4.1 Парсер SKILL.md
- [ ] Реализовать парсер YAML frontmatter
- [ ] Реализовать парсер Markdown тела
- [ ] Создать валидацию схемы
- [ ] Реализовать сканирование директорий навыков

**Структура Skill:**
```go
type Skill struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Emoji       string                 `json:"emoji"`
    Version     string                 `json:"version"`
    Author      string                 `json:"author"`
    Homepage    string                 `json:"homepage"`
    Location    string                 `json:"location"`
    Tags        []string               `json:"tags"`
    Category    string                 `json:"category"`
    Permissions []string              `json:"permissions"`
    EnvRequired bool                   `json:"env_required"`
    Metadata    map[string]interface{} `json:"metadata"`
    Requirements SkillRequirements      `json:"requirements"`
    Instructions string                `json:"instructions"`
}

type SkillRequirements struct {
    Binaries []string `json:"binaries"`
    Files    []string `json:"files"`
    Env      []string `json:"env"`
}
```

#### 4.2 Runtime для навыков
- [ ] Создать runtime для Bash навыков
- [ ] Создать runtime для Python навыков
- [ ] Создать runtime для Node.js навыков
- [ ] Реализовать sandbox (контейнеры или изоляция процессов)
- [ ] Реализовать таймауты и лимиты ресурсов

**Интерфейс SkillRuntime:**
```go
type SkillRuntime interface {
    Execute(ctx context.Context, skill Skill, input map[string]interface{}) (map[string]interface{}, error)
    Validate(skill Skill) error
}
```

#### 4.3 Базовые навыки (MVP)
- [ ] `shell-run` - выполнение shell команд
- [ ] `file-read` - чтение файлов
- [ ] `file-write` - запись файлов
- [ ] `http-request` - HTTP клиент
- [ ] `git-basic` - базовые git операции (clone, status, add, commit, push)
- [ ] `reminder` - напоминания

**Пример SKILL.md (shell-run):**
```yaml
---
name: shell-run
description: Safely run shell commands on the local machine
emoji: 🐚
version: 1.0.0
author: nexflow
location: ./run.sh
tags: [system, shell, cli]
category: system
permissions: [shell, filesystem]
env_required: false
metadata: {"timeoutSec": 30, "maxOutputKb": 64}
requirements:
  binaries: [bash]
  files: [./run.sh]
  env: []
---

# Shell Run Skill

## Purpose
This skill lets you safely execute simple shell commands.

## When to use
- User asks to "run a shell command"
- Need fresh system information
- Single non-interactive command

## How to use
1. Restate user's goal
2. Propose safe command
3. Ask for confirmation if command modifies data
4. Call skill with command and cwd

## Input schema
- command (string, required): Shell command
- cwd (string, optional): Working directory

## Output schema
- exit_code (integer): Process exit code
- stdout (string): Standard output
- stderr (string): Standard error
```

---

## Фаза 5: LLM Integration (2-3 недели)

### Задачи

#### 5.1 Абстрактный интерфейс LLM
- [ ] Определить интерфейс `LLMProvider`
- [ ] Создать factory для провайдеров

**Интерфейс LLMProvider:**
```go
type LLMProvider interface {
    // Generate completion
    Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)

    // Generate with tools (tool calling)
    GenerateWithTools(ctx context.Context, req GenerateRequest, tools []Tool) (*GenerateResponse, error)

    // Stream completion
    Stream(ctx context.Context, req GenerateRequest) (<-chan string, error)

    // Estimate cost
    EstimateCost(req GenerateRequest) (float64, error)
}

type GenerateRequest struct {
    Messages   []Message `json:"messages"`
    Model      string    `json:"model"`
    MaxTokens  int       `json:"max_tokens"`
    Temperature float64  `json:"temperature"`
}

type GenerateResponse struct {
    Message   Message `json:"message"`
    ToolCalls []ToolCall `json:"tool_calls,omitempty"`
    Tokens    Usage    `json:"tokens"`
}
```

#### 5.2 OpenAI Provider
- [ ] Реализовать OpenAI API клиент
- [ ] Добавить поддержку streaming
- [ ] Реализовать tool calling
- [ ] Добавить cost estimation

#### 5.3 Anthropic Provider
- [ ] Реализовать Anthropic API клиент (Claude)
- [ ] Добавить поддержку streaming
- [ ] Реализовать tool calling
- [ ] Добавить cost estimation

#### 5.4 Ollama Provider
- [ ] Реализовать Ollama API клиент
- [ ] Добавить поддержку streaming
- [ ] Реализовать tool calling
- [ ] Добавить cost estimation (0 for local)

#### 5.5 Google Gemini Provider
- [ ] Реализовать Gemini API клиент
- [ ] Добавить поддержку streaming
- [ ] Реализовать tool calling
- [ ] Добавить cost estimation

#### 5.6 z.ai Provider
- [ ] Реализовать z.ai API клиент
- [ ] Добавить поддержку streaming
- [ ] Реализовать tool calling
- [ ] Добавить cost estimation

#### 5.7 OpenRouter Provider
- [ ] Реализовать OpenRouter API клиент (агрегатор провайдеров)
- [ ] Добавить поддержку streaming
- [ ] Реализовать tool calling
- [ ] Добавить поддержку множественных моделей через один endpoint
- [ ] Добавить cost estimation

#### 5.8 Custom Provider
- [ ] Реализовать generic OpenAI-compatible client
- [ ] Поддержка кастомных endpoint'ов
- [ ] Настройка через конфигурацию

#### 5.9 Memory и контекст
- [ ] Реализовать SQLite хранилище памяти
- [ ] Создать Markdown профили (USER.md, WORKSPACE.md)
- [ ] Реализовать semantic search (векторный поиск)
- [ ] Создать менеджер контекста (context window management)

---

## Фаза 6: Web UI (1-2 недели)

### Задачи

#### 6.1 Frontend setup
- [ ] Создать Svelte проект
- [ ] Настроить структуру компонентов
- [ ] Интегрировать с Go backend

#### 6.2 Чат компонент
- [ ] Создать компонент чата
- [ ] Реализовать WebSocket подключение
- [ ] Добавить Markdown рендеринг
- [ ] Создать историю сообщений

#### 6.3 Дашборд
- [ ] Создать компонент дашборда
- [ ] Отображение активных сессий
- [ ] Просмотр логов
- [ ] Управление навыками

#### 6.4 Управление конфигурацией
- [ ] Создать компонент конфигурации
- [ ] Редактирование config.yml
- [ ] Управление каналами
- [ ] Управление LLM провайдерами

---

## Фаза 7: Observability и тестирование (1-2 недели)

### Задачи

#### 7.1 Observability
- [ ] Реализовать `/metrics` endpoint (Prometheus)
- [ ] Создать health checks
- [ ] Добавить алерты (например, через Telegram)
- [ ] Создать dashboard для метрик

#### 7.2 Тестирование навыков
- [ ] Создать unit тесты для shell-run
- [ ] Создать unit тесты для file-read/write
- [ ] Создать unit тесты для http-request
- [ ] Создать интеграционные тесты для sandbox

#### 7.3 Тестирование коннекторов
- [ ] Создать unit тесты для Telegram connector
- [ ] Создать unit тесты для Discord connector
- [ ] Создать E2E тесты для полного цикла

#### 7.4 Тестирование LLM
- [ ] Создать mock LLM для тестов
- [ ] Создать unit тесты для провайдеров
- [ ] Создать интеграционные тесты с реальными API

---

## Фаза 8: v1.0 расширения (4-6 недель)

### Задачи

#### 8.1 Дополнительные коннекторы
- [ ] Email connector (IMAP/SMTP)
- [ ] Webhook connector
- [ ] Slack connector (опционально)

#### 8.2 Расширенные навыки
- [ ] GitHub/GitLab API навыки
- [ ] AWS/GCP/Azure навыки
- [ ] Home Assistant навыки
- [ ] Kubernetes навыки
- [ ] Docker навыки
- [ ] Monitoring навыки (Prometheus, Grafana)

#### 8.3 Marketplace навыков
- [ ] Создать общественный репозиторий навыков
- [ ] Реализовать установку навыков из репозитория
- [ ] Создать систему рейтингов и отзывов

#### 8.4 Документация
- [ ] Quickstart guide
- [ ] API reference
- [ ] Guide для создания навыков
- [ ] Примеры конфигураций
- [ ] Troubleshooting guide

---

## Дополнительные задачи

### Безопасность
- [ ] Реализовать JWT аутентификацию для API
- [ ] Создать RBAC (role-based access control)
- [ ] Добавить rate limiting
- [ ] Реализовать аудит логов

### Производительность
- [ ] Кэширование LLM ответов (Redis)
- [ ] Кэширование описаний навыков
- [ ] Оптимизация SQL запросов
- [ ] Индексы для БД

### Масштабирование (v2.0+)
- [ ] Horizontal scaling для core
- [ ] Redis для распределенных блокировок
- [ ] Message queue (RabbitMQ/Kafka) для навыков
- [ ] Load balancer

---

## Зависимости

```
Phase 1 (Инфраструктура)
  ├── Config
  ├── Database
  └── Logging

Phase 2 (Core Gateway)
  ├── Phase 1
  ├── Message Router
  ├── Orchestrator
  └── API

Phase 3 (Connectors)
  ├── Phase 2
  ├── Telegram Connector
  ├── Discord Connector
  └── Web UI Connector

Phase 4 (Skills Layer)
  ├── Phase 1
  ├── SKILL.md Parser
  ├── Skill Runtime
  └── Basic Skills

Phase 5 (LLM Integration)
  ├── Phase 2
  ├── LLM Interface
  ├── OpenAI Provider
  ├── Ollama Provider
  ├── Gemini Provider
  ├── Custom Provider
  └── Memory/Context

Phase 6 (Web UI)
  ├── Phase 2
  ├── Frontend Setup
  ├── Chat Component
  ├── Dashboard
  └── Config Manager

Phase 7 (Observability & Testing)
  ├── All previous phases
  ├── Metrics
  ├── Health Checks
  └── Testing

Phase 8 (v1.0 Extensions)
  ├── All previous phases
  ├── Additional Connectors
  ├── Advanced Skills
  ├── Marketplace
  └── Documentation
```

---

## Приоритеты

### P0 (MVP блокеры)
- Проектная структура и конфигурация
- База данных
- Message Router
- HTTP API
- LLM Provider (Anthropic + OpenAI + Ollama)
- Telegram Connector
- Базовые навыки (shell, files, http)
- Web UI (базовый чат)

### P1 (Критические для MVP)
- WebSocket API
- Orchestrator
- Skill Runtime (Bash)
- Memory (базовая)
- Логирование
- Базовое тестирование

### P2 (Важные для MVP)
- Discord Connector
- Skill Runtime (Python, Node.js)
- Semantic search
- Observability (базовая)
- Расширенное тестирование

### P3 (v1.0)
- Email Connector
- Webhook Connector
- Расширенные навыки (CI/CD, облака)
- Marketplace навыков
- Документация

---

## Оценка сроков

| Фаза | Задачи | Оценка | Старт | Финиш |
|------|--------|--------|-------|-------|
| Phase 1 | Базовая инфраструктура | 1 неделя | - | - |
| Phase 2 | Core Gateway и API | 1-2 недели | - | - |
| Phase 3 | Connectors | 2-3 недели | - | - |
| Phase 4 | Skills Layer | 2-3 недели | - | - |
| Phase 5 | LLM Integration | 2-3 недели | - | - |
| Phase 6 | Web UI | 1-2 недели | - | - |
| Phase 7 | Observability и тестирование | 1-2 недели | - | - |
| **MVP Total** | | **4-6 недель** | | |
| Phase 8 | v1.0 расширения | 4-6 недель | - | - |
| **v1.0 Total** | | **8-12 недель** | | |

---

## Следующие шаги

1. Утвердить план с командой
2. Разбить на задачи в трекере (GitHub Issues/Jira)
3. Начать с Phase 1 (Базовая инфраструктура)
4. Еженедельные sync для отслеживания прогресса
