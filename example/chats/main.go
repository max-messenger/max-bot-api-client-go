// Command chats — примеры методов клиента max-bot-api-client-go/v2: Chats.* (15 методов).
//
// Запуск: go run ./chats <команда> [аргументы...]
//
// Реестр команд и их тела лежат в internal/examplecmds/chats: они должны быть
// импортируемы, чтобы тест мог перечислить команды и прогнать их на подставном
// транспорте. Точка входа остаётся здесь.
package main

import (
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/chats"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
)

func main() {
	setup.Run("chats", chats.Commands...)
}
