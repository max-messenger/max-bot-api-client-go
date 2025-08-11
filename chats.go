package maxbot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type chats struct {
	client *client
}

func newChats(client *client) *chats {
	return &chats{client: client}
}

// GetChats returns information about chats that bot participated in: a result list and marker points to the next page
func (a *chats) GetChats(ctx context.Context, count, marker int64) (*schemes.ChatList, error) {
	result := new(schemes.ChatList)
	values := url.Values{}
	if count > 0 {
		values.Set(paramCount, strconv.Itoa(int(count)))
	}
	if marker > 0 {
		values.Set(paramMarker, strconv.Itoa(int(marker)))
	}
	body, err := a.client.request(ctx, http.MethodGet, "chats", values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// GetChat returns info about chat
func (a *chats) GetChat(ctx context.Context, chatID int64) (*schemes.Chat, error) {
	result := new(schemes.Chat)
	values := url.Values{}
	body, err := a.client.request(ctx, http.MethodGet, fmt.Sprintf("chats/%d", chatID), values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// GetChatMembership returns chat membership info for current bot
func (a *chats) GetChatMembership(ctx context.Context, chatID int64) (*schemes.ChatMember, error) {
	result := new(schemes.ChatMember)
	values := url.Values{}
	body, err := a.client.request(ctx, http.MethodGet, fmt.Sprintf("chats/%d/members/me", chatID), values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// GetChatMembers returns users participated in chat
func (a *chats) GetChatMembers(ctx context.Context, chatID, count, marker int64) (*schemes.ChatMembersList, error) {
	result := new(schemes.ChatMembersList)
	values := url.Values{}
	if count > 0 {
		values.Set(paramCount, strconv.Itoa(int(count)))
	}
	if marker > 0 {
		values.Set(paramMarker, strconv.Itoa(int(marker)))
	}
	body, err := a.client.request(ctx, http.MethodGet, fmt.Sprintf("chats/%d/members", chatID), values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// LeaveChat removes bot from chat members
func (a *chats) LeaveChat(ctx context.Context, chatID int64) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	body, err := a.client.request(ctx, http.MethodDelete, fmt.Sprintf("chats/%d/members/me", chatID), values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// EditChat edits chat info: title, icon, etc…
func (a *chats) EditChat(ctx context.Context, chatID int64, update *schemes.ChatPatch) (*schemes.Chat, error) {
	result := new(schemes.Chat)
	values := url.Values{}
	body, err := a.client.request(ctx, http.MethodPatch, fmt.Sprintf("chats/%d", chatID), values, false, update)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// AddMember adds members to chat. Additional permissions may require.
func (a *chats) AddMember(ctx context.Context, chatID int64, users schemes.UserIdsList) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	body, err := a.client.request(ctx, http.MethodPost, fmt.Sprintf("chats/%d/members", chatID), values, false, users)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// RemoveMember removes member from chat. Additional permissions may require.
func (a *chats) RemoveMember(ctx context.Context, chatID int64, userID int64) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	values.Set(paramUserID, strconv.Itoa(int(userID)))
	body, err := a.client.request(ctx, http.MethodDelete, fmt.Sprintf("chats/%d/members", chatID), values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// SendAction send bot action to chat
func (a *chats) SendAction(ctx context.Context, chatID int64, action schemes.SenderAction) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	body, err := a.client.request(ctx, http.MethodPost, fmt.Sprintf("chats/%d/actions", chatID), values, false, schemes.ActionRequestBody{Action: action})
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}
