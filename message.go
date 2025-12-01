package maxbot

import "github.com/max-messenger/max-bot-api-client-go/schemes"

type Message struct {
	userID  int64
	chatID  int64
	reset   bool
	message *schemes.NewMessageBody
}

func NewMessage() *Message {
	return &Message{userID: 0, chatID: 0, message: &schemes.NewMessageBody{Attachments: []interface{}{}}}
}

func (m *Message) SetUser(userID int64) *Message {
	m.userID = userID
	return m
}

func (m *Message) SetChat(chatID int64) *Message {
	m.chatID = chatID
	return m
}

func (m *Message) SetReset(reset bool) *Message {
	m.reset = reset
	return m
}

func (m *Message) SetText(text string) *Message {
	m.message.Text = text
	return m
}

func (m *Message) SetFormat(format string) *Message {
	m.message.Format = format
	return m
}
func (m *Message) SetNotify(notify bool) *Message {
	m.message.Notify = notify
	return m
}

func (m *Message) SetReply(text, id string) *Message {
	m.message.Text = text
	m.message.Link = &schemes.NewMessageLink{Type: schemes.REPLY, Mid: id}
	return m
}

func (m *Message) Reply(text string, reply schemes.Message) *Message {
	m.message.Text = text
	if reply.Recipient.UserId != 0 {
		m.userID = reply.Recipient.UserId
	}
	if reply.Recipient.ChatId != 0 {
		m.chatID = reply.Recipient.ChatId
	}
	m.message.Link = &schemes.NewMessageLink{Type: schemes.REPLY, Mid: reply.Body.Mid}
	return m
}

func (m *Message) AddMarkUp(user int64, from int, len int) *Message {
	m.message.Markups = append(m.message.Markups, schemes.MarkUp{UserId: user, From: from, Length: len, Type: schemes.MarkupUser})
	return m
}

func (m *Message) AddKeyboard(keyboard *Keyboard) *Message {
	m.message.Attachments = append(m.message.Attachments, schemes.NewInlineKeyboardAttachmentRequest(keyboard.Build()))
	return m
}

func (m *Message) AddPhoto(photo *schemes.PhotoTokens) *Message {
	m.message.Attachments = append(m.message.Attachments, schemes.NewPhotoAttachmentRequest(schemes.PhotoAttachmentRequestPayload{
		Photos: photo.Photos,
	}))
	return m
}

func (m *Message) AddAudio(audio *schemes.UploadedInfo) *Message {
	m.message.Attachments = append(m.message.Attachments, schemes.NewAudioAttachmentRequest(*audio))
	return m
}

func (m *Message) AddVideo(video *schemes.UploadedInfo) *Message {
	m.message.Attachments = append(m.message.Attachments, schemes.NewVideoAttachmentRequest(*video))
	return m
}

func (m *Message) AddFile(file *schemes.UploadedInfo) *Message {
	m.message.Attachments = append(m.message.Attachments, schemes.NewFileAttachmentRequest(*file))
	return m
}

func (m *Message) AddLocation(lat float64, lon float64) *Message {
	m.message.Attachments = append(m.message.Attachments, schemes.NewLocationAttachmentRequest(lat, lon))
	return m
}

func (m *Message) AddContact(name string, contactID int64, vcfInfo string, vcfPhone string) *Message {
	m.message.Attachments = append(m.message.Attachments, schemes.NewContactAttachmentRequest(schemes.ContactAttachmentRequestPayload{
		Name:      name,
		ContactId: contactID,
		VcfInfo:   vcfInfo,
		VcfPhone:  vcfPhone,
	}))
	return m
}

func (m *Message) AddSticker(code string) *Message {
	m.message.Attachments = append(m.message.Attachments, schemes.NewStickerAttachmentRequest(schemes.StickerAttachmentRequestPayload{
		Code: code,
	}))
	return m
}
