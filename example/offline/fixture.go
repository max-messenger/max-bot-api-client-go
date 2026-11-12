package offline

import (
	"os"

	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
)

// Fixture читает стабилизированную копию стаба библиотеки из testdata этого
// пакета. Путь резолвится от корня модуля, поэтому читателю не нужно запускать
// пример из какой-то конкретной директории.
//
// Значения id в этих файлах заведомо фейковые: использовать их как образцы
// нельзя, см. testdata/README.md.
func Fixture(name string) (string, error) {
	path, err := setup.Fixture("examples/offline/testdata/" + name)
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
