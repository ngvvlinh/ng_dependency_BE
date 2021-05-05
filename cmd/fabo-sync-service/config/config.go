package config

import (
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/sync"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	Databases   cc.Databases       `yaml:",inline"`
	HTTP        cc.HTTP            `yaml:"http"`
	TelegramBot cc.TelegramBot     `yaml:"telegram_bot"`
	Redis       cc.Redis           `yaml:"redis"`
	FacebookApp fbclient.AppConfig `yaml:"facebook_app"`
	SyncConfig  sync.Config        `yaml:",inline"`
	Env         string             `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		Databases: map[string]*cc.Postgres{
			"postgres": cc.PtrDefaultPostgres(),
		},
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
		Env: cmenv.EnvDev.String(),
		SyncConfig: sync.Config{
			TimeLimit: 3,
			TimeToRun: "0 0 * * *", // 00:00 everyday
		},
	}
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(cfg.Databases["postgres"])
	cfg.TelegramBot.MustLoadEnv()
	cfg.FacebookApp.MustLoadEnv()
	return cfg, err
}
