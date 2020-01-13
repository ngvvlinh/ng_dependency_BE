package convert

import (
	"etop.vn/api/external/haravan"
	"etop.vn/api/main/location"
	"etop.vn/api/top/external/types"
)

func ToPbExternalAddress(in *haravan.Address, loc *location.LocationQueryResult) *types.OrderAddress {
	if in == nil {
		return nil
	}
	return &types.OrderAddress{
		FullName: in.Name,
		Phone:    in.Phone,
		Province: loc.Province.Name,
		District: loc.District.Name,
		Ward:     loc.Ward.Name,
		Address1: in.Address1,
		Address2: in.Address2,
	}
}

func ToPbExternalCreateOrderLine(in *haravan.Item) *types.OrderLine {
	if in == nil {
		return nil
	}
	return &types.OrderLine{
		ProductName:  in.Name,
		Quantity:     in.Quantity,
		ListPrice:    int(in.Price),
		RetailPrice:  int(in.Price),
		PaymentPrice: int(in.Price),
	}
}

func ToPbExternalCreateOrderLines(ins []*haravan.Item) (outs []*types.OrderLine) {
	for _, in := range ins {
		outs = append(outs, ToPbExternalCreateOrderLine(in))
	}
	return
}
