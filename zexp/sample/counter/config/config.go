package config

import (
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
)

type Config struct {
	HTTP     cc.HTTP     `yaml:"http"`
	Env      string      `yaml:"env"`
	Postgres cc.Postgres `yaml:"postgres"`
}

func Default() Config {
	cfg := Config{
		HTTP: cc.HTTP{
			Host: "127.0.0.1",
			Port: 8080,
		},
		Env:      cmenv.EnvDev.String(),
		Postgres: cc.DefaultPostgres(),
	}
	return cfg
}

func Load() (cfg Config, _ error) {
	err := cc.LoadWithDefault(&cfg, Default())
	return cfg, err
}
