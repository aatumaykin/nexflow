# Руководство по миграции базы данных Nexflow

Документ описывает как управлять миграциями базы данных в проекте Nexflow.

## Введение

Nexflow использует два варианта баз данных:
- **SQLite** (по умолчанию) - для разработки и небольших деплоев
- **PostgreSQL** - для продакшна и больших инстансов

## Структура миграций

```
migrations/
├── 000001_init_schema.up.sql
├── 000001_init_schema.down.sql
├── 000002_add_skill_table.up.sql
├── 000002_add_skill_table.down.sql
└── ...
```

## Нейминг файлов

Формат: `{YYYYMMDD}_{description}.{direction}.sql`

- **YYYYMMDD** - дата создания миграции
- **description** - краткое описание изменений
- **direction** - `up` (применить) или `down` (отменить)

## Примеры миграций

### Initial Schema

`migrations/000001_init_schema.up.sql`:
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    channel TEXT NOT NULL,
    channel_user_id TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(channel, channel_user_id)
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    skill TEXT NOT NULL,
    input TEXT NOT NULL,
    output TEXT,
    status TEXT NOT NULL,
    error TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE TABLE skills (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    version TEXT NOT NULL,
    location TEXT NOT NULL,
    permissions TEXT NOT NULL,
    metadata TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE schedules (
    id TEXT PRIMARY KEY,
    skill TEXT NOT NULL,
    cron_expression TEXT NOT NULL,
    input TEXT NOT NULL,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (skill) REFERENCES skills(name) ON DELETE CASCADE
);

CREATE TABLE logs (
    id TEXT PRIMARY KEY,
    level TEXT NOT NULL,
    source TEXT NOT NULL,
    message TEXT NOT NULL,
    metadata TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
```

### Добавление новой таблицы

`migrations/000002_add_indexes.up.sql`:
```sql
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_messages_session_id ON messages(session_id);
CREATE INDEX idx_tasks_session_id ON tasks(session_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_logs_level ON logs(level);
```

`migrations/000002_add_indexes.down.sql`:
```sql
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_messages_session_id;
DROP INDEX IF EXISTS idx_tasks_session_id;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_logs_level;
```

### Добавление колонки

`migrations/000003_add_user_metadata.up.sql`:
```sql
ALTER TABLE users ADD COLUMN metadata TEXT;
```

`migrations/000003_add_user_metadata.down.sql`:
```sql
ALTER TABLE users DROP COLUMN IF EXISTS metadata;
```

## Применение миграций

### Вручную

```bash
# Применить все миграции
go run cmd/migrate/main.go up

# Применить конкретную миграцию
go run cmd/migrate/main.go up 000001_init_schema

# Отменить последнюю миграцию
go run cmd/migrate/main.go down

# Отменить конкретную миграцию
go run cmd/migrate/migrate.go down 000001_init_schema
```

### Автоматически при старте

```go
// internal/shared/database/database.go
func New(cfg *DBConfig) (*DB, error) {
    db, err := sql.Open(cfg.Type, cfg.Path)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Enable foreign keys
    if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
        return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
    }
    
    // Run migrations
    if err := runMigrations(db); err != nil {
        return nil, fmt.Errorf("failed to run migrations: %w", err)
    }
    
    return &DB{db: db}, nil
}

func runMigrations(db *sql.DB) error {
    // Read migration files
    files, err := filepath.Glob("migrations/*.up.sql")
    if err != nil {
        return err
    }
    
    // Sort by filename
    sort.Strings(files)
    
    // Apply each migration
    for _, file := range files {
        content, err := os.ReadFile(file)
        if err != nil {
            return err
        }
        
        if _, err := db.Exec(string(content)); err != nil {
            return fmt.Errorf("failed to apply migration %s: %w", file, err)
        }
        
        log.Info("Migration applied", "file", filepath.Base(file))
    }
    
    return nil
}
```

## Отслеживание версий

### Таблица schema_migrations

```sql
CREATE TABLE schema_migrations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version TEXT NOT NULL UNIQUE,
    applied_at TEXT NOT NULL DEFAULT (datetime('now'))
);
```

### Запись примененных миграций

```go
func applyMigration(db *sql.DB, version string) error {
    // Read migration file
    content, err := os.ReadFile(fmt.Sprintf("migrations/%s.up.sql", version))
    if err != nil {
        return err
    }
    
    // Begin transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Apply migration
    if _, err := tx.Exec(string(content)); err != nil {
        return err
    }
    
    // Record migration
    _, err = tx.Exec(
        `INSERT INTO schema_migrations (version, applied_at) VALUES (?, datetime('now'))`,
        version,
    )
    if err != nil {
        return err
    }
    
    return tx.Commit()
}

func checkMigrationApplied(db *sql.DB, version string) (bool, error) {
    var count int
    err := db.QueryRow(
        `SELECT COUNT(*) FROM schema_migrations WHERE version = ?`,
        version,
    ).Scan(&count)
    
    if err != nil {
        return false, err
    }
    
    return count > 0, nil
}
```

## Rollback

### Отмена миграции

```bash
# Отменить последнюю миграцию
go run cmd/migrate/main.go down

# Отменить конкретную миграцию
go run cmd/migrate/main.go down 000002_add_indexes
```

### Down migration файл

`migrations/000002_add_indexes.down.sql`:
```sql
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_messages_session_id;
DROP INDEX IF EXISTS idx_tasks_session_id;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_logs_level;
```

## Best Practices

### Правила написания миграций

1. **Безопасные изменения** - делайте изменения обратимыми
2. **Данные** - никогда не удаляйте данные в down migrations
3. **Индексы** - создавайте индексы после создания таблиц
4. **Foreign Keys** - используйте внешние ключи для целостности данных
5. **Default Values** - задавайте разумные default значения
6. **NOT NULL** - используйте где применимо
7. **UNIQUE** - для уникальных полей (email, username)

### Разработка

```sql
-- Bad: удаление колонки с данными
ALTER TABLE users DROP COLUMN email;

-- Good: переименование колонки
ALTER TABLE users RENAME COLUMN email TO old_email;
ALTER TABLE users ADD COLUMN email TEXT;
UPDATE users SET email = old_email;
ALTER TABLE users DROP COLUMN old_email;
```

### Тестирование миграций

```go
func TestMigration_000001_init_schema(t *testing.T) {
    db := setupTestDB(t)
    
    // Read up migration
    up, err := os.ReadFile("migrations/000001_init_schema.up.sql")
    require.NoError(t, err)
    
    // Apply migration
    _, err = db.Exec(string(up))
    require.NoError(t, err)
    
    // Verify tables exist
    tables, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
    require.NoError(t, err)
    
    var tableNames []string
    for tables.Next() {
        var name string
        tables.Scan(&name)
        tableNames = append(tableNames, name)
    }
    
    assert.Contains(t, tableNames, "users")
    assert.Contains(t, tableNames, "sessions")
    // ...
}
```

## Продакшен миграции

### Резервное копирование

```bash
# Backup перед миграцией
pg_dump nexflow > backup_$(date +%Y%m%d_%H%M%S).sql
```

### Применение в транзакции

```bash
# Начало транзакции
BEGIN;

-- Применение миграций
-- ...

# Если все успешно
COMMIT;

# Если ошибка
ROLLBACK;
```

## PostgreSQL специфично

### Отличия от SQLite

| Фича | SQLite | PostgreSQL |
|------|--------|------------|
| AUTO_INCREMENT | INTEGER PRIMARY KEY AUTOINCREMENT | SERIAL/BIGSERIAL |
| Boolean | INTEGER | BOOLEAN |
| DateTime | TEXT (ISO8601) | TIMESTAMPTZ |
| JSON | TEXT | JSONB |

### Migration пример для PostgreSQL

`migrations/000001_init_schema.postgres.sql`:
```sql
CREATE TYPE log_level AS ENUM ('debug', 'info', 'warn', 'error');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel VARCHAR(50) NOT NULL,
    channel_user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(channel, channel_user_id)
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
```

## Troubleshooting

### Распространенные проблемы

#### Migration lock

**Проблема:** Два процесса одновременно пытаются применить миграции

**Решение:** Используйте advisory lock

```go
func runMigrations(db *sql.DB) error {
    // Try to acquire lock
    _, err := db.Exec(`SELECT GET_LOCK('migrations', 60)`)
    if err != nil {
        return errors.New("another migration is running")
    }
    defer db.Exec(`SELECT RELEASE_LOCK('migrations')`)
    
    // Apply migrations...
}
```

#### Ошибка SQL syntax

**Проблема:** Несовместимость SQL диалектов

**Решение:** Используйте отдельные файлы для SQLite и PostgreSQL

```
migrations/
├── 000001_init_schema.sqlite.sql
├── 000001_init_schema.postgres.sql
└── ...
```

#### Не применяется миграция

**Проблема:** Migration уже применена

**Решение:** Проверьте таблицу `schema_migrations`

```sql
SELECT * FROM schema_migrations ORDER BY applied_at DESC;
```

## Дополнительная документация

- [Database Configuration](./database-config.md)
- [Development Guide](./development-guide.md)
- [Testing Guide](./testing-guide.md)

---

**Последнее обновление:** Январь 2026
