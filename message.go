package maxbot

import "github.com/max-messenger/max-bot-api-client-go/v2/model"

type Message struct {
	userID             int64
	chatID             int64
	disableLinkPreview bool
	message            model.NewMessageBody
}

func NewMessage() *Message {
	return &Message{userID: 0, chatID: 0, message: model.NewMessageBody{}}
}

func (m *Message) SetUser(userID int64) *Message {
	m.userID = userID

	return m
}

func (m *Message) SetChat(chatID int64) *Message {
	m.chatID = chatID

	return m
}

func (m *Message) SetDisableLinkPreview(mode bool) *Message {
	m.disableLinkPreview = mode

	return m
}

func (m *Message) SetText(text string) *Message {
	m.message.Text = text

	return m
}

func (m *Message) AddImageUrl(photoUrl string) *Message {
	attach := model.Attachment{
		Type: model.AttachImage,
		Payload: model.Payload{
			URL: photoUrl,
		},
	}
	m.message.Attachments = append(m.message.Attachments, attach)

	return m
}

func (m *Message) SetFormat(format model.TextFormat) *Message {
	m.message.Format = format

	return m
}

// AddSticker добавляет стикер. Сообщение не должно содержать текста и других вложений.
func (m *Message) AddSticker(stickerCode string) *Message {
	attach := model.Attachment{
		Type: model.AttachSticker,
		Payload: model.Payload{
			Code: stickerCode,
		},
	}
	m.message.Attachments = append(m.message.Attachments, attach)

	return m
}

// AddContact добавляет контакт. Сообщение не должно содержать текста и других вложений.
func (m *Message) AddContact(userID int64) *Message {
	attach := model.Attachment{
		Type: model.AttachContact,
		Payload: model.Payload{
			ContactID: userID,
		},
	}
	m.message.Attachments = append(m.message.Attachments, attach)

	return m
}

func (m *Message) AddLocation(lat, lot float64) *Message {
	attach := model.Attachment{
		Type:      model.AttachLocation,
		Latitude:  lat,
		Longitude: lot,
	}
	m.message.Attachments = append(m.message.Attachments, attach)

	return m
}

func (m *Message) AddShare(link string) *Message {
	attach := model.Attachment{
		Type: model.AttachShare,
		Payload: model.Payload{
			URL: link,
		},
	}
	m.message.Attachments = append(m.message.Attachments, attach)

	return m
}

// WithoutNotify Если false, участники чата не будут уведомлены (по умолчанию true).
func (m *Message) WithoutNotify() *Message {
	notify := false
	m.message.Notify = &notify

	return m
}

func (m *Message) SetReply(text, id string) *Message {
	m.message.Text = text
	m.message.Link = &model.NewMessageLink{Type: model.LinkTypeReply, Mid: id}

	return m
}

func (m *Message) AddKeyboard(keyboard *model.Keyboard) *Message {
	m.message.Attachments = append(m.message.Attachments, keyboard.Build())

	return m
}

func (m *Message) AddAttachByToken(fileToken string, at model.AttachmentType) *Message {
	attach := model.Attachment{
		Type: at,
		Payload: model.Payload{
			Token: fileToken,
		},
	}
	m.message.Attachments = append(m.message.Attachments, attach)

	return m
}

func (m *Message) MessageBody() model.NewMessageBody {
	return m.message
}
