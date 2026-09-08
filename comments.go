package maxbot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type Comments struct {
	client *client
}

func (c *Comments) GetComments(ctx context.Context, postID string, before, after, count int64, commentIds []string) (res model.CommentList, err error) {
	values := url.Values{}

	if len(commentIds) > 0 {
		values.Set(paramCommentIDs, strings.Join(commentIds, ","))
	}
	if before > 0 {
		values.Set(paramBefore, strconv.FormatInt(before, 10))
	}
	if after > 0 {
		values.Set(paramAfter, strconv.FormatInt(after, 10))
	}
	if count > 0 {
		values.Set(paramCount, strconv.FormatInt(count, 10))
	}
	err = c.client.raw(ctx, http.MethodGet, fmt.Sprintf(formatPathComments, postID), values, nil, &res)

	return
}

func (c *Comments) GetCommentByID(ctx context.Context, postID, commentID string) (res model.Comment, err error) {
	err = c.client.raw(ctx, http.MethodGet, fmt.Sprintf(formatPathCommentByID, postID, commentID), nil, nil, &res)

	return
}

func (c *Comments) Send(ctx context.Context, postID string, comment *Comment) (res model.SendMessageResult, err error) {
	if comment == nil {
		err = fmt.Errorf("nil comment")

		return
	}
	err = c.client.rawWithRetry(ctx, http.MethodPost, fmt.Sprintf(formatPathComments, postID), nil, comment, &res)

	return
}

func (c *Comments) Edit(ctx context.Context, postID, commentID string, comment *Comment) (res model.SendMessageResult, err error) {
	if comment == nil {
		err = fmt.Errorf("nil comment")

		return
	}

	values := url.Values{}
	values.Set(paramCommentID, commentID)

	err = c.client.rawWithRetry(ctx, http.MethodPut, fmt.Sprintf(formatPathComments, postID), values, comment, &res)

	return
}

func (c *Comments) Delete(ctx context.Context, postID, commentID string) (res model.SendMessageResult, err error) {
	values := url.Values{}
	values.Set(paramCommentID, commentID)

	err = c.client.rawWithRetry(ctx, http.MethodDelete, fmt.Sprintf(formatPathComments, postID), values, nil, &res)

	return
}

func newComments(client *client) *Comments {
	return &Comments{
		client: client,
	}
}
