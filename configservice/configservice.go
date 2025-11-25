package configservice

import (
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=configservice.go -destination=./mocks/configservice_mock.go -package=mocks
type ConfigInterface interface {
	GetHttpBotAPIUrl() string
	GetHttpBotAPITimeOut() int
	GetHttpBotAPIVersion() string
	BotTokenCheckInInputSteam() bool
	BotTokenCheckString() string

	GetDebugLogMode() bool
	GetDebugLogChat() int64
}

func NewConfigInterface(configPath string) ConfigInterface {
	cs := Config{}
	if err := cs.readCompositeYamlConfigFile(configPath); err != nil {
		log.Err(err).Msg("NewConfigService loadConfigFromYaml")
		return nil
	}

	if err := cs.loadConfigFromEnv(); err != nil {
		log.Err(err).Msg("NewConfigService loadConfigFromEnv")
		return nil
	}
	return &cs
}
