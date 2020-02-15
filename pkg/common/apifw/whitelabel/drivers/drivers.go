package drivers

import (
	"etop.vn/backend/pkg/common/apifw/whitelabel"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/common/l"
)

func Drivers(env cmenv.EnvType) []*whitelabel.WL {
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

func baseHost(env cmenv.EnvType) string {
	switch env {
	case cmenv.EnvDev:
		return "d.etop.vn"

	case cmenv.EnvSandbox:
		return "sandbox.etop.vn"

	case cmenv.EnvStag:
		return "g.etop.vn"

	default:
		ll.S.Panicf("unexpected env: %v", env)
		return ""
	}
}

func (c config) host(env cmenv.EnvType) string {
	if env == cmenv.EnvProd {
		return c.prodHost
	}
	return c.key + "." + baseHost(env)
}

func (c config) siteUrl(env cmenv.EnvType, path string) string {
	return "https://" + c.host(env) + path
}
