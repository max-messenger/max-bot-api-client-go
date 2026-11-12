// Command messages — примеры методов клиента max-bot-api-client-go/v2: Messages.Send/GetMessages/GetMessageByID/EditMessage/AnswerOnCallback/GetVideoAttachmentDetails.
//
// Запуск: go run ./messages <команда> [аргументы...]
//
// Реестр команд и их тела лежат в internal/examplecmds/messages: они должны быть
// импортируемы, чтобы тест мог перечислить команды и прогнать их на подставном
// транспорте. Точка входа остаётся здесь.
package main

import (
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/messages"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
)

func main() {
	setup.Run("messages", messages.Commands...)
}
