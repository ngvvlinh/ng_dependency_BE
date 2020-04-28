package sqlstore

import (
	"o.o/backend/com/main/inventory/convert"
	"o.o/backend/pkg/common/conversion"
)

var scheme = conversion.Build(convert.RegisterConversions)
