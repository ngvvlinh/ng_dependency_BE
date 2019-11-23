package convert

import (
	"etop.vn/api/external/haravan"
	"etop.vn/api/main/location"
	pbexternal "etop.vn/api/pb/external"
	cm "etop.vn/backend/pkg/common"
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
		Quantity:     int32(in.Quantity),
		ListPrice:    int32(in.Price),
		RetailPrice:  int32(in.Price),
		PaymentPrice: cm.PIntToInt32(int(in.Price)),
	}
}

func ToPbExternalCreateOrderLines(ins []*haravan.Item) (outs []*pbexternal.OrderLine) {
	for _, in := range ins {
		outs = append(outs, ToPbExternalCreateOrderLine(in))
	}
	return
}
