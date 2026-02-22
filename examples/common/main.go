//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"

	"github.com/rs/zerolog/log"
)

func main() {
	opts := []maxbot.Option{
		maxbot.WithHTTPClient(&http.Client{Timeout: time.Second}),
		maxbot.WithClientTimeout(30 * time.Second),
		maxbot.WithDebugMode(),
	}
	api, err := maxbot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("NewWithConfig failed. Stop.")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	info, err := api.Bots.GetBot(ctx) // Простой метод
	log.Printf("Get me: %#v %#v", info, err)

	info, err = api.Bots.PatchBot(ctx, &schemes.BotPatch{Commands: []schemes.BotCommand{{Name: "shutdown", Description: "Перезапускает бота"}}}) // Простой метод
	log.Printf("Get me: %#v %#v", info, err)

	chatList, err := api.Chats.GetChats(ctx, 0, 0)
	if err != nil {
		fmt.Printf("Unknown type: %#v", err)
	}
	for _, chat := range chatList.Chats {
		fmt.Printf("Bot is members at the chat: %#v", chat.Title)
		fmt.Printf("	: %#v", chat.ChatId)
	}

	for update := range api.GetUpdates(ctx) { // Чтение из канала с обновлениями
		_ = api.Debugs.Send(ctx, update)
		var outText string
		switch upd := update.(type) { // Определение типа пришедшего обновления
		case *schemes.MessageCreatedUpdate:
			out := "bot прочитал текст: " + upd.Message.Body.Text

			switch upd.GetCommand() {
			case "/chats":
				out = "команда : " + upd.GetCommand()
				err = api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText(out))
				log.Printf("Answer: %#v", err)
				continue
			case "/chats_full":
				cl, err := api.Chats.GetChats(ctx, 0, 0)
				if err != nil {
					log.Printf("Unknown type: %#v", err)
				}

				outText = "List of chats\n"
				for _, chat := range cl.Chats {
					outText += fmt.Sprintf(" 	   title: %#v\n", chat.Title)
					outText += fmt.Sprintf("	      id: %#v\n", chat.ChatId)
					outText += fmt.Sprintf(" description: %#v\n", chat.Description)
					outText += fmt.Sprintf("   is public: %#v\n", chat.IsPublic)
					outText += fmt.Sprintf("   		link: %#v\n", chat.Link)
					outText += fmt.Sprintf("   	  status: %#v\n", chat.Status)
					outText += fmt.Sprintf("       owner: %#v\n", chat.OwnerId)
					outText += fmt.Sprintf("       type: %#v\n", chat.Type)
					outText += "______\n"
				}

				_ = api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText(outText))
				continue
			}

			keyboard := api.Messages.NewKeyboardBuilder()
			keyboard.
				AddRow().
				AddGeolocation("Прислать геолокацию", true).
				AddContact("Прислать контакт")
			keyboard.
				AddRow().
				AddLink("Ссылка", schemes.POSITIVE, "https://max.ru").
				AddCallback("Аудио", schemes.NEGATIVE, "audio").
				AddCallback("Видео", schemes.NEGATIVE, "video")
			keyboard.
				AddRow().
				AddCallback("Картинка", schemes.POSITIVE, "picture")

			_ = api.Messages.Send(ctx, maxbot.NewMessage().SetUser(upd.Message.Sender.UserId).SetReply("И вам привет!(в личку!)", upd.Message.Body.Mid))

			_ = api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetReply("И вам привет! (в чат)", upd.Message.Body.Mid))

			// Отправка сообщения с клавиатурой
			_ = api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).AddKeyboard(keyboard).SetText(outText))
			_ = api.Messages.Send(ctx, maxbot.NewMessage().Reply("**Reply** universal", upd.Message).SetFormat("markdown"))

		case *schemes.MessageCallbackUpdate:
			// Ответ на callback
			msg := maxbot.NewMessage()
			if upd.Message.Recipient.UserId != 0 {
				msg.SetUser(upd.Message.Recipient.UserId)
			}
			if upd.Message.Recipient.ChatId != 0 {
				msg.SetChat(upd.Message.Recipient.ChatId)
			}

			if upd.Callback.Payload == "picture" {
				photo, err := api.Uploads.UploadPhotoFromFile(ctx, "./big-logo.png")
				if err != nil {
					log.Err(err).Msg("Uploads.UploadPhotoFromFile")
					break
				}

				msg.AddPhoto(photo) // прикрепляем к сообщению изображение
				if err = api.Messages.Send(ctx, msg); err != nil {
					log.Err(err).Msg("Messages.Send")
				}
			}

			if upd.Callback.Payload == "audio" {
				if audio, err := api.Uploads.UploadMediaFromFile(ctx, schemes.AUDIO, "./music.mp3"); err == nil {
					msg.AddAudio(audio) // прикрепляем к сообщению mp3
				} else {
					log.Err(err).Msg("Uploads.UploadPhotoFromFile")
					break
				}

				if err = api.Messages.Send(ctx, msg); err != nil {
					log.Err(err).Msg("Messages.Send")
				}
			}

			if upd.Callback.Payload == "video" {
				if video, err := api.Uploads.UploadMediaFromFile(ctx, schemes.VIDEO, "./video.mp4"); err == nil {
					msg.AddVideo(video) // прикрепляем к сообщению mp4
				} else {
					log.Err(err).Msg("Uploads.UploadPhotoFromFile")
					break
				}

				if err = api.Messages.Send(ctx, msg); err != nil {
					log.Err(err).Msg("Messages.Send")
				}
			}

			if upd.Callback.Payload == "file" {
				if doc, err := api.Uploads.UploadMediaFromFile(ctx, schemes.FILE, "./max.pdf"); err == nil {
					msg.AddFile(doc) // прикрепляем к сообщению pdf file
				} else {
					log.Err(err).Msg("Uploads.UploadPhotoFromFile")
					break
				}

				if err = api.Messages.Send(ctx, msg); err != nil {
					log.Err(err).Msg("Messages.Send")
				}
			}

		default:
			log.Printf("Unknown type: %#v", upd)
		}
	}
}
