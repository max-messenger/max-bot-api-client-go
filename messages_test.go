package maxbot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestMessages(t *testing.T) {
	suite.Run(t, new(messagesTest))
}

type messagesTest struct {
	suite.Suite
}

func (t *messagesTest) SetupTest() {

}

func (t *messagesTest) TestGetMessages() {
	data, err := stabs.ReadFile("stabs/messages/get-messages.json")
	t.NoError(err)

	expect := model.MessageList{
		Messages: []model.Message{
			{
				Sender: model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775053249000,
				},
				Recipient: model.Recipient{
					ChatID:   182182182,
					ChatType: model.ChatTypeChat,
				},
				Timestamp: 1775053255737,
				Body: model.MessageBody{
					Mid:  "mid.ffffbdb48e6c3775019d496b34394b84",
					Seq:  116327994376978687,
					Text: "...",
				},
				Stat: model.MessageStat{
					Views: 1,
				},
				URL: "https://max.ru/chats/182182182/messages/mid.ffffbdb48e6c3775019d496b34394b84",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, pathMessages)

		query := r.URL.Query()
		t.Equal("182182182", query.Get(paramChatID))
		t.Equal("1775025604499", query.Get(paramFrom))
		t.Equal("1775025580000", query.Get(paramTo))
		t.Equal("50", query.Get(paramCount))
		t.Equal("mid.a,mid.b", query.Get(paramMessageIDs))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Messages.GetMessages(context.Background(), 182182182, 1775025604499, 1775025580000, 50, []string{"mid.a", "mid.b"})
	t.NoError(err)

	t.Equal(expect, res)
}

func (t *messagesTest) TestGetMessagesV2() {
	data, err := stabs.ReadFile("stabs/messages/get-messages.json")
	t.NoError(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, pathMessages)

		query := r.URL.Query()
		t.Equal("182182182", query.Get(paramChatID))
		t.Equal("mid.a,mid.b", query.Get(paramMessageIDs))
		t.Equal("1775025604499", query.Get(paramBefore))
		t.Equal("1775025580000", query.Get(paramAfter))
		t.Equal("50", query.Get(paramCount))
		t.Equal("", query.Get(paramFrom))
		t.Equal("", query.Get(paramTo))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Messages.GetMessagesV2(
		context.Background(), 182182182,
		WithMessageIDs([]string{"mid.a", "mid.b"}),
		WithBefore(1775025604499),
		WithAfter(1775025580000),
		WithCount(50),
	)
	t.NoError(err)

	t.Require().Len(res.Messages, 1)
	t.Equal("mid.ffffbdb48e6c3775019d496b34394b84", res.Messages[0].Body.Mid)
}

func (t *messagesTest) TestAnswerOnCallback() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPost)
		t.Equal(r.URL.Path, pathAnswers)

		query := r.URL.Query()
		t.Equal("callback-id", query.Get(paramCallbackID))
		t.Equal("true", query.Get(paramDisableLinkPreview))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Messages.AnswerOnCallback(
		context.Background(), "callback-id", model.CallbackAnswer{},
		WithDisableLinkPreview(),
	)
	t.NoError(err)

	t.True(res.Success)
}
