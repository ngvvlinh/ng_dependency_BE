// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package substate

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumSubstateName = map[int]string{
	1:  "default",
	17: "pick_fail",
	25: "delivery_fail",
	37: "devivery_giveup",
	43: "return_fail",
	51: "cancelling",
}

var enumSubstateValue = map[string]int{
	"default":         1,
	"pick_fail":       17,
	"delivery_fail":   25,
	"devivery_giveup": 37,
	"return_fail":     43,
	"cancelling":      51,
}

var enumSubstateMapLabel = map[string]map[string]string{
	"default": {
		"RefName": "",
	},
	"pick_fail": {
		"RefName": "Lấy hàng thất bại",
	},
	"delivery_fail": {
		"RefName": "Giao hàng thất bại",
	},
	"devivery_giveup": {
		"RefName": "Giao hàng thất bại. Chờ trả hàng",
	},
	"return_fail": {
		"RefName": "Trả hàng thất bại",
	},
	"cancelling": {
		"RefName": "",
	},
}

func (e Substate) GetLabelRefName() string {
	val := enumSubstateName[int(e)]
	nameVal := enumSubstateMapLabel[val]
	return nameVal["RefName"]
}
func ParseSubstate(s string) (Substate, bool) {
	val, ok := enumSubstateValue[s]
	return Substate(val), ok
}

func ParseSubstateWithDefault(s string, d Substate) Substate {
	val, ok := enumSubstateValue[s]
	if !ok {
		return d
	}
	return Substate(val)
}

func (e Substate) Apply(d Substate) Substate {
	if e == 0 {
		return d
	}
	return e
}

func (e Substate) Enum() int {
	return int(e)
}

func (e Substate) Name() string {
	return enumSubstateName[e.Enum()]
}

func (e Substate) String() string {
	s, ok := enumSubstateName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("Substate(%v)", e.Enum())
}

func (e Substate) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumSubstateName[e.Enum()] + "\""), nil
}

func (e *Substate) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumSubstateValue, data, "Substate")
	if err != nil {
		return err
	}
	*e = Substate(value)
	return nil
}

func (e Substate) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *Substate) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumSubstateValue, src, "Substate")
	*e = (Substate)(value)
	return err
}

func (e Substate) Wrap() NullSubstate {
	return WrapSubstate(e)
}

func ParseSubstateWithNull(s dot.NullString, d Substate) NullSubstate {
	if !s.Valid {
		return NullSubstate{}
	}
	val, ok := enumSubstateValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return Substate(val).Wrap()
}

func WrapSubstate(enum Substate) NullSubstate {
	return NullSubstate{Enum: enum, Valid: true}
}

func (n NullSubstate) Apply(s Substate) Substate {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullSubstate) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullSubstate) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullSubstate) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullSubstate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
