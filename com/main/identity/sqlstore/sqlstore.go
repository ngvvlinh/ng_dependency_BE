package sqlstore

import (
	"etop.vn/backend/com/main/identity/convert"
	"etop.vn/backend/pkg/common/conversion"
)

var scheme = conversion.Build(convert.RegisterConversions)
