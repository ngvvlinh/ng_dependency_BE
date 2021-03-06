package templatemessages

import (
	"html/template"
)

var ResetPasswordTpl = template.Must(template.New("reset-password").Parse(`
Gửi {{.FullName}},<br><br>

Bạn (hoặc một ai đó) đang muốn khôi phục mật khẩu của tài khoản <b>{{.Email}}</b>.<br>
Nếu là bạn, hãy bấm vào liên kết bên dưới để khôi phục mật khẩu: (có hiệu lực trong 24 giờ)<br><br>

<a href="{{.URL}}">{{.URL}}</a><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ {{.WlName}}
`))

var UpdateEmailTpl = template.Must(template.New("update-email").Parse(`
Gửi {{.FullName}},<br><br>

Bạn (hoặc một ai đó) đang muốn thay đổi email của tài khoản <b>{{.Email}}</b>.<br>
Nếu là bạn, hãy sử dụng mã bên dưới để tiếp tục thực hiện thay đổi: (có hiệu lực trong 2 giờ)<br><br>

<div style="font-size:30px;margin-left:60px;padding:5px 20px;border:solid 2px #aaa;background-color:#eee;display:inline-block">{{.Code}}</div><br><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ {{.WlName}}
`))

var UpdateEmailTplConfirm = template.Must(template.New("update-email").Parse(`
Gửi {{.FullName}},<br><br>

Bạn (hoặc một ai đó) đang muốn xác nhận địa chỉ email của tài khoản <b>{{.Email}}</b>.<br>
Nếu là bạn, hãy sử dụng mã bên dưới để tiếp tục thực hiện thay đổi: (có hiệu lực trong 2 giờ)<br><br>

<div style="font-size:30px;margin-left:60px;padding:5px 20px;border:solid 2px #aaa;background-color:#eee;display:inline-block">{{.Code}}</div><br><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ {{.WlName}}
`))

var UpdatePhoneTpl = template.Must(template.New("update-phone").Parse(`
Gửi {{.FullName}},<br><br>

Bạn (hoặc một ai đó) đang muốn thay đổi thông tin số điện thoại của tài khoản <b>{{.Email}}</b>.<br>
Nếu là bạn, hãy sử dụng mã bên dưới để tiếp tục thực hiện thay đổi: (có hiệu lực trong 2 giờ)<br><br>

<div style="font-size:30px;margin-left:60px;padding:5px 20px;border:solid 2px #aaa;background-color:#eee;display:inline-block">{{.Code}}</div><br><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ {{.WlName}}
`))

var EmailVerificationTpl = template.Must(template.New("email-verification").Parse(`
Gửi {{.FullName}},<br><br>

Bạn (hoặc một ai đó) đang muốn xác nhận địa chỉ email của tài khoản <b>{{.Email}}</b>.<br>
Nếu là bạn, hãy bấm vào liên kết bên dưới để xác nhận địa chỉ email: (có hiệu lực trong 24 giờ)<br><br>

<a href="{{.URL}}">{{.URL}}</a><br><br>

Nếu không phải bạn, hãy bỏ qua email này. Bạn cũng có thể sử dụng chức năng khôi phục mật khẩu để lấy lại tài khoản.<br><br>

Đội ngũ {{.WlName}}
`))

var EmailInvitationTpl = template.Must(template.New("email-verification").Parse(`
Gửi <b>{{.FullName}}</b>,<br><br>

Bạn được <b>{{.InvitedUsername}}</b> mời tham gia cửa hàng <b>{{.ShopName}}</b> với vai trò <b>{{.ShopRoles}}</b>.<br>
Hãy bấm vào liên kết bên dưới để xác nhận lời mời: (có hiệu lực trong 24 giờ)<br><br>
<a href="{{.URL}}">{{.URL}}</a><br><br>

Nếu bạn không nhận ra cửa hàng trên, hãy bỏ qua email này.<br><br>

Đội ngũ {{.WlName}}
`))

var PhoneInvitationTpl = template.Must(template.New("phone-verification").Parse(`Bạn được mời vào cửa hàng {{.ShopName}} trên etop. Hãy bấm vào liên kết sau để xác nhận {{.URL}}`))

var PhoneInvitationTplRepeat = template.Must(template.New("phone-verification").Parse(`Bạn được mời vào cửa hàng {{.ShopName}} trên etop. Hãy bấm vào liên kết sau để xác nhận {{.URL}} (gửi lần {{.SendTime}})`))

var EmailSTokenTpl = template.Must(template.New("email-verification").Parse(`
Gửi {{.FullName}},<br><br>

Bạn (hoặc một ai đó) đang muốn thay đổi thông tin {{.AccountType}} <b>{{.AccountName}}</b> được quản lý bởi tài khoản <b>{{.Email}}</b>.<br>
Nếu là bạn, hãy sử dụng mã bên dưới để tiếp tục thực hiện thay đổi: (có hiệu lực trong 2 giờ)<br><br>

<div style="font-size:30px;margin-left:60px;padding:5px 20px;border:solid 2px #aaa;background-color:#eee;display:inline-block">{{.Code}}</div><br><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ {{.WlName}}
`))

var SmsVerificationTpl = `EB2B - Nhập mã %v để xác nhận thông tin tài khoản của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất ky ai.`

var SmsRegisterSimplifyTpl = `EB2B - Nhập mã %v để truy cập vào ứng dụng của bạn. Mã có hiệu lực trong 5 phút. Vui lòng không chia sẻ cho bất ky ai.`

var SmsResetPasswordTpl = `EB2B - Nhập mã %v để khôi phục mật khẩu tài khoản của bạn. Vui lòng không chia sẻ cho bất ky ai.`

var SmsChangeEmailTpl = `EB2B - Nhập mã %v để thay đổi thông tin email tài khoản %v của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất ky ai.`

var SmsChangeEmailTplRepeat = `EB2B - Nhập mã %v để thay đổi thông tin email tài khoản %v của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất ky ai. (gửi lần %v)`

var SmsChangePhoneTpl = `EB2B - Nhập mã %v để thay đổi số điện thoại tài khoản %v của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất ky ai.`

var SmsChangePhoneTplRepeat = `EB2B - Nhập mã %v để thay đổi số điện thoại tài khoản %v của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất ky ai. (gửi lần %v)`

var SmsChangePhoneTplConfirm = `EB2B - Nhập mã %v để xác nhận số điện thoại tài khoản %v của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất ky ai.`

var SmsChangePhoneTplConfirmRepeat = `EB2B - Nhập mã %v để xác nhận số điện thoại tài khoản %v của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất ky ai. (gửi lần %v)`

var SmsVerificationTplRepeat = `EB2B - Nhập mã %v để xác nhận thông tin tài khoản của bạn. Mã có hiệu lực trong 1 giờ. Vui lòng không chia sẻ cho bất ky ai. (gửi lần %v)`

var SmsResetPasswordTplRepeat = `EB2B - Nhập mã %v để khôi phục mật khẩu tài khoản của bạn. Vui lòng không chia sẻ cho bất ky ai. (gửi lần %v)`

var RequestLoginEmailTpl = template.Must(template.New("request-login-email").Parse(`
{{.Hello}},<br><br>

Bạn (hoặc một ai đó) đang muốn đăng nhập vào tài khoản <b>{{.Email}}</b> thông qua hệ thống <b>{{.PartnerPublicName}}</b> ({{.PartnerWebsite}}). Nếu là bạn, hãy sử dụng mã bên dưới để đăng nhập: (có hiệu lực trong 2 giờ)<br><br>

<div style="font-size:30px;margin-left:60px;padding:5px 20px;border:solid 2px #aaa;background-color:#eee;display:inline-block">{{.Code}}</div><br><br><br>

<b>{{.Notice}}</b><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ {{.WlName}}
`))

var EmailVerificationByOTPTpl = template.Must(template.New("verify-email-by-otp").Parse(`
Bạn (hoặc một ai đó) đang muốn xác nhận địa chỉ email <b>{{.Email}}</b>. Nếu là bạn, hãy sử dụng mã bên dưới để tiếp tục xác thực: (có hiệu lực trong 2 giờ)<br><br>

<div style="font-size:30px;margin-left:60px;padding:5px 20px;border:solid 2px #aaa;background-color:#eee;display:inline-block">{{.Code}}</div><br><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ {{.WlName}}
{{.Extra}}
`))

var NewAccountViaPartnerEmailTpl = template.Must(template.New("register-email").Parse(`
Gửi {{.FullName}},<br><br>

Chào mừng bạn đã tạo tài khoản eTop.vn mới thông qua hệ thống <b>{{.PartnerPublicName}}</b> ({{.PartnerWebsite}}). Bạn có thể sử dụng thông tin bên dưới để đăng nhập vào eTop.vn và đổi mật khẩu:<br><br>

<div style="margin: 0 0 0 30px">
  <div style="float:left; width:110px">Đăng nhập: </div>
  <div style="display:inline-block; width:50%"><a href="https://etop.vn/login">https://etop.vn/login</a></div>
</div>
<div style="margin: 0 0 0 30px">
  <div style="float:left; width:110px">{{.LoginLabel}}: </div>
  <div style="display:inline-block; width:50%"><b>{{.Login}}</b></div>
</div>
<div style="margin: 0 0 0 30px">
  <div style="float:left; width:110px">Mật khẩu: </div>
  <div style="display:inline-block; width:50%"><b>{{.Password}}</b></div>
</div><br><br>

<b>Vui lòng chỉ sử dụng mật khẩu này ở https://etop.vn và không chia sẻ cho người khác. Nhân viên và đối tác của eTop sẽ không bao giờ hỏi mật khẩu của bạn.</b><br><br>

Đội ngũ eTop
`))
