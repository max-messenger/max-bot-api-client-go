package maxbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage()

	assert.NotNil(t, msg)
	assert.Equal(t, int64(0), msg.userID)
	assert.Equal(t, int64(0), msg.chatID)
	assert.Equal(t, model.NewMessageBody{}, msg.message)
	assert.False(t, msg.disableLinkPreview)
}

func TestMessage_SetUser(t *testing.T) {
	msg := NewMessage()
	userID := int64(12345)

	result := msg.SetUser(userID)

	assert.Equal(t, msg, result) // check fluent interface
	assert.Equal(t, userID, msg.userID)
}

func TestMessage_SetChat(t *testing.T) {
	msg := NewMessage()
	chatID := int64(67890)

	result := msg.SetChat(chatID)

	assert.Equal(t, msg, result)
	assert.Equal(t, chatID, msg.chatID)
}

func TestMessage_SetDisableLinkPreview(t *testing.T) {
	tests := []struct {
		name string
		mode bool
	}{
		{"enable preview", false},
		{"disable preview", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage()
			result := msg.SetDisableLinkPreview(tt.mode)

			assert.Equal(t, msg, result)
			assert.Equal(t, tt.mode, msg.disableLinkPreview)
		})
	}
}

func TestMessage_SetText(t *testing.T) {
	msg := NewMessage()
	text := "Hello, World!"

	result := msg.SetText(text)

	assert.Equal(t, msg, result)
	assert.Equal(t, text, msg.message.Text)
}

func TestMessage_AddImageUrl(t *testing.T) {
	msg := NewMessage()
	photoURL := "https://example.com/image.jpg"

	result := msg.AddImageUrl(photoURL)

	assert.Equal(t, msg, result)
	require.Len(t, msg.message.Attachments, 1)

	attach := msg.message.Attachments[0]
	assert.Equal(t, model.AttachImage, attach.Type)
	assert.Equal(t, photoURL, attach.Payload.URL)
}

func TestMessage_SetFormat(t *testing.T) {
	msg := NewMessage()
	format := model.FormatMarkdown

	result := msg.SetFormat(format)

	assert.Equal(t, msg, result)
	assert.Equal(t, format, msg.message.Format)
}

func TestMessage_AddSticker(t *testing.T) {
	msg := NewMessage()
	stickerCode := "sticker123"

	result := msg.AddSticker(stickerCode)

	assert.Equal(t, msg, result)
	require.Len(t, msg.message.Attachments, 1)

	attach := msg.message.Attachments[0]
	assert.Equal(t, model.AttachSticker, attach.Type)
	assert.Equal(t, stickerCode, attach.Payload.Code)
}

func TestMessage_AddContact(t *testing.T) {
	msg := NewMessage()
	contactID := int64(99999)

	result := msg.AddContact(contactID)

	assert.Equal(t, msg, result)
	require.Len(t, msg.message.Attachments, 1)

	attach := msg.message.Attachments[0]
	assert.Equal(t, model.AttachContact, attach.Type)
	assert.Equal(t, contactID, attach.Payload.ContactID)
}

func TestMessage_AddLocation(t *testing.T) {
	msg := NewMessage()
	lat, lon := 55.751244, 37.618423

	result := msg.AddLocation(lat, lon)

	assert.Equal(t, msg, result)
	require.Len(t, msg.message.Attachments, 1)

	attach := msg.message.Attachments[0]
	assert.Equal(t, model.AttachLocation, attach.Type)
	assert.Equal(t, lat, attach.Latitude)
	assert.Equal(t, lon, attach.Longitude)
}

func TestMessage_AddShare(t *testing.T) {
	msg := NewMessage()
	link := "https://example.com/shared"

	result := msg.AddShare(link)

	assert.Equal(t, msg, result)
	require.Len(t, msg.message.Attachments, 1)

	attach := msg.message.Attachments[0]
	assert.Equal(t, model.AttachShare, attach.Type)
	assert.Equal(t, link, attach.Payload.URL)
}

func TestMessage_WithoutNotify(t *testing.T) {
	msg := NewMessage()

	result := msg.WithoutNotify()

	assert.Equal(t, msg, result)
	require.NotNil(t, msg.message.Notify)
	assert.False(t, *msg.message.Notify)
}

func TestMessage_SetReply(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		replyID string
	}{
		{
			name:    "simple reply",
			text:    "Reply message",
			replyID: "msg_123",
		},
		{
			name:    "empty reply",
			text:    "",
			replyID: "msg_456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage()

			result := msg.SetReply(tt.text, tt.replyID)

			assert.Equal(t, msg, result)
			assert.Equal(t, tt.text, msg.message.Text)
			require.NotNil(t, msg.message.Link)
			assert.Equal(t, model.LinkTypeReply, msg.message.Link.Type)
			assert.Equal(t, tt.replyID, msg.message.Link.Mid)
		})
	}
}

func TestMessage_AddKeyboard(t *testing.T) {
	msg := NewMessage()
	keyboard := &model.Keyboard{} // adjust based on your actual Keyboard implementation
	// Mock keyboard.Build() if needed
	// For this test, assume keyboard.Build() returns some buttons

	result := msg.AddKeyboard(keyboard)

	assert.Equal(t, msg, result)
	require.Len(t, msg.message.Attachments, 1)

	attach := msg.message.Attachments[0]
	assert.Equal(t, model.AttachInlineKeyboard, attach.Type)
	assert.NotNil(t, attach.Payload.Buttons)
}

func TestMessage_AddAttachByToken(t *testing.T) {
	tests := []struct {
		name       string
		fileToken  string
		attachType model.AttachmentType
	}{
		{
			name:       "image token",
			fileToken:  "token_image_123",
			attachType: model.AttachImage,
		},
		{
			name:       "file token",
			fileToken:  "token_file_456",
			attachType: model.AttachFile,
		},
		{
			name:       "video token",
			fileToken:  "token_video_789",
			attachType: model.AttachVideo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage()

			result := msg.AddAttachByToken(tt.fileToken, tt.attachType)

			assert.Equal(t, msg, result)
			require.Len(t, msg.message.Attachments, 1)

			attach := msg.message.Attachments[0]
			assert.Equal(t, tt.attachType, attach.Type)
			assert.Equal(t, tt.fileToken, attach.Payload.Token)
		})
	}
}

// Integration tests for multiple attachments
func TestMessage_MultipleAttachments(t *testing.T) {
	msg := NewMessage()

	msg.AddImageUrl("img1.jpg").
		AddImageUrl("img2.jpg").
		AddShare("https://example.com").
		AddLocation(55.75, 37.62)

	assert.Len(t, msg.message.Attachments, 4)

	assert.Equal(t, model.AttachImage, msg.message.Attachments[0].Type)
	assert.Equal(t, model.AttachImage, msg.message.Attachments[1].Type)
	assert.Equal(t, model.AttachShare, msg.message.Attachments[2].Type)
	assert.Equal(t, model.AttachLocation, msg.message.Attachments[3].Type)
}

// Test chaining all methods
func TestMessage_FluentChaining(t *testing.T) {
	msg := NewMessage()
	keyboard := &model.Keyboard{}

	result := msg.
		SetUser(123).
		SetChat(456).
		SetDisableLinkPreview(true).
		SetText("Hello").
		AddImageUrl("img.jpg").
		SetFormat(model.FormatMarkdown).
		AddSticker("sticker1").
		AddContact(789).
		AddLocation(55.75, 37.62).
		AddShare("https://example.com").
		WithoutNotify().
		AddKeyboard(keyboard).
		AddAttachByToken("token123", model.AttachFile)

	assert.Equal(t, msg, result)
	assert.Equal(t, int64(123), msg.userID)
	assert.Equal(t, int64(456), msg.chatID)
	assert.True(t, msg.disableLinkPreview)
	assert.Equal(t, "Hello", msg.message.Text)
	assert.Equal(t, model.FormatMarkdown, msg.message.Format)
	assert.NotNil(t, msg.message.Notify)
	assert.False(t, *msg.message.Notify)
	assert.Len(t, msg.message.Attachments, 7) // image, sticker, contact, location, share, keyboard
}

// Test edge cases
func TestMessage_EdgeCases(t *testing.T) {
	t.Run("empty text", func(t *testing.T) {
		msg := NewMessage()
		msg.SetText("")
		assert.Equal(t, "", msg.message.Text)
	})

	t.Run("empty sticker code", func(t *testing.T) {
		msg := NewMessage()
		msg.AddSticker("")
		assert.Len(t, msg.message.Attachments, 1)
		assert.Equal(t, "", msg.message.Attachments[0].Payload.Code)
	})

	t.Run("zero coordinates", func(t *testing.T) {
		msg := NewMessage()
		msg.AddLocation(0, 0)
		assert.Len(t, msg.message.Attachments, 1)
		assert.Equal(t, 0.0, msg.message.Attachments[0].Latitude)
		assert.Equal(t, 0.0, msg.message.Attachments[0].Longitude)
	})

	t.Run("empty reply id", func(t *testing.T) {
		msg := NewMessage()
		msg.SetReply("test", "")
		assert.Equal(t, "", msg.message.Link.Mid)
	})
}
