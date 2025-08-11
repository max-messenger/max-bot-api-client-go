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
func (a *subscriptions) Subscribe(ctx context.Context, subscribeURL string, updateTypes []string) (*schemes.SimpleQueryResult, error) {
	subscription := &schemes.SubscriptionRequestBody{
		Url:         subscribeURL,
		UpdateTypes: updateTypes,
		Version:     a.client.version,
	}
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	body, err := a.client.request(ctx, http.MethodPost, "subscriptions", values, false, subscription)
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

// Unsubscribe unsubscribes bot from receiving updates via WebHook
func (a *subscriptions) Unsubscribe(ctx context.Context, subscriptionURL string) (*schemes.SimpleQueryResult, error) {
	result := new(schemes.SimpleQueryResult)
	values := url.Values{}
	values.Set(paramURL, subscriptionURL)
	body, err := a.client.request(ctx, http.MethodDelete, "subscriptions", values, false, nil)
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
