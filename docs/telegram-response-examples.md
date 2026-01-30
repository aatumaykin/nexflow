# Telegram Response Examples

Этот документ демонстрирует, как отправлять различные типы сообщений в Telegram через Nexflow.

## Основные типы ответов

### Текстовые сообщения

```go
import (
    "github.com/atumaikin/nexflow/internal/infrastructure/channels"
)

// Простое текстовое сообщение
response := &channels.Response{
    Type:    channels.ResponseTypeText,
    Content: "Привет, мир!",
}

// Текстовое сообщение с HTML форматированием
response := &channels.Response{
    Type:    channels.ResponseTypeText,
    Content: "<b>Жирный текст</b> и <i>курсив</i>",
    Metadata: map[string]interface{}{
        "parse_mode": "HTML",
    },
}

// Текстовое сообщение с MarkdownV2 форматированием
response := &channels.Response{
    Type:    channels.ResponseTypeText,
    Content: "*Жирный* и _курсив_",
    Metadata: map[string]interface{}{
        "parse_mode": "MarkdownV2",
    },
}
```

### Фото

```go
// Фото с подписью
response := &channels.Response{
    Type:    channels.ResponseTypePhoto,
    Content: "Посмотрите на это фото",
    Caption: "Прекрасный закат",
    Media: &channels.MediaContent{
        URL: "https://example.com/photo.jpg",
    },
}

// Фото по file_id (переиспользование существующего файла)
response := &channels.Response{
    Type:    channels.ResponseTypePhoto,
    Caption: "Красивый закат",
    Media: &channels.MediaContent{
        FileID: "AgACAgIAAxkBAAIC...",
    },
}

// Фото из байтов
photoData := []byte{...} // JPEG данные
response := &channels.Response{
    Type:    channels.ResponseTypePhoto,
    Content: "Загруженное фото",
    Caption: "Из камеры",
    Media: &channels.MediaContent{
        FileData: photoData,
    },
}
```

### Документы

```go
// Документ с подписью
response := &channels.Response{
    Type:    channels.ResponseTypeDocument,
    Content: "Вот документ",
    Caption: "Важный файл",
    Media: &channels.MediaContent{
        URL:      "https://example.com/report.pdf",
        FileName: "report.pdf",
    },
}

// Документ по file_id
response := &channels.Response{
    Type:    channels.ResponseTypeDocument,
    Caption: "Мой отчет",
    Media: &channels.MediaContent{
        FileID:   "BQACAgIAAxkBAAI...",
        FileName: "report.pdf",
    },
}
```

### Аудио

```go
// Аудио файл
response := &channels.Response{
    Type:    channels.ResponseTypeAudio,
    Content: "Послушайте это аудио",
    Caption: "Отличная песня",
    Media: &channels.MediaContent{
        FileID: "AwACAgIAAxkBAAI...",
    },
}
```

### Видео

```go
// Видео файл
response := &channels.Response{
    Type:    channels.ResponseTypeVideo,
    Content: "Посмотрите это видео",
    Caption: "Смешной клип",
    Media: &channels.MediaContent{
        URL: "https://example.com/video.mp4",
    },
}
```

### Стикеры

```go
// Отправка стикера
response := &channels.Response{
    Type: channels.ResponseTypeSticker,
    Media: &channels.MediaContent{
        FileID: "CAACAgIAAxkBAAI...",
    },
}
```

## Inline кнопки

```go
// Ответ с кнопками
response := &channels.Response{
    Type:    channels.ResponseTypeText,
    Content: "Выберите опцию:",
    Buttons: []channels.InlineButton{
        {
            Text: "Опция 1",
            Data: "opt1",
        },
        {
            Text: "Опция 2",
            Data: "opt2",
        },
    },
}
```

### Кнопка с URL

```go
// Кнопка-ссылка
response := &channels.Response{
    Type:    channels.ResponseTypeText,
    Content: "Нажмите для перехода на сайт",
    Buttons: []channels.InlineButton{
        {
            Text: "Открыть сайт",
            URL:  "https://example.com",
        },
    },
}
```

### Редактирование сообщения

```go
// Редактирование существующего сообщения
response := &channels.Response{
    Type:       channels.ResponseTypeText,
    Content:    "Обновленное сообщение",
    MessageID:  "123", // ID сообщения для редактирования
}
```

## Лимиты Telegram

- **Текст**: максимальная длина 4096 символов
- **Подпись (caption)**: максимальная длина 1024 символа
- **Фотографии**: до 10MB для ботов
- **Видео**: до 50MB для ботов
- **Документы**: до 50MB для ботов
- **Аудио**: до 50MB для ботов

## Rate Limiting

Nexflow автоматически применяет rate limiting для соблюдения ограничений Telegram API:
- Максимум 30 сообщений в секунду
- Длинные сообщения автоматически разбиваются на части с задержкой между ними

## Обработка ошибок

Connector обрабатывает распространенные ошибки Telegram API:
- Бот заблокирован пользователем
- Пользователь деактивирован
- Чат не найден
- Превышен rate limit
- Слишком длинное сообщение

Все ошибки логируются с контекстом для отладки.
