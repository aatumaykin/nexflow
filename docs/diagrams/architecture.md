# Архитектурные диаграммы Nexflow

## Обзор архитектуры

В этом документе представлены архитектурные диаграммы проекта Nexflow.

## High-Level архитектура

```mermaid
graph TB
    subgraph "Channels Layer"
        TG[Telegram Connector]
        DC[Discord Connector]
        WB[Web Connector]
        OT[Other Connectors]
    end

    subgraph "Core Gateway"
        MR[Message Router]
        OR[Orchestrator]
        LR[LLM Router]
        SR[Skill Registry]
    end

    subgraph "Storage & Execution"
        DB[(Database)]
        FS[File System]
        SB[Sandbox]
    end

    subgraph "External Services"
        LLM1[Anthropic]
        LLM2[OpenAI]
        LLM3[Ollama]
        LLM4[Others]
    end

    TG --> MR
    DC --> MR
    WB --> MR
    OT --> MR

    MR --> OR
    OR --> LR
    OR --> SR

    LR --> LLM1
    LR --> LLM2
    LR --> LLM3
    LR --> LLM4

    SR --> FS
    SR --> SB
    OR --> DB
    MR --> DB
    SR --> DB

    style TG fill:#4CAF50
    style DC fill:#4CAF50
    style WB fill:#4CAF50
    style OT fill:#4CAF50
    style DB fill:#2196F3
    style FS fill:#FF9800
    style SB fill:#FF9800
    style LLM1 fill:#9C27B0
    style LLM2 fill:#9C27B0
    style LLM3 fill:#9C27B0
    style LLM4 fill:#9C27B0
```

## Детальная архитектура слоёв

```mermaid
graph TB
    subgraph "Application Layer"
        UC[Use Cases]
        DTO[DTOs]
        PRT[Ports]
    end

    subgraph "Domain Layer"
        ENT[Entities]
        REP[Repository Interfaces]
    end

    subgraph "Infrastructure Layer"
        CHN[Channels]
        LLM[LLM Providers]
        DB[Database]
        SKL[Skills]
    end

    subgraph "Shared Layer"
        CFG[Config]
        LOG[Logging]
        UTL[Utils]
    end

    PRT --> CHN
    PRT --> LLM
    PRT --> SKL
    UC --> REP
    UC --> ENT
    CHN --> CFG
    CHN --> LOG
    DB --> REP
    LLM --> CFG
    SKL --> CFG
    SKL --> FS[File System]
    SKL --> SB[Sandbox]

    UC --> CFG
    UC --> LOG
    CHN --> UTL
    DB --> LOG
    LLM --> LOG

    style ApplicationLayer fill:#E3F2FD
    style DomainLayer fill:#C8E6C9
    style InfrastructureLayer fill:#FFECB3
    style SharedLayer fill:#F5F5F5
```

## Flow Diagram: Обработка сообщения

```mermaid
sequenceDiagram
    participant User
    participant Channel
    participant Router
    participant Orchestrator
    participant LLMRouter
    participant LLM
    participant Skill
    participant DB

    User->>Channel: Отправляет сообщение
    Channel->>Router: Event(message)
    Router->>DB: Get/Create User
    DB-->>Router: User
    Router->>DB: Get/Create Session
    DB-->>Router: Session
    Router->>Router: Save Message
    Router->>Orchestrator: Process(message)
    Orchestrator->>DB: Get Conversation History
    DB-->>Orchestrator: Messages
    Orchestrator->>LLMRouter: Generate Completion
    LLMRouter->>LLM: Send Request
    LLM-->>LLMRouter: Response (with tool calls)
    LLMRouter-->>Orchestrator: Response
    alt Tool Call Required
        Orchestrator->>Skill: Execute Skill
        Skill-->>Orchestrator: Result
        Orchestrator->>LLMRouter: Continue with Result
        LLMRouter->>LLM: Send Request
        LLM-->>LLMRouter: Final Response
    end
    Orchestrator->>DB: Save Assistant Message
    Orchestrator->>Channel: SendMessage
    Channel->>User: Response
```

## Flow Diagram: Выполнение навыка

```mermaid
flowchart TD
    Start([Начало]) --> Validate{Skill Valid?}
    Validate -->|No| Error1[Return Error]
    Validate -->|Yes| CheckPerm{Requires Sandbox?}

    CheckPerm -->|Yes| SB[Create Sandbox]
    CheckPerm -->|No| Exec[Execute Skill]

    SB --> Exec
    Exec --> Timeout{Timeout?}

    Timeout -->|Yes| Error2[Return Timeout Error]
    Timeout -->|No| Success{Success?}

    Success -->|No| Error3[Return Error]
    Success -->|Yes| Log[Log Result]

    Log --> SaveDB[Save to DB]
    SaveDB --> Return([Return Result])

    style Start fill:#4CAF50
    style Return fill:#4CAF50
    style Error1 fill:#F44336
    style Error2 fill:#F44336
    style Error3 fill:#F44336
    style Validate fill:#2196F3
    style CheckPerm fill:#2196F3
    style Timeout fill:#2196F3
    style Success fill:#2196F3
```

## Class Diagram: Domain Entities

```mermaid
classDiagram
    class User {
        +string ID
        +string Channel
        +string ChannelID
        +time.Time CreatedAt
        +NewUser(channel, channelID) User
        +CanAccessSession(sessionID) bool
        +IsSameChannel(other) bool
    }

    class Session {
        +string ID
        +string UserID
        +time.Time CreatedAt
        +time.Time UpdatedAt
        +NewSession(userID) Session
        +UpdateTimestamp()
        +IsOwnedBy(userID) bool
    }

    class Message {
        +string ID
        +string SessionID
        +string Role
        +string Content
        +time.Time CreatedAt
        +NewUserMessage(sessionID, content) Message
        +NewAssistantMessage(sessionID, content) Message
        +IsFromUser() bool
        +IsFromAssistant() bool
    }

    class Task {
        +string ID
        +string SessionID
        +string Skill
        +string Input
        +string Output
        +string Status
        +time.Time CreatedAt
        +time.Time UpdatedAt
        +NewTask(sessionID, skill, input) Task
        +SetRunning()
        +SetCompleted(output)
        +SetFailed(err)
    }

    class Skill {
        +string ID
        +string Name
        +string Version
        +string Location
        +string Permissions
        +time.Time CreatedAt
        +NewSkill(name, version, location, permissions, metadata) Skill
        +RequiresPermission(permission) bool
        +RequiresSandbox() bool
        +GetTimeout() int
    }

    class Schedule {
        +string ID
        +string Skill
        +string CronExpression
        +string Input
        +bool Enabled
        +time.Time CreatedAt
        +NewSchedule(skill, cron, input) Schedule
        +Enable()
        +Disable()
    }

    class Log {
        +string ID
        +string Level
        +string Source
        +string Message
        +string Metadata
        +time.Time CreatedAt
        +NewLog(level, source, message, metadata) Log
        +IsDebug() bool
        +IsError() bool
    }

    User "1" --> "*" Session : owns
    Session "1" --> "*" Message : contains
    Session "1" --> "*" Task : tracks
    Skill "1" --> "*" Schedule : scheduled
    User "0..*" --> "*" Log : generates
```

## Component Diagram: Infrastructure Layer

```mermaid
graph TB
    subgraph "Channels"
        TGC[Telegram Connector]
        DSC[Discord Connector]
        WBC[Web Connector]
    end

    subgraph "LLM Providers"
        ANTP[Anthropic Provider]
        OPNP[OpenAI Provider]
        OLLP[Ollama Provider]
    end

    subgraph "Persistence"
        SQLDB[SQLite Database]
        PGDB[PostgreSQL Database]
        MAP[Mappers]
        MIG[Migrations]
    end

    subgraph "Skills"
        SKR[Skill Runtime]
        SKV[Skill Validator]
        SKX[Skill Executor]
    end

    TGC --> EV[Events Channel]
    DSC --> EV
    WBC --> EV

    ANTP --> REQ[HTTP Client]
    OPNP --> REQ
    OLLP --> REQ

    MAP --> SQLDB
    MAP --> PGDB
    MIG --> SQLDB
    MIG --> PGDB

    SKR --> SKV
    SKR --> SKX
    SKX --> SBX[Box Sandbox]
    SKX --> CMD[Command Executor]

    style Channels fill:#4CAF50
    style LLM fill:#9C27B0
    style Persistence fill:#2196F3
    style Skills fill:#FF9800
```

## Deployment Diagram

```mermaid
graph TB
    subgraph "Production Environment"
        subgraph "Kubernetes/Docker"
            POD1[Nexflow Pod 1]
            POD2[Nexflow Pod 2]
            POD3[Nexflow Pod N]
        end

        subgraph "Database"
            PGDB[(PostgreSQL)]
            REDIS[(Redis)]
        end

        subgraph "Storage"
            S3[Object Storage]
            NFS[NFS Volume]
        end

        subgraph "Load Balancer"
            LB[NGINX/Traefik]
        end

        subgraph "Monitoring"
            PROM[Prometheus]
            GRAF[Grafana]
        end
    end

    LB --> POD1
    LB --> POD2
    LB --> POD3

    POD1 --> PGDB
    POD1 --> REDIS
    POD1 --> S3

    POD2 --> PGDB
    POD2 --> REDIS
    POD2 --> S3

    POD3 --> PGDB
    POD3 --> REDIS
    POD3 --> S3

    POD1 --> NFS
    POD2 --> NFS
    POD3 --> NFS

    PROM --> POD1
    PROM --> POD2
    PROM --> POD3
    PROM --> PGDB

    GRAF --> PROM

    style ProductionEnvironment fill:#E8F5E9
    style POD1 fill:#4CAF50
    style POD2 fill:#4CAF50
    style POD3 fill:#4CAF50
```

## ER Diagram: Database Schema

```mermaid
erDiagram
    USERS ||--o{ SESSIONS : owns
    USERS ||--o{ LOGS : generates

    SESSIONS ||--o{ MESSAGES : contains
    SESSIONS ||--o{ TASKS : tracks

    SKILLS ||--o{ TASKS : executes
    SKILLS ||--o{ SCHEDULES : scheduled

    USERS {
        string id PK
        string channel
        string channel_id
        string created_at
    }

    SESSIONS {
        string id PK
        string user_id FK
        string created_at
        string updated_at
    }

    MESSAGES {
        string id PK
        string session_id FK
        string role
        string content
        string created_at
    }

    TASKS {
        string id PK
        string session_id FK
        string skill
        string input
        string output
        string status
        string error
        string created_at
        string updated_at
    }

    SKILLS {
        string id PK
        string name UK
        string version
        string location
        string permissions
        string metadata
        string created_at
    }

    SCHEDULES {
        string id PK
        string skill FK
        string cron_expression
        string input
        bool enabled
        string created_at
    }

    LOGS {
        string id PK
        string level
        string source
        string message
        string metadata
        string created_at
    }
```

## State Diagram: Task Lifecycle

```mermaid
stateDiagram-v2
    [*] --> Pending: Create
    Pending --> Running: Start
    Running --> Completed: Success
    Running --> Failed: Error
    Completed --> [*]
    Failed --> [*]
    Pending --> Failed: Cancel
```

## Network Diagram: External Connections

```mermaid
graph LR
    subgraph "Nexflow"
        NF[Core Gateway]
    end

    subgraph "Messaging Platforms"
        TG[Telegram API]
        DC[Discord API]
        SL[Slack API]
    end

    subgraph "LLM Providers"
        AN[Anthropic API]
        OA[OpenAI API]
        GO[Google API]
        OL[Ollama]
        OR[OpenRouter]
    end

    subgraph "Infrastructure"
        AWS[AWS S3]
        GCP[GCP Storage]
        LOC[Local FS]
    end

    NF <--> TG
    NF <--> DC
    NF <--> SL

    NF <--> AN
    NF <--> OA
    NF <--> GO
    NF <--> OL
    NF <--> OR

    NF --> AWS
    NF --> GCP
    NF --> LOC

    style Nexflow fill:#4CAF50
    style Messaging fill:#2196F3
    style LLM fill:#9C27B0
    style Infrastructure fill:#FF9800
```
