 # Nexflow Implementation Plan

## Обзор прогресса (на Январь 2026)

**Общий статус проекта: ~30% завершено**

### ✅ ЗАВЕРШЕНО (MVP блокеры):
- Проектная структура и конфигурация (YAML/JSON + ENV)
- База данных (SQLite/Postgres) с миграциями и SQLC
- Логирование (slog) с JSON форматом и маскированием секретов
- Domain layer (entities, repositories, value objects) - ПОЛНОСТЬЮ
- Application layer (use cases, DTOs, ports) - ПОЛНОСТЬЮ
- Infrastructure layer (database, http, llm, skills, channels) - БАЗОВАЯ СТРУКТУРА
- DI контейнер с mock реализациями
- Все репозитории (User, Session, Message, Task, Skill, Schedule)
- Все use cases (Chat, User, Skill, Schedule)
- Шаблоны bootstrap файлов (AGENTS.md, SOUL.md, USER.md, NOTES.md)
- Unit тесты для domain, use cases, database

### 🔄 ЧАСТИЧНО ЗАВЕРШЕНО:
- Message Router (интерфейсы готовы)
- LLM Provider (интерфейсы + mock реализация)
- Telegram/Discord/Web Connectors (интерфейсы + mock реализация)
- Orchestrator (ChatUseCase с базовой логикой)
- Skills Runtime (интерфейс + mock реализация)
- HTTP API (базовая инфраструктура с middleware)
- Workspace & Memory System (конфигурация + шаблоны, но инъекция не реализована)
- Тестирование (unit тесты есть, но нет для connectors и LLM)

### ❌ НЕ ЗАВЕРШЕНО:
- Реальные LLM провайдеры (Anthropic, OpenAI, Ollama, etc.)
- Реальные коннекторы (Telegram, Discord, Web)
- Реальный skills runtime (Bash, Python, Node.js)
- Message router и event bus
- HTTP API endpoints
- WebSocket API
- Web UI (Svelte frontend)
- Supervised Mode
- Bootstrap Injection Module (загрузка bootstrap файлов)
- Memory Manager (Markdown файлы)
- Setup wizard (nexflow setup)
- Quick Actions & Slash Commands
- Templates System
- Observability (metrics, health checks)
- Heartbeats & Proactive Work

## Обзор

Документ описывает план реализации Nexflow — self-hosted персонального ИИ-агента. План разделен на фазы:
- **MVP** (2-3 недели): Telegram + LLM (один провайдер)
- **MVP+** (2-3 недели): Web UI + несколько провайдеров + базовые навыки
- **v1.0** (4-6 недель): Расширенные фичи из clawgo/gru
- **v1.1+** (дополнительно): mDNS, FIFO, Routing plugins

## Концепции из проектов clawgo и gru

| Концепция | Источник | Приоритет | Фаза |
|-----------|----------|-----------|------|
| **Supervised Mode** | gru | P0 | MVP |
| **Quick Actions** | clawgo | P1 | MVP+ |
| **Slash Commands** | gru | P1 | MVP+ |
| **Templates** | gru | P1 | MVP+ |
| **MCP Client** | gru | P1 | v1.0 |
| **Delivery Providers** | clawgo | P1 | v1.0 |
| **TTS Engines** | clawgo | P2 | v1.0 |
| **Ralph Loops** | gru | P2 | v1.0 |
| **Screenshot Handling** | gru | P2 | v1.0 |
| **Live Deploy** | gru | P2 | v1.0 |
| **mDNS Advertising** | clawgo | P3 | v1.1+ |
| **FIFO Streaming** | clawgo | P3 | v1.1+ |
| **Routing Plugins** | clawgo | P3 | v1.1+ |

## Стек технологий

- **Ядро:** Go 1.22+
- **БД:** SQLite (по умолчанию), Postgres (опционально)
- **Frontend:** Svelte
- **LLM:** Anthropic (Claude), OpenAI, Ollama, Google Gemini, z.ai, OpenRouter + кастомный провайдер
- **Конфигурация:** YAML + JSON
- **Навыки:** Bash, Python, Node.js
- **Деплой:** Docker/Docker Compose

 ---

## MVP Фаза: Telegram + LLM (2-3 недели)

**Цель:** Минимальный рабочий прототип: Telegram бота, который общается с LLM

**MVP Scope:**
- Telegram connector (бот)
- Один LLM провайдер (Anthropic OR OpenAI OR Ollama)
- Простое message routing
- Базовое логирование
- Supervised Mode (безопасность)
- SQLite (минимальная схема)

### MVP.1 Проектная структура и конфигурация (2-3 дня) ✅ ЗАВЕРШЕНО

#### Задачи

- [x] Создать Go проект с модулями
- [x] Настроить структуру директорий: `cmd/`, `internal/`, `pkg/`, `docs/`
- [x] Инициализировать Go modules
- [x] Создать базовый config struct
- [x] Реализовать парсер YAML (минимальный)
- [x] Создать пример `config.yml`

**Статус:** Реализовано в `internal/shared/config/`

**Пример config.yml (MVP):**
```yaml
server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "./data/nexflow.db"

llm:
  default_provider: "anthropic"  # или "openai", "ollama"
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

channels:
  telegram:
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users: [123456789]

supervised_mode:
  enabled: true

agents:
  defaults:
    workspace: "${NEXFLOW_WORKSPACE:~/nexflow}"
    skip_bootstrap: false
    bootstrap_max_chars: 20000

logging:
  level: "info"
  format: "json"
```

### MVP.2 База данных (1 день) ✅ ЗАВЕРШЕНО

#### Задачи

- [x] Создать SQLite схему (минимальная: users, sessions, messages)
- [x] Реализовать basic Go DB layer
- [x] Создать миграции

**Статус:** Реализовано в `internal/infrastructure/persistence/database/` с SQLC, есть SQLite и Postgres схемы

**Минимальная схема SQLite:**
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
```

### MVP.3 Логирование (1 день) ✅ ЗАВЕРШЕНО

#### Задачи

- [x] Настроить структурированный logger (slog)
- [x] Реализовать JSON формат
- [x] Реализовать маскирование секретов (ключи: token, key, password, secret)

**Статус:** Реализовано в `internal/shared/logging/` с маскированием секретов

### MVP.4 Message Router (1-2 дня) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

#### Задачи

- [x] Определить интерфейс `Event` (в `internal/infrastructure/channels/connector.go`)
- [ ] Реализовать basic router
- [ ] Реализовать event bus

**Статус:** Интерфейсы определены, но router и event bus не реализованы

**Интерфейс Event:**
```go
type Event struct {
    ID        string            `json:"id"`
    Channel   string            `json:"channel"` // "telegram"
    UserID    string            `json:"user_id"`
    Message   string            `json:"message"`
    Metadata  map[string]string `json:"metadata"`
    Timestamp time.Time         `json:"timestamp"`
}

type EventHandler interface {
    Handle(ctx context.Context, event Event) error
}
```

### MVP.5 LLM Provider (один) (2-3 дня) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

#### Задачи

- [x] Определить интерфейс `LLMProvider` (в `internal/application/ports/llm_provider.go`)
- [ ] Реализовать один провайдер (Anthropic OR OpenAI OR Ollama)
- [ ] Реализовать basic message generation

**Статус:** Интерфейс и mock реализация готовы, реальные провайдеры не реализованы

**Интерфейс LLMProvider:**
```go
type LLMProvider interface {
    Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
}

type GenerateRequest struct {
    Messages   []Message `json:"messages"`
    Model      string    `json:"model"`
    MaxTokens  int       `json:"max_tokens"`
}

type GenerateResponse struct {
    Message Message `json:"message"`
    Tokens  Usage    `json:"tokens"`
}
```

### MVP.6 Telegram Connector (2-3 дня) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

#### Задачи

- [ ] Интегрировать Telegram Bot API (go-telegram-bot-api)
- [ ] Реализовать обработку сообщений
- [ ] Реализовать отправку ответов
- [ ] Создать Whitelist пользователей
- [ ] Интегрировать с message router

**Статус:** Интерфейсы и mock реализации готовы в `internal/infrastructure/channels/`, реальная интеграция не выполнена

**Интерфейс Connector:**
```go
type Connector interface {
    Start(ctx context.Context) error
    Stop() error
    Events() <-chan Event
    SendMessage(ctx context.Context, userID, message string) error
}
```

### MVP.7 Supervised Mode (1 день) ❌ НЕ ЗАВЕРШЕНО

#### Задачи

- [ ] Реализовать механизм подтверждения действий
- [ ] Создать базовые правила подтверждения
- [ ] Интегрировать в LLM response handling

**Статус:** Не реализовано

**Конфигурация:**
```yaml
supervised_mode:
  enabled: true
  rules:
    - pattern: "rm.*"
      require_confirmation: true
    - pattern: "delete.*"
      require_confirmation: true
```

### MVP.8 Orchestrator (basic) (1-2 дня) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

#### Задачи

- [ ] Реализовать basic orchestrator
- [ ] Простой flow: Event → Load Bootstrap → LLM → Response
- [ ] Управление контекстом (минимальное)
- [ ] Создать базовые промпт-темплейты
- [ ] Интеграция с Memory Manager
- [ ] Загрузка bootstrap файлов в начале сессии

**Статус:** ChatUseCase реализован в `internal/application/usecase/chat_usecase.go` с базовой логикой, но загрузка bootstrap файлов и Memory Manager не реализованы

**Интерфейс Orchestrator:**
```go
type Orchestrator interface {
    ProcessMessage(ctx context.Context, event Event) (string, error)
    LoadBootstrapFiles(ctx context.Context, sessionID string) (BootstrapContext, error)
    DetermineSessionType(ctx context.Context, sessionID string) (SessionType, error)
}

type BootstrapContext struct {
    Soul    string // из SOUL.md
    User     string // из USER.md
    Notes    string // из NOTS.md
    Memory   string // из memory/memory.md (только main session)
    DailyLog string // из memory/YYYY-MM-DD.md (today + yesterday)
}

type SessionType int
const (
    SessionTypeMain SessionType = iota // Direct chat with human
    SessionTypeGroup             // Group chat or shared context
)
```

**Процесс:**
1. Определить тип сессии (main vs group)
2. Загрузить SOUL.md, USER.md, NOTS.md (всегда)
3. Если main session: загрузить memory/memory.md
4. Загрузить memory/YYYY-MM-DD.md (today + yesterday)
5. Сформировать system prompt с контекстом
6. Отправить запрос в LLM

### MVP.10 Workspace & Memory System (2-3 дня) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

#### Задачи

#### 10.1 Структура Workspace
- [x] Создать структуру директорий workspace
- [ ] Реализовать создание workspace с шаблонами
- [x] Добавить настройку `agents.defaults.workspace` в конфиг
- [x] Парсинг ENV подстановок `${VAR_NAME:default_value}`

**Статус:** Конфигурация готова, шаблоны bootstrap файлов есть в `docs/templates/`

**Структура workspace:**
```
~/nexflow/
├── AGENTS.md          # Инструкции агента + "memory"
├── SOUL.md            # Личность, границы, тон, имя/emoji
├── USER.md            # Профиль пользователя
├── NOTS.md           # Локальные заметки пользователя
└── memory/
    ├── memory.md       # Долгосрочная память
    └── YYYY-MM-DD.md  # Ежедневные логи
```

#### 10.2 Bootstrap Injection Module
- [ ] Реализовать загрузку bootstrap файлов
- [ ] Определить порядок инъекции: SOUL.md → USER.md → NOTS.md → memory files
- [ ] Обработка пустых файлов (skip)
- [ ] Обработка отсутствующих файлов (marker line)
- [ ] Лимиты размера файлов (bootstrap_max_chars)
- [ ] Инъекция в system prompt перед началом сессии
- [x] Настройка через `agents.defaults.skip_bootstrap: true`
- [ ] Маркер усечения: `... [file truncated]`

**Статус:** Конфигурация готова, но инъекция не реализована

#### 10.3 Memory Manager
- [ ] Реализовать чтение/запись memory/YYYY-MM-DD.md
- [ ] Реализовать чтение/запись memory/memory.md
- [ ] Автоматическое создание memory/ директории
- [ ] Загрузка today + yesterday files
- [ ] Защита: не загружать MEMORY.md в групповых чатах
- [x] Настройка через `agents.defaults.workspace`
- [ ] Логирование всех операций с памятью

**Статус:** Не реализовано

#### 10.4 Initial Setup Wizard
- [ ] Команда `nexflow setup` или `nexflow onboard`
- [ ] Создание workspace с шаблонами
- [ ] Интерактивные вопросы пользователю:
  - Workspace путь (или дефолтный)
  - Имя пользователя
  - Как обращаться
  - Timezone
  - Имя агента
  - Emoji агента
- [ ] Генерация SOUL.md из ответов
- [ ] Генерация USER.md из ответов
- [ ] Копирование AGENTS.md из шаблона
- [ ] Создание NOTS.md пустым
- [ ] Создание memory/ директории
- [ ] Создание memory/memory.md пустым
- [ ] Skip если workspace уже существует (с флагом --force для перезаписи)

**Статус:** Не реализовано, но шаблоны готовы в `docs/templates/`

#### 10.5 Templates for Bootstrap Files
- [x] Создать шаблон AGENTS.md (docs/templates/AGENTS.md)
- [x] Создать шаблон SOUL.md (docs/templates/SOUL.md)
- [x] Создать шаблон USER.md (docs/templates/USER.md)
- [x] Создать шаблон NOTS.md (docs/templates/NOTS.md)

**Статус:** Реализовано в `docs/templates/`

#### 10.6 Security Rules
- [ ] Реализовать проверку: main session vs other session
- [ ] Не инъектировать MEMORY.md в не-main сессиях
- [ ] Не делиться личными данными в группах
- [ ] Логирование всех операций с памятью
- [ ] Предупреждение в логах при попытке загрузки MEMORY.md в группах

**Конфигурация:**
```yaml
agents:
  defaults:
    workspace: "${NEXFLOW_WORKSPACE:~/nexflow}"
    skip_bootstrap: false
    bootstrap_max_chars: 20000
```

**ENV подстановка:**
```bash
# Дефолтный путь
export NEXFLOW_WORKSPACE=""  # Использует ~/nexflow

# Кастомный путь
export NEXFLOW_WORKSPACE="/custom/path/to/workspace"
```

### MVP.9 Тестирование и документация (1-2 дня) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

#### Задачи

- [ ] Unit тесты для LLM provider
- [ ] Unit тесты для Telegram connector
- [ ] Интеграционный тест: Telegram → LLM → Telegram
- [x] Создать README с quickstart guide
- [ ] Создать `.env.example`

**Статус:** README.md есть, есть unit тесты для domain, use cases, database, но не для LLM provider и Telegram connector

---

## MVP+ Фаза: Расширение MVP (2-3 недели)

**Цель:** Добавить Web UI, несколько LLM провайдеров, базовые навыки

**MVP+ Scope:**
- Web UI (basic чат)
- Несколько LLM провайдеров (Anthropic + OpenAI + Ollama)
- Базовые навыки (shell, files, http)
- Quick Actions & Slash Commands
- Templates (2-3 базовых)

### MVP+.1 Web UI (1 неделя)

#### Задачи

#### 1.1 Frontend setup
- [ ] Создать Svelte проект
- [ ] Настроить структуру компонентов
- [ ] Интегрировать с Go backend

**Статус:** Не реализовано

#### 1.2 Чат компонент
- [ ] Создать компонент чата
- [ ] Реализовать WebSocket подключение
- [ ] Добавить Markdown рендеринг
- [ ] Создать историю сообщений

**Статус:** Не реализовано

#### 1.3 HTTP API для Web UI
- [x] Создать HTTP сервер (chi/gin) - базовая структура есть в `internal/infrastructure/http/`
- [ ] Реализовать endpoints:
  - `POST /api/v1/chat` - отправить сообщение
  - `GET /api/v1/sessions` - список сессий
  - `GET /health` - health check
- [x] Реализовать middleware (logging, cors)

**Статус:** Базовая инфраструктура готова, но endpoints не реализованы

#### 1.4 WebSocket API
- [ ] Создать WebSocket сервер (gorilla/websocket)
- [ ] Реализовать endpoint: `ws://host/ws/chat/{session}`
- [ ] Создать manager для управления соединениями

**Статус:** Директория есть, но не реализовано

### MVP+.2 Несколько LLM провайдеров (3-4 дня) ❌ НЕ ЗАВЕРШЕНО

#### Задачи

- [ ] Реализовать OpenAI API клиент
- [ ] Реализовать Ollama API клиент (если не был в MVP)
- [ ] Реализовать Anthropic API клиент (если не был в MVP)
- [ ] Создать factory для провайдеров
- [ ] Реализовать выбор провайдера по конфигурации
- [ ] Добавить cost estimation

**Статус:** Не реализовано

### MVP+.3 Базовые навыки (4-5 дней) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

#### Задачи

#### 3.1 Парсер SKILL.md
- [ ] Реализовать парсер YAML frontmatter
- [ ] Реализовать парсер Markdown тела
- [ ] Создать валидацию схемы
- [ ] Реализовать сканирование директорий навыков

**Статус:** Не реализовано, но есть формат в `docs/формат SKILL.md`

#### 3.2 Runtime для навыков
- [ ] Создать runtime для Bash навыков
- [ ] Реализовать sandbox (изоляция процессов)
- [ ] Реализовать таймауты и лимиты ресурсов

**Статус:** Интерфейс `SkillRuntime` определен в `internal/application/ports/skill_runtime.go`, mock реализация есть, реальные runtime не реализованы

**Интерфейс SkillRuntime:**
```go
type SkillRuntime interface {
    Execute(ctx context.Context, skill Skill, input map[string]interface{}) (map[string]interface{}, error)
    Validate(skill Skill) error
}
```

#### 3.3 Базовые навыки
- [ ] `shell-run` - выполнение shell команд (с подтверждением)
- [ ] `file-read` - чтение файлов
- [ ] `file-write` - запись файлов
- [ ] `http-request` - HTTP клиент

**Статус:** Не реализовано

### MVP+.4 Quick Actions & Slash Commands (2-3 дня)

#### Задачи

- [ ] Реализовать систему Quick Actions
- [ ] Добавить предустановленные команды: `/status`, `/health`, `/ping`
- [ ] Реализовать Slash Commands
- [ ] Создать базовые команды: `/create`, `/doctor`

**Конфигурация Quick Actions:**
```yaml
quick_actions:
  - name: "status"
    message: "Покажи статус системы"
  - name: "health"
    message: "Проверь здоровье системы"
```

### MVP+.5 Templates (2-3 дня)

#### Задачи

- [ ] Реализовать систему шаблонов
- [ ] Создать 2-3 базовых шаблона:
  - `python-bot` - Python бот на FastAPI
  - `go-service` - Go микросервис
  - `lambda-function` - AWS Lambda
- [ ] Добавить команду `/create <template>`

**Конфигурация Templates:**
```yaml
templates:
  python-bot:
    description: "Бот на Python с FastAPI"
    files:
      - main.py
      - requirements.txt
      - config.yaml
    commands:
      - "pip install -r requirements.txt"
      - "python main.py"
```

### MVP+.6 Расширенное тестирование (1-2 дня)

#### Задачи

- [ ] Unit тесты для навыков
- [ ] E2E тесты для полного цикла
- [ ] Тестирование Web UI
- [ ] Интеграционные тесты с реальными LLM API

---

## Фаза 1: Базовая инфраструктура (1 неделя) ✅ ЗАВЕРШЕНО

### Задачи

#### 1.1 Проектная структура
- [x] Создать Go проект с модулями
- [x] Настроить структуру директорий: `cmd/`, `internal/`, `pkg/`, `skills/`, `docs/`
- [x] Инициализировать Go modules
- [ ] Настроить CI/CD (GitHub Actions)

**Статус:** Структура готова, CI/CD не настроен

#### 1.2 Конфигурация
- [x] Создать структуру конфигурации (struct)
- [x] Реализовать парсер YAML
- [x] Реализовать парсер JSON
- [x] Создать пример `config.yml` и `config.json`
- [x] Реализовать загрузку из ENV переменных

**Статус:** Реализовано в `internal/shared/config/`

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
- [x] Создать SQLite схему (таблицы: users, sessions, messages, tasks, skills, schedules, logs)
- [x] Создать Postgres схему (для прод)
- [x] Реализовать Go ORM (sqlc)
- [x] Создать миграции

**Статус:** Реализовано в `internal/infrastructure/persistence/database/` с SQLC

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
- [x] Настроить структурированный logger (slog/logrus)
- [x] Реализовать JSON формат
- [x] Реализовать маскирование секретов
- [x] Создать логирование по уровням (DEBUG, INFO, WARN, ERROR)

**Статус:** Реализовано в `internal/shared/logging/`

---

## Фаза 2: Core Gateway и API (1-2 недели) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

### Задачи

#### 2.1 Message Router
- [x] Определить интерфейс `Event`
- [ ] Реализовать router для обработки событий
- [ ] Реализовать диспетчеризацию по каналам
- [ ] Создать event bus для pub/sub

**Статус:** Интерфейсы готовы в `internal/infrastructure/channels/connector.go`

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

**Статус:** ChatUseCase есть, но без полноценной логики orchestrator'а

**Интерфейс Orchestrator:**
```go
type Orchestrator interface {
    ProcessMessage(ctx context.Context, event Event) (string, error)
    SelectLLM(taskType string) LLMProvider
    SelectSkills(ctx context.Context, task string) ([]Skill, error)
}
```

#### 2.3 HTTP API
- [x] Создать HTTP сервер (chi/gin) - базовая структура
- [ ] Реализовать endpoints:
  - `POST /api/v1/chat` - отправить сообщение
  - `GET /api/v1/sessions` - список сессий
  - `GET /api/v1/sessions/{id}` - детали сессии
  - `POST /api/v1/skills/{name}` - выполнить навык
  - `GET /api/v1/skills` - список навыков
  - `GET /api/v1/metrics` - метрики
  - `GET /health` - health check
- [x] Реализовать middleware (auth, logging, cors)

**Статус:** Базовая инфраструктура готова в `internal/infrastructure/http/`

#### 2.4 WebSocket API
- [ ] Создать WebSocket сервер (gorilla/websocket)
- [ ] Реализовать endpoints:
  - `ws://host/ws/chat/{session}` - чат в реальном времени
  - `ws://host/ws/logs` - логи в реальном времени
- [ ] Создать manager для управления соединениями

**Статус:** Директория есть, но не реализовано

---

## Фаза 3: Connectors (2-3 недели) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

### Задачи

#### 3.1 Общий интерфейс
- [x] Определить интерфейс `Connector`
- [ ] Создать registry для коннекторов
- [ ] Реализовать lifecycle management (start/stop)

**Статус:** Интерфейс определен в `internal/infrastructure/channels/connector.go`

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

**Статус:** Mock реализация есть в `internal/infrastructure/channels/mock/telegram.go`, реальная интеграция не выполнена

#### 3.3 Discord Connector (v1.0)
- [ ] Интегрировать Discord Bot API (discordgo)
- [ ] Реализовать обработку сообщений
- [ ] Реализовать поддержку embed сообщений
- [ ] Создать Whitelist ролей и каналов

**Статус:** Mock реализация есть в `internal/infrastructure/channels/mock/discord.go`, реальная интеграция не выполнена

#### 3.4 Web UI Connector
- [ ] Реализовать HTTP endpoint для чата
- [ ] Создать WebSocket для реального времени
- [ ] Интегрировать с message router

**Статус:** Mock реализация есть в `internal/infrastructure/channels/mock/web.go`, реальная интеграция не выполнена

---

## Фаза 4: Skills Layer (2-3 недели) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

### Задачи

#### 4.1 Парсер SKILL.md
- [ ] Реализовать парсер YAML frontmatter
- [ ] Реализовать парсер Markdown тела
- [ ] Создать валидацию схемы
- [ ] Реализовать сканирование директорий навыков

**Статус:** Не реализовано, формат описан в `docs/формат SKILL.md`

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

**Статус:** Интерфейс и mock реализация есть в `internal/infrastructure/skills/runtime_adapter.go`, реальные runtime не реализованы

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

**Статус:** Не реализовано

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

## Фаза 5: LLM Integration (2-3 недели) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

### Задачи

#### 5.1 Абстрактный интерфейс LLM
- [x] Определить интерфейс `LLMProvider`
- [ ] Создать factory для провайдеров

**Статус:** Интерфейс определен в `internal/application/ports/llm_provider.go`, mock реализация есть

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

**Статус:** Mock реализация есть для Anthropic, OpenAI, Ollama в `internal/infrastructure/llm/mock/`, реальные провайдеры не реализованы

#### 5.9 Memory и контекст
- [x] Реализовать SQLite хранилище памяти
- [x] Создать Markdown профили (USER.md, WORKSPACE.md)
- [x] Реализовать semantic search (векторный поиск)
- [x] Создать менеджер контекста (context window management)
- [ ] Реализовать Bootstrap Injection Module
- [ ] Реализовать Memory Manager (Markdown files)
- [ ] Интеграция с Orchestrator

**Статус:** SQLite хранилище есть, но Bootstrap Injection и Memory Manager не реализованы. Шаблоны bootstrap файлов готовы в `docs/templates/`

**Bootstrap Files:**
- [ ] AGENTS.md - инструкции для агента
- [ ] SOUL.md - личность агента (+ имя/emoji)
- [ ] USER.md - профиль пользователя
- [ ] NOTS.md - локальные заметки пользователя

**Memory Files:**
- [ ] memory/memory.md - долгосрочная память
- [ ] memory/YYYY-MM-DD.md - ежедневные логи

**Конфигурация:**
```yaml
agents:
  defaults:
    workspace: "${NEXFLOW_WORKSPACE:~/nexflow}"
    skip_bootstrap: false
    bootstrap_max_chars: 20000
```

**Setup Wizard:**
- [ ] `nexflow setup` команда
- [ ] Интерактивное создание workspace
- [ ] Генерация bootstrap файлов из шаблонов

---

## Фаза 6: Web UI (1-2 недели) ❌ НЕ ЗАВЕРШЕНО

### Задачи

#### 6.1 Frontend setup
- [ ] Создать Svelte проект
- [ ] Настроить структуру компонентов
- [ ] Интегрировать с Go backend

**Статус:** Не реализовано

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

**Статус:** Не реализовано

---

## Фаза 7: Observability и тестирование (1-2 недели) 🔄 ЧАСТИЧНО ЗАВЕРШЕНО

### Задачи

#### 7.1 Observability
- [ ] Реализовать `/metrics` endpoint (Prometheus)
- [ ] Создать health checks
- [ ] Добавить алерты (например, через Telegram)
- [ ] Создать dashboard для метрик

**Статус:** Логирование реализовано, но metrics endpoint и health checks не реализованы

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
- [x] Создать mock LLM для тестов
- [ ] Создать unit тесты для провайдеров
- [ ] Создать интеграционные тесты с реальными API

**Статус:** Mock LLM реализован в `internal/infrastructure/llm/mock/`, unit тесты есть для domain и use cases, но не для connectors и real LLM providers

 ---

## v1.0 Фаза: Расширенные возможности (4-6 недель)

**Цель:** Полнофункциональный агент с фичами из clawgo/gru

**v1.0 Scope:**
- Discord connector
- Email connector
- Webhook connector
- MCP Client
- Delivery Providers (WhatsApp/Signal/iMessage)
- TTS Engines (espeak-ng, Piper)
- Ralph Loops (итеративная разработка)
- Screenshot Handling
- Live Deploy (Vercel/Railway)
- Расширенные навыки (CI/CD, облака, HA)
- Observability (metrics, health checks)
- Полная документация

### v1.0.1 Дополнительные коннекторы (1 неделя)

#### Задачи

#### 1.1 Discord Connector
- [ ] Интегрировать Discord Bot API (discordgo)
- [ ] Реализовать обработку сообщений
- [ ] Реализовать поддержку embed сообщений
- [ ] Создать Whitelist ролей и каналов
- [ ] Интегрировать с message router

#### 1.2 Email Connector
- [ ] Реализовать IMAP для чтения почты
- [ ] Реализовать SMTP для отправки уведомлений
- [ ] Создать правила для email-триггеров
- [ ] Интегрировать с message router

#### 1.3 Webhook Connector
- [ ] Реализовать HTTP endpoint для webhooks
- [ ] Создать секретный ключ для валидации
- [ ] Интегрировать с message router
- [ ] Создать документацию по вебхукам

### v1.0.2 MCP Client (3-4 дня)

#### Задачи

- [ ] Реализовать MCP клиент (Model Context Protocol)
- [ ] Подключить стандартные MCP servers:
  - `filesystem` - работа с файловой системой
  - `github` - GitHub API
  - `search` - поиск
  - `database` - работа с БД
- [ ] Создать конфигурацию MCP servers
- [ ] Интегрировать с LLM для tool calling

**Конфигурация MCP:**
```json
{
  "mcpServers": {
    "filesystem": {
      "command": "node",
      "args": ["@modelcontextprotocol/server-filesystem"]
    },
    "github": {
      "command": "node",
      "args": ["@modelcontextprotocol/server-github"]
    }
  }
}
```

### v1.0.3 Delivery Providers (2-3 дня)

#### Задачи

- [ ] Реализовать абстракцию для Delivery Providers
- [ ] Добавить поддержку WhatsApp
- [ ] Добавить поддержку Signal
- [ ] Добавить поддержку iMessage
- [ ] Создать конфигурацию доставки по каналу

**Интерфейс DeliveryProvider:**
```go
type DeliveryProvider interface {
    SendMessage(ctx context.Context, to string, msg string) error
}
```

**Конфигурация:**
```yaml
delivery:
  default_channel: "telegram"
  providers:
    telegram:
      bot_token: "${TELEGRAM_BOT_TOKEN}"
    whatsapp:
      api_key: "${WHATSAPP_API_KEY}"
    signal:
      phone_number: "+1234567890"
    imessage:
      email: "user@icloud.com"
```

### v1.0.4 TTS Engines (2-3 дня)

#### Задачи

- [ ] Реализовать абстракцию для TTS engines
- [ ] Добавить поддержку espeak-ng (системный)
- [ ] Добавить поддержку Piper (быстрые нейронные голоса)
- [ ] Создать конфигурацию TTS
- [ ] Интегрировать с Telegram (voice messages)

**Интерфейс TTSEngine:**
```go
type TTSEngine interface {
    Synthesize(ctx context.Context, text string) ([]byte, error)
}
```

**Конфигурация:**
```yaml
tts:
  engine: "piper"  # "espeak", "piper", "elevenlabs", "none"
  voice: "ru-ru"
  rate: 200
  piper:
    model_path: "./models/piper"
  espeak:
    voice: "ru"
```

### v1.0.5 Ralph Loops (3-4 дня)

#### Задачи

- [ ] Реализовать итеративный цикл разработки
- [ ] Агент выполняет код
- [ ] Запускает тесты
- [ ] Если ошибки → исправляет и повторяет
- [ ] Если нет ошибок → завершает
- [ ] Создать логирование итераций

**Конфигурация:**
```yaml
ralph_loops:
  enabled: true
  max_iterations: 5
  auto_fix: true
  test_command: "go test ./..."
```

### v1.0.6 Screenshot Handling (2-3 дня)

#### Задачи

- [ ] Реализовать загрузку изображений
- [ ] Интегрировать Vision API (Claude/GPT-4V)
- [ ] Создать prompt для анализа UI
- [ ] Генерация HTML/CSS по скриншоту
- [ ] Предпросмотр результата

### v1.0.7 Live Deploy (2-3 дня)

#### Задачи

- [ ] Интегрировать Vercel API
- [ ] Интегрировать Railway API
- [ ] Реализовать команду `/deploy vercel`
- [ ] Реализовать команду `/deploy railway`
- [ ] Создать возврат URL после деплоя

### v1.0.8 Расширенные навыки (1 неделя)

#### Задачи

- [ ] GitHub/GitLab API навыки
- [ ] AWS/GCP/Azure навыки
- [ ] Home Assistant навыки
- [ ] Kubernetes навыки
- [ ] Docker навыки
- [ ] Monitoring навыки (Prometheus, Grafana)

### v1.0.9 Observability (2-3 дня)

#### Задачи

- [ ] Реализовать `/metrics` endpoint (Prometheus)
- [ ] Создать health checks
- [ ] Добавить алерты (через Telegram)
- [ ] Создать dashboard для метрик
- [ ] Структурированные JSON-логи
- [ ] Маскирование секретов

### v1.0.10 Документация (3-4 дня)

#### Задачи

- [ ] Quickstart guide
- [ ] API reference
- [ ] Guide для создания навыков
- [ ] Примеры конфигураций
- [ ] Troubleshooting guide
- [ ] MCP integration guide
- [ ] Templates guide

---

## v1.0.11 Heartbeats & Memory Maintenance (3-4 дня) ❌ НЕ ЗАВЕРШЕНО

#### Задачи

#### 11.1 Heartbeat System
- [ ] Реализовать heartbeat polling
- [ ] Настройка heartbeat prompt в конфиге
- [ ] HEARTBEAT.md файл для задач
- [ ] Проверка почты, календаря, уведомлений
- [ ] batch проверки для оптимизации API calls
- [ ] Настройка интервала проверки (interval_minutes)
- [ ] Quiet hours (quiet_hours)

#### 11.2 Memory Synthesis
- [ ] Автоматическое обновление MEMORY.md из daily logs
- [ ] Review recent `memory/YYYY-MM-DD.md` files
- [ ] Идентификация значимых событий/уроков
- [ ] Удаление устаревшей информации из memory.md
- [ ] Триггер по heartbeat или manual

#### 11.3 Proactive Work
- [ ] Проверка git status проектов
- [ ] Обновление документации
- [ ] Коммит и push собственных изменений
- [ ] Respect quiet time (23:00-08:00 по умолчанию)
- [ ] Tracking состояния проверок (heartbeat-state.json)

**Статус:** Не реализовано

**Конфигурация:**
```yaml
heartbeats:
  enabled: true
  interval_minutes: 30
  prompt: "Read HEARTBEAT.md if it exists. Follow it strictly. Do not infer or repeat old tasks. If nothing needs attention, reply HEARTBEAT_OK."
  quiet_hours: "23:00-08:00"
  checks:
    - email
    - calendar
    - notifications
    - weather
```

**Пример heartbeat-state.json:**
```json
{
  "lastChecks": {
    "email": 1703275200,
    "calendar": 1703260800,
    "weather": null
  }
}
```

---

## v1.1+ Фаза: Дополнительные возможности (дополнительно)

**Цель:** Расширенные фичи из clawgo для продвинутых пользователей

### v1.1.1 mDNS Advertising (2-3 дня)

#### Задачи

- [ ] Реализовать mDNS объявление сервиса
- [ ] Создать сервис `_nexflow-agent._tcp`
- [ ] Реализовать команду `/discover`
- [ ] Конфигурация через `-mdns-service`

### v1.1.2 FIFO Streaming (2-3 дня)

#### Задачи

- [ ] Реализовать поддержку named pipes
- [ ] Конфигурация через `-stdin` и `-stdin-file`
- [ ] Интеграция с голосовыми инструментами
- [ ] Документация по голосовому вводу

### v1.1.3 Routing Plugin System (3-4 дня)

#### Задачи

- [ ] Создать плагинную архитектуру маршрутизации
- [ ] Реализовать правила: по ключевым словам, контексту, конфигурации
- [ ] Создать стандартный плагин "default"
- [ ] Документация по созданию кастомных плагинов

**Конфигурация:**
```yaml
router: "smart"
rules:
  - pattern: ".*код.*"
    destination: "code-agent"
  - pattern: ".*дом.*"
    destination: "home-assistant"
```

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
MVP (Telegram + LLM)
  ├── Проектная структура и конфигурация
  ├── База данных (минимальная)
  ├── Логирование
  ├── Message Router (basic)
  ├── LLM Provider (один)
  ├── Telegram Connector
  ├── Orchestrator (basic)
  └── Supervised Mode

MVP+ (расширение)
  ├── MVP
  ├── Web UI (basic)
  ├── Несколько LLM провайдеров
  ├── Базовые навыки
  ├── Quick Actions
  ├── Slash Commands
  └── Templates (базовые)

v1.0 (полнофункциональный)
  ├── MVP+
  ├── Дополнительные коннекторы (Discord, Email, Webhook)
  ├── MCP Client
  ├── Delivery Providers
  ├── TTS Engines
  ├── Ralph Loops
  ├── Screenshot Handling
  ├── Live Deploy
  ├── Расширенные навыки
  ├── Observability
  └── Документация

v1.1+ (дополнительно)
  ├── v1.0
  ├── mDNS Advertising
  ├── FIFO Streaming
  └── Routing Plugins
```

---

## Приоритеты

### P0 (MVP блокеры)
- Проектная структура и конфигурация
- База данных (минимальная)
- Логирование
- Message Router (basic)
- LLM Provider (один: Anthropic OR OpenAI OR Ollama)
- Telegram Connector
- Orchestrator (basic)
- Supervised Mode

### P1 (Критические для MVP+)
- Web UI (basic)
- Несколько LLM провайдеров
- Базовые навыки (shell, files, http)
- Quick Actions
- Slash Commands
- Templates (2-3 базовых)
- Базовое тестирование

### P2 (Важные для v1.0)
- Discord Connector
- Email Connector
- Webhook Connector
- MCP Client
- Delivery Providers
- TTS Engines (espeak-ng, Piper)
- Ralph Loops
- Screenshot Handling
- Live Deploy
- Расширенные навыки
- Observability (metrics, health checks)
- Полная документация

### P3 (v1.1+)
- mDNS Advertising
- FIFO Streaming
- Routing Plugins
- Расширенные TTS (ElevenLabs)
- Масштабирование

 ---

## Оценка сроков

| Фаза | Задачи | Оценка | Старт | Финиш |
|------|--------|--------|-------|-------|
| MVP | Telegram + LLM (один провайдер) | 2-3 недели | Январь 2026 | Январь 2026 |
| MVP+ | Web UI + навыки + шаблоны | 2-3 недели | Январь 2026 | Февраль 2026 |
| **MVP+ Total** | | **4-6 недель** | | |
| v1.0 | Полнофункциональный агент | 4-6 недель | Март 2026 | Апрель 2026 |
| **v1.0 Total** | | **8-12 недель** | | |
| v1.1+ | Дополнительно | TBD | Q2 2026 | Q3 2026 |

---

## Следующие шаги

1. Утвердить план с командой
2. Разбить на задачи в трекере (GitHub Issues/Jira)
3. Начать с MVP Фазы (Telegram + LLM)
4. Еженедельные sync для отслеживания прогресса
5. После MVP: переход к MVP+ (Web UI + навыки)
6. После MVP+: переход к v1.0 (расширенные возможности)
