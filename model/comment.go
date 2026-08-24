package model

type CommentList struct {
	Messages []Comment `json:"messages"`
}

type Comment struct {
	Recipient Recipient   `json:"recipient"`
	Timestamp int64       `json:"timestamp"`
	Body      MessageBody `json:"body"`
}
