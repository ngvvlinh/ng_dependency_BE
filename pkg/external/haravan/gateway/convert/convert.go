package convert

import (
	"etop.vn/api/external/haravan"
	pborder "etop.vn/backend/pb/etop/order"
	"etop.vn/backend/pkg/etop/model"
	shippingmodel "etop.vn/backend/pkg/services/shipping/model"
)

func ToPbOrderCustomer(in *haravan.Address) *pborder.OrderCustomer {
	if in == nil {
		return nil
	}
	return &pborder.OrderCustomer{
		FirstName: "",
		LastName:  "",
		FullName:  in.Name,
		Email:     "",
		Phone:     in.Phone,
		Gender:    0,
	}
}

func ToPbOrderAddress(in *haravan.Address) *pborder.OrderAddress {
	if in == nil {
		return nil
	}
	return &pborder.OrderAddress{
		FullName: in.Name,
		Phone:    in.Phone,
		Country:  in.Country,
		Province: in.Province,
		District: in.District,
		Ward:     in.Ward,
		Zip:      in.Zip,
		Address1: in.Address1,
		Address2: in.Address2,
	}
}

func ToPbCreateOrderLine(in *haravan.Item) *pborder.CreateOrderLine {
	if in == nil {
		return nil
	}
	return &pborder.CreateOrderLine{
		ProductName:  in.Name,
		Quantity:     int32(in.Quantity),
		ListPrice:    int32(in.Price),
		RetailPrice:  int32(in.Price),
		PaymentPrice: int32(in.Price),
	}
}

func ToPbCreateOrderLines(ins []*haravan.Item) (outs []*pborder.CreateOrderLine) {
	for _, in := range ins {
		outs = append(outs, ToPbCreateOrderLine(in))
	}
	return
}

func ToFulfillmentStatus(in model.ShippingState) haravan.FulfillmentState {
	switch in {
	case model.StateCreated:
		return haravan.PendingState
	case model.StatePicking:
		return haravan.PickingState
	case model.StateHolding:
		return haravan.DeliveringState
	case model.StateDelivering:
		return haravan.DeliveringState
	case model.StateDelivered:
		return haravan.DeliveredState
	case model.StateReturning:
		return haravan.WaitingForReturnState
	case model.StateReturned:
		return haravan.ReturnState
	case model.StateCancelled:
		return haravan.CancelState
	default:
		return ""
	}
}

func ToCODStatus(ffm *shippingmodel.Fulfillment) haravan.CODStatus {
	if ffm == nil {
		return ""
	}
	if ffm.TotalCODAmount == 0 {
		return haravan.NoneStatus
	}
	if !ffm.CODEtopTransferedAt.IsZero() {
		return haravan.CODReceiptStatus
	}
	switch ffm.EtopPaymentStatus {
	case model.S4Zero:
		return haravan.CODPendingStatus
	case model.S4Positive:
		return haravan.CODPaidStatus
	case model.S4Negative:
		return ""
	case model.S4SuperPos:
		return haravan.CODPendingStatus
	}
	return ""
}
