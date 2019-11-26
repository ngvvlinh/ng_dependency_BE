package shipping

import (
	"etop.vn/common/jsonx"
)

type State int32

const (
	State_default       State = 0
	State_created       State = 1
	State_confirmed     State = 2
	State_processing    State = 3
	State_picking       State = 4
	State_holding       State = 5
	State_returning     State = 6
	State_returned      State = 7
	State_delivering    State = 8
	State_delivered     State = 9
	State_unknown       State = 101
	State_undeliverable State = 126
	State_cancelled     State = 127
)

var State_name = map[int32]string{
	0:   "default",
	1:   "created",
	2:   "confirmed",
	3:   "processing",
	4:   "picking",
	5:   "holding",
	6:   "returning",
	7:   "returned",
	8:   "delivering",
	9:   "delivered",
	101: "unknown",
	126: "undeliverable",
	127: "cancelled",
}

var State_value = map[string]int32{
	"default":       0,
	"created":       1,
	"confirmed":     2,
	"processing":    3,
	"picking":       4,
	"holding":       5,
	"returning":     6,
	"returned":      7,
	"delivering":    8,
	"delivered":     9,
	"unknown":       101,
	"undeliverable": 126,
	"cancelled":     127,
}

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return jsonx.EnumName(State_name, int32(x))
}

func (x *State) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(State_value, data, "State")
	if err != nil {
		return err
	}
	*x = State(value)
	return nil
}
