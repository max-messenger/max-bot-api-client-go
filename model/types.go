package model

type ChatType string

const (
	ChatTypeDialog  ChatType = "dialog"
	ChatTypeChat    ChatType = "chat"
	ChatTypeChannel ChatType = "channel"
)

type ChatStatus string

const (
	ChatStatusActive    ChatStatus = "active"
	ChatStatusRemoved   ChatStatus = "removed"
	ChatStatusLeft      ChatStatus = "left"
	ChatStatusClosed    ChatStatus = "closed"
	ChatStatusSuspended ChatStatus = "suspended"
)

type ChatAdminPermission string

const (
	PermReadAllMessages  ChatAdminPermission = "read_all_messages"
	PermAddRemoveMembers ChatAdminPermission = "add_remove_members"
	PermAddAdmins        ChatAdminPermission = "add_admins"
	PermChangeChatInfo   ChatAdminPermission = "change_chat_info"
	PermPinMessage       ChatAdminPermission = "pin_message"
	PermEditLink         ChatAdminPermission = "edit_link"
	PermWrite            ChatAdminPermission = "write"
	PermEdit             ChatAdminPermission = "edit"
	PermDelete           ChatAdminPermission = "delete"
	PermCanCall          ChatAdminPermission = "can_call"
	PermViewStats        ChatAdminPermission = "view_stats"
)

type MessageLinkType string

const (
	LinkTypeForward MessageLinkType = "forward"
	LinkTypeReply   MessageLinkType = "reply"
)

type TextFormat string

const (
	FormatMarkdown TextFormat = "markdown"
	FormatHTML     TextFormat = "html"
)

type UploadType string

const (
	UploadImage UploadType = "image"
	UploadVideo UploadType = "video"
	UploadAudio UploadType = "audio"
	UploadFile  UploadType = "file"
)

type SenderAction string

const (
	ActionTypingOn     SenderAction = "typing_on"
	ActionSendingPhoto SenderAction = "sending_photo"
	ActionSendingVideo SenderAction = "sending_video"
	ActionSendingAudio SenderAction = "sending_audio"
	ActionSendingFile  SenderAction = "sending_file"
	ActionMarkSeen     SenderAction = "mark_seen"
)

type Intent string

const (
	IntentPositive Intent = "positive"
	IntentNegative Intent = "negative"
	IntentDefault  Intent = "default"
)

type UpdateType string

const (
	UpdateMessageCreated   UpdateType = "message_created"
	UpdateMessageCallback  UpdateType = "message_callback"
	UpdateMessageEdited    UpdateType = "message_edited"
	UpdateMessageRemoved   UpdateType = "message_removed"
	UpdateBotAdded         UpdateType = "bot_added"
	UpdateBotRemoved       UpdateType = "bot_removed"
	UpdateUserAdded        UpdateType = "user_added"
	UpdateUserRemoved      UpdateType = "user_removed"
	UpdateBotStarted       UpdateType = "bot_started"
	UpdateBotStopped       UpdateType = "bot_stopped"
	UpdateDialogCleared    UpdateType = "dialog_cleared"
	UpdateDialogRemoved    UpdateType = "dialog_removed"
	UpdateDialogMuted      UpdateType = "dialog_muted"
	UpdateDialogUnmuted    UpdateType = "dialog_unmuted"
	UpdateChatTitleChanged UpdateType = "chat_title_changed"
)

type MarkupType string

const (
	MarkupStrong        MarkupType = "strong"
	MarkupEmphasized    MarkupType = "emphasized"
	MarkupMonospaced    MarkupType = "monospaced"
	MarkupLink          MarkupType = "link"
	MarkupStrikethrough MarkupType = "strikethrough"
	MarkupUnderline     MarkupType = "underline"
	MarkupUserMention   MarkupType = "user_mention"
	MarkupHeading       MarkupType = "heading"
	MarkupHighlighted   MarkupType = "highlighted"
	MarkupQuote         MarkupType = "quote"
)

type ButtonType string

const (
	ButtonCallback       ButtonType = "callback"
	ButtonLink           ButtonType = "link"
	ButtonRequestGeo     ButtonType = "request_geo_location"
	ButtonRequestContact ButtonType = "request_contact"
	ButtonChat           ButtonType = "chat"
	ButtonMessage        ButtonType = "message"
	ButtonOpenApp        ButtonType = "open_app"
	ButtonClipboard      ButtonType = "clipboard"
)

type AttachmentType string

const (
	AttachImage          AttachmentType = "image"
	AttachVideo          AttachmentType = "video"
	AttachAudio          AttachmentType = "audio"
	AttachFile           AttachmentType = "file"
	AttachSticker        AttachmentType = "sticker"
	AttachContact        AttachmentType = "contact"
	AttachInlineKeyboard AttachmentType = "inline_keyboard"
	AttachReplyKeyboard  AttachmentType = "reply_keyboard"
	AttachLocation       AttachmentType = "location"
	AttachShare          AttachmentType = "share"
	AttachData           AttachmentType = "data"
)
