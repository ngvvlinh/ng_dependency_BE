// Code generated by split-pbgo-def. DO NOT EDIT.

package v1

import (
	types "etop.vn/api/main/ordering/v1/types"
	types1 "etop.vn/api/main/shipping/v1/types"
	v1 "etop.vn/api/meta/v1"
)

type EventDataEnum int32

const (
	Created               EventDataEnum = 1
	ConfirmationRequested EventDataEnum = 4
	ConfirmationAccepted  EventDataEnum = 5
	ConfirmationRejected  EventDataEnum = 6
	CancellationRequested EventDataEnum = 7
	CancellationAccepted  EventDataEnum = 8
	CancellationRejected  EventDataEnum = 9
)

func (_ EventData_Created) GetEnumTag() EventDataEnum               { return 1 }
func (_ EventData_ConfirmationRequested) GetEnumTag() EventDataEnum { return 4 }
func (_ EventData_ConfirmationAccepted) GetEnumTag() EventDataEnum  { return 5 }
func (_ EventData_ConfirmationRejected) GetEnumTag() EventDataEnum  { return 6 }
func (_ EventData_CancellationRequested) GetEnumTag() EventDataEnum { return 7 }
func (_ EventData_CancellationAccepted) GetEnumTag() EventDataEnum  { return 8 }
func (_ EventData_CancellationRejected) GetEnumTag() EventDataEnum  { return 9 }

type ShipnowFulfillment struct {
	Id                  int64            `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId              int64            `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	PartnerId           int64            `protobuf:"varint,3,opt,name=partner_id,json=partnerId" json:"partner_id"`
	PickupAddress       *types.Address   `protobuf:"bytes,4,opt,name=pickup_address,json=pickupAddress" json:"pickup_address,omitempty"`
	DeliveryPoints      []*DeliveryPoint `protobuf:"bytes,5,rep,name=delivery_points,json=deliveryPoints" json:"delivery_points,omitempty"`
	Carrier             string           `protobuf:"bytes,6,opt,name=carrier" json:"carrier"`
	ShippingServiceCode string           `protobuf:"bytes,7,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingServiceFee  int32            `protobuf:"varint,8,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	types1.WeightInfo   `protobuf:"bytes,9,opt,name=weight_info,json=weightInfo,embedded=weight_info" json:"weight_info"`
	ValueInfo           types1.ValueInfo `protobuf:"bytes,10,opt,name=value_info,json=valueInfo" json:"value_info"`
	ShippingNote        string           `protobuf:"bytes,11,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	RequestPickupAt     *v1.Timestamp    `protobuf:"bytes,12,opt,name=request_pickup_at,json=requestPickupAt" json:"request_pickup_at,omitempty"`
}

type DeliveryPoint struct {
	ShippingAddress   *types.Address    `protobuf:"bytes,1,opt,name=shipping_address,json=shippingAddress" json:"shipping_address,omitempty"`
	Lines             []*types.ItemLine `protobuf:"bytes,2,rep,name=lines" json:"lines,omitempty"`
	ShippingNote      string            `protobuf:"bytes,3,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	OrderId           int64             `protobuf:"varint,4,opt,name=order_id,json=orderId" json:"order_id"`
	types1.WeightInfo `protobuf:"bytes,5,opt,name=weight_info,json=weightInfo,embedded=weight_info" json:"weight_info"`
	types1.ValueInfo  `protobuf:"bytes,6,opt,name=value_info,json=valueInfo,embedded=value_info" json:"value_info"`
	TryOn             types1.TryOnCode `protobuf:"varint,11,opt,name=try_on,json=tryOn,enum=etop.vn.api.main.shipping.v1.tryon.TryOnCode" json:"try_on"`
}

type CreateShipnowFulfillmentCommand struct {
	ShipnowFulfillment `protobuf:"bytes,1,opt,name=data,embedded=data" json:"data"`
}

type ConfirmShipnowFulfillmentCommand struct {
	Id int64 `protobuf:"varint,1,opt,name=id" json:"id"`
}

type CancelShipnowFulfillmentCommand struct {
	Id           int64  `protobuf:"varint,1,opt,name=id" json:"id"`
	CancelReason string `protobuf:"bytes,2,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
}

type GetShipnowFulfillmentQueryArgs struct {
	Id int64 `protobuf:"varint,1,opt,name=id" json:"id"`
}

type GetShipnowFulfillmentQueryResult struct {
}

type ShipnowEvent struct {
	Id                   int64      `protobuf:"varint,1,opt,name=id" json:"id"`
	Uuid                 v1.UUID    `protobuf:"bytes,2,opt,name=uuid" json:"uuid"`
	CorrelationId        v1.UUID    `protobuf:"bytes,3,opt,name=correlation_id,json=correlationId" json:"correlation_id"`
	ShipnowFulfillmentId int64      `protobuf:"varint,4,opt,name=shipnow_fulfillment_id,json=shipnowFulfillmentId" json:"shipnow_fulfillment_id"`
	Type                 int32      `protobuf:"varint,10,opt,name=type" json:"type"`
	Data                 *EventData `protobuf:"bytes,11,opt,name=data" json:"data,omitempty"`
}

type EventData struct {
	// Types that are valid to be assigned to Data:
	//	*EventData_Created
	//	*EventData_ConfirmationRequested
	//	*EventData_ConfirmationAccepted
	//	*EventData_ConfirmationRejected
	//	*EventData_CancellationRequested
	//	*EventData_CancellationAccepted
	//	*EventData_CancellationRejected
	Data isEventData_Data `protobuf_oneof:"data"`
}

type isEventData_Data interface {
	isEventData_Data()
}

type EventData_Created struct {
	Created *CreatedData `protobuf:"bytes,1,opt,name=Created,oneof"`
}

type EventData_ConfirmationRequested struct {
	ConfirmationRequested *ConfirmationRequestedData `protobuf:"bytes,4,opt,name=ConfirmationRequested,oneof"`
}

type EventData_ConfirmationAccepted struct {
	ConfirmationAccepted *ConfirmationAcceptedData `protobuf:"bytes,5,opt,name=ConfirmationAccepted,oneof"`
}

type EventData_ConfirmationRejected struct {
	ConfirmationRejected *ConfirmationRejectedData `protobuf:"bytes,6,opt,name=ConfirmationRejected,oneof"`
}

type EventData_CancellationRequested struct {
	CancellationRequested *CancellationRequestedData `protobuf:"bytes,7,opt,name=CancellationRequested,oneof"`
}

type EventData_CancellationAccepted struct {
	CancellationAccepted *CancellationAcceptedData `protobuf:"bytes,8,opt,name=CancellationAccepted,oneof"`
}

type EventData_CancellationRejected struct {
	CancellationRejected *CancellationRejectedData `protobuf:"bytes,9,opt,name=CancellationRejected,oneof"`
}

type CreatedData struct {
}

type ConfirmationRequestedData struct {
}

type ConfirmationAcceptedData struct {
}

type ConfirmationRejectedData struct {
}

type CancellationRequestedData struct {
}

type CancellationAcceptedData struct {
}

type CancellationRejectedData struct {
}
