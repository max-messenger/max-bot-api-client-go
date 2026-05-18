package maxbot

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestClient(t *testing.T) {
	suite.Run(t, new(testClient))
}

type testClient struct {
	suite.Suite
}

func (t *testClient) SetupTest() {

}

func (t *testClient) TestRawOnce() {
	data, err := stabs.ReadFile("stabs/botInfo.ok.json")
	t.NoError(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPost)
		t.Equal(r.URL.Path, pathMe)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()
	api, err := url.Parse(srv.URL)
	t.NoError(err)

	cli := newClient(testToken, api.Host)
	cli.baseURL.Scheme = api.Scheme

	err = cli.raw(context.Background(), http.MethodPost, pathMe, url.Values{}, nil, nil)
	t.NoError(err)
}

func (t *testClient) TestRetry() {
	data, err := stabs.ReadFile("stabs/botInfo.ok.json")
	t.NoError(err)

	dataError, err := stabs.ReadFile("stabs/error.attachment.not.ready.json")
	t.NoError(err)

	var count int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPost)
		t.Equal(r.URL.Path, pathMe)

		if count < 3 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write(dataError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()
	api, err := url.Parse(srv.URL)
	t.NoError(err)

	cli := newClient(testToken, api.Host)
	cli.baseURL.Scheme = api.Scheme

	err = cli.rawWithRetry(context.Background(), http.MethodPost, pathMe, url.Values{}, nil, nil)
	t.NoError(err)
}

func (t *testClient) TestDo() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		t.NoError(err)
		t.Equal(body, []byte("[1,2,3]"))

		w.WriteHeader(http.StatusCreated)
	}))

	defer srv.Close()
	api, err := url.Parse(srv.URL)
	t.NoError(err)

	cli := newClient(testToken, api.Host)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, bytes.NewBufferString(`[1,2,3]`))
	t.NoError(err)

	response, err := cli.do(req)
	t.NoError(err)
	t.Equal(response.StatusCode, http.StatusCreated)
}
