// Package setup — общий setup примеров из examples/: окружение, клиент,
// диспетчер подкоманд и печать уровня доказанности.
//
// Роль каждого ключа окружения описана здесь, потому что из имён она не читается.
package setup

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
)

// Уровень доказанности примера. Печатается до запроса, чтобы читатель
// не принял неподготовленный пример за проверенный.
type Proof int

const (
	// L1 — ассерт по ответу живого API.
	L1 Proof = iota
	// L2 — результат виден только глазами (второй аккаунт), не в ответе API.
	L2
	// L3 — воспроизведена только форма запроса (фикстуры + fake-транспорт),
	// поведение реального API не подтверждено.
	L3
)

func (p Proof) String() string {
	switch p {
	case L1:
		return "L1 ассерт по ответу API"
	case L2:
		return "L2 подтверждение глазами (аккаунт B)"
	case L3:
		return "L3 только форма запроса (поведение API не подтверждено)"
	}
	return "L? уровень не заявлен"
}

// Ключи окружения. BOT_TOKEN — единственный секрет; остальные — идентификаторы,
// они не секрет, но без них пример невалидируем.
const (
	KeyToken    = "BOT_TOKEN"      // токен тестового бота
	KeyUserA    = "TEST_USER_ID"   // A: адресат сообщений и владелец тестовой группы
	KeyMember   = "TEST_MEMBER_ID" // B: подопытный AddMembers/RemoveMember и наблюдатель L2
	KeyChat     = "CHAT_ID"        // id тестовой группы, отрицательный int64
	KeyChannel  = "CHAN_ID"        // id тестового публичного канала, отрицательный int64
	DocsBaseURL = "https://platform-api.max.ru"
)

// Env — разобранные идентификаторы окружения.
type Env struct {
	Token   string
	UserA   int64
	Member  int64
	Chat    int64
	Channel int64
}

// need описывает, какие ключи требуются команде, и что делать, если их нет.
var needHelp = map[string]string{
	KeyToken:   "получите токен бота на платформе MAX для партнёров",
	KeyUserA:   "user_id аккаунта A (получатель сообщений, владелец тестовой группы)",
	KeyMember:  "user_id аккаунта B (второй аккаунт: его добавляют/удаляют и с него смотрят)",
	KeyChat:    "создайте группу, добавьте в неё бота и выпишите chat_id (см. examples/subscriptions groups)",
	KeyChannel: "заведите публичный канал, добавьте в него бота и выпишите chat_id канала",
}

// id читает целое из v. Значения chat_id в MAX отрицательные и длинные,
// поэтому разбор строго int64: 32-битный вариант портит валидное значение.
func id(key, v string) (int64, error) {
	if v == "" {
		return 0, fmt.Errorf("%s не задан — %s", key, needHelp[key])
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s=%q — ожидалось целое int64, например -70801090403050 (%s)", key, v, needHelp[key])
	}
	return n, nil
}

// LoadEnv разбирает перечисленные ключи. Возвращает ошибку с именем ключа и
// инструкции по его получению — не панику и не ошибку API.
func LoadEnv(keys ...string) (Env, error) {
	var e Env
	var err error
	for _, k := range keys {
		v := os.Getenv(k)
		switch k {
		case KeyToken:
			if v == "" {
				return e, fmt.Errorf("%s не задан — %s", k, needHelp[k])
			}
			e.Token = v
		case KeyUserA:
			e.UserA, err = id(k, v)
		case KeyMember:
			e.Member, err = id(k, v)
		case KeyChat:
			e.Chat, err = id(k, v)
		case KeyChannel:
			e.Channel, err = id(k, v)
		default:
			return e, fmt.Errorf("неизвестный ключ окружения %s", k)
		}
		if err != nil {
			return e, err
		}
	}
	return e, nil
}

// Client собирает API-клиент. BaseURL пиннится значением из туториалов, а не
// дефолтом библиотеки (v2.4.7 дефолтит на platform-api2.max.ru); пауза и
// таймаут поллинга задаются явно, чтобы примеры не зависали в ожидании человека.
func Client(token string, opts ...maxbot.Opt) (*maxbot.Api, error) {
	opts = append([]maxbot.Opt{
		maxbot.WithBaseURL(DocsBaseURL),
		maxbot.WithPollingPause(200 * time.Millisecond),
		maxbot.WithPollingTimeout(2 * time.Second),
	}, opts...)
	return maxbot.NewApi(token, opts...)
}

// Fixture возвращает путь к файлу фикстуры относительно корня модуля.
// Корень ищется подъёмом по go.mod, потому что CWD различается:
// `go run` ставит его в каталог модуля, `go test` — в каталог пакета,
// а IDE-раннер может поставить куда угодно.
func Fixture(rel string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 6; i++ {
		cand := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(cand); err == nil {
			p := filepath.Join(dir, rel)
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}

			return "", fmt.Errorf("фикстура %s не найдена в корне модуля %s", rel, dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("не найден корень модуля (go.mod) от %s", dir)
}

// Command — одна подкоманда домена.
type Command struct {
	Name  string
	What  string // какой метод клиента показывает
	Proof Proof  // что утверждаем читателю (печатается в usage и перед запуском)

	// Live — факт: по этой команде был живой прогон против реального API.
	// Отделено от Proof намеренно: Proof — заявление, Live — основание для
	// заявления. Тест инвентаря требует Proof == L1 только при Live == true.
	// Поле не печатается: читателю нужен уровень, а не журнал.
	Live bool

	Needs []string
	Force bool // требует -force для исполнения
	Run   func(context.Context, *maxbot.Api, Env, []string) error
}

// Options настраивает выполнение реестра без запуска процесса: куда писать,
// подтверждено ли разрушающее действие и каким клиентом обслуживать команду.
// Нулевые значения дают поведение командной строки: stdout, без -force, живой клиент.
type Options struct {
	Out   io.Writer
	Force bool

	// Client собирает API-клиент по токену. Тесты M-слоя подставляют сюда
	// сборщик с подменным транспортом, чтобы проверять форму запроса без сети.
	Client func(token string) (*maxbot.Api, error)
}

// Result исход выполнения: код возврата и ошибка (nil => успех).
type Result struct {
	Code int
	Err  error
}

// Execute ядро диспетчера: без flag.Parse и без os.Exit, чтобы вызывать его
// можно было из теста. Запускает первую команду с указанным именем.
func Execute(domain string, cmds []Command, argv []string, opts Options) Result {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	if opts.Client == nil {
		opts.Client = func(token string) (*maxbot.Api, error) { return Client(token) }
	}

	usage := func() {
		_, _ = fmt.Fprintf(opts.Out, "запуск: go run ./examples/%s [-force] <команда> [аргументы]\n\nкоманды:\n", domain)
		for _, c := range cmds {
			_, _ = fmt.Fprintf(opts.Out, "  %-14s %s [%s]\n", c.Name, c.What, c.Proof)
		}
	}

	var rest []string

	for _, a := range argv {
		switch a {
		case "-force":
			opts.Force = true
		case "-h", "--help":
			usage()

			return Result{Code: 2}
		default:
			if strings.HasPrefix(a, "-") {
				_, _ = fmt.Fprintf(opts.Out, "неизвестный флаг %q\n\n", a)
				usage()

				return Result{Code: 2}
			}

			rest = append(rest, a)
		}
	}

	if len(rest) == 0 {
		usage()

		return Result{Code: 2}
	}

	name, args := rest[0], rest[1:]
	for _, c := range cmds {
		if c.Name != name {
			continue
		}

		return runOne(opts.Out, opts.Force, opts.Client, domain, c, args)
	}
	_, _ = fmt.Fprintf(opts.Out, "неизвестная команда %q\n\n", name)
	usage()

	return Result{Code: 2}
}

func runOne(out io.Writer, force bool, newClient func(string) (*maxbot.Api, error),
	domain string, c Command, args []string) Result {
	_, _ = fmt.Fprintf(out, "# %s/%s — %s\n# уровень: %s\n", domain, c.Name, c.What, c.Proof)

	// Разрушающее действие отсекается до создания клиента: ни одного запроса
	// не уходит, и это можно проверить тестом.
	if c.Force && !force {
		_, _ = fmt.Fprintf(out, "НЕ ВЫПОЛНЕНО: действие меняет или удаляет среду подготовки; повторите с -force.")

		return Result{}
	}

	env, err := LoadEnv(c.Needs...)
	if err != nil {
		_, _ = fmt.Fprintf(out, "НЕ ВАЛИДИРОВАНО:", err)

		return Result{Code: 1, Err: err}
	}

	api, err := newClient(env.Token)
	if err != nil {
		err = fmt.Errorf("не собран клиент: %w", err)
		_, _ = fmt.Fprintf(out, "НЕ ВАЛИДИРОВАНО:", err)

		return Result{Code: 1, Err: err}
	}

	if err := c.Run(context.Background(), api, env, args); err != nil {
		_, _ = fmt.Fprintf(out, "ошибка:", err)

		return Result{Code: 1, Err: err}
	}

	return Result{}
}

// Run — терминальная обёртка для main: `go run ./examples/<домен> <команда> [аргументы...]`
func Run(domain string, cmds ...Command) {
	if res := RunWith(domain, cmds, os.Args[1:], os.Stdout); res.Code != 0 {
		os.Exit(res.Code)
	}
}

// RunWith — тот же путь, что и Run, но с переданными argv и потоком: тест может
// доказать эквивалентность обёртки и ядра, не подменяя os.Stdout/os.Args.
func RunWith(domain string, cmds []Command, argv []string, out io.Writer) Result {
	return Execute(domain, cmds, argv, Options{Out: out})
}
