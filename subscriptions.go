package maxbot

import (
	"context"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type subscriptions struct {
	client *client
}

func newSubscriptions(client *client) *subscriptions {
	return &subscriptions{client: client}
}

// GetSubscriptions returns the list of all subscriptions.
func (a *subscriptions) GetSubscriptions(ctx context.Context) (*schemes.GetSubscriptionsResult, error) {
	result := new(schemes.GetSubscriptionsResult)
	values := url.Values{}

	body, err := a.client.request(ctx, http.MethodGet, pathSubscriptions, values, false, nil)
	if err != nil {
		return result, err
	}
	defer a.client.closer("getSubscriptions body", body)

	return result, jsoniter.NewDecoder(body).Decode(result)
}

// Subscribe subscribes the bot to receive updates via WebHook.
func (a *subscriptions) Subscribe(ctx context.Context, subscribeURL string, updateTypes []string, secret string) (*schemes.SimpleQueryResult, error) {
	subscription := &schemes.SubscriptionRequestBody{
		Secret:      secret,
		Url:         subscribeURL,
		UpdateTypes: updateTypes,
		Version:     a.client.version,
	}
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}

	body, err := a.client.request(ctx, http.MethodPost, pathSubscriptions, values, false, subscription)
	if err != nil {
		return result, err
	}
	defer a.client.closer("getSubscribe body", body)

	return result, jsoniter.NewDecoder(body).Decode(result)
}

// Unsubscribe unsubscribes the bot from receiving updates via WebHook.
func (a *subscriptions) Unsubscribe(ctx context.Context, subscriptionURL string) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	values.Set(paramURL, subscriptionURL)

	body, err := a.client.request(ctx, http.MethodDelete, pathSubscriptions, values, false, nil)
	if err != nil {
		return result, err
	}
	defer a.client.closer("unSubscribe body", body)

	return result, jsoniter.NewDecoder(body).Decode(result)
}
