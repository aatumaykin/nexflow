# Диаграммы потока данных Nexflow

## Обзор потока данных

В этом документе представлены диаграммы потока данных в проекте Nexflow.

## User Message Flow

```mermaid
sequenceDiagram
    participant User
    participant Channel
    participant MessageRouter
    participant Orchestrator
    participant UserRepo
    participant SessionRepo
    participant MessageRepo
    participant LLMRouter
    participant LLMProvider
    participant SkillRuntime

    User->>Channel: Отправляет сообщение
    Channel->>MessageRouter: Event(userID, message)
    MessageRouter->>UserRepo: GetOrCreateUser(channel, channelID)
    UserRepo-->>MessageRouter: User
    MessageRouter->>SessionRepo: GetOrCreateSession(userID)
    SessionRepo-->>MessageRouter: Session
    MessageRouter->>MessageRepo: Create Message(user, content)
    MessageRepo-->>MessageRouter: Message
    MessageRouter->>MessageRepo: Save Message
    MessageRouter->>Orchestrator: Process(session, message)
    Orchestrator->>MessageRepo: GetMessagesBySessionID(sessionID)
    MessageRepo-->>Orchestrator: Messages (history)
    Orchestrator->>LLMRouter: GenerateCompletion(messages)
    LLMRouter->>LLMProvider: Send request
    LLMProvider-->>LLMRouter: Response
    alt Tool Call in Response
        LLMRouter-->>Orchestrator: ToolCall(skill, input)
        Orchestrator->>SkillRuntime: Execute(skill, input)
        SkillRuntime-->>Orchestrator: Result
        Orchestrator->>LLMRouter: GenerateCompletion(messages + result)
        LLMRouter->>LLMProvider: Send request
        LLMProvider-->>LLMRouter: Final Response
    end
    LLMRouter-->>Orchestrator: Completion(content)
    Orchestrator->>MessageRepo: Create Message(assistant, content)
    MessageRepo-->>Orchestrator: Message
    Orchestrator->>SessionRepo: Update Session timestamp
    Orchestrator->>Channel: SendMessage(userID, content)
    Channel->>User: Ответ от AI
```

## Skill Execution Flow

```mermaid
flowchart LR
    Start([Event: Skill Request]) --> Validate[Validate Skill]
    Validate -->|Valid| CheckPerm{Requires Permissions?}

    CheckPerm -->|Yes| GetPerm[Get Permissions]
    CheckPerm -->|No| ParseCfg[Parse Config]

    GetPerm --> Sandbox{Dangerous?}
    Sandbox -->|Yes| CreateBox[Create Sandbox Box]
    Sandbox -->|No| Prepare[Prepare Execution]

    CreateBox --> Prepare
    Prepare --> ParseCfg

    ParseCfg --> GetTimeout{Has Timeout?}
    GetTimeout -->|Yes| SetTimeout[Set Execution Timeout]
    GetTimeout -->|No| Execute[Execute Skill]

    SetTimeout --> Execute
    Execute --> Monitor[Monitor Execution]
    Monitor --> Success{Success?}

    Success -->|No| LogErr[Log Error]
    Success -->|Yes| LogRes[Log Result]

    LogErr --> SaveDB[Save to DB]
    LogRes --> SaveDB

    SaveDB --> Notify[Notify Orchestrator]
    Notify --> Cleanup[Cleanup Sandbox]
    Cleanup --> Return([Return Result])

    style Start fill:#4CAF50
    style Return fill:#4CAF50
    style LogErr fill:#F44336
    style Validate fill:#2196F3
    style CheckPerm fill:#2196F3
    style Sandbox fill:#2196F3
    style GetTimeout fill:#2196F3
    style Success fill:#2196F3
```

## Data Flow: Session Management

```mermaid
stateDiagram-v2
    [*] --> Idle

    Idle --> UserMessage: User sends message
    UserMessage --> Processing: Create user message
    Processing --> LLMRequest: Send to LLM
    LLMRequest --> LLMResponse: Receive response
    LLMResponse --> AssistantMessage: Create assistant message
    AssistantMessage --> Idle: Wait for next message

    Idle --> Timeout: Inactivity timeout
    Timeout --> Idle: Keep alive
```

## Data Flow: User Creation

```mermaid
flowchart TD
    Start([User interaction begins]) --> GetID{User exists?}

    GetID -->|No| CreateID[Generate UUID]
    GetID -->|Yes| RetrieveUser[Get User from DB]

    CreateID --> CreateUser[Create User entity]
    CreateUser --> SaveUser[Save to DB]
    SaveUser --> Success

    RetrieveUser --> UpdateTimestamp[Update Last Seen]
    UpdateTimestamp --> ReturnUser([Return User])

    CreateUser --> Success
    SaveUser --> Success

    style Start fill:#4CAF50
    style Success fill:#4CAF50
    style ReturnUser fill:#4CAF50
    style GetID fill:#2196F3
```

## Data Flow: LLM Routing

```mermaid
flowchart TD
    Start([LLM Request]) --> CheckPolicy{Policy Check}

    CheckPolicy -->|Anthropic| Anthropic[Send to Anthropic]
    CheckPolicy -->|OpenAI| OpenAI[Send to OpenAI]
    CheckPolicy -->|Ollama| Ollama[Send to Ollama]
    CheckPolicy -->|Custom| Custom[Send to Custom Provider]

    Anthropic --> CheckCost{Estimate Cost?}
    OpenAI --> CheckCost
    Ollama --> CheckCost
    Custom --> CheckCost

    CheckCost -->|Yes| Estimate[Calculate Tokens Cost]
    CheckCost -->|No| Execute

    Estimate --> Budget{Within Budget?}

    Budget -->|No| Reject[Reject Request]
    Budget -->|Yes| Execute[Execute Request]

    Execute --> Retry{Success?}

    Retry -->|No| CheckRetry{Max Retries?}
    CheckRetry -->|No| Reject
    CheckRetry -->|Yes| Wait[Wait before retry]
    Wait --> Execute

    Retry -->|Yes| Return([Return Response])

    Reject --> Fail([Return Error])

    style Start fill:#4CAF50
    style Return fill:#4CAF50
    style Fail fill:#F44336
    style CheckPolicy fill:#2196F3
    style CheckCost fill:#2196F3
    style Budget fill:#2196F3
    style Retry fill:#2196F3
    style CheckRetry fill:#2196F3
```

## Data Flow: Logging

```mermaid
flowchart LR
    App[Application Layer] --> |Write Log| Logger[Structured Logger]
    Domain[Domain Layer] --> |Write Log| Logger
    Infra[Infrastructure Layer] --> |Write Log| Logger

    Logger --> Mask{Secret Fields?}
    Mask --> |Yes| MaskSecrets[Mask token, key, password]
    Mask --> |No| Format[Format Message]

    MaskSecrets --> Format
    Format --> JSON{JSON Format?}
    Format --> |Yes| JSONOutput[JSON Output]
    Format --> |No| TextOutput[Text Output]

    JSONOutput --> Stdout[Write to Stdout]
    TextOutput --> Stdout

    Stdout --> File[Log File]
    Stdout --> Remote[Remote Logging]

    style App fill:#E3F2FD
    style Domain fill:#C8E6C9
    style Infra fill:#FFECB3
    style Logger fill:#4CAF50
```

## Data Flow: Configuration Loading

```mermaid
flowchart TD
    Start([Application Start]) --> ReadFile[Read Config File]
    ReadFile --> Parse{Parse Format?}

    Parse --> |YAML| ParseYAML[Parse YAML]
    Parse --> |JSON| ParseJSON[Parse JSON]

    ParseYAML --> EnvVars[Replace ENV Vars]
    ParseJSON --> EnvVars

    EnvVars --> Validate[Validate Config]
    Validate --> Valid{Valid?}

    Valid --> |No| Error[Return Error]
    Valid --> |Yes| Apply[Apply Defaults]
    Apply --> Build[Build Config Struct]
    Build --> Success([Return Config])

    Error --> Fail([Exit with Error])

    style Start fill:#4CAF50
    style Success fill:#4CAF50
    style Fail fill:#F44336
    style Parse fill:#2196F3
    style Valid fill:#2196F3
```

## Data Flow: Database Operations

```mermaid
sequenceDiagram
    participant UseCase
    participant Repository
    participant Mapper
    participant QueryBuilder
    participant Database
    participant SQLC

    UseCase->>Repository: CreateEntity(entity)
    Repository->>Mapper: EntityToDB(entity)
    Mapper-->>Repository: DBModel
    Repository->>QueryBuilder: BuildQuery(params)
    QueryBuilder-->>Repository: Query
    Repository->>SQLC: Execute(query, args)
    SQLC->>Database: Send SQL
    Database-->>SQLC: Result
    SQLC-->>Repository: DBModel (result)
    Repository->>Mapper: DBModelToEntity(dbModel)
    Mapper-->>Repository: Entity
    Repository-->>UseCase: Entity

    Note over Repository,SQLC: Using SQLC for type-safe<br/>database operations
```

## Data Flow: Error Handling

```mermaid
flowchart TD
    Error([Error Occurs]) --> Wrap[Wrap with Context]
    Wrap --> Log[Log Error]
    Log --> Type{Error Type?}

    Type --> |Validation| Validation[Return 400]
    Type --> |Not Found| NotFound[Return 404]
    Type --> |Permission| Permission[Return 403]
    Type --> |Conflict| Conflict[Return 409]
    Type --> |Internal| Internal[Return 500]

    Validation --> Notify[Notify User]
    NotFound --> Notify
    Permission --> Notify
    Conflict --> Notify
    Internal --> LogError[Log Full Error]
    LogError --> Notify

    Notify --> Response([Return HTTP Response])

    style Error fill:#F44336
    style Response fill:#4CAF50
    style Type fill:#2196F3
    style LogError fill:#FF9800
```

## Data Flow: Skill Discovery

```mermaid
flowchart TD
    Start([System Start]) --> ScanDir[Scan Skills Directory]
    ScanDir --> FindFiles[Find SKILL.md files]
    FindFiles --> Parse[Parse SKILL.md]

    Parse --> ValidateMD{Valid Markdown?}
    ValidateMD --> |No| LogWarn[Log Warning]
    ValidateMD --> |Yes| Extract[Extract Metadata]

    Extract --> CheckName{Name exists?}
    CheckName --> |No| Skip[Skip Skill]
    CheckName --> |Yes| CheckPerm{Permissions defined?}

    CheckPerm --> |No| Default[Use default permissions]
    CheckPerm --> |Yes| ParsePerm[Parse permissions]

    Default --> Register
    ParsePerm --> Register[Register Skill]

    Register --> Next{More Skills?}
    Next --> |Yes| Parse
    Next --> |No| Build[Build Registry]
    Build --> Done([Registry Ready])

    style Start fill:#4CAF50
    style Done fill:#4CAF50
    style LogWarn fill:#FF9800
    style CheckName fill:#2196F3
    style CheckPerm fill:#2196F3
    style Next fill:#2196F3
```

## Data Flow: Scheduled Task Execution

```mermaid
flowchart TD
    Start([Cron Trigger]) --> GetSchedule[Get Schedule]
    GetSchedule --> Enabled{Schedule Enabled?}

    Enabled --> |No| Skip[Skip Execution]
    Enabled --> |Yes| ParseInput[Parse Input JSON]

    ParseInput --> ValidateInput{Input Valid?}
    ValidateInput --> |No| LogErr[Log Error]
    ValidateInput --> |Yes| CreateTask[Create Task]

    LogErr --> Skip
    CreateTask --> SaveTask[Save Task to DB]

    SaveTask --> SetRunning[Set Task to Running]
    SetRunning --> ExecuteSkill[Execute Skill]

    ExecuteSkill --> Result{Execution Result?}

    Result --> |Success| SetCompleted[Set Task to Completed]
    Result --> |Failure| SetFailed[Set Task to Failed]

    SetCompleted --> SaveResult[Save Output to DB]
    SetFailed --> SaveResult

    SaveResult --> LogExecution[Log Execution]
    LogExecution --> Wait([Wait for Next Trigger])

    style Start fill:#4CAF50
    style Wait fill:#4CAF50
    style LogErr fill:#F44336
    style Enabled fill:#2196F3
    style ValidateInput fill:#2196F3
    style Result fill:#2196F3
```

## Data Flow: Webhook Processing

```mermaid
flowchart LR
    Start([Webhook Request]) --> ValidateSig{Valid Signature?}
    ValidateSig --> |No| Reject[Return 401]
    ValidateSig --> |Yes| ParsePayload[Parse JSON Payload]

    ParsePayload --> Extract[Extract Event Type]
    Extract --> Route{Event Type?}

    Route --> |Message| ProcessMsg[Process Message]
    Route --> |User| ProcessUser[Process User Update]
    Route --> |System| ProcessSys[Process System Event]

    ProcessMsg --> GetContext[Get Session Context]
    GetContext --> Process[Process with Orchestrator]
    ProcessUser --> UpdateUser[Update User Data]
    UpdateUser --> SaveUser[Save to DB]
    ProcessSys --> LogEvent[Log System Event]

    Process --> SaveMsg[Save Message to DB]
    SaveMsg --> Response([Return 200 OK])
    SaveUser --> Response
    LogEvent --> Response

    Reject --> Fail([Return 401 Unauthorized])

    style Start fill:#4CAF50
    style Response fill:#4CAF50
    style Reject fill:#F44336
    style Fail fill:#F44336
    style ValidateSig fill:#2196F3
    style Route fill:#2196F3
```

## Data Flow: Session State Management

```mermaid
stateDiagram-v2
    [*] --> Inactive

    Inactive --> UserMessage: User sends message
    UserMessage --> WaitingLLM: Waiting for LLM

    WaitingLLM --> LLMProcessing: LLM processing
    LLMProcessing --> ToolCall: Tool call requested
    LLMProcessing --> AssistantMessage: Assistant response

    ToolCall --> ToolExecution: Executing skill
    ToolExecution --> WaitingLLM: Back to LLM

    AssistantMessage --> UserMessage: Waiting for user
    AssistantMessage --> Inactive: Session timeout

    UserMessage --> Inactive: User inactive
```
