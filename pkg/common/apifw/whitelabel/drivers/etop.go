package drivers

import (
	"etop.vn/api/main/identity"
	"etop.vn/backend/pkg/common/apifw/whitelabel"
	"etop.vn/backend/pkg/common/cmenv"
)

func ETop(env cmenv.EnvType) *whitelabel.WL {
	return &whitelabel.WL{
		Partner: identity.Partner{
			ID:         0,
			Name:       "eTop",
			PublicName: "eTop",
			ImageURL:   "",
			WebsiteURL: "",
		},
		Config: whitelabel.Config{
			Key: "etop",
			Host: ternary(env == cmenv.EnvProd,
				"etop.vn",
				"shop."+baseHost(env),
			),
			RootURL: ternary(env == cmenv.EnvProd,
				"https://etop.vn/register",
				"https://shop."+baseHost(env),
			),
			AuthURL: ternary(env == cmenv.EnvProd,
				"https://auth.etop.vn",
				"https://auth."+baseHost(env),
			),
			SiteName:        "eTop",
			CompanyName:     "eTop",
			CompanyFullName: "Công ty cổ phần công nghệ eTop",
			CSEmail:         "hotro@etop.vn",
		},
		Driver: &etopDriver{},
	}
}

type etopDriver struct {
}
