package convert

import (
	"regexp"

	"etop.vn/api/etc/logging/smslog"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/common/l"
)

// +gen:convert: etop.vn/backend/com/etc/logging/smslog/model -> etop.vn/api/etc/logging/smslog
// +gen:convert: etop.vn/api/etc/logging/smslog

var ll = l.New()

var re = regexp.MustCompile(`[0-9]`)

func createSmsLog(args *smslog.CreateSmsArgs, out *smslog.SmsLog) {
	apply_smslog_CreateSmsArgs_smslog_SmsLog(args, out)
	out.Content = re.ReplaceAllString(args.Content, "*")
	out.ID = cm.NewID()
}
