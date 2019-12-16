package config

import (
	cc "etop.vn/backend/pkg/common/config"
)

type Config struct {
	Postgres     cc.Postgres `yaml:"postgres"`
	PostgresTest cc.Postgres `yaml:"postgres_target"`
	HTTP         cc.HTTP     `yaml:"http"`
	Env          string      `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		Postgres:     cc.DefaultPostgres(),
		PostgresTest: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8080,
		},
		Env: "dev",
	}
	cfg.Postgres.Database = "etop_crm"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cc.PostgresMustLoadEnv(&cfg.PostgresTest)
	return cfg, err
}
