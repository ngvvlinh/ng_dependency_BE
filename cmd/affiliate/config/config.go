package config

import (
	cc "etop.vn/backend/pkg/common/config"
)

type Config struct {
	Postgres    cc.Postgres `yaml:"postgres"`
	AffPostgres cc.Postgres `yaml:"postgres"`
	Redis       cc.Redis    `yaml:"redis"`
	HTTP        cc.HTTP     `yaml:"http"`
	Env         string      `yaml:"env"`
	Secret      string      `yaml:"secret"`

	SAdminToken string `yaml:"sadmin_token"`
}

func Default() Config {
	cfg := Config{
		Postgres:    cc.DefaultPostgres(),
		AffPostgres: cc.DefaultPostgres(),
		Redis:       cc.DefaultRedis(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8081,
		},
		Env:         "dev",
		Secret:      "secret",
		SAdminToken: "PZJvDAY2.sadmin.HXnnEkdV",
	}
	cfg.Postgres.Database = "etop_dev"
	cfg.AffPostgres.Database = "etop_aff"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	return cfg, err
}
