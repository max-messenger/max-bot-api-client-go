package maxbot

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type subscriptions struct {
	client *client
}

func newSubscriptions(client *client) *subscriptions {
	return &subscriptions{client: client}
}

// GetSubscriptions returns the list of all subscriptions
func (a *subscriptions) GetSubscriptions(ctx context.Context) (*schemes.GetSubscriptionsResult, error) {
	result := new(schemes.GetSubscriptionsResult)
	values := url.Values{}
	body, err := a.client.request(ctx, http.MethodGet, "subscriptions", values, false, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return result, json.NewDecoder(body).Decode(result)
}

// Subscribe subscribes bot to receive updates via WebHook
func (a *subscriptions) Subscribe(ctx context.Context, subscribeURL string, updateTypes []string, secret string) (*schemes.SimpleQueryResult, error) {
	subscription := &schemes.SubscriptionRequestBody{
		Secret:      secret,
		Url:         subscribeURL,
		UpdateTypes: updateTypes,
		Version:     a.client.version,
	}
	values := url.Values{}
	body, err := a.client.request(ctx, http.MethodPost, "subscriptions", values, false, subscription)
	if err != nil {
		return nil, err
	}
	res, raw, err := decodeSimpleQueryResult(body)
	if err != nil {
		return res, err
	}
	if apiErr := newSimpleQueryAPIError("subscribe", res, raw); apiErr != nil {
		return res, apiErr
	}
	return res, nil
}


// Unsubscribe unsubscribes bot from receiving updates via WebHook
func (a *subscriptions) Unsubscribe(ctx context.Context, subscriptionURL string) (*schemes.SimpleQueryResult, error) {
	values := url.Values{}
	values.Set("url", subscriptionURL)
	body, err := a.client.request(ctx, http.MethodDelete, "subscriptions", values, false, nil)
	if err != nil {
		return nil, err
	}
	res, raw, err := decodeSimpleQueryResult(body)
	if err != nil {
		return res, err
	}
	if apiErr := newSimpleQueryAPIError("unsubscribe", res, raw); apiErr != nil {
		return res, apiErr
	}
	return res, nil
}

