// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	identitymodel "o.o/backend/com/main/identity/model"
	conversion "o.o/backend/pkg/common/conversion"
	shopmodel "o.o/backend/zexp/etl/main/shop/model"
	dot "o.o/capi/dot"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shopmodel.ShippingServiceSelectStrategyItem)(nil), (*identitymodel.ShippingServiceSelectStrategyItem)(nil), func(arg, out interface{}) error {
		Convert_shopmodel_ShippingServiceSelectStrategyItem_identitymodel_ShippingServiceSelectStrategyItem(arg.(*shopmodel.ShippingServiceSelectStrategyItem), out.(*identitymodel.ShippingServiceSelectStrategyItem))
		return nil
	})
	s.Register(([]*shopmodel.ShippingServiceSelectStrategyItem)(nil), (*[]*identitymodel.ShippingServiceSelectStrategyItem)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopmodel_ShippingServiceSelectStrategyItems_identitymodel_ShippingServiceSelectStrategyItems(arg.([]*shopmodel.ShippingServiceSelectStrategyItem))
		*out.(*[]*identitymodel.ShippingServiceSelectStrategyItem) = out0
		return nil
	})
	s.Register((*identitymodel.ShippingServiceSelectStrategyItem)(nil), (*shopmodel.ShippingServiceSelectStrategyItem)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_ShippingServiceSelectStrategyItem_shopmodel_ShippingServiceSelectStrategyItem(arg.(*identitymodel.ShippingServiceSelectStrategyItem), out.(*shopmodel.ShippingServiceSelectStrategyItem))
		return nil
	})
	s.Register(([]*identitymodel.ShippingServiceSelectStrategyItem)(nil), (*[]*shopmodel.ShippingServiceSelectStrategyItem)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_ShippingServiceSelectStrategyItems_shopmodel_ShippingServiceSelectStrategyItems(arg.([]*identitymodel.ShippingServiceSelectStrategyItem))
		*out.(*[]*shopmodel.ShippingServiceSelectStrategyItem) = out0
		return nil
	})
	s.Register((*shopmodel.Shop)(nil), (*identitymodel.Shop)(nil), func(arg, out interface{}) error {
		Convert_shopmodel_Shop_identitymodel_Shop(arg.(*shopmodel.Shop), out.(*identitymodel.Shop))
		return nil
	})
	s.Register(([]*shopmodel.Shop)(nil), (*[]*identitymodel.Shop)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopmodel_Shops_identitymodel_Shops(arg.([]*shopmodel.Shop))
		*out.(*[]*identitymodel.Shop) = out0
		return nil
	})
	s.Register((*identitymodel.Shop)(nil), (*shopmodel.Shop)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_Shop_shopmodel_Shop(arg.(*identitymodel.Shop), out.(*shopmodel.Shop))
		return nil
	})
	s.Register(([]*identitymodel.Shop)(nil), (*[]*shopmodel.Shop)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_Shops_shopmodel_Shops(arg.([]*identitymodel.Shop))
		*out.(*[]*shopmodel.Shop) = out0
		return nil
	})
	s.Register((*shopmodel.SurveyInfo)(nil), (*identitymodel.SurveyInfo)(nil), func(arg, out interface{}) error {
		Convert_shopmodel_SurveyInfo_identitymodel_SurveyInfo(arg.(*shopmodel.SurveyInfo), out.(*identitymodel.SurveyInfo))
		return nil
	})
	s.Register(([]*shopmodel.SurveyInfo)(nil), (*[]*identitymodel.SurveyInfo)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopmodel_SurveyInfoes_identitymodel_SurveyInfoes(arg.([]*shopmodel.SurveyInfo))
		*out.(*[]*identitymodel.SurveyInfo) = out0
		return nil
	})
	s.Register((*identitymodel.SurveyInfo)(nil), (*shopmodel.SurveyInfo)(nil), func(arg, out interface{}) error {
		Convert_identitymodel_SurveyInfo_shopmodel_SurveyInfo(arg.(*identitymodel.SurveyInfo), out.(*shopmodel.SurveyInfo))
		return nil
	})
	s.Register(([]*identitymodel.SurveyInfo)(nil), (*[]*shopmodel.SurveyInfo)(nil), func(arg, out interface{}) error {
		out0 := Convert_identitymodel_SurveyInfoes_shopmodel_SurveyInfoes(arg.([]*identitymodel.SurveyInfo))
		*out.(*[]*shopmodel.SurveyInfo) = out0
		return nil
	})
}

//-- convert o.o/backend/com/main/identity/model.ShippingServiceSelectStrategyItem --//

func Convert_shopmodel_ShippingServiceSelectStrategyItem_identitymodel_ShippingServiceSelectStrategyItem(arg *shopmodel.ShippingServiceSelectStrategyItem, out *identitymodel.ShippingServiceSelectStrategyItem) *identitymodel.ShippingServiceSelectStrategyItem {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identitymodel.ShippingServiceSelectStrategyItem{}
	}
	convert_shopmodel_ShippingServiceSelectStrategyItem_identitymodel_ShippingServiceSelectStrategyItem(arg, out)
	return out
}

func convert_shopmodel_ShippingServiceSelectStrategyItem_identitymodel_ShippingServiceSelectStrategyItem(arg *shopmodel.ShippingServiceSelectStrategyItem, out *identitymodel.ShippingServiceSelectStrategyItem) {
	out.Key = arg.Key     // simple assign
	out.Value = arg.Value // simple assign
}

func Convert_shopmodel_ShippingServiceSelectStrategyItems_identitymodel_ShippingServiceSelectStrategyItems(args []*shopmodel.ShippingServiceSelectStrategyItem) (outs []*identitymodel.ShippingServiceSelectStrategyItem) {
	if args == nil {
		return nil
	}
	tmps := make([]identitymodel.ShippingServiceSelectStrategyItem, len(args))
	outs = make([]*identitymodel.ShippingServiceSelectStrategyItem, len(args))
	for i := range tmps {
		outs[i] = Convert_shopmodel_ShippingServiceSelectStrategyItem_identitymodel_ShippingServiceSelectStrategyItem(args[i], &tmps[i])
	}
	return outs
}

func Convert_identitymodel_ShippingServiceSelectStrategyItem_shopmodel_ShippingServiceSelectStrategyItem(arg *identitymodel.ShippingServiceSelectStrategyItem, out *shopmodel.ShippingServiceSelectStrategyItem) *shopmodel.ShippingServiceSelectStrategyItem {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopmodel.ShippingServiceSelectStrategyItem{}
	}
	convert_identitymodel_ShippingServiceSelectStrategyItem_shopmodel_ShippingServiceSelectStrategyItem(arg, out)
	return out
}

func convert_identitymodel_ShippingServiceSelectStrategyItem_shopmodel_ShippingServiceSelectStrategyItem(arg *identitymodel.ShippingServiceSelectStrategyItem, out *shopmodel.ShippingServiceSelectStrategyItem) {
	out.Key = arg.Key     // simple assign
	out.Value = arg.Value // simple assign
}

func Convert_identitymodel_ShippingServiceSelectStrategyItems_shopmodel_ShippingServiceSelectStrategyItems(args []*identitymodel.ShippingServiceSelectStrategyItem) (outs []*shopmodel.ShippingServiceSelectStrategyItem) {
	if args == nil {
		return nil
	}
	tmps := make([]shopmodel.ShippingServiceSelectStrategyItem, len(args))
	outs = make([]*shopmodel.ShippingServiceSelectStrategyItem, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_ShippingServiceSelectStrategyItem_shopmodel_ShippingServiceSelectStrategyItem(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/backend/com/main/identity/model.Shop --//

func Convert_shopmodel_Shop_identitymodel_Shop(arg *shopmodel.Shop, out *identitymodel.Shop) *identitymodel.Shop {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identitymodel.Shop{}
	}
	convert_shopmodel_Shop_identitymodel_Shop(arg, out)
	return out
}

func convert_shopmodel_Shop_identitymodel_Shop(arg *shopmodel.Shop, out *identitymodel.Shop) {
	out.ID = arg.ID                                       // simple assign
	out.OwnerID = arg.OwnerID                             // simple assign
	out.IsTest = 0                                        // zero value
	out.Name = arg.Name                                   // simple assign
	out.NameNorm = ""                                     // zero value
	out.AddressID = arg.AddressID                         // simple assign
	out.ShipToAddressID = arg.ShipToAddressID             // simple assign
	out.ShipFromAddressID = arg.ShipFromAddressID         // simple assign
	out.Phone = arg.Phone                                 // simple assign
	out.BankAccount = arg.BankAccount                     // simple assign
	out.WebsiteURL = dot.NullString{}                     // types do not match
	out.ImageURL = arg.ImageURL                           // simple assign
	out.Email = arg.Email                                 // simple assign
	out.Code = arg.Code                                   // simple assign
	out.AutoCreateFFM = arg.AutoCreateFFM                 // simple assign
	out.OrderSourceID = 0                                 // zero value
	out.Status = arg.Status                               // simple assign
	out.CreatedAt = arg.CreatedAt                         // simple assign
	out.UpdatedAt = arg.UpdatedAt                         // simple assign
	out.DeletedAt = time.Time{}                           // zero value
	out.Address = arg.Address                             // simple assign
	out.RecognizedHosts = arg.RecognizedHosts             // simple assign
	out.GhnNoteCode = arg.GhnNoteCode                     // simple assign
	out.TryOn = arg.TryOn                                 // simple assign
	out.CompanyInfo = arg.CompanyInfo                     // simple assign
	out.MoneyTransactionRRule = arg.MoneyTransactionRRule // simple assign
	out.SurveyInfo = Convert_shopmodel_SurveyInfoes_identitymodel_SurveyInfoes(arg.SurveyInfo)
	out.ShippingServiceSelectStrategy = nil // zero value
	out.InventoryOverstock = dot.NullBool{} // zero value
	out.WLPartnerID = 0                     // zero value
	out.Rid = arg.Rid                       // simple assign
}

func Convert_shopmodel_Shops_identitymodel_Shops(args []*shopmodel.Shop) (outs []*identitymodel.Shop) {
	if args == nil {
		return nil
	}
	tmps := make([]identitymodel.Shop, len(args))
	outs = make([]*identitymodel.Shop, len(args))
	for i := range tmps {
		outs[i] = Convert_shopmodel_Shop_identitymodel_Shop(args[i], &tmps[i])
	}
	return outs
}

func Convert_identitymodel_Shop_shopmodel_Shop(arg *identitymodel.Shop, out *shopmodel.Shop) *shopmodel.Shop {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopmodel.Shop{}
	}
	convert_identitymodel_Shop_shopmodel_Shop(arg, out)
	return out
}

func convert_identitymodel_Shop_shopmodel_Shop(arg *identitymodel.Shop, out *shopmodel.Shop) {
	out.ID = arg.ID                                       // simple assign
	out.Name = arg.Name                                   // simple assign
	out.OwnerID = arg.OwnerID                             // simple assign
	out.AddressID = arg.AddressID                         // simple assign
	out.ShipToAddressID = arg.ShipToAddressID             // simple assign
	out.ShipFromAddressID = arg.ShipFromAddressID         // simple assign
	out.Phone = arg.Phone                                 // simple assign
	out.BankAccount = arg.BankAccount                     // simple assign
	out.WebsiteURL = ""                                   // types do not match
	out.ImageURL = arg.ImageURL                           // simple assign
	out.Email = arg.Email                                 // simple assign
	out.Code = arg.Code                                   // simple assign
	out.AutoCreateFFM = arg.AutoCreateFFM                 // simple assign
	out.Status = arg.Status                               // simple assign
	out.CreatedAt = arg.CreatedAt                         // simple assign
	out.UpdatedAt = arg.UpdatedAt                         // simple assign
	out.Address = arg.Address                             // simple assign
	out.RecognizedHosts = arg.RecognizedHosts             // simple assign
	out.GhnNoteCode = arg.GhnNoteCode                     // simple assign
	out.TryOn = arg.TryOn                                 // simple assign
	out.CompanyInfo = arg.CompanyInfo                     // simple assign
	out.MoneyTransactionRRule = arg.MoneyTransactionRRule // simple assign
	out.SurveyInfo = Convert_identitymodel_SurveyInfoes_shopmodel_SurveyInfoes(arg.SurveyInfo)
	out.Rid = arg.Rid // simple assign
}

func Convert_identitymodel_Shops_shopmodel_Shops(args []*identitymodel.Shop) (outs []*shopmodel.Shop) {
	if args == nil {
		return nil
	}
	tmps := make([]shopmodel.Shop, len(args))
	outs = make([]*shopmodel.Shop, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_Shop_shopmodel_Shop(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/backend/com/main/identity/model.SurveyInfo --//

func Convert_shopmodel_SurveyInfo_identitymodel_SurveyInfo(arg *shopmodel.SurveyInfo, out *identitymodel.SurveyInfo) *identitymodel.SurveyInfo {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &identitymodel.SurveyInfo{}
	}
	convert_shopmodel_SurveyInfo_identitymodel_SurveyInfo(arg, out)
	return out
}

func convert_shopmodel_SurveyInfo_identitymodel_SurveyInfo(arg *shopmodel.SurveyInfo, out *identitymodel.SurveyInfo) {
	out.Key = arg.Key           // simple assign
	out.Question = arg.Question // simple assign
	out.Answer = arg.Answer     // simple assign
}

func Convert_shopmodel_SurveyInfoes_identitymodel_SurveyInfoes(args []*shopmodel.SurveyInfo) (outs []*identitymodel.SurveyInfo) {
	if args == nil {
		return nil
	}
	tmps := make([]identitymodel.SurveyInfo, len(args))
	outs = make([]*identitymodel.SurveyInfo, len(args))
	for i := range tmps {
		outs[i] = Convert_shopmodel_SurveyInfo_identitymodel_SurveyInfo(args[i], &tmps[i])
	}
	return outs
}

func Convert_identitymodel_SurveyInfo_shopmodel_SurveyInfo(arg *identitymodel.SurveyInfo, out *shopmodel.SurveyInfo) *shopmodel.SurveyInfo {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopmodel.SurveyInfo{}
	}
	convert_identitymodel_SurveyInfo_shopmodel_SurveyInfo(arg, out)
	return out
}

func convert_identitymodel_SurveyInfo_shopmodel_SurveyInfo(arg *identitymodel.SurveyInfo, out *shopmodel.SurveyInfo) {
	out.Key = arg.Key           // simple assign
	out.Question = arg.Question // simple assign
	out.Answer = arg.Answer     // simple assign
}

func Convert_identitymodel_SurveyInfoes_shopmodel_SurveyInfoes(args []*identitymodel.SurveyInfo) (outs []*shopmodel.SurveyInfo) {
	if args == nil {
		return nil
	}
	tmps := make([]shopmodel.SurveyInfo, len(args))
	outs = make([]*shopmodel.SurveyInfo, len(args))
	for i := range tmps {
		outs[i] = Convert_identitymodel_SurveyInfo_shopmodel_SurveyInfo(args[i], &tmps[i])
	}
	return outs
}
