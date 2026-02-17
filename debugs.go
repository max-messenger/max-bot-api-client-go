package maxbot

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type debugs struct {
	client *client
	chat   int64
}

func newDebugs(client *client, chat int64) *debugs {
	return &debugs{client: client, chat: chat}
}

// Send sends a message to a chat. As a result for this method new message identifier returns.
func (a *debugs) Send(ctx context.Context, upd schemes.UpdateInterface) error {
	return a.sendMessage(ctx, false, a.chat, 0, &schemes.NewMessageBody{Text: upd.GetDebugRaw()})
}

// SendErr sends a message to a chat. As a result for this method new message identifier returns.
func (a *debugs) SendErr(ctx context.Context, err error) error {
	return a.sendMessage(ctx, false, a.chat, 0, &schemes.NewMessageBody{Text: err.Error()})
}

func (a *debugs) sendMessage(ctx context.Context, reset bool, chatID int64, userID int64, message *schemes.NewMessageBody) error {
	result := new(schemes.Error)
	values := url.Values{}
	if chatID != 0 {
		values.Set(paramChatID, strconv.Itoa(int(chatID)))
	}
	if userID != 0 {
		values.Set(paramUserID, strconv.Itoa(int(userID)))
	}

	body, err := a.client.request(ctx, http.MethodPost, pathMessages, values, reset, message)
	if err != nil {
		return err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if err = json.NewDecoder(body).Decode(result); err != nil {
		return nil
	}
	if result.Code == "" {
		return nil
	}

	return result
}
