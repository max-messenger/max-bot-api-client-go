package model

type Sender struct {
	UserID           int64
	FirstName        string
	LastName         string
	Username         string
	IsBot            bool
	LastActivityTime int64
	Name             string
}

type UserApp struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhotoURL     string `json:"photo_url"`
}

type User struct {
	UserID           int64  `json:"user_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	IsBot            bool   `json:"is_bot"`
	LastActivityTime int64  `json:"last_activity_time"`
	AvatarURL        string `json:"avatar_url"`
	FullAvatarURL    string `json:"full_avatar_url"`
	Name             string `json:"name"`
	Username         string `json:"username"` // need for bot
}

type BotInfo struct {
	UserID           int64        `json:"user_id"`
	FirstName        string       `json:"first_name"`
	Username         string       `json:"username"`
	IsBot            bool         `json:"is_bot"`
	LastActivityTime int64        `json:"last_activity_time"`
	Description      string       `json:"description,omitempty"`
	AvatarURL        string       `json:"avatar_url,omitempty"`
	FullAvatarURL    string       `json:"full_avatar_url,omitempty"`
	IsOfficial       bool         `json:"is_official,omitempty"`
	Commands         []BotCommand `json:"commands,omitempty"`
}

type BotPatch struct {
	FirstName   string       `json:"first_name,omitempty"`
	Description string       `json:"description,omitempty"`
	Commands    []BotCommand `json:"commands,omitempty"`
	Photo       *Payload     `json:"photo,omitempty"`
}

type BotCommand struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
