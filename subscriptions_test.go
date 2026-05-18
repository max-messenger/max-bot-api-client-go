package maxbot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestUpdate(t *testing.T) {
	suite.Run(t, new(subscriptionTest))
}

type subscriptionTest struct {
	suite.Suite
}

func (t *subscriptionTest) SetupTest() {}

func (t *subscriptionTest) TestUpdate() {
	cases := []struct {
		fileName string
		expected model.Update
	}{
		{
			fileName: "stabs/update.bot_added.json",
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
			fileName: "stabs/update.bot_removed.json",
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
			fileName: "stabs/update.bot_started.json",
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
			fileName: "stabs/update.bot_stopped.json",
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
			fileName: "stabs/update.chat_title_changed.json",
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
			fileName: "stabs/update.dialog_cleared.json",
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
			fileName: "stabs/update.dialog_muted.json",
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
			fileName: "stabs/update.dialog_removed.json",
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
			fileName: "stabs/update.dialog_unmuted.json",
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
			fileName: "stabs/update.message_callback.json",
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
			fileName: "stabs/update.message_created.json",
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
			fileName: "stabs/update.message_edited.json",
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
			fileName: "stabs/update.message_removed.json",
			expected: model.Update{
				Timestamp:  1775025604499,
				ChatID:     182182182,
				UserID:     123456789,
				UpdateType: model.UpdateMessageRemoved,
				MessageID:  "mid.000000000adf429c019d47b1ce4600ff",
			},
		},
		{
			fileName: "stabs/update.user_added.json",
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
			fileName: "stabs/update.user_removed.json",
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

	for _, c := range cases {
		t.T().Run(c.fileName, func(tr *testing.T) {
			data, err := stabs.ReadFile(c.fileName)
			t.NoError(err)

			updateList := &updateList{}
			err = json.Unmarshal(data, updateList)
			t.NoError(err)

			t.Require().Len(updateList.Updates, 1)

			t.Equal(c.expected, updateList.Updates[0].FromRaw())
		})
	}
}
