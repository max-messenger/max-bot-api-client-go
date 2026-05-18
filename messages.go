package maxbot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type Messages struct {
	client *client
}

func (m *Messages) GetMessages(ctx context.Context, chatID, from, to, count int64, messageIDs []string) (res model.MessageList, err error) {
	values := url.Values{}
	if from > to {
		from, to = to, from
	}
	values.Set(paramChatID, strconv.FormatInt(chatID, 10))

	if len(messageIDs) > 0 {
		values.Set(paramMessageIDs, strings.Join(messageIDs, ","))
	}
	if to > 0 {
		values.Set(paramTo, strconv.FormatInt(to, 10))
	}
	if from > 0 {
		values.Set(paramFrom, strconv.FormatInt(from, 10))
	}
	if count > 0 {
		values.Set(paramCount, strconv.FormatInt(count, 10))
	}

	err = m.client.raw(ctx, http.MethodGet, pathMessages, values, nil, &res)

	return
}

func (m *Messages) GetMessageByID(ctx context.Context, messageID string) (res model.Message, err error) {
	err = m.client.raw(ctx, http.MethodGet, fmt.Sprintf(formatPathMessageId, messageID), nil, nil, &res)

	return
}

func (m *Messages) Send(ctx context.Context, msg *Message) (res model.SendMessageResult, err error) {
	if msg == nil {
		err = fmt.Errorf("nil message")

		return
	}
	values := url.Values{}
	if msg.userID > 0 {
		values.Set(paramUserID, strconv.FormatInt(msg.userID, 10))
	}
	if msg.chatID != 0 {
		values.Set(paramChatID, strconv.FormatInt(msg.chatID, 10))
	}
	if msg.disableLinkPreview {
		values.Set(paramDisableLinkPreview, "true")
	}
	err = m.client.rawWithRetry(ctx, http.MethodPost, pathMessages, values, msg.message, &res)

	return
}

func (m *Messages) EditMessage(ctx context.Context, messageID string, body model.NewMessageBody) (res model.SimpleQueryResult, err error) {
	values := url.Values{}
	values.Set(paramMessageID, messageID)
	err = m.client.raw(ctx, http.MethodPut, pathMessages, values, body, &res)

	return
}

func (m *Messages) DeleteMessage(ctx context.Context, messageID string) (res model.SimpleQueryResult, err error) {
	values := url.Values{}
	values.Set(paramMessageID, messageID)
	err = m.client.raw(ctx, http.MethodDelete, pathMessages, values, nil, &res)

	return
}

func (m *Messages) AnswerOnCallback(ctx context.Context, id string, answer model.CallbackAnswer) (res model.SimpleQueryResult, err error) {
	values := url.Values{}
	values.Set(paramCallbackID, id)
	err = m.client.raw(ctx, http.MethodPost, pathAnswers, values, answer, &res)

	return
}

func (m *Messages) GetVideoAttachmentDetails(ctx context.Context, videoToken string) (res model.VideoAttachmentDetails, err error) {
	err = m.client.raw(ctx, http.MethodGet, fmt.Sprintf(formatPathVideoAttachmentDetails, videoToken), nil, nil, &res)

	return
}

func newMessages(client *client) *Messages {
	return &Messages{
		client: client,
	}
}
