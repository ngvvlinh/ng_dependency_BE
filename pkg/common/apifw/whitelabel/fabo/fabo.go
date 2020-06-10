package fabo

import (
	"o.o/api/main/identity"
	"o.o/backend/pkg/common/apifw/whitelabel"
	"o.o/backend/pkg/common/cmenv"
	"o.o/common/l"
	_ "o.o/common/l"
)

func Drivers(env cmenv.EnvType) []*whitelabel.WL {
	return []*whitelabel.WL{
		Fabo(env),
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

// MUSTDO this is temporary config. Define later
func baseHost(env cmenv.EnvType) string {
	switch env {
	case cmenv.EnvDev:
		return "d.etop.vn"

	case cmenv.EnvSandbox:
		return "sandbox.etop.vn"

	case cmenv.EnvStag:
		return "g.etop.vn"

	case cmenv.EnvProd:
		return "faboshop.vn"

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

func Fabo(env cmenv.EnvType) *whitelabel.WL {
	InitTemplateMsg()
	cfg := config{
		prodHost: "faboshop.vn",
		key:      "fabo",
	}
	return &whitelabel.WL{
		Partner: identity.Partner{
			ID:         0,
			Name:       "FaboShop",
			PublicName: "FaboShop",
			ImageURL:   "",
			WebsiteURL: "",
		},
		Config: whitelabel.Config{
			Key:     cfg.key,
			Host:    cfg.host(env),
			RootURL: cfg.siteUrl(env, ""),
			AuthURL: ternary(env == cmenv.EnvProd,
				"https://auth."+cfg.prodHost,
				"https://auth."+baseHost(env),
			),
			InviteUserURLByEmail: cfg.siteUrl(env, "/invitation"),
			InviteUserURLByPhone: cfg.siteUrl(env, "/i"),
			SiteName:             "Faboshop",
			CompanyName:          "Faboshop",
			CompanyFullName:      "Công ty cổ phần công nghệ Faboshop",
			CSEmail:              "hotro@fabo.vn",
			Templates:            &whitelabel.Templates{},
		},
		Driver: &FaboDriver{},
	}
}

type FaboDriver struct {
}
