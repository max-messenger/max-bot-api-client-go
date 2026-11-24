package maxbot

import "github.com/max-messenger/max-bot-api-client-go/v3/model"

type Comment struct {
	Text               string           `json:"text"`
	Format             model.TextFormat `json:"format,omitempty"`
	Link               *CommentLink     `json:"link,omitempty"`
	disableLinkPreview bool
}

type CommentLink struct {
	Type model.MessageLinkType `json:"type"`
	Mid  string                `json:"mid"`
}

func NewComment(text string) *Comment {
	return &Comment{
		Text:   text,
		Format: model.FormatMarkdown,
	}
}

func (c *Comment) SetFormat(format model.TextFormat) *Comment {
	c.Format = format

	return c
}
func (c *Comment) SetLink(link CommentLink) *Comment {
	c.Link = &link

	return c
}

func (c *Comment) SetDisableLinkPreview(mode bool) *Comment {
	c.disableLinkPreview = mode

	return c
}
