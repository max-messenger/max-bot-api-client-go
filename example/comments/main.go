// Command comments — примеры методов клиента max-bot-api-client-go/v2: Comments.* + Messages.GetMessages для post_id.
//
// Запуск: go run ./comments <команда> [аргументы...]
//
// Реестр команд и их тела лежат в internal/examplecmds/comments: они должны быть
// импортируемы, чтобы тест мог перечислить команды и прогнать их на подставном
// транспорте. Точка входа остаётся здесь.
package main

import (
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/comments"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
)

func main() {
	setup.Run("comments", comments.Commands...)
}
