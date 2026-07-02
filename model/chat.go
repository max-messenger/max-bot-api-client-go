package model

type ChatProp struct {
	Title      string
	MutedUntil int64
	AdminID    int64
	InviterID  int64
}

type Chat struct {
	ChatID            int64            `json:"chat_id"`
	Type              ChatType         `json:"type"`
	Status            ChatStatus       `json:"status"`
	Title             string           `json:"title"`
	Icon              Image            `json:"icon"`
	LastEventTime     int64            `json:"last_event_time"`
	ParticipantsCount int32            `json:"participants_count"`
	OwnerID           int64            `json:"owner_id,omitempty"`
	Participants      map[string]int64 `json:"participants,omitempty"`
	IsPublic          bool             `json:"is_public"`
	Link              string           `json:"link,omitempty"`
	Description       string           `json:"description,omitempty"`
	DialogWithUser    *User            `json:"dialog_with_user,omitempty"`
	MessagesCount     int              `json:"messages_count,omitempty"`
	ChatMessageID     string           `json:"chat_message_id,omitempty"`
	PinnedMessage     Message          `json:"pinned_message"`
}

type ChatList struct {
	Chats  []Chat `json:"chats"`
	Marker *int64 `json:"marker"`
}

type ChatPatch struct {
	Icon   *Payload `json:"icon,omitempty"`
	Title  string   `json:"title,omitempty"`
	Pin    string   `json:"pin,omitempty"`
	Notify *bool    `json:"notify,omitempty"`
}

type ChatMember struct {
	LastAccessTime   int64                 `json:"last_access_time"`
	IsOwner          bool                  `json:"is_owner"`
	IsAdmin          bool                  `json:"is_admin"`
	JoinTime         int64                 `json:"join_time"`
	Permissions      []ChatAdminPermission `json:"permissions"`
	Alias            string                `json:"alias"`
	UserID           int64                 `json:"user_id"`
	FirstName        string                `json:"first_name"`
	LastName         string                `json:"last_name"`
	IsBot            bool                  `json:"is_bot"`
	LastActivityTime int64                 `json:"last_activity_time"`
	AvatarURL        string                `json:"avatar_url"`
	FullAvatarURL    string                `json:"full_avatar_url"`
	Name             string                `json:"name"`
	Username         string                `json:"username"` // need for bot
	Description      string                `json:"description"`
}

type ChatMembersList struct {
	Members []ChatMember `json:"members"`
	Marker  int64        `json:"marker,omitempty"`
}

type Image struct {
	URL string `json:"url"`
}

type ChatAdmin struct {
	UserID      int64                 `json:"user_id"`
	Permissions []ChatAdminPermission `json:"permissions"`
	Alias       string                `json:"alias,omitempty"`
}

type ChatAdminsList struct {
	Admins []ChatAdmin `json:"admins"`
}
