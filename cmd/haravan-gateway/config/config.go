package config

import (
	"errors"

	cc "etop.vn/backend/pkg/common/config"
	haravanclient "etop.vn/backend/pkg/integration/haravan/client"
	"etop.vn/backend/pkg/integration/shipping/ghn"
	"etop.vn/backend/pkg/integration/shipping/ghtk"
	"etop.vn/backend/pkg/integration/shipping/vtpost"
)

type Config struct {
	Postgres    cc.Postgres          `yaml:"postgres"`
	Redis       cc.Redis             `yaml:"redis"`
	HTTP        cc.HTTP              `yaml:"http"`
	TelegramBot cc.TelegramBot       `yaml:"telegram_bot"`
	Haravan     haravanclient.Config `yaml:"haravan"`
	GHN         ghn.Config           `yaml:"ghn"`
	GHTK        ghtk.Config          `yaml:"ghtk"`
	VTPost      vtpost.Config        `yaml:"vtpost"`
	URL         struct {
		MainSite string `yaml:"main_site"`
	} `yaml:"url"`
	Env string `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		Redis:    cc.DefaultRedis(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8085,
		},
		GHN:    ghn.DefaultConfig(),
		GHTK:   ghtk.DefaultConfig(),
		VTPost: vtpost.DefaultConfig(),
		Env:    "dev",
	}
	cfg.URL.MainSite = "http://localhost:8080"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cfg.Redis.MustLoadEnv()
	cfg.TelegramBot.MustLoadEnv()
	cfg.GHN.MustLoadEnv()
	cfg.GHTK.MustLoadEnv()
	cfg.VTPost.MustLoadEnv()
	cfg.Haravan.MustLoadEnv()

	if cfg.Haravan.Secret == "" {
		return cfg, errors.New("Empty Haravan secret")
	}
	if cfg.URL.MainSite == "" {
		return cfg, errors.New("Empty MainSite URL")
	}
	return cfg, err
}
