// Command upload — примеры методов клиента max-bot-api-client-go/v2: Upload.Upload.
//
// Запуск: go run ./upload <команда> [аргументы...]
//
// Реестр команд и их тела лежат в internal/examplecmds/upload: они должны быть
// импортируемы, чтобы тест мог перечислить команды и прогнать их на подставном
// транспорте. Точка входа остаётся здесь.
package main

import (
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/upload"
)

func main() {
	setup.Run("upload", upload.Commands...)
}
