package config

import (
	cc "o.o/backend/pkg/common/config"
)

type Config struct {
	Databases   cc.Databases   `yaml:",inline"`
	Redis       cc.Redis       `yaml:"redis"`
	HTTP        cc.HTTP        `yaml:"http"`
	Kafka       cc.Kafka       `yaml:"kafka"`
	TelegramBot cc.TelegramBot `yaml:"telegram_bot"`
	Secret      string         `yaml:"secret"`
	Env         string         `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		// TODO(vu): automatically map default config
		Databases: map[string]*cc.Postgres{
			"postgres":         cc.PtrDefaultPostgres(),
			"postgres_webhook": cc.PtrDefaultPostgres(),
		},
		Redis: cc.DefaultRedis(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8081,
		},
		TelegramBot: cc.TelegramBot{
			// TODO(vu): automatically map default config
			Chats: map[string]int64{
				"default": 0,
				"deploy":  0,
				"high":    0,
			},
		},
		Kafka:  cc.DefaultKafka(),
		Secret: "secret",
		Env:    "dev",
	}
	cfg.Databases["postgres"].Database = "etop_dev"
	cfg.Databases["postgres_webhook"].Database = "etop_dev"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(cfg.Databases["postgres"])
	cc.PostgresMustLoadEnv(cfg.Databases["postgres_webhook"], "ET_POSTGRES_WEBHOOK")
	cc.RedisMustLoadEnv(&cfg.Redis)
	cfg.TelegramBot.MustLoadEnv()
	cc.EnvMap{
		"ET_SECRET": &cfg.Secret,
	}.MustLoad()
	return cfg, err
}
