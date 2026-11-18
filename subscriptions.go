package maxbot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type UpdateHandler func(context.Context, model.Update)

type Subscriptions struct {
	client      *client
	pause       time.Duration
	timeout     time.Duration
	updateTypes []model.UpdateType
}

func newSubscriptions(client *client) *Subscriptions {
	return &Subscriptions{
		client:      client,
		pause:       client.pollPause,
		timeout:     client.pollTimeout,
		updateTypes: client.updateTypes,
	}
}

func (s *Subscriptions) GetSubscriptions(ctx context.Context) (res model.GetSubscriptionsResult, err error) {
	err = s.client.raw(ctx, http.MethodGet, pathSubscriptions, nil, nil, &res)

	return
}

func (s *Subscriptions) Subscribe(ctx context.Context, u, st string, ut []string) (res model.SimpleQueryResult, err error) {
	data := model.SubscriptionRequestBody{
		URL:         u,
		Secret:      st,
		UpdateTypes: ut,
	}
	err = s.client.raw(ctx, http.MethodPost, pathSubscriptions, nil, data, &res)

	return
}

func (s *Subscriptions) Unsubscribe(ctx context.Context, u string) (res model.SimpleQueryResult, err error) {
	values := url.Values{}
	values.Add(paramURL, u)
	err = s.client.raw(ctx, http.MethodDelete, pathSubscriptions, values, nil, &res)

	return
}

// GetUpdates returns a list of updates from the API.
func (s *Subscriptions) GetUpdates(ctx context.Context, marker int64) ([]model.Update, int64, error) {
	res := make([]model.Update, 0)

	result, err := s.getUpdatesWithRetry(ctx, maxUpdatesLimit, int(s.timeout.Seconds()), marker)
	if err != nil {
		return nil, 0, err
	}

	if len(result.Updates) == 0 {
		return res, result.Marker, nil
	}

	for _, rawUpdate := range result.Updates {
		res = append(res, rawUpdate.FromRaw())
	}

	return res, result.Marker, nil
}

func (s *Subscriptions) getUpdatesWithRetry(ctx context.Context, limit, timeout int, marker int64) (res updateList, err error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		res, err = s.getUpdates(ctx, limit, timeout, marker)
		if err == nil {
			return
		}

		// Остановить retry если context отменен/истек
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return
		}

		if attempt < maxRetries-1 {
			retryWait := time.Duration(1<<uint(attempt)) * time.Second
			select {
			case <-ctx.Done():
				err = ctx.Err()

				return
			case <-time.After(retryWait):
			}
		}
	}

	err = fmt.Errorf("failed after %d attempts: %w", maxRetries, err)

	return
}

func (s *Subscriptions) getUpdates(ctx context.Context, limit, timeout int, marker int64) (res updateList, err error) {
	values := url.Values{}

	if limit > 0 {
		values.Set(paramLimit, strconv.Itoa(limit))
	}
	if timeout > 0 {
		values.Set(paramTimeout, strconv.Itoa(timeout))
	}
	if marker > 0 {
		values.Set(paramMarker, strconv.FormatInt(marker, 10))
	}
	if len(s.updateTypes) > 0 {
		types := make([]string, len(s.updateTypes))
		for i, t := range s.updateTypes {
			types[i] = string(t)
		}
		values.Set(paramTypes, strings.Join(types, ","))
	}

	err = s.client.raw(ctx, http.MethodGet, pathUpdates, values, nil, &res)

	return
}
