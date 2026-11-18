package maxbot

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestComments(t *testing.T) {
	suite.Run(t, new(commentsTest))
}

type commentsTest struct {
	suite.Suite
}

func (t *commentsTest) SetupTest() {

}

func (t *commentsTest) TestGetCommentByID() {
	data, err := stabs.ReadFile("stabs/comments/get-comment-by-id.json")
	t.NoError(err)

	expect := model.Comment{
		Recipient: model.CommentRecipient{
			ChatID:   -70801090403050,
			ChatType: model.ChatTypeChannel,
			PostID:   "mid.post",
		},
		Timestamp: 1775025604499,
		Body: model.CommentBody{
			CommentID: "mid.comment",
			Seq:       116327994376978687,
			Text:      "comment text",
			MarkUp: []model.Markup{
				{Type: "strong", From: 0, Length: 7},
			},
		},
		Sender: &model.User{
			UserID:           123456789,
			FirstName:        "John",
			LastName:         "Doe",
			IsBot:            false,
			LastActivityTime: 1775025580000,
		},
		Stat: &model.MessageStat{
			Views: 42,
		},
		Link: &model.CommentLinkedMessage{
			Type: model.LinkTypeReply,
			Sender: &model.User{
				UserID:           229229229,
				FirstName:        "UniBot",
				IsBot:            true,
				LastActivityTime: 1775025580000,
				Username:         "unit_bot",
			},
			ChatID: -70801090403050,
			Message: model.CommentBody{
				CommentID: "mid.parent-comment",
				Seq:       116327994376978686,
				Text:      "parent comment",
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodGet)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathCommentByID, "mid.post", "mid.comment"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	res, err := api.Comments.GetCommentByID(context.Background(), "mid.post", "mid.comment")
	t.NoError(err)

	t.Equal(expect, res)
}

func (t *commentsTest) TestSend() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Equal(r.Header.Get(AuthorizationHeader), testToken)
		t.Equal(r.Method, http.MethodPost)
		t.Equal(r.URL.Path, fmt.Sprintf(formatPathComments, "mid.post"))
		t.Equal("true", r.URL.Query().Get(paramDisableLinkPreview))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":{}}`))
	}))

	defer srv.Close()

	api, err := NewApi(testToken, WithBaseURL(srv.URL))
	t.NoError(err)

	comment := NewComment("comment text").SetDisableLinkPreview(true)

	res, err := api.Comments.Send(context.Background(), "mid.post", comment)
	t.NoError(err)

	t.Equal(model.SendMessageResult{}, res)
}
