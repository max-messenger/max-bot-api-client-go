package maxbot

import "time"

const (
	defaultScheme = "https"

	// Deprecated: not allowed
	DefaultHost = "platform-api.max.ru"

	DefaultHostV2   = "platform-api2.max.ru"
	defaultFileName = "file"

	SecretHeader        = "X-Max-Bot-Api-Secret"
	AuthorizationHeader = "Authorization"

	maxRetries      = 3
	defaultTimeout  = 30 * time.Second
	defaultPause    = time.Second
	maxUpdatesLimit = 50
)

const (
	pathMe            = "/me"
	pathAnswers       = "/answers"
	pathUpdates       = "/updates"
	pathUpload        = "/uploads"
	pathMessages      = "/messages"
	pathSubscriptions = "/subscriptions"
	pathChats         = "/chats"

	formatPathMessageId               = "/messages/%s"
	formatPathVideoAttachmentDetails  = "/videos/%s"
	formatPathChatsID                 = "/chats/%d"
	formatPathChatPin                 = "/chats/%d/pin"
	formatPathChatsActions            = "/chats/%d/actions"
	formatPathChatsMembers            = "/chats/%d/members"
	formatPathChatsMembersMe          = "/chats/%d/members/me"
	formatPathChatsMembersAdmin       = "/chats/%d/members/admins"
	formatPathChatsMembersAdminDelete = "/chats/%d/members/admins/%d"
)

const (
	paramURL    = "url"
	paramType   = "type"
	paramMarker = "marker"

	paramUser           = "user"
	paramHash           = "hash"
	paramBlock          = "block"
	paramChatID         = "chat_id"
	paramUserID         = "user_id"
	paramUserIDs        = "user_ids"
	paramMessageID      = "message_id"
	paramMessageIDs     = "message_ids"
	paramCallbackID     = "callback_id"
	paramWebAppPlatform = "WebAppPlatform"
	paramWebAppVersion  = "WebAppVersion"
	paramWebAppData     = "WebAppData"

	fieldData = "data"

	paramTo                 = "to"
	paramCount              = "count"
	paramFrom               = "from"
	paramLimit              = "limit"
	paramTimeout            = "timeout"
	paramDisableLinkPreview = "disable_link_preview"
)
