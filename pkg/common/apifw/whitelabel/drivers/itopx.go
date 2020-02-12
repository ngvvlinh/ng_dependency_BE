package drivers

import (
	"etop.vn/api/main/identity"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/backend/pkg/common/apifw/whitelabel"
	"etop.vn/backend/pkg/common/cmenv"
)

const ITopXID = 1057192413421863086

func ITopX(env cmenv.EnvType) *whitelabel.WL {
	cfg := config{
		prodHost: "itopx.vn",
		key:      "itopx",
	}
	return &whitelabel.WL{
		Partner: identity.Partner{
			ID:         ITopXID,
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
			Shipment: &whitelabel.ShipmentConfig{
				Topship: []connection_type.ConnectionProvider{
					connection_type.ConnectionProviderGHN,
				},
			},
		},
		Driver: &itopxDriver{},
	}
}

type itopxDriver struct {
}
