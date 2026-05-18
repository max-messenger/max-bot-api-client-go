package maxbot

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Code    string `json:"code"`
	Err     string `json:"error,omitempty"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Code + " " + e.Err + ": " + e.Message
}

func (e Error) IsAttachmentNotReady() bool {
	return e.Code == "attachment.not.ready"
}

func parseResponseError(resp *http.Response) error {
	responseErr := &Error{}

	err := json.NewDecoder(resp.Body).Decode(responseErr)
	if err != nil {
		return fmt.Errorf("parse response error: %w", err)
	}

	return responseErr
}

type TimeoutError struct {
	Op     string
	Reason string
}

func (e *TimeoutError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("timeout error during %s: %s", e.Op, e.Reason)
	}

	return fmt.Sprintf("timeout error during %s", e.Op)
}

type NetworkError struct {
	Op  string
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error during %s: %v", e.Op, e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}
