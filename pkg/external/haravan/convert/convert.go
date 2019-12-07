package convert

import (
	"etop.vn/api/top/types/etc/status4"
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

func ToCODStatus(etopPaymentStatus status4.Status) haravanclient.PaymentStatus {
	switch etopPaymentStatus {
	case status4.Z:
		return haravanclient.PaymentPendingStatus
	case status4.P:
		return haravanclient.PaymentReceiptStatus
	case status4.N:
		return ""
	case status4.S:
		return haravanclient.PaymentPendingStatus
	}
	return ""
}
