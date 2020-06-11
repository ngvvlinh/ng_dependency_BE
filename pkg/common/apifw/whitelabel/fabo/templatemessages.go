package fabo

import (
	"o.o/backend/pkg/common/apifw/whitelabel/templatemessages"
)

func InitTemplateMsg() {

	// Override Messages
	templatemessages.SmsVerificationTpl = `Nhập mã %v để xác nhận thông tin tài khoản Faboshop của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất kỳ ai.`

	templatemessages.SmsResetPasswordTpl = `Nhập mã %v để khôi phục mật khẩu tài khoản Faboshop của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất kỳ ai.`

	templatemessages.SmsVerificationTplRepeat = `Nhập mã %v để xác nhận thông tin tài khoản Faboshop của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất kỳ ai. (gửi lần %v)`

	templatemessages.SmsResetPasswordTplRepeat = `Nhập mã %v để khôi phục mật khẩu tài khoản Faboshop của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất kỳ ai. (gửi lần %v)`
}