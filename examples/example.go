package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/configservice"
	"github.com/max-messenger/max-bot-api-client-go/schemes"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var configPath string
	var maxenv = os.Getenv("MAXBOT_ENV")
	// Customize ConsoleWriter

	if maxenv != "" {
		configPath = "/go/bin/config/app-" + maxenv + ".yaml"
		// zerolog.SetGlobalLevel(zerolog.InfoLevel)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		// log.Logger = log.Output(consoleWriter).With().Caller().Logger()
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}).With().Timestamp().Caller().Logger()
	} else if 2 <= len(os.Args) {
		configPath = os.Args[1]
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		// log a human-friendly, colorized output
		//		log.Logger = log.Output(consoleWriter).With().Caller().Logger()
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	} else {
		log.Error().Msg("maxenv environment variable not found. Stop.")
		return
	}

	configService := configservice.NewConfigInterface(configPath)
	if configService == nil {
		log.Fatal().Str("configPath", configPath).Msg("NewConfigInterface failed. Stop.")
	}

	api, err := maxbot.NewWithConfig(configService)
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

	for upd := range api.GetUpdates(ctx) { // Чтение из канала с обновлениями
		api.Debugs.Send(ctx, upd)
		switch upd := upd.(type) { // Определение типа пришедшего обновления
		case *schemes.MessageCreatedUpdate:
			out := "bot прочитал текст: " + upd.Message.Body.Text

			switch upd.GetCommand() {
			case "/chats":
				out = "команда : " + upd.GetCommand()
				err = api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText(out))
				log.Printf("Answer: %#v", err)
				continue
			case "/chats_full":
				chatList, err := api.Chats.GetChats(ctx, 0, 0)
				if err != nil {
					log.Printf("Unknown type: %#v", err)
				}

				out := "List of chats\n"
				for _, chat := range chatList.Chats {
					out += fmt.Sprintf(" 	   title: %#v\n", chat.Title)
					out += fmt.Sprintf("	      id: %#v\n", chat.ChatId)
					out += fmt.Sprintf(" description: %#v\n", chat.Description)
					out += fmt.Sprintf("   is public: %#v\n", chat.IsPublic)
					out += fmt.Sprintf("   		link: %#v\n", chat.Link)
					out += fmt.Sprintf("   	  status: %#v\n", chat.Status)
					out += fmt.Sprintf("       owner: %#v\n", chat.OwnerId)
					out += fmt.Sprintf("       type: %#v\n", chat.Type)
					out += "______\n"
				}

				api.Messages.Send(ctx, maxbot.NewMessage().SetReply("И вам привет!", upd.Message.Body.Mid))
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
			keyboard.
				AddRow().
				AddMessage("Привет!")

			api.Messages.Send(ctx, maxbot.NewMessage().SetUser(upd.Message.Sender.UserId).SetReply("И вам привет!(в личку!)", upd.Message.Body.Mid))

			api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetReply("И вам привет! (в чат)", upd.Message.Body.Mid))

			// Отправка сообщения с клавиатурой
			api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).AddKeyboard(keyboard).SetText(out))
			api.Messages.Send(ctx, maxbot.NewMessage().Reply("**Reply** universal", upd.Message).SetFormat("markdown"))

		case *schemes.MessageCallbackUpdate:
			// Ответ на коллбек
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
				if err := api.Messages.Send(ctx, msg); err != nil {
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
				if err := api.Messages.Send(ctx, msg); err != nil {
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
				if err := api.Messages.Send(ctx, msg); err != nil {
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
				if err := api.Messages.Send(ctx, msg); err != nil {
					log.Err(err).Msg("Messages.Send")
				}
			}

		default:
			log.Printf("Unknown type: %#v", upd)
		}
	}
}
