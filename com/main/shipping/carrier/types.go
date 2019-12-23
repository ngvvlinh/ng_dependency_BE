package carrier

import (
	carriertypes "etop.vn/backend/com/main/shipping/carrier/types"
	"etop.vn/capi/dot"
)

type ConnectionSignInArgs struct {
	ConnectionID dot.ID
	Email        string
	Password     string
}

type ShopConnectionSignInArgs struct {
	ConnectionID dot.ID
	ShopID       dot.ID
	Email        string
	Password     string
}

type ConnectionSignUpArgs struct {
	ConnectionID dot.ID
	Name         string
	Email        string
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
	Email        string
	Password     string
	Phone        string
	Province     string
	District     string
	Address      string
}

type GetShippingServicesArgs struct {
	ConnectionIDs    []dot.ID
	FromDistrictCode string
	ToDistrictCode   string

	ChargeableWeight int
	Length           int
	Width            int
	Height           int
	IncludeInsurance bool
	BasketValue      int
	CODAmount        int
}

func (a *GetShippingServicesArgs) ToShipmentServiceArgs() *carriertypes.GetShippingServicesArgs {
	return &carriertypes.GetShippingServicesArgs{
		FromDistrictCode:       a.FromDistrictCode,
		ToDistrictCode:         a.ToDistrictCode,
		ChargeableWeight:       a.ChargeableWeight,
		Length:                 a.Length,
		Width:                  a.Width,
		Height:                 a.Height,
		IncludeInsurance:       a.IncludeInsurance,
		BasketValue:            a.BasketValue,
		CODAmount:              a.CODAmount,
		ArbitraryID:            0,     // fill it
		AccountID:              0,     // fill it
		IncludeTopshipServices: false, // fill it
	}
}
