package schemes

import (
	"encoding/json"
	"strings"
	"time"
)

type ActionRequestBody struct {
	Action SenderAction `json:"action"`
}

type AttachmentType string

const (
	AttachmentImage    AttachmentType = "image"
	AttachmentVideo                   = "video"
	AttachmentAudio                   = "audio"
	AttachmentFile                    = "file"
	AttachmentContact                 = "contact"
	AttachmentSticker                 = "sticker"
	AttachmentShare                   = "share"
	AttachmentLocation                = "location"
	AttachmentKeyboard                = "inline_keyboard"
)

// Generic schema representing message attachment
type Attachment struct {
	Type AttachmentType `json:"type"`
}

func (a Attachment) GetAttachmentType() AttachmentType {
	return a.Type
}

type AttachmentInterface interface {
	GetAttachmentType() AttachmentType
}

type AttachmentPayload struct {
	// Media attachment URL
	Url string `json:"url"`
}

// Request to attach some data to message
type AttachmentRequest struct {
	Type AttachmentType `json:"type"`
}

type AudioAttachment struct {
	Attachment
	Payload MediaAttachmentPayload `json:"payload"`
}

// Request to attach audio to message. MUST be the only attachment in message
type AudioAttachmentRequest struct {
	AttachmentRequest
	Payload UploadedInfo `json:"payload"`
}

func NewAudioAttachmentRequest(payload UploadedInfo) *AudioAttachmentRequest {
	return &AudioAttachmentRequest{Payload: payload, AttachmentRequest: AttachmentRequest{Type: AttachmentAudio}}
}

type BotCommand struct {
	Name        string `json:"name"`                  // Command name
	Description string `json:"description,omitempty"` // Optional command description
}

type BotInfo struct {
	UserId        int64        `json:"user_id"`                   // Users identifier
	Name          string       `json:"name"`                      // Users visible name
	Username      string       `json:"username,omitempty"`        // Unique public user name. Can be `null` if user is not accessible or it is not set
	AvatarUrl     string       `json:"avatar_url,omitempty"`      // URL of avatar
	FullAvatarUrl string       `json:"full_avatar_url,omitempty"` // URL of avatar of a bigger size
	Commands      []BotCommand `json:"commands,omitempty"`        // Commands supported by bots
	Description   string       `json:"description,omitempty"`     // Bot description
}

type BotPatch struct {
	Name        string                         `json:"name,omitempty"`        // Visible name of bots
	Username    string                         `json:"username,omitempty"`    // Bot unique identifier. It can be any string 4-64 characters long containing any digit, letter or special symbols: \"-\" or \"_\". It **must** starts with a letter
	Description string                         `json:"description,omitempty"` // Bot description up to 16k characters long
	Commands    []BotCommand                   `json:"commands,omitempty"`    // Commands supported by bots. Pass empty list if you want to remove commands
	Photo       *PhotoAttachmentRequestPayload `json:"photo,omitempty"`       // Request to set bots photo
}

type Button struct {
	Type ButtonType `json:"type"`
	Text string     `json:"text"` // Visible text of button
}

func (b Button) GetType() ButtonType {
	return b.Type
}

func (b Button) GetText() string {
	return b.Text
}

type ButtonInterface interface {
	GetType() ButtonType
	GetText() string
}

// Send this object when your bots wants to react to when a button is pressed
type CallbackAnswer struct {
	Message      *NewMessageBody `json:"message,omitempty"`      // Fill this if you want to modify current message
	Notification string          `json:"notification,omitempty"` // Fill this if you just want to send one-time notification to user
}

// After pressing this type of button client sends to server payload it contains
type CallbackButton struct {
	Button
	Payload string `json:"payload"`          // Button payload
	Intent  Intent `json:"intent,omitempty"` // Intent of button. Affects clients representation
}

type CallbackButtonAllOf struct {
	Payload string `json:"payload"`          // Button payload
	Intent  Intent `json:"intent,omitempty"` // Intent of button. Affects clients representation
}

type Chat struct {
	ChatId            int64           `json:"chat_id"`                // Chats identifier
	Type              ChatType        `json:"type"`                   // Type of chat. One of: dialog, chat, channel
	Status            ChatStatus      `json:"status"`                 // Chat status. One of:  - active: bots is active member of chat  - removed: bots was kicked  - left: bots intentionally left chat  - closed: chat was closed
	Title             string          `json:"title,omitempty"`        // Visible title of chat. Can be null for dialogs
	Icon              *Image          `json:"icon"`                   // Icon of chat
	LastEventTime     int             `json:"last_event_time"`        // Time of last event occurred in chat
	ParticipantsCount int             `json:"participants_count"`     // Number of people in chat. Always 2 for `dialog` chat type
	OwnerId           int64           `json:"owner_id,omitempty"`     // Identifier of chat owner. Visible only for chat admins
	Participants      *map[string]int `json:"participants,omitempty"` // Participants in chat with time of last activity. Can be *null* when you request list of chats. Visible for chat admins only
	IsPublic          bool            `json:"is_public"`              // Is current chat publicly available. Always `false` for dialogs
	Link              string          `json:"link,omitempty"`         // Link on chat if it is public
	Description       *string         `json:"description"`            // Chat description
}

// ChatAdminPermission : Chat admin permissions
type ChatAdminPermission string

// List of ChatAdminPermission
const (
	READ_ALL_MESSAGES  ChatAdminPermission = "read_all_messages"
	ADD_REMOVE_MEMBERS ChatAdminPermission = "add_remove_members"
	ADD_ADMINS         ChatAdminPermission = "add_admins"
	CHANGE_CHAT_INFO   ChatAdminPermission = "change_chat_info"
	PIN_MESSAGE        ChatAdminPermission = "pin_message"
	WRITE              ChatAdminPermission = "write"
)

type ChatList struct {
	Chats  []Chat `json:"chats"`  // List of requested chats
	Marker *int   `json:"marker"` // Reference to the next page of requested chats
}

type ChatMember struct {
	UserId         int64                 `json:"user_id"`                   // Users identifier
	Name           string                `json:"name"`                      // Users visible name
	Username       string                `json:"username,omitempty"`        // Unique public user name. Can be `null` if user is not accessible or it is not set
	AvatarUrl      string                `json:"avatar_url,omitempty"`      // URL of avatar
	FullAvatarUrl  string                `json:"full_avatar_url,omitempty"` // URL of avatar of a bigger size
	LastAccessTime int                   `json:"last_access_time"`
	IsOwner        bool                  `json:"is_owner"`
	IsAdmin        bool                  `json:"is_admin"`
	JoinTime       int                   `json:"join_time"`
	Permissions    []ChatAdminPermission `json:"permissions,omitempty"` // Permissions in chat if member is admin. `null` otherwise
}

type ChatMembersList struct {
	Members []ChatMember `json:"members"` // Participants in chat with time of last activity. Visible only for chat admins
	Marker  *int64       `json:"marker"`  // Pointer to the next data page
}

type ChatPatch struct {
	Icon  *PhotoAttachmentRequestPayload `json:"icon,omitempty"`
	Title string                         `json:"title,omitempty"`
}

// ChatStatus : Chat status for current bots
type ChatStatus string

// List of ChatStatus
const (
	ACTIVE    ChatStatus = "active"
	REMOVED   ChatStatus = "removed"
	LEFT      ChatStatus = "left"
	CLOSED    ChatStatus = "closed"
	SUSPENDED ChatStatus = "suspended"
)

// ChatType : Type of chat. Dialog (one-on-one), chat or channel
type ChatType string

// List of ChatType
const (
	DIALOG  ChatType = "dialog"
	CHAT    ChatType = "chat"
	CHANNEL ChatType = "channel"
)

type ContactAttachment struct {
	Attachment
	Payload ContactAttachmentPayload `json:"payload"`
}

type ContactAttachmentPayload struct {
	VcfInfo string `json:"vcf_info,omitempty"` // User info in VCF format
	TamInfo *User  `json:"max_info"`           // User info
}

// Request to attach contact card to message. MUST be the only attachment in message
type ContactAttachmentRequest struct {
	AttachmentRequest
	Payload ContactAttachmentRequestPayload `json:"payload"`
}

func NewContactAttachmentRequest(payload ContactAttachmentRequestPayload) *ContactAttachmentRequest {
	return &ContactAttachmentRequest{Payload: payload, AttachmentRequest: AttachmentRequest{Type: AttachmentContact}}
}

type ContactAttachmentRequestPayload struct {
	Name      string `json:"name,omitempty"`       // Contact name
	ContactId int64  `json:"contact_id,omitempty"` // Contact identifier
	VcfInfo   string `json:"vcf_info,omitempty"`   // Full information about contact in VCF format
	VcfPhone  string `json:"vcf_phone,omitempty"`  // Contact phone in VCF format
}

// Server returns this if there was an exception to your request
type Error struct {
	ErrorText   string    `json:"error,omitempty"`                  // Error
	Code        string    `json:"code,omitempty"`                   // Error code
	Message     Message   `json:"message,omitempty"`                // Human-readable description
	Results     []Results `json:"results,omitempty"`                // phones
	NumberExist []string  `json:"existing_phone_numbers,omitempty"` // extistphones

}

type Results struct {
	PhoneNumber string `json:"phone_number,omitempty"` // Phones
	Status      string `json:"status,omitempty"`       // Status delivery
}

func (e Error) Error() string {
	return e.ErrorText
}

type FileAttachment struct {
	Attachment
	Payload  FileAttachmentPayload `json:"payload"`
	Filename string                `json:"filename"` // Uploaded file name
	Size     int64                 `json:"size"`     // File size in bytes
}

type FileAttachmentPayload struct {
	Url   string `json:"url"`   // Media attachment URL
	Token string `json:"token"` // Use `token` in case when you are trying to reuse the same attachment in other message
}

// Request to attach file to message. MUST be the only attachment in message
type FileAttachmentRequest struct {
	AttachmentRequest
	Payload UploadedInfo `json:"payload"`
}

func NewFileAttachmentRequest(payload UploadedInfo) *FileAttachmentRequest {
	return &FileAttachmentRequest{Payload: payload, AttachmentRequest: AttachmentRequest{Type: AttachmentFile}}
}

// List of all WebHook subscriptions
type GetSubscriptionsResult struct {
	Subscriptions []Subscription `json:"subscriptions"` // Current subscriptions
}

// Generic schema describing image object
type Image struct {
	Url string `json:"url"` // URL of image
}

// Buttons in messages
type InlineKeyboardAttachment struct {
	Attachment
	Payload Keyboard `json:"payload"`
}

// Request to attach keyboard to message
type InlineKeyboardAttachmentRequest struct {
	AttachmentRequest
	Payload Keyboard `json:"payload"`
}

func NewInlineKeyboardAttachmentRequest(payload Keyboard) *InlineKeyboardAttachmentRequest {
	return &InlineKeyboardAttachmentRequest{Payload: payload, AttachmentRequest: AttachmentRequest{Type: AttachmentKeyboard}}
}

type ButtonType string

const (
	LINK        ButtonType = "link"
	CALLBACK    ButtonType = "callback"
	CONTACT     ButtonType = "request_contact"
	GEOLOCATION ButtonType = "request_geo_location"
)

// Intent : Intent of button
type Intent string

// List of Intent
const (
	POSITIVE Intent = "positive"
	NEGATIVE Intent = "negative"
	DEFAULT  Intent = "default"
)

// Keyboard is two-dimension array of buttons
type Keyboard struct {
	Buttons [][]ButtonInterface `json:"buttons"`
}

// After pressing this type of button user follows the link it contains
type LinkButton struct {
	Button
	Url string `json:"url"`
}

type LinkedMessage struct {
	Type    MessageLinkType `json:"type"`              // Type of linked message
	Sender  User            `json:"sender,omitempty"`  // User sent this message. Can be `null` if message has been posted on behalf of a channel
	ChatId  int64           `json:"chat_id,omitempty"` // Chat where message has been originally posted. For forwarded messages only
	Message MessageBody     `json:"message"`
}

type LocationAttachment struct {
	Attachment
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Request to attach keyboard to message
type LocationAttachmentRequest struct {
	AttachmentRequest
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewLocationAttachmentRequest(latitude float64, longitude float64) *LocationAttachmentRequest {
	return &LocationAttachmentRequest{Latitude: latitude, Longitude: longitude, AttachmentRequest: AttachmentRequest{Type: AttachmentLocation}}
}

type MediaAttachmentPayload struct {
	Url   string `json:"url"`   // Media attachment URL
	Token string `json:"token"` // Use `token` in case when you are trying to reuse the same attachment in other message
}

// Message in chat
type Message struct {
	Sender    User           `json:"sender,omitempty"` // User that sent this message. Can be `null` if message has been posted on behalf of a channel
	Recipient Recipient      `json:"recipient"`        // Message recipient. Could be user or chat
	Timestamp int64          `json:"timestamp"`        // Unix-time when message was created
	Link      *LinkedMessage `json:"link,omitempty"`   // Forwarder or replied message
	Body      MessageBody    `json:"body"`             // Body of created message. Text + attachments. Could be null if message contains only forwarded message
	Stat      *MessageStat   `json:"stat,omitempty"`   // Message statistics. Available only for channels in [GET:/messages](#operation/getMessages) context
}

// Schema representing body of message
type MessageBody struct {
	Mid            string            `json:"mid"`            // Unique identifier of message
	Seq            int64             `json:"seq"`            // Sequence identifier of message in chat
	Text           string            `json:"text,omitempty"` // Message text
	RawAttachments []json.RawMessage `json:"attachments"`    // Message attachments. Could be one of `Attachment` type. See description of this schema
	Attachments    []interface{}
	ReplyTo        string `json:"reply_to,omitempty"` // In case this message is reply to another, it is the unique identifier of the replied message
}

type UpdateType string

const (
	TypeMessageCallback  UpdateType = "message_callback"
	TypeMessageCreated   UpdateType = "message_created"
	TypeMessageRemoved   UpdateType = "message_removed"
	TypeMessageEdited    UpdateType = "message_edited"
	TypeBotAdded         UpdateType = "bot_added"
	TypeBotRemoved       UpdateType = "bot_removed"
	TypeUserAdded        UpdateType = "user_added"
	TypeUserRemoved      UpdateType = "user_removed"
	TypeBotStarted       UpdateType = "bot_started"
	TypeChatTitleChanged UpdateType = "chat_title_changed"
)

// MessageLinkType : Type of linked message
type MessageLinkType string

// List of MessageLinkType
const (
	FORWARD MessageLinkType = "forward"
	REPLY   MessageLinkType = "reply"
)

// MarkupType : Type of markups
type MarkupType string

// List of MarkupType
const (
	MarkupUser MarkupType = "user_mention"
	MarkupBot  MarkupType = "bot_mention"
)

// Paginated list of messages
type MessageList struct {
	Messages []Message `json:"messages"` // List of messages
}

// Message statistics
type MessageStat struct {
	Views int `json:"views"`
}

type NewMessageBody struct {
	BotToken     string          `json:"bot_token,omitempty"`     // bot
	Text         string          `json:"text,omitempty"`          // Message text
	Attachments  []interface{}   `json:"attachments,omitempty"`   // Message attachments. See `AttachmentRequest` and it's inheritors for full information
	Link         *NewMessageLink `json:"link,omitempty"`          // Link to Message
	Format       string          `json:"format,omitempty"`        // Format to Message
	PhoneNumbers []string        `json:"phone_numbers,omitempty"` // PhoneNumber to Message
	Notify       bool            `json:"notify,omitempty"`        // If false, chat participants wouldn't be notified
	Markups      []MarkUp        `json:"markup,omitempty"`        // mention users
}

// Generic schema representing message markup
type Markup struct {
	Type MarkupType `json:"type"` // Type of markup
}

func (a Markup) GetMarkupType() MarkupType {
	return a.Type
}

type MarkupInterface interface {
	GetMarkupType() MarkupType
}

type MarkUpUser struct {
	Markup
	From   int   `json:"from"`              // where starts
	Length int   `json:"length"`            // length of marker message
	UserId int64 `json:"user_id,omitempty"` // User identifier, if message was sent to user
}

// old
type MarkUp struct {
	From   int        `json:"from"`              // where starts
	Length int        `json:"length"`            // length of marker message
	UserId int64      `json:"user_id,omitempty"` // User identifier, if message was sent to user
	Type   MarkupType `json:"type"`              // Type of markup
}

type NewMessageLink struct {
	Type MessageLinkType `json:"type"` // Type of message link
	Mid  string          `json:"mid"`  // Message identifier of original message
}

// Image attachment
type PhotoAttachment struct {
	Attachment
	Payload PhotoAttachmentPayload `json:"payload"`
}

type PhotoAttachmentPayload struct {
	PhotoId int64  `json:"photo_id"` // Unique identifier of this image
	Token   string `json:"token"`
	Url     string `json:"url"` // Image URL
}

type PhotoAttachmentRequest struct {
	AttachmentRequest
	Payload PhotoAttachmentRequestPayload `json:"payload"`
}

func NewPhotoAttachmentRequest(payload PhotoAttachmentRequestPayload) *PhotoAttachmentRequest {
	return &PhotoAttachmentRequest{Payload: payload, AttachmentRequest: AttachmentRequest{Type: AttachmentImage}}
}

type PhotoAttachmentRequestAllOf struct {
	Payload PhotoAttachmentRequestPayload `json:"payload"`
}

// Request to attach image. All fields are mutually exclusive
type PhotoAttachmentRequestPayload struct {
	Url    string                `json:"url,omitempty"`    // Any external image URL you want to attach
	Token  string                `json:"token,omitempty"`  // Token of any existing attachment
	Photos map[string]PhotoToken `json:"photos,omitempty"` // Tokens were obtained after uploading images
}

type PhotoToken struct {
	Token string `json:"token"` // Encoded information of uploaded image
}

// This is information you will receive as soon as an image uploaded
type PhotoTokens struct {
	Photos map[string]PhotoToken `json:"photos"`
}

// PinMessageBody defines model for PinMessageBody.
type PinMessageBody struct {
	// MessageId Identifier of message to be pinned in chat
	MessageId string `json:"message_id"`

	// Notify If `true`, participants will be notified with system message in chat/channel
	Notify *bool `json:"notify"`
}

// Recipient New message recipient. Could be user or chat
type Recipient struct {
	ChatId   int64    `json:"chat_id,omitempty"` // Chat identifier
	ChatType ChatType `json:"chat_type"`         // Chat type
	UserId   int64    `json:"user_id,omitempty"` // User identifier, if message was sent to user
}

// After pressing this type of button client sends new message with attachment of current user contact
type RequestContactButton struct {
	Button
}

// After pressing this type of button client sends new message with attachment of current user geo location
type RequestGeoLocationButton struct {
	Button
	Quick bool `json:"quick,omitempty"` // If *true*, sends location without asking user's confirmation
}

type SendMessageResult struct {
	Message Message `json:"message"`
}

// SenderAction : Different actions to send to chat members
type SenderAction string

// List of SenderAction
const (
	TYPING_ON     SenderAction = "typing_on"
	TYPING_OFF    SenderAction = "typing_off"
	SENDING_PHOTO SenderAction = "sending_photo"
	SENDING_VIDEO SenderAction = "sending_video"
	SENDING_AUDIO SenderAction = "sending_audio"
	MARK_SEEN     SenderAction = "mark_seen"
)

type ShareAttachment struct {
	Attachment
	Payload AttachmentPayload `json:"payload"`
}

// Simple response to request
type SimpleQueryResult struct {
	Success bool   `json:"success"`           // `true` if request was successful. `false` otherwise
	Message string `json:"message,omitempty"` // Explanatory message if the result is not successful
}

type StickerAttachment struct {
	Attachment
	Payload StickerAttachmentPayload `json:"payload"`
	Width   int                      `json:"width"`  // Sticker width
	Height  int                      `json:"height"` // Sticker height
}

type StickerAttachmentPayload struct {
	Url  string `json:"url"`  // Media attachment URL
	Code string `json:"code"` // Sticker identifier
}

// Request to attach sticker. MUST be the only attachment request in message
type StickerAttachmentRequest struct {
	AttachmentRequest
	Payload StickerAttachmentRequestPayload `json:"payload"`
}

func NewStickerAttachmentRequest(payload StickerAttachmentRequestPayload) *StickerAttachmentRequest {
	return &StickerAttachmentRequest{Payload: payload, AttachmentRequest: AttachmentRequest{Type: AttachmentSticker}}
}

type StickerAttachmentRequestPayload struct {
	Code string `json:"code"` // Sticker code
}

// Schema to describe WebHook subscription
type Subscription struct {
	Url         string   `json:"url"`                    // Webhook URL
	Time        int64    `json:"time"`                   // Unix-time when subscription was created
	UpdateTypes []string `json:"update_types,omitempty"` // Update types bots subscribed for
	Version     string   `json:"version,omitempty"`
}

// Request to set up WebHook subscription
type SubscriptionRequestBody struct {
	// Secret A secret to be sent in a header “X-Max-Bot-Api-Secret” in every webhook request, 5-256 characters. Only characters A-Z, a-z, 0-9, _ and - are allowed. The header is useful to ensure that the request comes from a webhook set by you.
	Secret      string   `json:"secret,omitempty"`
	Url         string   `json:"url"`                    // URL of HTTP(S)-endpoint of your bots. Must starts with http(s)://
	UpdateTypes []string `json:"update_types,omitempty"` // List of update types your bots want to receive. See `Update` object for a complete list of types
	Version     string   `json:"version,omitempty"`      // Version of API. Affects model representation
}

// List of all updates in chats your bots participated in
type UpdateList struct {
	Updates []json.RawMessage `json:"updates"` // Page of updates
	Marker  *int64            `json:"marker"`  // Pointer to the next data page
}

// Endpoint you should upload to your binaries
type UploadEndpoint struct {
	// Token Video or audio token for send message
	Token string `json:"token,omitempty"`
	Url   string `json:"url"` // URL to upload
}

// UploadType : Type of file uploading
type UploadType string

// List of UploadType
const (
	PHOTO UploadType = "photo"
	VIDEO UploadType = "video"
	AUDIO UploadType = "audio"
	FILE  UploadType = "file"
)

// This is information you will receive as soon as audio/video is uploaded
type UploadedInfo struct {
	FileID int64  `json:"file_id,omitempty"`
	Token  string `json:"token,omitempty"` // Token is unique uploaded media identfier
}

type User struct {
	UserId                int64         `json:"user_id"`            // Users identifier
	Name                  string        `json:"name"`               // Users visible name
	Username              string        `json:"username,omitempty"` // Unique public user name. Can be `null` if user is not accessible or it is not set
	FirstName             string        `json:"first_name,omitempty"`
	LastName              string        `json:"last_name,omitempty"`
	IsBot                 bool          `json:"is_bot,omitempty"`
	LastActivityTimeIsBot time.Duration `json:"last_activity_time,omitempty"`
}

type UserIdsList struct {
	UserIds []int `json:"user_ids"`
}

type UserWithPhoto struct {
	User
	AvatarUrl     string `json:"avatar_url,omitempty"`      // URL of avatar
	FullAvatarUrl string `json:"full_avatar_url,omitempty"` // URL of avatar of a bigger size
}

type VideoAttachment struct {
	Attachment
	Payload MediaAttachmentPayload `json:"payload"`
}

// Request to attach video to message
type VideoAttachmentRequest struct {
	AttachmentRequest
	Payload UploadedInfo `json:"payload"`
}

func NewVideoAttachmentRequest(payload UploadedInfo) *VideoAttachmentRequest {
	return &VideoAttachmentRequest{Payload: payload, AttachmentRequest: AttachmentRequest{Type: AttachmentVideo}}
}

// `Update` object represents different types of events that happened in chat. See its inheritors
type Update struct {
	UpdateType UpdateType `json:"update_type"`
	Timestamp  int        `json:"timestamp"` // Unix-time when event has occurred
	DebugRaw   string
}

func (u Update) GetUpdateType() UpdateType {
	return u.UpdateType
}

func (u Update) GetDebugRaw() string {
	return u.DebugRaw
}

func (u Update) GetUpdateTime() time.Time {
	return time.Unix(int64(u.Timestamp/1000), 0)
}

type UpdateInterface interface {
	GetDebugRaw() string
	GetUpdateType() UpdateType
	GetUpdateTime() time.Time
	GetUserID() int64
	GetChatID() int64
}

// You will receive this update when bots has been added to chat
type BotAddedToChatUpdate struct {
	Update
	ChatId int64 `json:"chat_id"` // Chat id where bots was added
	User   User  `json:"user"`    // User who added bots to chat
}

func (b BotAddedToChatUpdate) GetUserID() int64 {
	return b.User.UserId
}

func (b BotAddedToChatUpdate) GetChatID() int64 {
	return b.ChatId
}

// You will receive this update when bots has been removed from chat
type BotRemovedFromChatUpdate struct {
	Update
	ChatId int64 `json:"chat_id"` // Chat identifier bots removed from
	User   User  `json:"user"`    // User who removed bots from chat
}

func (b BotRemovedFromChatUpdate) GetUserID() int64 {
	return b.User.UserId
}

func (b BotRemovedFromChatUpdate) GetChatID() int64 {
	return b.ChatId
}

// Bot gets this type of update as soon as user pressed `Start` button
type BotStartedUpdate struct {
	Update
	ChatId int64 `json:"chat_id"` // Dialog identifier where event has occurred
	User   User  `json:"user"`    // User pressed the 'Start' button
}

func (b BotStartedUpdate) GetUserID() int64 {
	return b.User.UserId
}

func (b BotStartedUpdate) GetChatID() int64 {
	return b.ChatId
}

// Object sent to bots when user presses button
type Callback struct {
	Timestamp  int64  `json:"timestamp"` // Unix-time when event has occurred
	CallbackID string `json:"callback_id"`
	Payload    string `json:"payload,omitempty"` // Button payload
	User       User   `json:"user"`              // User pressed the button
}

func (b Callback) GetUserID() int64 {
	return b.User.UserId
}

func (b Callback) GetChatID() int64 {
	return 0
}

// Bot gets this type of update as soon as title has been changed in chat
type ChatTitleChangedUpdate struct {
	Update
	ChatId int64  `json:"chat_id"` // Chat identifier where event has occurred
	User   User   `json:"user"`    // User who changed title
	Title  string `json:"title"`   // New title
}

func (b ChatTitleChangedUpdate) GetUserID() int64 {
	return b.User.UserId
}

func (b ChatTitleChangedUpdate) GetChatID() int64 {
	return b.ChatId
}

// You will get this `update` as soon as user presses button
type MessageCallbackUpdate struct {
	Update
	Callback Callback `json:"callback"`
	Message  *Message `json:"message"` // Original message containing inline keyboard. Can be `null` in case it had been deleted by the moment a bots got this update
}

func (b MessageCallbackUpdate) GetUserID() int64 {
	return b.Callback.User.UserId
}

func (b MessageCallbackUpdate) GetChatID() int64 {
	return 0
}

// You will get this `update` as soon as message is created
type MessageCreatedUpdate struct {
	Update
	Message Message `json:"message"` // Newly created message
}

func (b MessageCreatedUpdate) GetUserID() int64 {
	return b.Message.Sender.UserId
}

func (b MessageCreatedUpdate) GetChatID() int64 {
	return b.Message.Recipient.ChatId
}

func (b MessageCreatedUpdate) GetText() string {
	return b.Message.Body.Text
}

func (b MessageCreatedUpdate) GetCommand() string {
	if strings.Index(b.Message.Body.Text, "/") == 0 {
		if strings.Contains(b.Message.Body.Text, ":") {
			return strings.Split(b.Message.Body.Text, ":")[0]
		}
		return b.Message.Body.Text
	}
	return "undefened"
}

func (b MessageCreatedUpdate) GetParam() string {
	if strings.Index(b.Message.Body.Text, "/") == 0 {
		if strings.Contains(b.Message.Body.Text, ":") {
			return strings.Split(b.Message.Body.Text, ":")[1]
		}
		return ""
	}
	return ""
}

// You will get this `update` as soon as message is edited
type MessageEditedUpdate struct {
	Update
	Message Message `json:"message"` // Edited message
}

func (b MessageEditedUpdate) GetUserID() int64 {
	return b.Message.Sender.UserId
}

func (b MessageEditedUpdate) GetChatID() int64 {
	return b.Message.Recipient.ChatId
}

// You will get this `update` as soon as message is removed
type MessageRemovedUpdate struct {
	Update
	MessageId string `json:"message_id"` // Identifier of removed message
}

func (b MessageRemovedUpdate) GetUserID() int64 {
	return 0
}

func (b MessageRemovedUpdate) GetChatID() int64 {
	return 0
}

// You will receive this update when user has been added to chat where bots is administrator
type UserAddedToChatUpdate struct {
	Update
	ChatId    int64 `json:"chat_id"`    // Chat identifier where event has occurred
	User      User  `json:"user"`       // User added to chat
	InviterId int64 `json:"inviter_id"` // User who added user to chat
}

func (b UserAddedToChatUpdate) GetUserID() int64 {
	return b.User.UserId
}

func (b UserAddedToChatUpdate) GetChatID() int64 {
	return b.ChatId
}

// You will receive this update when user has been removed from chat where bots is administrator
type UserRemovedFromChatUpdate struct {
	Update
	ChatId  int64 `json:"chat_id"`  // Chat identifier where event has occurred
	User    User  `json:"user"`     // User removed from chat
	AdminId int64 `json:"admin_id"` // Administrator who removed user from chat
}

func (b UserRemovedFromChatUpdate) GetUserID() int64 {
	return b.User.UserId
}

func (b UserRemovedFromChatUpdate) GetChatID() int64 {
	return b.ChatId
}
