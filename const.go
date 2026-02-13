package maxbot

import "time"

const (
	version = "1.2.5"

	defaultAPIURL   = "https://platform-api.max.ru/"
	defaultTimeout  = 30 * time.Second
	defaultPause    = 1 * time.Second
	maxUpdatesLimit = 50

	maxRetries = 3
)

const (
	notifyExists = "notify/exists"
	envToken     = "TOKEN"
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
)

const (
	paramVersion      = "v"
	paramURL          = "url"
	paramType         = "type"
	paramTypes        = "types"
	paramMarker       = "marker"
	paramAccessToken  = "access_token"
	paramPhoneNumbers = "phone_numbers"

	paramChatID     = "chat_id"
	paramUserID     = "user_id"
	paramMessageID  = "message_id"
	paramMessageIDs = "message_ids"
	paramCallbackID = "callback_id"

	paramTo      = "to"
	paramCount   = "count"
	paramFrom    = "from"
	paramLimit   = "limit"
	paramTimeout = "timeout"
)
