package model

type Recipient struct {
	ChatID   int64    `json:"chat_id"`
	ChatType ChatType `json:"chat_type"`
	UserID   int64    `json:"user_id"`
}

type Message struct {
	Sender    User           `json:"sender"`
	Recipient Recipient      `json:"recipient"`
	Timestamp int64          `json:"timestamp"`
	Link      *LinkedMessage `json:"link,omitempty"`
	Body      MessageBody    `json:"body"`
	Stat      MessageStat    `json:"stat"`
	URL       string         `json:"url"`
}

type MessageUpdate struct {
	Timestamp int64
	Recipient Recipient
	Sender    Sender
	Body      MessageBody
	Link      *LinkedMessage
}

type MessageStat struct {
	Views int `json:"views"`
}

type MessageBody struct {
	Mid         string       `json:"mid"`
	Seq         int64        `json:"seq"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type MessageList struct {
	Messages []Message `json:"messages"`
}

type NewMessageBody struct {
	Text        string          `json:"text,omitempty"`
	Attachments []Attachment    `json:"attachments"`
	Link        *NewMessageLink `json:"link,omitempty"`
	Notify      *bool           `json:"notify,omitempty"`
	Format      TextFormat      `json:"format,omitempty"`
}

type NewMessageLink struct {
	Type MessageLinkType `json:"type"`
	Mid  string          `json:"mid"`
}

type LinkedMessage struct {
	Type    MessageLinkType `json:"type"`
	Sender  *User           `json:"sender,omitempty"`
	ChatID  int64           `json:"chat_id,omitempty"`
	Message MessageBody     `json:"message"`
}

type SendMessageResult struct {
	Message Message `json:"message"`
}

type PinMessageBody struct {
	MessageID string `json:"message_id"`
	Notify    *bool  `json:"notify,omitempty"`
}

type GetPinnedMessageResult struct {
	Message Message `json:"message"`
}
