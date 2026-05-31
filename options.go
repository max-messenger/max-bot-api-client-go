package maxbot

import (
	"net/url"
	"time"
)

type Option func(api *Api)

func WithBaseURL(baseURL string) Option {
	return func(api *Api) {
		u, err := url.Parse(baseURL)
		if err != nil {
			return
		}

		api.client.baseURL = u
	}
}

func WithHTTPClient(httpClient HttpClient) Option {
	return func(api *Api) {
		api.client.httpClient = httpClient
	}
}

func WithApiTimeout(timeout time.Duration) Option {
	return func(api *Api) {
		api.timeout = timeout
	}
}

func WithPauseTimeout(timeout time.Duration) Option {
	return func(api *Api) {
		api.pause = timeout
	}
}

func WithVersion(version string) Option {
	return func(api *Api) {
		api.client.version = version
	}
}

func WithDebugMode() Option {
	return func(api *Api) {
		api.debug = true
	}
}

func WithDebugChat(chat int64) Option {
	return func(api *Api) {
		api.Debugs.chat = chat
	}
}

func WithUpdateHandler(f UpdateHandler) Option {
	return func(api *Api) {
		api.updateHandler = f
	}
}

// WithErrorBufferSize задаёт ёмкость канала ошибок, возвращаемого GetErrors().
// Больший буфер снижает риск потери ошибок при всплесках (повторные сбои long
// polling и т.п.). Значения меньше 1 игнорируются.
func WithErrorBufferSize(size int) Option {
	return func(api *Api) {
		if size < 1 {
			return
		}
		api.client.errors = make(chan error, size)
	}
}
