package maxbot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type Bots struct {
	client *client
}

func newBots(cli *client) *Bots {
	return &Bots{
		client: cli,
	}
}

func (b *Bots) GetMyInfo(ctx context.Context) (info model.BotInfo, err error) {
	err = b.client.raw(ctx, http.MethodGet, pathMe, nil, nil, &info)
	if err != nil {
		err = fmt.Errorf(`GetMyInfo: %w`, err)
	}

	return
}

func (b *Bots) EditMyInfo(ctx context.Context, botPath model.BotPatch) (info model.BotInfo, err error) {
	err = b.client.raw(ctx, http.MethodPatch, pathMe, nil, botPath, &info)
	if err != nil {
		err = fmt.Errorf(`EditMyInfo: %w`, err)
	}

	return
}
