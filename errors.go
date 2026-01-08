package maxbot

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyToken = errors.New("bot token is empty")
	ErrInvalidURL = errors.New("invalid API URL")
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`

	// Extra context (filled by client; not part of API schema)
	HTTPStatus string `json:"-"`
	Method     string `json:"-"`
	URL        string `json:"-"`
	RawBody    string `json:"-"`
}

func (e *APIError) Error() string {
	// Human-friendly error message that keeps MAX API details.
	base := ""
	if e.HTTPStatus != "" {
		base = fmt.Sprintf("HTTP %d: %s", e.Code, e.HTTPStatus)
	} else {
		base = fmt.Sprintf("HTTP %d", e.Code)
	}

	msg := e.Message
	if msg == "" {
		msg = "request failed"
	}

	parts := []string{base, msg}

	if e.Details != "" {
		parts = append(parts, "("+e.Details+")")
	}

	if e.Method != "" || e.URL != "" {
		mu := strings.TrimSpace(strings.TrimSpace(e.Method) + " " + strings.TrimSpace(e.URL))
		if mu != "" {
			parts = append(parts, "["+mu+"]")
		}
	}

	if e.RawBody != "" {
		raw := strings.TrimSpace(e.RawBody)
		raw = strings.ReplaceAll(raw, "\n", " ")
		raw = strings.ReplaceAll(raw, "\r", " ")
		if len(raw) > 500 {
			raw = raw[:500] + "…"
		}
		parts = append(parts, "body="+raw)
	}

	return strings.Join(parts, " ")
}

func (e *APIError) Is(target error) bool {
	if t, ok := target.(*APIError); ok {
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
