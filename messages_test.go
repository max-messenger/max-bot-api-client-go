package maxbot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetMessagesPassesFromToAsIs проверяет, что GetMessages передаёт параметры
// from/to ровно как заданы вызывающим кодом, без свапа (issue #136). По
// контракту API from — верхняя граница (t <= from), to — нижняя (t >= to),
// поэтому корректный диапазон имеет from > to и не должен инвертироваться.
func TestGetMessagesPassesFromToAsIs(t *testing.T) {
	var gotFrom, gotTo string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotFrom = r.URL.Query().Get("from")
		gotTo = r.URL.Query().Get("to")
		_, _ = w.Write([]byte(`{"messages":[]}`))
	}))
	defer server.Close()

	api, err := New("token", WithBaseURL(server.URL))
	require.NoError(t, err)

	_, err = api.Messages.GetMessages(context.Background(), 1, nil, 2000, 1000, 50)
	require.NoError(t, err)

	require.Equal(t, "2000", gotFrom, "from must be passed unchanged")
	require.Equal(t, "1000", gotTo, "to must be passed unchanged")
}

// TestGetMessagesOmitsZeroBounds проверяет, что нулевые from/to не отправляются
// (открытая граница диапазона).
func TestGetMessagesOmitsZeroBounds(t *testing.T) {
	var hasFrom, hasTo bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasFrom = r.URL.Query()["from"]
		_, hasTo = r.URL.Query()["to"]
		_, _ = w.Write([]byte(`{"messages":[]}`))
	}))
	defer server.Close()

	api, err := New("token", WithBaseURL(server.URL))
	require.NoError(t, err)

	_, err = api.Messages.GetMessages(context.Background(), 1, nil, 0, 1000, 50)
	require.NoError(t, err)

	require.False(t, hasFrom, "from=0 must be omitted")
	require.True(t, hasTo, "to must be present")
}
