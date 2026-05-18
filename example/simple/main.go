package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

//go:embed all:upload
var uploadStore embed.FS

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []maxbot.Opt{
		maxbot.WithHTTPClient(newHttpClient()),
	}

	api, err := maxbot.NewApi(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		log.Fatal(err)
	}

	info, err := api.Bots.GetMyInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("info: %+v", info)

	handle := func(ctx context.Context, update model.Update) {
		fmt.Printf("Received: [%s] %#v\n", update.UpdateType, update)
		switch update.UpdateType {
		case model.UpdateMessageCreated:
			text := update.GetMessage().Body.Text
			switch text {
			case "/image":
				imageHandler(ctx, api, update)
			case "/video":
				videoHandler(ctx, api, update)
			case "/file":
				fileHandler(ctx, api, update)
			case "/audio":
				audioHandler(ctx, api, update)
			case "/html":
				htmlHandler(ctx, api, update)
			case "/md":
				mdHandler(ctx, api, update)
			case "/sticker":
				stickerHandler(ctx, api, update)
			case "/contact":
				contactHandler(ctx, api, update)
			case "/location":
				locationHandler(ctx, api, update)
			case "/share":
				shareHandler(ctx, api, update)
			default:
				textHandler(ctx, api, update)
			}

		default:
			log.Printf("Unknown type: %#v\n", update)
		}
	}

	var updates []model.Update
	var marker int64
	for {
		select {
		case <-ctx.Done():
		default:
			updates, marker, err = api.Subscriptions.GetUpdates(ctx, marker)
			if err != nil {
				log.Println("GetUpdates: ", err)
				return
			}

			for _, update := range updates {
				handle(ctx, update)
			}
		}
	}
}

func textHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	keyboard := model.NewKeyboard()

	//// Предварительно url miniApp должен быть установлен в настройках masterBot
	// var botId int64
	// keyboard.
	//	AddRow().
	//	AddOpenApp("open", botId)

	keyboard.
		AddRow().
		AddLink("link", "http://max.ru").
		AddGeoLocation("link", true).
		AddContact("contact")

	keyboard.
		AddRow().
		AddMessage("message").
		AddClipboard("скопировать url", "https://max.ru").
		AddCallback("callback button", model.IntentDefault, "callback")

	msg := maxbot.NewMessage().
		SetText("hello").
		AddKeyboard(keyboard).
		SetChat(update.ChatID).
		SetUser(update.UserID)

	fmt.Printf("--text: %s. %v\n", update.Message.Body.Text, update)
	fmt.Println(msg)
	res, _ := api.Messages.Send(ctx, msg)

	log.Printf("%v\n", res)
}

func imageHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	// отправить изображения без текста
	msg := maxbot.NewMessage().
		SetChat(update.ChatID).
		AddImageUrl("https://raw.githubusercontent.com/max-messenger/max-bot-api-client-go/refs/heads/main/examples/common/big-logo.png").
		AddImageUrl("https://raw.githubusercontent.com/max-messenger/max-bot-api-client-go/refs/heads/main/examples/common/big-logo.png")

	_, _ = api.Messages.Send(ctx, msg)
}

func htmlHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	msg := maxbot.NewMessage().
		SetUser(update.UserID).
		SetFormat(model.FormatHTML).
		SetText(fmt.Sprintf(`hello <b><i>%s</i></b>`, update.GetUser().Name))
	_, _ = api.Messages.Send(ctx, msg)
}

func mdHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	msg := maxbot.NewMessage().
		SetUser(update.UserID).
		SetFormat(model.FormatMarkdown).
		SetText(fmt.Sprintf(`hi **_%s_**`, update.GetUser().Name))
	_, _ = api.Messages.Send(ctx, msg)
}

func videoHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	f, err := uploadStore.Open("upload/video.mp4")
	if err != nil {
		log.Fatal(err)
	}

	info, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	token, err := api.Upload.Upload(ctx, model.UploadVideo, f, info.Name(), info.Size())
	if err != nil {
		log.Println("upload error:", err)
	}

	msg := maxbot.NewMessage().
		SetUser(update.UserID).
		AddAttachByToken(token, model.AttachVideo)

	_, _ = api.Messages.Send(ctx, msg)
}

func fileHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	f, err := uploadStore.Open("upload/video.mp4")
	if err != nil {
		log.Fatal(err)
	}

	info, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	token, err := api.Upload.Upload(ctx, model.UploadFile, f, info.Name(), info.Size())
	if err != nil {
		log.Println("upload error:", err)

		return
	}

	msg := maxbot.NewMessage().
		SetUser(update.UserID).
		AddAttachByToken(token, model.AttachFile)

	_, _ = api.Messages.Send(ctx, msg)
}

func audioHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	f, err := uploadStore.Open("upload/music.mp3")
	if err != nil {
		log.Fatal(err)
	}

	info, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	token, err := api.Upload.Upload(ctx, model.UploadAudio, f, info.Name(), info.Size())
	if err != nil {
		log.Println("upload error:", err)

		return
	}

	msg := maxbot.NewMessage().
		SetUser(update.UserID).
		AddAttachByToken(token, model.AttachAudio)

	_, _ = api.Messages.Send(ctx, msg)
}

func stickerHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	msg := maxbot.NewMessage().
		SetUser(update.UserID).
		AddSticker("1a27781bb")

	_, _ = api.Messages.Send(ctx, msg)
}

func contactHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	contactID := update.UserID
	msg := maxbot.NewMessage().
		SetUser(update.UserID).
		AddContact(contactID)

	_, _ = api.Messages.Send(ctx, msg)
}

func locationHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	msg := maxbot.NewMessage().
		SetUser(update.UserID).
		SetText("Букачача").
		AddLocation(52.981201, 116.902013)

	_, _ = api.Messages.Send(ctx, msg)
}

func shareHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
	msg := maxbot.NewMessage().
		SetUser(update.UserID).
		SetText("Документация").
		AddShare("https://dev.max.ru/docs-api/methods/POST/messages")

	_, _ = api.Messages.Send(ctx, msg)
}
