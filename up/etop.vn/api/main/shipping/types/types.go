package types

import (
	"errors"
	"time"

	orderingtypes "etop.vn/api/main/ordering/types"
)

type TryOn int

const (
	TryOnUnknown TryOn = 0
	TryOnNone    TryOn = 1
	TryOnOpen    TryOn = 2
	TryOnTry     TryOn = 3
)

type FeeLineType int

const (
	FeeLineTypeOther         FeeLineType = 0
	FeeLineTypeMain          FeeLineType = 1
	FeeLineTypeReturn        FeeLineType = 2
	FeeLineTypeAdjustment    FeeLineType = 3
	FeeLineTypeCods          FeeLineType = 4
	FeeLineTypeInsurance     FeeLineType = 5
	FeeLineTypeAddressChange FeeLineType = 6
	FeeLineTypeDiscount      FeeLineType = 7
)

var TryOnCode_name = map[int]string{
	0: "unknown",
	1: "none",
	2: "open",
	3: "try",
}

var TryOnCode_value = map[string]int{
	"unknown": 0,
	"none":    1,
	"open":    2,
	"try":     3,
}

func TryOnFromString(s string) (TryOn, error) {
	t := TryOnCode_value[s]
	if t == 0 {
		return 0, errors.New("invalid tryon code")
	}
	return TryOn(t), nil
}

var FeeLineType_name = map[int]string{
	0: "other",
	1: "main",
	2: "return",
	3: "adjustment",
	4: "cods",
	5: "insurance",
	6: "address_change",
	7: "discount",
}

var FeeLineType_value = map[string]int{
	"other":          0,
	"main":           1,
	"return":         2,
	"adjustment":     3,
	"cods":           4,
	"insurance":      5,
	"address_change": 6,
	"discount":       7,
}

func FeelineTypeFromString(s string) FeeLineType {
	f, ok := FeeLineType_value[s]
	if !ok {
		f = 0
	}
	return FeeLineType(f)
}

func (t FeeLineType) String() string {
	return FeeLineType_name[int(t)]
}

type ShippingInfo struct {
	PickupAddress       *orderingtypes.Address
	ReturnAddress       *orderingtypes.Address
	ShippingServiceName string
	ShippingServiceCode string
	ShippingServiceFee  int
	Carrier             string
	IncludeInsurance    bool
	TryOn               TryOn
	ShippingNote        string
	CODAmount           int
	GrossWeight         int
	Length              int
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
	CodAmount        int
	IncludeInsurance bool
}

type FeeLine struct {
	ShippingFeeType     FeeLineType
	Cost                int
	ExternalServiceName string
	ExternalServiceType string
}

func TotalFee(lines []*FeeLine) int {
	res := 0
	for _, line := range lines {
		res += int(line.Cost)
	}
	return res
}
