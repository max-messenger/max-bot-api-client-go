# `4` Клавиатура

Для упрощения работы с клавиатурой вы можете использовать `NewKeyboardBuilder`.

```go
keyboard := api.Messages.NewKeyboardBuilder()
keyboard.
	AddRow().   // 1-я строка с 2-мя кнопками
	AddGeolocation("Прислать геолокацию", true).
	AddContact("Прислать контакт")
keyboard.
	AddRow().   // 2-я строка с 3-мя кнопками
	AddLink("Открыть Max", schemes.POSITIVE, "https://max.ru").
	AddCallback("Аудио", schemes.NEGATIVE, "audio").
	AddCallback("Видео", schemes.NEGATIVE, "video")
keyboard.
	AddRow().   // 3-я строка с кнопкой
	AddCallback("Картинка", schemes.POSITIVE, "picture")
```

## Типы кнопок

### Callback

```go
// AddCallback button
func (k *KeyboardRow) AddCallback(text string, intent schemes.Intent, payload string) *KeyboardRow
```

Добавляет callback-кнопку. При нажатии на неё сервер Max отправляет обновление `message_callback`.

### Link

```go
// AddLink button
func (k *KeyboardRow) AddLink(text string, intent schemes.Intent, url string) *KeyboardRow
```

Добавляет кнопку-ссылку. При нажатии на неё пользователю будет предложено открыть ссылку в новой вкладке.

### RequestContact

```go
// AddContact button
func (k *KeyboardRow) AddContact(text string) *KeyboardRow
```

Добавляет кнопку запроса контакта. При нажатии на неё боту будет отправлено сообщение с номером телефона, полным именем и почтой пользователя во вложении в формате `VCF`.

### RequestGeoLocation

```go
// AddGeolocation button
func (k *KeyboardRow) AddGeolocation(text string, quick bool) *KeyboardRow
```

Добавляет кнопку запроса геолокации. При нажатии на неё боту будет отправлено сообщение с геолокацией, которую укажет пользователь.

### MessageButton
```go
func (k *KeyboardRow) AddMessage(text string) *KeyboardRow
```

Добавляет текстовую кнопку. При ее нажатии в чат отправится сообщение из кнопки.

### Chat

```go
// Отправка сообщения с клавиатурой
message := maxbot.NewMessage().
	SetChat(upd.Message.Recipient.ChatId).
	AddKeyboard(keyboard).
	SetText(out)
err := api.Messages.Send(ctx, message)
```

Отправляет сообщение в чат с текстом out и клавиатурой `keyboard := api.Messages.NewKeyboardBuilder()`. При нажатии на неё будет создано событие `schemes.MessageCallbackUpdate`.
