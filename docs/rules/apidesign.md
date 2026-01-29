# REST/Web API Design для Nexflow

Этот модуль описывает правила проектирования REST/Web API.

## Общие принципы

### RESTful API

- HTTP методы по назначению (GET, POST, PUT, DELETE, PATCH)
- Ресурсо-ориентированный дизайн
- Правильные HTTP status codes
- JSON формат запросов/ответов
- Версионирование через URL путь

### WebSocket API

- Реальное время (чаты, логи, метрики)
- Reconnection логика на клиенте
- Heartbeat для keep-alive

## Базовые эндпоинты

### Health Check

**GET** `/health` — проверка работоспособности

```json
{
  "status": "ok",
  "version": "1.0.0",
  "timestamp": "2025-01-29T12:00:00Z"
}
```

**Status:** 200 OK, 503 Service Unavailable

### Metrics

**GET** `/metrics` — метрики Prometheus

**Status:** 200 OK

## Chat API

### Send Message

**POST** `/api/v1/chat`

**Request:**
```json
{
  "message": "привет, как дела?",
  "session_id": "optional-session-id",
  "user_id": "user-identifier",
  "channel": "telegram|discord|web",
  "metadata": {"key": "value"}
}
```

**Response:**
```json
{
  "session_id": "session-uuid",
  "message_id": "message-uuid",
  "response": "Привет! Всё хорошо.",
  "model_used": "claude-3-opus-20240229",
  "tokens_used": {"input": 15, "output": 20, "total": 35},
  "timestamp": "2025-01-29T12:00:00Z"
}
```

**Status:** 200, 400, 401, 429, 500

### Get Session

**GET** `/api/v1/sessions/{id}`

**Response:**
```json
{
  "id": "session-uuid",
  "user_id": "user-identifier",
  "channel": "telegram",
  "created_at": "2025-01-29T10:00:00Z",
  "updated_at": "2025-01-29T12:00:00Z",
  "message_count": 15,
  "status": "active"
}
```

**Status:** 200, 404

### List Sessions

**GET** `/api/v1/sessions?user_id=xxx&limit=50&offset=0`

**Response:**
```json
{
  "sessions": [...],
  "total": 2,
  "limit": 50,
  "offset": 0
}
```

**Status:** 200

## Skills API

### List Skills

**GET** `/api/v1/skills`

**Response:**
```json
{
  "skills": [
    {
      "name": "shell-run",
      "description": "Выполнение shell команд",
      "category": "system",
      "permissions": ["shell"],
      "enabled": true
    }
  ]
}
```

**Status:** 200

### Execute Skill

**POST** `/api/v1/skills/{name}`

**Request:**
```json
{
  "input": "параметры для навыка",
  "timeout": 30
}
```

**Response:**
```json
{
  "skill_name": "shell-run",
  "output": "результат выполнения",
  "error": null,
  "execution_time_ms": 123,
  "timestamp": "2025-01-29T12:00:00Z"
}
```

**Status:** 200, 400, 404, 408, 500

## WebSocket API

### Chat WebSocket

**WS** `/ws/chat/{session}` — чат в реальном времени

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/chat/session-uuid');

ws.send(JSON.stringify({
  type: 'message',
  data: {message: 'привет', user_id: 'user-123', channel: 'web'}
}));

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  // msg.type: 'message' | 'error' | 'heartbeat'
};
```

**Типы сообщений:** message, error, heartbeat

### Logs WebSocket

**WS** `/ws/logs` — поток логов в реальном времени

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/logs');

ws.send(JSON.stringify({level: 'info'})); // debug, info, warn, error

ws.onmessage = (event) => {
  const log = JSON.parse(event.data);
  // log.level, log.message, log.timestamp, log.fields
};
```

## Аутентификация и авторизация

### Gateway Token

**Header:** `Authorization: Bearer {gateway_token}`

**Query param:** `?token={gateway_token}`

### Channels Auth

Встроенная аутентификация провайдера (bot token, webhook secret)

## Rate Limiting

### Правила

- Базовый лимит: 100 запросов/минута на пользователя
- WebSocket: максимум 5 соединений на пользователя
- Skill execution: 10/минута на пользователя

### 429 Too Many Requests

```json
{
  "error": "rate_limit_exceeded",
  "message": "Too many requests.",
  "retry_after": 60
}
```

**Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1640902800
Retry-After: 60
```

## Error Handling

### Формат ошибок

```json
{
  "error": "error_code",
  "message": "Human readable error message",
  "details": {"field": "Additional details"},
  "request_id": "uuid-request-id"
}
```

### Коды ошибок

- `validation_error` — ошибка валидации
- `authentication_failed` — ошибка аутентификации
- `authorization_failed` — ошибка авторизации
- `not_found` — ресурс не найден
- `rate_limit_exceeded` — превышен rate limit
- `internal_error` — внутренняя ошибка

### Примеры ошибок

**400 Bad Request:**
```json
{
  "error": "validation_error",
  "message": "Invalid request payload",
  "details": {"field": "message", "reason": "required field is empty"},
  "request_id": "req-uuid"
}
```

**401 Unauthorized:**
```json
{
  "error": "authentication_failed",
  "message": "Invalid or missing authentication token",
  "request_id": "req-uuid"
}
```

**404 Not Found:**
```json
{
  "error": "not_found",
  "message": "Session not found: session-uuid",
  "request_id": "req-uuid"
}
```

**500 Internal Server Error:**
```json
{
  "error": "internal_error",
  "message": "An unexpected error occurred.",
  "request_id": "req-uuid"
}
```

## CORS

### Разрешённые origins

**Development:** `*`

**Production:** указывать явно в конфигурации

```yaml
server:
  cors:
    allowed_origins: ["https://example.com"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allowed_headers: ["Content-Type", "Authorization"]
    exposed_headers: ["X-Request-ID"]
    max_age: 86400
}
```

## Request ID

Каждый запрос имеет уникальный ID для трассировки.

**Header:** `X-Request-ID: uuid-request-id`

## Версионирование API

### URL-версионирование

Все endpoints: `/api/v1/`

Новая версия: `/api/v2/`

### Совместимость

- v1 поддерживается минимум 12 месяцев после релиза v2
- Deprecation warnings в заголовках

**Header:**
```
Deprecation: true
Sunset: Sat, 01 Jan 2026 00:00:00 GMT
Link: <https://api.example.com/api/v2/chat>; rel="successor-version"
```

## Документация API

### OpenAPI/Swagger

Файл спецификации: `docs/api/openapi.yaml` или `docs/api/openapi.json`

**Пример спецификации:**
```yaml
openapi: 3.0.0
info:
  title: Nexflow API
  version: 1.0.0
servers:
  - url: http://localhost:8080
paths:
  /api/v1/chat:
    post:
      summary: Send message for processing
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [message, user_id, channel]
              properties:
                message: {type: string}
                user_id: {type: string}
                channel: {type: string, enum: [telegram, discord, web]}
      responses:
        '200':
          description: Message processed successfully
```

## Тестирование API

### Unit тесты

Тестируйте handlers изолировано с mock зависимостями.

### Интеграционные тесты

Тестируйте API end-to-end с реальным сервером.

```go
func TestChatAPI_Integration(t *testing.T) {
    server := setupTestServer()
    defer server.Close()

    // Создание запроса
    reqBody := ChatRequest{Message: "test", UserID: "user-123", Channel: "web"}
    body, _ := json.Marshal(reqBody)

    // Отправка запроса
    resp, err := http.Post(server.URL+"/api/v1/chat", "application/json", bytes.NewReader(body))
    require.NoError(t, err)
    defer resp.Body.Close()

    // Проверка ответа
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var result ChatResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    require.NoError(t, err)
    assert.NotEmpty(t, result.SessionID)
}
```

## Критические правила

1. ВСЕГДА используйте правильные HTTP методы
2. НИКОГДА не игнорируйте коды статусов
3. ВСЕГДА валидируйте входные данные
4. НИКОГДА не раскрывайте секреты в ошибках
5. ВСЕГДА используйте request ID
6. НИКОГДА не делайте бесконечные запросы (timeout)
7. ВСЕГДА используйте rate limiting
8. НИКОГДА не игнорируйте CORS для web UI
9. ВСЕГДА документируйте API (OpenAPI/Swagger)
10. НИКОГДА не меняйте API без backward compatibility

---

**Памятка:** API — это контракт с клиентами. Изменения должны быть обратимыми.
