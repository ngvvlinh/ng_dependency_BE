package config

import (
	_telebot "o.o/backend/cogs/base/telebot"
	telecom_all "o.o/backend/cogs/telecom/_all"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
)

type Config struct {
	Databases cc.Databases `yaml:",inline"`

	Redis        cc.Redis                       `yaml:"redis"`
	HTTP         cc.HTTP                        `yaml:"http"`
	TelegramBot  cc.TelegramBot                 `yaml:"telegram_bot"`
	Env          string                         `yaml:"env"`
	AdminPortsip telecom_all.AdminPortsipConfig `yaml:"admin_portsip"`
}

func Default() Config {
	cfg := Config{
		Databases: map[string]*cc.Postgres{
			"postgres":          cc.PtrDefaultPostgres(),
			"postgres_etelecom": cc.PtrDefaultPostgres(),
		},
		HTTP:        cc.HTTP{Port: 8200},
		Redis:       cc.DefaultRedis(),
		TelegramBot: _telebot.DefaultConfig(),
		Env:         cmenv.EnvDev.String(),
	}
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	if err != nil {
		return cfg, err
	}

	cc.PostgresMustLoadEnv(cfg.Databases["postgres"])
	cc.PostgresMustLoadEnv(cfg.Databases["postgres_etelecom"])
	cc.RedisMustLoadEnv(&cfg.Redis)
	cfg.TelegramBot.MustLoadEnv()

	return cfg, err
}
