package try_on

import (
	"etop.vn/common/jsonx"
)

// +enum
type TryOnCode int

const (
	// +enum=unknown
	TryOnCode_unknown TryOnCode = 0

	// +enum=none
	TryOnCode_none TryOnCode = 1

	// +enum=open
	TryOnCode_open TryOnCode = 2

	// +enum=try
	TryOnCode_try TryOnCode = 3
)

var TryOnCode_name = map[int]string{
	0: "unknown",
	1: "none",
	2: "open",
	3: "try",
}

var TryOnCode_value = map[string]int{
	"unknown": 0,
	"none":    1,
	"open":    2,
	"try":     3,
}

func (x TryOnCode) Enum() *TryOnCode {
	p := new(TryOnCode)
	*p = x
	return p
}

func (x TryOnCode) String() string {
	return jsonx.EnumName(TryOnCode_name, int(x))
}

func (x *TryOnCode) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(TryOnCode_value, data, "TryOnCode")
	if err != nil {
		return err
	}
	*x = TryOnCode(value)
	return nil
}