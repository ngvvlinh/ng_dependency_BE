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
	Kafka       cc.Kafka           `yaml:"kafka"`
	Env         string             `yaml:"env"`
	FacebookApp fbclient.AppConfig `yaml:"facebook_app"`
	Webhook     WebhookConfig      `yaml:"webhook"`
}

type WebhookConfig struct {
	HTTP        cc.HTTP `yaml:"http"`
	VerifyToken string  `yaml:"verify_token"`
}

func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8080,
		},
		Kafka: cc.DefaultKafka(),
		Redis: cc.DefaultRedis(),
		TelegramBot: cc.TelegramBot{
			Chats: map[string]int64{
				"default": 0,
				"webhook": 0,
				"import":  0,
				"sms":     0,
			},
		},
		FacebookApp: fbclient.AppConfig{
			ID:          "1581362285363031",
			Secret:      "b3962ddf033b295c2bd0b543fff904f7",
			AccessToken: "1581362285363031|eLuNU9-1KNA0AMNucV9PQIHCF1A",
		},
		Webhook: WebhookConfig{
			HTTP: cc.HTTP{
				Host: "",
				Port: 8081,
			},
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
