package config

import (
	cc "o.o/backend/pkg/common/config"
)

// Config ...
type Config struct {
	Database cc.Postgres `yaml:"database"`
}

// Default ...
func Default() Config {
	cfg := Config{
		Database: cc.DefaultPostgres(),
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
