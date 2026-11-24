package maxbot

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/v3/model"
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

// WithUpdateTypes ограничивает набор апдейтов, получаемых при поллинге (query-параметр types у GET /updates).
func WithUpdateTypes(types ...model.UpdateType) Opt {
	return func(c *client) error {
		c.updateTypes = types

		return nil
	}
}

// QueryOpt — необязательный query-параметр запроса.
type QueryOpt func(url.Values)

// WithMessageIDs — список идентификаторов сообщений (query-параметр message_ids).
func WithMessageIDs(ids []string) QueryOpt {
	return func(values url.Values) {
		if len(ids) > 0 {
			values.Set(paramMessageIDs, strings.Join(ids, ","))
		}
	}
}

// WithCount — максимальное число сообщений в ответе (query-параметр count).
func WithCount(v int64) QueryOpt {
	return func(values url.Values) {
		if v > 0 {
			values.Set(paramCount, strconv.FormatInt(v, 10))
		}
	}
}

// WithBefore — сообщения до указанной метки времени (query-параметр before).
func WithBefore(v int64) QueryOpt {
	return func(values url.Values) {
		values.Set(paramBefore, strconv.FormatInt(v, 10))
	}
}

// WithAfter — сообщения после указанной метки времени (query-параметр after).
func WithAfter(v int64) QueryOpt {
	return func(values url.Values) {
		values.Set(paramAfter, strconv.FormatInt(v, 10))
	}
}

// WithDisableLinkPreview отключает превью ссылок (query-параметр disable_link_preview).
func WithDisableLinkPreview() QueryOpt {
	return func(values url.Values) {
		values.Set(paramDisableLinkPreview, "true")
	}
}
