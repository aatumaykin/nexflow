# Качество кода для Nexflow

Этот модуль описывает принципы качества кода, включая Clean Code, SOLID, DRY, KISS, YAGNI и антипаттерны.

## Основные принципы

### DRY (Don't Repeat Yourself)

Повторяющаяся логика → отдельные функции/методы.

**✅ ХОРОШО:**
```go
func validateUser(user User) error {
    if user.Name == "" {
        return fmt.Errorf("name is required")
    }
    return nil
}

func processUser1(user User) error {
    if err := validateUser(user); err != nil {
        return err
    }
    return nil
}
```

**❌ ПЛОХО:**
```go
func processUser1(user User) error {
    if user.Name == "" {
        return fmt.Errorf("name is required")
    }
    return nil
}

func processUser2(user User) error {
    if user.Name == "" {
        return fmt.Errorf("name is required")
    }
    return nil
}
```

### KISS (Keep It Simple, Stupid)

Пишите простой, понятный код. Избегайте излишней сложности.

### YAGNI (You Aren't Gonna Need It)

Не пишите код, который вам не нужен сейчас. Избегайте преждевременной абстракции.

**✅ ХОРОШО:**
```go
type Database interface {
    CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
    // ... только методы, которые мы используем
}
```

**❌ ПЛОХО:**
```go
type DatabaseDriver interface {
    Connect() error
    Query(query string) Result
    Execute(query string) error
    Migrate(migration Migration) error
    Backup(path string) error
    Restore(path string) error
    // ... десятки методов, которые нам не нужны
}
```

## SOLID принципы

### Single Responsibility Principle (SRP)

Каждая функция/метод имеет только одну причину для изменения.

**❌ ПЛОХО:**
```go
func (u *User) Save() error { ... }      // БД
func (u *User) SendEmail() error { ... } // Email
func (u *User) Validate() error { ... }  // Валидация
```

**✅ ХОРОШО:**
```go
type UserRepository interface { Save(ctx context.Context, user User) error }
type EmailService interface { SendWelcomeEmail(user User) error }
func ValidateUser(user User) error { ... }
```

### Open/Closed Principle (OCP)

Сущности открыты для расширения, закрыты для модификации.

**❌ ПЛОХО:**
```go
func (p *Processor) Process(dataType string, data interface{}) error {
    switch dataType {
    case "json": return p.processJSON(data)
    case "xml":  return p.processXML(data)
    // Нужно модифицировать для добавления новых форматов!
    }
}
```

**✅ ХОРОШО:**
```go
type Processor interface {
    Process(data interface{}) error
    SupportedType() string
}

type ProcessorRegistry struct {
    processors map[string]Processor
}

func (r *ProcessorRegistry) Register(processor Processor) {
    r.processors[processor.SupportedType()] = processor
}
```

### Liskov Substitution Principle (LSP)

Подтипы взаимозаменяемы с базовыми типами.

**❌ ПЛОХО:**
```go
type Bird interface { Fly() error }
type Penguin struct { ... }
func (p *Penguin) Fly() error { return fmt.Errorf("penguins can't fly") }
```

**✅ ХОРОШО:**
```go
type Bird interface { Move() error }
type FlyingBird interface { Bird; Fly() error }
type Penguin struct { ... }
func (p *Penguin) Move() error { return p.Swim() }
```

### Interface Segregation Principle (ISP)

Клиенты не зависят от интерфейсов, которые не используют.

**❌ ПЛОХО (огромный интерфейс):**
```go
type UserService interface {
    CreateUser(user User) error
    GetUserByID(id string) (User, error)
    UpdateUser(user User) error
    // ... десятки других методов
}
```

**✅ ХОРОШО (маленькие интерфейсы):**
```go
type UserReader interface { GetUserByID(id string) (User, error) }
```

### Dependency Inversion Principle (DIP)

Модули верхних уровней зависят от абстракций.

**❌ ПЛОХО:**
```go
type UserService struct {
    db     *sql.DB
    email  *EmailClient
    logger *FileLogger
}
```

**✅ ХОРОШО:**
```go
type UserService struct {
    db     Database
    email  EmailService
    logger Logger
}
```

## Clean Code

### Именование

**Переменные и функции:**
```go
// ✅ Хорошо
userName := "john"
func getUserByID(id string) (User, error) { ... }

// ❌ Плохо
n := "john"
func get(x string) (User, error) { ... }
```

**Булевы переменные:**
```go
// ✅ Хорошо
isValid := true
hasPermission := false

// ❌ Плохо
valid := true
permission := false
```

**Константы:**
```go
// ✅ Хорошо
const MaxRetries = 3

// ❌ Плохо
const MAX_RETRIES = 3
```

### Функции

**Краткость (< 50 строк)**

**Параметры (максимум 3-4):**
```go
// ✅ Хорошо
func CreateUser(name, email string) (User, error) { ... }

// ✅ Хорошо — структура для группировки
func CreateUser(ctx context.Context, req CreateUserRequest) (User, error) { ... }
```

## Ошибки

### Error wrapping

**✅ ХОРОШО:**
```go
func ProcessUser(id string) error {
    user, err := db.GetUserByID(id)
    if err != nil {
        return fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return nil
}
```

**❌ ПЛОХО:**
```go
func ProcessUser(id string) error {
    user, err := db.GetUserByID(id)
    if err != nil {
        return err // Потерян контекст
    }
    return nil
}
```

## Антипаттерны

### Magic Numbers

**❌ ПЛОХО:**
```go
if len(data) > 100 { ... }
time.Sleep(5 * time.Second)
```

**✅ ХОРОШО:**
```go
const MaxDataSize = 100
if len(data) > MaxDataSize { ... }

const RequestTimeout = 5 * time.Second
time.Sleep(RequestTimeout)
```

### God Functions

**❌ ПЛОХО:** 500 строк кода, всё подряд

**✅ ХОРОШО:** разбейте на маленькие функции

### Deep Nesting

**❌ ПЛОХО:**
```go
if user != nil {
    if user.Email != "" {
        if isValidEmail(user.Email) {
            // ...
        }
    }
}
```

**✅ ХОРОШО (guard clauses):**
```go
if user == nil {
    return fmt.Errorf("user is nil")
}
if user.Email == "" {
    return fmt.Errorf("email is required")
}
if !isValidEmail(user.Email) {
    return fmt.Errorf("invalid email")
}
// Основная логика без глубокой вложенности
```

### Premature Optimization

**❌ ПЛОХО:** оптимизируйте то, что не является узким местом

**✅ ХОРОШО:** пишите просто и понятно, оптимизируйте потом с бенчмарками

## Комментарии

### Когда писать

**✅ ПЛОХО (избыточно):**
```go
// Получаем пользователя по ID
user, err := db.GetUserByID(id)
// Если ошибка — возвращаем её
if err != nil {
    return err
}
```

**✅ ХОРОШО (объясняет "почему"):**
```go
// Используем XOR 0xFF для инверсии битов (требование протокола)
result[i] = (data[i] ^ 0xFF)
```

### Документация для экспортируемых функций

```go
// NewDatabase creates a new database connection.
// By default, it uses a NoopLogger. Use WithLogger option to provide a custom logger.
// Returns an error if connection cannot be established.
func NewDatabase(cfg *config.DatabaseConfig, opts ...Option) (Database, error) { ... }
```

## Рефакторинг

### Когда рефакторить

- Перед добавлением новой функциональности
- Когда видите дублирование (DRY)
- Когда функция становится слишком большой
- Когда видите антипаттерны

### Как рефакторить

1. Напишите тесты для кода
2. Рефакторите маленькими шагами
3. Запускайте тесты после каждого шага
4. Не смешивайте рефакторинг с добавлением функциональности
5. Коммитьте часто

## Критические правила

1. ВСЕГДА обрабатывайте ошибки
2. НИКОГДА не игнорируйте ошибки без причины
3. ВСЕГДА используйте понятные имена
4. НИКОГДА не используйте magic numbers
5. ВСЕГДА пишите короткие функции (< 50 строк)
6. НИКОГДА не создавайте God Functions
7. ВСЕГДА избегайте глубокой вложенности (> 3 уровней)
8. НИКОГДА не оптимизируйте преждевременно
9. ВСЕГДА пишите комментарии для "почему", не для "что"
10. НИКОГДА не нарушайте SOLID принципы

---

**Памятка:** Код пишется один раз, читается много раз. Делайте его понятным.
