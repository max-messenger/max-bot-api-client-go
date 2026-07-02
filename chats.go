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

type Chats struct {
	client *client
}

func (c *Chats) GetChats(ctx context.Context, count, marker int64) (res model.ChatList, err error) {
	values := url.Values{}
	if count > 0 {
		values.Set(paramCount, strconv.Itoa(int(count)))
	}
	if marker > 0 {
		values.Set(paramMarker, strconv.Itoa(int(marker)))
	}

	err = c.client.raw(ctx, http.MethodGet, pathChats, values, nil, &res)

	return
}

func (c *Chats) GetChat(ctx context.Context, chatID int64) (res model.Chat, err error) {
	err = c.client.raw(ctx, http.MethodGet, fmt.Sprintf(formatPathChatsID, chatID), nil, nil, &res)

	return
}

func (c *Chats) EditChat(ctx context.Context, chatID int64, patch model.ChatPatch) (res model.Chat, err error) {
	err = c.client.raw(ctx, http.MethodPatch, fmt.Sprintf(formatPathChatsID, chatID), nil, patch, &res)

	return
}

func (c *Chats) DeleteChat(ctx context.Context, chatID int64) (res model.SimpleQueryResult, err error) {
	err = c.client.raw(ctx, http.MethodDelete, fmt.Sprintf(formatPathChatsID, chatID), nil, nil, &res)

	return
}

func (c *Chats) SendAction(ctx context.Context, chatID int64, action model.SenderAction) (res model.SimpleQueryResult, err error) {
	err = c.client.raw(ctx, http.MethodPost, fmt.Sprintf(formatPathChatsActions, chatID), nil, model.ActionRequestBody{Action: action}, &res)

	return
}

func (c *Chats) GetPinnedMessage(ctx context.Context, chatID int64) (res model.GetPinnedMessageResult, err error) {
	err = c.client.raw(ctx, http.MethodGet, fmt.Sprintf(formatPathChatPin, chatID), nil, nil, &res)

	return
}

func (c *Chats) PinMessage(ctx context.Context, chatID int64, messageID string, notify bool) (res model.SimpleQueryResult, err error) {
	data := model.PinMessageBody{
		MessageID: messageID,
		Notify:    &notify,
	}

	err = c.client.raw(ctx, http.MethodPut, fmt.Sprintf(formatPathChatPin, chatID), nil, data, &res)

	return
}

func (c *Chats) UnpinMessage(ctx context.Context, chatID int64) (res model.SimpleQueryResult, err error) {
	err = c.client.raw(ctx, http.MethodDelete, fmt.Sprintf(formatPathChatPin, chatID), nil, nil, &res)

	return
}

func (c *Chats) GetMembership(ctx context.Context, chatID int64) (res model.ChatMember, err error) {
	err = c.client.raw(ctx, http.MethodGet, fmt.Sprintf(formatPathChatsMembersMe, chatID), nil, nil, &res)

	return
}

func (c *Chats) LeaveChat(ctx context.Context, chatID int64) (res model.SimpleQueryResult, err error) {
	err = c.client.raw(ctx, http.MethodDelete, fmt.Sprintf(formatPathChatsMembersMe, chatID), nil, nil, &res)

	return
}

func (c *Chats) GetAdmins(ctx context.Context, chatID int64) (res model.ChatMembersList, err error) {
	err = c.client.raw(ctx, http.MethodGet, fmt.Sprintf(formatPathChatsMembersAdmin, chatID), nil, nil, &res)

	return
}

func (c *Chats) SetAdmins(ctx context.Context, chatID int64, admins []model.ChatAdmin) (res model.SimpleQueryResult, err error) {
	data := model.ChatAdminsList{
		Admins: admins,
	}
	err = c.client.raw(ctx, http.MethodPost, fmt.Sprintf(formatPathChatsMembersAdmin, chatID), nil, data, &res)

	return
}

func (c *Chats) DeleteAdmins(ctx context.Context, chatID, userID int64) (res model.SimpleQueryResult, err error) {
	err = c.client.raw(ctx, http.MethodDelete, fmt.Sprintf(formatPathChatsMembersAdminDelete, chatID, userID), nil, nil, &res)

	return
}

func (c *Chats) GetMembers(ctx context.Context, chatID, marker, count int64, userIDs []int64) (res model.ChatMembersList, err error) {
	values := url.Values{}

	ids := make([]string, len(userIDs))
	for i, id := range userIDs {
		ids[i] = strconv.FormatInt(id, 10)
	}

	if len(ids) > 0 {
		values.Set(paramUserIDs, strings.Join(ids, ","))
	}
	if count > 0 {
		values.Set(paramCount, strconv.FormatInt(count, 10))
	}
	if marker != 0 {
		values.Set(paramMarker, strconv.FormatInt(marker, 10))
	}

	err = c.client.raw(ctx, http.MethodGet, fmt.Sprintf(formatPathChatsMembers, chatID), values, nil, &res)

	return
}

func (c *Chats) AddMembers(ctx context.Context, chatID int64, userIDs []int64) (res model.SimpleQueryResult, err error) {
	data := model.UserIdsList{
		UserIds: userIDs,
	}
	err = c.client.raw(ctx, http.MethodPost, fmt.Sprintf(formatPathChatsMembers, chatID), nil, data, &res)

	return
}

func (c *Chats) RemoveMember(ctx context.Context, chatID, userID int64, block bool) (res model.SimpleQueryResult, err error) {
	values := url.Values{}
	if userID > 0 {
		values.Set(paramUserID, strconv.FormatInt(userID, 10))
	}
	if block {
		values.Set(paramBlock, strconv.FormatBool(block))
	}

	err = c.client.raw(ctx, http.MethodDelete, fmt.Sprintf(formatPathChatsMembers, chatID), values, nil, &res)

	return
}

func newChats(cli *client) *Chats {
	return &Chats{
		client: cli,
	}
}
