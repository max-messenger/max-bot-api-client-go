# `3` Отправка сообщений с вложениями

Для упрощения работы с вложениями существует модуль `Uploads`.

## Отправка файлов

### Загрузка новых файлов

Подходит для файлов на диске:

```go
// Ответ на коллбек
ctx := context.Background()

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
	if _, err := api.Messages.Send(msg); err != nil {
		log.Err(err).Msg("Messages.Send")
	}
}
if upd.Callback.Payload == "audio" {
	if audio, err := api.Uploads.UploadMediaFromFile(ctx, schemes.AUDIO, "./music.mp3"); err == nil {
		msg.AddAudio(audio) // прикрепляем к сообщению mp3
	} else {
		log.Err(err).Msg("Uploads.UploadMediaFromFile")
		break
	}
	if _, err := api.Messages.Send(msg); err != nil {
		log.Err(err).Msg("Messages.Send")
	}
}
if upd.Callback.Payload == "video" {
	if video, err := api.Uploads.UploadMediaFromFile(ctx, schemes.VIDEO, "./video.mp4"); err == nil {
		msg.AddVideo(video) // прикрепляем к сообщению mp4
	} else {
		log.Err(err).Msg("Uploads.UploadMediaFromFile")
		break
	}
	if _, err := api.Messages.Send(msg); err != nil {
		log.Err(err).Msg("Messages.Send")
	}
}
if upd.Callback.Payload == "file" {
	if doc, err := api.Uploads.UploadMediaFromFile(ctx, schemes.FILE, "./max.pdf"); err == nil {
		msg.AddFile(doc) // прикрепляем к сообщению pdf file
	} else {
		log.Err(err).Msg("Uploads.UploadMediaFromFile")
		break
	}
	if _, err := api.Messages.Send(msg); err != nil {
		log.Err(err).Msg("Messages.Send")
	}
}
```

### При помощи ссылки

```go
// Ответ на коллбек
msg := maxbot.NewMessage()
if upd.Message.Recipient.UserId != 0 {
	msg.SetUser(upd.Message.Recipient.UserId)
}
if upd.Message.Recipient.ChatId != 0 {
	msg.SetChat(upd.Message.Recipient.ChatId)
}
if upd.Callback.Payload == "picture" {
	photo, err := api.Uploads.UploadMediaFromUrl(ctx, schemes.PHOTO, "https://max.ru/s/img/big-logo.png")
	if err != nil {
		log.Err(err).Msg("Uploads.UploadMediaFromUrl")
		break
	}
	msg.AddPhoto(photo) // прикрепляем к сообщению изображение
	if _, err := api.Messages.Send(msg); err != nil {
		log.Err(err).Msg("Messages.Send")
	}
}
if upd.Callback.Payload == "audio" {
	if audio, err := api.Uploads.UploadMediaFromUrl(ctx, schemes.AUDIO, "https://max.ru/s/audio/music.mp3"); err == nil {
		msg.AddAudio(audio) // прикрепляем к сообщению mp3
	} else {
		log.Err(err).Msg("Uploads.UploadPhotoFromFile")
		break
	}
	if _, err := api.Messages.Send(msg); err != nil {
		log.Err(err).Msg("Messages.Send")
	}
}
if upd.Callback.Payload == "video" {
	if video, err := api.Uploads.UploadMediaFromUrl(ctx, schemes.VIDEO, "https://max.ru/s/video/reactions.mp4"); err == nil {
		msg.AddVideo(video) // прикрепляем к сообщению mp4
	} else {
		log.Err(err).Msg("Uploads.UploadPhotoFromFile")
		break
	}
	if _, err := api.Messages.Send(msg); err != nil {
		log.Err(err).Msg("Messages.Send")
	}
}
if upd.Callback.Payload == "file" {
	if doc, err := api.Uploads.UploadMediaFromUrl(ctx, schemes.FILE, "https://max.ru/s/docs/tekhnicheskaya-dokumentatsiya.zip"); err == nil {
		msg.AddFile(doc) // прикрепляем к сообщению zip file
	} else {
		log.Err(err).Msg("Uploads.UploadPhotoFromFile")
		break
	}
	if _, err := api.Messages.Send(msg); err != nil {
		log.Err(err).Msg("Messages.Send")
	}
}
```
