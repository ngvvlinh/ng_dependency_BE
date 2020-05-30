package config

import (
	cc "o.o/backend/pkg/common/config"
)

type Config struct {
	Postgres        cc.Postgres    `yaml:"postgres"`
	PostgresWebhook cc.Postgres    `yaml:"postgres_webhook"`
	Redis           cc.Redis       `yaml:"redis"`
	HTTP            cc.HTTP        `yaml:"http"`
	Kafka           cc.Kafka       `yaml:"kafka"`
	TelegramBot     cc.TelegramBot `yaml:"telegram_bot"`
	Secret          string         `yaml:"secret"`
	Env             string         `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		Postgres:        cc.DefaultPostgres(),
		PostgresWebhook: cc.DefaultPostgres(),
		Redis:           cc.DefaultRedis(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8081,
		},
		Kafka:  cc.DefaultKafka(),
		Secret: "secret",
		Env:    "dev",
	}
	cfg.Postgres.Database = "etop_dev"
	cfg.PostgresWebhook.Database = "etop_webhook"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cc.PostgresMustLoadEnv(&cfg.PostgresWebhook, "ET_POSTGRES_WEBHOOK")
	cc.RedisMustLoadEnv(&cfg.Redis)
	cfg.TelegramBot.MustLoadEnv()
	cc.EnvMap{
		"ET_SECRET": &cfg.Secret,
	}.MustLoad()
	return cfg, err
}
