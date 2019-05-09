package config

import (
	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/integration/kiotviet"
)

var ll = l.New()

const ChannelWebhook = "webhook"

// Config ...
type Config struct {
	Postgres    cc.Postgres    `yaml:"postgres" env:"PG"`
	Redis       cc.Redis       `yaml:"redis"`
	HTTP        cc.HTTP        `yaml:"http" env:"HTTP"`
	Webhook     cc.HTTP        `yaml:"webhook"`
	TelegramBot cc.TelegramBot `yaml:"telegram_bot"`

	Sync kiotviet.SyncConfig `yaml:"sync"`

	Secret   string `yaml:"secret"`
	ServeDoc bool   `yaml:"serve_doc"`
	Env      string `yaml:"env"`
}

// Default ...
func Default() Config {
	cfg := Config{
		Redis:    cc.DefaultRedis(),
		Postgres: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 9001,
		},
		Webhook: cc.HTTP{
			Host: "",
			Port: 9002,
		},
		Sync: kiotviet.SyncConfig{
			WebhookBaseURL: "",
		},
		Secret:   "9RZEnRx3kBb7wm4z",
		ServeDoc: true,
		Env:      cm.EnvDev,
	}
	cfg.Postgres.Database = "etop_dev"
	return cfg
}

// Load loads config from file
func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cfg.Redis.MustLoadEnv()
	cfg.TelegramBot.MustLoadEnv()
	cc.MustLoadEnv("ET_SECRET", &cfg.Secret)
	return cfg, err
}
