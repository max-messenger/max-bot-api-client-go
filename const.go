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
	pathMeCommands    = "/me/commands"
	pathAnswers       = "/answers"
	pathUpdates       = "/updates"
	pathUpload        = "/uploads"
	pathMessages      = "/messages"
	pathSubscriptions = "/subscriptions"

	formatPathMessageId               = "/messages/%s"
	formatPathComments                = "/messages/%s/comments"
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

	paramUser       = "user"
	paramChat       = "chat"
	paramIP         = "ip"
	paramQueryID    = "query_id"
	paramAuthDate   = "auth_date"
	paramStartParam = "start_param"
	paramHash       = "hash"
	paramBlock      = "block"
	paramChatID     = "chat_id"
	paramUserID     = "user_id"
	paramUserIDs    = "user_ids"
	paramMessageID  = "message_id"
	paramMessageIDs = "message_ids"
	paramCommentIDs = "comment_ids"
	paramCallbackID = "callback_id"
	paramWebAppData = "WebAppData"

	fieldData = "data"

	paramBefore             = "before"
	paramAfter              = "after"
	paramTo                 = "to"
	paramCount              = "count"
	paramFrom               = "from"
	paramLimit              = "limit"
	paramTimeout            = "timeout"
	paramDisableLinkPreview = "disable_link_preview"
)
