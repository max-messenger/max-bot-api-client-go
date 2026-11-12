// Command upload — единственный метод группы UploadAPI.
//
//	go run ./examples/upload upload [video|audio|file]
package upload

import (
	"context"
	"fmt"
	"os"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

var Commands = []setup.Command{
	{
		Name: "upload", What: "Upload.Upload (4 типа) -> AddAttachByToken", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken, setup.KeyUserA},
		Run:   upload,
	},
}

// kinds — фикстуры существующие, новых бинарей не плодим (см. Non-Goals плана).
// Тип image в test_data/ отсутствует: он проверяется встроенным PNG в
// tests/doc03_attachments_test.go, поэтому здесь вынесен в отдельную ветку.
var kinds = map[string]struct {
	path  string
	ut    model.UploadType
	at    model.AttachmentType
	token string // откуда приходит токен: ответ POST /uploads или JSON upload-хоста
}{
	"video": {"test_data/video.mp4", model.UploadVideo, model.AttachVideo, "ответ POST /uploads"},
	"audio": {"test_data/music_test.mp3", model.UploadAudio, model.AttachAudio, "ответ POST /uploads"},
	"file":  {"test_data/test_file.txt", model.UploadFile, model.AttachFile, "JSON upload-хоста"},
}

func upload(_ context.Context, api *maxbot.Api, env setup.Env, args []string) error {
	kind := "file"
	if len(args) > 0 {
		kind = args[0]
	}
	if kind == "image" {
		fmt.Println("НЕ ВАЛИДИРОВАНО: для image нет локальной фикстуры в test_data/ (type image")
		fmt.Println("покрыт встроенным PNG в tests/doc03_attachments_test.go). Добавьте файл")
		fmt.Println("test_data/logo.png и пример станет запускаемым — новых бинарей план не предполагает.")

		return nil
	}
	k, ok := kinds[kind]
	if !ok {
		return fmt.Errorf("неизвестный тип %q; доступны video|audio|file и image (не валидирован)", kind)
	}
	path, err := setup.Fixture(k.path)
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("открыть %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()
	info, err := f.Stat()
	if err != nil {
		return err
	}

	ctx := context.Background()
	token, err := api.Upload.Upload(ctx, k.ut, f, info.Name(), info.Size())
	if err != nil {
		return fmt.Errorf("Upload.Upload(%s): %w", k.ut, err)
	}
	if token == "" {
		return fmt.Errorf("Upload.Upload(%s) вернул пустой токен", k.ut)
	}
	// Значение токена — доступ к контенту, в вывод оно не попадает.
	fmt.Printf("загружено %s (%d байт), токен: тип=%s длина=%d, источник=%s\n",
		k.ut, info.Size(), k.at, len(token), k.token)

	res, err := api.Messages.Send(ctx, maxbot.NewMessage().
		SetUser(env.UserA).SetText("пример: загрузка "+string(k.ut)).
		AddAttachByToken(token, k.at))
	if err != nil {
		return fmt.Errorf("Messages.Send с вложением %s: %w", k.at, err)
	}
	mid := res.Message.Body.Mid
	defer func() {
		if _, err := api.Messages.DeleteMessage(ctx, mid); err != nil {
			fmt.Printf("cleanup: DeleteMessage(%s) не удался: %v\n", mid, err)

			return
		}
		fmt.Printf("cleanup: удалено mid=%s\n", mid)
	}()

	var found bool
	for _, a := range res.Message.Body.Attachments {
		if a.Type == k.at {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("сообщение %s не содержит вложение типа %s", mid, k.at)
	}
	fmt.Printf("OK: токен принят, сообщение с вложением %s создано и удалено\n", k.at)

	return nil
}
