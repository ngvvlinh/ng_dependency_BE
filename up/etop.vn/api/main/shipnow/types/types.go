package types

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/ordering/types"
	v1 "etop.vn/api/main/shipnow/carrier/types"
	shippingtypes "etop.vn/api/main/shipping/types"
)

type DeliveryPoint struct {
	ShippingAddress *types.Address
	Lines           []*types.ItemLine
	ShippingNote    string
	OrderId         int64
	OrderCode       string
	shippingtypes.WeightInfo
	shippingtypes.ValueInfo
	TryOn shippingtypes.TryOn
}

type ShipnowService struct {
	Carrier            v1.Carrier
	Name               string
	Code               string
	Fee                int32
	ExpectedPickupAt   time.Time
	ExpectedDeliveryAt time.Time
	Description        string
}

type State int32

const (
	StateDefault       State = 0
	StateCreated       State = 1
	StateAssigning     State = 2
	StatePicking       State = 3
	StateDelivering    State = 4
	StateDelivered     State = 5
	StateReturning     State = 6
	StateReturned      State = 7
	StateUnknown       State = 101
	StateUndeliverable State = 126
	StateCancelled     State = 127
)

func (s State) String() string {
	return State_name[int32(s)]
}

var State_name = map[int32]string{
	0:   "default",
	1:   "created",
	2:   "assigning",
	3:   "picking",
	4:   "delivering",
	5:   "delivered",
	6:   "returning",
	7:   "returned",
	101: "unknown",
	126: "undeliverable",
	127: "cancelled",
}

var State_value = map[string]int32{
	"default":       0,
	"created":       1,
	"assigning":     2,
	"picking":       3,
	"delivering":    4,
	"delivered":     5,
	"returning":     6,
	"returned":      7,
	"unknown":       101,
	"undeliverable": 126,
	"cancelled":     127,
}

func StateFromString(s string) State {
	st, ok := State_value[s]
	if !ok {
		return StateUnknown
	}
	return State(st)
}

func StateToString(s State) string {
	if s == 0 {
		return ""
	}
	return s.String()
}

func StateToStatus5(s State) etop.Status5 {
	switch s {
	case StateDefault:
	case StateCreated:
		return etop.S5Zero
	case StateCancelled:
		return etop.S5Negative
	case StateReturned:
	case StateReturning:
		return etop.S5NegSuper
	case StateDelivered:
		return etop.S5Positive
	default:
		return etop.S5SuperPos
	}
	return etop.S5SuperPos
}
