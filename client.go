package maxbot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type client struct {
	token       string
	baseURL     url.URL
	httpClient  HttpClient
	pollPause   time.Duration
	pollTimeout time.Duration
}

func newClient(token, host string) *client {
	return &client{
		token: token,
		baseURL: url.URL{
			Scheme: defaultScheme,
			Host:   host,
		},
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		pollPause:   defaultPause,
		pollTimeout: defaultTimeout,
	}
}

func (c *client) rawWithRetry(ctx context.Context, method, path string, query url.Values, in, out any) error {
	var err error
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = c.raw(ctx, method, path, query, in, out)
		if err == nil {
			return nil
		}

		apiErr := &Error{}
		if errors.As(err, &apiErr) && !apiErr.IsAttachmentNotReady() {
			return fmt.Errorf("sending message failed: %w", err)
		}

		retryWait := time.Duration(1<<uint(attempt)) * time.Second
		if attempt < maxRetries-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryWait):
			}
		}
	}

	return err
}

func (c *client) raw(ctx context.Context, method, path string, query url.Values, in, out any) error {
	u := c.baseURL
	u.Path = path

	u.RawQuery = query.Encode()

	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(AuthorizationHeader, c.token)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			if urlErr.Timeout() {
				return &TimeoutError{
					fmt.Sprintf("%s %s", method, path),
					"request timeout exceeded",
				}
			}
		}

		return &NetworkError{
			Op:  fmt.Sprintf("%s %s", method, path),
			Err: err,
		}
	}

	defer func() { _ = resp.Body.Close() }()
	if c.isNotOk(resp.StatusCode) {
		return parseResponseError(resp)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}

func (c *client) do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *client) isNotOk(statusCode int) bool {
	return statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices
}
