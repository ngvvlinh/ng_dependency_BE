package config

import (
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/shipping/ghn"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/common/l"
)

var ll = l.New()

// Config ...
type Config struct {
	Postgres    cc.Postgres    `yaml:"postgres"`
	HTTP        cc.HTTP        `yaml:"http"`
	TelegramBot cc.TelegramBot `yaml:"telegram_bot"`
	GHN         ghn.Config     `yaml:"ghn"`
	GHTK        ghtk.Config    `yaml:"ghtk"`

	Env string `yaml:"env"`
}

// Default ...
func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 9021,
		},
		TelegramBot: cc.TelegramBot{
			Chats: map[string]int64{
				"default": 0,
			},
		},
		GHN:  ghn.DefaultConfig(),
		GHTK: ghtk.DefaultConfig(),
		Env:  cmenv.EnvDev.String(),
	}
	return cfg
}

// Load loads config from file
func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cfg.TelegramBot.MustLoadEnv()
	cfg.GHN.MustLoadEnv()
	cfg.GHTK.MustLoadEnv()
	return cfg, err
}
