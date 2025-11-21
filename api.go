// Package maxbot implements MAX Bot API.
// Official documentation: https://dev.max.ru/
package maxbot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/configservice"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

const (
	version = "1.2.5"

	defaultTimeout  = 30 * time.Second
	defaultAPIURL   = "https://botapi.max.ru/"
	defaultPause    = 1 * time.Second
	maxUpdatesLimit = 50

	maxRetries = 3
)

// Api represents the MAX Bot API client
type Api struct {
	Bots          *bots
	Chats         *chats
	Debugs        *debugs
	Messages      *messages
	Subscriptions *subscriptions
	Uploads       *uploads

	client  *client
	timeout time.Duration
	pause   time.Duration
	debug   bool
}

// New creates a new Max Bot API client with the provided token
func New(token string) (*Api, error) {
	if token == "" {
		return nil, ErrEmptyToken
	}

	u, err := url.Parse(defaultAPIURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	cl := newClient(token, version, u, &http.Client{
		Timeout: defaultTimeout,
	})

	api := &Api{
		client:  cl,
		timeout: defaultTimeout,
		pause:   defaultPause,
		debug:   false,
	}

	// Initialize sub-clients
	api.Bots = newBots(cl)
	api.Chats = newChats(cl)
	api.Uploads = newUploads(cl)
	api.Messages = newMessages(cl)
	api.Subscriptions = newSubscriptions(cl)
	api.Debugs = newDebugs(cl, 0)

	return api, nil
}

// NewWithConfig creates a new Max Bot API client from configuration service
func NewWithConfig(cfg configservice.ConfigInterface) (*Api, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	token := cfg.BotTokenCheckString()
	if token == "" {
		token = os.Getenv("TOKEN")
		if token == "" {
			return nil, ErrEmptyToken
		}
	}

	timeout := time.Duration(cfg.GetHttpBotAPITimeOut()) * time.Second
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	baseURL := cfg.GetHttpBotAPIUrl()
	if baseURL == "" {
		baseURL = defaultAPIURL
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	apiVersion := cfg.GetHttpBotAPIVersion()
	if apiVersion == "" {
		apiVersion = version
	}

	cl := newClient(token, apiVersion, u, &http.Client{
		Timeout: timeout,
	})

	api := &Api{
		client:  cl,
		timeout: timeout,
		pause:   defaultPause,
		debug:   cfg.GetDebugLogMode(),
	}

	// Initialize sub-clients
	api.Bots = newBots(cl)
	api.Chats = newChats(cl)
	api.Uploads = newUploads(cl)
	api.Messages = newMessages(cl)
	api.Subscriptions = newSubscriptions(cl)
	api.Debugs = newDebugs(cl, cfg.GetDebugLogChat())

	return api, nil
}

// updateTypeMap maps update types to their corresponding struct constructors
var updateTypeMap = map[schemes.UpdateType]func(debugRaw string) schemes.UpdateInterface{
	schemes.TypeMessageCallback: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.MessageCallbackUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
	schemes.TypeMessageCreated: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.MessageCreatedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
	schemes.TypeMessageRemoved: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.MessageRemovedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
	schemes.TypeMessageEdited: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.MessageEditedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
	schemes.TypeBotAdded: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.BotAddedToChatUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
	schemes.TypeBotRemoved: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.BotRemovedFromChatUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
	schemes.TypeUserAdded: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.UserAddedToChatUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
	schemes.TypeUserRemoved: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.UserRemovedFromChatUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
	schemes.TypeBotStarted: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.BotStartedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
	schemes.TypeChatTitleChanged: func(debugRaw string) schemes.UpdateInterface {
		return &schemes.ChatTitleChangedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
	},
}

// attachmentTypeMap maps attachment types to their corresponding struct constructors
var attachmentTypeMap = map[schemes.AttachmentType]func() schemes.AttachmentInterface{
	schemes.AttachmentAudio:    func() schemes.AttachmentInterface { return new(schemes.AudioAttachment) },
	schemes.AttachmentContact:  func() schemes.AttachmentInterface { return new(schemes.ContactAttachment) },
	schemes.AttachmentFile:     func() schemes.AttachmentInterface { return new(schemes.FileAttachment) },
	schemes.AttachmentImage:    func() schemes.AttachmentInterface { return new(schemes.PhotoAttachment) },
	schemes.AttachmentKeyboard: func() schemes.AttachmentInterface { return new(schemes.InlineKeyboardAttachment) },
	schemes.AttachmentLocation: func() schemes.AttachmentInterface { return new(schemes.LocationAttachment) },
	schemes.AttachmentShare:    func() schemes.AttachmentInterface { return new(schemes.ShareAttachment) },
	schemes.AttachmentSticker:  func() schemes.AttachmentInterface { return new(schemes.StickerAttachment) },
	schemes.AttachmentVideo:    func() schemes.AttachmentInterface { return new(schemes.VideoAttachment) },
}

// bytesToProperUpdate converts raw JSON bytes to the appropriate update type
func (a *Api) bytesToProperUpdate(data []byte) (schemes.UpdateInterface, error) {
	baseUpdate := &schemes.Update{}
	if err := json.Unmarshal(data, baseUpdate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base update: %w", err)
	}

	debugRaw := ""
	if a.debug {
		debugRaw = string(data)
	}

	updateType := baseUpdate.GetUpdateType()
	constructor, exists := updateTypeMap[updateType]
	if !exists {
		return nil, fmt.Errorf("unknown update type: %s", updateType)
	}

	update := constructor(debugRaw)
	if err := json.Unmarshal(data, update); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update of type %s: %w", updateType, err)
	}

	if err := a.processMessageAttachments(update); err != nil {
		return nil, fmt.Errorf("failed to process message attachments: %w", err)
	}

	return update, nil
}

// processMessageAttachments processes attachments for message-type updates
func (a *Api) processMessageAttachments(update schemes.UpdateInterface) error {
	switch u := update.(type) {
	case *schemes.MessageCreatedUpdate:
		if u.Message.Body.RawAttachments != nil {
			for _, rawAttachment := range u.Message.Body.RawAttachments {
				attachment, err := a.bytesToProperAttachment([]byte(rawAttachment))
				if err != nil {
					return fmt.Errorf("failed to process attachment: %w", err)
				}

				u.Message.Body.Attachments = append(u.Message.Body.Attachments, attachment)
			}
		}
	case *schemes.MessageEditedUpdate:
		if u.Message.Body.RawAttachments != nil {
			for _, rawAttachment := range u.Message.Body.RawAttachments {
				attachment, err := a.bytesToProperAttachment([]byte(rawAttachment))
				if err != nil {
					return fmt.Errorf("failed to process attachment: %w", err)
				}

				u.Message.Body.Attachments = append(u.Message.Body.Attachments, attachment)
			}
		}
	default:
		return nil // No attachments to process
	}

	return nil
}

// bytesToProperAttachment converts raw JSON bytes to the appropriate attachment type
func (a *Api) bytesToProperAttachment(data []byte) (schemes.AttachmentInterface, error) {
	baseAttachment := &schemes.Attachment{}
	if err := json.Unmarshal(data, baseAttachment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base attachment: %w", err)
	}

	attachmentType := baseAttachment.GetAttachmentType()
	constructor, exists := attachmentTypeMap[attachmentType]
	if !exists {
		// Return base attachment for unknown types
		return baseAttachment, nil
	}

	attachment := constructor()
	if err := json.Unmarshal(data, attachment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attachment of type %s: %w", attachmentType, err)
	}

	return attachment, nil
}

// UpdatesParams holds parameters for getting updates
type UpdatesParams struct {
	Limit   int
	Timeout time.Duration
	Marker  int64
	Types   []string
}

// getUpdates fetches updates from the API
func (a *Api) getUpdates(ctx context.Context, params *UpdatesParams) (*schemes.UpdateList, error) {
	if params == nil {
		params = &UpdatesParams{}
	}

	values := url.Values{}

	if params.Limit > 0 {
		values.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Timeout > 0 {
		values.Set("timeout", strconv.Itoa(int(params.Timeout.Seconds())))
	}
	if params.Marker > 0 {
		values.Set("marker", strconv.FormatInt(params.Marker, 10))
	}
	for _, t := range params.Types {
		values.Add("types", t)
	}

	body, err := a.client.request(ctx, http.MethodGet, "updates", values, false, nil)
	if err != nil {
		var te *TimeoutError
		// Обрабатывать timeout как пустую страницу (ожидается при длительном опросе)
		if errors.As(err, &te) {
			return &schemes.UpdateList{}, nil
		}
		// Учитывать отмену контекста/истечение срока
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}

	defer func() {
		if closeErr := body.Close(); closeErr != nil {
			log.Printf("failed to close response body: %v", closeErr)
		}
	}()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	result := &schemes.UpdateList{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updates: %w", err)
	}

	return result, nil
}

func (a *Api) getUpdatesWithRetry(ctx context.Context, params *UpdatesParams) (*schemes.UpdateList, error) {
	if params == nil {
		params = &UpdatesParams{}
	}

	var result *schemes.UpdateList
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		result, lastErr = a.getUpdates(ctx, params)
		if lastErr == nil {
			return result, nil
		}

		// Остановить retry если context отменен/истек
		if errors.Is(lastErr, context.Canceled) || errors.Is(lastErr, context.DeadlineExceeded) {
			return nil, lastErr
		}

		if attempt < maxRetries-1 {
			retryWait := time.Duration(1<<uint(attempt)) * time.Second
			log.Printf("Attempt %d failed, retrying in %v: %v", attempt+1, retryWait, lastErr)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryWait):
			}
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// GetUpdates returns a channel that delivers updates from the API
func (a *Api) GetUpdates(ctx context.Context) <-chan schemes.UpdateInterface {
	ch := make(chan schemes.UpdateInterface, 100)

	go func() {
		defer close(ch)

		var marker int64
		ticker := time.NewTicker(a.pause)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for {
					params := &UpdatesParams{
						Limit:   maxUpdatesLimit,
						Timeout: a.timeout,
						Marker:  marker,
					}

					updateList, err := a.getUpdatesWithRetry(ctx, params)
					if err != nil {
						log.Printf("failed to get updates: %v", err)
						break
					}

					if len(updateList.Updates) == 0 {
						break
					}

					for _, rawUpdate := range updateList.Updates {
						update, err := a.bytesToProperUpdate(rawUpdate)
						if err != nil {
							continue
						}

						select {
						case ch <- update:
						case <-ctx.Done():
							return
						}
					}

					if updateList.Marker != nil {
						marker = *updateList.Marker
					}
				}
			}
		}
	}()

	return ch
}

// GetHandler returns an http.HandlerFunc for webhook handling
func (a *Api) GetHandler(updates chan<- schemes.UpdateInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		update, err := a.bytesToProperUpdate(body)
		if err != nil {
			http.Error(w, "Failed to parse update", http.StatusBadRequest)
			return
		}

		select {
		case updates <- update:
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Updates channel is full", http.StatusServiceUnavailable)
		}
	}
}
