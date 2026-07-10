package maxbot

import (
	"fmt"
	"testing"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func TestCommand(t *testing.T) {
	cases := []struct {
		command  string
		expected string
	}{
		{
			command:  "/help",
			expected: "/help",
		},
		{
			command:  "/help:id-773",
			expected: "/help",
		},
		{
			command:  "/help:id-773",
			expected: "/help",
		},
		{
			command:  "/help-me",
			expected: "/help-me",
		},
		{
			command:  "/help/me",
			expected: "/help/me",
		},
		{
			command:  "/help.me",
			expected: "/help.me",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			m := model.Update{
				Message: &model.MessageUpdate{
					Body: model.MessageBody{
						Text: c.command,
					},
				},
			}
			if c.expected != GetCommand(m) {
				t.Errorf("expected %s, got %s", c.expected, GetCommand(m))
			}
		})
	}
}
