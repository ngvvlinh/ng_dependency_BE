package convert

import (
	"etop.vn/backend/pkg/etop/model"
	haravanclient "etop.vn/backend/pkg/integration/haravan/client"
)

func ToFulfillmentState(in model.ShippingState) haravanclient.FulfillmentState {
	switch in {
	case model.StateCreated:
		return haravanclient.PendingState
	case model.StatePicking:
		return haravanclient.PickingState
	case model.StateHolding:
		return haravanclient.DeliveringState
	case model.StateDelivering:
		return haravanclient.DeliveringState
	case model.StateDelivered:
		return haravanclient.DeliveredState
	case model.StateReturning:
		return haravanclient.WaitingForReturnState
	case model.StateReturned:
		return haravanclient.ReturnState
	case model.StateCancelled:
		return haravanclient.CancelState
	default:
		return ""
	}
}

func ToCODStatus(etopPaymentStatus model.Status4) haravanclient.PaymentStatus {
	switch etopPaymentStatus {
	case model.S4Zero:
		return haravanclient.PaymentPendingStatus
	case model.S4Positive:
		return haravanclient.PaymentReceiptStatus
	case model.S4Negative:
		return ""
	case model.S4SuperPos:
		return haravanclient.PaymentPendingStatus
	}
	return ""
}
