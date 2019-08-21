package httpreq

import (
	"encoding/json"

	"etop.vn/common/l"
)

var ll = l.New()

func IsNullJsonRaw(data json.RawMessage) bool {
	return len(data) == 0 ||
		len(data) == 4 && string(data) == "null"
}
