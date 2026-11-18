package maxbot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
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
