package maxbot

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-messenger/max-bot-api-client-go/v3/model"
)

func TestSubscriptions(t *testing.T) {
	suite.Run(t, new(subscriptionsTest))
}

type subscriptionsTest struct {
	suite.Suite
}

func (t *subscriptionsTest) SetupTest() {

}

func (t *subscriptionsTest) TestGetUpdatesWithTypes() {
	data, err := stabs.ReadFile("stabs/update.bot_added.json")
	t.NoError(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, pathUpdates)
		t.Equal("message_created,message_callback", r.URL.Query().Get(paramTypes))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(
		testToken,
		WithBaseURL(srv.URL),
		WithUpdateTypes(model.UpdateMessageCreated, model.UpdateMessageCallback),
	)
	t.NoError(err)

	updates, marker, err := api.Subscriptions.GetUpdates(context.Background(), 0)
	t.NoError(err)

	t.Equal(int64(25493970), marker)
	t.Require().Len(updates, 1)
	t.Equal(model.UpdateBotAdded, updates[0].UpdateType)
	t.Equal(int64(-70801090403050), updates[0].ChatID)
}

func (t *subscriptionsTest) TestGetUpdatesWithoutTypes() {
	data, err := stabs.ReadFile("stabs/update.bot_added.json")
	t.NoError(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal("", r.URL.Query().Get(paramTypes))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	_, _, err = api.Subscriptions.GetUpdates(context.Background(), 0)
	t.NoError(err)
}

func (t *subscriptionsTest) TestSubscribe() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPost)
		t.Equal(r.URL.Path, pathSubscriptions)

		body, err := io.ReadAll(r.Body)
		t.NoError(err)
		t.JSONEq(`{"url":"https://example.com/webhook","secret":"secret","update_types":["message_created"]}`, string(body))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Subscriptions.Subscribe(context.Background(), "https://example.com/webhook", "secret", []string{"message_created"})
	t.NoError(err)

	t.True(res.Success)
}
