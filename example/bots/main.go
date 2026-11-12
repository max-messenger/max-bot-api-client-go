// Command bots — примеры методов клиента max-bot-api-client-go/v2: Bots.GetMyInfo, Bots.PatchCommands, Bots.EditMyInfo.
//
// Запуск: go run ./bots <команда> [аргументы...]
//
// Реестр команд и их тела лежат в internal/examplecmds/bots: они должны быть
// импортируемы, чтобы тест мог перечислить команды и прогнать их на подставном
// транспорте. Точка входа остаётся здесь.
package main

import (
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/bots"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
)

func main() {
	setup.Run("bots", bots.Commands...)
}
