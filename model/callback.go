package model

import (
	"regexp"
)

var callBackReg = regexp.MustCompile(`^(?P<payload>[^:]+?)(?::(?P<param>.+))?$`)

type Callback struct {
	Timestamp  int64  `json:"timestamp"`
	CallbackID string `json:"callback_id"`
	Payload    string `json:"payload"`
	User       User   `json:"user"`
}

type CallbackAnswer struct {
	Message      *NewMessageBody `json:"message,omitempty"`
	Notification *string         `json:"notification,omitempty"`
}

type CallbackPayload struct {
	Payload string
	Param   string
}

func NewCallbackPayload(input string) (res CallbackPayload) {
	if input == "" {
		return
	}

	matches := callBackReg.FindStringSubmatch(input)
	if matches == nil {
		return
	}

	groupNames := callBackReg.SubexpNames()
	result := make(map[string]string, len(groupNames))
	for i, name := range groupNames {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	res.Payload = result["payload"]
	res.Param = result["param"]

	return
}
