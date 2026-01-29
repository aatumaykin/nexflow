# Руководство по тестированию Nexflow

Это руководство описывает лучшие практики и подходы к тестированию в проекте Nexflow.

## Содержание

- [Философия тестирования](#философия-тестирования)
- [Типы тестов](#типы-тестов)
- [Структура тестов](#структура-тестов)
- [Unit тесты](#unit-тесты)
- [Integration тесты](#integration-тесты)
- [Моки и стабы](#моки-и-стабы)
- [Тестирование базы данных](#тестирование-базы-данных)
- [Бенчмарки](#бенчмарки)
- [CI/CD](#cicd)

## Философия тестирования

### Принципы

1. **Pyramid Testing** — больше unit тестов, меньше integration/e2e тестов
2. **Fast Feedback** — тесты должны выполняться быстро (<5 мин для всего проекта)
3. **Isolation** — каждый тест должен быть независимым
4. **Readability** — тесты должны быть читаемыми и понятными
5. **Maintainability** — тесты должны легко поддерживаться

### Золотое правило

> "Тестируйте поведение, а не реализацию."

**❌ Плохо (тестирование реализации):**
```go
func TestUser_Create(t *testing.T) {
    user := NewUser("telegram", "123")
    assert.Equal(t, "usr_", user.ID[:4]) // Тестирует деталь реализации ID
}
```

**✅ Хорошо (тестирование поведения):**
```go
func TestUser_Create(t *testing.T) {
    user := NewUser("telegram", "123")
    assert.NotEmpty(t, user.ID) // Тестирует что ID не пустой
    assert.Equal(t, "telegram", user.Channel)
    assert.Equal(t, "123", user.ChannelID)
}
```

## Типы тестов

| Тип          | Цель                         | Скорость  | Пример               |
|--------------|------------------------------|-----------|----------------------|
| Unit         | Ликальная логика              | <1ms      | Валидация entity     |
| Integration  | Взаимодействие компонентов    | <100ms    | Repository + БД      |
| E2E          | Пользовательские сценарии     | >1s       | HTTP API             |
| Benchmark    | Производительность            | >10ms     | Сортировка списков   |

## Структура тестов

### Именование файлов

```bash
# Тестовый файл всегда рядом с исходным
internal/domain/entity/user.go
internal/domain/entity/user_test.go
```

### Именование функций

```go
// Формат: Test<TypeName>_<FunctionName>_<Scenario>
func TestUser_Create_Success(t *testing.T) {}
func TestUser_Create_InvalidInput(t *testing.T) {}
func TestUser_GetByID_NotFound(t *testing.T) {}
```

## Unit тесты

### Базовая структура

```go
package entity

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
    tests := []struct {
        name      string
        channel   string
        channelID string
        want      *User
    }{
        {
            name:      "telegram user",
            channel:   "telegram",
            channelID: "123",
            want: func() *User {
                u := &User{
                    Channel:   "telegram",
                    ChannelID: "123",
                }
                // ID и CreatedAt генерируются динамически
                u.ID = u.ID // Placeholder
                return u
            }(),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := NewUser(tt.channel, tt.channelID)

            assert.Equal(t, tt.channel, got.Channel)
            assert.Equal(t, tt.channelID, got.ChannelID)
            assert.NotEmpty(t, got.ID)
            assert.NotZero(t, got.CreatedAt)
        })
    }
}
```

### Тестирование ошибок

```go
func TestUser_Validate(t *testing.T) {
    tests := []struct {
        name    string
        user    *User
        wantErr bool
    }{
        {
            name:    "valid user",
            user:    NewUser("telegram", "123"),
            wantErr: false,
        },
        {
            name:    "empty channel",
            user:    NewUser("", "123"),
            wantErr: true,
        },
        {
            name:    "empty channel ID",
            user:    NewUser("telegram", ""),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Тестирование методов

```go
func TestUser_IsSameChannel(t *testing.T) {
    user1 := NewUser("telegram", "123")
    user2 := NewUser("telegram", "456")
    user3 := NewUser("discord", "789")

    assert.True(t, user1.IsSameChannel(user2))
    assert.False(t, user1.IsSameChannel(user3))
}
```

## Integration тесты

### Тестирование с реальной БД

```go
package database_test

import (
    "context"
    "os"
    "testing"
    "github.com/atumaikin/nexflow/internal/domain/entity"
    dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
)

func TestUserRepository_Create(t *testing.T) {
    // Создаём временную БД
    tmpDB, err := os.CreateTemp("", "test-*.db")
    require.NoError(t, err)
    defer os.Remove(tmpDB.Name())

    // Инициализируем repository
    repo := NewSQLiteRepository(tmpDB.Name())
    defer repo.Close()

    // Запускаем миграции
    require.NoError(t, repo.Migrate(context.Background()))

    // Создаём пользователя
    user := entity.NewUser("telegram", "123")
    err = repo.Create(context.Background(), user)
    require.NoError(t, err)

    // Проверяем, что пользователь создан
    found, err := repo.GetByID(context.Background(), user.ID)
    require.NoError(t, err)
    assert.Equal(t, user.ID, found.ID)
    assert.Equal(t, user.Channel, found.Channel)
}
```

### Использование Testify для assertions

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
    // assert — продолжает выполнение даже при ошибке
    assert.Equal(t, expected, actual)
    assert.NoError(t, err)
    assert.True(t, condition)
    assert.Contains(t, slice, value)

    // require — останавливает тест при ошибке
    require.NoError(t, err)
    require.NotNil(t, value)
    require.Len(t, slice, expectedLen)
}
```

## Моки и стабы

### Использование testify/mock

```go
package usecase_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock реализации
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

// Тест с mock
func TestUserUseCase_Create(t *testing.T) {
    tests := []struct {
        name    string
        setup   func(*MockUserRepository)
        input   string
        wantErr bool
    }{
        {
            name: "success",
            setup: func(repo *MockUserRepository) {
                repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).
                    Return(nil)
            },
            input:   "test-user",
            wantErr: false,
        },
        {
            name: "repository error",
            setup: func(repo *MockUserRepository) {
                repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).
                    Return(assert.AnError)
            },
            input:   "test-user",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockUserRepository)
            tt.setup(mockRepo)

            uc := NewUserUseCase(mockRepo)
            _, err := uc.Create(context.Background(), tt.input)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

### Создание стабов (stubs)

```go
// Stub для Logger
type NoopLogger struct{}

func (l *NoopLogger) Info(msg string, args ...any)               {}
func (l *NoopLogger) Error(msg string, args ...any)              {}
func (l *NoopLogger) With(args ...any) logging.Logger            { return l }
func (l *NoopLogger) WithContext(ctx context.Context) logging.Logger { return l }

// Использование в тестах
func TestService_DoSomething(t *testing.T) {
    service := NewService(&NoopLogger{})
    // ...
}
```

## Тестирование базы данных

### Временная БД для тестов

```go
func setupTestDB(t *testing.T) (*sql.DB, func()) {
    // Создаём временную БД
    tmpDB, err := os.CreateTemp("", "test-*.db")
    require.NoError(t, err)

    // Открываем соединение
    db, err := sql.Open("sqlite", tmpDB.Name())
    require.NoError(t, err)

    // Возвращаем cleanup функцию
    cleanup := func() {
        db.Close()
        os.Remove(tmpDB.Name())
    }

    return db, cleanup
}

func TestMigrations(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    // Запускаем миграции
    err := Migrate(db, "migrations/sqlite")
    require.NoError(t, err)

    // Проверяем, что таблицы созданы
    var tables []string
    rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
    require.NoError(t, err)
    defer rows.Close()

    for rows.Next() {
        var name string
        rows.Scan(&name)
        tables = append(tables, name)
    }

    assert.Contains(t, tables, "users")
    assert.Contains(t, tables, "sessions")
}
```

### Тестирование mappers

```go
func TestUserMapper_ToDomain(t *testing.T) {
    now := time.Now()

    dbUser := &dbmodel.User{
        ID:            "test-id",
        Channel:       "telegram",
        ChannelUserID: "123",
        CreatedAt:     now.Format(time.RFC3339),
    }

    user := UserToDomain(dbUser)

    assert.Equal(t, dbUser.ID, user.ID)
    assert.Equal(t, dbUser.Channel, user.Channel)
    assert.Equal(t, dbUser.ChannelUserID, user.ChannelID)
    assert.WithinDuration(t, now, user.CreatedAt, time.Second)
}

func TestUserMapper_ToDB(t *testing.T) {
    now := time.Now()

    user := &entity.User{
        ID:        "test-id",
        Channel:   "telegram",
        ChannelID: "123",
        CreatedAt: now,
    }

    dbUser := UserToDB(user)

    assert.Equal(t, user.ID, dbUser.ID)
    assert.Equal(t, user.Channel, dbUser.Channel)
    assert.Equal(t, user.ChannelID, dbUser.ChannelUserID)
}
```

## Бенчмарки

### Создание бенчмарков

```go
func BenchmarkUser_New(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewUser("telegram", "123")
    }
}

func BenchmarkUser_Validate(b *testing.B) {
    user := NewUser("telegram", "123")
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        user.Validate()
    }
}
```

### Запуск бенчмарков

```bash
# Запуск всех бенчмарков
go test -bench=. ./...

# С детализацией памяти
go test -bench=. -benchmem ./...

# С конкретным фильтром
go test -bench=BenchmarkUser ./...
```

### Анализ профилей

```bash
# CPU профиль
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof

# Memory профиль
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof
```

## CI/CD

### GitHub Actions

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.25.5', '1.26.0']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test -v -race -cover ./...

      - name: Run benchmarks
        run: go test -bench=. -run=^$ ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
```

### Локальный CI

```bash
# Запуск всех проверок
make test-all

# Makefile
.PHONY: test-all
test-all:
	go test ./...
	go vet ./...
	golangci-lint run
	go fmt ./...
```

## Полезные техники

### Тестирование с таблицами (Table-Driven Tests)

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "case 1",
            input:    "input1",
            expected: "output1",
            wantErr:  false,
        },
        {
            name:     "case 2",
            input:    "input2",
            expected: "output2",
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    "",
            expected: "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := DoSomething(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}
```

### Тестирование с context

```go
func TestService_DoSomething(t *testing.T) {
    ctx := context.Background()

    service := NewService(ctx, mockRepo, logger)
    result, err := service.DoSomething(ctx)

    assert.NoError(t, err)
    assert.NotNil(t, result)
}

func TestService_DoSomething_Timeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer cancel()

    service := NewService(ctx, slowRepo, logger)
    _, err := service.DoSomething(ctx)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "timeout")
}
```

### Тестирование concurrency

```go
func TestCache_Concurrent(t *testing.T) {
    cache := NewCache()

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            cache.Set(fmt.Sprintf("key%d", n), fmt.Sprintf("value%d", n))
        }(i)
    }

    wg.Wait()

    assert.Equal(t, 100, cache.Len())
}
```

## Чеклист для тестов

### Unit тесты

- [ ] Тестирует поведение, а не реализацию
- [ ] Не зависит от порядка выполнения
- [ ] Использует описательные имена
- [ ] Покрывает все ветки логики
- [ ] Изолирован от других тестов

### Integration тесты

- [ ] Использует реальные зависимости (БД, API)
- [ ] Создаёт и удаляет тестовые данные
- [ ] Работает быстро (<5 сек)
- [ ] Помечен тегом `//go:build integration`

### Code review

- [ ] Тесты читаемы и понятны
- [ ] Покрытие кода >80%
- [ ] Нет race conditions
- [ ] Бенчмарки не ухудшились

## Ресурсы

- [Go Testing](https://golang.org/pkg/testing/)
- [Testify](https://github.com/stretchr/testify)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [Go Concurrency Patterns](https://blog.golang.org/pipelines)
