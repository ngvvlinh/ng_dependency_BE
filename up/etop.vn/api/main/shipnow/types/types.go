package types

import (
	shipnowv1types "etop.vn/api/main/shipnow/v1/types"
)

type DeliveryPoint = shipnowv1types.DeliveryPoint

type State = shipnowv1types.State

func StateFromString(s string) State {
	st, ok := shipnowv1types.State_value[s]
	if !ok {
		return shipnowv1types.State_unknown
	}
	return State(st)
}

const (
	StateDefault       = shipnowv1types.State_default
	StateCreated       = shipnowv1types.State_created
	StateAssigning     = shipnowv1types.State_assigning
	StatePicking       = shipnowv1types.State_picking
	StateDelivering    = shipnowv1types.State_delivering
	StateDelivered     = shipnowv1types.State_delivered
	StateReturning     = shipnowv1types.State_returning
	StateReturned      = shipnowv1types.State_returned
	StateUnknown       = shipnowv1types.State_unknown
	StateUndeliverable = shipnowv1types.State_undeliverable
	StateCancelled     = shipnowv1types.State_cancelled
)
