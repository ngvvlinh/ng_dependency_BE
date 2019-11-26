package status3

import (
	"etop.vn/common/jsonx"
)

type Status int32

const (
	Status_Z Status = 0
	Status_P Status = 1
	Status_N Status = 127
)

var Status_name = map[int32]string{
	0:   "Z",
	1:   "P",
	127: "N",
}

var Status_value = map[string]int32{
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
	return jsonx.EnumName(Status_name, int32(x))
}

func (x *Status) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(Status_value, data, "Status")
	if err != nil {
		return err
	}
	*x = Status(value)
	return nil
}
