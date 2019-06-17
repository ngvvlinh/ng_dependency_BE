package config

import (
	"errors"
	"strings"

	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/l"
)

var ll = l.New()

const PathAhamoveUserVerification = "/upload/ahamove/user_verification"

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
		HTTP:                         cc.HTTP{Port: 8180},
		Redis:                        cc.DefaultRedis(),
		UploadDirImg:                 "/tmp/upload",
		URLPrefix:                    "http://localhost:8180/img",
		UploadDirAhamoveVerification: "/tmp/upload",
		URLPrefixAhamoveVerification: "http://localhost:8180",
		Env:                          cm.EnvDev,
	}
	return cfg
}

// Load loads config from file
func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	if err != nil {
		return cfg, err
	}

	cfg.Redis.MustLoadEnv()
	cfg.TelegramBot.MustLoadEnv()

	if cfg.UploadDirImg == "" {
		return cfg, errors.New("empty upload_dir")
	}
	if cfg.URLPrefix == "" {
		return cfg, errors.New("empty url_prefix")
	}
	cfg.URLPrefix = strings.TrimSuffix(cfg.URLPrefix, "/")
	return cfg, err
}
