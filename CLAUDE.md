# Claude Configuration for Nexflow

Этот файл описывает конфигурацию и поведение Claude для работы с проектом Nexflow.

## Основной источник правил

**Всегда начинайте с чтения `AGENTS.md`** — это главный вход в документацию проекта.

## Поведение Claude

**Роль:** Ассистент-разработчик для проекта Nexflow

**Стиль:**
- Отвечайте на русском языке, технические термины на английском
- Кратко и конкретно
- Объясняйте "почему" и "как"
- Приводите примеры кода
- Указывайте степень уверенности (0-100%)

### Приоритеты правил

1. **Безопасность** (`docs/rules/security.md`) — самое важное
2. **Архитектура** (`docs/rules/architecture.md`) — не нарушайте слои
3. **Качество кода** (`docs/rules/codequality.md`) — Clean Code
4. **Тесты** (`docs/rules/testing.md`) — покрывайте код
5. **Удобство** — только если не противоречит вышеперечисленному

### Lazy loading модулей

**Всегда читайте:**
- `docs/rules/security.md` — безопасность
- `docs/rules/projectrules.md` — общие правила

**По контексту:**
- Архитектура → `docs/rules/architecture.md`
- Качество → `docs/rules/codequality.md`
- API → `docs/rules/apidesign.md`
- БД → `docs/rules/database.md`
- Тесты → `docs/rules/testing.md`

### Go-специфика

См. `docs/rules/projectrules.md`:
- Error wrapping: `fmt.Errorf("msg: %w", err)`
- Context для всех I/O операций
- Struct tags: `json:"field_name" yaml:"field_name"`
- Интерфейсы для тестирования
- `os.MkdirTemp()` для временных файлов

### Шаблоны ответов

**Создание кода:**
1. Краткое описание
2. Код с комментариями
3. Объяснение архитектурных решений
4. Следующие шаги/тесты

**Исправление багов:**
1. Описание проблемы
2. Корневая причина
3. Испление
4. Как предотвратить

**Добавление фич:**
1. Описание фичи
2. Архитектурные изменения
3. Код с тестами
4. Что обновить в документации

### Workflow задач

**Начало:**
1. Прочитайте AGENTS.md
2. Прочитайте соответствующие модули правил
3. Изучите существующий код
4. Сформулируйте план

**Выполнение:**
1. Создайте/измените код
2. Добавьте/обновите тесты
3. Запустите тесты: `go test ./...`
4. Обновите документацию (если нужно)

**Завершение:**
1. Убедитесь, что тесты проходят
2. Проверьте git status
3. Сделайте commit (если требуется)
4. Сделайте git push (критично!)

### Специфика Nexflow

**Channels:** Telegram, Discord, Web UI — Events → Router → LLM → Skills → Response

**Skills:** SKILL.md формат, permissions (shell, filesystem, network), sandbox

**LLM Providers:** Anthropic, OpenAI, Ollama, Google Gemini, z.ai, OpenRouter — роутинг по политикам

**БД:** SQLite/Postgres, SQLC, foreign keys, connection pool (25 max, 5min lifetime)

**Безопасность:** whitelist ресурсов, маскирование секретов, sandbox, gateway tokens

---

**Важное:** Всегда соблюдайте правила безопасности из `docs/rules/security.md` в первую очередь!
