package config

import (
	"github.com/jinzhu/configor"
)

type Config struct {
	Slack struct {
		APIToken string `required:"true" env:"SLACK_API_TOKEN"`
		Team     string `required:"true" env:"SLACK_TEAM"`
	}
	Redirect struct {
		FromChannel string `required:"true" env:"REDIRECT_FROM_CHANNEL"`
		ToChannel   string `required:"true" env:"REDIRECT_TO_CHANNEL"`
	}
	Filter struct {
		NgWords []string `env:"FILTER_NG_WORDS"`
		OkUsers []string `env:"FILTER_OK_USERS"`
	}
}

func Load(file string) *Config {
	c := &Config{}
	configor.New(&configor.Config{Silent: true}).Load(c, file)
	return c
}
