package config

import (
	"errors"
	"strings"

	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	Redis          cc.Redis `yaml:"redis"`
	HTTP           cc.HTTP  `yaml:"http"`
	cc.TelegramBot `yaml:"telegram_bot"`

	UploadDirImg                 string `yaml:"upload_dir_img"`
	URLPrefix                    string `yaml:"url_prefix"`
	UploadDirAhamoveVerification string `yaml:"upload_dir_ahamove_verification"`
	URLPrefixAhamoveVerification string `yaml:"url_prefix_ahamove_verification"`

	Env string `yaml:"env"`
}

func Default() Config {
	cfg := Config{
		HTTP:         cc.HTTP{Port: 8180},
		Redis:        cc.DefaultRedis(),
		UploadDirImg: "/tmp/upload",
		TelegramBot: cc.TelegramBot{
			Chats: map[string]int64{
				"default": 0,
			},
		},
		URLPrefix:                    "http://localhost:8180/img",
		UploadDirAhamoveVerification: "/tmp/upload",
		URLPrefixAhamoveVerification: "http://localhost:8180",
		Env:                          cmenv.EnvDev.String(),
	}
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

	if cfg.UploadDirImg == "" {
		return cfg, errors.New("empty upload_dir")
	}
	if cfg.URLPrefix == "" {
		return cfg, errors.New("empty url_prefix")
	}
	cfg.URLPrefix = strings.TrimSuffix(cfg.URLPrefix, "/")

	if cfg.UploadDirAhamoveVerification == "" {
		return cfg, errors.New("Empty upload_dir_ahamove_verification")
	}
	if cfg.URLPrefixAhamoveVerification == "" {
		return cfg, errors.New("Empty url_prefix_ahamove_verification")
	}
	cfg.URLPrefixAhamoveVerification = strings.TrimSuffix(cfg.URLPrefixAhamoveVerification, "/")

	return cfg, err
}
