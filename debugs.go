package maxbot

import (
	"context"
	"encoding/json"
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
func (a *debugs) Send(ctx context.Context, upd schemes.UpdateInterface) (string, error) {
	return a.sendMessage(ctx, false, false, a.chat, 0, &schemes.NewMessageBody{Text: upd.GetDebugRaw()})
}

// SendErr sends a message to a chat. As a result for this method new message identifier returns.
func (a *debugs) SendErr(ctx context.Context, err error) (string, error) {
	return a.sendMessage(ctx, false, false, a.chat, 0, &schemes.NewMessageBody{Text: err.Error()})
}

func (a *debugs) sendMessage(ctx context.Context, vip bool, reset bool, chatID int64, userID int64, message *schemes.NewMessageBody) (string, error) {
	result := new(schemes.Error)
	values := url.Values{}
	if chatID != 0 {
		values.Set(paramChatID, strconv.Itoa(int(chatID)))
	}
	if userID != 0 {
		values.Set(paramUserID, strconv.Itoa(int(userID)))
	}
	if reset {
		values.Set(paramAccessToken, message.BotToken)
	}
	mode := "messages"
	if vip {
		mode = "notify"
	}
	body, err := a.client.request(ctx, http.MethodPost, mode, values, reset, message)
	if err != nil {
		return "heir", err
	}
	defer body.Close()
	if err := json.NewDecoder(body).Decode(result); err != nil {
		// Message sent without errors
		return "err", err
	}
	if result.Code == "" {
		if mode == "notify" {
			return "ok", result
		} else {
			return "", nil
		}

	}
	return "", result
}
