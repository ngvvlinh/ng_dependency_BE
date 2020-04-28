package carrier

import (
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	"o.o/capi/dot"
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
	BasketValue      int
	CODAmount        int
}

func (a *GetShippingServicesArgs) ToShipmentServiceArgs(arbitraryID, accountID dot.ID) *carriertypes.GetShippingServicesArgs {
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
		ArbitraryID:            arbitraryID,
		AccountID:              accountID,
		IncludeTopshipServices: false, // fill it
	}
}
