package jsonx

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func UnmarshalJSONEnum(m map[string]int, data []byte, enumName string) (int, error) {
	if data[0] == '"' {
		var repr string
		if err := json.Unmarshal(data, &repr); err != nil {
			return -1, err
		}
		val, ok := m[repr]
		if !ok {
			return 0, fmt.Errorf("unrecognized enum %s value %q", enumName, repr)
		}
		return val, nil
	}

	var val int
	if err := json.Unmarshal(data, &val); err != nil {
		return 0, fmt.Errorf("cannot unmarshal %#q into enum %s", data, enumName)
	}
	return val, nil
}

func EnumName(m map[int]string, v int) string {
	s, ok := m[v]
	if ok {
		return s
	}
	return strconv.Itoa(v)
}
