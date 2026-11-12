package offline

import (
	"fmt"
	"strings"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

// Dispatch — документный switch-паттерн из docs/02, в виде кода, а не прозы.
// Возвращает человекочитаемое описание апдейта; ни одно поле не берётся "на глаз":
// у message_created текст лежит в GetMessage(), у message_callback — в GetCallback().
func Dispatch(u model.Update) string {
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "%s", u.UpdateType)

	switch u.UpdateType {
	case model.UpdateMessageCreated, model.UpdateMessageEdited:
		m := u.GetMessage()
		cmd := u.GetCommand()
		_, _ = fmt.Fprintf(&b, ": chat_id=%d text=%q", u.ChatID, m.Body.Text)
		if cmd.Command != "" {
			_, _ = fmt.Fprintf(&b, " команда=%s params=%v rest=%q", cmd.Command, cmd.Params, cmd.RemainingText)
		}
	case model.UpdateMessageCallback:
		cb := u.GetCallback()
		payload := u.GetCallbackPayload()
		_, _ = fmt.Fprintf(&b, ": callback_id=%.8s… payload=%q param=%q", cb.CallbackID, cb.Payload, payload.Param)
	case model.UpdateCommentCreated, model.UpdateCommentEdited, model.UpdateCommentRemoved:
		_, _ = fmt.Fprintf(&b, ": post_id=%s comment_id=%s", u.PostID, u.CommentID)
	case model.UpdateMessageRemoved:
		// В событии остаётся только ссылка на удалённое сообщение.
		_, _ = fmt.Fprintf(&b, ": chat_id=%d message_id=%s", u.ChatID, u.MessageID)
	case model.UpdateBotAdded, model.UpdateBotRemoved, model.UpdateUserAdded, model.UpdateUserRemoved:
		_, _ = fmt.Fprintf(&b, ": chat_id=%d user_id=%d", u.ChatID, u.UserID)
	case model.UpdateDialogMuted, model.UpdateDialogUnmuted, model.UpdateDialogCleared,
		model.UpdateDialogRemoved, model.UpdateChatTitleChanged:
		_, _ = fmt.Fprintf(&b, ": chat_id=%d title=%q", u.ChatID, u.GetChat().Title)
	case model.UpdateBotStarted, model.UpdateBotStopped:
		_, _ = fmt.Fprintf(&b, ": user_id=%d", u.UserID)
	default:
		_, _ = fmt.Fprintf(&b, ": (не разобранный тип)")
	}

	return b.String()
}
