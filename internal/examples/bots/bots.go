// Command bots — методы группы BotsAPI клиента max-bot-api-client-go/v2.
//
//	go run ./examples/bots get          # Bots.GetMyInfo
//	go run ./examples/bots commands     # Bots.PatchCommands
//	go run ./examples/bots edit-info    # Bots.EditMyInfo (только с -force)
package bots

import (
	"context"
	"fmt"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

// ourCommands — набор слаш-команд, который показывают примеры остальных доменов.
var ourCommands = []model.BotCommand{
	{Name: "start", Description: "знакомое сообщение"},
	{Name: "chats", Description: "id групп, где есть бот"},
	{Name: "whoami", Description: "информация о боте"},
}

var Commands = []setup.Command{
	{
		Name: "get", What: "Bots.GetMyInfo", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken},
		Run:   get,
	},
	{
		Name: "commands", What: "Bots.PatchCommands", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken},
		Run:   commands,
	},
	{
		Name: "edit-info", What: "Bots.EditMyInfo (односторонняя)", Proof: setup.L3,
		Live:  false,
		Needs: []string{setup.KeyToken},
		Force: true,
		Run:   editInfo,
	},
}

func get(_ context.Context, api *maxbot.Api, _ setup.Env, _ []string) error {
	info, err := api.Bots.GetMyInfo(context.Background())
	if err != nil {
		return err
	}
	// user_id печатать не нужно, токен — тем более.
	fmt.Printf("имя=%q username=%q is_bot=%v команд=%d\n",
		info.FirstName, info.Username, info.IsBot, len(info.Commands))
	if info.FirstName == "" || info.UserID == 0 {
		return fmt.Errorf("GetMyInfo вернул пустой профиль")
	}
	fmt.Println("OK: профиль бота получен")

	return nil
}

func commands(_ context.Context, api *maxbot.Api, _ setup.Env, _ []string) error {
	patch := model.BotPatchCommands{Commands: ourCommands}
	res, err := api.Bots.PatchCommands(context.Background(), patch)
	if err != nil {
		return err
	}
	fmt.Printf("запрошено %d команд, ответ вернул %d\n", len(ourCommands), len(res.Commands))

	// Повторный запуск должен давать тот же набор: команды заменяются, а не накапливаются.
	again, err := api.Bots.PatchCommands(context.Background(), patch)
	if err != nil {
		return err
	}
	if len(again.Commands) != len(ourCommands) {
		return fmt.Errorf("после повторной отправки команд %d, ожидалось %d — набор накапливается",
			len(again.Commands), len(ourCommands))
	}
	fmt.Println("OK: идемпотентно — повторная отправка не добавила дублей")

	return nil
}

// editInfo меняет живой профиль бота, поэтому по умолчанию только показывает
// патч; фактическая запись - по -force.
//
// Откат здесь невозможен как класс: MAX требует длину description >= 1
// (PATCH "" -> 400 proto.payload "Field 'description' size (0) must be at
// least 1"), а null сервер игнорирует. Плюс у model.BotPatch.Description стоит
// omitempty, так что пустое исходное значение клиентом не восстанавливается.
// Поэтому команда честно односторонняя и стартует как L3.
func editInfo(_ context.Context, api *maxbot.Api, _ setup.Env, _ []string) error {
	ctx := context.Background()
	before, err := api.Bots.GetMyInfo(ctx)
	if err != nil {
		return err
	}
	want := before.Description + " (example)"
	fmt.Printf("было: описание=%q\nбудет: описание=%q\n", before.Description, want)
	fmt.Println("ВНИМАНИЕ: откатить это нельзя - API не принимает description короче одного")
	fmt.Println("символа и игнорирует null. Патч односторонний, отсюда L3, а не L1.")

	patched := model.BotPatch{FirstName: before.FirstName, Description: want, Commands: before.Commands}
	info, err := api.Bots.EditMyInfo(ctx, patched)
	if err != nil {
		return fmt.Errorf("EditMyInfo: %w", err)
	}

	// Подтверждение перечитыванием, а не доверием к эху ответа.
	after, err := api.Bots.GetMyInfo(ctx)
	if err != nil {
		return fmt.Errorf("GetMyInfo после правки: %w", err)
	}
	if after.Description != want {
		return fmt.Errorf("EditMyInfo вернул %q, перечитано %q", info.Description, after.Description)
	}
	fmt.Printf("применено и перечитано: описание=%q\n", after.Description)
	fmt.Printf("Чтобы вернуть исходное значение, нужна длина >= 1; минимум, которого можно\n")
	fmt.Printf("достичь через клиент, - пробел (сейчас в профиле %q)\n", after.Description)

	return nil
}
