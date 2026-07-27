# `4` Клавиатура

Для упрощения работы с клавиатурой вы можете использовать `NewKeyboardBuilder`.

```go
	keyboard := model.NewKeyboard()
    keyboard.
        AddRow().   // 1-я строка с 2-мя кнопками
        AddGeoLocation("Прислать геолокацию", true).
        AddContact("Прислать контакт")
    keyboard.
        AddRow().   // 2-я строка с 3-мя кнопками
        AddLink("Открыть Max", "https://max.ru").
        AddCallBack("Аудио", "audio").
        AddCallBack("Видео", "video")
    keyboard.
        AddRow().   // 3-я строка с кнопкой
        AddCallBack("Картинка", "picture")
```

## Типы кнопок

### Callback

```go
func (k *KeyboardRow) AddCallBack(text string, payload string) *KeyboardRow
```

Добавляет callback-кнопку. При нажатии на неё сервер Max отправляет обновление `message_callback`.

### Link

```go
func (k *KeyboardRow) AddLink(text, link string) *KeyboardRow
```

Добавляет кнопку-ссылку. При нажатии на неё пользователю будет предложено открыть ссылку в новой вкладке.

### RequestContact

```go
func (k *KeyboardRow) AddCallBack(text string, payload string) *KeyboardRow
```

Добавляет кнопку запроса контакта. При нажатии на неё боту будет отправлено сообщение с номером телефона, полным именем и почтой пользователя во вложении в формате `VCF`.

### RequestGeoLocation

```go
func (k *KeyboardRow) AddGeoLocation(text string, quick bool) *KeyboardRow
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
msg := maxbot.NewMessage().
    SetText("hello").
    AddKeyboard(keyboard).
    SetChat(update.ChatID).
    SetUser(update.UserID)

res, _ := api.Messages.Send(ctx, msg)
```

Отправляет сообщение в чат с текстом и клавиатурой `keyboard := model.NewKeyboard()`. При нажатии на неё будет создано событие `model.UpdateMessageCallback`.
