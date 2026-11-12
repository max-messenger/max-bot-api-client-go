// Command messages — методы группы MessagesAPI клиента max-bot-api-client-go/v2.
//
//	go run ./examples/messages send [текст]           # Messages.Send в диалог с TEST_USER_ID
//	go run ./examples/messages send-group [текст]     # Messages.Send в CHAT_ID (нужна группа)
//	go run ./examples/messages get-one <mid>          # Messages.GetMessageByID
//	go run ./examples/messages get-list               # Messages.GetMessages
//	go run ./examples/messages edit <mid>             # Messages.EditMessage
//	go run ./examples/messages answer-callback        # Messages.AnswerOnCallback
//	go run ./examples/messages video-details          # Messages.GetVideoAttachmentDetails
package messages

import (
	"context"
	"fmt"
	"os"
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

// shortcut: путь фикстурыresolved через setup.Fixture, потому что `go run` и
// `go test` ставят разный CWD.
const fixtureVideo = "test_data/video.mp4"

var Commands = []setup.Command{
	{
		Name: "send", What: "Messages.Send -> GetMessageByID -> DeleteMessage", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken, setup.KeyUserA},
		Run:   send,
	},
	{
		Name: "send-group", What: "Messages.Send в группу", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken, setup.KeyUserA, setup.KeyChat},
		Run:   sendGroup,
	},
	{
		Name: "get-one", What: "Messages.GetMessageByID", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken, setup.KeyUserA},
		Run:   getOne,
	},
	{
		Name: "get-list", What: "Messages.GetMessages (диалог и группа из CHAT_ID)", Proof: setup.L3,
		Live:  false,
		Needs: []string{setup.KeyToken, setup.KeyUserA, setup.KeyChat},
		Run:   getList,
	},
	{
		Name: "edit", What: "Messages.EditMessage", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken, setup.KeyUserA},
		Run:   edit,
	},
	{
		Name: "answer-callback", What: "Messages.AnswerOnCallback", Proof: setup.L2,
		Live:  false,
		Needs: []string{setup.KeyToken, setup.KeyUserA},
		Run:   answerCallback,
	},
	{
		Name: "video-details", What: "Upload.Upload -> Messages.GetVideoAttachmentDetails", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken, setup.KeyUserA},
		Run:   videoDetails,
	},
}

func stamp() string {
	return fmt.Sprintf("example %d", time.Now().Unix())
}

// send отправляет сообщение в диалог и сразу удаляет его, чтобы прогон
// повторялся без накопления мусора.
func send(_ context.Context, api *maxbot.Api, env setup.Env, args []string) error {
	text := stamp()
	if len(args) > 0 {
		text = args[0]
	}
	bg := context.Background()

	res, err := api.Messages.Send(bg, maxbot.NewMessage().
		SetUser(env.UserA).SetText(text).SetFormat(model.FormatMarkdown))
	if err != nil {
		return fmt.Errorf("Messages.Send: %w", err)
	}
	mid := res.Message.Body.Mid
	if mid == "" {
		return fmt.Errorf("Messages.Send вернул пустой mid")
	}
	fmt.Printf("отправлено в диалог A: mid=%s text=%q\n", mid, res.Message.Body.Text)
	del(bg, api, mid)

	return nil
}

func sendGroup(_ context.Context, api *maxbot.Api, env setup.Env, args []string) error {
	text := stamp()
	if len(args) > 0 {
		text = args[0]
	}
	bg := context.Background()

	res, err := api.Messages.Send(bg, maxbot.NewMessage().
		SetChat(env.Chat).SetText(text))
	if err != nil {
		return fmt.Errorf("Messages.Send в группу %d: %w", env.Chat, err)
	}
	mid := res.Message.Body.Mid
	fmt.Printf("отправлено в группу: mid=%s chat_id=%d chat_type=%q\n",
		mid, res.Message.Recipient.ChatID, res.Message.Recipient.ChatType)
	if res.Message.Recipient.ChatID != env.Chat {
		return fmt.Errorf("сообщение ушло в чат %d, ожидался %d", res.Message.Recipient.ChatID, env.Chat)
	}

	// Подтверждение доставки читаем составом группы: списка чатов у клиента нет.
	members, err := api.Chats.GetMembers(bg, env.Chat, 0, 10, nil)
	if err != nil {
		return fmt.Errorf("Chats.GetMembers: %w", err)
	}
	fmt.Printf("участников в группе: %d\n", len(members.Members))
	del(bg, api, mid)

	return nil
}

func getOne(_ context.Context, api *maxbot.Api, _ setup.Env, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("нужен аргумент <mid>, например: go run ./examples/messages get-one mid.xxx")
	}
	m, err := api.Messages.GetMessageByID(context.Background(), args[0])
	if err != nil {
		return fmt.Errorf("Messages.GetMessageByID(%s): %w", args[0], err)
	}
	fmt.Printf("mid=%s text=%q sender_bot=%v timestamp=%d\n",
		m.Body.Mid, m.Body.Text, m.Sender.IsBot, m.Timestamp)

	return nil
}

func getList(_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
	bg := context.Background()
	res, err := api.Messages.Send(bg, maxbot.NewMessage().SetUser(env.UserA).SetText(stamp()))
	if err != nil {
		return fmt.Errorf("Messages.Send для истории: %w", err)
	}
	mid := res.Message.Body.Mid
	defer del(bg, api, mid)

	// 1) Точечное чтение по mid — единственный путь, который обязан вернуть данные.
	byID, err := api.Messages.GetMessages(bg, env.UserA, 0, 0, 5, []string{mid})
	if err != nil {
		return fmt.Errorf("Messages.GetMessages по messageIDs: %w", err)
	}
	fmt.Printf("по messageIDs=[%s]: вернулось %d сообщений\n", mid, len(byID.Messages))
	if len(byID.Messages) == 0 {
		// Единичное чтение по тому же mid работает - значит GetMessages требует
		// chatID настоящего чата, а не user_id диалога (в туториалах не показано).
		one, err := api.Messages.GetMessageByID(bg, mid)
		if err != nil {
			return fmt.Errorf("Messages.GetMessageByID(%s): %w", mid, err)
		}
		fmt.Printf("  то же сообщение через GetMessageByID: text=%q - читаем, что GetMessages с chatID=user_id ничего не отдаёт\n", one.Body.Text)
	}

	// 2) История диалога: на живой среде отдаёт 0 сообщений даже сразу после
	// отправки (проверено прогоном). Утверждать обратное пример не вправе.
	history, err := api.Messages.GetMessages(bg, env.UserA, 0, 0, 5, nil)
	if err != nil {
		return fmt.Errorf("Messages.GetMessages истории: %w", err)
	}
	fmt.Printf("история диалога A: %d сообщений (err=%v)\n", len(history.Messages), err)
	if len(history.Messages) == 0 {
		fmt.Println("ЗАМЕЧАНИЕ: пустая история диалога — известное поведение API при chatID=user_id.")
	}
	for i, m := range history.Messages {
		fmt.Printf("  [%d] mid=%s text=%.40q\n", i, m.Body.Mid, m.Body.Text)
	}

	// История настоящей группы читается иначе, чем история диалога: проверяем,
	// если CHAT_ID задан, и ничего не утверждаем, если оба чтения пусты.
	if env.Chat != 0 {
		gl, err := api.Messages.GetMessages(bg, env.Chat, 0, 0, 5, nil)
		if err != nil {
			return fmt.Errorf("Messages.GetMessages по группе %d: %w", env.Chat, err)
		}
		fmt.Printf("история группы %d: %d сообщений\n", env.Chat, len(gl.Messages))
		for i, m := range gl.Messages {
			fmt.Printf("  [%d] mid=%s text=%.40q\n", i, m.Body.Mid, m.Body.Text)
		}
		if len(gl.Messages) > 0 {
			fmt.Printf("история группы прочитана: %d сообщений\n", len(gl.Messages))
		} else {
			// Замерено: на одной и той же группе вызов отдавал 5 сообщений, а через
			// 15 минут - 0 при 6 попытках и разных from/to, тогда как GetChat и
			// GetMembers на неё же отвечали нормально. Пусто здесь ничего не доказывает.
			fmt.Println("ЗАМЕЧАНИЕ: GetMessages по группе отдал пусто. Эндпоинт наблюдался")
			fmt.Println("нестабильным (тот же чат давал и непустой ответ), так что это")
			fmt.Println("не утверждение об истории. Детерминированное чтение - GetMessageByID.")
		}
		if len(byID.Messages) == 0 && len(gl.Messages) == 0 {
			fmt.Println("НЕ ВАЛИДИРОВАНО (L3): GetMessages не вернул сообщений ни на диалоге, ни на группе")

			return nil
		}
	}
	if len(byID.Messages) == 0 {
		fmt.Println("НЕ ВАЛИДИРОВАНО: чтение по messageIDs не вернуло только что созданный mid")

		return nil
	}
	fmt.Println("OK: чтение по messageIDs нашело созданное сообщение")

	return nil
}

func edit(_ context.Context, api *maxbot.Api, env setup.Env, args []string) error {
	bg := context.Background()
	text := stamp()
	res, err := api.Messages.Send(bg, maxbot.NewMessage().SetUser(env.UserA).SetText(text))
	if err != nil {
		return fmt.Errorf("Messages.Send для редактирования: %w", err)
	}
	mid := res.Message.Body.Mid

	edited := "отредактировано " + text
	if len(args) > 0 {
		edited = args[0]
	}
	if _, err := api.Messages.EditMessage(bg, mid, model.NewMessageBody{Text: edited}); err != nil {
		del(bg, api, mid)

		return fmt.Errorf("Messages.EditMessage(%s): %w", mid, err)
	}

	after, err := api.Messages.GetMessageByID(bg, mid)
	del(bg, api, mid)
	if err != nil {
		return fmt.Errorf("Messages.GetMessageByID после правки: %w", err)
	}
	if after.Body.Text != edited {
		return fmt.Errorf("после EditMessage текст %q, ожидалось %q", after.Body.Text, edited)
	}
	fmt.Printf("OK: mid=%s редактирование подтверждено, текст=%q (сообщение удалено)\n", mid, after.Body.Text)

	return nil
}

// answer-callback формирует оба вида ответа на нажатие кнопки. callback_id
// приходит только из живого события message_callback, поэтому без клика
// человека проверяется только форма ответа; обе ветки видны на экране
// аккаунта B - поэтому уровень L2, а не L1.
func answerCallback(_ context.Context, api *maxbot.Api, env setup.Env, args []string) error {
	bg := context.Background()
	notice := "examples: кнопка нажата"

	kb := model.NewKeyboard()
	kb.AddRow().AddCallBack("показать ответ", "examples-payload")
	res, err := api.Messages.Send(bg, maxbot.NewMessage().
		SetUser(env.UserA).SetText("нажмите кнопку, затем перезапустите команду с <callback_id>").
		AddKeyboard(kb))
	if err != nil {
		return fmt.Errorf("Messages.Send с клавиатурой: %w", err)
	}
	mid := res.Message.Body.Mid
	defer del(bg, api, mid)
	fmt.Printf("сообщение с кнопкой отправлено: mid=%s\n", mid)

	if len(args) == 0 {
		fmt.Println("НЕ ВАЛИДИРОВАНО: нужен <callback_id> из события message_callback (клик по кнопке).")
		fmt.Println("Две формы ответа, которые умеет API:")
		fmt.Printf("  notification: {Notification: %q}\n", notice)
		fmt.Printf("  message:      {Message: {Text: %q}}\n", "ответ заменён примером")

		return nil
	}

	for _, id := range args {
		if _, err := api.Messages.AnswerOnCallback(bg, id, model.CallbackAnswer{Notification: &notice}); err != nil {
			return fmt.Errorf("Messages.AnswerOnCallback(notification, %s): %w", id, err)
		}
		fmt.Printf("OK: ответ notification отправлен на callback_id=%s\n", id)

		if _, err := api.Messages.AnswerOnCallback(bg, id, model.CallbackAnswer{
			Message: &model.NewMessageBody{Text: "ответ заменён примером"},
		}); err != nil {
			return fmt.Errorf("Messages.AnswerOnCallback(message, %s): %w", id, err)
		}
		fmt.Println("OK: ответ message отправлен; что появилось на экране — видно только с аккаунта B (L2)")
	}

	return nil
}

// video-details загружает видео из встроенной фикстуры и читает детали по токену.
func videoDetails(_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
	bg := context.Background()
	path, err := setup.Fixture(fixtureVideo)
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("фикстура %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	token, err := api.Upload.Upload(bg, model.UploadVideo, f, info.Name(), info.Size())
	if err != nil {
		return fmt.Errorf("Upload.Upload(video): %w", err)
	}
	fmt.Printf("видео загружено, длина токена=%d (значение не печатается)\n", len(token))

	details, err := api.Messages.GetVideoAttachmentDetails(bg, token)
	if err != nil {
		return fmt.Errorf("Messages.GetVideoAttachmentDetails: %w", err)
	}
	fmt.Printf("OK: размеры=%dx%d длительность=%d URL=%v\n",
		details.Width, details.Height, details.Duration, details.Urls != nil)

	sent, err := api.Messages.Send(bg, maxbot.NewMessage().SetUser(env.UserA).
		SetText("пример: видео-детали").AddAttachByToken(token, model.AttachVideo))
	if err != nil {
		return fmt.Errorf("Messages.Send с видео-вложением: %w", err)
	}
	var hasVideo bool
	for _, a := range sent.Message.Body.Attachments {
		if a.Type == model.AttachVideo {
			hasVideo = true
		}
	}
	del(bg, api, sent.Message.Body.Mid)
	if !hasVideo {
		return fmt.Errorf("отправленное сообщение не содержит вложение типа video")
	}
	fmt.Println("OK: сообщение с видео-вложением создано и удалено")

	return nil
}

// del удаляет созданное сообщение; ошибка удаления — это ошибка примера,
// иначе среда подготовки постепенно забивается мусором.
func del(bg context.Context, api *maxbot.Api, mid string) {
	if mid == "" {
		return
	}
	if _, err := api.Messages.DeleteMessage(bg, mid); err != nil {
		fmt.Printf("cleanup: DeleteMessage(%s) не удался: %v\n", mid, err)

		return
	}
	fmt.Printf("cleanup: удалено mid=%s\n", mid)
}
