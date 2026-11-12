package setup_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/internal/examples/setup"
)

// stubTransport обслуживает ответы без сети, чтобы тест ядра не зависел от API.
type stubTransport struct{ requests []string }

func (s *stubTransport) Do(req *http.Request) (*http.Response, error) {
	s.requests = append(s.requests, req.Method+" "+req.URL.Path)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"success":true}`)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

// newClient собирает клиент с подменным транспортом и поднимает called, когда
// клиента вообще попросили создать: так видно, дошло ли исполнение до тела команды.
func newClient(tr *stubTransport, called *bool) func(string) (*maxbot.Api, error) {
	return func(token string) (*maxbot.Api, error) {
		if called != nil {
			*called = true
		}

		return maxbot.NewApi(token, maxbot.WithHTTPClient(tr))
	}
}

func testCmds(_ *stubTransport, clientCalled *bool) []setup.Command {
	noop := func(context.Context, *maxbot.Api, setup.Env, []string) error { return nil }

	return []setup.Command{
		{Name: "ok", What: "Bots.GetMyInfo", Proof: setup.L1,
			Needs: []string{setup.KeyToken}, Run: noop},
		{Name: "needs-chat", What: "Chats.GetChat", Proof: setup.L1,
			Needs: []string{setup.KeyToken, setup.KeyChat}, Run: noop},
		{Name: "boom", What: "Messages.Send", Proof: setup.L3,
			Needs: []string{setup.KeyToken},
			Run: func(context.Context, *maxbot.Api, setup.Env, []string) error {
				return errors.New("намеренная ошибка команды")
			}},
		{Name: "wipe", What: "Chats.DeleteChat", Proof: setup.L3, Force: true,
			Needs: []string{setup.KeyToken},
			Run: func(context.Context, *maxbot.Api, setup.Env, []string) error {
				if clientCalled != nil {
					*clientCalled = true
				}

				return errors.New("тело разрушающей команды не должно было исполниться")
			}},
	}
}

// TestExecuteNeverExits доказывает главное: Execute возвращает исход, а не
// завершает процесс-хост. Если бы здесь остался os.Exit, тест до последней
// ассерции не дожил бы.
func TestExecuteNeverExits(t *testing.T) {
	// OK-команды требуют BOT_TOKEN: без него они упёрлись бы в диагностику
	// окружения вместо проверяемых веток диспетчера.
	t.Setenv(setup.KeyToken, "token")

	tr := &stubTransport{}
	var out bytes.Buffer
	cmds := testCmds(tr, nil)

	for _, tc := range []struct {
		name string
		argv []string
		want int
	}{
		{"неизвестная команда", []string{"nope"}, 2},
		{"пустой argv", nil, 2},
		{"неизвестный флаг", []string{"-x"}, 2},
		{"падающая команда", []string{"boom"}, 1},
	} {
		out.Reset()
		res := setup.Execute("d", cmds, tc.argv, setup.Options{Out: &out, Client: newClient(tr, nil)})
		if res.Code != tc.want {
			t.Errorf("%s: код %d, ожидался %d (вывод %q)", tc.name, res.Code, tc.want, out.String())
		}
	}

	if !strings.Contains(out.String(), "ошибка:") {
		t.Errorf("вывод падающей команды не содержит «ошибка:»: %q", out.String())
	}
	t.Log("процесс жив, хотя четырежды был запрошен ненулевой код")
}

// TestExecuteDiagnostics проверяет текст вместо паники при отсутствии окружения.
func TestExecuteDiagnostics(t *testing.T) {
	t.Setenv(setup.KeyToken, "token")
	t.Setenv(setup.KeyChat, "")

	var out bytes.Buffer
	res := setup.Execute("d", testCmds(&stubTransport{}, nil), []string{"needs-chat"},
		setup.Options{Out: &out, Client: newClient(nil, nil)})

	if res.Code != 1 {
		t.Fatalf("код %d, ожидался 1", res.Code)
	}
	for _, want := range []string{"НЕ ВАЛИДИРОВАНО", setup.KeyChat, "добавьте в неё бота"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("в выводе нет %q:\n%s", want, out.String())
		}
	}
}

// TestForceGateDoesNotTouchEnvironment - ключевой ассерт защиты среды: без
// подтверждения тело разрушающей команды не вызывается. Отказ наступает до
// создания клиента, поэтому ни один запрос не может уйти.
func TestForceGateDoesNotTouchEnvironment(t *testing.T) {
	clientCalled := false
	tr := &stubTransport{}
	var out bytes.Buffer

	res := setup.Execute("d", testCmds(tr, &clientCalled), []string{"wipe"},
		setup.Options{Out: &out, Client: newClient(tr, &clientCalled)})

	if res.Code != 0 || res.Err != nil {
		t.Errorf("отказ обязан быть кодом 0 без ошибки, получили %+v", res)
	}
	if clientCalled {
		t.Error("тело команды вызвано - гейт не сработал")
	}
	if len(tr.requests) != 0 {
		t.Errorf("разрушающая команда без -force отправила запросы: %v", tr.requests)
	}
	if !strings.Contains(out.String(), "НЕ ВЫПОЛНЕНО") {
		t.Errorf("вывод не содержит «НЕ ВЫПОЛНЕНО»: %q", out.String())
	}
}

// TestRunWithEqualsExecute закрепляет инвариант CLI: обёртка даёт тот же вывод и
// тот же код, что и ядро с теми же аргументами.
func TestRunWithEqualsExecute(t *testing.T) {
	t.Setenv(setup.KeyToken, "token")
	tr := &stubTransport{}
	cmds := testCmds(tr, nil)

	var wantBuf, gotBuf bytes.Buffer
	want := setup.Execute("d", cmds, []string{"ok"}, setup.Options{Out: &wantBuf, Client: newClient(tr, nil)})
	got := setup.RunWith("d", cmds, []string{"ok"}, &gotBuf)

	if want.Code != 0 || want.Err != nil {
		t.Fatalf("базовый Execute: %+v", want)
	}
	if got.Code != want.Code {
		t.Errorf("код RunWith %d != Execute %d", got.Code, want.Code)
	}
	if gotBuf.String() != wantBuf.String() {
		t.Errorf("вывод разошёлся:\nRunWith:  %q\nExecute: %q", gotBuf.String(), wantBuf.String())
	}

	// -force разбирается из argv так же, как флаг в прежнем flag.Parse.
	var forceBuf bytes.Buffer
	if res := setup.Execute("d", cmds, []string{"-force", "ok"},
		setup.Options{Out: &forceBuf, Client: newClient(tr, nil)}); res.Code != 0 {
		t.Errorf("-force в argv: %+v", res)
	}
	var forceBuf2 bytes.Buffer
	if res := setup.Execute("d", cmds, []string{"ok"},
		setup.Options{Out: &forceBuf2, Force: true, Client: newClient(tr, nil)}); res.Code != 0 {
		t.Errorf("Force через Options: %+v", res)
	}
}

// TestUsageListsCommands проверяет, что список команд печатается для домена без
// аргументов и содержит каждую команду с её уровнем.
func TestUsageListsCommands(t *testing.T) {
	var out bytes.Buffer
	if res := setup.Execute("d", testCmds(&stubTransport{}, nil), nil, setup.Options{Out: &out}); res.Code != 2 {
		t.Fatalf("код %d, ожидался 2", res.Code)
	}
	for _, name := range []string{"ok", "needs-chat", "boom", "wipe"} {
		if !strings.Contains(out.String(), name) {
			t.Errorf("в usage нет команды %q", name)
		}
	}
	if !strings.Contains(out.String(), "L1 ассерт") || !strings.Contains(out.String(), "L3 только форма") {
		t.Error("usage не печатает уровни доказанности")
	}
}
