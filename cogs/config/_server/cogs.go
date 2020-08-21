package config_server

import (
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/authorize/tokens"
)

type SharedConfig struct {
	HTTP        cc.HTTP `yaml:"http"`
	ServeDoc    bool    `yaml:"serve_doc"`
	SAdminToken string  `yaml:"sadmin_token"`

	Env string `yaml:"env"`
}

func DefaultConfig() SharedConfig {
	return SharedConfig{
		HTTP:        cc.HTTP{Port: 8080},
		ServeDoc:    true,
		SAdminToken: "PZJvDAY2.sadmin.HXnnEkdV",
		Env:         cmenv.EnvDev.String(),
	}
}

func NewSession(cfg SharedConfig, redisStore redis.Store) session.Session {
	return session.New(
		session.OptValidator(tokens.NewTokenStore(redisStore)),
		session.OptSuperAdmin(cfg.SAdminToken),
	)
}

func WireSAdminToken(cfg SharedConfig) middleware.SAdminToken {
	return middleware.SAdminToken(cfg.SAdminToken)
}
