package maxbot

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyToken = errors.New("bot token is empty")
	ErrInvalidURL = errors.New("invalid API URL")
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error %d: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

func (e *APIError) Is(target error) bool {
	var t *APIError
	if errors.As(target, &t) {
		return e.Code == t.Code
	}

	return false
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

func (e *TimeoutError) Timeout() bool {
	return true
}

type SerializationError struct {
	Op   string
	Type string
	Err  error
}

func (e *SerializationError) Error() string {
	return fmt.Sprintf("serialization error during %s of %s: %v", e.Op, e.Type, e.Err)
}

func (e *SerializationError) Unwrap() error {
	return e.Err
}
