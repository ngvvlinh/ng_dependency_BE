package config

import (
	"o.o/backend/cmd/mc-hub/service/client"
	_telebot "o.o/backend/cogs/base/telebot"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
)

type Config struct {
	InlineConfig `yaml:",inline"`

	Redis cc.Redis `yaml:"redis"`
	HTTP  cc.HTTP  `yaml:"http"`

	cc.TelegramBot `yaml:"telegram_bot"`
	Env            string `yaml:"env"`
}

type InlineConfig struct {
	Endpoints Endpoints `yaml:"endpoints"`
}

type Endpoints struct {
	Main Endpoint `yaml:"main"`
}

type Endpoint = client.EndpointConfig

func Default() Config {
	endpoints := Endpoints{
		Main: Endpoint{BaseURL: "https://api.sandbox.etop.vn"},
	}
	cfg := Config{
		HTTP:        cc.HTTP{Port: 8200},
		Redis:       cc.DefaultRedis(),
		TelegramBot: _telebot.DefaultConfig(),
		Env:         cmenv.EnvDev.String(),

		InlineConfig: InlineConfig{
			Endpoints: endpoints,
		},
	}
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	if err != nil {
		return cfg, err
	}

	cc.RedisMustLoadEnv(&cfg.Redis)
	cfg.TelegramBot.MustLoadEnv()
	return cfg, err
}
