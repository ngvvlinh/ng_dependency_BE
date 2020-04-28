package convert

import (
	"regexp"

	"o.o/api/etc/logging/smslog"
	cm "o.o/backend/pkg/common"
	"o.o/common/l"
)

// +gen:convert: o.o/backend/com/etc/logging/smslog/model  -> o.o/api/etc/logging/smslog
// +gen:convert:  o.o/api/etc/logging/smslog

var ll = l.New()

var re = regexp.MustCompile(`[0-9]`)

func createSmsLog(args *smslog.CreateSmsArgs, out *smslog.SmsLog) {
	apply_smslog_CreateSmsArgs_smslog_SmsLog(args, out)
	out.Content = re.ReplaceAllString(args.Content, "*")
	out.ID = cm.NewID()
}
