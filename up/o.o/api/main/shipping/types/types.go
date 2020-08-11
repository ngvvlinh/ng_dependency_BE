package types

import (
	"time"

	orderingtypes "o.o/api/main/ordering/types"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

type ShippingInfo struct {
	PickupAddress       *orderingtypes.Address
	ReturnAddress       *orderingtypes.Address
	ShippingServiceName string
	ShippingServiceCode string
	ShippingServiceFee  int
	Carrier             shipping_provider.ShippingProvider
	IncludeInsurance    bool
	TryOn               try_on.TryOnCode
	ShippingNote        string
	CODAmount           int
	GrossWeight         int
	Length              int
	Width               int
	Height              int
	ChargeableWeight    int
}

type ShippingService struct {
	Carrier             string
	Code                string
	Fee                 int
	Name                string
	EstimatedPickupAt   time.Time
	EstimatedDeliveryAt time.Time
}

type WeightInfo struct {
	GrossWeight      int
	ChargeableWeight int
	Length           int
	Width            int
	Height           int
}

type ValueInfo struct {
	BasketValue      int
	CODAmount        int
	IncludeInsurance bool
	InsuranceValue   dot.NullInt
}

type FeeLine struct {
	ShippingFeeType     shipping_fee_type.ShippingFeeType
	Cost                int
	ExternalServiceName string
	ExternalServiceType string
}

func TotalFee(lines []*FeeLine) int {
	res := 0
	for _, line := range lines {
		res += line.Cost
	}
	return res
}
