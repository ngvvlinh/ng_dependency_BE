package config

import (
	"o.o/backend/com/fabo/pkg/fbclient"
	cc "o.o/backend/pkg/common/config"
)

type Config struct {
	Postgres    cc.Postgres        `yaml:"postgres"`
	HTTP        cc.HTTP            `yaml:"http"`
	TelegramBot cc.TelegramBot     `yaml:"telegram_bot"`
	Redis       cc.Redis           `yaml:"redis"`
	Env         string             `yaml:"env"`
	FacebookApp fbclient.AppConfig `yaml:"facebook_app"`
}

func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8080,
		},
		TelegramBot: cc.TelegramBot{},
		Redis:       cc.DefaultRedis(),
		FacebookApp: fbclient.AppConfig{
			ID:          "1581362285363031",
			Secret:      "b3962ddf033b295c2bd0b543fff904f7",
			AccessToken: "1581362285363031|eLuNU9-1KNA0AMNucV9PQIHCF1A",
		},
		Env: "dev",
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
