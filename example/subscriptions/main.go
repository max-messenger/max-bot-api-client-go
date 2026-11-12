// Command subscriptions — примеры методов клиента max-bot-api-client-go/v2: Subscriptions.GetUpdates/GetSubscriptions.
//
// Запуск: go run ./subscriptions <команда> [аргументы...]
//
// Реестр команд и их тела лежат в internal/examplecmds/subscriptions: они должны быть
// импортируемы, чтобы тест мог перечислить команды и прогнать их на подставном
// транспорте. Точка входа остаётся здесь.
package main

import (
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/subscriptions"
)

func main() {
	setup.Run("subscriptions", subscriptions.Commands...)
}
