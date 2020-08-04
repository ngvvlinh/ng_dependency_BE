package config

import (
	_telebot "o.o/backend/cogs/base/telebot"
	storage_all "o.o/backend/cogs/storage/_all"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/storage"
)

type Purpose = storage.Purpose
type DirConfigs = storage.DirConfigs
type DirConfig = storage.DirConfig

const (
	PurposeDefault             Purpose = "default"
	PurposeAhamoveVerification Purpose = "ahamove_verification"
)

func SupportedPurposes() []Purpose {
	return []Purpose{
		PurposeDefault,
		PurposeAhamoveVerification,
	}
}

func DefaultDirConfigs(httpAddress string) DirConfigs {
	dirs := DirConfigs{}
	dirs[PurposeDefault] =
		DirConfig{
			Path:      "img",
			URLPath:   "/img",
			URLPrefix: httpAddress + "/img",
		}
	dirs[PurposeAhamoveVerification] = DirConfig{
		Path:      "ahamove/user_verification",
		URLPath:   "/ahamove/user_verification",
		URLPrefix: httpAddress + "/ahamove/user_verification",
	}
	return dirs
}

type Config struct {
	Redis          cc.Redis `yaml:"redis"`
	HTTP           cc.HTTP  `yaml:"http"`
	cc.TelegramBot `yaml:"telegram_bot"`

	Dirs          DirConfigs               `yaml:"upload_dirs"`
	StorageDriver storage_all.DriverConfig `yaml:"storage_driver"`

	Env string `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		HTTP:        cc.HTTP{Port: 8180},
		Redis:       cc.DefaultRedis(),
		TelegramBot: _telebot.DefaultConfig(),
		Env:         cmenv.EnvDev.String(),
	}
	cfg.Dirs = DefaultDirConfigs("http://localhost:8180")
	cfg.StorageDriver = storage_all.DefaultDriver()
	return cfg
}

// Load loads config from file
func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	if err != nil {
		return cfg, err
	}

	cc.RedisMustLoadEnv(&cfg.Redis)
	cfg.TelegramBot.MustLoadEnv()
	return cfg, err
}
