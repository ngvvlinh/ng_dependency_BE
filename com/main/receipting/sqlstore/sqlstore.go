package sqlstore

import (
	"o.o/backend/com/main/receipting/convert"
	"o.o/backend/pkg/common/conversion"
)

var scheme = conversion.Build(convert.RegisterConversions)
