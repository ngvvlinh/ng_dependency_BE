// Code generated by split-pbgo-def. DO NOT EDIT.

package v1

import (
	v11 "etop.vn/api/main/etop/v1"
	types "etop.vn/api/main/ordering/v1/types"
	types1 "etop.vn/api/main/shipnow/v1/types"
	types2 "etop.vn/api/main/shipping/v1/types"
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
	Id                  int64                   `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId              int64                   `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
	PartnerId           int64                   `protobuf:"varint,3,opt,name=partner_id,json=partnerId" json:"partner_id"`
	PickupAddress       *types.Address          `protobuf:"bytes,4,opt,name=pickup_address,json=pickupAddress" json:"pickup_address,omitempty"`
	DeliveryPoints      []*types1.DeliveryPoint `protobuf:"bytes,5,rep,name=delivery_points,json=deliveryPoints" json:"delivery_points,omitempty"`
	Carrier             string                  `protobuf:"bytes,6,opt,name=carrier" json:"carrier"`
	ShippingServiceCode string                  `protobuf:"bytes,7,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingServiceFee  int32                   `protobuf:"varint,8,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	types2.WeightInfo   `protobuf:"bytes,9,opt,name=weight_info,json=weightInfo,embedded=weight_info" json:"weight_info"`
	ValueInfo           types2.ValueInfo `protobuf:"bytes,10,opt,name=value_info,json=valueInfo" json:"value_info"`
	ShippingNote        string           `protobuf:"bytes,11,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	RequestPickupAt     *v1.Timestamp    `protobuf:"bytes,12,opt,name=request_pickup_at,json=requestPickupAt" json:"request_pickup_at,omitempty"`
	Status              v11.Status5      `protobuf:"varint,13,opt,name=status,enum=etop.vn.api.main.etop.v1.status5.Status5" json:"status"`
	ShippingStatus      v11.Status5      `protobuf:"varint,14,opt,name=shipping_status,json=shippingStatus,enum=etop.vn.api.main.etop.v1.status5.Status5" json:"shipping_status"`
	ShippingCode        string           `protobuf:"bytes,15,opt,name=shipping_code,json=shippingCode" json:"shipping_code"`
	ShippingState       types1.State     `protobuf:"varint,16,opt,name=shipping_state,json=shippingState,enum=etop.vn.api.main.shipnow.v1.state.State" json:"shipping_state"`
	ConfirmStatus       v11.Status3      `protobuf:"varint,17,opt,name=confirm_status,json=confirmStatus,enum=etop.vn.api.main.etop.v1.status3.Status3" json:"confirm_status"`
}

type CreateShipnowFulfillmentCommand struct {
	OrderIds            []int64        `protobuf:"varint,1,rep,name=order_ids,json=orderIds" json:"order_ids,omitempty"`
	Carrier             string         `protobuf:"bytes,2,opt,name=carrier" json:"carrier"`
	ShopId              int64          `protobuf:"varint,3,opt,name=shop_id,json=shopId" json:"shop_id"`
	ShippingServiceCode string         `protobuf:"bytes,4,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingServiceFee  int32          `protobuf:"varint,5,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	ShippingNote        string         `protobuf:"bytes,6,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	RequestPickupAt     *v1.Timestamp  `protobuf:"bytes,7,opt,name=request_pickup_at,json=requestPickupAt" json:"request_pickup_at,omitempty"`
	PickupAddress       *types.Address `protobuf:"bytes,8,opt,name=pickup_address,json=pickupAddress" json:"pickup_address,omitempty"`
}

type UpdateShipnowFulfillmentCommand struct {
	Id                  int64          `protobuf:"varint,1,opt,name=id" json:"id"`
	OrderIds            []int64        `protobuf:"varint,2,rep,name=order_ids,json=orderIds" json:"order_ids,omitempty"`
	Carrier             string         `protobuf:"bytes,3,opt,name=carrier" json:"carrier"`
	ShopId              int64          `protobuf:"varint,4,opt,name=shop_id,json=shopId" json:"shop_id"`
	ShippingServiceCode string         `protobuf:"bytes,5,opt,name=shipping_service_code,json=shippingServiceCode" json:"shipping_service_code"`
	ShippingServiceFee  int32          `protobuf:"varint,6,opt,name=shipping_service_fee,json=shippingServiceFee" json:"shipping_service_fee"`
	ShippingNote        string         `protobuf:"bytes,7,opt,name=shipping_note,json=shippingNote" json:"shipping_note"`
	RequestPickupAt     *v1.Timestamp  `protobuf:"bytes,8,opt,name=request_pickup_at,json=requestPickupAt" json:"request_pickup_at,omitempty"`
	PickupAddress       *types.Address `protobuf:"bytes,9,opt,name=pickup_address,json=pickupAddress" json:"pickup_address,omitempty"`
}

type ConfirmShipnowFulfillmentCommand struct {
	Id     int64 `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId int64 `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
}

type CancelShipnowFulfillmentCommand struct {
	Id           int64  `protobuf:"varint,1,opt,name=id" json:"id"`
	CancelReason string `protobuf:"bytes,2,opt,name=cancel_reason,json=cancelReason" json:"cancel_reason"`
}

type GetShipnowFulfillmentQueryArgs struct {
	Id     int64 `protobuf:"varint,1,opt,name=id" json:"id"`
	ShopId int64 `protobuf:"varint,2,opt,name=shop_id,json=shopId" json:"shop_id"`
}

type GetShipnowFulfillmentsQueryArgs struct {
	ShopId int64 `protobuf:"varint,1,opt,name=shop_id,json=shopId" json:"shop_id"`
}

type GetShipnowFulfillmentQueryResult struct {
	ShipnowFulfillment *ShipnowFulfillment `protobuf:"bytes,1,opt,name=shipnow_fulfillment,json=shipnowFulfillment" json:"shipnow_fulfillment,omitempty"`
}

type GetShipnowFulfillmentsQueryResult struct {
	ShipnowFulfillments []*ShipnowFulfillment `protobuf:"bytes,1,rep,name=shipnow_fulfillments,json=shipnowFulfillments" json:"shipnow_fulfillments,omitempty"`
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
