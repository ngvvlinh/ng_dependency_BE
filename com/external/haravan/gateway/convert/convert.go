package convert

import (
	"etop.vn/api/external/haravan"
	"etop.vn/api/main/location"
	pbexternal "etop.vn/api/pb/external"
	"etop.vn/capi/dot"
)

func ToPbExternalAddress(in *haravan.Address, loc *location.LocationQueryResult) *pbexternal.OrderAddress {
	if in == nil {
		return nil
	}
	return &pbexternal.OrderAddress{
		FullName: in.Name,
		Phone:    in.Phone,
		Province: loc.Province.Name,
		District: loc.District.Name,
		Ward:     loc.Ward.Name,
		Address1: in.Address1,
		Address2: in.Address2,
	}
}

func ToPbExternalCreateOrderLine(in *haravan.Item) *pbexternal.OrderLine {
	if in == nil {
		return nil
	}
	return &pbexternal.OrderLine{
		ProductName:  in.Name,
		Quantity:     in.Quantity,
		ListPrice:    int(in.Price),
		RetailPrice:  int(in.Price),
		PaymentPrice: dot.Int(int(in.Price)),
	}
}

func ToPbExternalCreateOrderLines(ins []*haravan.Item) (outs []*pbexternal.OrderLine) {
	for _, in := range ins {
		outs = append(outs, ToPbExternalCreateOrderLine(in))
	}
	return
}
