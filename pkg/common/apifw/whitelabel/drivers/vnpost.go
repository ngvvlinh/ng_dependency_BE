package drivers

import (
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/wl_type"
	"o.o/backend/pkg/common/apifw/whitelabel"
	"o.o/backend/pkg/common/cmenv"
	"o.o/capi/dot"
)

const VNPostID = dot.ID(1156518020386448488)
const VNPostKey = "vnpost"

func VNPost(env cmenv.EnvType) *whitelabel.WL {
	cfg := config{
		prodHost: "vnpost",
		key:      VNPostKey,
	}
	return &whitelabel.WL{
		Partner: identity.Partner{
			ID:         VNPostID,
			Name:       "VNPost",
			PublicName: "VNPost",
			ImageURL:   "",
			WebsiteURL: "",
		},
		Config: whitelabel.Config{
			Key:                  cfg.key,
			Host:                 cfg.host(env),
			RootURL:              cfg.siteUrl(env, ""),
			AuthURL:              cfg.siteUrl(env, "/welcome/login"),
			InviteUserURLByEmail: cfg.siteUrl(env, "/invitation"),
			InviteUserURLByPhone: cfg.siteUrl(env, "/i"),
			SiteName:             "vnpost",
			CompanyName:          "VNPost",
			CompanyFullName:      "VNPost",
			CSEmail:              "",
			WLType:               wl_type.POS,

			// vnpost hiện tại không mở ra như một whitelabel partner xài etop/topship.
			IgnoreParseFromHost: true,
		},
		Driver: &vnpostDriver{},
	}
}

type vnpostDriver struct {
}
