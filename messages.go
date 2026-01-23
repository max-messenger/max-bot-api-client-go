package maxbot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type messages struct {
	client *client
}

func newMessages(client *client) *messages {
	return &messages{client: client}
}

// MessageResponse represents the response wrapper when a message is sent
type MessageResponse struct {
	Message schemes.Message `json:"message"`
}

// isAttachmentNotReadyError checks whether the error is "attachment.not.ready"
// According to the MAX API documentation, files are processed on the server after downloading
// and may not be available immediately for sending in a message
func (a *messages) isAttachmentNotReadyError(err error) bool {
	if err == nil {
		return false
	}

	// Check if error is typed APIError (from client.go)
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return strings.Contains(apiErr.Message, "attachment.not.ready") ||
			strings.Contains(apiErr.Message, "attachment.file.not.processed") ||
			strings.Contains(apiErr.Message, "errors.process.attachment.file.not.processed")
	}

	// Check if error is typed schemes.Error
	var schemesErr *schemes.Error
	if errors.As(err, &schemesErr) {
		return schemesErr.Code == "attachment.not.ready" ||
			schemesErr.Code == "attachment.file.not.processed" ||
			schemesErr.Code == "errors.process.attachment.file.not.processed"
	}

	// Fallback to string matching for non-typed errors
	errText := err.Error()
	return strings.Contains(errText, "attachment.not.ready") ||
		strings.Contains(errText, "attachment.file.not.processed") ||
		strings.Contains(errText, "errors.process.attachment.file.not.processed")
}

// GetMessages returns messages in chat: result page and marker referencing to the next page. Messages traversed in reverse direction so the latest message in chat will be first in result array. Therefore if you use from and to parameters, to must be less than from
func (a *messages) GetMessages(ctx context.Context, chatID int64, messageIDs []string, from int, to int, count int) (*schemes.MessageList, error) {
	result := new(schemes.MessageList)
	values := url.Values{}
	if chatID != 0 {
		values.Set("chat_id", strconv.Itoa(int(chatID)))
	}
	if len(messageIDs) > 0 {
		for _, mid := range messageIDs {
			values.Add("message_ids", mid)
		}
	}
	if from != 0 {
		values.Set("from", strconv.Itoa(from))
	}
	if to != 0 {
		values.Set("to", strconv.Itoa(to))
	}
	if count > 0 {
		values.Set("count", strconv.Itoa(count))
	}
	body, err := a.client.request(ctx, http.MethodGet, "messages", values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

func (a *messages) GetMessage(ctx context.Context, messageID string) (*schemes.Message, error) {
	result := new(schemes.Message)
	path := "messages/" + url.PathEscape(messageID)
	body, err := a.client.request(ctx, http.MethodGet, path, nil, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// EditMessage updates message by id
// Returns error if the operation fails
func (a *messages) EditMessage(ctx context.Context, messageID string, message *Message) error {
	_, err := a.editMessage(ctx, messageID, message.message)
	return err
}

// DeleteMessage deletes message by id
func (a *messages) DeleteMessage(ctx context.Context, messageID string) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	values.Set("message_id", messageID)
	body, err := a.client.request(ctx, http.MethodDelete, "messages", values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// AnswerOnCallback should be called to send an answer after a user has clicked the button. The answer may be an updated message or/and a one-time user notification.
func (a *messages) AnswerOnCallback(ctx context.Context, callbackID string, callback *schemes.CallbackAnswer) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	values.Set("callback_id", callbackID)
	body, err := a.client.request(ctx, http.MethodPost, "answers", values, false, callback)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// NewKeyboardBuilder returns new keyboard builder helper
func (a *messages) NewKeyboardBuilder() *Keyboard {
	return &Keyboard{
		rows: make([]*KeyboardRow, 0),
	}
}

// Send sends a message to a chat. As a result for this method new message identifier returns.
func (a *messages) Send(ctx context.Context, m *Message) error {
	_, err := a.sendMessage(ctx, m.reset, m.chatID, m.userID, m.message)
	return err
}

// SendWithResult sends a message to a chat and returns the created message along with any error.
func (a *messages) SendWithResult(ctx context.Context, m *Message) (*schemes.Message, error) {
	return a.sendMessage(ctx, m.reset, m.chatID, m.userID, m.message)
}

// sendMessage sends a message with automatic retry on attachment.not.ready error
// According to the MAX API documentation (https://dev.max.ru/), after uploading a file
// the server processes it for some time, and sending a message immediately may result in an error.
// This method automatically retries with exponential backoff.
func (a *messages) sendMessage(ctx context.Context, reset bool, chatID int64, userID int64, message *schemes.NewMessageBody) (*schemes.Message, error) {
	values := url.Values{}
	if chatID != 0 {
		values.Set("chat_id", strconv.Itoa(int(chatID)))
	}
	if userID != 0 {
		values.Set("user_id", strconv.Itoa(int(userID)))
	}

	var lastErr error
	retryDelay := schemes.InitialRetryDelay

	// Main attempt + retry on attachment.not.ready
	for attempt := 0; attempt <= schemes.MaxAttachmentRetries; attempt++ {
		body, err := a.client.request(ctx, http.MethodPost, "messages", values, reset, message)

		// Handle error case
		if err != nil {
			// body might still be available even on error - try to read it
			if body != nil {
				// Try to decode as SimpleQueryResult first
				result := new(schemes.SimpleQueryResult)
				if decodeErr := json.NewDecoder(body).Decode(result); decodeErr == nil && !result.Success {
					body.Close()
					lastErr = errors.New(result.Message)

					// Check if error message contains attachment.not.ready
					if a.isAttachmentNotReadyError(lastErr) && attempt < schemes.MaxAttachmentRetries {
						// Wait before next attempt
						select {
						case <-ctx.Done():
							return nil, fmt.Errorf("context cancelled while waiting for attachment: %w", ctx.Err())
						case <-time.After(retryDelay):
							// Increase delay exponentially (300ms -> 600ms -> 1200ms -> 2400ms -> 3000ms)
							retryDelay *= 2
							if retryDelay > schemes.MaxRetryDelay {
								retryDelay = schemes.MaxRetryDelay
							}
							continue
						}
					}
					// Not attachment error - return immediately
					return nil, lastErr
				}
				body.Close()
			}

			// Check if the error itself contains attachment.not.ready
			lastErr = err
			if a.isAttachmentNotReadyError(err) && attempt < schemes.MaxAttachmentRetries {
				// Wait before next attempt
				select {
				case <-ctx.Done():
					return nil, fmt.Errorf("context cancelled while waiting for attachment: %w", ctx.Err())
				case <-time.After(retryDelay):
					// Increase delay exponentially (300ms -> 600ms -> 1200ms -> 2400ms -> 3000ms)
					retryDelay *= 2
					if retryDelay > schemes.MaxRetryDelay {
						retryDelay = schemes.MaxRetryDelay
					}
					continue
				}
			}

			// For other errors or when attempts exhausted - return error
			return nil, err
		}

		// Successful HTTP request - decode response
		defer func() {
			if err := body.Close(); err != nil {
				slog.Error("failed to close response body", "error", err)
			}
		}()

		wrapper := new(MessageResponse)
		if err := json.NewDecoder(body).Decode(wrapper); err != nil {
			return nil, err
		}

		return &wrapper.Message, nil
	}

	return nil, fmt.Errorf("failed to send message after %d retries: %w", schemes.MaxAttachmentRetries, lastErr)
}

// editMessage updates message by id with automatic retry on attachment.not.ready error
// Note: MAX API has inconsistent behavior - Send returns HTTP 400 on attachment.not.ready,
// while EditMessage returns HTTP 200 with success=false
func (a *messages) editMessage(ctx context.Context, messageID string, message *schemes.NewMessageBody) (*schemes.SimpleQueryResult, error) {
	values := url.Values{}
	values.Set("message_id", messageID)

	var lastErr error
	retryDelay := schemes.InitialRetryDelay

	// Main attempt + retry on attachment.not.ready
	for attempt := 0; attempt <= schemes.MaxAttachmentRetries; attempt++ {
		body, err := a.client.request(ctx, http.MethodPut, "messages", values, false, message)

		// Handle error case
		if err != nil {
			// Check if the error itself contains attachment.not.ready
			lastErr = err
			if a.isAttachmentNotReadyError(err) && attempt < schemes.MaxAttachmentRetries {
				// Wait before next attempt
				select {
				case <-ctx.Done():
					return nil, fmt.Errorf("context cancelled while waiting for attachment: %w", ctx.Err())
				case <-time.After(retryDelay):
					// Increase delay exponentially
					retryDelay *= 2
					if retryDelay > schemes.MaxRetryDelay {
						retryDelay = schemes.MaxRetryDelay
					}
					continue
				}
			}

			// For other errors or when attempts exhausted - return error
			return nil, err
		}

		// Successful HTTP request - decode response
		defer func() {
			if err := body.Close(); err != nil {
				slog.Error("failed to close response body", "error", err)
			}
		}()

		result := new(schemes.SimpleQueryResult)
		if err := json.NewDecoder(body).Decode(result); err != nil {
			return nil, err
		}

		// If response indicates failure, check if it's attachment.not.ready
		if !result.Success {
			err := errors.New(result.Message)

			// Check if error is attachment.not.ready
			if a.isAttachmentNotReadyError(err) && attempt < schemes.MaxAttachmentRetries {
				// Wait before next attempt
				select {
				case <-ctx.Done():
					return nil, fmt.Errorf("context cancelled while waiting for attachment: %w", ctx.Err())
				case <-time.After(retryDelay):
					// Increase delay exponentially
					retryDelay *= 2
					if retryDelay > schemes.MaxRetryDelay {
						retryDelay = schemes.MaxRetryDelay
					}
					continue
				}
			}

			// Not attachment error - return immediately
			return result, err
		}

		return result, nil
	}

	return nil, fmt.Errorf("failed to edit message after %d retries: %w", schemes.MaxAttachmentRetries, lastErr)
}

// Check posiable to send a message to a chat.
func (a *messages) Check(ctx context.Context, m *Message) (bool, error) {
	return a.checkUser(ctx, m.reset, m.message)
}

func (a *messages) checkUser(ctx context.Context, reset bool, message *schemes.NewMessageBody) (bool, error) {
	result := new(schemes.Error)
	values := url.Values{}
	if reset {
		values.Set("access_token", message.BotToken)
	}
	mode := "notify/exists"

	if message.PhoneNumbers != nil {
		values.Set("phone_numbers", strings.Join(message.PhoneNumbers, ","))
	}

	body, err := a.client.request(ctx, http.MethodGet, mode, values, reset, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if err := json.NewDecoder(body).Decode(result); err != nil {
		return false, err
	}

	if len(result.NumberExist) > 0 {
		return true, result
	}

	return false, result
}

// Check posiable to send a message to a chat.
func (a *messages) ListExist(ctx context.Context, m *Message) ([]string, error) {
	return a.checkNumberExist(ctx, m.reset, m.message)
}

func (a *messages) checkNumberExist(ctx context.Context, reset bool, message *schemes.NewMessageBody) ([]string, error) {
	result := new(schemes.Error)
	values := url.Values{}
	if reset {
		values.Set("access_token", message.BotToken)
	}
	mode := "notify/exists"

	if message.PhoneNumbers != nil {
		values.Set("phone_numbers", strings.Join(message.PhoneNumbers, ","))
	}

	body, err := a.client.request(ctx, http.MethodGet, mode, values, reset, nil)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	if err := json.NewDecoder(body).Decode(result); err != nil {
		// Message sent without errors
		return nil, err
	}
	if len(result.NumberExist) > 0 {
		return result.NumberExist, result
	}
	return nil, result
}
