package drivers

import (
	"etop.vn/api/main/identity"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel"
)

func ITopX(env cm.EnvType) *whitelabel.WL {
	cfg := config{
		prodHost: "itopx.vn",
		key:      "itopx",
	}
	return &whitelabel.WL{
		Partner: identity.Partner{
			ID:         1057192413421863086,
			Name:       "IM Group",
			PublicName: "IM Group",
			ImageURL:   "",
			WebsiteURL: "",
		},
		Config: whitelabel.Config{
			Key:             cfg.key,
			Host:            cfg.host(env),
			RootURL:         cfg.siteUrl(env, ""),
			AuthURL:         cfg.siteUrl(env, "/welcome"),
			SiteName:        "iTopX",
			CompanyName:     "IM Group",
			CompanyFullName: "Công ty cổ phần đầu tư và phát triển IM",
			CSEmail:         "support@imgroup.vn",
		},
		Driver: &itopxDriver{},
	}
}

type itopxDriver struct {
}
