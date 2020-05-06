package config

import (
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	Postgres    cc.Postgres        `yaml:"postgres"`
	HTTP        cc.HTTP            `yaml:"http"`
	TelegramBot cc.TelegramBot     `yaml:"telegram_bot"`
	FacebookApp fbclient.AppConfig `yaml:"facebook_app"`
	TimeLimit   int                `yaml:"time_limit"` // days

	Env string `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8081,
		},
		Env:       cmenv.EnvDev.String(),
		TimeLimit: 3,
	}
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cfg.TelegramBot.MustLoadEnv()
	cfg.FacebookApp.MustLoadEnv()
	return cfg, err
}
