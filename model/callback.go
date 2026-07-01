package model

type Callback struct {
	Timestamp  int64  `json:"timestamp"`
	CallbackID string `json:"callback_id"`
	Payload    string `json:"payload"`
	User       User   `json:"user"`
}

type CallbackAnswer struct {
	Message      *NewMessageBody `json:"message,omitempty"`
	Notification *string         `json:"notification,omitempty"`
}
