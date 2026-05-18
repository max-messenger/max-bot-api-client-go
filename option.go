package maxbot

import (
	"net/url"
	"time"
)

type Opt func(cli *client) error

func WithHTTPClient(cli HttpClient) Opt {
	return func(c *client) error {
		c.httpClient = cli

		return nil
	}
}

func WithBaseURL(baseURL string) Opt {
	u, err := url.Parse(baseURL)

	return func(c *client) error {
		if err != nil {
			return err
		}
		c.baseURL.Host = u.Host
		c.baseURL.Scheme = u.Scheme

		return nil
	}
}

func WithPollingPause(d time.Duration) Opt {
	return func(c *client) error {
		c.pollPause = d

		return nil
	}
}

func WithPollingTimeout(d time.Duration) Opt {
	return func(c *client) error {
		c.pollTimeout = d

		return nil
	}
}
