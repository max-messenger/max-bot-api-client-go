# `2` Прослушивание обновлений и реакция на них

После запуска бота Max начнёт отправлять вам обновления.

> Подробности обо всех обновлениях смотрите в [официальной документации](https://dev.max.ru/).

Max Bot API позволяет прослушивать эти обновления, например:

```go
func Handle(_ context.Context, update model.Update) {
    switch update.UpdateType { // Определение типа пришедшего обновления
	case model.UpdateBotStarted:   // Обработчик начала диалога с ботом
		/* ... */
	case model.UpdateMessageCreated: // Обработчик новых сообщений
		/* ... */
	case model.UpdateBotAdded: // Обработчик добавления пользователя в беседу
		/* ... */
	}
}
```

Вы можете посмотреть [model.types](../model.types), чтобы увидеть все доступные типы обновлений. `UpdateType` содержит актуальный список типов.

## Получение сообщений

Вы можете [подписаться на обновление](https://dev.max.ru/docs-api/methods/POST/subscriptions) `message_created`:

```go
func Handle(_ context.Context, update model.Update) {
    switch update.UpdateType { // Определение типа пришедшего обновления
    case model.UpdateMessageCreated: // Обработчик новых сообщений
        message := update.GetMessage().Body.Text // полученное сообщение
    }
}
```

Или воспользоваться специальными методами:

```go
func Handle(_ context.Context, update model.Update) {
    switch update.UpdateType { // Определение типа пришедшего обновления
    case model.UpdateMessageCreated: // Обработчик новых сообщений
	    out := "bot прочитал текст: " + update.GetMessage().Body.Text
        cmd := update.GetCommand()
        switch cmd.Command {
            case "/start": // Обработчик команды '/start'
                out = "команда : " + cmd.Command
                /* ... */
            }
        }
    }
}
```

Сравнение текста сообщения со строкой или регулярным выражением производится стандартными средствами golang
Например, пакет strings в Golang, функции Contains

```go
if strings.Contains(update.GetMessage().Body.Text, "hello") {
	/* ... */
}
```

Для обработки нажатия на callback-кнопку с указанным payload используете событие `model.UpdateMessageCallback`:

```go
func Handle(_ context.Context, update model.Update) {
    switch update.UpdateType { // Определение типа пришедшего обновления
    case model.UpdateMessageCallback: // Обработчик нажатия на callback-кнопку с указанным payload
        // Ответ на коллбек
        if update.Callback.Payload == "picture" { // Обработчик callback-кнопки с указанным payload
            /* ... */
        }
    }
}
```

## Отправка сообщений

Вы можете воспользоваться методами:

```go
// Отправить сообщение пользователю с id=12345
maxbot.Messages.Send(ctx, maxbot.NewMessage().SetUser(update.UserID).SetText("Привет!"))
// Отправить сообщение в чат с id=54321
maxbot.Messages.Send(ctx, maxbot.NewMessage().SetChat(54321).SetText("Всем привет!"))

// Получить id отправленного сообщения
message, err := maxbot.Messages.Send(ctx, maxbot.NewMessage().SetChat(54321).SetText("Всем привет!"))
fmt.Printf("message_id: %v", message.MessageID)
```

Отправить ответ на сообщение можно с помощью метода `SetReply`:

```go
maxbot.Messages.Send(ctx, maxbot.NewMessage().SetChat(update.ChatID).SetReply("И вам привет!", update.MessageID))
```

или более короткая форма `Reply`:

```go
maxbot.Messages.Send(ctx, maxbot.NewMessage().Reply("Re: И вам привет!", message)) // reply on reply
```

## Форматирование сообщений

> Подробности про форматирование смотрите в [официальной документации](https://dev.max.ru/).

Вы можете отправлять сообщения, используя **жирный** или _курсивный_ текст, ссылки и многое другое. Есть два типа форматирования: `markdown` и `html`.

### Markdown

```go
maxbot.Messages.Send(ctx, maxbot.NewMessage().SetUser(12345).SetFormat(model.FormatMarkdown).SetText("**Привет!** _Добро пожаловать_ в [Max](https://dev.max.ru)."))
```

### HTML

```go
maxbot.Messages.Send(ctx, maxbot.NewMessage().SetUser(12345).SetFormat(model.FormatHTML).SetText(`<b>Привет!</b> <i>Добро пожаловать</i> в <a href="https://dev.max.ru">Max</a>.`))
```
