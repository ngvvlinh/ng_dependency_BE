package status3

import (
	"etop.vn/common/jsonx"
)

// +enum
type Status int

const (
	// +enum=Z
	Z Status = 0

	// +enum=P
	P Status = 1

	// +enum=N
	N Status = 127
)

var Status_name = map[int]string{
	0:   "Z",
	1:   "P",
	127: "N",
}

var Status_value = map[string]int{
	"Z": 0,
	"P": 1,
	"N": 127,
}

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return jsonx.EnumName(Status_name, int(x))
}

func (x *Status) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(Status_value, data, "Status")
	if err != nil {
		return err
	}
	*x = Status(value)
	return nil
}
