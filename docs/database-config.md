# Пример конфигурации базы данных

## SQLite (для разработки)

```yaml
database:
  type: "sqlite"
  path: "./data/nexflow.db"
```

## PostgreSQL (для продакшена)

```yaml
database:
  type: "postgres"
  path: "postgres://user:password@localhost:5432/nexflow?sslmode=disable"
```

## Полный пример конфигурации

```yaml
database:
  type: "sqlite"  # или "postgres"
  path: "./data/nexflow.db"

# Для PostgreSQL:
# database:
#   type: "postgres"
#   path: "postgres://user:password@localhost:5432/nexflow?sslmode=disable"
#
# Для подключения через URL:
# postgres://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
#
# Примеры:
# - postgres://user:pass@localhost:5432/nexflow
# - postgres://user@localhost:5432/nexflow?sslmode=require
# - postgres:///nexflow (использует Unix domain socket)
```

## Использование в коде

```go
import (
    "github.com/atumaikin/nexflow/internal/config"
    "github.com/atumaikin/nexflow/internal/database"
)

// Загрузка конфигурации
cfg, err := config.Load("config.yaml")
if err != nil {
    log.Fatal(err)
}

// Создание подключения к базе данных
db, err := database.NewDatabase(&cfg.Database)
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Запуск миграций
ctx := context.Background()
if err := db.Migrate(ctx); err != nil {
    log.Fatal(err)
}
```

## Примечания

### SQLite
- Подходит для локальной разработки
- Не требует отдельного сервера
- Файл базы данных будет создан автоматически
- Поддержка транзакций и foreign keys включена по умолчанию

### PostgreSQL
- Рекомендуется для продакшена
- Поддерживает параллельные операции
- Лучше масштабируется
- Требует запущенного PostgreSQL сервера

## Переменные окружения

Вы можете использовать переменные окружения в конфигурации:

```yaml
database:
  type: "postgres"
  path: "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
```

Пример запуска:

```bash
export DB_USER=nexflow
export DB_PASSWORD=secret
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=nexflow_db
export DB_SSLMODE=disable

./nexflow
```
