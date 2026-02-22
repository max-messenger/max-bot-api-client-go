# `2` Прослушивание обновлений и реакция на них

После запуска бота Max начнёт отправлять вам обновления.

> Подробности обо всех обновлениях смотрите в [официальной документации](https://dev.max.ru/).

Max Bot API позволяет прослушивать эти обновления, например:

```go
for upd := range api.GetUpdates(ctx) { // Чтение из канала с обновлениями
	switch upd := upd.(type) { // Определение типа пришедшего обновления
	case *schemes.BotStartedUpdate:   // Обработчик начала диалога с ботом
		/* ... */
	case *schemes.MessageCreatedUpdate: // Обработчик новых сообщений
		/* ... */
	case *schemes.UserAddedToChatUpdate: // Обработчик добавления пользователя в беседу
		/* ... */
	}
}
```

Вы можете посмотреть модуль schemes, чтобы увидеть все доступные типы обновлений. UpdateType содержит актуальный список типов.

## Получение сообщений

Вы можете подписаться на обновление `message_created`:

```go
for update := range api.GetUpdates(ctx) { // Чтение из канала с обновлениями
	switch upd := update.(type) { // Определение типа пришедшего обновления
	case *schemes.MessageCreatedUpdate: // Обработчик новых сообщений
		message := upd.Message // полученное сообщение
	}
}
```

Или воспользоваться специальными методами:

```go
for update := range api.GetUpdates(ctx) { // Чтение из канала с обновлениями
	switch upd := update.(type) { // Определение типа пришедшего обновления
	case *schemes.MessageCreatedUpdate: // Обработчик новых сообщений и команд
		out := "bot прочитал текст: " + upd.Message.Body.Text
		switch upd.GetCommand() {
		case "/start": // Обработчик команды '/start'
 			out = "команда : " + upd.GetCommand()
			/* ... */
		}
	}
}
```

Аналогичный код со специальным методом GetText

```go
for update := range api.GetUpdates(ctx) { // Чтение из канала с обновлениями
	switch upd := update.(type) { // Определение типа пришедшего обновления
	case *schemes.MessageCreatedUpdate: // Обработчик новых сообщений и команд
		out := "bot прочитал текст: " + upd.GetText()
		switch upd.GetCommand() {
		case "/start": // Обработчик команды '/start'
			out = "команда : " + upd.GetCommand()
			/* ... */
		}
	}
}
```

Сравнение текста сообщения со строкой или регулярным выражением производится стандартными средствами golang
Например, пакет strings в Golang, функции Contains

```go
if strings.Contains(upd.GetText(), "hello") {
	/* ... */
}
```

Для обработки нажатия на callback-кнопку с указанным payload используете событие schemes.MessageCallbackUpdate:

```go
for update := range api.GetUpdates(ctx) { // Чтение из канала с обновлениями
	switch upd := update.(type) { // Определение типа пришедшего обновления
	case *schemes.MessageCallbackUpdate: // Обработчик нажатия на callback-кнопку с указанным payload
		// Ответ на коллбек
		if upd.Callback.Payload == "picture" { // Обработчик callback-кнопки с указанным payload
			/* ... */
		}
	}
}
```

## Отправка сообщений

Вы можете воспользоваться методами:

```go
// Отправить сообщение пользователю с id=12345
api.Messages.Send(ctx, maxbot.NewMessage().SetUser(12345).SetText("Привет!"))
// Отправить сообщение в чат с id=54321
api.Messages.Send(ctx, maxbot.NewMessage().SetChat(54321).SetText("Всем привет!"))

// Получить id отправленного сообщения
id, err := api.Messages.Send(ctx, maxbot.NewMessage().SetChat(54321).SetText("Всем привет!"))
fmt.Printf("message_id: %v", id)
// Получить полный ответ со структурой schemes.Message
message, err := api.Messages.SendMessageResult(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText("Всем привет!"))
fmt.Printf("message_id: %v", message.Body.Mid)
```

Отправить ответ на сообщение можно с помощью метода `SetReply`:

```go
message, err := api.Messages.SendMessageResult(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetReply("И вам привет!", upd.Message.Body.Mid))
```

или более короткая форма `Reply`:

```go
message, err := api.Messages.SendMessageResult(ctx, maxbot.NewMessage().Reply("И вам привет! на ", upd.Message))
api.Messages.Send(ctx, maxbot.NewMessage().Reply("Re: И вам привет!", message)) // reply on reply
```

## Форматирование сообщений

> Подробности про форматирование смотрите в [официальной документации](https://dev.max.ru/).

Вы можете отправлять сообщения, используя **жирный** или _курсивный_ текст, ссылки и многое другое. Есть два типа форматирования: `markdown` и `html`.

### Markdown

```go
message := maxbot.NewMessage().SetUser(12345).SetText('**Привет!** _Добро пожаловать_ в [Max](https://dev.max.ru).').SetFormat("markdown")
```

### HTML

```go
message := maxbot.NewMessage().SetUser(12345).SetText('<b>Привет!</b> <i>Добро пожаловать</i> в <a href="https://dev.max.ru">Max</a>.').SetFormat('html'),
```
