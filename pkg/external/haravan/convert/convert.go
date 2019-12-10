package convert

import (
	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/status4"
	haravanclient "etop.vn/backend/pkg/integration/haravan/client"
)

func ToFulfillmentState(in shipping.State) haravanclient.FulfillmentState {
	switch in {
	case shipping.Created:
		return haravanclient.PendingState
	case shipping.Picking:
		return haravanclient.PickingState
	case shipping.Holding:
		return haravanclient.DeliveringState
	case shipping.Delivering:
		return haravanclient.DeliveringState
	case shipping.Delivered:
		return haravanclient.DeliveredState
	case shipping.Returning:
		return haravanclient.WaitingForReturnState
	case shipping.Returned:
		return haravanclient.ReturnState
	case shipping.Cancelled:
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
