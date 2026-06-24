package maxbot

import "time"

const (
	// Deprecated: not allowed
	Version = "1.2.5"

	SecretHeader = "X-Max-Bot-Api-Secret"

	DefaultAPIURL = "https://platform-api2.max.ru/"

	// Deprecated
	defaultAPIURL   = "https://platform-api.max.ru/"
	defaultTimeout  = 30 * time.Second
	defaultPause    = 1 * time.Second
	maxUpdatesLimit = 50

	// defaultErrorBufferSize — ёмкость канала ошибок по умолчанию. Буфер в 1
	// элемент терял ошибки при всплесках (например, при повторных сбоях long
	// polling): лишние ошибки уходили только в log, минуя GetErrors().
	defaultErrorBufferSize = 16

	maxRetries = 3
)

const (
	notifyExists = "notify/exists"
)

const (
	pathMe            = "me"
	pathChats         = "chats"
	pathAnswers       = "answers"
	pathUpdates       = "updates"
	pathUpload        = "uploads"
	pathMessages      = "messages"
	pathSubscriptions = "subscriptions"

	formatPathChatsID           = "chats/%d"
	formatPathChatsActions      = "chats/%d/actions"
	formatPathChatsMembers      = "chats/%d/members"
	formatPathChatsMembersMe    = "chats/%d/members/me"
	formatPathChatsMembersAdmin = "chats/%d/members/admins"
	formatPathChatsPin          = "chats/%d/pin"
)

const (
	paramURL          = "url"
	paramType         = "type"
	paramTypes        = "types"
	paramMarker       = "marker"
	paramAccessToken  = "access_token"
	paramPhoneNumbers = "phone_numbers"

	paramChatID             = "chat_id"
	paramUserID             = "user_id"
	paramMessageID          = "message_id"
	paramMessageIDs         = "message_ids"
	paramCallbackID         = "callback_id"
	paramDisableLinkPreview = "disable_link_preview"

	paramTo      = "to"
	paramCount   = "count"
	paramFrom    = "from"
	paramLimit   = "limit"
	paramTimeout = "timeout"
)
