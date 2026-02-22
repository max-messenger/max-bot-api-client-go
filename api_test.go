package maxbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func TestBytesToProperUpdate(t *testing.T) {
	api, err := New("test")
	require.NoError(t, err)

	tests := []struct {
		name       string
		data       func(t *testing.T) []byte
		wantType   reflect.Type
		err        error
		wantUpdate schemes.UpdateInterface
	}{
		{
			name: "message created",
			data: func(t *testing.T) []byte {
				return mustMarshal(t, &schemes.MessageCreatedUpdate{
					Update: schemes.Update{UpdateType: schemes.TypeMessageCreated, Timestamp: 1234567890},
					Message: schemes.Message{
						Sender:    schemes.User{UserId: 100},
						Recipient: schemes.Recipient{ChatId: 1, ChatType: schemes.ChatType("dialog"), UserId: 200},
						Timestamp: 1234567890,
						Body:      schemes.MessageBody{Mid: "mid1", Seq: 1, Text: "hello"},
					},
				})
			},
			wantType: reflect.TypeOf(&schemes.MessageCreatedUpdate{}),
			wantUpdate: &schemes.MessageCreatedUpdate{
				Update: schemes.Update{UpdateType: schemes.TypeMessageCreated, Timestamp: 1234567890},
				Message: schemes.Message{
					Sender:    schemes.User{UserId: 100},
					Recipient: schemes.Recipient{ChatId: 1, ChatType: schemes.ChatType("dialog"), UserId: 200},
					Timestamp: 1234567890,
					Body:      schemes.MessageBody{Mid: "mid1", Seq: 1, Text: "hello"},
				},
			},
		},
		{
			name: "unknown type",
			data: func(t *testing.T) []byte { return mustMarshal(t, schemes.Update{UpdateType: "unknown"}) },
			err:  fmt.Errorf("unknown update type: unknown"),
		},
		{
			name: "bot added",
			data: func(t *testing.T) []byte {
				return mustMarshal(t, &schemes.BotAddedToChatUpdate{
					Update: schemes.Update{UpdateType: schemes.TypeBotAdded, Timestamp: 1234567890},
					ChatId: 1,
					User:   schemes.User{UserId: 100},
				})
			},
			wantType: reflect.TypeOf(&schemes.BotAddedToChatUpdate{}),
			wantUpdate: &schemes.BotAddedToChatUpdate{
				Update: schemes.Update{UpdateType: schemes.TypeBotAdded, Timestamp: 1234567890},
				ChatId: 1,
				User:   schemes.User{UserId: 100},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.data(t)
			got, err := api.bytesToProperUpdate(data)
			if tt.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantType, reflect.TypeOf(got))
			require.Equal(t, tt.wantUpdate, got)
		})
	}
}

func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()

	data, err := jsoniter.Marshal(v)
	require.NoError(t, err)

	return data
}

func TestBytesToProperAttachment(t *testing.T) {
	api, err := New("test")
	require.NoError(t, err)

	tests := []struct {
		name     string
		attach   schemes.AttachmentInterface
		wantType reflect.Type
		wantErr  bool
	}{
		{
			name: "image",
			attach: &schemes.PhotoAttachment{
				Attachment: schemes.Attachment{Type: schemes.AttachmentImage},
				Payload:    schemes.PhotoAttachmentPayload{Url: "http://example.com/img.jpg"},
			},
			wantType: reflect.TypeOf(&schemes.PhotoAttachment{}),
		},
		{
			name:     "unknown",
			attach:   &schemes.Attachment{Type: "unknown"},
			wantType: reflect.TypeOf(&schemes.Attachment{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := jsoniter.Marshal(tt.attach)
			require.NoError(t, err)

			got, err := api.bytesToProperAttachment(data)
			require.Equal(t, tt.wantType, reflect.TypeOf(got))
		})
	}
}

func TestGetUpdates(t *testing.T) {
	wantUpdate := &schemes.MessageCreatedUpdate{
		Update: schemes.Update{UpdateType: schemes.TypeMessageCreated, Timestamp: 1234567890},
		Message: schemes.Message{
			Sender:    schemes.User{UserId: 100},
			Recipient: schemes.Recipient{ChatId: 1, ChatType: schemes.ChatType("dialog"), UserId: 200},
			Timestamp: 1234567890,
			Body:      schemes.MessageBody{Mid: "mid1", Seq: 1, Text: "test message"},
		},
	}
	updateJSON, err := jsoniter.Marshal(wantUpdate)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/updates", r.URL.Path)

		markerStr := r.URL.Query().Get(paramMarker)
		marker, _ := strconv.ParseInt(markerStr, 10, 64)

		var updateList schemes.UpdateList
		if marker == 0 {
			updateList = schemes.UpdateList{
				Updates: []json.RawMessage{updateJSON},
				Marker:  new(int64),
			}
			*updateList.Marker = 1
		} else {
			updateList = schemes.UpdateList{
				Updates: []json.RawMessage{},
				Marker:  new(int64),
			}
			*updateList.Marker = marker
		}

		_ = jsoniter.NewEncoder(w).Encode(updateList)
	}))
	defer server.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api, err := New("token", WithBaseURL(server.URL))
	require.NoError(t, err)

	api.pause = 10 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch := api.GetUpdates(ctx)

	select {
	case got := <-ch:
		require.Equal(t, wantUpdate, got)
	case <-ctx.Done():
		t.Error("no update received in time")
	}
}

func TestGetHandler(t *testing.T) {
	api, err := New("test")
	require.NoError(t, err)

	ch := make(chan schemes.UpdateInterface, 1)
	handler := api.GetHandler(ch)

	wantUpdate := &schemes.MessageCreatedUpdate{
		Update: schemes.Update{UpdateType: schemes.TypeMessageCreated, Timestamp: 1234567890},
		Message: schemes.Message{
			Sender:    schemes.User{UserId: 100},
			Recipient: schemes.Recipient{ChatId: 1, ChatType: schemes.ChatType("dialog"), UserId: 200},
			Timestamp: 1234567890,
			Body:      schemes.MessageBody{Mid: "mid1", Seq: 1, Text: "test message"},
		},
	}
	updateJSON, err := jsoniter.Marshal(wantUpdate)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	select {
	case got := <-ch:
		require.Equal(t, wantUpdate, got)
	case <-time.After(time.Second):
		t.Error("no update received")
	}
}
