package schemes

import (
	"strconv"
	"testing"
	"time"
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

func TestUpdate_GetUpdateTime(t *testing.T) {
	cases := []struct {
		timestamp int
		expect    time.Time
	}{
		{
			timestamp: 1739184000000,
			expect:    time.Date(2025, 2, 10, 10, 40, 0, 0, time.UTC).Local(),
		},
		{
			timestamp: int(time.Now().UnixMilli()),
			expect:    time.Now().Truncate(time.Second),
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b := &Update{
				Timestamp: c.timestamp,
			}

			if c.expect != b.GetUpdateTime() {
				t.Errorf("GetUpdateTime returned %+v, want %+v", b.GetUpdateTime(), c.expect)
			}
		})
	}

}
