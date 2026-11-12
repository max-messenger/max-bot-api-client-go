// Command demo — герметичный прогон разбора апдейтов и формы ответов без сети.
//
//	go run ./examples/offline/demo
//
// Работает без BOT_TOKEN и без доступа к API: запросы формируются боевым
// клиентом библиотеки, а ответы приходят из examples/offline/testdata/.
// Уровень доказанности — L3 (форма запроса), поведение реального API не проверяется.
package main

import (
	"context"
	"fmt"
	"net/http"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/example/offline"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func main() {
	names := []string{
		"update.message_created.json",
		"update.message_callback.json",
		"update.comment_created.json",
		"update.bot_added.json",
	}
	for _, n := range names {
		body, err := offline.Fixture(n)
		if err != nil {
			fmt.Println("ошибка:", err)

			return
		}
		upd, err := updatesFrom(body)
		if err != nil {
			fmt.Printf("%s: %v\n", n, err)

			return
		}
		fmt.Printf("%-30s -> %s\n", n, offline.Dispatch(upd[0]))
	}

	shapes, err := requestShapes()
	if err != nil {
		fmt.Println("ошибка:", err)

		return
	}
	fmt.Println("формы запросов, которые генерирует клиент (L3):")
	for _, s := range shapes {
		fmt.Printf("  %s\n", s)
	}
}

// updatesFrom прогоняет тело ответа через настоящий GetUpdates клиента:
// разбор делает библиотека, а не пример.
func updatesFrom(body string) ([]model.Update, error) {
	tr := offline.NewTransport(map[string]offline.Stub{
		"GET /updates": func(int) (int, string) { return http.StatusOK, body },
	})
	api, err := maxbot.NewApi("offline-token", maxbot.WithHTTPClient(tr))
	if err != nil {
		return nil, err
	}

	upd, _, err := api.Subscriptions.GetUpdates(context.Background(), 0)
	if err != nil {
		return nil, err
	}

	return upd, nil
}

// requestShapes показывает, какие HTTP-запросы клиент строит для ответа на
// кнопку и для отправки вложения с ретраем на attachment.not.ready.
func requestShapes() ([]string, error) {
	notReady, err := offline.Fixture("error.attachment.not.ready.json")
	if err != nil {
		return nil, err
	}
	tr := offline.NewTransport(map[string]offline.Stub{
		"POST /answers": func(int) (int, string) { return http.StatusOK, `{"success":true}` },
		"POST /messages": func(attempt int) (int, string) {
			if attempt == 1 {
				return http.StatusBadRequest, notReady
			}

			return http.StatusOK, `{"message":{"body":{"mid":"mid.stub-accepted"}}}`
		},
	})
	api, err := maxbot.NewApi("offline-token", maxbot.WithHTTPClient(tr))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	notice := "пример: принято"
	if _, err := api.Messages.AnswerOnCallback(ctx, "cb-id",
		model.CallbackAnswer{Notification: &notice}); err != nil {
		return nil, err
	}
	if _, err := api.Messages.Send(ctx, maxbot.NewMessage().
		SetUser(123456789).SetText("пример").AddAttachByToken("stub", model.AttachFile)); err != nil {
		return nil, err
	}

	return tr.Requests, nil
}
