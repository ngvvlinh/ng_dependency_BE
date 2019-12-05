package change_type

import (
	"etop.vn/common/jsonx"
)

type ChangeType int

const (
	ChangeType_unknown ChangeType = 0
	ChangeType_update  ChangeType = 1
	ChangeType_create  ChangeType = 2
	ChangeType_delete  ChangeType = 3
)

var ChangeType_name = map[int]string{
	0: "unknown",
	1: "update",
	2: "create",
	3: "delete",
}

var ChangeType_value = map[string]int{
	"unknown": 0,
	"update":  1,
	"create":  2,
	"delete":  3,
}

func (x ChangeType) Enum() *ChangeType {
	p := new(ChangeType)
	*p = x
	return p
}

func (x ChangeType) String() string {
	return jsonx.EnumName(ChangeType_name, int(x))
}

func (x *ChangeType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(ChangeType_value, data, "ChangeType")
	if err != nil {
		return err
	}
	*x = ChangeType(value)
	return nil
}
