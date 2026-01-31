# Coverage Guide

## Обзор

Документ описывает работу с покрытием кода (code coverage) в проекте Nexflow.

## Цели покрытия

- **Минимальный порог:** 60%
- **Целевой порог:** 70%
- **Текущее покрытие:** ~44.5% (на 2026-01-31)

## CI/CD Integration

### Автоматическая проверка

CI pipeline в GitHub Actions автоматически:
1. Запускает тесты с покрытием
2. Генерирует HTML отчеты
3. Проверяет порог 60%
4. Генерирует отчеты по каждому пакету
5. Загружает результаты в Codecov

### Пороги

- **Проект:** Минимум 60% (необходимо для прохождения CI)
- **Patch:** Минимум 60% (для каждого PR)
- Если покрытие ниже порога, CI не проходит

## Локальная работа

### Быстрая проверка

```bash
# Использование скрипта
./scripts/coverage-check.sh

# Использование Makefile
make coverage-check
```

### Детальная проверка

```bash
# Сгенерировать HTML отчет
make coverage-html
open coverage.html

# Покрытие по функциям
make coverage-func

# Покрытие по пакетам
make coverage-packages
```

### Запуск тестов

```bash
# Все тесты
make test

# С покрытием
make test-cover

# С race detection
make test-race

# Все проверки (как в CI)
make ci
```

## Отчеты

### HTML отчет

Генерируется локально и в CI:

```bash
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out -o coverage.html
```

### Покрытие по функциям

```bash
go tool cover -func=coverage.out
```

Пример вывода:

```
github.com/atumaikin/nexflow/internal/application/dto/session.go:18:	NewSessionDTO	100.0%
github.com/atumaikin/nexflow/internal/application/dto/session.go:28:	ToEntity	100.0%
...
total:	(statements)	44.5%
```

### Покрытие по пакетам

```bash
go tool cover -func=coverage.out | grep "^github.com" | awk '{
    pkg = $1
    gsub(/\/github.com\/atumaikin\/nexflow\//, "", pkg)
    cov = $3
    print pkg ": " cov
}' | sort
```

## Codecov

### Настройка

Конфигурация в `.codecov.yml`:

```yaml
coverage:
  precision: 2
  round: down
  range: "60...100"

status:
  project:
    default:
      target: auto
      threshold: 60%
```

### CI Integration

Codecov автоматически:
- Отслеживает тренд покрытия
- Показывает изменения в PR
- Генерирует отчеты по каждому коммиту
- Комментирует PR с изменениями покрытия

## Badges

Coverage badge добавлен в README:

```markdown
![codecov](https://codecov.io/gh/aatumaykin/nexflow/branch/main/graph/badge.svg)
```

## Улучшение покрытия

### Шаги

1. **Запустите анализ:**
   ```bash
   ./scripts/coverage-check.sh
   ```

2. **Найдите наименее покрытый код:**
   ```bash
   make coverage-func | sort -k2 -n
   ```

3. **Откройте HTML отчет:**
   ```bash
   make coverage-html
   open coverage.html
   ```

4. **Напишите тесты** для непокрытого кода

5. **Проверьте результат:**
   ```bash
   make coverage-check
   ```

### Лучшие практики

- Пишите тесты одновременно с кодом (TDD)
- Используйте табличные тесты для сложной логики
- Мокайте внешние зависимости (database, API)
- Проверяйте все ветки (branch coverage)
- Покрывайте ошибки и edge cases

### Приоритеты

Покрытие в порядке приоритета:

1. **Domain Layer** - бизнес-логика (приоритет: высокий)
2. **Application Layer** - use cases (приоритет: высокий)
3. **Infrastructure Layer** - адаптеры (приоритет: средний)
4. **Presentation Layer** - handlers/commands (приоритет: низкий)

## Исключения

Следующие директории исключены из покрытия:

- `**/*_test.go` - тестовые файлы
- `**/mock*.go` - моки
- `**/cmd/genmapper/**` - утилиты генерации
- `**/cmd/validate-config/**` - утилиты валидации
- `**/examples/**` - примеры

## Траблшутинг

### Coverage.out пуст или 0%

Убедитесь, что тесты запускаются:

```bash
go test ./...
go tool cover -func=coverage.out
```

### Coverage ниже порога

1. Проверьте, какие пакеты не покрыты
2. Напишите тесты для критического кода
3. Используйте `make coverage-packages` для анализа

### Codecov не обновляется

Проверьте CI status и убедитесь, что `codecov-action` успешно завершается.

## Дополнительно

- [Go Test Coverage](https://go.dev/blog/cover)
- [Codecov Documentation](https://docs.codecov.com/)
- [Testing Guide](testing-guide.md)
