package config

import (
	cc "o.o/backend/pkg/common/config"
)

// Config ...
type Config struct {
	PostgresEtelecom cc.Postgres `yaml:"postgres_etelecom"`
}

// Default ...
func Default() Config {
	cfg := Config{
		PostgresEtelecom: cc.DefaultPostgres(),
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
