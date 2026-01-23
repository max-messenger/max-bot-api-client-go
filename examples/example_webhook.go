//go:build ignore

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
	"os/signal"
	"syscall"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	// Initialisation
	api, _ := maxbot.New(os.Getenv("TOKEN"))
	host := os.Getenv("HOST")

	// Some methods demo:
	info, err := api.Bots.GetBot(ctx)
	log.Printf("Get me: %#v %#v", info, err)

	subs, _ := api.Subscriptions.GetSubscriptions(ctx)
	for _, s := range subs.Subscriptions {
		_, _ = api.Subscriptions.Unsubscribe(ctx, s.Url)
	}
	subscriptionResp, err := api.Subscriptions.Subscribe(ctx, host+"/webhook", []string{}, "my-secret-phrase")
	log.Printf("Subscription: %#v %#v", subscriptionResp, err)

	ch := make(chan schemes.UpdateInterface) // Channel with updates from Max

	http.HandleFunc("/webhook", api.GetHandler(ch))
	go func() {
		for {
			upd := <-ch
			log.Printf("Received: %#v", upd)
			switch upd := upd.(type) {
			case *schemes.MessageCreatedUpdate:
				message := maxbot.NewMessage().
					SetUser(upd.Message.Sender.UserId).
					SetText(fmt.Sprintf("Hello, %s! Your message: %s", upd.Message.Sender.Name, upd.Message.Body.Text))

				err := api.Messages.Send(ctx, message)
				log.Printf("Answer: %#v", err)
			default:
				log.Printf("Unknown type: %#v", upd)
			}
		}
	}()

	http.ListenAndServe(":10888", nil)
}
