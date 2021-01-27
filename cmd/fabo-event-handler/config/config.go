package config

import (
	"o.o/backend/cmd/fabo-server/config"
	database_min "o.o/backend/cogs/database/_min"
	"o.o/backend/com/fabo/pkg/fbclient"
	cc "o.o/backend/pkg/common/config"
)

type Config struct {
	Databases    database_min.Config        `yaml:",inline"`
	Redis       cc.Redis           `yaml:"redis"`
	HTTP        cc.HTTP            `yaml:"http"`
	Kafka       cc.Kafka           `yaml:"kafka"`
	TelegramBot cc.TelegramBot     `yaml:"telegram_bot"`
	Onesignal   cc.OnesignalConfig `yaml:"onesignal"`
	Secret      string             `yaml:"secret"`
	Env         string             `yaml:"env"`
	URL         struct {
		MainSite string `yaml:"main_site"`
	} `yaml:"url"`
	FacebookApp fbclient.AppConfig `yaml:"facebook_app"`
	Webhook     config.WebhookConfig      `yaml:"webhook"`
}

func Default() Config {
	cfg := Config{
		// TODO(vu): automatically map default config
		Databases:     database_min.DefaultConfig(),
		Redis: cc.DefaultRedis(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8081,
		},
		TelegramBot: cc.TelegramBot{
			Chats: map[string]int64{
				"default": 0,
			},
		},
		Kafka:  cc.DefaultKafka(),
		Secret: "secret",
		Env:    "dev",
	}
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cfg.Databases.MustLoadEnv()
	cc.RedisMustLoadEnv(&cfg.Redis)
	cfg.TelegramBot.MustLoadEnv()
	cfg.Onesignal.MustLoadEnv()
	cc.EnvMap{
		"ET_SECRET": &cfg.Secret,
	}.MustLoad()
	cfg.FacebookApp.MustLoadEnv()
	return cfg, err
}
