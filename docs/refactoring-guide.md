# Руководство по рефакторингу Nexflow

Это руководство описывает процессы и лучшие практики для рефакторинга кода в проекте Nexflow.

## Содержание

- [Принципы рефакторинга](#принципы-рефакторинга)
- [Подготовка к рефакторингу](#подготовка-к-рефакторингу)
- [Безопасное рефакторинг](#безопасное-рефакторинг)
- [Типичные сценарии](#типичные-сценарии)
- [Инструменты](#инструменты)
- [Чеклист](#чеклист)

## Принципы рефакторинга

### Когда делать рефакторинг

**✅ ДЕЛАЙТЕ:**
- Перед добавлением новой функциональности ("Make it work, then make it right")
- Когда дублирование кода обнаружено (Rule of Three)
- После добавления тестов (Red-Green-Refactor)
- При выявлении нарушений принципов SOLID
- Когда код сложно понять или поддерживать

**❌ НЕ ДЕЛАЙТЕ:**
- На продакшне без тестов
- Изменив всё сразу (Big Bang refactoring)
- Без сохранения функциональности (не меняйте имена без причин)
- Когда вы не понимаете, что делает код

### Правило Rule of Three

1. **Первый раз** — это "просто сделай" (Just do it)
2. **Второй раз** — подумаешь о дублировании (Deja vu)
3. **Третий раз** — рефакторинг (Refactor)

### Критерии качества кода

Перед рефакторингом убедитесь, что:
- ✅ Тесты покрывают изменяемый код (>80%)
- ✅ Тесты проходят успешно
- ✅ Код форматируется (`go fmt`)
- ✅ Линтер не находит критических ошибок

## Подготовка к рефакторингу

### 1. Понимание контекста

Прежде чем что-то менять:

```bash
# Изучите структуру пакета
tree internal/domain/entity/

# Прочитайте godoc
go doc ./internal/domain/entity

# Изучите использование кода
grep -r "User" ./internal/
```

### 2. Создание ветки

```bash
git checkout -b refactor/feature-name
```

### 3. Запуск тестов

```bash
# Базовая линия тестов
go test ./... -cover

# Сохраните результаты
go test ./... -cover > test-coverage-before.txt
```

### 4. Коммит текущего состояния

```bash
git add .
git commit -m "Starting refactoring: <description>
- Tests passing: 100%
- Coverage: XX%"
```

## Безопасное рефакторинг

### Микро-коммиты

Разделяйте рефакторинг на маленькие, атомарные изменения:

```bash
# 1. Добавьте новое поле (не удаляя старое)
git commit -m "refactor: add new field (BC)"

# 2. Мигрируйте использование
git commit -m "refactor: migrate to new field"

# 3. Удалите старое поле
git commit -m "refactor: remove old field"
```

### Red-Green-Refactor

```go
// 1. RED: Напишите тест для нового поведения
func TestUser_NewFormat(t *testing.T) {
    user := NewUser("telegram", "123")
    assert.Equal(t, "telegram", user.Channel)
}

// 2. GREEN: Сделайте тест проходящим (минимальные изменения)
// 3. REFACTOR: Улучшите код, сохраняя зелёный тест
```

### Feature Flags

Для безопасного деплоя больших изменений:

```go
// internal/shared/config/config.go
type Config struct {
    // ...
    NewBehaviorEnabled bool `yaml:"new_behavior_enabled"`
}

// internal/domain/service.go
func (s *Service) Process(ctx context.Context) error {
    if s.cfg.NewBehaviorEnabled {
        return s.processNew(ctx)
    }
    return s.processOld(ctx)
}
```

## Типичные сценарии

### Сценарий 1: Устранение дублирования

**Проблема:**

```go
// ❌ Дублирование кода
func TaskToDomain(dbTask *dbmodel.Task) *entity.Task {
    createdAt, _ := time.Parse(time.RFC3339, dbTask.CreatedAt)
    updatedAt, _ := time.Parse(time.RFC3339, dbTask.UpdatedAt)
    // ...
}

func UserToDomain(dbUser *dbmodel.User) *entity.User {
    createdAt, _ := time.Parse(time.RFC3339, dbUser.CreatedAt)
    // ...
}
```

**Решение:**

```go
// ✅ Helper функция
// internal/shared/utils/utils.go
func ParseTimeRFC3339(s string) time.Time {
    t, err := time.Parse(time.RFC3339, s)
    if err != nil {
        return time.Time{}
    }
    return t
}

func FormatTimeRFC3339(t time.Time) string {
    return t.Format(time.RFC3339)
}

// Использование
func TaskToDomain(dbTask *dbmodel.Task) *entity.Task {
    return &entity.Task{
        CreatedAt: utils.ParseTimeRFC3339(dbTask.CreatedAt),
        UpdatedAt: utils.ParseTimeRFC3339(dbTask.UpdatedAt),
        // ...
    }
}
```

### Сценарий 2: Выделение интерфейса

**Проблема:**

```go
// ❌ Жёсткая зависимость
type UserService struct {
    db *sql.DB // Жёсткая зависимость от конкретной реализации
}
```

**Решение:**

```go
// ✅ Интерфейс
// internal/domain/repository/user_repository.go
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id string) (*entity.User, error)
}

// internal/application/usecase/user_usecase.go
type UserUseCase struct {
    repo repository.UserRepository // Зависимость от интерфейса
}

// Реализация в Infrastructure
type SQLUserRepository struct {
    db *sql.DB
}
```

### Сценарий 3: Валидация в entity

**Проблема:**

```go
// ❌ Валидация разбросана
func CreateUser(name, email string) error {
    if name == "" {
        return errors.New("name is required")
    }
    // ...
}
```

**Решение:**

```go
// ✅ Валидация в entity
type User struct {
    Name  string
    Email string
}

func (u *User) Validate() error {
    if u.Name == "" {
        return fmt.Errorf("name is required")
    }
    if u.Email == "" {
        return fmt.Errorf("email is required")
    }
    return nil
}

func CreateUser(name, email string) (*User, error) {
    user := &User{Name: name, Email: email}
    if err := user.Validate(); err != nil {
        return nil, err
    }
    return user, nil
}
```

### Сценарий 4: Логика ошибок

**Проблема:**

```go
// ❌ Потеря контекста ошибки
func GetUser(id string) (*User, error) {
    user, err := db.Query("SELECT ... WHERE id = ?", id)
    if err != nil {
        return nil, err // Потерян контекст
    }
    return user, nil
}
```

**Решение:**

```go
// ✅ Обёртывание ошибок
func GetUser(id string) (*User, error) {
    user, err := db.Query("SELECT ... WHERE id = ?", id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return user, nil
}
```

### Сценарий 5: Константы вместо magic numbers

**Проблема:**

```go
// ❌ Magic numbers
if len(name) > 100 {
    return errors.New("name too long")
}
```

**Решение:**

```go
// ✅ Константы
const (
    MaxNameLength = 100
    MinNameLength = 1
)

func ValidateName(name string) error {
    if len(name) < MinNameLength {
        return fmt.Errorf("name too short (min %d)", MinNameLength)
    }
    if len(name) > MaxNameLength {
        return fmt.Errorf("name too long (max %d)", MaxNameLength)
    }
    return nil
}
```

## Инструменты

### Статический анализ

```bash
# go vet
go vet ./...

# golangci-lint
golangci-lint run --timeout 5m

# go fmt
go fmt ./...

# goimports
goimports -w .
```

### Тестирование

```bash
# Базовое тестирование
go test ./...

# С покрытием
go test -cover ./...

# С race detector
go test -race ./...

# Бенчмарки
go test -bench=. ./...

# Профилирование
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...
```

### Визуализация зависимостей

```bash
# Graph зависимостей
go mod graph | grep "nexflow" | dot -Tpng > deps.png

# Циклические зависимости
go mod graph | grep -v std | grep -v "nexflow" | grep -v "github.com"
```

### Refactoring в IDE

**VS Code / GoLand:**

1. **Rename Symbol** (F2) — безопасное переименование
2. **Extract Method** — выделение метода
3. **Extract Interface** — выделение интерфейса
4. **Inline Variable** — встраивание переменной
5. **Safe Delete** — безопасное удаление (проверка использования)

## Чеклист

### До рефакторинга

- [ ] Создана отдельная ветка
- [ ] Все тесты проходят
- [ ] Покрытие кода >80%
- [ ] Понято текущее поведение кода
- [ ] Создан базовый коммит

### Во время рефакторинга

- [ ] Изменения разбиты на микро-коммиты
- [ ] После каждого коммита тесты проходят
- [ ] Не смешиваются функциональные изменения с рефакторингом
- [ ] godoc comments обновлены
- [ ] Обновлена документация (если нужно)

### После рефакторинга

- [ ] Все тесты проходят
- [ ] Покрытие кода не упало
- [ ] go vet не находит ошибок
- [ ] golangci-lint проходит
- [ ] go fmt применён
- [ ] Бенчмарки не ухудшились
- [ ] Документация обновлена
- [ ] Code review пройден

### Перед мержем

```bash
# Финальная проверка
go test ./...
go vet ./...
golangci-lint run
go fmt ./...
go build ./...
```

## Пример полного процесса

### Шаг 1: Подготовка

```bash
# Создаём ветку
git checkout -b refactor/user-validation

# Запускаем тесты
go test ./... -cover > before-coverage.txt

# Коммитим текущее состояние
git add .
git commit -m "Before refactoring: user validation"
```

### Шаг 2: Добавляем валидацию в entity

```bash
# Изменяем internal/domain/entity/user.go
git add internal/domain/entity/user.go
git commit -m "refactor: add validation to User entity"
```

### Шаг 3: Обновляем use cases

```bash
# Обновляем internal/application/usecase/user_usecase.go
git add internal/application/usecase/
git commit -m "refactor: use User.Validate() in use cases"
```

### Шаг 4: Запускаем тесты

```bash
go test ./... -cover > after-coverage.txt

# Сравниваем покрытие
diff before-coverage.txt after-coverage.txt
```

### Шаг 5: Финальный коммит

```bash
# Форматируем код
go fmt ./...

# Запускаем линтер
golangci-lint run

# Коммитим все изменения
git add .
git commit -m "refactor: completed user validation refactoring

- Added Validate() method to User entity
- Updated all use cases to use Validate()
- Maintained 100% test coverage
- All tests passing"
```

## Рекомендации

1. **Не меняйте всё сразу** — делайте маленькие, безопасные изменения
2. **Тестируйте каждый шаг** — не переходите к следующему, пока текущий не работает
3. **Документируйте изменения** — обновляйте godoc comments
4. **Используйте инструменты** — let tools help you (IDE, lint, formatting)
5. **Просите code review** — вторая пара глаз всегда полезна

## Ресурсы

- [Refactoring.Guru](https://refactoring.guru/)
- [Clean Code by Robert C. Martin](https://www.oreilly.com/library/view/clean-code-a/9780136083238/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)
