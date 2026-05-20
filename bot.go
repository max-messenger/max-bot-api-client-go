package maxbot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

var commandReg = regexp.MustCompile(`^(/\w+)`)

type BotsAPI interface {
	GetMyInfo(ctx context.Context) (model.BotInfo, error)
	EditMyInfo(ctx context.Context, patch model.BotPatch) (model.BotInfo, error)
}

type ChatsAPI interface {
	GetChats(ctx context.Context, count, marker int64) (model.ChatList, error)
	GetChat(ctx context.Context, chatID int64) (model.Chat, error)
	EditChat(ctx context.Context, chatID int64, patch model.ChatPatch) (model.Chat, error)
	DeleteChat(ctx context.Context, chatID int64) (model.SimpleQueryResult, error)
	SendAction(ctx context.Context, chatID int64, action model.SenderAction) (model.SimpleQueryResult, error)
	GetPinnedMessage(ctx context.Context, chatID int64) (model.GetPinnedMessageResult, error)
	PinMessage(ctx context.Context, chatID int64, messageID string, notify bool) (model.SimpleQueryResult, error)
	UnpinMessage(ctx context.Context, chatID int64) (model.SimpleQueryResult, error)
	GetMembership(ctx context.Context, chatID int64) (model.ChatMember, error)
	LeaveChat(ctx context.Context, chatID int64) (model.SimpleQueryResult, error)
	GetAdmins(ctx context.Context, chatID int64) (model.ChatMembersList, error)
	SetAdmins(ctx context.Context, chatID int64, admins []model.ChatAdmin) (model.SimpleQueryResult, error)
	DeleteAdmins(ctx context.Context, chatID, userID int64) (model.SimpleQueryResult, error)
	GetMembers(ctx context.Context, chatID, marker, count int64, userIDs []int64) (model.ChatMembersList, error)
	AddMembers(ctx context.Context, chatID int64, userIDs []int64) (model.SimpleQueryResult, error)
	RemoveMember(ctx context.Context, chatID, userID int64, block bool) (model.SimpleQueryResult, error)
}

type MessagesAPI interface {
	GetMessages(ctx context.Context, chatID, from, to, count int64, messageIDs []string) (model.MessageList, error)
	GetMessageByID(ctx context.Context, messageID string) (model.Message, error)
	Send(ctx context.Context, msg *Message) (res model.SendMessageResult, err error)
	EditMessage(ctx context.Context, messageID string, body model.NewMessageBody) (model.SimpleQueryResult, error)
	DeleteMessage(ctx context.Context, messageID string) (model.SimpleQueryResult, error)
	AnswerOnCallback(ctx context.Context, callbackID string, answer model.CallbackAnswer) (model.SimpleQueryResult, error)
	GetVideoAttachmentDetails(ctx context.Context, videoToken string) (model.VideoAttachmentDetails, error)
}

type SubscriptionsAPI interface {
	GetSubscriptions(ctx context.Context) (model.GetSubscriptionsResult, error)
	Subscribe(ctx context.Context, url, secret string, updateTypes []string, version string) (model.SimpleQueryResult, error)
	Unsubscribe(ctx context.Context, url string) (model.SimpleQueryResult, error)
	GetUpdates(ctx context.Context, marker int64) ([]model.Update, int64, error)
}

type UploadAPI interface {
	Upload(ctx context.Context, uploadType model.UploadType, reader io.Reader, name string, size int64) (string, error)
}

func NewApi(token string, opt ...Opt) (*Api, error) {
	cli := newClient(token, DefaultHost)
	var err error
	for _, o := range opt {
		err = o(cli)
		if err != nil {
			return nil, err
		}
	}

	b := &Api{
		Bots:          newBots(cli),
		Upload:        newUpload(cli),
		Chats:         newChats(cli),
		Messages:      newMessages(cli),
		Subscriptions: newSubscriptions(cli),
	}

	return b, nil
}

type Api struct {
	Bots          BotsAPI
	Upload        UploadAPI
	Chats         ChatsAPI
	Messages      MessagesAPI
	Subscriptions SubscriptionsAPI
}

func (a *Api) GetHandler(handler UpdateHandler, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler == nil {
			http.Error(w, "handler is nil", http.StatusInternalServerError)

			return
		}

		if secret != r.Header.Get(SecretHeader) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)

			return
		}

		update := updateRaw{}
		err = json.Unmarshal(body, &update)
		if err != nil {
			http.Error(w, "Failed to parse update", http.StatusBadRequest)

			return
		}

		handler(r.Context(), update.FromRaw())
	}
}

// ValidateInitData Проверяет подпись запроса от Max MiniApp. Возвращает пользователя.
func ValidateInitData(initData string, botToken string) (res model.UserApp, err error) {
	if initData == "" {
		err = fmt.Errorf("initData cannot be empty")

		return
	}

	if botToken == "" {
		err = fmt.Errorf("botToken cannot be empty")

		return
	}

	decodedInitData, err := url.QueryUnescape(initData)
	if err != nil {
		decodedInitData = initData
	}

	values, err := url.ParseQuery(decodedInitData)
	if err != nil {
		err = fmt.Errorf("failed to parse initData: %w", err)

		return
	}

	hashValues := values[paramHash]
	if len(hashValues) == 0 {
		err = fmt.Errorf("hash parameter is missing")

		return
	}

	receivedHash := hashValues[0]
	values.Del(paramHash)

	values.Del(paramWebAppPlatform)
	values.Del(paramWebAppVersion)

	var sortedParams []string
	for key := range values {
		value := values.Get(key)
		sortedParams = append(sortedParams, fmt.Sprintf("%s=%s", key, value))
	}
	sort.Strings(sortedParams)

	dataCheckString := strings.Join(sortedParams, "\n")

	mac1 := hmac.New(sha256.New, []byte(paramWebAppData))
	mac1.Write([]byte(botToken))
	secretKey := mac1.Sum(nil)

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(receivedHash), []byte(expectedHash)) != 1 {
		err = fmt.Errorf("hash verification failed")

		return
	}

	err = json.Unmarshal([]byte(values.Get(paramUser)), &res)
	if err != nil {
		err = fmt.Errorf("json decode err: %w", err)

		return
	}

	return
}

func GetCommand(u model.Update) string {
	match := commandReg.FindAllString(u.Message.Body.Text, -1)
	if len(match) > 0 {
		return match[0]
	}

	return ""
}
