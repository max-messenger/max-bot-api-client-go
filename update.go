package maxbot

import "github.com/max-messenger/max-bot-api-client-go/v2/model"

type updateRaw struct {
	Timestamp  int64            `json:"timestamp"`
	ChatID     int64            `json:"chat_id"`
	UserID     int64            `json:"user_id"`
	UserLocale string           `json:"user_locale"`
	IsChannel  bool             `json:"is_channel"`
	Title      string           `json:"title"`
	MutedUntil int64            `json:"muted_until"`
	InviterID  int64            `json:"inviter_id"`
	AdminID    int64            `json:"admin_id"`
	MessageID  string           `json:"message_id"`
	UpdateType model.UpdateType `json:"update_type"`
	User       model.User       `json:"user"`
	Message    model.Message    `json:"message"`
	Callback   model.Callback   `json:"callback"`
}

type updateList struct {
	Updates []updateRaw `json:"updates"`
	Marker  int64       `json:"marker"`
}

func (u updateRaw) FromRaw() model.Update {
	update := model.Update{
		Timestamp:  u.Timestamp,
		UpdateType: u.UpdateType,
		UserLocale: u.UserLocale,
		ChatID:     u.ChatID,
		UserID:     u.User.UserID,
	}

	switch u.UpdateType {
	case model.UpdateMessageCallback:
		update.ChatID = u.Message.Recipient.ChatID
		update.UserID = u.Message.Recipient.UserID
		update.MessageID = u.Message.Body.Mid
		update.Message = &model.MessageUpdate{
			Timestamp: u.Message.Timestamp,
			Recipient: model.Recipient{
				ChatID:   u.Message.Recipient.ChatID,
				UserID:   u.Message.Recipient.UserID,
				ChatType: u.Message.Recipient.ChatType,
			},
			Sender: u.senderFromRaw(),
			Body: model.MessageBody{
				Mid:         u.Message.Body.Mid,
				Seq:         u.Message.Body.Seq,
				Text:        u.Message.Body.Text,
				Attachments: u.Message.Body.Attachments,
			},
		}
		update.Callback = &u.Callback

	case model.UpdateMessageCreated, model.UpdateMessageEdited:
		update.ChatID = u.Message.Recipient.ChatID
		update.UserID = u.Message.Sender.UserID
		update.MessageID = u.Message.Body.Mid
		update.User = &u.Message.Sender
		update.Message = &model.MessageUpdate{
			Timestamp: u.Message.Timestamp,
			Recipient: model.Recipient{
				ChatID:   u.Message.Recipient.ChatID,
				UserID:   u.Message.Recipient.UserID,
				ChatType: u.Message.Recipient.ChatType,
			},
			Sender: u.senderFromRaw(),
			Body: model.MessageBody{
				Mid:         u.Message.Body.Mid,
				Seq:         u.Message.Body.Seq,
				Text:        u.Message.Body.Text,
				Attachments: u.Message.Body.Attachments,
			},
			Link: u.Message.Link,
		}
	case model.UpdateDialogCleared:
		update.User = &u.User
	case model.UpdateUserRemoved:
		update.ChatProp = &model.ChatProp{
			AdminID: u.AdminID,
		}
		update.User = &u.User
	case model.UpdateUserAdded:
		update.ChatProp = &model.ChatProp{
			InviterID: u.InviterID,
		}
		update.User = &u.User
	case model.UpdateDialogRemoved:
		update.User = &u.User
	case model.UpdateDialogMuted, model.UpdateDialogUnmuted:
		update.ChatProp = &model.ChatProp{
			MutedUntil: u.MutedUntil,
		}
		update.User = &u.User
	case model.UpdateChatTitleChanged:
		update.ChatProp = &model.ChatProp{
			Title: u.Title,
		}
		update.User = &u.User
	case model.UpdateBotAdded, model.UpdateBotRemoved, model.UpdateBotStarted, model.UpdateBotStopped:
		update.IsChannel = u.IsChannel
		update.User = &u.User
	case model.UpdateMessageRemoved:
		update.MessageID = u.MessageID
		update.UserID = u.UserID
	}

	return update
}

func (u updateRaw) senderFromRaw() model.Sender {
	return model.Sender{
		UserID:           u.Message.Sender.UserID,
		FirstName:        u.Message.Sender.FirstName,
		LastName:         u.Message.Sender.LastName,
		IsBot:            u.Message.Sender.IsBot,
		LastActivityTime: u.Message.Sender.LastActivityTime,
		Name:             u.Message.Sender.Name,
		Username:         u.Message.Sender.Username,
	}
}
