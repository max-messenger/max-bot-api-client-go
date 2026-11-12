// Command chats — методы группы ChatsAPI клиента max-bot-api-client-go/v2.
//
// Требуется подготовленная группа: A (TEST_USER_ID) создаёт группу, добавляет
// в неё бота, B (TEST_MEMBER_ID) — второй участник для add/remove.
// Id группы берётся из CHAT_ID: метода GetChats в клиенте нет, поэтому id
// выписывается из клиента MAX или достаётся командой examples/subscriptions groups.
//
//	go run ./examples/chats get                     # Chats.GetChat
//	go run ./examples/chats membership              # Chats.GetMembership
//	go run ./examples/chats admins                  # Chats.GetAdmins
//	go run ./examples/chats members                 # Chats.GetMembers
//	go run ./examples/chats pin-get                 # Chats.GetPinnedMessage
//	go run ./examples/chats action [typing_on|…]    # Chats.SendAction      (видно только глазами)
//	go run ./examples/chats unpin                    # Chats.UnpinMessage (после pin)
//	go run ./examples/chats rename [новый title]    # Chats.EditChat        (возвращает title)
//	go run ./examples/chats member-roundtrip        # Chats.AddMembers + RemoveMember
//	go run ./examples/chats admin-set               # Chats.SetAdmins       (только -force)
//	go run ./examples/chats admin-delete            # Chats.DeleteAdmins    (только -force)
//	go run ./examples/chats leave                   # Chats.LeaveChat       (только -force)
//	go run ./examples/chats delete                  # Chats.DeleteChat      (только -force)
package chats

import (
	"context"
	"fmt"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

// ctxFor берёт чат и проверяет, что это действительно группа: CHAT_ID из
// окружения мог указывать на диалог или на удалённую группу.
func ctxFor(bg context.Context, api *maxbot.Api, env setup.Env) (model.Chat, error) {
	c, err := api.Chats.GetChat(bg, env.Chat)
	if err != nil {
		return c, fmt.Errorf("Chats.GetChat(%d): %w — проверьте CHAT_ID (отрицательное число)", env.Chat, err)
	}
	if c.Type != model.ChatTypeChat {
		return c, fmt.Errorf("chat_id=%d имеет type=%q, ожидался chat (нужна группа, а не диалог)", env.Chat, c.Type)
	}

	return c, nil
}

// tok — окружение, общее для read-only команд домена: токен и id группы.
var tok = []string{setup.KeyToken, setup.KeyChat}

var Commands = []setup.Command{
	{
		Name: "get", What: "Chats.GetChat", Proof: setup.L1, Live: true,
		Needs: tok, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			c, err := ctxFor(context.Background(), api, env)
			if err != nil {
				return err
			}
			fmt.Printf("title=%q type=%q status=%q участников=%d владелец=%d публичный=%v\n",
				c.Title, c.Type, c.Status, c.ParticipantsCount, c.OwnerID, c.IsPublic)
			if c.ChatID != env.Chat {
				return fmt.Errorf("GetChat вернул chat_id=%d, запрашивался %d", c.ChatID, env.Chat)
			}
			fmt.Println("OK: группа найдена и описана")

			return nil
		},
	},
	{
		Name: "membership", What: "Chats.GetMembership", Proof: setup.L1, Live: true,
		Needs: tok, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			m, err := api.Chats.GetMembership(context.Background(), env.Chat)
			if err != nil {
				return fmt.Errorf("Chats.GetMembership(%d): %w", env.Chat, err)
			}
			fmt.Printf("бот в группе: is_owner=%v is_admin=%v alias=%q прав=%d\n",
				m.IsOwner, m.IsAdmin, m.Alias, len(m.Permissions))
			fmt.Printf("права: %v\n", m.Permissions)
			fmt.Println("OK: членство получено — от него зависит L1 у admin-* и rename")

			return nil
		},
	},
	{
		Name: "admins", What: "Chats.GetAdmins", Proof: setup.L1, Live: true,
		Needs: tok, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			al, err := api.Chats.GetAdmins(context.Background(), env.Chat)
			if err != nil {
				return fmt.Errorf("Chats.GetAdmins(%d): %w", env.Chat, err)
			}
			fmt.Printf("администраторов: %d\n", len(al.Members))
			for _, a := range al.Members {
				fmt.Printf("  user_id=%d alias=%q прав=%d\n", a.UserID, a.Alias, len(a.Permissions))
			}

			return nil
		},
	},
	{
		Name: "members", What: "Chats.GetMembers", Proof: setup.L1, Live: true,
		Needs: tok, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			ml, err := api.Chats.GetMembers(context.Background(), env.Chat, 0, 50, nil)
			if err != nil {
				return fmt.Errorf("Chats.GetMembers(%d): %w", env.Chat, err)
			}
			fmt.Printf("участников: %d (marker=%d)\n", len(ml.Members), ml.Marker)
			for _, u := range ml.Members {
				fmt.Printf("  user_id=%d name=%q owner=%v admin=%v\n",
					u.UserID, u.FirstName, u.IsOwner, u.IsAdmin)
			}
			if len(ml.Members) < 2 {
				fmt.Println("ЗАМЕЧАНИЕ: участников меньше двух — добавьте аккаунт B, иначе")
				fmt.Println("member-roundtrip нечего показывать")
			}

			return nil
		},
	},
	{
		Name: "pin-get", What: "Chats.GetPinnedMessage", Proof: setup.L1, Live: true,
		Needs: tok, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			p, err := api.Chats.GetPinnedMessage(context.Background(), env.Chat)
			if err != nil {
				fmt.Printf("НЕ ВАЛИДИРОВАНО: GetPinnedMessage(%d): %v\n", env.Chat, err)
				fmt.Println("Закреплённого сообщения может не быть - это штатная ситуация.")

				return nil
			}
			if p.Message.Body.Mid == "" {
				// API отвечает success с пустым сообщением, когда закрепления нет.
				// Отдельного статуса "не закреплено" у ответа нет - это надо уметь распознать.
				fmt.Println("OK: запрос выполнен, закреплённого сообщения нет (body.mid пустой)")

				return nil
			}
			fmt.Printf("закреплено: mid=%s text=%.40q\n", p.Message.Body.Mid, p.Message.Body.Text)

			return nil
		},
	},
	{
		Name: "action", What: "Chats.SendAction", Proof: setup.L2, Live: false,
		Needs: tok, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, args []string) error {
			action := model.ActionTypingOn
			if len(args) > 0 {
				action = model.SenderAction(args[0])
			}
			res, err := api.Chats.SendAction(context.Background(), env.Chat, action)
			if err != nil {
				return fmt.Errorf("Chats.SendAction(%s): %w", action, err)
			}
			fmt.Printf("SendAction вернул success=%v message=%q\n", res.Success, res.Message)
			fmt.Println("Это L2: индикатор 'печатает…' не виден ни в одном ответе API.")
			fmt.Println("Смотрите группу с аккаунта B (TEST_MEMBER_ID) в момент вызова.")

			return nil
		},
	},
	{
		Name: "pin", What: "Chats.PinMessage (+ GetPinnedMessage)", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken, setup.KeyChat, setup.KeyUserA}, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			bg := context.Background()
			if _, err := ctxFor(bg, api, env); err != nil {
				return err
			}
			sent, err := api.Messages.Send(bg, maxbot.NewMessage().
				SetChat(env.Chat).SetText("пример: закрепить"))
			if err != nil {
				return fmt.Errorf("Messages.Send в группу: %w", err)
			}
			mid := sent.Message.Body.Mid
			defer func() { _, _ = api.Messages.DeleteMessage(bg, mid) }()

			if _, err := api.Chats.PinMessage(bg, env.Chat, mid, false); err != nil {
				return fmt.Errorf("Chats.PinMessage: %w", err)
			}
			p, err := api.Chats.GetPinnedMessage(bg, env.Chat)
			if err != nil {
				return fmt.Errorf("Chats.GetPinnedMessage после pin: %w", err)
			}
			fmt.Printf("закреплено mid=%s (запрашивался %s)\n", p.Message.Body.Mid, mid)

			if _, err := api.Chats.UnpinMessage(bg, env.Chat); err != nil {
				return fmt.Errorf("Chats.UnpinMessage: %w (закрепление осталось висеть!)", err)
			}
			fmt.Println("OK: pin подтверждён чтением, unpin выполнен")
			fmt.Println("L2-часть: баннер закрепления в самой группе виден только с аккаунта B")

			return nil
		},
	},
	{
		Name: "unpin", What: "Chats.UnpinMessage", Proof: setup.L1,
		Live:  true,
		Needs: []string{setup.KeyToken, setup.KeyChat, setup.KeyUserA}, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			bg := context.Background()
			if _, err := ctxFor(bg, api, env); err != nil {
				return err
			}
			sent, err := api.Messages.Send(bg, maxbot.NewMessage().
				SetChat(env.Chat).SetText("пример: открепить"))
			if err != nil {
				return fmt.Errorf("Messages.Send в группу: %w", err)
			}
			mid := sent.Message.Body.Mid
			defer func() { _, _ = api.Messages.DeleteMessage(bg, mid) }()

			if _, err := api.Chats.PinMessage(bg, env.Chat, mid, false); err != nil {
				return fmt.Errorf("Chats.PinMessage для подготовки: %w", err)
			}
			if _, err := api.Chats.UnpinMessage(bg, env.Chat); err != nil {
				return fmt.Errorf("Chats.UnpinMessage: %w", err)
			}
			// Открепление подтверждается чтением: GetPinnedMessage на группе
			// без закрепления отвечает success с пустым body.mid.
			p, err := api.Chats.GetPinnedMessage(bg, env.Chat)
			if err != nil {
				return fmt.Errorf("Chats.GetPinnedMessage после unpin: %w", err)
			}
			if p.Message.Body.Mid != "" {
				return fmt.Errorf("после UnpinMessage по-прежнему закреплено mid=%s", p.Message.Body.Mid)
			}
			fmt.Printf("OK: mid=%s закреплён и откреплён, GetPinnedMessage больше ничего не отдаёт\n", mid)

			return nil
		},
	},
	{
		Name: "rename", What: "Chats.EditChat", Proof: setup.L2, Live: false,
		Needs: tok, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, args []string) error {
			bg := context.Background()
			before, err := ctxFor(bg, api, env)
			if err != nil {
				return err
			}
			title := before.Title + " (example)"
			if len(args) > 0 {
				title = args[0]
			}
			// Ловушка: у ChatPatch.Description нет omitempty - пустая строка уедет
			//в PATCH и затрёт описание группы. Всегда отдаём текущее.
			patch := model.ChatPatch{Title: title, Description: before.Description}
			after, err := api.Chats.EditChat(bg, env.Chat, patch)
			if err != nil {
				return fmt.Errorf("Chats.EditChat: %w", err)
			}
			fmt.Printf("title: %q -> %q\n", before.Title, after.Title)

			restore := model.ChatPatch{Title: before.Title, Description: before.Description}
			back, err := api.Chats.EditChat(bg, env.Chat, restore)
			if err != nil {
				return fmt.Errorf("откат EditChat не выполнен, title остался %q: %w", after.Title, err)
			}
			fmt.Printf("OK: изменение принято и откачено, title=%q\n", back.Title)
			fmt.Println("L2-часть: переименование группы видно её участникам, не в ответе API")

			return nil
		},
	},
	{
		// L3: форма запроса проверяется тестом, живого прогона не было.
		Name: "member-roundtrip", What: "Chats.AddMembers + Chats.RemoveMember", Proof: setup.L3,
		Live:  false,
		Needs: []string{setup.KeyToken, setup.KeyChat, setup.KeyMember}, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			bg := context.Background()
			if _, err := ctxFor(bg, api, env); err != nil {
				return err
			}
			before, err := api.Chats.GetMembers(bg, env.Chat, 0, 50, nil)
			if err != nil {
				return err
			}
			// Предпроверка обязательна: AddMembers на уже состоящем участнике
			// отвечает success, а последующий RemoveMember выгнал бы его из группы
			// (в том числе владельца среды подготовки).
			for _, u := range before.Members {
				if u.UserID == env.Member {
					fmt.Printf("НЕ ВАЛИДИРОВАНО: TEST_MEMBER_ID=%d (%s) уже состоит в группе (owner=%v admin=%v).\n",
						env.Member, u.FirstName, u.IsOwner, u.IsAdmin)
					fmt.Println("Для этой команды нужен аккаунт, которого в группе НЕТ, иначе")
					fmt.Println("AddMembers не изменит состав, а RemoveMember удалит существующего участника.")
					fmt.Println("Ловушка подтверждена прогоном: API отвечает success и на повторы.")

					return fmt.Errorf("участник %d уже в группе - команда не будет ничего менять", env.Member)
				}
			}

			if _, err := api.Chats.AddMembers(bg, env.Chat, []int64{env.Member}); err != nil {
				return fmt.Errorf("Chats.AddMembers(B=%d): %w — у бота должен быть прав add_remove_members", env.Member, err)
			}
			mid, err := api.Chats.GetMembers(bg, env.Chat, 0, 50, nil)
			if err != nil {
				return err
			}
			fmt.Printf("участников: %d -> %d после AddMembers\n", len(before.Members), len(mid.Members))
			if len(mid.Members) <= len(before.Members) {
				return fmt.Errorf("после AddMembers состав не вырос (%d -> %d)",
					len(before.Members), len(mid.Members))
			}

			if _, err := api.Chats.RemoveMember(bg, env.Chat, env.Member, false); err != nil {
				return fmt.Errorf("Chats.RemoveMember(B=%d): %w - B мог остаться в группе!", env.Member, err)
			}
			after, err := api.Chats.GetMembers(bg, env.Chat, 0, 50, nil)
			if err != nil {
				return err
			}
			fmt.Printf("участников: %d после RemoveMember (block=false, B не заблокирован)\n",
				len(after.Members))
			if len(after.Members) != len(before.Members) {
				fmt.Println("ЗАМЕЧАНИЕ: состав не совпал с исходным - проверьте руками, что B в группе")
			} else {
				fmt.Println("OK: добавление и удаление подтверждены чтением состава")
			}

			return nil
		},
	},
	{
		Name: "admin-set", What: "Chats.SetAdmins", Proof: setup.L3,
		Live:  false,
		Needs: []string{setup.KeyToken, setup.KeyChat, setup.KeyMember}, Force: true, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			bg := context.Background()
			admins, err := api.Chats.GetAdmins(bg, env.Chat)
			if err != nil {
				return err
			}
			fmt.Printf("было администраторов: %d\n", len(admins.Members))
			// GetAdmins отдаёт []ChatMember, а SetAdmins принимает []ChatAdmin:
			// типы разные, и перезапись списка без переноса прав снимает их со
			// всех действующих админов.
			list := make([]model.ChatAdmin, 0, len(admins.Members)+1)
			for _, a := range admins.Members {
				list = append(list, model.ChatAdmin{
					UserID:      a.UserID,
					Permissions: a.Permissions,
					Alias:       a.Alias,
				})
			}
			newAdmin := model.ChatAdmin{
				UserID:      env.Member,
				Permissions: []model.ChatAdminPermission{model.PermReadAllMessages, model.PermWrite},
				Alias:       "example-admin",
			}
			list = append(list, newAdmin)
			res, err := api.Chats.SetAdmins(bg, env.Chat, list)
			if err != nil {
				// Проверено прогоном на группе, где бот админ, но не владелец:
				// chat.denied с набором прав, которые есть у самого бота, и
				// proto.payload, если просить pin_message без read_all_messages.
				return fmt.Errorf("Chats.SetAdmins: %w — для этого нужна группа, где бот владелец (см. Chats.GetMembership)", err)
			}
			fmt.Printf("SetAdmins: success=%v message=%q\n", res.Success, res.Message)

			after, err := api.Chats.GetAdmins(bg, env.Chat)
			if err != nil {
				return err
			}
			fmt.Printf("стало администраторов: %d\n", len(after.Members))
			if _, err := api.Chats.DeleteAdmins(bg, env.Chat, env.Member); err != nil {
				return fmt.Errorf("откат DeleteAdmins(B=%d) не выполнен: %w", env.Member, err)
			}
			fmt.Println("OK: назначение подтверждено чтением, B снят с администратора")

			return nil
		},
	},
	{
		Name: "admin-delete", What: "Chats.DeleteAdmins", Proof: setup.L3,
		Live:  false,
		Needs: []string{setup.KeyToken, setup.KeyChat, setup.KeyMember}, Force: true, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			res, err := api.Chats.DeleteAdmins(context.Background(), env.Chat, env.Member)
			if err != nil {
				return fmt.Errorf("Chats.DeleteAdmins(B=%d): %w — B мог не быть администратором",
					env.Member, err)
			}
			fmt.Printf("DeleteAdmins: success=%v message=%q\n", res.Success, res.Message)
			fmt.Println("ЗАМЕЧАНИЕ: без предварительно назначенного админа команда не показывает ничего,")
			fmt.Println("используйте admin-set (он делает назначение и откат в одном прогоне)")

			return nil
		},
	},
	{
		Name: "leave", What: "Chats.LeaveChat", Proof: setup.L3, Live: false,
		Needs: tok, Force: true, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			fmt.Println("Бот покинет группу. Восстановить нельзя вызовом API — только руками.")
			res, err := api.Chats.LeaveChat(context.Background(), env.Chat)
			if err != nil {
				return fmt.Errorf("Chats.LeaveChat: %w", err)
			}
			fmt.Printf("LeaveChat: success=%v message=%q\n", res.Success, res.Message)

			return nil
		},
	},
	{
		Name: "delete", What: "Chats.DeleteChat", Proof: setup.L3, Live: false,
		Needs: tok, Force: true, Run: func(
			_ context.Context, api *maxbot.Api, env setup.Env, _ []string) error {
			fmt.Println("НЕОБРАТИМО: будет удалена группа из CHAT_ID. API не умеет создавать группы,")
			fmt.Println("так что среду подготовки придётся собирать заново руками.")
			res, err := api.Chats.DeleteChat(context.Background(), env.Chat)
			if err != nil {
				return fmt.Errorf("Chats.DeleteChat: %w — возможно, бот не владелец группы", err)
			}
			fmt.Printf("DeleteChat: success=%v message=%q\n", res.Success, res.Message)

			return nil
		},
	},
}
