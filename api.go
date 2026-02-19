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

// Api represents the MAX Bot API client.
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

// New creates a new Max Bot API client with the provided token.
func New(token string, opts ...Option) (*Api, error) {
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

	for _, o := range opts {
		o(api)
	}

	return api, nil
}

// NewWithConfig creates a new Max Bot API client from the configuration service.
func NewWithConfig(cfg configservice.ConfigInterface) (*Api, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	token := cfg.BotTokenCheckString()
	if token == "" {
		token = os.Getenv(envToken)
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

func getUpdateType(updateType schemes.UpdateType) func(debugRaw string) schemes.UpdateInterface {
	switch updateType {
	case schemes.TypeMessageCallback:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.MessageCallbackUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	case schemes.TypeMessageCreated:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.MessageCreatedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	case schemes.TypeMessageRemoved:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.MessageRemovedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	case schemes.TypeMessageEdited:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.MessageEditedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	case schemes.TypeBotAdded:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.BotAddedToChatUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	case schemes.TypeBotRemoved:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.BotRemovedFromChatUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	case schemes.TypeUserAdded:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.UserAddedToChatUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	case schemes.TypeUserRemoved:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.UserRemovedFromChatUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	case schemes.TypeBotStarted:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.BotStartedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	case schemes.TypeChatTitleChanged:
		return func(debugRaw string) schemes.UpdateInterface {
			return &schemes.ChatTitleChangedUpdate{Update: schemes.Update{DebugRaw: debugRaw}}
		}
	}

	return nil
}

func getAttachmentType(attachmentType schemes.AttachmentType) func() schemes.AttachmentInterface {
	switch attachmentType {
	case schemes.AttachmentAudio:
		return func() schemes.AttachmentInterface { return new(schemes.AudioAttachment) }
	case schemes.AttachmentContact:
		return func() schemes.AttachmentInterface { return new(schemes.ContactAttachment) }
	case schemes.AttachmentFile:
		return func() schemes.AttachmentInterface { return new(schemes.FileAttachment) }
	case schemes.AttachmentImage:
		return func() schemes.AttachmentInterface { return new(schemes.PhotoAttachment) }
	case schemes.AttachmentKeyboard:
		return func() schemes.AttachmentInterface { return new(schemes.InlineKeyboardAttachment) }
	case schemes.AttachmentLocation:
		return func() schemes.AttachmentInterface { return new(schemes.LocationAttachment) }
	case schemes.AttachmentShare:
		return func() schemes.AttachmentInterface { return new(schemes.ShareAttachment) }
	case schemes.AttachmentSticker:
		return func() schemes.AttachmentInterface { return new(schemes.StickerAttachment) }
	case schemes.AttachmentVideo:
		return func() schemes.AttachmentInterface { return new(schemes.VideoAttachment) }
	}

	return nil
}

// bytesToProperUpdate converts raw JSON bytes to the appropriate update type.
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
	constructor := getUpdateType(updateType)
	if constructor == nil {
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

// processMessageAttachments processes attachments for message-type updates.
func (a *Api) processMessageAttachments(update schemes.UpdateInterface) error {
	switch u := update.(type) {
	case *schemes.MessageCreatedUpdate:
		if u.Message.Body.RawAttachments != nil {
			attachments, err := a.convertRawAttachments(u.Message.Body.RawAttachments)
			if err != nil {
				return err
			}

			u.Message.Body.Attachments = append(u.Message.Body.Attachments, attachments...)
		} else if u.Message.Link != nil && u.Message.Link.Message.RawAttachments != nil {
			attachments, err := a.convertRawAttachments(u.Message.Link.Message.RawAttachments)
			if err != nil {
				return err
			}

			u.Message.Link.Message.Attachments = append(u.Message.Link.Message.Attachments, attachments...)
		}
	case *schemes.MessageEditedUpdate:
		if u.Message.Body.RawAttachments != nil {
			attachments, err := a.convertRawAttachments(u.Message.Body.RawAttachments)
			if err != nil {
				return err
			}

			u.Message.Body.Attachments = append(u.Message.Body.Attachments, attachments...)
		} else if u.Message.Link != nil && u.Message.Link.Message.RawAttachments != nil {
			attachments, err := a.convertRawAttachments(u.Message.Link.Message.RawAttachments)
			if err != nil {
				return err
			}

			u.Message.Link.Message.Attachments = append(u.Message.Link.Message.Attachments, attachments...)
		}
	default:
		return nil // No attachments to process
	}

	return nil
}

// bytesToProperAttachment converts raw JSON bytes to the appropriate attachment type.
func (a *Api) bytesToProperAttachment(data []byte) (schemes.AttachmentInterface, error) {
	baseAttachment := &schemes.Attachment{}
	if err := json.Unmarshal(data, baseAttachment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base attachment: %w", err)
	}

	attachmentType := baseAttachment.GetAttachmentType()
	constructor := getAttachmentType(attachmentType)
	if constructor == nil {
		// Return base attachment for unknown types
		return baseAttachment, nil
	}

	attachment := constructor()
	if err := json.Unmarshal(data, attachment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attachment of type %s: %w", attachmentType, err)
	}

	return attachment, nil
}

func (a *Api) convertRawAttachments(rawAttachments []json.RawMessage) ([]any, error) {
	result := make([]any, 0, len(rawAttachments))
	for _, rawAttachment := range rawAttachments {
		attachment, err := a.bytesToProperAttachment([]byte(rawAttachment))
		if err != nil {
			return nil, fmt.Errorf("failed to process attachment: %w", err)
		}

		result = append(result, attachment)
	}

	return result, nil
}

// UpdatesParams holds parameters for getting updates.
type UpdatesParams struct {
	Limit   int
	Timeout time.Duration
	Marker  int64
	Types   []string
}

// getUpdates fetches updates from the API.
func (a *Api) getUpdates(ctx context.Context, params *UpdatesParams) (*schemes.UpdateList, error) {
	if params == nil {
		params = &UpdatesParams{}
	}

	values := url.Values{}

	if params.Limit > 0 {
		values.Set(paramLimit, strconv.Itoa(params.Limit))
	}
	if params.Timeout > 0 {
		values.Set(paramTimeout, strconv.Itoa(int(params.Timeout.Seconds())))
	}
	if params.Marker > 0 {
		values.Set(paramMarker, strconv.FormatInt(params.Marker, 10))
	}
	for _, t := range params.Types {
		values.Add(paramTypes, t)
	}

	body, err := a.client.request(ctx, http.MethodGet, pathUpdates, values, false, nil)
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

// GetUpdates returns a channel that delivers updates from the API.
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

// GetHandler returns an http.HandlerFunc for webhook handling.
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
