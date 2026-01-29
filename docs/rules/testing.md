# Тестирование для Nexflow

Этот модуль описывает подход к тестированию, включая unit, integration и e2e тесты.

## Стратегия тестирования

### Три уровня

```
E2E Tests (5-10%)    — полные сценарии через public API
    ↓
Integration Tests (20-30%) — тесты интеграций (БД, внешние API)
    ↓
Unit Tests (60-75%)   — тесты отдельных функций и методов
```

### Целевое покрытие

- Unit тесты: > 70%
- Integration тесты: критические пути
- E2E тесты: основные user flows

## Unit тесты

### Структура файла тестов

```go
package config

import (
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
    // Arrange (подготовка)
    tmpDir := t.TempDir()
    configPath := tmpDir + "/config.yml"

    // Act (выполнение)
    config, err := Load(configPath)

    // Assert (проверка)
    require.NoError(t, err)
    assert.NotNil(t, config)
}
```

### Табличные тесты

Используйте для множественных сценариев.

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid config",
            config: &Config{
                Server: ServerConfig{Host: "localhost", Port: 8080},
                Database: DatabaseConfig{Type: "sqlite", Path: "test.db"},
            },
            wantErr: false,
        },
        {
            name: "missing host",
            config: &Config{
                Server: ServerConfig{Host: "", Port: 8080},
            },
            wantErr: true,
            errMsg:  "server.host is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()

            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

### Использование testify

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// require — останавливает тест если неудача
require.NotNil(t, config)
require.NoError(t, err)

// assert — продолжает тест если неудача
assert.Equal(t, "expected", actual)
assert.True(t, condition)
```

### Mock интерфейсов

```go
// Интерфейс
type Logger interface {
    Info(msg string, args ...any)
    Error(msg string, args ...any)
}

// Mock
type MockLogger struct {
    messages []string
}

func (m *MockLogger) Info(msg string, args ...any) {
    m.messages = append(m.messages, msg)
}

// Тест с mock
func TestProcess_WithMockLogger(t *testing.T) {
    mockLogger := &MockLogger{}
    processor := NewProcessor(mockLogger)

    processor.Process("test")

    assert.Len(t, mockLogger.messages, 1)
}
```

## Integration тесты

### Тестирование с БД

Используйте временную БД.

```go
func TestDatabase_Integration(t *testing.T) {
    // Arrange — временная БД
    tmpDir, _ := os.MkdirTemp("", "nexflow-test-*")
    defer os.RemoveAll(tmpDir)

    dbPath := tmpDir + "/test.db"

    db, _ := sql.Open("sqlite3", dbPath)
    defer db.Close()

    _, _ = db.Exec("PRAGMA foreign_keys = ON")

    queries := New(db)
    testDB := &DB{Queries: queries, db: db}

    ctx := context.Background()
    testDB.Migrate(ctx)

    // Act — создаем пользователя
    user, err := testDB.CreateUser(ctx, CreateUserParams{
        ID: "user-1", Name: "Test", Email: "test@example.com",
        CreatedAt: time.Now(), UpdatedAt: time.Now(),
    })
    require.NoError(t, err)

    // Assert
    assert.Equal(t, "user-1", user.ID)
}
```

### Тестирование с внешними API

Используйте mock серверы.

```go
func TestLLMProvider_Integration(t *testing.T) {
    // Arrange — mock сервер
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"response": "test response"}`))
    }))
    defer server.Close()

    provider := NewOpenAIProvider(Config{
        BaseURL: server.URL,
        APIKey:  "test-key",
    })

    // Act
    response, err := provider.CreateCompletion(ctx, CompletionRequest{
        Model: "gpt-4", Prompt: "test prompt",
    })

    // Assert
    require.NoError(t, err)
    assert.Equal(t, "test response", response.Response)
}
```

## E2E тесты

### Пример E2E теста

```go
func TestE2E_ChatFlow(t *testing.T) {
    // Arrange — тестовый сервер
    config := loadTestConfig(t)
    server, _ := NewServer(config)
    defer server.Close()

    client := server.Client()

    // Act 1 — создаем пользователя
    createUserResp, _ := client.Post(server.URL+"/api/v1/users", "application/json", strings.NewReader(`{
        "name": "Test User",
        "email": "test@example.com"
    }`))
    assert.Equal(t, http.StatusCreated, createUserResp.StatusCode)

    var user User
    json.NewDecoder(createUserResp.Body).Decode(&user)
    createUserResp.Body.Close()

    // Act 2 — отправляем сообщение
    chatResp, _ := client.Post(server.URL+"/api/v1/chat", "application/json", strings.NewReader(`{
        "message": "test message",
        "user_id": "`+user.ID+`",
        "channel": "web"
    }`))
    assert.Equal(t, http.StatusOK, chatResp.StatusCode)

    var chatRespData ChatResponse
    json.NewDecoder(chatResp.Body).Decode(&chatRespData)
    chatResp.Body.Close()

    // Assert
    assert.NotEmpty(t, chatRespData.SessionID)
    assert.NotEmpty(t, chatRespData.Response)
}
```

## Тестирование concurrently

### Race conditions

```bash
go test -race ./...
```

### Тесты с конкурентностью

```go
func TestConcurrentAccess(t *testing.T) {
    cache := NewCache()

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            cache.Set(fmt.Sprintf("key%d", n), n)
        }(i)
    }
    wg.Wait()

    // Assert — все ключи сохранены
    for i := 0; i < 100; i++ {
        val, ok := cache.Get(fmt.Sprintf("key%d", i))
        assert.True(t, ok)
        assert.Equal(t, i, val)
    }
}
```

## Бенчмарки

### Написание бенчмарков

```go
func BenchmarkCache_Get(b *testing.B) {
    cache := NewCache()
    cache.Set("key", "value")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        cache.Get("key")
    }
}
```

### Запуск бенчмарков

```bash
go test -bench=. ./...           # Все бенчмарки
go test -bench=. -benchmem ./... # С памятью
go test -bench=. -cpuprofile=cpu.prof ./... # С профилированием
go tool pprof cpu.prof
```

## Фикстуры и helpers

```go
// Создает временную БД
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()

    tmpDir, _ := os.MkdirTemp("", "nexflow-test-*")
    dbPath := tmpDir + "/test.db"

    db, _ := sql.Open("sqlite3", dbPath)
    _, _ = db.Exec("PRAGMA foreign_keys = ON")

    t.Cleanup(func() {
        db.Close()
        os.RemoveAll(tmpDir)
    })

    return db
}

// Создает тестовый конфиг
func setupTestConfig(t *testing.T) *Config {
    t.Helper()

    tmpDir := t.TempDir()
    configPath := tmpDir + "/config.yml"

    data := []byte(`server:
  host: "localhost"
  port: 8080
database:
  type: "sqlite"
  path: "test.db"`)
    _ = os.WriteFile(configPath, data, 0644)

    config, _ := Load(configPath)
    return config
}
```

## Покрытие кода

### Генерация отчёта

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out  # Проверка в терминале
go tool cover -html=coverage.out   # HTML отчёт
```

### Целевые показатели

- Общее покрытие: > 60%
- Критический код (БД, LLM): > 80%
- Утилитные функции: > 90%

### Исключения

- `main.go` — точки входа
- Generated code (SQLC, protobuf)
- Структуры данных (plain structs)

## Запуск тестов

### Локально

```bash
go test ./...                    # Все тесты
go test -v ./internal/database   # Тесты пакета
go test -race ./...              # С race detection
go test -cover ./...              # С покрытием
go test -timeout 30s ./...        # С таймаутом
```

### CI/CD

**.github/workflows/ci.yml:**
```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25.5'

      - name: Run tests
        run: go test -race -coverprofile=coverage.out ./...

      - name: Check coverage
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$COVERAGE < 60" | bc -l) )); then
            echo "Coverage is below 60%"
            exit 1
          fi
```

## Правила тестирования

### Всегда тестируйте

1. Unit тесты — для всех public функций/методов
2. Error paths — тестируйте не только success paths
3. Edge cases — пустые строки, nil значения, границы
4. Integration — для критических путей с БД и API
5. E2E — для основных user flows

### Никогда не тестируйте

1. Private функции — тестируйте через public API
2. External libs — доверяйте авторам библиотек
3. Generated code — SQLC, protobuf и т.д.
4. Trivial code — геттеры/сеттеры, простые структуры

## Критические правила

1. ВСЕГДА пишите unit тесты для новой функциональности
2. НИКОГДА не игнорируйте ошибки в тестах (используйте require)
3. ВСЕГДА очищайте ресурсы (defer, Cleanup)
4. НИКОГДА не зависите от внешних сервисов в unit тестах
5. ВСЕГДА используйте временные директории/файлы (t.TempDir)
6. НИКОГДА не тестируйте private функции напрямую
7. ВСЕГДА тестируйте error paths
8. НИКОГДА не делайте flaky тесты (нестабильные)
9. ВСЕГДА используйте табличные тесты для множественных сценариев
10. НИКОГДА не запускайте тесты с race detection если они слишком медленные

---

**Памятка:** Тесты — это документация и страховка кода. Инвестируйте время в тесты.
