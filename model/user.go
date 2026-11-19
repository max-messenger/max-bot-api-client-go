package model

type Sender struct {
	UserID           int64
	FirstName        string
	LastName         string
	Username         string
	IsBot            bool
	LastActivityTime int64
}

type User struct {
	UserID           int64  `json:"user_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	IsBot            bool   `json:"is_bot"`
	LastActivityTime int64  `json:"last_activity_time"`
	Username         string `json:"username"` // need for bot
}

type BotInfo struct {
	UserID           int64        `json:"user_id"`
	FirstName        string       `json:"first_name"`
	LastName         string       `json:"last_name,omitempty"`
	Username         string       `json:"username"`
	IsBot            bool         `json:"is_bot"`
	LastActivityTime int64        `json:"last_activity_time"`
	Description      string       `json:"description,omitempty"`
	AvatarURL        string       `json:"avatar_url,omitempty"`
	FullAvatarURL    string       `json:"full_avatar_url,omitempty"`
	Commands         []BotCommand `json:"commands,omitempty"`
}

type BotPatchCommands struct {
	Commands []BotCommand `json:"commands"`
}

type BotCommand struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
