package model

type CommentList struct {
	Messages []Comment `json:"messages"`
}

type Comment struct {
	Recipient CommentRecipient `json:"recipient"`
	Timestamp int64            `json:"timestamp"`
	Body      MessageBody      `json:"body"`
}

type CommentRecipient struct {
	ChatID   int64    `json:"chat_id"`
	ChatType ChatType `json:"chat_type"`
	PostID   string   `json:"post_id"`
}
