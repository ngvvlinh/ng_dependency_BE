package config

import (
	cc "o.o/backend/pkg/common/config"
)

// Config ...
type Config struct {
	NewDatabase cc.Postgres `yaml:"new_database"`
	OldDatabase cc.Postgres `yaml:"old_database"`
}

// Default ...
func Default() Config {
	cfg := Config{
		NewDatabase: cc.DefaultPostgres(),
		OldDatabase: cc.DefaultPostgres(),
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
