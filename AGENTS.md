# Правила работы с проектом Nexflow

Этот файл — главный вход для ИИ-агентов. Подробные правила в `docs/rules/*.md`.

## Обзор проекта

**Nexflow** — self-hosted ИИ-агент на Go, управляющий цифровыми потоками задач через multiple channels (Telegram, Discord, Web UI) с LLM-провайдерами (Anthropic, OpenAI, Ollama и др.) и навыками (skills).

**Стек:** Go 1.25.5, SQLite/Postgres, SQLC, slog, YAML/JSON + ENV, GitHub Actions

## Быстрый старт

```bash
go mod download && go test ./...
go run cmd/server/main.go
```

## Модули правил (lazy loading)

**Обязательные (всегда):**
- `docs/rules/security.md` — безопасность (КРИТИЧНО)
- `docs/rules/projectrules.md` — общие правила проекта

**По контексту:**
- Архитектура → `docs/rules/architecture.md`
- Качество кода → `docs/rules/codequality.md`
- REST/Web API → `docs/rules/apidesign.md`
- База данных → `docs/rules/database.md`
- Тестирование → `docs/rules/testing.md`

## Приоритеты

1. **Безопасность** > всё остальное
2. **Архитектура** > удобство
3. **Качество кода** > скорость
4. **Тесты** > фичи

## Структура проекта

```
nexflow/
├── cmd/              # Entry points (main.go)
├── internal/         # Private packages
│   ├── config/       # Конфигурация
│   ├── database/     # БД (SQLite/Postgres, SQLC)
│   ├── logging/      # Структурированное логирование (slog)
│   ├── channels/     # Connectors (Telegram, Discord, Web)
│   ├── llm/         # LLM провайдеры + роутинг
│   └── skills/      # Skills execution system
├── pkg/              # Public libraries
├── docs/
│   └── rules/        # Модули правил
├── migrations/       # SQL миграции
├── skills/           # Навыки (SKILL.md)
└── .cursor/rules/    # agent_behavior.mdc
```

## Go-специфика (см. `docs/rules/projectrules.md`)

- Пакеты в `internal/` недоступны извне
- Интерфейсы там, где нужны для тестирования
- Валидация в `Validate()` методе
- Ошибки: `fmt.Errorf("msg: %w", err)`
- Context для всех I/O операций
- Struct tags: `json:"field_name" yaml:"field_name"`

## Конфигурация

- ENV через `${VAR_NAME}` в YAML/JSON
- Валидация после загрузки
- Секреты — через ENV (см. `security.md`)

## База данных

- SQLC для типобезопасных запросов
- Interface Database для мокинга
- Connection pool: 25 max open/idle, 5min lifetime
- Foreign keys включены

## Логирование

- slog (structured logging)
- JSON для production, text для development
- Маскирование секретов (keys: token, key, password, secret)
- Контекст: `logger.Info("msg", "key", value)`

## Тестирование

- Unit тесты для каждого модуля
- Временные файлы: `os.MkdirTemp()`
- Mock интерфейсов
- Табличные тесты

## Безопасность (КРИТИЧНО)

См. `docs/rules/security.md`:
- Никаких секретов в коде
- ENV для секретов
- Маскирование в логах
- Whitelist ресурсов
- Sandbox для опасных навыков
- Валидация входных данных
- Параметризованные запросы (SQLC)

## Issue tracking

```bash
bd ready              # Найти задачи
bd show <id>          # Детали
bd update <id> --status in_progress  # Захватить
bd close <id>         # Завершить
bd sync               # Sync с git
```

## Landing the Plane (Завершение сессии)

**ВСЕГДА выполняйте:**
1. Создайте issues для незавершенного
2. Запустите quality gates
3. Обновите статусы задач
4. **ОБЯЗАТЕЛЬНО git push** (работа не завершена без push!)
5. Очистите stashes/ветки
6. Проверьте `git status` — должно быть "up to date"

**КРИТИЧНО:**
- Работа НЕ завершена до git push
- НЕ говорите "ready to push when you are" — ВЫ push
- Если push не удается — решите и повторите

## Создание компонентов

1. Создайте пакет в `internal/` или `pkg/`
2. Определите интерфейс (если нужен)
3. Реализуйте функции
4. Валидация (если применимо)
5. Напишите тесты
6. Добавьте логирование
7. Обновите конфиг (если нужно)

## Формат SKILL.md

См. `docs/формат SKILL.md` в проекте.

## CI/CD

`.github/workflows/ci.yml`: тесты, сборка, линтеры, security checks.

## Документация

- PRD: `docs/nexflow-prd-v2.md`
- Implementation: `docs/nexflow-implementation-plan.md`
- Database: `docs/database-config.md`
- Logging: `docs/LOGGING_INTEGRATION.md`

---

**Помните:** Начинайте с `docs/rules/security.md` при любых изменениях!
