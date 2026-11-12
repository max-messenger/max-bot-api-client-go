// Command comments — методы группы CommentsAPI: комментарии к постам публичного
// канала. Путь API — /messages/{postID}/comments, отдельного "COMMENT"-эндпоинта
// нет, поэтому post_id это mid поста канала.
//
// docs-gap: официальная страница «Модерация комментариев к постам в каналах»
// (dev.max.ru/docs-api/use-cases/comment-moderation) в шаге 4 удаляет комментарий
// через `DELETE /comments/{commentId}` — проверено: оба хоста отвечают
// `404 method.not.found`, путь не распознан. Реальный путь (шаг 2 той же страницы
// и comments.go:70) — `DELETE /messages/{messageId}/comments?comment_id=…`.
// docs-gap: та же страница обещает событие `comment_deleted`; в модели клиента
// такого UpdateType нет (model/types.go:100 — только comment_removed).
// docs-gap: все curl-примеры той страницы идут на platform-api2.max.ru, а
// туториалы 01–04 и наш клиент пинят platform-api.max.ru; оба хоста живые.
// docs-gap: в примере шага 5 JSON с висячей запятой ("text": "…", }), он невалиден.
//
// docs-gap: страница предписывает подписку через webhook (POST /subscriptions),
// но события comment_* приходят и поллингом (GetUpdates) — webhook для этого
// не обязателен, см. команду wait.
//
// Требуется публичный канал (CHAN_ID), в который добавлен бот, и включённая
// опция «Комментарии» в настройках канала — через API её включить нельзя.
//
//	go run ./examples/comments posts              # посты канала -> mid (он же post_id)
//	go run ./examples/comments wait               # post_id из события comment_created
//	go run ./examples/comments list <post_id>     # Comments.GetComments
//	go run ./examples/comments get-one <post_id> <comment_id>   # Comments.GetCommentByID
//	go run ./examples/comments roundtrip <post_id>              # Send -> Edit -> Delete
package comments

import (
	"context"
	"fmt"
	"strings"
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

var Commands = []setup.Command{
	{
		Name: "posts", What: "Messages.GetMessages(канал) -> mid как post_id", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken, setup.KeyChannel}, Run: posts,
	},
	{
		Name: "wait", What: "Subscriptions.GetUpdates -> comment_created -> PostID", Proof: setup.L2,
		Live:  false,
		Needs: []string{setup.KeyToken}, Run: wait,
	},
	{
		Name: "list", What: "Comments.GetComments", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken}, Run: list,
	},
	{
		// L3: форма запроса проверяется тестом, живого прогона не было.
		Name: "get-one", What: "Comments.GetCommentByID", Proof: setup.L3,
		Live:  false,
		Needs: []string{setup.KeyToken}, Run: getOne,
	},
	{
		// L3: Send/Edit/Delete блокированы comments.comments_disabled (в канале выключены комментарии).
		Name: "roundtrip", What: "Comments.Send -> Edit -> Delete", Proof: setup.L3,
		Live:  false,
		Needs: []string{setup.KeyToken}, Run: roundtrip,
	},
}

// posts — единственный способ получить post_id без участия человека: канал
// это чат, его сообщения читаются тем же GetMessages.
func posts(_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
	bg := context.Background()
	// Сначала тип чата, потом посты: GetMessages на чужом чате отвечает пусто, и
	// без этой проверки читатель получил бы «канал пуст» вместо «это не канал».
	c, err := api.Chats.GetChat(bg, env.Channel)
	if err != nil {
		return fmt.Errorf("Chats.GetChat(CHAN_ID=%d): %w — нужен chat_id публичного канала", env.Channel, err)
	}
	if c.Type != model.ChatTypeChannel {
		// Проверено прогоном: комментарий к сообщению обычной группы API отвечает
		// channel.not.found, поэтому post_id отсюда взять нельзя.
		fmt.Printf("НЕ ВАЛИДИРОВАНО: CHAN_ID=%d - это %q (type=%q, публичный=%v), а не публичный канал.\n",
			env.Channel, c.Title, c.Type, c.IsPublic)
		fmt.Println("Нужен chat_id публичного канала, в который добавлен бот.")

		return nil
	}
	ml, err := api.Messages.GetMessages(bg, env.Channel, 0, 0, 5, nil)
	if err != nil {
		return fmt.Errorf("Messages.GetMessages(CHAN_ID=%d): %w", env.Channel, err)
	}
	if len(ml.Messages) == 0 {
		fmt.Printf("НЕ ВАЛИДИРОВАНО: канал %q (%d) вернул пустой список постов.\n", c.Title, env.Channel)
		fmt.Println("Второй путь: go run ./examples/comments wait (нужен комментарий человека),")
		fmt.Println("либо опубликуйте пост в канале и повторите posts. Обратите внимание:")
		fmt.Println("GetMessages по группе наблюдался нестабильным, пусто не доказывает отсутствие постов.")

		return nil
	}
	fmt.Printf("постов получено: %d\n", len(ml.Messages))
	for i, m := range ml.Messages {
		// post_id это mid поста; конкретное значение нельзя копировать в другой
		// прогон - оно принадлежит только этому каналу и этому посту.
		fmt.Printf("  [%d] post_id=%s chat_type=%q text=%.40q\n",
			i, m.Body.Mid, m.Recipient.ChatType, m.Body.Text)
	}
	fmt.Println("OK: взять post_id отсюда и передать его остальным командам")

	return nil
}

// wait читает событие comment_created. Событие выдаётся один раз: если его уже
// забрал другой прогон, здесь будет пусто, и это не ошибка примера.
func wait(_ context.Context, api *maxbot.Api, _ setup.Env, _ []string) error {
	fmt.Println("GetUpdates потребляет очередь: комментарий должен быть оставлен после запуска команды.")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var marker int64
	for i := 0; i < 20; i++ {
		upd, next, err := api.Subscriptions.GetUpdates(ctx, marker)
		if err != nil {
			return fmt.Errorf("GetUpdates: %w", err)
		}
		marker = next
		for _, u := range upd {
			switch u.UpdateType {
			case model.UpdateCommentCreated:
				fmt.Printf("post_id=%s comment_id=%s -> go run ./examples/comments list %s\n",
					u.PostID, u.CommentID, u.PostID)
				fmt.Println("OK: пост найден в событии")

				return nil
			default:
				fmt.Printf("  (%s пропущен)\n", u.UpdateType)
			}
		}
	}
	fmt.Println("НЕ ВАЛИДИРОВАНО: за выборку не пришло ни одного comment_created.")
	fmt.Println("Причины: бот не состоит в канале, комментарии к постам недоступны боту,")
	fmt.Println("или событие уже потреблено другим проголом. Тогда post_id берётся из posts.")

	return nil
}

// postID достаёт post_id из аргумента; без аргумента — не угадывает.
func postID(args []string, min int) (string, error) {
	if len(args) < min {
		return "", fmt.Errorf("нужен post_id (см. go run ./examples/comments posts), аргументов %d", len(args))
	}

	return args[0], nil
}

func list(_ context.Context, api *maxbot.Api, _ setup.Env, args []string) error {
	pid, err := postID(args, 1)
	if err != nil {
		return err
	}
	cl, err := api.Comments.GetComments(context.Background(), pid, 0, 0, 10, nil)
	if err != nil {
		return fmt.Errorf("Comments.GetComments(post_id=%s): %w — убедитесь, что это mid поста канала", pid, err)
	}
	// Ловушка формы ответа: список называется CommentList, но ключ json — "messages",
	// а id комментария лежит в поле CommentBody.CommentID (тег "mid").
	fmt.Printf("OK: запрос выполнен, комментариев: %d\n", len(cl.Messages))
	for i, c := range cl.Messages {
		fmt.Printf("  [%d] comment_id=%s post_id=%s text=%.40q\n",
			i, c.Body.CommentID, c.Recipient.PostID, c.Body.Text)
	}
	// Пустой ответ здесь - подтверждённое чтение, а не отсутствие среды:
	// в отличие от нестабильного GetMessages, этот эндпоинт отдавал 0 стабильно.
	if len(cl.Messages) == 0 {
		fmt.Println("  (пусто, но запрос прошёл успешно: это подтверждённое чтение, а не отсутствие среды)")
	}

	return nil
}

func getOne(_ context.Context, api *maxbot.Api, _ setup.Env, args []string) error {
	pid, err := postID(args, 2)
	if err != nil {
		return err
	}
	c, err := api.Comments.GetCommentByID(context.Background(), pid, args[1])
	if err != nil {
		return fmt.Errorf("Comments.GetCommentByID(%s, %s): %w", pid, args[1], err)
	}
	fmt.Printf("comment_id=%s post_id=%s text=%q timestamp=%d\n",
		c.Body.CommentID, c.Recipient.PostID, c.Body.Text, c.Timestamp)

	return nil
}

// roundtrip создаёт комментарий, правит его и удаляет: среда подготовки не
// должна забиваться мусором от примеров.
func roundtrip(_ context.Context, api *maxbot.Api, _ setup.Env, args []string) error {
	pid, err := postID(args, 1)
	if err != nil {
		return err
	}
	bg := context.Background()
	text := fmt.Sprintf("пример-комментарий %d", time.Now().Unix())

	res, err := api.Comments.Send(bg, pid, maxbot.NewComment(text))
	if err != nil {
		if strings.Contains(err.Error(), "comments.comments_disabled") {
			return fmt.Errorf("Comments.Send(post_id=%s): %w — в канале выключены комментарии.\n"+
				"Включить можно только в интерфейсе MAX: канал -> аватар канала -> шестерёнка -> "+
				"«Комментарии» (метода API для этого нет). Пока опция выключена, бот читает, "+
				"правит и удаляет только старые комментарии (GetComments работает)", pid, err)
		}

		return fmt.Errorf("Comments.Send(post_id=%s): %w", pid, err)
	}
	mid := res.Message.Body.Mid
	if mid == "" {
		return fmt.Errorf("Comments.Send не вернул comment_id (body.mid пуст)")
	}
	fmt.Printf("создан комментарий mid=%s text=%q\n", mid, text)
	defer func() {
		if _, err := api.Comments.Delete(bg, pid, mid); err != nil {
			fmt.Printf("cleanup: Comments.Delete(%s) не удался: %v\n", mid, err)

			return
		}
		fmt.Printf("cleanup: комментарий удалён mid=%s\n", mid)
	}()

	edited := "отредактировано " + text
	if _, err := api.Comments.Edit(bg, pid, mid, maxbot.NewComment(edited)); err != nil {
		return fmt.Errorf("Comments.Edit(%s): %w", mid, err)
	}
	after, err := api.Comments.GetCommentByID(bg, pid, mid)
	if err != nil {
		return fmt.Errorf("Comments.GetCommentByID после правки: %w", err)
	}
	if after.Body.Text != edited {
		return fmt.Errorf("после Edit текст %q, ожидалось %q", after.Body.Text, edited)
	}
	fmt.Printf("OK: правка подтверждена чтением, текст=%q\n", after.Body.Text)

	return nil
}
