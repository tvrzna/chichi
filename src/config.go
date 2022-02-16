package src

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/tvrzna/go-utils/args"
	"github.com/tvrzna/go-utils/config"
)

type Config struct {
	NotifySendPath string
	ShortPeriod    int `config:"SHORT_PERIOD" parser:"ParseInt" default:"600"`
	ShortBreak     int `config:"SHORT_BREAK" parser:"ParseInt" default:"10"`
	LongPeriod     int `config:"LONG_PERIOD" parser:"ParseInt" default:"2100"`
	LongBreak      int `config:"LONG_BREAK" parser:"ParseInt" default:"60"`
}

func (c *Config) ParseInt(value, defaultValue string) (result int) {
	if value == "" {
		result, _ = strconv.Atoi(defaultValue)
	} else {
		result, _ = strconv.Atoi(strings.TrimSpace(value))
	}
	return result
}

func (c *Config) ParseIntWithIntDefault(value string, defaultValue int) (result int) {
	if value == "" {
		result = defaultValue
	} else {
		result, _ = strconv.Atoi(strings.TrimSpace(value))
	}
	return result
}

func loadConfig(path string, arguments []string) (*Config, error) {
	c := &Config{}
	config.LoadConfigFromFile(c, path, true)

	args.ParseArgs(arguments, func(value, nextValue string) {
		switch value {
		case "-sp", "--short-period":
			c.ShortPeriod = c.ParseIntWithIntDefault(nextValue, c.ShortPeriod)
		case "-sb", "--short-break":
			c.ShortBreak = c.ParseIntWithIntDefault(nextValue, c.ShortBreak)
		case "-lp", "--long-period":
			c.LongPeriod = c.ParseIntWithIntDefault(nextValue, c.LongPeriod)
		case "-lb", "--long-break":
			c.LongBreak = c.ParseIntWithIntDefault(nextValue, c.LongBreak)
		}
	})

	path, err := exec.LookPath("notify-send")
	c.NotifySendPath = path

	return c, err
}
