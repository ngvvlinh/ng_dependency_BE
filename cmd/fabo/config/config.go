package config

import (
	"fmt"

	cc "etop.vn/backend/pkg/common/config"
)

type AppInfo struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_id"`
}

type ApiInfo struct {
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
}

func (api ApiInfo) Url() string {
	return fmt.Sprintf("%s/%s", api.Host, api.Version)
}

type Config struct {
	Postgres    cc.Postgres    `yaml:"postgres"`
	HTTP        cc.HTTP        `yaml:"http"`
	TelegramBot cc.TelegramBot `yaml:"telegram_bot"`
	Env         string         `yaml:"env"`
	App         AppInfo        `yaml:"app"`
	ApiInfo     ApiInfo        `yaml:"api_info"`
}

func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8080,
		},
		TelegramBot: cc.TelegramBot{},
		App: AppInfo{
			AppID:     "1581362285363031",
			AppSecret: "b3962ddf033b295c2bd0b543fff904f7",
		},
		ApiInfo: ApiInfo{
			Host:    "https://graph.facebook.com",
			Version: "v6.0",
		},
		Env: "dev",
	}
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cfg.TelegramBot.MustLoadEnv()
	return cfg, err
}
