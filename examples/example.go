package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/configservice"
	"github.com/max-messenger/max-bot-api-client-go/schemes"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func url(user int64) string {
	buf := new(bytes.Buffer)
	byteOrder := binary.BigEndian

	binary.Write(buf, byteOrder, int64(user))
	fmt.Printf("uint64: %v\n", buf.Bytes())
	fmt.Printf("uint64: %v\n", string(buf.Bytes()))
	fmt.Printf("uint64b: %v\n", base64.StdEncoding.EncodeToString(buf.Bytes()))
	fmt.Printf("uint64bt: %v\n", strings.Trim(base64.StdEncoding.EncodeToString(buf.Bytes()), "="))
	return "https://max.ru/u/" + strings.Trim(base64.StdEncoding.EncodeToString(buf.Bytes()), "=")
}

func main() {
	var configPath string
	var maxenv = os.Getenv("MAXBOT_ENV")
	// Customize ConsoleWriter

	if maxenv != "" {
		configPath = "/go/bin/config/app-" + maxenv + ".yaml"
		//zerolog.SetGlobalLevel(zerolog.InfoLevel)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		//		log.Logger = log.Output(consoleWriter).With().Caller().Logger()
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

	api := maxbot.NewFormConfig(configService)

	info, err := api.Bots.GetBot() // Простой метод
	log.Printf("Get me: %#v %#v", info, err)

	info, err = api.Bots.PatchBot(&schemes.BotPatch{Commands: []schemes.BotCommand{{Name: "shutdown", Description: "Перезапускает бота"}}}) // Простой метод
	log.Printf("Get me: %#v %#v", info, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, syscall.SIGTERM, os.Interrupt)
		<-exit
		cancel()
	}()

	chatList, err := api.Chats.GetChats(0, 0)
	if err != nil {
		fmt.Printf("Unknown type: %#v", err)
	}
	for _, chat := range chatList.Chats {
		fmt.Printf("Bot is members at the chat: %#v", chat.Title)
		fmt.Printf("	: %#v", chat.ChatId)
	}

	for upd := range api.GetUpdates(ctx) { // Чтение из канала с обновлениями
		api.Debugs.Send(upd)
		switch upd := upd.(type) { // Определение типа пришедшего обновления
		case *schemes.MessageCreatedUpdate:
			out := "bot прочитал текст: " + upd.Message.Body.Text
			switch upd.GetCommand() {
			case "/chats":
				out = "команда : " + upd.GetCommand()
				_, err = api.Messages.Send(maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText(out))
				log.Printf("Answer: %#v", err)
				continue
			case "/chats_full":
				chatList, err := api.Chats.GetChats(0, 0)
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
					out += fmt.Sprintf("______\n")
				}
				api.Messages.SendMessageResult(maxbot.NewMessage().SetReply("И вам привет!", upd.Message.Body.Mid))
				mes, err := api.Messages.SendMessageResult(maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText(out))
				fmt.Printf("Answer: %v", mes.Body.Mid)
				continue
			}
			keyboard := api.Messages.NewKeyboardBuilder()
			keyboard.
				AddRow().
				AddGeolocation("Прислать геолокацию", true).
				AddContact("Прислать контакт")
			keyboard.
				AddRow().
				AddLink("Cсылка", schemes.POSITIVE, "https://max.ru").
				AddCallback("Аудио", schemes.NEGATIVE, "audio").
				AddCallback("Видео", schemes.NEGATIVE, "video")
			keyboard.
				AddRow().
				AddCallback("Картинка", schemes.POSITIVE, "picture")

			mes, _ := api.Messages.SendMessageResult(maxbot.NewMessage().SetUser(upd.Message.Sender.UserId).SetReply("И вам привет!(в личку!)", upd.Message.Body.Mid))
			api.Messages.SendMessageResult(maxbot.NewMessage().SetUser(upd.Message.Sender.UserId).SetReply("И вам привет!(в личку!)", mes.Body.Mid))
			reply_id, err := api.Messages.Send(maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetReply("И вам привет! (в чат)", upd.Message.Body.Mid))
			api.Messages.Send(maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetReply("И вам привет! (в чат) на rep", reply_id))
			// Отправка сообщения с клавиатурой
			id, err := api.Messages.Send(maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).AddKeyboard(keyboard).SetText(out))
			mes_rep, _ := api.Messages.SendMessageResult(maxbot.NewMessage().Reply("**Reply** univesal", upd.Message).SetFormat("markdown"))
			api.Messages.SendMessageResult(maxbot.NewMessage().Reply("<b>Привет!</b> <i>Добро пожаловать</i>", mes_rep).SetFormat("html"))
			fmt.Printf("Answer:%v : %v", id, err)

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
				photo, err := api.Uploads.UploadPhotoFromFile("./big-logo.png")
				if err != nil {
					log.Err(err).Msg("Uploads.UploadPhotoFromFile")
					break
				}
				msg.AddPhoto(photo) // прикрипляем к сообщению изображение
				if _, err := api.Messages.SendMessageResult(msg); err != nil {
					log.Err(err).Msg("Messages.Send")
				}
			}
			if upd.Callback.Payload == "audio" {
				if audio, err := api.Uploads.UploadMediaFromFile(schemes.AUDIO, "./music.mp3"); err == nil {
					msg.AddAudio(audio) // прикрипляем к сообщению mp3
				} else {
					log.Err(err).Msg("Uploads.UploadPhotoFromFile")
					break
				}
				if _, err := api.Messages.SendMessageResult(msg); err != nil {
					log.Err(err).Msg("Messages.Send")
				}
			}
			if upd.Callback.Payload == "video" {
				if video, err := api.Uploads.UploadMediaFromFile(schemes.VIDEO, "./video.mp4"); err == nil {
					msg.AddVideo(video) // прикрипляем к сообщению mp4
				} else {
					log.Err(err).Msg("Uploads.UploadPhotoFromFile")
					break
				}
				if _, err := api.Messages.SendMessageResult(msg); err != nil {
					log.Err(err).Msg("Messages.Send")
				}
			}
			if upd.Callback.Payload == "file" {
				if doc, err := api.Uploads.UploadMediaFromFile(schemes.FILE, "./max.pdf"); err == nil {
					msg.AddFile(doc) // прикрипляем к сообщению pdf file
				} else {
					log.Err(err).Msg("Uploads.UploadPhotoFromFile")
					break
				}
				if _, err := api.Messages.SendMessageResult(msg); err != nil {
					log.Err(err).Msg("Messages.Send")
				}
			}

		default:
			log.Printf("Unknown type: %#v", upd)
		}
	}
}
