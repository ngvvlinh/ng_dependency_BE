package config

import (
	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/integration/shipping/ghn"
	"etop.vn/backend/pkg/integration/shipping/ghtk"
	"etop.vn/common/l"
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
		GHN:  ghn.DefaultConfig(),
		GHTK: ghtk.DefaultConfig(),
		Env:  cm.EnvDev,
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
