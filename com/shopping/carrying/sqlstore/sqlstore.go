package sqlstore

import (
	"o.o/backend/com/shopping/carrying/convert"
	"o.o/backend/pkg/common/conversion"
)

var scheme = conversion.Build(convert.RegisterConversions)
