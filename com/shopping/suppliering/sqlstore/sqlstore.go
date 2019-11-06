package sqlstore

import (
	"etop.vn/backend/com/shopping/suppliering/convert"
	"etop.vn/backend/pkg/common/conversion"
)

var scheme = conversion.Build(convert.RegisterConversions)
