package maxbot

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type messages struct {
	client *client
}

func newMessages(client *client) *messages {
	return &messages{client: client}
}

// GetMessages returns messages in chat: result page and marker referencing to the next page. Messages traversed in reverse direction so the latest message in chat will be first in result array. Therefore if you use from and to parameters, to must be less than from
func (a *messages) GetMessages(ctx context.Context, chatID int64, messageIDs []string, from int, to int, count int) (*schemes.MessageList, error) {
	result := new(schemes.MessageList)
	values := url.Values{}
	if chatID != 0 {
		values.Set("chat_id", strconv.Itoa(int(chatID)))
	}
	if len(messageIDs) > 0 {
		for _, mid := range messageIDs {
			values.Add("message_ids", mid)
		}
	}
	if from != 0 {
		values.Set("from", strconv.Itoa(from))
	}
	if to != 0 {
		values.Set("to", strconv.Itoa(to))
	}
	if count > 0 {
		values.Set("count", strconv.Itoa(count))
	}
	body, err := a.client.request(ctx, http.MethodGet, "messages", values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *messages) GetMessage(ctx context.Context, messageID string) (*schemes.Message, error) {
	result := new(schemes.Message)
	path := "messages/" + url.PathEscape(messageID)
	body, err := a.client.request(ctx, http.MethodGet, path, nil, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// EditMessage updates message by id
func (a *messages) EditMessage(ctx context.Context, messageID string, message *Message) error {
	s, err := a.editMessage(ctx, messageID, message.message)
	if err != nil {
		return err
	}
	if !s.Success {
		return errors.New(s.Message)
	}
	return nil
}

// DeleteMessage deletes message by id
func (a *messages) DeleteMessage(ctx context.Context, messageID string) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	values.Set("message_id", messageID)
	body, err := a.client.request(ctx, http.MethodDelete, "messages", values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// AnswerOnCallback should be called to send an answer after a user has clicked the button. The answer may be an updated message or/and a one-time user notification.
func (a *messages) AnswerOnCallback(ctx context.Context, callbackID string, callback *schemes.CallbackAnswer) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	values.Set("callback_id", callbackID)
	body, err := a.client.request(ctx, http.MethodPost, "answers", values, false, callback)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// NewKeyboardBuilder returns new keyboard builder helper
func (a *messages) NewKeyboardBuilder() *Keyboard {
	return &Keyboard{
		rows: make([]*KeyboardRow, 0),
	}
}

// Send sends a message to a chat. As a result for this method new message identifier returns.
func (a *messages) Send(ctx context.Context, m *Message) error {
	_, err := a.sendMessage(ctx, m.reset, m.chatID, m.userID, m.message)
	return err
}

// SendWithResult sends a message to a chat and returns the created message along with any error.
func (a *messages) SendWithResult(ctx context.Context, m *Message) (*schemes.Message, error) {
	return a.sendMessage(ctx, m.reset, m.chatID, m.userID, m.message)
}

func (a *messages) sendMessage(ctx context.Context, reset bool, chatID int64, userID int64, message *schemes.NewMessageBody) (*schemes.Message, error) {
	result := new(schemes.Message)
	values := url.Values{}
	if chatID != 0 {
		values.Set("chat_id", strconv.Itoa(int(chatID)))
	}
	if userID != 0 {
		values.Set("user_id", strconv.Itoa(int(userID)))
	}

	body, err := a.client.request(ctx, http.MethodPost, "messages", values, reset, message)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if err := json.NewDecoder(body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *messages) editMessage(ctx context.Context, messageID string, message *schemes.NewMessageBody) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	values.Set("message_id", messageID)
	body, err := a.client.request(ctx, http.MethodPut, "messages", values, false, message)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// Check posiable to send a message to a chat.
func (a *messages) Check(ctx context.Context, m *Message) (bool, error) {
	return a.checkUser(ctx, m.reset, m.message)
}

func (a *messages) checkUser(ctx context.Context, reset bool, message *schemes.NewMessageBody) (bool, error) {
	result := new(schemes.Error)
	values := url.Values{}
	if reset {
		values.Set("access_token", message.BotToken)
	}
	mode := "notify/exists"

	if message.PhoneNumbers != nil {
		values.Set("phone_numbers", strings.Join(message.PhoneNumbers, ","))
	}

	body, err := a.client.request(ctx, http.MethodGet, mode, values, reset, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if err := json.NewDecoder(body).Decode(result); err != nil {
		return false, err
	}

	if len(result.NumberExist) > 0 {
		return true, result
	}

	return false, result
}

// Check posiable to send a message to a chat.
func (a *messages) ListExist(ctx context.Context, m *Message) ([]string, error) {
	return a.checkNumberExist(ctx, m.reset, m.message)
}

func (a *messages) checkNumberExist(ctx context.Context, reset bool, message *schemes.NewMessageBody) ([]string, error) {
	result := new(schemes.Error)
	values := url.Values{}
	if reset {
		values.Set("access_token", message.BotToken)
	}
	mode := "notify/exists"

	if message.PhoneNumbers != nil {
		values.Set("phone_numbers", strings.Join(message.PhoneNumbers, ","))
	}

	body, err := a.client.request(ctx, http.MethodGet, mode, values, reset, nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	if err := json.NewDecoder(body).Decode(result); err != nil {
		// Message sent without errors
		return nil, err
	}
	if len(result.NumberExist) > 0 {
		return result.NumberExist, result
	}
	return nil, result
}
