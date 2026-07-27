# `1` Ваш первый бот

Для начала вам нужно получить токен бота. Для этого авторизуйтесь на портале [Max для партнеров](https://business.max.ru/), создайте нового бота, следуя инструкциям, и получите его токен в разделе «Интеграция» после модерации.

Создайте новый проект `my-first-bot` и установите `github.com/max-messenger/max-bot-api-client-go/v2`. Для этого откройте терминал и выполните следующие команды:

```sh
# Создайте новую папку для исходного кода вашего модуля Go и перейдите в неё
mkdir my-first-bot
cd my-first-bot
# Запустите свой модуль с помощью команды go mod init
go mod init first-max-bot

# Установите библиотеку для работы с MAX API на golang
go get github.com/max-messenger/max-bot-api-client-go/v2
```

Команда go mod init создает файл go.mod для отслеживания зависимостей вашего кода. Пока что файл включает только имя вашего модуля и версию Go, которую поддерживает ваш код.

Теперь создайте файл, например, `bot.go`.

Весь код в языке Go организуется в пакеты. Пакеты представляют удобную организацию разделения кода на отдельные части или модули. Модульность позволяет определять один раз пакет с нужной функциональностью и потом использовать его многократно в различных программах.
Код пакета располагается в одном или нескольких файлах с расширением go. Для определения пакета применяется ключевое слово package. Поэтому наш файл `bot.go` будет иметь следующую структуру.

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello Max Bot Go")
}
```

В данном случае пакет называется main. Определение пакета должно идти в начале файла.

Есть два типа пакетов: исполняемые (executable) и библиотеки (reusable). Для создания исполняемых файлов пакет должен иметь имя main. Все остальные пакеты не являются исполняемыми. При этом пакет main должен содержать функцию main, которая является входной точкой в приложение.

Импортируем в наш пакет main установленный модуль `github.com/max-messenger/max-bot-api-client-go/v2`

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
)

func main() {
	ctx := context.Background()
	opts := []maxbot.Opt{
		maxbot.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	}
	api, err := maxbot.NewApi(os.Getenv("TOKEN"), opts...)
	if err != nil {
		fmt.Println("api initial err:", err)
		return
	}
	// Some methods demo:
	info, err := api.Bots.GetMyInfo(ctx)
	fmt.Printf("Get me: %#v %#v", info, err)
}
```

Код выше, создает объект `api`, передавая токен (и набор опций для конфигурирования) в конструктор NewApi. 
Мы рекомендуем передавать токен через переменные окружения, т.к. использовать токен в коде - плохая практика.

Данная программа выведет только информацию о вашем боте и закончит работу.
Чтобы бот заработал необходим обработчик событий из канала с обновлениями

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []maxbot.Opt{
		maxbot.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	}

	api, err := maxbot.NewApi(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		log.Fatal(err)
	}

	info, err := api.Bots.GetMyInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("info: %+v", info)

	handle := func(ctx context.Context, update model.Update) {
		fmt.Printf("Received: [%s] %#v\n", update.UpdateType, update)
		switch update.UpdateType {
		case model.UpdateMessageCreated:
			msg := maxbot.NewMessage().
				SetText("Hello from Bot").
				SetChat(update.ChatID).
				SetUser(update.UserID)

			res, cErr := api.Messages.Send(ctx, msg)
			if cErr != nil {
				/* ... */
			}
			log.Printf("%v\n", res)
		}
	}

	var updates []model.Update
	var marker int64
	for {
		select {
		case <-ctx.Done():
		default:
			updates, marker, err = api.Subscriptions.GetUpdates(ctx, marker)
			if err != nil {
				log.Println("GetUpdates: ", err)
				return
			}

			for _, update := range updates {
				handle(ctx, update)
			}
		}
	}
}
```

Поздравляем, вы написали первого бота! 🎉
