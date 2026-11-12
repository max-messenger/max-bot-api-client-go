// Package offline — герметичная часть примеров: разбор апдейтов и формирование
// запросов без сети и без токена.
//
// Транспорт подставляется через exported seam клиента:
// maxbot.WithHTTPClient接受ет интерфейс из одного метода Do
// (client.go:24), поэтому запросы реального клиента перехватываются и обслуживаются
// фикстурами из testdata/. Уровень доказанности — L3: подтверждается форма запроса
// и разбор ответа, но не поведение живого API.
package offline

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Stub отвечает на один запрос. attempt считается отдельно для каждого
// method+path и начинается с 1 — этого хватает показать ретрай клиента.
type Stub func(attempt int) (status int, body string)

// Transport — реализация maxbot.HttpClient, отвечающая из памяти.
type Transport struct {
	stubs map[string]Stub
	at    map[string]int
	// Requests записывает method path?query всех запросов — так пример
	// показывает форму запроса, а не только результат разбора.
	Requests []string
}

func NewTransport(stubs map[string]Stub) *Transport {
	return &Transport{stubs: stubs, at: map[string]int{}}
}

func (t *Transport) Do(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.Path
	t.at[key]++
	line := req.Method + " " + req.URL.Path
	if req.URL.RawQuery != "" {
		line += "?" + req.URL.RawQuery
	}
	t.Requests = append(t.Requests, line)

	stub, ok := t.stubs[key]
	if !ok {
		return nil, fmt.Errorf("offline: нет заглушки для %s", key)
	}
	status, body := stub(t.at[key])

	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

// Attempts возвращает число обращений к method+path.
func (t *Transport) Attempts(method, path string) int {
	return t.at[method+" "+path]
}
