package drivers

import (
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel"
	"etop.vn/common/l"
)

func Drivers(env cm.EnvType) []*whitelabel.WL {
	return []*whitelabel.WL{
		ETop(env),
		ITopX(env),
	}
}

var ll = l.New()

type config struct {
	prodHost string
	key      string
}

func ternary(cond bool, whenTrue string, whenFalse string) string {
	if cond {
		return whenTrue
	}
	return whenFalse
}

func baseHost(env cm.EnvType) string {
	switch env {
	case cm.EnvDev:
		return "d.etop.vn"

	case cm.EnvSandbox:
		return "s.etop.vn"

	case cm.EnvStag:
		return "g.etop.vn"

	default:
		ll.S.Panicf("unexpected env: %v", env)
		return ""
	}
}

func (c config) host(env cm.EnvType) string {
	if env == cm.EnvProd {
		return c.prodHost
	}
	return c.key + "." + baseHost(env)
}

func (c config) siteUrl(env cm.EnvType, path string) string {
	return "https://" + c.host(env) + path
}
