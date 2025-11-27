package admins

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

type Administrator struct {
	UserId      int64                 `json:"user_id"`               // Users identifier
	Name        string                `json:"name"`                  // Users visible name
	Username    string                `json:"username,omitempty"`    // Unique public user name. Can be `null` if user is not accessible or it is not set
	Permissions []ChatAdminPermission `json:"permissions,omitempty"` // Permissions in chat if member is admin. `null` otherwise
}

type AdminMembersList struct {
	Admins []Administrator `json:"admins"` // Participants in chat with time of last activity. Visible only for chat admins
	Marker *int64          `json:"marker"` // Pointer to the next data page
}
