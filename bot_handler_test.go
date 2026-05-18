package maxbot

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestBotHandler(t *testing.T) {
	suite.Run(t, new(testBotHandler))
}

type testBotHandler struct {
	suite.Suite
}

func (t *testBotHandler) SetupTest() {}

func (t *testBotHandler) TestHandlerSuccess() {
	cases := []struct {
		fileName string
		expected model.Update
	}{
		{
			fileName: "stabs/webhook.bot_added.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     -70801090403050,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateBotAdded,
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.bot_removed.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     -70801090403050,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateBotRemoved,
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.bot_started.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateBotStarted,
				UserLocale: "ru",
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.bot_stopped.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateBotStopped,
				UserLocale: "ru",
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.chat_title_changed.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     -70801090403050,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateChatTitleChanged,
				ChatProp: &model.ChatProp{
					Title: "Look at me",
				},
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.dialog_cleared.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateDialogCleared,
				UserLocale: "ru",
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.dialog_muted.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateDialogMuted,
				UserLocale: "ru",
				ChatProp: &model.ChatProp{
					MutedUntil: 1775027479470,
				},
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.dialog_removed.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateDialogRemoved,
				UserLocale: "ru",
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.dialog_unmuted.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateDialogUnmuted,
				UserLocale: "ru",
				ChatProp:   &model.ChatProp{},
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.message_callback.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				IsChannel:  false,
				UserLocale: "ru",
				UpdateType: model.UpdateMessageCallback,
				MessageID:  "mid.000000000adf429c019d47d58f2b3d9e",
				Message: &model.MessageUpdate{
					Timestamp: 1775026671403,
					Recipient: model.Recipient{
						UserID:   123456789,
						ChatID:   182182182,
						ChatType: model.ChatTypeDialog,
					},
					Body: model.MessageBody{
						Mid:  "mid.000000000adf429c019d47d58f2b3d9e",
						Seq:  116328147937082782,
						Text: "Hello, John Doe! Your message: empty",
					},
					Sender: model.Sender{
						UserID:           229229229,
						FirstName:        "UniBot",
						IsBot:            true,
						LastActivityTime: 1775026702261,
						Name:             "UniBot",
						Username:         "unit_bot",
					},
				},
				Callback: &model.Callback{
					Timestamp:  1775026702210,
					Payload:    "picture",
					CallbackID: "f9LHodD0cOJf4_DkJGeq8BkDgc5vgSwZocVrn44oirfMzUQ4mv5k_h1-yQvExmZNjV7gVcaO2Z3Gv6LJpZQ-nj_0HTcX7NwSdT4fDtXou9i0A51TjSj9",
				},
			},
		},
		{
			fileName: "stabs/webhook.message_created.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     -70801090403050,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateMessageCreated,
				MessageID:  "mid.ffffbdb48e6c3775019d496b34394b84",
				Message: &model.MessageUpdate{
					Timestamp: 1775053255737,
					Recipient: model.Recipient{
						ChatID:   -70801090403050,
						ChatType: model.ChatTypeChat,
					},
					Body: model.MessageBody{
						Mid:  "mid.ffffbdb48e6c3775019d496b34394b84",
						Seq:  116327994376978687,
						Text: "...",
					},
					Sender: model.Sender{
						UserID:           123456789,
						FirstName:        "John",
						LastName:         "Doe",
						IsBot:            false,
						LastActivityTime: 1775053249000,
						Name:             "John Doe",
					},
					Link: &model.LinkedMessage{
						Type: model.LinkTypeForward,
						Sender: &model.User{
							UserID:           398398398,
							FirstName:        "Tod",
							LastName:         "V",
							IsBot:            false,
							LastActivityTime: 1775755269000,
							Name:             "Tod V",
						},
						ChatID: -695695695695,
						Message: model.MessageBody{
							Mid:  "mid.sha-more",
							Seq:  116327994376978687,
							Text: "Лада седан - баклажан",
						},
					},
				},
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775053249000,
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.message_edited.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateMessageEdited,
				MessageID:  "mid.000000000adf429c019d47b1ce4600ff",
				Message: &model.MessageUpdate{
					Timestamp: 1775025603399,
					Recipient: model.Recipient{
						UserID:   229229229,
						ChatID:   182182182,
						ChatType: model.ChatTypeDialog,
					},
					Body: model.MessageBody{
						Mid:  "mid.000000000adf429c019d47b1ce4600ff",
						Seq:  116327994376978687,
						Text: "hi bot",
					},
					Sender: model.Sender{
						UserID:           123456789,
						FirstName:        "John",
						LastName:         "Doe",
						IsBot:            false,
						LastActivityTime: 1775024330000,
						Name:             "John Doe",
					},
				},
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775024330000,
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.message_removed.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				UpdateType: model.UpdateMessageRemoved,
				MessageID:  "mid.000000000adf429c019d47b1ce4600ff",
			},
		},
		{
			fileName: "stabs/webhook.user_added.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     -70801090403050,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateUserAdded,
				ChatProp: &model.ChatProp{
					InviterID: 123456789,
				},
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
		{
			fileName: "stabs/webhook.user_removed.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     -70801090403050,
				UserID:     123456789,
				IsChannel:  false,
				UpdateType: model.UpdateUserRemoved,
				ChatProp: &model.ChatProp{
					AdminID: 123456789,
				},
				User: &model.User{
					UserID:           123456789,
					FirstName:        "John",
					LastName:         "Doe",
					IsBot:            false,
					LastActivityTime: 1775025580000,
					AvatarURL:        "avatar.png",
					FullAvatarURL:    "avatar.full.png",
					Name:             "John Doe",
				},
			},
		},
	}

	api, err := NewApi(testToken)
	t.NoError(err)

	var data []byte
	for _, c := range cases {
		t.T().Run(c.fileName, func(tr *testing.T) {

			h := func(_ context.Context, upd model.Update) {
				t.Equal(c.expected, upd)
			}

			data, err = stabs.ReadFile(c.fileName)
			t.NoError(err)

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(data))
			req.Header.Set(SecretHeader, testSecret)
			t.NoError(err)

			webhookHandler := api.GetHandler(h, testSecret)
			responseWriter := httptest.NewRecorder()

			webhookHandler.ServeHTTP(responseWriter, req)

			t.Require().Equal(responseWriter.Code, http.StatusOK)
		})
	}
}

type errorReader struct {
	err error
}

func (e errorReader) Read(_ []byte) (int, error) {
	return 0, e.err
}

func (e errorReader) Close() error {
	return nil
}

func (t *testBotHandler) TestHandlerError() {
	cases := []struct {
		name    string
		status  int
		secret  string
		handler UpdateHandler
		expect  string
		method  string
		body    io.Reader
	}{
		{
			name:    "wrong secret",
			status:  http.StatusUnauthorized,
			handler: func(_ context.Context, _ model.Update) {},
			secret:  "wrong",
			expect:  "Unauthorized\n",
			method:  http.MethodPost,
			body:    bytes.NewBufferString(`{}`),
		},
		{
			name:    "empty secret",
			status:  http.StatusUnauthorized,
			handler: func(_ context.Context, _ model.Update) {},
			secret:  "wrong",
			expect:  "Unauthorized\n",
			method:  http.MethodPost,
			body:    bytes.NewBufferString(`{}`),
		},
		{
			name:   "empty handler",
			status: http.StatusInternalServerError,
			expect: "handler is nil\n",
			method: http.MethodPost,
			body:   bytes.NewBufferString(`{}`),
		},
		{
			name:   "empty handler",
			status: http.StatusInternalServerError,
			expect: "handler is nil\n",
			method: http.MethodPost,
			body:   bytes.NewBufferString(`{}`),
		},
		{
			name:    "wrong method",
			status:  http.StatusMethodNotAllowed,
			handler: func(_ context.Context, _ model.Update) {},
			secret:  testSecret,
			expect:  "Method not allowed\n",
			method:  http.MethodGet,
			body:    bytes.NewBufferString(`{}`),
		},
		{
			name:    "without body",
			status:  http.StatusBadRequest,
			handler: func(_ context.Context, _ model.Update) {},
			secret:  testSecret,
			expect:  "Failed to read request body\n",
			method:  http.MethodPost,
			body:    errorReader{err: errors.New("something went wrong")},
		},
		{
			name:    "wrong body",
			status:  http.StatusBadRequest,
			handler: func(_ context.Context, _ model.Update) {},
			secret:  testSecret,
			expect:  "Failed to parse update\n",
			method:  http.MethodPost,
			body:    bytes.NewBufferString(`[]`),
		},
	}

	api, err := NewApi(testToken)
	t.NoError(err)

	for _, c := range cases {
		t.T().Run(c.name, func(tr *testing.T) {
			req, err := http.NewRequest(c.method, "/", c.body)
			req.Header.Set(SecretHeader, testSecret)
			t.NoError(err)

			webhookHandler := api.GetHandler(c.handler, c.secret)
			responseWriter := httptest.NewRecorder()

			webhookHandler.ServeHTTP(responseWriter, req)

			t.Require().Equal(responseWriter.Code, c.status)
			t.Require().Equal(responseWriter.Body.String(), c.expect)
		})
	}
}
