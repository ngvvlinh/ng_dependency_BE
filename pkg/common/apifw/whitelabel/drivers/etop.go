package drivers

import (
	"etop.vn/api/main/identity"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel"
	etopmodel "etop.vn/backend/pkg/etop/model"
)

func ETop(env cm.EnvType) *whitelabel.WL {
	return &whitelabel.WL{
		Partner: identity.Partner{
			ID:         etopmodel.TagEtop,
			Name:       "eTop",
			PublicName: "eTop",
			ImageURL:   "",
			WebsiteURL: "",
		},
		Config: whitelabel.Config{
			Key: "etop",
			Host: ternary(env == cm.EnvProd,
				"etop.vn",
				"shop."+baseHost(env),
			),
			RootURL: ternary(env == cm.EnvProd,
				"https://etop.vn/register",
				"https://shop."+baseHost(env),
			),
			AuthURL: ternary(env == cm.EnvProd,
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
