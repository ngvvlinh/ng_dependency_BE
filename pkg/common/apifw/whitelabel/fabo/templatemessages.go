package fabo

import (
	"html/template"

	"o.o/backend/pkg/common/apifw/whitelabel/templatemessages"
)

func InitTemplateMsg() {
	// TODO(Nam) Đợi A Vũ có deal template
	// example override messages
	templatemessages.UpdateEmailTpl = template.Must(template.New("update-email").Parse(`
	Gửi {{.FullName}},<br><br>
	
	Bạn (hoặc một ai đó) đang muốn thay đổi email của tài khoản <b>{{.Email}}</b>.<br>
	Nếu là bạn, hãy sử dụng mã bên dưới để tiếp tục thực hiện thay đổi: (có hiệu lực trong 2 giờ)<br><br>
	
	<div style="font-size:30px;margin-left:60px;padding:5px 20px;border:solid 2px #aaa;background-color:#eee;display:inline-block">{{.Code}}</div><br><br><br>
	
	Nếu không phải bạn, hãy bỏ qua email này.<br><br>
	
	Đội ngũ {{.WlName}}
	`))
}
