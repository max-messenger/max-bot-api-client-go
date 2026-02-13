package maxbot

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type bots struct {
	client *client
}

func newBots(client *client) *bots {
	return &bots{client: client}
}

// GetBot returns info about current bot.
// Current bot can be identified by access token. Method returns bot identifier, name and avatar (if any).
func (a *bots) GetBot(ctx context.Context) (*schemes.BotInfo, error) {
	result := new(schemes.BotInfo)
	values := url.Values{}

	body, err := a.client.request(ctx, http.MethodGet, pathMe, values, false, nil)
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

// PatchBot edits current bot info. Fill only the fields you want to update. All remaining fields will stay untouched.
func (a *bots) PatchBot(ctx context.Context, patch *schemes.BotPatch) (*schemes.BotInfo, error) {
	result := new(schemes.BotInfo)
	values := url.Values{}

	body, err := a.client.request(ctx, http.MethodPatch, pathMe, values, false, patch)
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
