# Changelog

Все значимые изменения проекта Nexflow будут задокументированы в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Добавлено
- Unit тесты для всех domain entities (User, Task, Skill, Session, Message, Log, Schedule)
- Unit тесты для use cases (UserUseCase, ChatUseCase) с mock repositories
- Integration тесты для repository implementations (UserRepository, SessionRepository, MessageRepository)
- README.md с полноценным описанием проекта и быстрый старт
- Architecture Overview документация
- Layer Architecture документация
- Обновлены AGENTS.md с описанием слоистой архитектуры

### Изменено
- Рефакторинг проекта на Clean Layered Architecture
- Структура проекта: Domain, Application, Infrastructure, Presentation слои
- Реализация repository interfaces в infrastructure слое

### Исправлено
- Дубликаты типов в llm_provider.go
- Unit тесты для domain entities с учетом актуального поведения функций
- Integration тесты с исправлением assertion для session updates

## [0.1.0] - 2026-01-30

### Добавлено
- Initial release of Nexflow
- Domain entities: User, Session, Message, Task, Skill, Schedule, Log
- Repository interfaces и SQLite implementations
- Application use cases: UserUseCase, ChatUseCase
- LLM Provider ports (OpenAI, Anthropic, Ollama)
- Конфигурация через YAML/JSON с ENV переменными
- Структурированное логирование с slog
- SQLC для type-safe SQL запросов

### Изменено
- Разделение на слои: Domain, Application, Infrastructure
- Dependency Inversion через interfaces
- DTOs для data transfer между слоями

---

## Формат changelog

### Категории

- **Добавлено** - новые функциональности
- **Изменено** - изменения в существующей функциональности
- **Устранено** - исправление багов
- **Удалено** - удаление функциональности
- **Безопасность** - уязвимости безопасности
- **Депрекейтед** - функциональность помечена как deprecated

### Схема версионирования

Версии следуют [Semantic Versioning](https://semver.org/lang/ru/):
- **MAJOR.MINOR.PATCH**
- MAJOR: несовместимые изменения API
- MINOR: новая функциональность, обратно совместимая
- PATCH: исправление багов, обратно совместимое

## Дополнительная информация

- [Development Guide](docs/development-guide.md)
- [Testing Guide](docs/testing-guide.md)
- [Architecture Overview](docs/architecture-overview.md)
- [Layer Architecture](docs/layer-architecture.md)
- [Refactoring Guide](docs/refactoring-guide.md)
