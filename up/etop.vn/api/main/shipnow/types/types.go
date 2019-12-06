package types

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/ordering/types"
	v1 "etop.vn/api/main/shipnow/carrier/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/capi/dot"
)

type DeliveryPoint struct {
	ShippingAddress *types.Address
	Lines           []*types.ItemLine
	ShippingNote    string
	OrderId         dot.ID
	OrderCode       string
	shippingtypes.WeightInfo
	shippingtypes.ValueInfo
	TryOn shippingtypes.TryOn
}

type ShipnowService struct {
	Carrier            v1.Carrier
	Name               string
	Code               string
	Fee                int
	ExpectedPickupAt   time.Time
	ExpectedDeliveryAt time.Time
	Description        string
}

// +enum
type State int

const (
	// +enum=default
	StateDefault State = 0

	// +enum=created
	StateCreated State = 1

	// +enum=assigning
	StateAssigning State = 2

	// +enum=picking
	StatePicking State = 3

	// +enum=delivering
	StateDelivering State = 4

	// +enum=delivered
	StateDelivered State = 5

	// +enum=returning
	StateReturning State = 6

	// +enum=returned
	StateReturned State = 7

	// +enum=unknown
	StateUnknown State = 101

	// +enum=undeliverable
	StateUndeliverable State = 126

	// +enum=cancelled
	StateCancelled State = 127
)

func StateFromString(s string) State {
	st, ok := enumStateValue[s]
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
