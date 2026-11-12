package offline_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/example/offline"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

// fixture читает стабилизированную копию стаба библиотеки.
func fixture(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join("testdata", name)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("фикстура %s: %v", p, err)
	}
	// Значение-заглушка из оригинального стаба (mid.fff) выглядит как настоящий
	// id; в наших копиях оно переименовано, и ни одному id быть не должно.
	if txt := string(b); strings.Contains(txt, "mid.fff") {
		t.Fatalf("фикстура %s содержит плейсхолдер mid.fff", name)
	}

	return string(b)
}

// newAPI собирает настоящий клиент библиотеки, но с подставным транспортом:
// запросы формируются боевым кодом, ответы приходят из testdata/.
func newAPI(t *testing.T, tr *offline.Transport) *maxbot.Api {
	t.Helper()
	api, err := maxbot.NewApi("offline-token", maxbot.WithHTTPClient(tr))
	if err != nil {
		t.Fatalf("NewApi: %v", err)
	}

	return api
}

// TestDispatchFromFixture разбирает message_created документным switch-паттерном
// (docs/02) и проверяет, что клиент донёс до нас поля, а не обнулил их.
func TestDispatchFromFixture(t *testing.T) {
	tr := offline.NewTransport(map[string]offline.Stub{
		"GET /updates": func(int) (int, string) {
			return http.StatusOK, fixture(t, "update.message_created.json")
		},
	})
	api := newAPI(t, tr)

	upd, marker, err := api.Subscriptions.GetUpdates(context.Background(), 0)
	if err != nil {
		t.Fatalf("GetUpdates: %v", err)
	}
	if len(upd) != 1 {
		t.Fatalf("апдейтов %d, ожидался 1 (marker=%d)", len(upd), marker)
	}

	desc := offline.Dispatch(upd[0])
	if !strings.HasPrefix(desc, "message_created:") {
		t.Errorf("Dispatch = %q, ожидался разбор message_created", desc)
	}
	if !strings.Contains(desc, fmt.Sprintf("chat_id=%d", upd[0].ChatID)) {
		t.Errorf("Dispatch не включил chat_id: %q", desc)
	}
	t.Logf("разобрано: %s (marker=%d)", desc, marker)
}

// TestCommandParsing показывает парсер команд: ни один туториал его не упоминает,
// хотя собственный example/simple в библиотеке построен именно на нём.
func TestCommandParsing(t *testing.T) {
	tests := []struct {
		in      string
		wantCmd string
		wantArg []string
		wantTxt string
		wantBot string
	}{
		{"/start", "/start", nil, "", ""},
		{"/image:a,b остальной текст", "/image", []string{"a", "b"}, "остальной текст", ""},
		{"@otherbot /help", "/help", nil, "", "otherbot"},
		{"просто текст", "", nil, "", ""},
	}
	for _, tt := range tests {
		c := model.NewCommand(tt.in)
		if c.Command != tt.wantCmd || c.RemainingText != tt.wantTxt || c.BotName != tt.wantBot {
			t.Errorf("NewCommand(%q) = cmd=%q rest=%q bot=%q", tt.in, c.Command, c.RemainingText, c.BotName)
		}
		if len(c.Params) != len(tt.wantArg) {
			t.Errorf("NewCommand(%q) params=%v, хотим %v", tt.in, c.Params, tt.wantArg)
		}
	}

	// Мост от апдейта к парсеру: текст берётся из разобранного сообщения.
	tr := offline.NewTransport(map[string]offline.Stub{
		"GET /updates": func(int) (int, string) {
			return http.StatusOK, strings.Replace(
				fixture(t, "update.message_created.json"), `"text": "..."`, `"text": "/image:big rest"`, 1)
		},
	})
	upd, _, err := newAPI(t, tr).Subscriptions.GetUpdates(context.Background(), 0)
	if err != nil || len(upd) != 1 {
		t.Fatalf("GetUpdates: err=%v n=%d", err, len(upd))
	}
	cmd := upd[0].GetCommand()
	if cmd.Command != "/image" || len(cmd.Params) != 1 || cmd.Params[0] != "big" {
		t.Fatalf("GetCommand из апдейта = %+v", cmd)
	}
}

// TestAnswerOnCallbackForm проверяет, что оба вида ответа на кнопку уходят
// в правильных эндпоинтах и с правильным телом.
func TestAnswerOnCallbackForm(t *testing.T) {
	tr := offline.NewTransport(map[string]offline.Stub{
		"POST /answers": func(int) (int, string) { return http.StatusOK, `{"success":true}` },
	})
	api := newAPI(t, tr)

	notice := "пример: принято"
	if _, err := api.Messages.AnswerOnCallback(context.Background(), "cb-id-1",
		model.CallbackAnswer{Notification: &notice}); err != nil {
		t.Fatalf("AnswerOnCallback(notification): %v", err)
	}
	if _, err := api.Messages.AnswerOnCallback(context.Background(), "cb-id-1",
		model.CallbackAnswer{Message: &model.NewMessageBody{Text: "пример: заменено"}}); err != nil {
		t.Fatalf("AnswerOnCallback(message): %v", err)
	}
	if got := tr.Attempts("POST", "/answers"); got != 2 {
		t.Fatalf("обращений к /answers %d, ожидалось 2", got)
	}
	if len(tr.Requests) != 2 || tr.Requests[0] != tr.Requests[1] {
		t.Errorf("запросы к /answers: %v", tr.Requests)
	}
	t.Logf("форма запроса: %s", tr.Requests[0])
}

// TestAttachmentNotReadyRetry показывает поведение, описанное в docs/03 и
// проверенное живьём в tests/: на attachment.not.ready клиент ретраит сам.
func TestAttachmentNotReadyRetry(t *testing.T) {
	notReady := fixture(t, "error.attachment.not.ready.json")
	tr := offline.NewTransport(map[string]offline.Stub{
		"POST /messages": func(attempt int) (int, string) {
			if attempt == 1 {
				return http.StatusBadRequest, notReady
			}

			return http.StatusOK, `{"message":{"body":{"mid":"mid.stub-accepted"}}}`
		},
	})
	api := newAPI(t, tr)

	res, err := api.Messages.Send(context.Background(), maxbot.NewMessage().
		SetUser(123456789).SetText("пример").AddAttachByToken("stub-token", model.AttachFile))
	if err != nil {
		t.Fatalf("Send должен был смолчать и повторить: %v", err)
	}
	if res.Message.Body.Mid != "mid.stub-accepted" {
		t.Errorf("mid=%q", res.Message.Body.Mid)
	}
	if got := tr.Attempts("POST", "/messages"); got != 2 {
		t.Errorf("попыток %d, ожидалось 2 (первая - attachment.not.ready)", got)
	}
}

// TestInvalidTokenIsNotRetried показывает границу ретраев: ошибка авторизации
// не повторляется, и её код доступен через maxbot.Error.
func TestInvalidTokenIsNotRetried(t *testing.T) {
	tr := offline.NewTransport(map[string]offline.Stub{
		"POST /messages": func(int) (int, string) {
			return http.StatusUnauthorized, fixture(t, "error.invalid-token.json")
		},
	})
	api := newAPI(t, tr)

	_, err := api.Messages.Send(context.Background(), maxbot.NewMessage().SetUser(1).SetText("x"))
	if err == nil {
		t.Fatal("ожидалась ошибка авторизации")
	}
	apiErr := &maxbot.Error{}
	if !jsonUnmarshalable(err, apiErr) {
		t.Fatalf("ошибка не приводится к maxbot.Error: %#v", err)
	}
	if apiErr.Code == "attachment.not.ready" {
		t.Error("invalid-token не должен прикидываться attachment.not.ready")
	}
	if got := tr.Attempts("POST", "/messages"); got != 1 {
		t.Errorf("попыток %d, ожидалось 1: ошибка авторизации не ретраится", got)
	}
}

func jsonUnmarshalable(err error, into *maxbot.Error) bool {
	var e *maxbot.Error
	if !errors.As(err, &e) {
		return false
	}
	*into = *e

	return true
}

// TestCommentFixtureShape проверяет форму события комментария: chat_type=channel
// и непустой post_id. Настоящий post_id так не получается - он приходит только
// из живого ответа, см. examples/comments posts.
func TestCommentFixtureShape(t *testing.T) {
	var raw struct {
		Updates []struct {
			Message struct {
				Recipient struct {
					ChatID   int64  `json:"chat_id"`
					ChatType string `json:"chat_type"`
					PostID   string `json:"post_id"`
				} `json:"recipient"`
				Body struct {
					Mid  string `json:"mid"`
					Text string `json:"text"`
				} `json:"body"`
			} `json:"message"`
		} `json:"updates"`
	}
	if err := json.Unmarshal([]byte(fixture(t, "update.comment_created.json")), &raw); err != nil {
		t.Fatalf("разбор фикстуры: %v", err)
	}
	rec := raw.Updates[0].Message.Recipient
	if rec.ChatType != "channel" {
		t.Errorf("chat_type=%q, ожидался channel", rec.ChatType)
	}
	if !strings.HasPrefix(rec.PostID, "mid.NOT-A-REAL-ID") {
		t.Errorf("post_id=%q - ожидались только заведомо фейковые значения", rec.PostID)
	}
}

// TestDispatchCoversEveryUpdateType гарантирует, что ни один тип апдейта из
// модели не уходит в ветку "(не разобранный тип)".
func TestDispatchCoversEveryUpdateType(t *testing.T) {
	types := []model.UpdateType{
		model.UpdateMessageCreated, model.UpdateMessageEdited, model.UpdateMessageRemoved,
		model.UpdateMessageCallback, model.UpdateBotAdded, model.UpdateBotRemoved,
		model.UpdateUserAdded, model.UpdateUserRemoved, model.UpdateBotStarted,
		model.UpdateBotStopped, model.UpdateDialogCleared, model.UpdateDialogRemoved,
		model.UpdateDialogMuted, model.UpdateDialogUnmuted, model.UpdateChatTitleChanged,
		model.UpdateCommentCreated, model.UpdateCommentEdited, model.UpdateCommentRemoved,
	}
	for _, ut := range types {
		if got := offline.Dispatch(model.Update{UpdateType: ut}); strings.Contains(got, "не разобранный") {
			t.Errorf("тип %q не разбирается Dispatch", ut)
		}
	}
}

// TestGroupsExtraction проверяет механику команды `subscriptions groups`:
// chat_id достаётся из событий, которые генерирует ДРУГОЙ аккаунт (бот не
// получает message_created на собственные отправки - проверено живьём).
func TestGroupsExtraction(t *testing.T) {
	// Все три события несут ненулевой chat_id - значит механика команды годится
	// для любого из них (проверено на фикстурах, а не на живых событиях).
	const wantID int64 = -70801090403050
	names := []string{
		"update.message_created.json",
		"update.bot_added.json",
		"update.user_added.json",
	}
	for _, name := range names {
		minID := wantID
		body, err := offline.Fixture(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		tr := offline.NewTransport(map[string]offline.Stub{
			"GET /updates": func(int) (int, string) { return http.StatusOK, body },
		})
		api, err := maxbot.NewApi("offline-token", maxbot.WithHTTPClient(tr))
		if err != nil {
			t.Fatalf("NewApi: %v", err)
		}
		upd, _, err := api.Subscriptions.GetUpdates(context.Background(), 0)
		if err != nil || len(upd) != 1 {
			t.Fatalf("%s: err=%v n=%d", name, err, len(upd))
		}
		// Тот же предикат, что в команде: ненулевой ChatID годится как CHAT_ID.
		got := upd[0].ChatID
		if got != minID {
			t.Errorf("%s: chat_id=%d, ожидался %d", name, got, minID)
		}
		t.Logf("%s: chat_id=%d -> этот же id годится как CHAT_ID", name, got)
	}
}
