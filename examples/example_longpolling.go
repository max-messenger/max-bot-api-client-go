//go:build ignore

/**
 * Updates loop example
 */
package main

import (
	"context"
	"fmt"
	"log"
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

	// Some methods demo:
	info, err := api.Bots.GetBot(ctx)
	log.Printf("Get me: %#v %#v", info, err)

	go func() {
		for upd := range api.GetUpdates(ctx) {
			log.Printf("Received: %#v", upd)
			switch upd := upd.(type) {
			case *schemes.MessageCreatedUpdate:
				message := maxbot.NewMessage().
					SetUser(upd.Message.Sender.UserId).
					SetText(fmt.Sprintf("Hello, %s! Your message: %s", upd.Message.Sender.Name, upd.Message.Body.Text))

				_, err := api.Messages.Send(ctx, message)
				if err != nil {
					log.Printf("Error: %#v", err)
				}
			default:
				log.Printf("Unknown type: %#v", upd)
			}
		}
	}()
	<-ctx.Done()
}
