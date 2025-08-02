package configservice

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type YamlConfig struct {
	HttpBotAPIUrl             string `yaml:"http.bot.api.url"`
	HttpBotAPITimeOut         int    `yaml:"http.bot.api.timeout"`
	HttpBotAPIVersion         string `yaml:"http.bot.api.version"`
	BotTokenCheckInInputSteam bool   `yaml:"bot.token.check.in.input.stream"`
	BotTokenCheckString       string `yaml:"bot.token.check.string"`

	DebugLogMode bool  `yaml:"debug.log.mode"`
	DebugLogChat int64 `yaml:"debug.log.chat"`
}

type Config struct {
	config YamlConfig
}

func (c *Config) SetEnvVariables(str string) string {
	re := regexp.MustCompile(`(\$\((\w+)\))`)
	res := re.FindAllStringSubmatch(str, -1)
	for _, v := range res {
		if len(v) == 3 {
			if os.Getenv(v[2]) == "" {
				log.Error().Msg("variable " + v[1] + " is not defined!")
			}
			str = strings.ReplaceAll(str, v[1], os.Getenv(v[2]))
		}
	}
	return str
}

func (c *Config) GetHttpBotAPIUrl() string {
	return c.config.HttpBotAPIUrl
}

func (c *Config) GetHttpBotAPITimeOut() int {
	return c.config.HttpBotAPITimeOut
}

func (c *Config) GetHttpBotAPIVersion() string {
	return c.config.HttpBotAPIVersion
}

func (c *Config) BotTokenCheckInInputSteam() bool {
	return c.config.BotTokenCheckInInputSteam
}

func (c *Config) BotTokenCheckString() string {
	return c.config.BotTokenCheckString
}

func (c *Config) GetJsonConfigMarshalled() ([]byte, error) {
	return json.Marshal(c.config)
}

func (c *Config) GetDebugLogMode() bool {
	return c.config.DebugLogMode // debug.log.mode
}

func (c *Config) GetDebugLogChat() int64 {
	return c.config.DebugLogChat // debug.log.chat
}

func (c *Config) WriteJSON(w io.Writer) error {
	js, err := json.Marshal(c.config)
	if err != nil {
		return err
	}

	_, err = w.Write(js)
	return err
}

func (s *Config) loadConfigFromEnv() error {
	return env.Parse(&s.config, env.Options{TagName: "yaml"})
}

func (c *Config) readYamlConfigFile(path string) error {
	filename, err := os.Open(path)
	if err != nil {
		log.Err(err).Msg("readYamlConfigFile os.Open")
		return err
	}
	defer func() {
		err = filename.Close()
		if err != nil {
			log.Err(err).Msg("filename.Close()")
		}
	}()

	source, err := io.ReadAll(filename)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}
	unsource := source

	err = yaml.Unmarshal(unsource, &c.config)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal yaml config")
	}

	return err
}

func (c *Config) readCompositeYamlConfigFile(configPath string) error {
	configPath = c.SetEnvVariables(configPath)
	dir, file := path.Split(configPath)
	ext := path.Ext(file)
	base := strings.TrimSuffix(file, ext)
	baseParts := strings.Split(base, "-")

	for i := 0; i < len(baseParts); i++ {
		composed := dir
		for j := 0; j <= i; j++ {
			if 0 < j {
				composed += "-"
			}
			composed += baseParts[j]
		}
		composed += ext

		if err := c.readYamlConfigFile(composed); err != nil {
			log.Err(err).Msg("ReadCompositeYamlConfigFile")
			return err
		}
	}
	return nil
}
