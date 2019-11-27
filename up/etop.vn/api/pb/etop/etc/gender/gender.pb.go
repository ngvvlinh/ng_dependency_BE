package gender

import (
	"etop.vn/common/jsonx"
)

type Gender int

const (
	Gender_unknown Gender = 0
	Gender_male    Gender = 1
	Gender_female  Gender = 2
	Gender_other   Gender = 3
)

var Gender_name = map[int]string{
	0: "unknown",
	1: "male",
	2: "female",
	3: "other",
}

var Gender_value = map[string]int{
	"unknown": 0,
	"male":    1,
	"female":  2,
	"other":   3,
}

func (x Gender) Enum() *Gender {
	p := new(Gender)
	*p = x
	return p
}

func (x Gender) String() string {
	return jsonx.EnumName(Gender_name, int(x))
}

func (x *Gender) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(Gender_value, data, "Gender")
	if err != nil {
		return err
	}
	*x = Gender(value)
	return nil
}
