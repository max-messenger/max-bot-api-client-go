package model

type CommentList struct {
	Messages []Comment `json:"messages"`
}

type SendCommentResult struct {
	Message Comment `json:"message"`
}

type Comment struct {
	Recipient Recipient             `json:"recipient"`
	Timestamp int64                 `json:"timestamp"`
	Body      CommentBody           `json:"body"`
	Sender    *User                 `json:"sender,omitempty"`
	Stat      *MessageStat          `json:"stat,omitempty"`
	Link      *CommentLinkedMessage `json:"link,omitempty"`
}

type CommentBody struct {
	CommentID string   `json:"mid"`
	Seq       int64    `json:"seq"`
	Text      string   `json:"text"`
	MarkUp    []Markup `json:"markup,omitempty"`
}

type CommentLinkedMessage struct {
	Type    MessageLinkType `json:"type"`
	Sender  *User           `json:"sender,omitempty"`
	ChatID  int64           `json:"chat_id,omitempty"`
	Message CommentBody     `json:"message"`
}
