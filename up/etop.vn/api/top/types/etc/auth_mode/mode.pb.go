package auth_mode

import (
	"etop.vn/common/jsonx"
)

// +enum
type AuthMode int

const (
	// +enum=default
	AuthMode_default AuthMode = 0

	// +enum=manual
	AuthMode_manual AuthMode = 1
)

var AuthMode_name = map[int]string{
	0: "default",
	1: "manual",
}

var AuthMode_value = map[string]int{
	"default": 0,
	"manual":  1,
}

func (x AuthMode) Enum() *AuthMode {
	p := new(AuthMode)
	*p = x
	return p
}

func (x AuthMode) String() string {
	return jsonx.EnumName(AuthMode_name, int(x))
}

func (x *AuthMode) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(AuthMode_value, data, "AuthMode")
	if err != nil {
		return err
	}
	*x = AuthMode(value)
	return nil
}
