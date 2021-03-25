package drivers

import (
	"o.o/api/main/identity"
	"o.o/backend/pkg/common/apifw/whitelabel"
	"o.o/backend/pkg/common/cmenv"
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
			Key:                  "etop",
			Host:                 "shop." + baseHost(env),
			RootURL:              "https://shop." + baseHost(env) + "/register",
			AuthURL:              "https://auth." + baseHost(env),
			InviteUserURLByEmail: "https://shop." + baseHost(env) + "/invitation",
			InviteUserURLByPhone: "https://shop." + baseHost(env) + "/i",
			SiteName:             "eTop",
			CompanyName:          "eTop",
			CompanyFullName:      "Công ty cổ phần công nghệ eTop",
			CSEmail:              "hotro@etop.vn",
			Templates: &whitelabel.Templates{
				RequestLoginSmsTpl: whitelabel.MustParseTemplate("request-login-sms",
					"Nhập mã {{.Code}} để đăng nhập vào tài khoản eTop thông qua hệ thống của đối tác. Mã có hiệu lực trong 2 giờ. Vui lòng không chia sẻ cho bất kỳ ai."),
				NewAccountViaPartnerSmsTpl: whitelabel.MustParseTemplate("register-sms",
					`Sử dụng mật khẩu {{.Password}} để đăng nhập vào tài khoản eTop của bạn. Vui lòng chỉ sử dụng mật khẩu này ở etop.vn và không chia sẻ cho bất kỳ ai.`),
			},
		},
		Driver: &etopDriver{},
	}
}

type etopDriver struct {
}
