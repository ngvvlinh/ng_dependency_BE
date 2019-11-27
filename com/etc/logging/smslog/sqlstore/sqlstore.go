package sqlstore

import (
	"etop.vn/backend/com/etc/logging/smslog/convert"
	"etop.vn/backend/pkg/common/conversion"
)

var scheme = conversion.Build(convert.RegisterConversions)
