package config

import (
	"strings"

	cc "etop.vn/backend/pkg/common/config"
)

type Database struct {
	cc.Postgres `yaml:"postgres"`
	Tables      []string `yaml:"tables"`
}

type Config struct {
	MapDB       map[string]Database `yaml:"databases"`
	HTTP        cc.HTTP             `yaml:"http"`
	TelegramBot cc.TelegramBot      `yaml:"telegram_bot"`
	Env         string              `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		// key = key of whitelabel
		MapDB: map[string]Database{
			"etop": {
				Postgres: cc.DefaultPostgres(),
			},
			"itopx": {
				Postgres: cc.Postgres{
					Protocol:       "",
					Host:           "postgres",
					Port:           5432,
					Username:       "postgres",
					Password:       "postgres",
					Database:       "test_dst",
					SSLMode:        "",
					Timeout:        15,
					GoogleAuthFile: "",
				},
				Tables: []string{
					"user",
					"account",
					"customer",
					"order",
				},
			},
		},
		HTTP: cc.HTTP{
			Host: "",
			Port: 8081,
		},
		Env: "dev",
	}
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	for key, cfgDB := range cfg.MapDB {
		cc.PostgresMustLoadEnv(&cfgDB.Postgres, "ET_WL_POSTGRES_"+strings.ToUpper(key))
		// it's not a pointer, so remember to assign back
		cfg.MapDB[key] = cfgDB
	}
	cfg.TelegramBot.MustLoadEnv()
	return cfg, err
}
