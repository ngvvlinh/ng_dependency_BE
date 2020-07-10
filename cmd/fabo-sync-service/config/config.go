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
	Redis       cc.Redis           `yaml:"redis"`
	FacebookApp fbclient.AppConfig `yaml:"facebook_app"`
	TimeLimit   int                `yaml:"time_limit"`    // days
	TimeToCrawl int                `yaml:"time_to_crawl"` // mins

	Env string `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8081,
		},
		Redis: cc.DefaultRedis(),
		TelegramBot: cc.TelegramBot{
			Chats: map[string]int64{
				"default": 0,
			},
		},
		Env:         cmenv.EnvDev.String(),
		TimeLimit:   3,
		TimeToCrawl: 60,
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
