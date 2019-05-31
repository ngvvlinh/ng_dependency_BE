// Code generated by split-pbgo-def. DO NOT EDIT.

package types

import (
	v1 "etop.vn/api/meta/v1"
)

type ShippingService struct {
	Carrier             string        `protobuf:"bytes,1,opt,name=carrier" json:"carrier"`
	Code                string        `protobuf:"bytes,2,opt,name=code" json:"code"`
	Fee                 int32         `protobuf:"varint,3,opt,name=fee" json:"fee"`
	Name                string        `protobuf:"bytes,4,opt,name=name" json:"name"`
	EstimatedPickupAt   *v1.Timestamp `protobuf:"bytes,5,opt,name=estimated_pickup_at,json=estimatedPickupAt" json:"estimated_pickup_at,omitempty"`
	EstimatedDeliveryAt *v1.Timestamp `protobuf:"bytes,6,opt,name=estimated_delivery_at,json=estimatedDeliveryAt" json:"estimated_delivery_at,omitempty"`
}

type WeightInfo struct {
	GrossWeight      int32 `protobuf:"varint,1,opt,name=gross_weight,json=grossWeight" json:"gross_weight"`
	ChargeableWeight int32 `protobuf:"varint,2,opt,name=chargeable_weight,json=chargeableWeight" json:"chargeable_weight"`
	Length           int32 `protobuf:"varint,3,opt,name=length" json:"length"`
	Width            int32 `protobuf:"varint,4,opt,name=width" json:"width"`
	Height           int32 `protobuf:"varint,5,opt,name=height" json:"height"`
}

type ValueInfo struct {
	BasketValue      int32 `protobuf:"varint,1,opt,name=basket_value,json=basketValue" json:"basket_value"`
	CodAmount        int32 `protobuf:"varint,2,opt,name=cod_amount,json=codAmount" json:"cod_amount"`
	IncludeInsurance bool  `protobuf:"varint,3,opt,name=include_insurance,json=includeInsurance" json:"include_insurance"`
}

type FeeLine struct {
	ShippingFeeType     FeeLineType `protobuf:"varint,1,opt,name=shipping_fee_type,json=shippingFeeType,enum=etop.vn.api.main.shipping.v1.feeline.FeeLineType" json:"shipping_fee_type"`
	Cost                int32       `protobuf:"varint,2,opt,name=cost" json:"cost"`
	ExternalServiceName string      `protobuf:"bytes,3,opt,name=external_service_name,json=externalServiceName" json:"external_service_name"`
	ExternalServiceType string      `protobuf:"bytes,4,opt,name=external_service_type,json=externalServiceType" json:"external_service_type"`
}
