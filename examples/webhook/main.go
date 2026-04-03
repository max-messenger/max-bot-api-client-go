//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type httpClient struct {
	httpClient *http.Client
}

// Do use as middleware
func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	// your trace and metrics
	log.Printf("sending request to %s", req.URL.String())
	r, err := c.httpClient.Do(req)
	log.Printf("received response from %s", req.URL.String())
	return r, err
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	secret := os.Getenv("MAX_BOT_API_SECRET")
	httpCli := &httpClient{
		httpClient: &http.Client{
			Timeout: time.Second * 35,
		},
	}

	// Initialisation
	opts := []maxbot.Option{
		maxbot.WithDebugMode(),
		maxbot.WithHTTPClient(httpCli),
	}
	api, err := maxbot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	host := os.Getenv("HOST")

	errChan := api.GetErrors()
	go func() {
		for errMessage := range errChan {
			log.Println(errMessage) // use your favorite logger
		}
	}()

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

	http.HandleFunc("/webhook", api.GetUpdateHandler(ch, secret))
	go func() {
		for {
			update := <-ch
			log.Printf("Received: %#v", update)
			switch upd := update.(type) {
			case *schemes.MessageCreatedUpdate:
				message := maxbot.NewMessage().
					SetUser(upd.Message.Sender.UserId).
					SetText(fmt.Sprintf("Hello, %s! Your message: %s", upd.Message.Sender.Name, upd.Message.Body.Text))

				err = api.Messages.Send(ctx, message)
				log.Printf("Answer: %#v", err)
			default:
				log.Printf("Unknown type: %#v", upd)
			}
		}
	}()

	_ = http.ListenAndServe(":10888", nil)
}
