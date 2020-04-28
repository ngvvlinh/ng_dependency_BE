package sqlstore

import (
	"o.o/backend/com/main/identity/convert"
	"o.o/backend/pkg/common/conversion"
)

var scheme = conversion.Build(convert.RegisterConversions)
