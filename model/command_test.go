package model

import (
	"fmt"
	"testing"
)

func TestCommand(tt *testing.T) {
	cases := []struct {
		message string
		botName string
		command string
		params  []string
		text    string
	}{
		{
			message: "@botName/command_line:param1 aaa",
			botName: "botName",
			command: "/command_line",
			params:  []string{"param1"},
			text:    "aaa",
		},
		{
			message: "@botName/command_line:param1,param2 aaa",
			botName: "botName",
			command: "/command_line",
			params:  []string{"param1", "param2"},
			text:    "aaa",
		},
		{
			message: "@botName        /command_line:param1,param2 aaa",
			botName: "botName",
			command: "/command_line",
			params:  []string{"param1", "param2"},
			text:    "aaa",
		},
		{
			message: "/command_line:param1,param2 aaa",
			command: "/command_line",
			params:  []string{"param1", "param2"},
			text:    "aaa",
		},
		{
			message: "/command_line aaa",
			command: "/command_line",
			text:    "aaa",
		},
		{
			message: "@botName /command aaa",
			botName: "botName",
			command: "/command",
			text:    "aaa",
		},
		{
			message: "/command",
			command: "/command",
		},
		{
			message: "@botName /command:param1,param2",
			botName: "botName",
			command: "/command",
			params:  []string{"param1", "param2"},
		},
		{
			message: "/command:param1,param2",
			command: "/command",
			params:  []string{"param1", "param2"},
		},
		{
			message: "invalid command",
		},
	}

	for i, tc := range cases {
		tt.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			parsed := NewCommand(tc.message)

			if parsed.BotName != tc.botName {
				t.Errorf("botName want %q, got %q", tc.botName, parsed.BotName)
			}
			if parsed.Command != tc.command {
				t.Errorf("command want %q, got %q", tc.command, parsed.Command)
			}
			if len(parsed.Params) != len(tc.params) {
				t.Errorf("params want %d (%v), got %d (%v)", len(tc.params), tc.params, len(parsed.Params), parsed.Params)
			} else {
				for j := range parsed.Params {
					if parsed.Params[j] != tc.params[j] {
						t.Errorf("params[%d] want %q, got %q", j, tc.params[j], parsed.Params[j])
					}
				}
			}
			if parsed.RemainingText != tc.text {
				t.Errorf("remaining text want %q, got %q", tc.text, parsed.RemainingText)
			}
		})
	}
}
