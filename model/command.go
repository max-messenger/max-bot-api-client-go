package model

import (
	"regexp"
	"strings"
)

var commandReg = regexp.MustCompile(`^(?:@(?P<bot>\S+))?(?:\s+)?(?P<cmd>/\S+?)(?::(?P<params>[^:\s]+))?(?:\s+(?P<text>.+))?$`)

type Command struct {
	BotName       string
	Command       string
	Params        []string
	RemainingText string
}

func NewCommand(input string) (res Command) {
	if input == "" {
		return
	}

	matches := commandReg.FindStringSubmatch(input)
	if matches == nil {
		return
	}

	groupNames := commandReg.SubexpNames()
	result := make(map[string]string, len(groupNames))
	for i, name := range groupNames {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	res.BotName = result["bot"]
	res.Command = result["cmd"]
	res.RemainingText = result["text"]

	if result["params"] != "" {
		res.Params = strings.Split(result["params"], ",")
	}

	return
}
