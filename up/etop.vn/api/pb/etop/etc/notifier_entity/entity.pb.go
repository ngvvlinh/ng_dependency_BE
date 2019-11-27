package notifier_entity

import (
	"etop.vn/common/jsonx"
)

type NotifierEntity int

const (
	NotifierEntity_unknown                    NotifierEntity = 0
	NotifierEntity_fulfillment                NotifierEntity = 1
	NotifierEntity_money_transaction_shipping NotifierEntity = 2
)

var NotifierEntity_name = map[int]string{
	0: "unknown",
	1: "fulfillment",
	2: "money_transaction_shipping",
}

var NotifierEntity_value = map[string]int{
	"unknown":                    0,
	"fulfillment":                1,
	"money_transaction_shipping": 2,
}

func (x NotifierEntity) Enum() *NotifierEntity {
	p := new(NotifierEntity)
	*p = x
	return p
}

func (x NotifierEntity) String() string {
	return jsonx.EnumName(NotifierEntity_name, int(x))
}

func (x *NotifierEntity) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(NotifierEntity_value, data, "NotifierEntity")
	if err != nil {
		return err
	}
	*x = NotifierEntity(value)
	return nil
}
