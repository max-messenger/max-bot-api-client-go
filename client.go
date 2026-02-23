package maxbot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type client struct {
	key        string
	version    string
	baseURL    *url.URL
	httpClient *http.Client
	errors     chan error
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
		errors:     make(chan error, 1),
	}
}

func (cl *client) notifyError(err error) {
	if err == nil {
		return
	}
	select {
	case cl.errors <- err:
	default:
		log.Println(err)
	}
}

func (cl *client) closer(name string, c io.Closer) {
	if c == nil {
		return
	}
	if err := c.Close(); err != nil {
		cl.notifyError(fmt.Errorf("failed to close %s: %w", name, err))
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

	data, err := jsoniter.Marshal(body)
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

	query.Set(paramVersion, cl.version)
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

	resp, err := cl.do(req)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
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
		defer cl.closer("requestReader body", resp.Body)

		apiErr := &schemes.Error{}
		if decodeErr := jsoniter.NewDecoder(resp.Body).Decode(apiErr); decodeErr != nil {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		}

		return nil, &APIError{
			Code:    resp.StatusCode,
			Message: apiErr.Code,
			Details: apiErr.Message,
		}
	}

	return resp.Body, nil
}

func (cl *client) do(req *http.Request) (*http.Response, error) {
	return cl.httpClient.Do(req)
}

// Close closes the HTTP client.
func (cl *client) Close() error {
	if transport, ok := cl.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}

	return nil
}
