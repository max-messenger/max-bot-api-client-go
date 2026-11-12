# Фикстуры examples/offline

Копии стабов из `github.com/max-messenger/max-bot-api-client-go/v2/stabs/`
(версия v2.4.7). Каталог называется `testdata`, поэтому `go build ./...` и
`go vet ./...` его игнорируют — читать его вправе только код этого пакета.

## Что изменено по сравнению с оригиналом

В оригинале у `mid`/`post_id` значение `mid.fff`. Оно выглядит как настоящий
идентификатор сообщения, но им не является: скопированный в реальный запрос, оно
даёт 404, и читатель принимает это за баг API. Здесь такие значения переименованы в
`mid.NOT-A-REAL-ID`. **Никакое значение из этих файлов нельзя использовать как
образец id** — настоящие id приходят только из ответов API.

## Назначение файлов

| Файл | Что демонстрируется |
|---|---|
| `update.message_created.json` | разбор `message_created` документным switch-паттерном (docs/02) |
| `update.message_callback.json` | `callback_id` + `payload` нажатой кнопки, ответ через `Messages.AnswerOnCallback` |
| `update.comment_created.json` | форма комментария канала: `chat_type=channel` и `post_id` |
| `update.bot_added.json` | событие, из которого читатель достаёт `chat_id` своей группы |
| `update.user_added.json` | `ChatProp.InviterID` |
| `error.attachment.not.ready.json` | ответ, на котором клиент ретраит отправку вложения |
| `error.invalid-token.json` | форма ошибки авторизации (`maxbot.Error.Code`) |
| `botInfo.ok.json` | ответ `/me` |

Уровень доказанности этих файлов — L3: они подтверждают форму запроса и разбора,
но не поведение реального API.
