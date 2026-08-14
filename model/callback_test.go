package model

import (
	"fmt"
	"testing"
)

func TestCallback(tt *testing.T) {
	cases := []struct {
		message string
		payload string
		param   string
	}{
		{
			message: "command_line:param",
			payload: "command_line",
			param:   "param",
		},
		{
			message: "command_line:333",
			payload: "command_line",
			param:   "333",
		},
		{
			message: "command_line",
			payload: "command_line",
		},
		{
			message: "command:param1,param2",
			payload: "command",
			param:   "param1,param2",
		},
		{
			message: "/command:1/2/3",
			payload: "/command",
			param:   "1/2/3",
		},
		{
			message: "some command",
			payload: "some command",
		},
	}

	for i, tc := range cases {
		tt.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			parsed := NewCallbackPayload(tc.message)

			if parsed.Payload != tc.payload {
				t.Errorf("payload want %q, got %q", tc.payload, parsed.Payload)
			}
			if parsed.Param != tc.param {
				t.Errorf("param want %q, got %q", tc.param, parsed.Param)
			}

			// command_line
			// command_line
		})
	}
}
