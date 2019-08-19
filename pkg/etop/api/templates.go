package api

import "html/template"

var resetPasswordTpl = template.Must(template.New("reset-password").Parse(`
Gửi {{.FullName}},<br><br>

Bạn (hoặc một ai đó) đang muốn khôi phục mật khẩu của tài khoản <b>{{.Email}}</b>.<br>
Nếu là bạn, hãy bấm vào đường link bên dưới để khôi phục mật khẩu: (có hiệu lực trong 24 giờ)<br><br>

<a href="{{.URL}}">{{.URL}}</a><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ eTop
`))

var emailVerificationTpl = template.Must(template.New("email-verification").Parse(`
Gửi {{.FullName}},<br><br>

Bạn (hoặc một ai đó) đang muốn xác nhận địa chỉ email của tài khoản <b>{{.Email}}</b>.<br>
Nếu là bạn, hãy bấm vào đường link bên dưới để xác nhận địa chỉ email: (có hiệu lực trong 24 giờ)<br><br>

<a href="{{.URL}}">{{.URL}}</a><br><br>

Nếu không phải bạn, hãy bỏ qua email này. Bạn cũng có thể sử dụng chức năng khôi phục mật khẩu để lấy lại tài khoản.<br><br>

Đội ngũ eTop
`))

var emailSTokenTpl = template.Must(template.New("email-verification").Parse(`
Gửi {{.FullName}},<br><br>

Bạn (hoặc một ai đó) đang muốn thay đổi thông tin {{.AccountType}} <b>{{.AccountName}}</b> được quản lý bởi tài khoản <b>{{.Email}}</b>.<br>
Nếu là bạn, hãy sử dụng mã bên dưới để tiếp tục thực hiện thay đổi: (có hiệu lực trong 2 giờ)<br><br>

<div style="font-size:30px;margin-left:60px;padding:5px 20px;border:solid 2px #aaa;background-color:#eee;display:inline-block">{{.Code}}</div><br><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ eTop
`))

var smsVerificationTpl = `Nhập mã %v để xác nhận thông tin tài khoản %v của bạn trên eTop.vn. Mã có hiệu lực trong 2 giờ.`

var RequestLoginEmailTpl = template.Must(template.New("request-login-email").Parse(`
{{.Hello}},<br><br>

Bạn (hoặc một ai đó) đang muốn đăng nhập vào tài khoản <b>{{.Email}}</b> thông qua hệ thống <b>{{.PartnerPublicName}}</b> ({{.PartnerWebsite}}). Nếu là bạn, hãy sử dụng mã bên dưới để đăng nhập: (có hiệu lực trong 2 giờ)<br><br>

<div style="font-size:30px;margin-left:60px;padding:5px 20px;border:solid 2px #aaa;background-color:#eee;display:inline-block">{{.Code}}</div><br><br><br>

<b>{{.Notice}}</b><br><br>

Nếu không phải bạn, hãy bỏ qua email này.<br><br>

Đội ngũ eTop
{{.Extra}}
`))

var RequestLoginSmsTpl = template.Must(template.New("request-login-sms").Parse(`
Nhập mã {{.Code}} để đăng nhập vào tài khoản của bạn trên eTop.vn thông qua hệ thống của đối tác. Mã có hiệu lực trong 2 giờ. {{.Notice}}
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

var NewAccountViaPartnerSmsTpl = template.Must(template.New("register-sms").Parse(`
Sử dụng mật khẩu {{.Password}} để đăng nhập vào tài khoản của bạn trên eTop.vn. Vui lòng chỉ sử dụng mật khẩu này ở https://etop.vn và không chia sẻ cho người khác.
`))
