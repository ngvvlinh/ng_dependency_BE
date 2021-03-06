package config

import (
	cc "o.o/backend/pkg/common/config"
)

type Config struct {
	Postgres         cc.Postgres        `yaml:"postgres"`
	PostgresNotifier cc.Postgres        `yaml:"postgres_notifier"`
	Redis            cc.Redis           `yaml:"redis"`
	HTTP             cc.HTTP            `yaml:"http"`
	Kafka            cc.Kafka           `yaml:"kafka"`
	TelegramBot      cc.TelegramBot     `yaml:"telegram_bot"`
	Onesignal        cc.OnesignalConfig `yaml:"onesignal"`
	Env              string             `yaml:"env"`
	URL              struct {
		MainSite string `yaml:"main_site"`
	} `yaml:"url"`
}

func Default() Config {
	cfg := Config{
		Postgres:         cc.DefaultPostgres(),
		PostgresNotifier: cc.DefaultPostgres(),
		Redis:            cc.DefaultRedis(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8083,
		},
		TelegramBot: cc.TelegramBot{
			Chats: map[string]int64{
				"default": 0,
			},
		},
		Kafka:     cc.DefaultKafka(),
		Env:       "dev",
		Onesignal: cc.DefaultOnesignal(),
	}
	cfg.PostgresNotifier.Database = "etop_notifier"
	cfg.URL.MainSite = "http://localhost:8080"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cc.PostgresMustLoadEnv(&cfg.PostgresNotifier, "ET_POSTGRES_NOTIFIER")
	cc.RedisMustLoadEnv(&cfg.Redis)
	cfg.TelegramBot.MustLoadEnv()
	cfg.Onesignal.MustLoadEnv("FABO_ONESIGNAL")
	return cfg, err
}
