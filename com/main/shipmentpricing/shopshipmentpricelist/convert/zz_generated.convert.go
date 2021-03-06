// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	shopshipmentpricelist "o.o/api/main/shipmentpricing/shopshipmentpricelist"
	shopshipmentpricelistmodel "o.o/backend/com/main/shipmentpricing/shopshipmentpricelist/model"
	conversion "o.o/backend/pkg/common/conversion"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shopshipmentpricelistmodel.ShopShipmentPriceList)(nil), (*shopshipmentpricelist.ShopShipmentPriceList)(nil), func(arg, out interface{}) error {
		Convert_shopshipmentpricelistmodel_ShopShipmentPriceList_shopshipmentpricelist_ShopShipmentPriceList(arg.(*shopshipmentpricelistmodel.ShopShipmentPriceList), out.(*shopshipmentpricelist.ShopShipmentPriceList))
		return nil
	})
	s.Register(([]*shopshipmentpricelistmodel.ShopShipmentPriceList)(nil), (*[]*shopshipmentpricelist.ShopShipmentPriceList)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopshipmentpricelistmodel_ShopShipmentPriceLists_shopshipmentpricelist_ShopShipmentPriceLists(arg.([]*shopshipmentpricelistmodel.ShopShipmentPriceList))
		*out.(*[]*shopshipmentpricelist.ShopShipmentPriceList) = out0
		return nil
	})
	s.Register((*shopshipmentpricelist.ShopShipmentPriceList)(nil), (*shopshipmentpricelistmodel.ShopShipmentPriceList)(nil), func(arg, out interface{}) error {
		Convert_shopshipmentpricelist_ShopShipmentPriceList_shopshipmentpricelistmodel_ShopShipmentPriceList(arg.(*shopshipmentpricelist.ShopShipmentPriceList), out.(*shopshipmentpricelistmodel.ShopShipmentPriceList))
		return nil
	})
	s.Register(([]*shopshipmentpricelist.ShopShipmentPriceList)(nil), (*[]*shopshipmentpricelistmodel.ShopShipmentPriceList)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopshipmentpricelist_ShopShipmentPriceLists_shopshipmentpricelistmodel_ShopShipmentPriceLists(arg.([]*shopshipmentpricelist.ShopShipmentPriceList))
		*out.(*[]*shopshipmentpricelistmodel.ShopShipmentPriceList) = out0
		return nil
	})
	s.Register((*shopshipmentpricelist.CreateShopShipmentPriceListArgs)(nil), (*shopshipmentpricelist.ShopShipmentPriceList)(nil), func(arg, out interface{}) error {
		Apply_shopshipmentpricelist_CreateShopShipmentPriceListArgs_shopshipmentpricelist_ShopShipmentPriceList(arg.(*shopshipmentpricelist.CreateShopShipmentPriceListArgs), out.(*shopshipmentpricelist.ShopShipmentPriceList))
		return nil
	})
	s.Register((*shopshipmentpricelist.UpdateShopShipmentPriceListArgs)(nil), (*shopshipmentpricelist.ShopShipmentPriceList)(nil), func(arg, out interface{}) error {
		Apply_shopshipmentpricelist_UpdateShopShipmentPriceListArgs_shopshipmentpricelist_ShopShipmentPriceList(arg.(*shopshipmentpricelist.UpdateShopShipmentPriceListArgs), out.(*shopshipmentpricelist.ShopShipmentPriceList))
		return nil
	})
}

//-- convert o.o/api/main/shipmentpricing/shopshipmentpricelist.ShopShipmentPriceList --//

func Convert_shopshipmentpricelistmodel_ShopShipmentPriceList_shopshipmentpricelist_ShopShipmentPriceList(arg *shopshipmentpricelistmodel.ShopShipmentPriceList, out *shopshipmentpricelist.ShopShipmentPriceList) *shopshipmentpricelist.ShopShipmentPriceList {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopshipmentpricelist.ShopShipmentPriceList{}
	}
	convert_shopshipmentpricelistmodel_ShopShipmentPriceList_shopshipmentpricelist_ShopShipmentPriceList(arg, out)
	return out
}

func convert_shopshipmentpricelistmodel_ShopShipmentPriceList_shopshipmentpricelist_ShopShipmentPriceList(arg *shopshipmentpricelistmodel.ShopShipmentPriceList, out *shopshipmentpricelist.ShopShipmentPriceList) {
	out.ShopID = arg.ShopID                           // simple assign
	out.ShipmentPriceListID = arg.ShipmentPriceListID // simple assign
	out.ConnectionID = arg.ConnectionID               // simple assign
	out.Note = arg.Note                               // simple assign
	out.CreatedAt = arg.CreatedAt                     // simple assign
	out.UpdatedAt = arg.UpdatedAt                     // simple assign
	out.DeletedAt = arg.DeletedAt                     // simple assign
	out.UpdatedBy = arg.UpdatedBy                     // simple assign
}

func Convert_shopshipmentpricelistmodel_ShopShipmentPriceLists_shopshipmentpricelist_ShopShipmentPriceLists(args []*shopshipmentpricelistmodel.ShopShipmentPriceList) (outs []*shopshipmentpricelist.ShopShipmentPriceList) {
	if args == nil {
		return nil
	}
	tmps := make([]shopshipmentpricelist.ShopShipmentPriceList, len(args))
	outs = make([]*shopshipmentpricelist.ShopShipmentPriceList, len(args))
	for i := range tmps {
		outs[i] = Convert_shopshipmentpricelistmodel_ShopShipmentPriceList_shopshipmentpricelist_ShopShipmentPriceList(args[i], &tmps[i])
	}
	return outs
}

func Convert_shopshipmentpricelist_ShopShipmentPriceList_shopshipmentpricelistmodel_ShopShipmentPriceList(arg *shopshipmentpricelist.ShopShipmentPriceList, out *shopshipmentpricelistmodel.ShopShipmentPriceList) *shopshipmentpricelistmodel.ShopShipmentPriceList {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopshipmentpricelistmodel.ShopShipmentPriceList{}
	}
	convert_shopshipmentpricelist_ShopShipmentPriceList_shopshipmentpricelistmodel_ShopShipmentPriceList(arg, out)
	return out
}

func convert_shopshipmentpricelist_ShopShipmentPriceList_shopshipmentpricelistmodel_ShopShipmentPriceList(arg *shopshipmentpricelist.ShopShipmentPriceList, out *shopshipmentpricelistmodel.ShopShipmentPriceList) {
	out.ShopID = arg.ShopID                           // simple assign
	out.ShipmentPriceListID = arg.ShipmentPriceListID // simple assign
	out.ConnectionID = arg.ConnectionID               // simple assign
	out.Note = arg.Note                               // simple assign
	out.CreatedAt = arg.CreatedAt                     // simple assign
	out.UpdatedAt = arg.UpdatedAt                     // simple assign
	out.DeletedAt = arg.DeletedAt                     // simple assign
	out.UpdatedBy = arg.UpdatedBy                     // simple assign
	out.WLPartnerID = 0                               // zero value
}

func Convert_shopshipmentpricelist_ShopShipmentPriceLists_shopshipmentpricelistmodel_ShopShipmentPriceLists(args []*shopshipmentpricelist.ShopShipmentPriceList) (outs []*shopshipmentpricelistmodel.ShopShipmentPriceList) {
	if args == nil {
		return nil
	}
	tmps := make([]shopshipmentpricelistmodel.ShopShipmentPriceList, len(args))
	outs = make([]*shopshipmentpricelistmodel.ShopShipmentPriceList, len(args))
	for i := range tmps {
		outs[i] = Convert_shopshipmentpricelist_ShopShipmentPriceList_shopshipmentpricelistmodel_ShopShipmentPriceList(args[i], &tmps[i])
	}
	return outs
}

func Apply_shopshipmentpricelist_CreateShopShipmentPriceListArgs_shopshipmentpricelist_ShopShipmentPriceList(arg *shopshipmentpricelist.CreateShopShipmentPriceListArgs, out *shopshipmentpricelist.ShopShipmentPriceList) *shopshipmentpricelist.ShopShipmentPriceList {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopshipmentpricelist.ShopShipmentPriceList{}
	}
	apply_shopshipmentpricelist_CreateShopShipmentPriceListArgs_shopshipmentpricelist_ShopShipmentPriceList(arg, out)
	return out
}

func apply_shopshipmentpricelist_CreateShopShipmentPriceListArgs_shopshipmentpricelist_ShopShipmentPriceList(arg *shopshipmentpricelist.CreateShopShipmentPriceListArgs, out *shopshipmentpricelist.ShopShipmentPriceList) {
	out.ShopID = arg.ShopID                           // simple assign
	out.ShipmentPriceListID = arg.ShipmentPriceListID // simple assign
	out.ConnectionID = arg.ConnectionID               // simple assign
	out.Note = arg.Note                               // simple assign
	out.CreatedAt = time.Time{}                       // zero value
	out.UpdatedAt = time.Time{}                       // zero value
	out.DeletedAt = time.Time{}                       // zero value
	out.UpdatedBy = arg.UpdatedBy                     // simple assign
}

func Apply_shopshipmentpricelist_UpdateShopShipmentPriceListArgs_shopshipmentpricelist_ShopShipmentPriceList(arg *shopshipmentpricelist.UpdateShopShipmentPriceListArgs, out *shopshipmentpricelist.ShopShipmentPriceList) *shopshipmentpricelist.ShopShipmentPriceList {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopshipmentpricelist.ShopShipmentPriceList{}
	}
	apply_shopshipmentpricelist_UpdateShopShipmentPriceListArgs_shopshipmentpricelist_ShopShipmentPriceList(arg, out)
	return out
}

func apply_shopshipmentpricelist_UpdateShopShipmentPriceListArgs_shopshipmentpricelist_ShopShipmentPriceList(arg *shopshipmentpricelist.UpdateShopShipmentPriceListArgs, out *shopshipmentpricelist.ShopShipmentPriceList) {
	out.ShopID = arg.ShopID                           // simple assign
	out.ShipmentPriceListID = arg.ShipmentPriceListID // simple assign
	out.ConnectionID = arg.ConnectionID               // simple assign
	out.Note = arg.Note                               // simple assign
	out.CreatedAt = out.CreatedAt                     // no change
	out.UpdatedAt = out.UpdatedAt                     // no change
	out.DeletedAt = out.DeletedAt                     // no change
	out.UpdatedBy = arg.UpdatedBy                     // simple assign
}
