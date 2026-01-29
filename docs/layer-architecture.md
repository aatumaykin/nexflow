# Слоистая архитектура Nexflow

Документ описывает детальную слоистую архитектуру проекта Nexflow, включая каждый слой, его ответственность и взаимодействия.

## Введение

Nexflow реализует Clean Layered Architecture с четким разделением на 4 слоя:
1. **Presentation Layer** - пользовательский интерфейс и API
2. **Application Layer** - бизнес-логика и use cases
3. **Domain Layer** - сущности и бизнес-правила
4. **Infrastructure Layer** - технические реализации

## Presentation Layer

### Расположение
`cmd/` - entry points

### Ответственность
- Обработка HTTP/WebSocket запросов
- Маршрутизация API endpoints
- Интеграции с Telegram/Discord bots
- Преобразование HTTP request/response

### Компоненты

#### HTTP Server
```go
// cmd/server/main.go
func main() {
    config := config.Load("config.yml")
    db := database.New(config.Database)
    
    // Инициализация repositories
    userRepo := sqlite.NewUserRepository(db)
    sessionRepo := sqlite.NewSessionRepository(db)
    // ...
    
    // Инициализация use cases
    userUC := usecase.NewUserUseCase(userRepo, logger)
    chatUC := usecase.NewChatUseCase(userRepo, sessionRepo, ...)
    
    // Setup HTTP handlers
    router := setupRouter(userUC, chatUC, ...)
    
    // Start server
    server := &http.Server{Addr: config.Server.Address}
    server.ListenAndServe()
}
```

#### Handlers
```go
// Presentation layer handles HTTP requests
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // Delegate to use case
    resp, err := h.userUC.CreateUser(ctx, req)
    
    // Convert to HTTP response
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(resp)
}
```

### Правила
- Никакой бизнес-логики
- Минимум валидации (базовая проверка формата)
- Преобразование только HTTP <-> DTO
- Никаких прямых вызовов к БД или внешних сервисов

## Application Layer

### Расположение
`internal/application/` - use cases и orchestration

### Ответственность
- Оркестрация бизнес-процессов
- Координация между domain entities
- Преобразование Domain <-> DTO
- Транзакционное управление

### Use Cases

#### UserUseCase
```go
type UserUseCase struct {
    userRepo repository.UserRepository
    logger   logging.Logger
}

func (uc *UserUseCase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    // 1. Check if user exists
    existing, err := uc.userRepo.FindByChannel(ctx, req.Channel, req.ChannelID)
    if err == nil && existing != nil {
        return &dto.UserResponse{
            Success: false,
            Error:   "user already exists",
        }, nil
    }
    
    // 2. Create new user
    user := entity.NewUser(req.Channel, req.ChannelID)
    if err := uc.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // 3. Convert to DTO
    return &dto.UserResponse{
        Success: true,
        User:    dto.UserDTOFromEntity(user),
    }, nil
}
```

#### ChatUseCase
```go
type ChatUseCase struct {
    userRepo    repository.UserRepository
    sessionRepo repository.SessionRepository
    messageRepo repository.MessageRepository
    taskRepo    repository.TaskRepository
    llmProvider ports.LLMProvider
    logger      logging.Logger
}

func (uc *ChatUseCase) SendMessage(ctx context.Context, req dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
    // 1. Find or create user
    user, err := uc.userRepo.FindByChannel(ctx, "web", req.UserID)
    if err != nil {
        user = entity.NewUser("web", req.UserID)
        uc.userRepo.Create(ctx, user)
    }
    
    // 2. Create session
    session := entity.NewSession(user.ID)
    uc.sessionRepo.Create(ctx, session)
    
    // 3. Save user message
    userMsg := entity.NewUserMessage(session.ID, req.Message.Content)
    uc.messageRepo.Create(ctx, userMsg)
    
    // 4. Get conversation history
    messages, _ := uc.messageRepo.FindBySessionID(ctx, session.ID)
    
    // 5. Call LLM
    llmResp, _ := uc.llmProvider.Generate(ctx, ports.CompletionRequest{
        Messages:  convertToLLMMessages(messages),
        Model:     req.Options.Model,
    })
    
    // 6. Save assistant message
    assistantMsg := entity.NewAssistantMessage(session.ID, llmResp.Message.Content)
    uc.messageRepo.Create(ctx, assistantMsg)
    
    return &dto.SendMessageResponse{
        Success:  true,
        Message:  dto.MessageDTOFromEntity(assistantMsg),
        Messages: convertToDTOs(messages),
    }, nil
}
```

### DTOs (Data Transfer Objects)
```go
type UserDTO struct {
    ID        string `json:"id"`
    Channel   string `json:"channel"`
    ChannelID string `json:"channel_id"`
    CreatedAt string `json:"created_at"`
}

func UserDTOFromEntity(user *entity.User) *UserDTO {
    return &UserDTO{
        ID:        user.ID,
        Channel:   user.Channel,
        ChannelID: user.ChannelID,
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
    }
}
```

### Правила
- Бизнес-логика и оркестрация здесь
- Зависимости только от интерфейсов (repository, ports)
- Использование domain entities и DTOs
- Никаких прямых вызовов к внешним сервисам

## Domain Layer

### Расположение
`internal/domain/` - бизнес-сущности

### Ответственность
- Определение бизнес-сущностей
- Бизнес-правила и валидация
- Repository interfaces (для dependency inversion)

### Entities

#### User Entity
```go
type User struct {
    ID        string    `json:"id"`
    Channel   string    `json:"channel"`
    ChannelID string    `json:"channel_id"`
    CreatedAt time.Time `json:"created_at"`
}

func NewUser(channel, channelID string) *User {
    return &User{
        ID:        utils.GenerateID(),
        Channel:   channel,
        ChannelID: channelID,
        CreatedAt: time.Now(),
    }
}

func (u *User) CanAccessSession(sessionID string) bool {
    // Business rule implementation
    return true
}
```

#### Session Entity
```go
type Session struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func NewSession(userID string) *Session {
    now := time.Now()
    return &Session{
        ID:        utils.GenerateID(),
        UserID:    userID,
        CreatedAt: now,
        UpdatedAt: now,
    }
}

func (s *Session) UpdateTimestamp() {
    s.UpdatedAt = time.Now()
}

func (s *Session) IsOwnedBy(userID string) bool {
    return s.UserID == userID
}
```

### Repository Interfaces

```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindByChannel(ctx context.Context, channel, channelID string) (*User, error)
    List(ctx context.Context) ([]*User, error)
    Delete(ctx context.Context, id string) error
}

type SessionRepository interface {
    Create(ctx context.Context, session *Session) error
    FindByID(ctx context.Context, id string) (*Session, error)
    FindByUserID(ctx context.Context, userID string) ([]*Session, error)
    Update(ctx context.Context, session *Session) error
    Delete(ctx context.Context, id string) error
}
```

### Правила
- Никаких зависимостей от других слоев
- Чистая бизнес-логика
- Entity-центричный подход
- Никаких framework-специфичных деталей

## Infrastructure Layer

### Расположение
`internal/infrastructure/` - технические реализации

### Ответственность
- Реализация repository interfaces
- Интеграции с БД, внешними API
- Технические детали (HTTP clients, database drivers)

### Database Implementations

```go
package sqlite

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
    dbUser := mappers.UserToDB(user)
    
    _, err := r.db.ExecContext(ctx,
        `INSERT INTO users (id, channel, channel_user_id, created_at) VALUES (?, ?, ?, ?)`,
        dbUser.ID, dbUser.Channel, dbUser.ChannelUserID, dbUser.CreatedAt,
    )
    
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}
```

### Ports (External Dependencies Interfaces)

```go
package ports

type LLMProvider interface {
    Generate(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
    GenerateWithTools(ctx context.Context, req CompletionRequest, tools []ToolDefinition) (*CompletionResponse, error)
    Stream(ctx context.Context, req CompletionRequest) (<-chan string, error)
    EstimateCost(req CompletionRequest) (float64, error)
}

type Connector interface {
    Listen(ctx context.Context, messages <-chan string) error
    SendMessage(ctx context.Context, chatID, message string) error
}
```

### Mappers

```go
package mappers

// UserToDB converts domain entity to database model
func UserToDB(user *entity.User) *dbmodel.User {
    return &dbmodel.User{
        ID:            user.ID,
        Channel:       user.Channel,
        ChannelUserID: user.ChannelID,
        CreatedAt:     user.CreatedAt.Format(time.RFC3339),
    }
}

// UserToDomain converts database model to domain entity
func UserToDomain(dbUser *dbmodel.User) *entity.User {
    createdAt, _ := time.Parse(time.RFC3339, dbUser.CreatedAt)
    return &entity.User{
        ID:        dbUser.ID,
        Channel:   dbUser.Channel,
        ChannelID: dbUser.ChannelUserID,
        CreatedAt: createdAt,
    }
}
```

### Правила
- Реализация интерфейсов из Domain/Application
- Никакой бизнес-логики (кроме технической валидации)
- Адаптация к внешним системам

## Взаимодействие между слоями

### Request Flow (Creation)

```
HTTP Request
    ↓
[Presentation] Handler parses request
    ↓
[Presentation] Creates DTO from request
    ↓
[Application] Use case validates and orchestrates
    ↓
[Application] Creates domain entity from DTO
    ↓
[Domain] Entity applies business rules
    ↓
[Infrastructure] Repository saves to DB
    ↓
[Domain] Entity returns result
    ↓
[Application] Use case creates response DTO
    ↓
[Presentation] Handler returns HTTP response
```

### Request Flow (Query)

```
HTTP Request with ID
    ↓
[Presentation] Handler extracts ID
    ↓
[Application] Use case calls repository
    ↓
[Infrastructure] Repository queries DB
    ↓
[Domain] Entity returned
    ↓
[Application] Use case converts to DTO
    ↓
[Presentation] Handler returns HTTP response
```

## Dependency Rules

### Allowed Dependencies

| Layer | Может зависеть от |
|--------|-------------------|
| Presentation | Application, shared (utils, config, logging) |
| Application | Domain, shared (utils, config, logging) |
| Domain | shared (utils, config, logging) |
| Infrastructure | Domain, shared (utils, config, logging) |

### Запрещенные Dependencies

- **Presentation** → Domain (только через Application)
- **Application** → Infrastructure (только через interfaces)
- **Domain** → Infrastructure (только через interfaces)
- **Infrastructure** → Application (только через callbacks/hooks)

## Тестирование

### Unit Tests
```go
// Domain entities tests
func TestUser_NewUser(t *testing.T) {
    user := entity.NewUser("telegram", "user123")
    
    assert.NotEmpty(t, user.ID)
    assert.Equal(t, "telegram", user.Channel)
    assert.Equal(t, "user123", user.ChannelID)
}

// Use case tests with mocks
func TestUserUseCase_CreateUser_Success(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockLogger := new(MockLogger)
    uc := usecase.NewUserUseCase(mockRepo, mockLogger)
    
    req := dto.CreateUserRequest{Channel: "telegram", ChannelID: "user123"}
    mockRepo.On("FindByChannel", ctx, "telegram", "user123").Return(nil, errors.New("not found"))
    mockRepo.On("Create", ctx, mock.Anything).Return(nil)
    
    resp, err := uc.CreateUser(ctx, req)
    
    require.NoError(t, err)
    assert.True(t, resp.Success)
    mockRepo.AssertExpectations(t)
}
```

### Integration Tests
```go
func TestUserRepository_Integration(t *testing.T) {
    db := setupTestDB(t) // In-memory SQLite
    repo := sqlite.NewUserRepository(db)
    
    user := entity.NewUser("telegram", "user123")
    err := repo.Create(ctx, user)
    require.NoError(t, err)
    
    found, err := repo.FindByID(ctx, user.ID)
    require.NoError(t, err)
    assert.Equal(t, user.ID, found.ID)
}
```

## Заключение

Слоистая архитектура обеспечивает:

1. **Разделение ответственности** - каждый слой имеет четкую роль
2. **Тестируемость** - легко unit/интеграционные тесты
3. **Масштабируемость** - независимая замена компонентов
4. **Поддерживаемость** - простое добавление функциональности
5. **Безопасность** - четкие границы и зависимости

---

**Последнее обновление:** Январь 2026
