# Nexflow - Product Requirements Document v2.0

## Введение

Nexflow — self-hosted персональный ИИ-агент, который управляет цифровыми потоками задач: от чатов и почты до DevOps-рутину, работая на вашей машине и в вашей инфраструктуре.

## Цели продукта

- Дать разработчику один "поток" (flow), где сходятся мессенджеры, CI/CD, репозитории, локальные скрипты и расписания
- Обеспечить безопасный локальный gateway от текстовых команд к реальным действиям: файлы, терминал, браузер, внешние API
- Сделать установку и расширение максимально простыми для инженеров

**Ключевые метрики:**
- Время холодного деплоя ≤ 10 минут
- Среднее время ответа ≤ 5 секунд для простых задач
- ≥ 20 стабильных интеграций (skills) в v1.0

## Целевая аудитория

- Senior backend/DevOps/infra инженеры
- Энтузиасты локальных агентов

Nexflow говорит на русском и английском, UX/доки ориентированы на инженеров из СНГ.

## Ключевые функции

### 1. Multi-channel Gateway

**Поддерживаемые каналы v1.0:**
- Telegram (бот + личные сообщения)
- Discord (сервер, личные сообщения)
- Web UI (встроенный чат для отладки)

**Требования:**
- Единый user identity через все каналы
- Конфигурация через YAML/JSON + ENV
- Сообщение → event → LLM → tools → response

### 2. ИИ-ядро и модели

**Поддержка провайдеров:**
- Облачные API: Anthropic (Claude), OpenAI, Google Gemini, z.ai (BYO key)
- Агрегаторы: OpenRouter (множественные провайдеры через один endpoint)
- Локальные модели: Ollama
- Кастомные провайдеры: OpenAI-compatible API

**Фичи:**
- System/persona-профили в Markdown (USER.md, WORKSPACE.md)
- Переключение модели по политике (сложные задачи → мощная модель, cheap → компактная)
- Tool calling и агентные цепочки

### 3. Память и контекст

**Типы памяти:**
- Долгосрочная: SQLite/Markdown (профиль, проекты, предпочтения)
- Сессионная: контекст текущего диалога
- Кэш: результаты повторяющихся запросов

**Функции:**
- Semantic-search по памяти
- Редактируемые "rules" и "skills hints"

### 4. Skills (инструменты)

**Базовый набор (MVP):**
- `shell-run`: выполнение shell команд с подтверждением
- `file-read/write`: файловые операции
- `http-request`: HTTP клиент
- `git-basic`: базовые git операции
- `reminder/cron`: напоминания

**Расширенный набор (v1.0):**
- CI/CD: GitHub/GitLab API, запуск тестов, деплой
- Облака: AWS/Azure/GCP управление
- Коммуникации: почта, календарь
- Домашняя автоматизация: Home Assistant
- DevOps: мониторинг, логи, алерты

**Формат навыков:**
- `SKILL.md` с YAML frontmatter + Markdown инструкция
- Совместимость с Moltbot/Clawdbot
- Плагины: Bash, Python, Node.js

### 5. Безопасность и наблюдаемость

**Безопасность:**
- Все действия логируются (время, вход, выход, статус)
- Sandboxing: ограничение директорий, команд, хостов
- Whitelist ресурсов и API ключей

**Observability:**
- Встроенный Web dashboard (логи, задачи, очереди)
- Структурированные JSON-логи
- `/metrics` endpoint для Prometheus
- Health checks

## Архитектура

### High-level схема

```
┌─────────────────────────────────────────────────────────────┐
│                      Channels Layer                         │
│  Telegram  │  Discord  │  Web UI  │  Email  │  Webhooks     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Core Gateway (Go)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────┐  │
│  │ Message      │  │ Orchestrator │  │ LLM Router     │  │
│  │ Router       │  │              │  │                │  │
│  └──────────────┘  └──────────────┘  └────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │            Skills Layer (Plugin System)              │  │
│  │  System  │  Git  │  HTTP  │  Custom (Markdown)       │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Storage & Execution                       │
│  SQLite/Postgres │  FS (Markdown) │  Sandbox Containers   │
└─────────────────────────────────────────────────────────────┘
```

### Компоненты

**Core сервис (Go):**
- HTTP/WebSocket API (порт 8080 по умолчанию)
- Message router: события каналов → LLM → skills
- Orchestrator: выбор модели, цепочки инструментов
- Cron планировщик

**Connectors:**
- Отдельные пакеты для каждого канала
- Единый event-интерфейс
- Сохранение состояния соединений

**Skills Layer:**
- Интерфейс для инструментов
- Плагины на Bash/Python/Node.js
- Markdown описания (`SKILL.md`)
- Sandbox для опасных навыков

**Хранилища:**
- SQLite/Postgres: события, задачи, очереди, настройки
- FS: профиль пользователя, инструкции, конфиги навыков

## Технические требования

### Платформа и системные требования

- **ОС:** Linux (Ubuntu/Debian), macOS 13+, Windows через WSL2
- **Ресурсы:** 2+ vCPU, 2–4 ГБ RAM (без локальных LLM), от 8 ГБ RAM (с локальными)
- **Хранилище:** ≥10 ГБ SSD
- **Сеть:** интернет для облачных LLM/каналов; поддержка офлайн режима

### Технологический стек

- **Ядро:** Go ≥1.22, монолитный бинарь
- **БД:** SQLite по умолчанию, опционально Postgres
- **Frontend:** Svelte
- **LLM:** Anthropic (Claude), OpenAI, Ollama, Google Gemini, z.ai, OpenRouter + кастомный провайдер
- **Конфигурация:** YAML + JSON
- **Навыки:** Bash, Python, Node.js
- **Деплой:** Docker/Docker Compose

### Сетевое взаимодействие

- **Режимы:** loopback по умолчанию (127.0.0.1), явное включение внешнего доступа
- **Аутентификация:** gateway token для каналов и UI
- **Авторизация:** права/роль для навыков (FS, shell, сеть, секреты)
- **Sandboxing:** опасные навыки в отдельных контейнерах

### API Gateway

**HTTP API:**
- `POST /api/v1/chat` - отправить сообщение
- `GET /api/v1/sessions` - список сессий
- `GET /api/v1/sessions/{id}` - детали сессии
- `POST /api/v1/skills/{name}` - выполнить навык
- `GET /api/v1/skills` - список зарегистрированных навыков
- `GET /api/v1/metrics` - метрики (Prometheus)
- `GET /health` - health check

**WebSocket:**
- `ws://host/ws/chat/{session}` - чат в реальном времени
- `ws://host/ws/logs` - логи в реальном времени

## Модуль Skills

### Структура навыка

Каждый навык - директория с `SKILL.md` + опциональные скрипты:

```
skills/
├── shell-run/
│   ├── SKILL.md
│   └── run.sh
├── git-basic/
│   ├── SKILL.md
│   └── main.py
└── github-issues/
    ├── SKILL.md
    └── skill.js
```

### Формат SKILL.md

**YAML Frontmatter (обязательные поля):**
```yaml
---
name: shell-run
description: Safely run shell commands on the local machine
emoji: 🐚
version: 1.0.0
author: your-name
homepage: https://github.com/youruser/nexflow-skills
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
```

**Markdown тело (обязательные секции):**
1. **Purpose** - что делает навык
2. **When to use** - когда использовать/не использовать
3. **How to use** - пошаговая инструкция
4. **Input schema** - параметры и типы
5. **Output schema** - что возвращает
6. **Examples** - 2–3 примера
7. **Implementation details** (опционально) - подсказки для рантайма

### Интерфейс рантайма

**Запуск навыка:**
- Рантайм выполняет `location` с параметрами из input schema
- Результат должен соответствовать output schema
- Таймауты из `metadata` или дефолтные (30s)
- Логирование входа/выхода/ошибок

**Sandbox:**
- Изолированный процесс/контейнер для навыков с `permissions: [shell, filesystem, network, secrets]`
- Ограничение по времени, памяти, IO
- Whitelist директорий и хостов

## Интеграции (Connectors)

### Общий интерфейс

```go
type Connector interface {
    Start(ctx context.Context) error
    Stop() error
    Events() <-chan Event
    SendMessage(ctx context.Context, userID, message string) error
}

type Event struct {
    ID        string
    Channel   string // "telegram", "discord", "web"
    UserID    string
    Message   string
    Metadata  map[string]string
    Timestamp time.Time
}
```

### Telegram Connector

**Функции:**
- Бот для групповых чатов
- Личные сообщения
- Markdown/HTML форматирование
- Загрузка файлов

**Конфигурация:**
```yaml
connectors:
  telegram:
    bot_token: "${TELEGRAM_BOT_TOKEN}"
    allowed_users: [123456789, 987654321]
    allowed_chats: [-1001234567890]
```

### Discord Connector

**Функции:**
- Server channels (каналы сервера)
- Direct messages (личные сообщения)
- Embed сообщения для богатого форматирования
- Webhook для интеграций

**Конфигурация:**
```yaml
connectors:
  discord:
    bot_token: "${DISCORD_BOT_TOKEN}"
    guild_id: "123456789012345678"
    allowed_roles: ["admin", "dev"]
    channels:
      general: "123456789012345678"
      dev: "987654321098765432"
```

### Web UI Connector

**Функции:**
- Встроенный SPA на Svelte
- WebSocket для реального времени
- Дашборд: логи, задачи, навыки
- Управление конфигурацией

## Требования к реализации

### MVP (4–6 недель)

**Ядро:**
- Один бинарь Nexflow с конфигурацией (env + config.yml)
- SQLite БД
- HTTP/WebSocket API
- Message router

**Каналы:**
- Telegram (бот)
- Web UI (basic чат)

**LLM:**
- Облачные провайдеры (Anthropic + OpenAI + z.ai + OpenRouter)
- Локальный провайдер (Ollama)
- Кастомный провайдер
- Базовое переключение моделей

**Память:**
- Markdown профили (USER.md, WORKSPACE.md)
- SQLite хранилище событий
- Базовый semantic search

**Навыки (v0):**
- `shell-run` (с подтверждением)
- `file-read/write`
- `http-request`
- `git-basic`
- `reminder/cron`

**UI:**
- Базовый чат в Web UI
- Просмотр последних действий
- Логи навыков

### v1.0 (8–12 недель)

**Каналы:**
- Discord
- Email (уведомления/триггеры)
- Webhooks для внешних систем

**Навыки:**
- ≥ 20–30 навыков (CI/CD, облака, домашняя автоматизация)
- Marketplace навыков (общий репозиторий)
- Кастомные навыки пользователей

**LLM:**
- Множественные провайдеры (Anthropic, OpenAI, Ollama, Google Gemini, z.ai, OpenRouter)
- Гибкий роутинг по политикам
- Файн-тюнинг промптов

**Observability:**
- Структурированные логи (JSON)
- `/metrics` endpoint (Prometheus)
- Алерты и уведомления
- Health checks

**Документация:**
- Quickstart guide
- Примеры конфигов
- Guide для создания навыков
- API reference

## Тестирование

### Тестирование навыков

**Единичные тесты:**
- Mock input/output
- Проверка контракта (input/output schema)
- Таймауты и обработка ошибок

**Интеграционные тесты:**
- Проверка sandbox
- Файловые операции
- HTTP запросы (с mock серверами)
- Shell команды (безопасные)

### Тестирование коннекторов

**Модульные тесты:**
- Event routing
- Message parsing
- Error handling

**E2E тесты:**
- Полный цикл: пользовательское сообщение → LLM → навык → ответ
- Тестовые аккаунты для Telegram/Discord

### Тестирование LLM интеграций

**Unit тесты:**
- Mock LLM ответы
- Проверка вызова навыков

**Интеграционные тесты:**
- Реальные LLM API (с лимитом запросов)
- Проверка токенов/цены

## Масштабирование и производительность

### Горизонтальное масштабирование

**Stateless core:**
- Message router и orchestrator могут быть горизонтально масштабированы
- Shared storage: Postgres вместо SQLite
- Redis для распределенных блокировок

**Очереди задач:**
- Отдельный сервис для очереди навыков (например, Redis Queue)
- Worker процессы для выполнения

### Кэширование

**Кэш LLM ответов:**
- Хранение повторяющихся запросов
- TTL для старения кэша
- In-memory (Redis)

**Кэш навыков:**
- Кэш `SKILL.md` описаний
- Кэш результатов навыков (idempotency key)

## Риски и ограничения

### Технические риски

- **Сложность настройки:** важна простота деплоя (одна команда)
- **Безопасность:** нужно явное whitelisting и sandbox
- **Переменчивость LLM:** разные модели требуют разных промптов

### Бизнес риски

- **Конкуренция:** Moltbot, Clawdbot
- **Поддержка:** нужна активная комьюнити и контрибьюторы
- **Документация:** критически важна для пользовательского опыта

## Roadmap

### Q1 2026
- MVP (4–6 недель)
- Базовые навыки
- Telegram + Web UI

### Q2 2026
- v1.0 (8–12 недель)
- Discord, Email
- Расширенные навыки
- Observability

### Q3 2026
- Marketplace навыков
- Горизонтальное масштабирование
- Fine-tuning промптов

### Q4 2026
- Кастомные LLM модели
- Enterprise features
- Облачный SaaS (опционально)
