package maxbot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestChats(t *testing.T) {
	suite.Run(t, new(chatsTest))
}

type chatsTest struct {
	suite.Suite
}

func (t *chatsTest) SetupTest() {

}

func (t *chatsTest) TestGetChat() {
	data, err := stabs.ReadFile("stabs/chats/get-chat-by-id.json")
	t.NoError(err)

	expect := model.Chat{
		ChatID:            -70000000000005,
		Type:              model.ChatTypeChat,
		Status:            model.ChatStatusActive,
		Title:             "chat title",
		LastEventTime:     1775628268494,
		ParticipantsCount: 3,
		IsPublic:          false,
		Description:       "chat description",
		OwnerID:           123123123,
		Link:              "https://max.ru/join/hash-chat-link",
		MessagesCount:     7,
		Participants: map[string]int64{
			"123123123": 1775628268494,
			"201201201": 1775628268494,
			"229229229": 1775628268494,
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.GetChat(context.Background(), -70000000000005)
	t.NoError(err)

	t.Equal(expect, res)
}

func (t *chatsTest) TestEditChat() {
	data, err := stabs.ReadFile("stabs/chats/chat-path-result.json")
	t.NoError(err)

	expectResult := model.Chat{
		ChatID:            -70000000000005,
		Type:              model.ChatTypeChat,
		Status:            model.ChatStatusActive,
		Title:             "chat title",
		LastEventTime:     1775626919577,
		ParticipantsCount: 3,
		IsPublic:          false,
		Description:       "chat description",
		OwnerID:           123123123,
		Link:              "https://max.ru/join/hash-chat-link",
		MessagesCount:     7,
		PinnedMessage: model.Message{
			Timestamp: 1775626919577,
			Body: model.MessageBody{
				Mid:  "mid.sha-more",
				Seq:  116328127605864690,
				Text: "hi",
			},
			Recipient: model.Recipient{
				ChatID:   -70000000000005,
				ChatType: model.ChatTypeChat,
			},
			Sender: model.User{
				UserID:           123123123,
				FirstName:        "John",
				LastName:         "Doe",
				IsBot:            false,
				LastActivityTime: 1775626919577,
				Name:             "John Doe",
			},
		},
	}

	chat := model.ChatPatch{
		Title: "chat title",
		Pin:   "mid.sha-more",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		rData, _ := io.ReadAll(r.Body)
		defer func() { _ = r.Body.Close() }()
		obj := model.ChatPatch{}

		t.NoError(json.Unmarshal(rData, &obj))
		t.Equal(obj, chat)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.EditChat(context.Background(), -70000000000005, chat)
	t.NoError(err)

	t.Equal(expectResult, res)
}

func (t *chatsTest) TestDeleteChat() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodDelete)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsID, -70000000000005))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.DeleteChat(context.Background(), -70000000000005)
	t.NoError(err)

	t.True(res.Success)
}

func (t *chatsTest) TestSendAction() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPost)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsActions, -70000000000005))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.SendAction(context.Background(), -70000000000005, model.ActionSendingPhoto)
	t.NoError(err)

	t.True(res.Success)
}

func (t *chatsTest) TestGetPinMessage() {
	data, err := stabs.ReadFile("stabs/chats/get-pin-message.json")
	t.NoError(err)

	expect := model.GetPinnedMessageResult{
		Message: model.Message{
			Timestamp: 1775626919577,
			Recipient: model.Recipient{
				ChatID:   -70000000000005,
				ChatType: model.ChatTypeChat,
			},
			Body: model.MessageBody{
				Mid:  "mid.sha-more",
				Seq:  116328127605864690,
				Text: "hi",
			},
			Sender: model.User{
				UserID:           123123123,
				FirstName:        "John",
				LastName:         "Doe",
				IsBot:            false,
				LastActivityTime: 1775628268494,
				Name:             "John Doe",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatPin, -70000000000005))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.GetPinnedMessage(context.Background(), -70000000000005)
	t.NoError(err)

	t.Equal(expect, res)
}

func (t *chatsTest) TestPinMessage() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPut)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatPin, -70000000000005))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.PinMessage(context.Background(), -70000000000005, "", true)
	t.NoError(err)

	t.True(res.Success)
}

func (t *chatsTest) TestUnpinMessage() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodDelete)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatPin, -70000000000005))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.UnpinMessage(context.Background(), -70000000000005)
	t.NoError(err)

	t.True(res.Success)
}

func (t *chatsTest) TestGetMembership() {
	data, err := stabs.ReadFile("stabs/chats/get-membership.json")
	t.NoError(err)

	expect := model.ChatMember{
		LastAccessTime: 1775626919577,
		IsOwner:        false,
		IsAdmin:        true,
		JoinTime:       1775626919577,
		Permissions: []model.ChatAdminPermission{
			model.PermAddAdmins,
			model.PermEditLink,
			model.PermPinMessage,
			model.PermChangeChatInfo,
			model.PermAddRemoveMembers,
			model.PermCanCall,
			model.PermWrite,
			model.PermReadAllMessages,
		},
		UserID:           229229229,
		FirstName:        "bot firstname",
		Username:         "test-bot",
		IsBot:            true,
		LastActivityTime: 1775626919577,
		Description:      "bot for unit-test",
		AvatarURL:        "https://localhost/avatar.png",
		FullAvatarURL:    "https://localhost/avatar.full.png",
		Name:             "bot firstname",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsMembersMe, -70000000000005))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.GetMembership(context.Background(), -70000000000005)
	t.NoError(err)

	t.Equal(expect, res)
}

func (t *chatsTest) TestLeaveChat() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodDelete)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsMembersMe, -70000000000005))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.LeaveChat(context.Background(), -70000000000005)
	t.NoError(err)

	t.True(res.Success)
}

func (t *chatsTest) TestGetAdmins() {
	data, err := stabs.ReadFile("stabs/chats/get-admins.json")
	t.NoError(err)

	expect := model.ChatMembersList{
		Members: []model.ChatMember{
			{
				LastAccessTime: 1775626919577,
				IsOwner:        false,
				IsAdmin:        true,
				JoinTime:       1775626919577,
				Permissions: []model.ChatAdminPermission{
					model.PermAddAdmins,
					model.PermEditLink,
					model.PermPinMessage,
					model.PermChangeChatInfo,
					model.PermAddRemoveMembers,
					model.PermCanCall,
					model.PermWrite,
					model.PermReadAllMessages,
				},
				UserID:           229229229,
				FirstName:        "bot firstname",
				Username:         "test-bot",
				IsBot:            true,
				LastActivityTime: 1775626919577,
				Description:      "bot for unit-test",
				AvatarURL:        "https://localhost/avatar.png",
				FullAvatarURL:    "https://localhost/avatar.full.png",
				Name:             "bot firstname",
			},
			{
				LastAccessTime: 1775626919577,
				IsOwner:        true,
				IsAdmin:        true,
				JoinTime:       1775626919577,
				Permissions: []model.ChatAdminPermission{
					model.PermAddAdmins,
					model.PermEditLink,
					model.PermEdit,
					model.PermPinMessage,
					model.PermChangeChatInfo,
					model.PermAddRemoveMembers,
					model.PermCanCall,
					model.PermDelete,
					model.PermWrite,
					model.PermReadAllMessages,
					model.PermViewStats,
				},
				UserID:           123123123,
				FirstName:        "John",
				LastName:         "Doe",
				IsBot:            false,
				LastActivityTime: 1775626919577,
				AvatarURL:        "https://localhost/avatar.png",
				FullAvatarURL:    "https://localhost/avatar.full.png",
				Name:             "John Doe",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsMembersAdmin, -70000000000005))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.GetAdmins(context.Background(), -70000000000005)
	t.NoError(err)

	t.Equal(expect, res)
}

func (t *chatsTest) TestSetAdmins() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPost)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsMembersAdmin, -70000000000005))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	admins := []model.ChatAdmin{
		{
			UserID: 123123123,
			Permissions: []model.ChatAdminPermission{
				model.PermPinMessage,
			},
			Alias: "Дежурный",
		},
	}
	res, err := api.Chats.SetAdmins(context.Background(), -70000000000005, admins)
	t.NoError(err)

	t.True(res.Success)
}

func (t *chatsTest) TestDeleteAdmins() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodDelete)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsMembersAdminDelete, -70000000000005, 123123123))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.DeleteAdmins(context.Background(), -70000000000005, 123123123)
	t.NoError(err)

	t.True(res.Success)
}

func (t *chatsTest) TestGetMembers() {
	data, err := stabs.ReadFile("stabs/chats/get-members.json")
	t.NoError(err)

	expect := model.ChatMembersList{
		Members: []model.ChatMember{
			{
				LastAccessTime: 1775626919577,
				IsOwner:        false,
				IsAdmin:        true,
				JoinTime:       1775626919577,
				Permissions: []model.ChatAdminPermission{
					model.PermAddAdmins,
					model.PermEditLink,
					model.PermPinMessage,
					model.PermChangeChatInfo,
					model.PermAddRemoveMembers,
					model.PermCanCall,
					model.PermWrite,
					model.PermReadAllMessages,
				},
				UserID:           229229229,
				FirstName:        "bot firstname",
				Username:         "test-bot",
				IsBot:            true,
				LastActivityTime: 1775626919577,
				Description:      "bot for unit-test",
				AvatarURL:        "https://localhost/avatar.png",
				FullAvatarURL:    "https://localhost/avatar.full.png",
				Name:             "bot firstname",
			},
			{
				LastAccessTime: 1775626919577,
				IsOwner:        true,
				IsAdmin:        true,
				JoinTime:       1775626919577,
				Permissions: []model.ChatAdminPermission{
					model.PermAddAdmins,
					model.PermEditLink,
					model.PermEdit,
					model.PermPinMessage,
					model.PermChangeChatInfo,
					model.PermAddRemoveMembers,
					model.PermCanCall,
					model.PermDelete,
					model.PermWrite,
					model.PermReadAllMessages,
					model.PermViewStats,
				},
				UserID:           123123123,
				FirstName:        "John",
				LastName:         "Doe",
				IsBot:            false,
				LastActivityTime: 1775626919577,
				AvatarURL:        "https://localhost/avatar.png",
				FullAvatarURL:    "https://localhost/avatar.full.png",
				Name:             "John Doe",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsMembers, -70000000000005))
		t.Equal(r.RequestURI, "/chats/-70000000000005/members?count=100&marker=89&user_ids=1%2C2%2C3")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.GetMembers(context.Background(), -70000000000005, 89, 100, []int64{1, 2, 3})
	t.NoError(err)

	t.Equal(expect, res)
}

func (t *chatsTest) TestAddMembers() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPost)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsMembers, -70000000000005))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.AddMembers(context.Background(), -70000000000005, []int64{123123123})
	t.NoError(err)

	t.True(res.Success)
}

func (t *chatsTest) TestRemoveMember() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodDelete)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathChatsMembers, -70000000000005))
		t.Equal(r.RequestURI, "/chats/-70000000000005/members?block=true&user_id=123123123")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Chats.RemoveMember(context.Background(), -70000000000005, 123123123, true)
	t.NoError(err)

	t.True(res.Success)
}
