package sqlstore

import (
	"o.o/backend/com/etc/logging/smslog/convert"
	"o.o/backend/pkg/common/conversion"
)

var scheme = conversion.Build(convert.RegisterConversions)
