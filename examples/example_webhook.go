//go:build ignore
// +build ignore

/**
 * Webhook example
 */
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	maxbot "github.com/max-messenger/max-bot-api-client-go"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func main() {
	ctx := context.Background()
	// Initialisation
	api, err := maxbot.New(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatalf("failed to create api client: %v", err)
	}
	host := os.Getenv("HOST")

	// Some methods demo:
	info, err := api.Bots.GetBot(ctx)
	log.Printf("Get me: %#v %#v", info, err)

	subs, _ := api.Subscriptions.GetSubscriptions(ctx)
	for _, s := range subs.Subscriptions {
		_, _ = api.Subscriptions.Unsubscribe(ctx, s.Url)
	}
	subscriptionResp, err := api.Subscriptions.Subscribe(ctx, host+"/webhook", []string{})
	log.Printf("Subscription: %#v %#v", subscriptionResp, err)

	ch := make(chan schemes.UpdateInterface) // Channel with updates from Max

	http.HandleFunc("/webhook", api.GetHandler(ch))
	go func() {
		for {
			upd := <-ch
			log.Printf("Received: %#v", upd)
			switch upd := upd.(type) {
			case *schemes.MessageCreatedUpdate:
				_, err := api.Messages.Send(
					ctx,
					maxbot.NewMessage().
						SetUser(upd.Message.Sender.UserId).
						SetText(fmt.Sprintf("Hello, %s! Your message: %s", upd.Message.Sender.Name, upd.Message.Body.Text)),
				)
				log.Printf("Answer: %#v", err)
			default:
				log.Printf("Unknown type: %#v", upd)
			}
		}
	}()

	http.ListenAndServe(":10888", nil)
}
