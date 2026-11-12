// Command subscriptions — приём апдейтов через polling (webhook в объём не входит).
//
//	go run ./examples/subscriptions poll       # Subscriptions.GetUpdates, два прохода
//	go run ./examples/subscriptions options    # WithPollingPause / WithPollingTimeout
//	go run ./examples/subscriptions list       # Subscriptions.GetSubscriptions
//	go run ./examples/subscriptions groups     # вспомогательная: chat_id из событий
package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

const queueWarning = "ВНИМАНИЕ: GetUpdates потребляет очередь апдейтов: событие выдаётся один раз " +
	"и при следующем вызове его уже не будет (проверено прогоном). Не запускайте эти примеры, " +
	"если бот живёт на webhook или поллинге в проде, и не ожидайте, что вторую команду события увидят."

var Commands = []setup.Command{
	{
		Name: "poll", What: "Subscriptions.GetUpdates (2 прохода по маркеру)", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken},
		Run:   poll,
	},
	{
		Name: "options", What: "WithPollingPause/WithPollingTimeout на одном проходе", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken},
		Run:   options,
	},
	{
		Name: "list", What: "Subscriptions.GetSubscriptions", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken},
		Run:   list,
	},
	{
		// L3: форма запроса проверяется тестом, живого прогона не было.
		Name: "groups", What: "вспомогательная на GetUpdates: chat_id из событий", Proof: setup.L3,
		Live:  false,
		Needs: []string{setup.KeyToken},
		Run:   groups,
	},
}

// twoPass делает ровно два вызова GetUpdates: маркер 0, затем маркер из ответа.
func twoPass(ctx context.Context, api *maxbot.Api) ([]model.Update, int64, int64, error) {
	upd1, marker1, err := api.Subscriptions.GetUpdates(ctx, 0)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("GetUpdates(marker=0): %w", err)
	}
	upd2, marker2, err := api.Subscriptions.GetUpdates(ctx, marker1)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("GetUpdates(marker=%d): %w", marker1, err)
	}

	return append(upd1, upd2...), marker1, marker2, nil
}

func poll(_ context.Context, api *maxbot.Api, _ setup.Env, _ []string) error {
	fmt.Println(queueWarning)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	upd, m1, m2, err := twoPass(ctx, api)
	if err != nil {
		return err
	}
	if m2 < m1 {
		return fmt.Errorf("маркер уменьшился: %d -> %d", m1, m2)
	}
	fmt.Printf("проход 1: marker=%d, проход 2: marker=%d, апдейтов суммарно=%d\n", m1, m2, len(upd))
	for i, u := range upd {
		fmt.Printf("  [%d] %s chat_id=%d user_id=%d is_channel=%v\n", i, u.UpdateType, u.ChatID, u.UserID, u.IsChannel)
	}
	if len(upd) == 0 {
		fmt.Println("Очередь пуста или события уже потреблены предыдущим запуском — это ожидаемо:")
		fmt.Println("ассертом этого прогона считается успешное получение маркера, а не число событий.")
	}
	fmt.Println("OK: оба вызова завершены без ошибки, маркер не уменьшился")

	return nil
}

// options показывает, чем управляется ожидание на пустой очереди: WithPollingTimeout
// уходит в long-poll как параметр ?timeout=, и он же решает, когда библиотека
// успеет вернуться до своих трёх ретраев. Меряем ПЕРВЫЙ заход с контекстом 12s:
// без опции полный прогон (30s × 3 попытки) уходит далеко за минуту и меряет
// уже network, а не опцию.
func options(_ context.Context, _ *maxbot.Api, env setup.Env, _ []string) error {
	fmt.Println(queueWarning)
	fmt.Println("Замер первого захода (контекст 12s); полный прогон с ретраями дольше.")

	var def, fast time.Duration
	for _, tuned := range []bool{false, true} {
		name, opts := "дефолт библиотеки", []maxbot.Opt{maxbot.WithBaseURL(setup.DocsBaseURL)}
		if tuned {
			name = "пример (timeout 2s)"
			opts = append(opts,
				maxbot.WithPollingPause(200*time.Millisecond), maxbot.WithPollingTimeout(2*time.Second))
		}
		api, err := maxbot.NewApi(env.Token, opts...)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		start := time.Now()
		upd, marker, err := api.Subscriptions.GetUpdates(ctx, 0)
		took := time.Since(start).Round(time.Millisecond)
		cancel()
		if tuned {
			fast = took
		} else {
			def = took
		}
		fmt.Printf("%-24s первый заход на пустой очереди=%-8v апдейтов=%d marker=%d err=%v\n",
			name, took, len(upd), marker, err)
		// Long-poll на пустой очереди штатно уходит в таймаут: это ожидаемый
		// результат замера, а не падение примера.
		if err != nil {
			var te *maxbot.TimeoutError
			if !errors.Is(err, context.DeadlineExceeded) && !errors.As(err, &te) {
				return fmt.Errorf("GetUpdates (%s): %w", name, err)
			}
		}
	}
	if fast >= def {
		return fmt.Errorf("WithPollingTimeout(2s) не быстрее дефолта: %v vs %v", fast, def)
	}
	fmt.Printf("OK: с опцией первый заход вернулся за %v, без него висит %v\n", fast, def)

	return nil
}

func list(_ context.Context, api *maxbot.Api, _ setup.Env, _ []string) error {
	res, err := api.Subscriptions.GetSubscriptions(context.Background())
	if err != nil {
		return fmt.Errorf("Subscriptions.GetSubscriptions: %w", err)
	}
	fmt.Printf("активных подписок (webhook): %d\n", len(res.Subscriptions))
	for _, s := range res.Subscriptions {
		// URL подписки может вести на внутренний хост — печатаем только домен-заглушку.
		fmt.Printf("  update_types=%d version=%q time=%d\n", len(s.UpdateTypes), s.Version, s.Time)
	}
	fmt.Println("OK: список подписок получен; Subscribe/Unsubscribe в объёме примеров не показываются")

	return nil
}

// groups — единственный доступный читателю способ узнать chat_id своей группы:
// метода GetChats в клиенте нет, id достаётся из события.
func groups(_ context.Context, api *maxbot.Api, _ setup.Env, _ []string) error {
	fmt.Println(queueWarning)
	fmt.Println("Добавьте бота в группу и напишите ему там сообщение, пока идёт выборка.")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	found := map[int64]string{}
	var marker int64
	for i := 0; i < 12; i++ {
		upd, next, err := api.Subscriptions.GetUpdates(ctx, marker)
		if err != nil {
			return err
		}
		marker = next
		for _, u := range upd {
			switch u.UpdateType {
			case model.UpdateBotAdded, model.UpdateMessageCreated, model.UpdateUserAdded,
				model.UpdateBotStarted, model.UpdateChatTitleChanged:
				if u.ChatID != 0 {
					found[u.ChatID] = string(u.UpdateType)
				}
			}
		}
		if len(found) > 0 {
			break
		}
		select {
		case <-ctx.Done():
		case <-time.After(1 * time.Second):
		}
	}

	if len(found) == 0 {
		fmt.Println("НЕ ВАЛИДИРОВАНО: за выборку не попало ни одного события с chat_id.")
		fmt.Println("Причины: бот ни в одной группе, либо события уже потреблены предыдущим")
		fmt.Println("запуском poll/options (GetUpdates отдаёт событие ровно один раз).")
		fmt.Println("Иначе - id группы можно выписать из клиента MAX и положить в CHAT_ID.")

		return nil
	}
	fmt.Println("Найденные chat_id — подпишите нужный как CHAT_ID (отрицательное число):")
	for id, via := range found {
		fmt.Printf("  CHAT_ID=%d   (из события %s)\n", id, via)
	}

	return nil
}
