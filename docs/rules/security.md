# Правила безопасности для Nexflow

Критические правила безопасности, которые **всегда** должны соблюдаться.

## Основные принципы

1. Никаких секретов в коде
2. Валидация всех входных данных
3. Принцип наименьших привилегий
4. Sandbox для опасных операций
5. Маскирование секретов в логах
6. Whitelist вместо blacklist

## Секреты (НИКОГДА не хардкодить)

**Категории:** API keys, passwords, tokens, secrets, private keys, connection strings, JWT, authorization

**✅ ХОРОШО:**
```go
apiKey := os.Getenv("OPENAI_API_KEY")
// config.yml: api_key: "${OPENAI_API_KEY}"
```

**❌ ПЛОХО:**
```go
apiKey := "sk-ant-..."  // Хардкод!
```

## Переменные окружения

Конфигурация поддерживает `${VAR_NAME}` в YAML/JSON:

```yaml
llm:
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
```

```go
apiKey := os.Getenv("OPENAI_API_KEY")
if apiKey == "" {
    return fmt.Errorf("OPENAI_API_KEY is required")
}
```

## Маскирование в логах

Автоматическое маскирование в slog для ключей: token, key, password, secret, apikey, privatekey, connectionstring, jwt, authorization

```go
logger.Info("Processing",
    "user_id", userID,
    "api_key", secretKey,  // Будет: "***"
)
```

## SQL Injection Prevention

Используйте SQLC или prepared statements:

**✅ ХОРОШО:**
```go
user, err := queries.GetUserByName(ctx, GetUserByNameParams{Name: userInput})
// или
db.Query("SELECT * FROM users WHERE name = $1", userInput)
```

**❌ ПЛОХО:**
```go
query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", userInput)
```

## Whitelist для ресурсов

### Файлы

```go
allowedDirs := []string{"/tmp/nexflow-sandbox"}

func isPathAllowed(path string) bool {
    absPath, _ := filepath.Abs(path)
    for _, dir := range allowedDirs {
        if strings.HasPrefix(absPath, dir) {
            return true
        }
    }
    return false
}
```

### Хосты

```go
allowedHosts := []string{"api.github.com", "api.openai.com"}

func isURLAllowed(urlStr string) (bool, error) {
    u, _ := url.Parse(urlStr)
    for _, host := range allowedHosts {
        if u.Host == host {
            return true, nil
        }
    }
    return false, nil
}
```

### Shell команды

```go
allowedCommands := map[string]bool{"git": true, "docker": true}

func isCommandAllowed(cmd string) bool {
    parts := strings.Fields(cmd)
    return allowedCommands[parts[0]]
}
```

## Sandbox для опасных навыков

**Категории:** shell, filesystem, network, secrets

**В SKILL.md:**
```yaml
permissions: [shell, filesystem, network]
```

**Принципы:**
- Изолированный контейнер
- Ограничение по времени/памяти/IO
- Виртуальная сеть

## Валидация входных данных

```go
func validateInput(input SkillInput) error {
    if input.Command == "" {
        return fmt.Errorf("command is required")
    }
    if !isValidCommandFormat(input.Command) {
        return fmt.Errorf("invalid command format")
    }
    return nil
}
```

## LLM провайдеры

```go
// ✅ ХОРОШО
apiKey := config.LLM.Providers["anthropic"].APIKey

// ❌ ПЛОХО
apiKey := "sk-ant-..."
```

Валидация промптов:
```go
func validatePrompt(prompt string) error {
    if len(prompt) > 100000 {
        return fmt.Errorf("prompt too long")
    }
    // ... проверка на инъекции
    return nil
}
```

## Channels безопасность

**Telegram:**
```yaml
telegram:
  bot_token: "${TELEGRAM_BOT_TOKEN}"
  allowed_users: [123456789, 987654321]
  allowed_chats: [-1001234567890]
```

**Discord:**
```yaml
discord:
  bot_token: "${DISCORD_BOT_TOKEN}"
  allowed_roles: ["admin", "dev"]
```

**Web UI:**
```yaml
web:
  enabled: true
  gateway_token: "${WEB_GATEWAY_TOKEN}"
```

## Vault/Secret Manager (опционально)

```go
import "github.com/hashicorp/vault/api"

func getSecretFromVault(path string) (string, error) {
    client, _ := api.NewClient(api.DefaultConfig())
    secret, err := client.Logical().Read(path)
    if err != nil {
        return "", err
    }
    return secret.Data["value"].(string), nil
}
```

## Сканеры секретов

**Инструменты:** gitleaks, git-secrets, truffleHog

**Настройка gitleaks:**
```bash
brew install gitleaks
gitleaks detect --source . --verbose
```

**GitHub Actions:**
```yaml
- uses: gitleaks/gitleaks-action@v2
```

**git-secrets:**
```bash
brew install git-secrets
git secrets --register-aws
git secrets --add 'api_key'
git secrets --scan
```

## Git контроль файлов (КРИТИЧНО)

**ПЕРЕД ЛЮБЫМ git commit:**

1. **Проверка файлов на коммит:**
   ```bash
   git status --short
   ```

2. **Запрещены к коммиту:**
   - Бинарные файлы: `*.exe`, `*.dll`, `*.so`, `*.dylib`, `server`, `*.test`
   - Coverage файлы: `*.out`, `coverage.out`
   - Базы данных: `*.db`, `*.sqlite`, `*.sqlite3`
   - Логи: `*.log`, `logs/`, `*.log.*`
   - Временные файлы: `tmp/`, `temp/`, `*.tmp`, `*.swp`, `*.swo`, `*~`
   - Файлы с локальными путями: содержащие `/Users/`, `/home/`
   - OS файлы: `.DS_Store`, `Thumbs.db`
   - Environment: `.env`, `.env.local`, `.env.*.local`
   - IDE: `.vscode/`, `.idea/`
   - Локальное состояние: `.beads/export-state/`, `.beads/daemon.log`, `.beads/daemon.lock`

3. **Правила проверки:**
   - НЕ добавляйте подозрительные файлы в commit
   - Уведомляйте пользователя о запрещенных файлах
   - Предлагайте добавить в `.gitignore`

4. **Проверка на локальные пути:**
   ```bash
   git diff --cached | grep -E "(\+|)\s*/(Users|home)/" || echo "OK"
   ```

5. **Пример уведомления:**
   ```
   ⚠️ Файл server не должен быть закоммичен. Добавьте его в .gitignore или удалите перед коммитом.
   ```

## Логирование безопасности

**✅ Логируйте:**
```go
logger.Warn("Authentication failed", "user_id", userID, "channel", "telegram")
logger.Warn("Request blocked", "resource", resource, "type", resourceType)
```

**❌ НЕ логируйте:** API ключи, passwords, tokens, секреты

## Тестирование безопасности

```go
func TestValidateInput(t *testing.T) {
    tests := []struct {
        input   SkillInput
        wantErr bool
    }{
        {SkillInput{Command: "ls"}, false},
        {SkillInput{Command: "rm -rf /"}, true},
    }
    for _, tt := range tests {
        err := validateInput(tt.input)
        if (err != nil) != tt.wantErr {
            t.Errorf("validateInput() error = %v, wantErr %v", err, tt.wantErr)
        }
    }
}
```

## КРИТИЧЕСКИЕ ПРАВИЛА (Never Break)

1. НИКОГДА не коммитьте секреты
2. ВСЕГДА используйте ENV для секретов
3. ВСЕГДА маскируйте секреты в логах
4. ВСЕГДА валидируйте входные данные
5. ВСЕГДА используйте whitelist
6. ВСЕГДА изолируйте опасные навыки
7. НИКОГДА не позволяйте произвольный shell
8. НИКОГДА не позволяйте произвольные файлы
9. НИКОГДА не позволяйте произвольные запросы
10. ВСЕГДА запускайте сканеры секретов в CI/CD

---

**Памятка:** Безопасность — это не опция, это требование.
