package httpreq

import (
	"encoding/json"
)

func IsNullJsonRaw(data json.RawMessage) bool {
	return len(data) == 0 ||
		len(data) == 4 && string(data) == "null"
}
