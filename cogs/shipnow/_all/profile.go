package shipnow_all

import (
	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	"o.o/backend/pkg/integration/shipnow/ahamove/server"
)

func AhamoveConfig(cfg config.Config) ahamove.URLConfig {
	return ahamove.URLConfig{
		ThirdPartyHost:       cfg.ThirdPartyHost,
		PathUserVerification: server.PathAhamoveUserVerification,
	}
}
