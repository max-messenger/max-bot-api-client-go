package schemes

import (
	"strconv"
	"testing"
)

func TestMessageCreatedUpdate_GetCommand(t *testing.T) {
	cases := []struct {
		text   string
		expect string
	}{
		{
			text:   "/command:paramA:paramB",
			expect: "/command",
		},
		{
			text:   "/command:paramA:/next",
			expect: "/command",
		},
		{
			text:   "any text",
			expect: CommandUndefined,
		},
		{
			text:   "run/command",
			expect: CommandUndefined,
		},
		{
			text:   "",
			expect: CommandUndefined,
		},
		{
			text:   "/",
			expect: "/",
		},
		{
			text:   "/:run",
			expect: "/",
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := &MessageCreatedUpdate{
				Message: Message{
					Body: MessageBody{
						Text: c.text,
					},
				},
			}

			if c.expect != b.GetCommand() {
				t.Errorf("GetCommand returned %+v, want %+v", b.GetCommand(), c.expect)
			}
		})
	}

}
