package carrier

import (
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	"o.o/capi/dot"
)

type ConnectionSignInArgs struct {
	ConnectionID dot.ID
	Identifier   string // email or phone
	Password     string
	OTP          string
}

type ShopConnectionSignInArgs struct {
	ConnectionID dot.ID
	ShopID       dot.ID
	Identifier   string // email or phone
	Password     string
}

type ShopConnectionSignInWithOTPArgs struct {
	ConnectionID dot.ID
	ShopID       dot.ID
	Identifier   string
	OTP          string
}

type ConnectionSignUpArgs struct {
	ConnectionID dot.ID
	Name         string
	Identifier   string
	Password     string
	Phone        string
	Province     string
	District     string
	Address      string
}

type ShopConnectionSignUpArgs struct {
	ConnectionID dot.ID
	ShopID       dot.ID
	Name         string
	Identifier   string
	Password     string
	Phone        string
	Province     string
	District     string
	Address      string
}

type GetShippingServicesArgs struct {
	AccountID           dot.ID
	ShipmentPriceListID dot.ID
	ConnectionIDs       []dot.ID
	FromDistrictCode    string
	FromProvinceCode    string
	FromWardCode        string
	ToDistrictCode      string
	ToProvinceCode      string
	ToWardCode          string

	ChargeableWeight int
	Length           int
	Width            int
	Height           int
	IncludeInsurance bool
	InsuranceValue   int
	BasketValue      int
	CODAmount        int

	Coupon string
}

type UpdateFulfillmentCODArgs struct {
	ConnectionID  dot.ID
	FulfillmentID dot.ID
	CODAmount     int
}

func (a *GetShippingServicesArgs) ToShipmentServiceArgs(arbitraryID, accountID dot.ID) *carriertypes.GetShippingServicesArgs {
	return &carriertypes.GetShippingServicesArgs{
		ArbitraryID:      arbitraryID,
		AccountID:        accountID,
		FromDistrictCode: a.FromDistrictCode,
		FromWardCode:     a.FromWardCode,
		ToDistrictCode:   a.ToDistrictCode,
		ToWardCode:       a.ToWardCode,
		ChargeableWeight: a.ChargeableWeight,
		Length:           a.Length,
		Width:            a.Width,
		Height:           a.Height,
		IncludeInsurance: a.IncludeInsurance,
		InsuranceValue:   a.InsuranceValue,
		BasketValue:      a.BasketValue,
		CODAmount:        a.CODAmount,
		Coupon:           a.Coupon,
	}
}

type GetPriceListPromotionArgs struct {
	ShopID           dot.ID
	FromProvinceCode string
	ConnectionID     dot.ID
}
