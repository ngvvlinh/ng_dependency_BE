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
			InviteUserURLByEmail: ternary(env == cmenv.EnvProd,
				"https://etop.vn/invitation",
				"https://shop."+baseHost(env)+"/invitation",
			),
			InviteUserURLByPhone: ternary(env == cmenv.EnvProd,
				"https://etop.vn/i",
				"https://shop."+baseHost(env)+"/i",
			),
			SiteName:        "eTop",
			CompanyName:     "eTop",
			CompanyFullName: "Công ty cổ phần công nghệ eTop",
			CSEmail:         "hotro@etop.vn",
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
