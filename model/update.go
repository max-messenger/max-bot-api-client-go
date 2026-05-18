package model

import "time"

type Update struct {
	Timestamp  int64
	ChatID     int64
	UserID     int64
	UserLocale string
	IsChannel  bool
	UpdateType UpdateType
	MessageID  string
	User       *User
	Callback   *Callback
	ChatProp   *ChatProp
	Message    *MessageUpdate
}

func (u Update) GetTimestampTime() time.Time {
	return time.UnixMilli(u.Timestamp)
}

func (u Update) GetUserActivityTime() time.Time {
	if u.User != nil && u.User.LastActivityTime > 0 {
		return time.UnixMilli(u.User.LastActivityTime)
	}

	return time.Time{}
}

func (u Update) GetMessage() MessageUpdate {
	if u.Message != nil {
		return *u.Message
	}

	return MessageUpdate{}
}

func (u Update) GetUser() User {
	if u.User != nil {
		return *u.User
	}

	return User{}
}
