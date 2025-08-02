//go:build ignore
// +build ignore

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
	ctx := context.Background()

	// Initialisation
	api, err := maxbot.New(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatalf("Error creating API client: %v", err)
	}

	// Some methods demo:
	info, err := api.Bots.GetBot(ctx)
	log.Printf("Get me: %#v %#v", info, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for upd := range api.GetUpdates(ctx) {
			log.Printf("Received: %#v", upd)
			switch upd := upd.(type) {
			case *schemes.MessageCreatedUpdate:
				_, err := api.Messages.Send(
					ctx,
					maxbot.NewMessage().
						SetUser(upd.Message.Sender.UserId).
						SetText(fmt.Sprintf("Hello, %s! Your message: %s", upd.Message.Sender.Name, upd.Message.Body.Text)),
				)
				if err != nil {
					log.Printf("Error: %#v", err)
				}
			default:
				log.Printf("Unknown type: %#v", upd)
			}
		}
	}()
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, os.Interrupt)
		select {
		case <-exit:
			cancel()
		case <-ctx.Done():
			return
		}
	}()
	<-ctx.Done()
}
