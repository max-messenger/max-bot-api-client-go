package maxbot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestBots(t *testing.T) {
	suite.Run(t, new(botsTest))
}

type botsTest struct {
	suite.Suite
}

func (t *botsTest) SetupTest() {

}

func (t *botsTest) TestInfoSuccess() {
	data, err := stabs.ReadFile("stabs/botInfo.ok.json")
	t.NoError(err)

	info := model.BotInfo{
		UserID:           123123123,
		FirstName:        "unit bot",
		Username:         "test-bot",
		IsBot:            true,
		LastActivityTime: 1774677196913,
		Description:      "bot for testing",
		AvatarURL:        "https://localhost/i?r=hash1",
		FullAvatarURL:    "https://localhost/i?r=hash1",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, pathMe)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Bots.GetMyInfo(context.Background())
	t.NoError(err)

	t.Equal(info, res)
}

func (t *botsTest) TestInfoError() {
	data, err := stabs.ReadFile("stabs/error.invalid-token.json")
	t.NoError(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, pathMe)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	_, err = api.Bots.GetMyInfo(context.Background())
	t.EqualError(err, "GetMyInfo: verify.token : Invalid access_token")
}

func (t *botsTest) TestPathSuccess() {
	data, err := stabs.ReadFile("stabs/botInfo.ok.json")
	t.NoError(err)

	expect := model.BotInfo{
		UserID:           123123123,
		FirstName:        "unit bot",
		Username:         "test-bot",
		IsBot:            true,
		LastActivityTime: 1774677196913,
		Description:      "bot for testing",
		AvatarURL:        "https://localhost/i?r=hash1",
		FullAvatarURL:    "https://localhost/i?r=hash1",
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPatch)
		t.Equal(r.URL.Path, pathMe)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Bots.EditMyInfo(context.Background(), model.BotPatch{FirstName: "unit bot"})
	t.NoError(err)

	t.Equal(expect, res)
}

func (t *botsTest) TestPathError() {
	data, err := stabs.ReadFile("stabs/error.invalid-token.json")
	t.NoError(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPatch)
		t.Equal(r.URL.Path, pathMe)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	_, err = api.Bots.EditMyInfo(context.Background(), model.BotPatch{FirstName: "unit bot"})
	t.EqualError(err, "EditMyInfo: verify.token : Invalid access_token")
}
