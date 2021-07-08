package config

import (
	cc "o.o/backend/pkg/common/config"
)

// Config ...
type Config struct {
	Postgres cc.Postgres `yaml:"postgres"`
}

// Default ...
func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
	}
	return cfg
}

// Load loads config from file
func Load() (Config, error) {
	var cfg Config
	err := cc.LoadWithDefault(&cfg, Default())
	if err != nil {
		return cfg, err
	}

	return cfg, err
}
