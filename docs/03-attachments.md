# `3` Отправка сообщений с вложениями

Для упрощения работы с вложениями существует модуль `Uploads`.

## Отправка файлов

### Загрузка новых файлов

Подходит для файлов на диске:

```go
func videoHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
    f, err := os.Open("upload/video.mp4")
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

```

### При помощи ссылки

```go
func imageHandler(ctx context.Context, api *maxbot.Api, update model.Update) {
    // отправить изображения без текста
    msg := maxbot.NewMessage().
    SetChat(update.ChatID).
    AddImageUrl("https://raw.githubusercontent.com/max-messenger/max-bot-api-client-go/refs/heads/main/examples/common/big-logo.png").
    AddImageUrl("https://raw.githubusercontent.com/max-messenger/max-bot-api-client-go/refs/heads/main/examples/common/big-logo.png")
    
    _, _ = api.Messages.Send(ctx, msg)
}
```
