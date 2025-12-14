package maxbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	errLongPollTimeout = &TimeoutError{
		Op:     "long polling",
		Reason: "request timeout exceeded",
	}
)

type client struct {
	key        string
	version    string
	baseURL    *url.URL
	httpClient *http.Client
}

func newClient(key string, version string, baseURL *url.URL, httpClient *http.Client) *client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	return &client{
		key:        key,
		version:    version,
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (cl *client) createTimeoutError(op string, reason string) *TimeoutError {
	return &TimeoutError{
		Op:     op,
		Reason: reason,
	}
}

func (cl *client) request(ctx context.Context, method, path string, query url.Values, reset bool, body interface{}) (io.ReadCloser, error) {
	if body == nil {
		return cl.requestReader(ctx, method, path, query, reset, nil)
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, &SerializationError{
			Op:   "marshal",
			Type: "request body",
			Err:  err,
		}
	}

	return cl.requestReader(ctx, method, path, query, reset, bytes.NewReader(data))
}

func (cl *client) requestReader(ctx context.Context, method, path string, query url.Values, reset bool, body io.Reader) (io.ReadCloser, error) {
	if query == nil {
		query = url.Values{}
	}

	u := *cl.baseURL
	u.Path = path

	/* DEPRECATED
	if !reset {
		query.Set("access_token", cl.key)
	}
	*/

	query.Set("v", cl.version)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", fmt.Sprintf("max-bot-api-client-go/%s", cl.version))
	if !reset {
		req.Header.Set("Authorization", cl.key)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() {
				return nil, cl.createTimeoutError(
					fmt.Sprintf("%s %s", method, path),
					fmt.Sprintf("request timeout exceeded (%v)", cl.httpClient.Timeout),
				)
			}
		}

		return nil, &NetworkError{
			Op:  fmt.Sprintf("%s %s", method, path),
			Err: err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		// Read body into buffer to decode error
		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if readErr != nil {
			return nil, fmt.Errorf("HTTP %d: failed to read error body: %w", resp.StatusCode, readErr)
		}

		// Use simplified structure as schemes.Error.Message expects Message type, not string
		type SimpleError struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}

		simpleErr := &SimpleError{}
		if decodeErr := json.Unmarshal(bodyBytes, simpleErr); decodeErr != nil {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		}

		return nil, &APIError{
			Code:    resp.StatusCode,
			Message: fmt.Sprintf("%s: %s", simpleErr.Code, simpleErr.Message),
		}
	}

	return resp.Body, nil
}

// Close closes the HTTP client
func (cl *client) Close() error {
	if transport, ok := cl.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}

	return nil
}
