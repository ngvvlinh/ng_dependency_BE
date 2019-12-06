package encode

import (
	"encoding/json"
	"fmt"
)

func UnmarshalJSONEnum(m map[string]int, data []byte, enumName string) (int, error) {
	return UnmarshalJSONEnumInt(m, data, enumName)
}

func UnmarshalJSONEnumInt(m map[string]int, data []byte, enumName string) (int, error) {
	if string(data) == "null" {
		return 0, nil
	}

	if data[0] == '"' {
		var repr string
		if err := json.Unmarshal(data, &repr); err != nil {
			return 0, err
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

func UnmarshalJSONEnumUint64(m map[string]uint64, data []byte, enumName string) (uint64, error) {
	if string(data) == "null" {
		return 0, nil
	}

	if data[0] == '"' {
		var repr string
		if err := json.Unmarshal(data, &repr); err != nil {
			return 0, err
		}
		val, ok := m[repr]
		if !ok {
			return 0, fmt.Errorf("unrecognized enum %s value %q", enumName, repr)
		}
		return val, nil
	}

	var val uint64
	if err := json.Unmarshal(data, &val); err != nil {
		return 0, fmt.Errorf("cannot unmarshal %#q into enum %s", data, enumName)
	}
	return val, nil
}

func ScanEnumInt(m map[string]int, src interface{}, enumName string) (int, error) {
	switch src := src.(type) {
	case nil:
		return 0, nil

	case []byte:
		value, ok := m[string(src)]
		if !ok {
			return 0, fmt.Errorf("can not read value %v into enum %v", string(src), enumName)
		}
		return value, nil

	case string:
		value, ok := m[src]
		if !ok {
			return 0, fmt.Errorf("can not read value %v into enum %v", src, enumName)
		}
		return value, nil

	case int64:
		return int(src), nil

	default:
		return 0, fmt.Errorf("can not read value of type %T into enum", src)
	}
}

func ScanEnumUint64(m map[string]uint64, src interface{}, enumName string) (uint64, error) {
	switch src := src.(type) {
	case nil:
		return 0, nil

	case []byte:
		value, ok := m[string(src)]
		if !ok {
			return 0, fmt.Errorf("can not read value %v into enum %v", string(src), enumName)
		}
		return value, nil

	case string:
		value, ok := m[src]
		if !ok {
			return 0, fmt.Errorf("can not read value %v into enum %v", src, enumName)
		}
		return value, nil

	case int64:
		return uint64(src), nil

	case uint64:
		return src, nil

	default:
		return 0, fmt.Errorf("can not read value of type %T into enum", src)
	}
}
