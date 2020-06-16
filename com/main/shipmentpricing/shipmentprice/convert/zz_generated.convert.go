// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	shipmentprice "o.o/api/main/shipmentpricing/shipmentprice"
	shipmentpricemodel "o.o/backend/com/main/shipmentpricing/shipmentprice/model"
	conversion "o.o/backend/pkg/common/conversion"
)

/*
Custom conversions:
    Convert_shipmentprice_ShippingFee_To_shippingsharemodel_ShippingFeeLine    // not use, no conversions between params

Ignored functions:
    Convert_shipmentprice_ShippingFees_To_shippingsharemodel_ShippingFeeLines    // params are not pointer to named types
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shipmentpricemodel.AdditionalFee)(nil), (*shipmentprice.AdditionalFee)(nil), func(arg, out interface{}) error {
		Convert_shipmentpricemodel_AdditionalFee_shipmentprice_AdditionalFee(arg.(*shipmentpricemodel.AdditionalFee), out.(*shipmentprice.AdditionalFee))
		return nil
	})
	s.Register(([]*shipmentpricemodel.AdditionalFee)(nil), (*[]*shipmentprice.AdditionalFee)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentpricemodel_AdditionalFees_shipmentprice_AdditionalFees(arg.([]*shipmentpricemodel.AdditionalFee))
		*out.(*[]*shipmentprice.AdditionalFee) = out0
		return nil
	})
	s.Register((*shipmentprice.AdditionalFee)(nil), (*shipmentpricemodel.AdditionalFee)(nil), func(arg, out interface{}) error {
		Convert_shipmentprice_AdditionalFee_shipmentpricemodel_AdditionalFee(arg.(*shipmentprice.AdditionalFee), out.(*shipmentpricemodel.AdditionalFee))
		return nil
	})
	s.Register(([]*shipmentprice.AdditionalFee)(nil), (*[]*shipmentpricemodel.AdditionalFee)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentprice_AdditionalFees_shipmentpricemodel_AdditionalFees(arg.([]*shipmentprice.AdditionalFee))
		*out.(*[]*shipmentpricemodel.AdditionalFee) = out0
		return nil
	})
	s.Register((*shipmentpricemodel.AdditionalFeeRule)(nil), (*shipmentprice.AdditionalFeeRule)(nil), func(arg, out interface{}) error {
		Convert_shipmentpricemodel_AdditionalFeeRule_shipmentprice_AdditionalFeeRule(arg.(*shipmentpricemodel.AdditionalFeeRule), out.(*shipmentprice.AdditionalFeeRule))
		return nil
	})
	s.Register(([]*shipmentpricemodel.AdditionalFeeRule)(nil), (*[]*shipmentprice.AdditionalFeeRule)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentpricemodel_AdditionalFeeRules_shipmentprice_AdditionalFeeRules(arg.([]*shipmentpricemodel.AdditionalFeeRule))
		*out.(*[]*shipmentprice.AdditionalFeeRule) = out0
		return nil
	})
	s.Register((*shipmentprice.AdditionalFeeRule)(nil), (*shipmentpricemodel.AdditionalFeeRule)(nil), func(arg, out interface{}) error {
		Convert_shipmentprice_AdditionalFeeRule_shipmentpricemodel_AdditionalFeeRule(arg.(*shipmentprice.AdditionalFeeRule), out.(*shipmentpricemodel.AdditionalFeeRule))
		return nil
	})
	s.Register(([]*shipmentprice.AdditionalFeeRule)(nil), (*[]*shipmentpricemodel.AdditionalFeeRule)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentprice_AdditionalFeeRules_shipmentpricemodel_AdditionalFeeRules(arg.([]*shipmentprice.AdditionalFeeRule))
		*out.(*[]*shipmentpricemodel.AdditionalFeeRule) = out0
		return nil
	})
	s.Register((*shipmentpricemodel.PricingDetail)(nil), (*shipmentprice.PricingDetail)(nil), func(arg, out interface{}) error {
		Convert_shipmentpricemodel_PricingDetail_shipmentprice_PricingDetail(arg.(*shipmentpricemodel.PricingDetail), out.(*shipmentprice.PricingDetail))
		return nil
	})
	s.Register(([]*shipmentpricemodel.PricingDetail)(nil), (*[]*shipmentprice.PricingDetail)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentpricemodel_PricingDetails_shipmentprice_PricingDetails(arg.([]*shipmentpricemodel.PricingDetail))
		*out.(*[]*shipmentprice.PricingDetail) = out0
		return nil
	})
	s.Register((*shipmentprice.PricingDetail)(nil), (*shipmentpricemodel.PricingDetail)(nil), func(arg, out interface{}) error {
		Convert_shipmentprice_PricingDetail_shipmentpricemodel_PricingDetail(arg.(*shipmentprice.PricingDetail), out.(*shipmentpricemodel.PricingDetail))
		return nil
	})
	s.Register(([]*shipmentprice.PricingDetail)(nil), (*[]*shipmentpricemodel.PricingDetail)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentprice_PricingDetails_shipmentpricemodel_PricingDetails(arg.([]*shipmentprice.PricingDetail))
		*out.(*[]*shipmentpricemodel.PricingDetail) = out0
		return nil
	})
	s.Register((*shipmentpricemodel.PricingDetailOverweight)(nil), (*shipmentprice.PricingDetailOverweight)(nil), func(arg, out interface{}) error {
		Convert_shipmentpricemodel_PricingDetailOverweight_shipmentprice_PricingDetailOverweight(arg.(*shipmentpricemodel.PricingDetailOverweight), out.(*shipmentprice.PricingDetailOverweight))
		return nil
	})
	s.Register(([]*shipmentpricemodel.PricingDetailOverweight)(nil), (*[]*shipmentprice.PricingDetailOverweight)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentpricemodel_PricingDetailOverweights_shipmentprice_PricingDetailOverweights(arg.([]*shipmentpricemodel.PricingDetailOverweight))
		*out.(*[]*shipmentprice.PricingDetailOverweight) = out0
		return nil
	})
	s.Register((*shipmentprice.PricingDetailOverweight)(nil), (*shipmentpricemodel.PricingDetailOverweight)(nil), func(arg, out interface{}) error {
		Convert_shipmentprice_PricingDetailOverweight_shipmentpricemodel_PricingDetailOverweight(arg.(*shipmentprice.PricingDetailOverweight), out.(*shipmentpricemodel.PricingDetailOverweight))
		return nil
	})
	s.Register(([]*shipmentprice.PricingDetailOverweight)(nil), (*[]*shipmentpricemodel.PricingDetailOverweight)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentprice_PricingDetailOverweights_shipmentpricemodel_PricingDetailOverweights(arg.([]*shipmentprice.PricingDetailOverweight))
		*out.(*[]*shipmentpricemodel.PricingDetailOverweight) = out0
		return nil
	})
	s.Register((*shipmentpricemodel.ShipmentPrice)(nil), (*shipmentprice.ShipmentPrice)(nil), func(arg, out interface{}) error {
		Convert_shipmentpricemodel_ShipmentPrice_shipmentprice_ShipmentPrice(arg.(*shipmentpricemodel.ShipmentPrice), out.(*shipmentprice.ShipmentPrice))
		return nil
	})
	s.Register(([]*shipmentpricemodel.ShipmentPrice)(nil), (*[]*shipmentprice.ShipmentPrice)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentpricemodel_ShipmentPrices_shipmentprice_ShipmentPrices(arg.([]*shipmentpricemodel.ShipmentPrice))
		*out.(*[]*shipmentprice.ShipmentPrice) = out0
		return nil
	})
	s.Register((*shipmentprice.ShipmentPrice)(nil), (*shipmentpricemodel.ShipmentPrice)(nil), func(arg, out interface{}) error {
		Convert_shipmentprice_ShipmentPrice_shipmentpricemodel_ShipmentPrice(arg.(*shipmentprice.ShipmentPrice), out.(*shipmentpricemodel.ShipmentPrice))
		return nil
	})
	s.Register(([]*shipmentprice.ShipmentPrice)(nil), (*[]*shipmentpricemodel.ShipmentPrice)(nil), func(arg, out interface{}) error {
		out0 := Convert_shipmentprice_ShipmentPrices_shipmentpricemodel_ShipmentPrices(arg.([]*shipmentprice.ShipmentPrice))
		*out.(*[]*shipmentpricemodel.ShipmentPrice) = out0
		return nil
	})
	s.Register((*shipmentprice.CreateShipmentPriceArgs)(nil), (*shipmentprice.ShipmentPrice)(nil), func(arg, out interface{}) error {
		Apply_shipmentprice_CreateShipmentPriceArgs_shipmentprice_ShipmentPrice(arg.(*shipmentprice.CreateShipmentPriceArgs), out.(*shipmentprice.ShipmentPrice))
		return nil
	})
	s.Register((*shipmentprice.UpdateShipmentPriceArgs)(nil), (*shipmentprice.ShipmentPrice)(nil), func(arg, out interface{}) error {
		Apply_shipmentprice_UpdateShipmentPriceArgs_shipmentprice_ShipmentPrice(arg.(*shipmentprice.UpdateShipmentPriceArgs), out.(*shipmentprice.ShipmentPrice))
		return nil
	})
}

//-- convert o.o/api/main/shipmentpricing/shipmentprice.AdditionalFee --//

func Convert_shipmentpricemodel_AdditionalFee_shipmentprice_AdditionalFee(arg *shipmentpricemodel.AdditionalFee, out *shipmentprice.AdditionalFee) *shipmentprice.AdditionalFee {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentprice.AdditionalFee{}
	}
	convert_shipmentpricemodel_AdditionalFee_shipmentprice_AdditionalFee(arg, out)
	return out
}

func convert_shipmentpricemodel_AdditionalFee_shipmentprice_AdditionalFee(arg *shipmentpricemodel.AdditionalFee, out *shipmentprice.AdditionalFee) {
	out.FeeType = arg.FeeType // simple assign
	out.Rules = Convert_shipmentpricemodel_AdditionalFeeRules_shipmentprice_AdditionalFeeRules(arg.Rules)
}

func Convert_shipmentpricemodel_AdditionalFees_shipmentprice_AdditionalFees(args []*shipmentpricemodel.AdditionalFee) (outs []*shipmentprice.AdditionalFee) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentprice.AdditionalFee, len(args))
	outs = make([]*shipmentprice.AdditionalFee, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentpricemodel_AdditionalFee_shipmentprice_AdditionalFee(args[i], &tmps[i])
	}
	return outs
}

func Convert_shipmentprice_AdditionalFee_shipmentpricemodel_AdditionalFee(arg *shipmentprice.AdditionalFee, out *shipmentpricemodel.AdditionalFee) *shipmentpricemodel.AdditionalFee {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentpricemodel.AdditionalFee{}
	}
	convert_shipmentprice_AdditionalFee_shipmentpricemodel_AdditionalFee(arg, out)
	return out
}

func convert_shipmentprice_AdditionalFee_shipmentpricemodel_AdditionalFee(arg *shipmentprice.AdditionalFee, out *shipmentpricemodel.AdditionalFee) {
	out.FeeType = arg.FeeType // simple assign
	out.Rules = Convert_shipmentprice_AdditionalFeeRules_shipmentpricemodel_AdditionalFeeRules(arg.Rules)
}

func Convert_shipmentprice_AdditionalFees_shipmentpricemodel_AdditionalFees(args []*shipmentprice.AdditionalFee) (outs []*shipmentpricemodel.AdditionalFee) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentpricemodel.AdditionalFee, len(args))
	outs = make([]*shipmentpricemodel.AdditionalFee, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentprice_AdditionalFee_shipmentpricemodel_AdditionalFee(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/main/shipmentpricing/shipmentprice.AdditionalFeeRule --//

func Convert_shipmentpricemodel_AdditionalFeeRule_shipmentprice_AdditionalFeeRule(arg *shipmentpricemodel.AdditionalFeeRule, out *shipmentprice.AdditionalFeeRule) *shipmentprice.AdditionalFeeRule {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentprice.AdditionalFeeRule{}
	}
	convert_shipmentpricemodel_AdditionalFeeRule_shipmentprice_AdditionalFeeRule(arg, out)
	return out
}

func convert_shipmentpricemodel_AdditionalFeeRule_shipmentprice_AdditionalFeeRule(arg *shipmentpricemodel.AdditionalFeeRule, out *shipmentprice.AdditionalFeeRule) {
	out.MinValue = arg.MinValue                   // simple assign
	out.MaxValue = arg.MaxValue                   // simple assign
	out.PriceModifierType = arg.PriceModifierType // simple assign
	out.Amount = arg.Amount                       // simple assign
	out.MinPrice = arg.MinPrice                   // simple assign
}

func Convert_shipmentpricemodel_AdditionalFeeRules_shipmentprice_AdditionalFeeRules(args []*shipmentpricemodel.AdditionalFeeRule) (outs []*shipmentprice.AdditionalFeeRule) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentprice.AdditionalFeeRule, len(args))
	outs = make([]*shipmentprice.AdditionalFeeRule, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentpricemodel_AdditionalFeeRule_shipmentprice_AdditionalFeeRule(args[i], &tmps[i])
	}
	return outs
}

func Convert_shipmentprice_AdditionalFeeRule_shipmentpricemodel_AdditionalFeeRule(arg *shipmentprice.AdditionalFeeRule, out *shipmentpricemodel.AdditionalFeeRule) *shipmentpricemodel.AdditionalFeeRule {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentpricemodel.AdditionalFeeRule{}
	}
	convert_shipmentprice_AdditionalFeeRule_shipmentpricemodel_AdditionalFeeRule(arg, out)
	return out
}

func convert_shipmentprice_AdditionalFeeRule_shipmentpricemodel_AdditionalFeeRule(arg *shipmentprice.AdditionalFeeRule, out *shipmentpricemodel.AdditionalFeeRule) {
	out.MinValue = arg.MinValue                   // simple assign
	out.MaxValue = arg.MaxValue                   // simple assign
	out.PriceModifierType = arg.PriceModifierType // simple assign
	out.Amount = arg.Amount                       // simple assign
	out.MinPrice = arg.MinPrice                   // simple assign
}

func Convert_shipmentprice_AdditionalFeeRules_shipmentpricemodel_AdditionalFeeRules(args []*shipmentprice.AdditionalFeeRule) (outs []*shipmentpricemodel.AdditionalFeeRule) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentpricemodel.AdditionalFeeRule, len(args))
	outs = make([]*shipmentpricemodel.AdditionalFeeRule, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentprice_AdditionalFeeRule_shipmentpricemodel_AdditionalFeeRule(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/main/shipmentpricing/shipmentprice.PricingDetail --//

func Convert_shipmentpricemodel_PricingDetail_shipmentprice_PricingDetail(arg *shipmentpricemodel.PricingDetail, out *shipmentprice.PricingDetail) *shipmentprice.PricingDetail {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentprice.PricingDetail{}
	}
	convert_shipmentpricemodel_PricingDetail_shipmentprice_PricingDetail(arg, out)
	return out
}

func convert_shipmentpricemodel_PricingDetail_shipmentprice_PricingDetail(arg *shipmentpricemodel.PricingDetail, out *shipmentprice.PricingDetail) {
	out.Weight = arg.Weight // simple assign
	out.Price = arg.Price   // simple assign
	out.Overweight = Convert_shipmentpricemodel_PricingDetailOverweights_shipmentprice_PricingDetailOverweights(arg.Overweight)
}

func Convert_shipmentpricemodel_PricingDetails_shipmentprice_PricingDetails(args []*shipmentpricemodel.PricingDetail) (outs []*shipmentprice.PricingDetail) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentprice.PricingDetail, len(args))
	outs = make([]*shipmentprice.PricingDetail, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentpricemodel_PricingDetail_shipmentprice_PricingDetail(args[i], &tmps[i])
	}
	return outs
}

func Convert_shipmentprice_PricingDetail_shipmentpricemodel_PricingDetail(arg *shipmentprice.PricingDetail, out *shipmentpricemodel.PricingDetail) *shipmentpricemodel.PricingDetail {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentpricemodel.PricingDetail{}
	}
	convert_shipmentprice_PricingDetail_shipmentpricemodel_PricingDetail(arg, out)
	return out
}

func convert_shipmentprice_PricingDetail_shipmentpricemodel_PricingDetail(arg *shipmentprice.PricingDetail, out *shipmentpricemodel.PricingDetail) {
	out.Weight = arg.Weight // simple assign
	out.Price = arg.Price   // simple assign
	out.Overweight = Convert_shipmentprice_PricingDetailOverweights_shipmentpricemodel_PricingDetailOverweights(arg.Overweight)
}

func Convert_shipmentprice_PricingDetails_shipmentpricemodel_PricingDetails(args []*shipmentprice.PricingDetail) (outs []*shipmentpricemodel.PricingDetail) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentpricemodel.PricingDetail, len(args))
	outs = make([]*shipmentpricemodel.PricingDetail, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentprice_PricingDetail_shipmentpricemodel_PricingDetail(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/main/shipmentpricing/shipmentprice.PricingDetailOverweight --//

func Convert_shipmentpricemodel_PricingDetailOverweight_shipmentprice_PricingDetailOverweight(arg *shipmentpricemodel.PricingDetailOverweight, out *shipmentprice.PricingDetailOverweight) *shipmentprice.PricingDetailOverweight {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentprice.PricingDetailOverweight{}
	}
	convert_shipmentpricemodel_PricingDetailOverweight_shipmentprice_PricingDetailOverweight(arg, out)
	return out
}

func convert_shipmentpricemodel_PricingDetailOverweight_shipmentprice_PricingDetailOverweight(arg *shipmentpricemodel.PricingDetailOverweight, out *shipmentprice.PricingDetailOverweight) {
	out.MinWeight = arg.MinWeight   // simple assign
	out.MaxWeight = arg.MaxWeight   // simple assign
	out.WeightStep = arg.WeightStep // simple assign
	out.PriceStep = arg.PriceStep   // simple assign
}

func Convert_shipmentpricemodel_PricingDetailOverweights_shipmentprice_PricingDetailOverweights(args []*shipmentpricemodel.PricingDetailOverweight) (outs []*shipmentprice.PricingDetailOverweight) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentprice.PricingDetailOverweight, len(args))
	outs = make([]*shipmentprice.PricingDetailOverweight, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentpricemodel_PricingDetailOverweight_shipmentprice_PricingDetailOverweight(args[i], &tmps[i])
	}
	return outs
}

func Convert_shipmentprice_PricingDetailOverweight_shipmentpricemodel_PricingDetailOverweight(arg *shipmentprice.PricingDetailOverweight, out *shipmentpricemodel.PricingDetailOverweight) *shipmentpricemodel.PricingDetailOverweight {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentpricemodel.PricingDetailOverweight{}
	}
	convert_shipmentprice_PricingDetailOverweight_shipmentpricemodel_PricingDetailOverweight(arg, out)
	return out
}

func convert_shipmentprice_PricingDetailOverweight_shipmentpricemodel_PricingDetailOverweight(arg *shipmentprice.PricingDetailOverweight, out *shipmentpricemodel.PricingDetailOverweight) {
	out.MinWeight = arg.MinWeight   // simple assign
	out.MaxWeight = arg.MaxWeight   // simple assign
	out.WeightStep = arg.WeightStep // simple assign
	out.PriceStep = arg.PriceStep   // simple assign
}

func Convert_shipmentprice_PricingDetailOverweights_shipmentpricemodel_PricingDetailOverweights(args []*shipmentprice.PricingDetailOverweight) (outs []*shipmentpricemodel.PricingDetailOverweight) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentpricemodel.PricingDetailOverweight, len(args))
	outs = make([]*shipmentpricemodel.PricingDetailOverweight, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentprice_PricingDetailOverweight_shipmentpricemodel_PricingDetailOverweight(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/main/shipmentpricing/shipmentprice.ShipmentPrice --//

func Convert_shipmentpricemodel_ShipmentPrice_shipmentprice_ShipmentPrice(arg *shipmentpricemodel.ShipmentPrice, out *shipmentprice.ShipmentPrice) *shipmentprice.ShipmentPrice {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentprice.ShipmentPrice{}
	}
	convert_shipmentpricemodel_ShipmentPrice_shipmentprice_ShipmentPrice(arg, out)
	return out
}

func convert_shipmentpricemodel_ShipmentPrice_shipmentprice_ShipmentPrice(arg *shipmentpricemodel.ShipmentPrice, out *shipmentprice.ShipmentPrice) {
	out.ID = arg.ID                                   // simple assign
	out.ShipmentPriceListID = arg.ShipmentPriceListID // simple assign
	out.ShipmentServiceID = arg.ShipmentServiceID     // simple assign
	out.Name = arg.Name                               // simple assign
	out.CustomRegionTypes = arg.CustomRegionTypes     // simple assign
	out.CustomRegionIDs = arg.CustomRegionIDs         // simple assign
	out.RegionTypes = arg.RegionTypes                 // simple assign
	out.ProvinceTypes = arg.ProvinceTypes             // simple assign
	out.UrbanTypes = arg.UrbanTypes                   // simple assign
	out.Details = Convert_shipmentpricemodel_PricingDetails_shipmentprice_PricingDetails(arg.Details)
	out.AdditionalFees = Convert_shipmentpricemodel_AdditionalFees_shipmentprice_AdditionalFees(arg.AdditionalFees)
	out.PriorityPoint = arg.PriorityPoint // simple assign
	out.CreatedAt = arg.CreatedAt         // simple assign
	out.UpdatedAt = arg.UpdatedAt         // simple assign
	out.DeletedAt = arg.DeletedAt         // simple assign
	out.WLPartnerID = arg.WLPartnerID     // simple assign
	out.Status = arg.Status               // simple assign
}

func Convert_shipmentpricemodel_ShipmentPrices_shipmentprice_ShipmentPrices(args []*shipmentpricemodel.ShipmentPrice) (outs []*shipmentprice.ShipmentPrice) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentprice.ShipmentPrice, len(args))
	outs = make([]*shipmentprice.ShipmentPrice, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentpricemodel_ShipmentPrice_shipmentprice_ShipmentPrice(args[i], &tmps[i])
	}
	return outs
}

func Convert_shipmentprice_ShipmentPrice_shipmentpricemodel_ShipmentPrice(arg *shipmentprice.ShipmentPrice, out *shipmentpricemodel.ShipmentPrice) *shipmentpricemodel.ShipmentPrice {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentpricemodel.ShipmentPrice{}
	}
	convert_shipmentprice_ShipmentPrice_shipmentpricemodel_ShipmentPrice(arg, out)
	return out
}

func convert_shipmentprice_ShipmentPrice_shipmentpricemodel_ShipmentPrice(arg *shipmentprice.ShipmentPrice, out *shipmentpricemodel.ShipmentPrice) {
	out.ID = arg.ID                                   // simple assign
	out.ShipmentPriceListID = arg.ShipmentPriceListID // simple assign
	out.ShipmentServiceID = arg.ShipmentServiceID     // simple assign
	out.Name = arg.Name                               // simple assign
	out.CustomRegionTypes = arg.CustomRegionTypes     // simple assign
	out.CustomRegionIDs = arg.CustomRegionIDs         // simple assign
	out.RegionTypes = arg.RegionTypes                 // simple assign
	out.ProvinceTypes = arg.ProvinceTypes             // simple assign
	out.UrbanTypes = arg.UrbanTypes                   // simple assign
	out.Details = Convert_shipmentprice_PricingDetails_shipmentpricemodel_PricingDetails(arg.Details)
	out.AdditionalFees = Convert_shipmentprice_AdditionalFees_shipmentpricemodel_AdditionalFees(arg.AdditionalFees)
	out.PriorityPoint = arg.PriorityPoint // simple assign
	out.CreatedAt = arg.CreatedAt         // simple assign
	out.UpdatedAt = arg.UpdatedAt         // simple assign
	out.DeletedAt = arg.DeletedAt         // simple assign
	out.WLPartnerID = arg.WLPartnerID     // simple assign
	out.Status = arg.Status               // simple assign
}

func Convert_shipmentprice_ShipmentPrices_shipmentpricemodel_ShipmentPrices(args []*shipmentprice.ShipmentPrice) (outs []*shipmentpricemodel.ShipmentPrice) {
	if args == nil {
		return nil
	}
	tmps := make([]shipmentpricemodel.ShipmentPrice, len(args))
	outs = make([]*shipmentpricemodel.ShipmentPrice, len(args))
	for i := range tmps {
		outs[i] = Convert_shipmentprice_ShipmentPrice_shipmentpricemodel_ShipmentPrice(args[i], &tmps[i])
	}
	return outs
}

func Apply_shipmentprice_CreateShipmentPriceArgs_shipmentprice_ShipmentPrice(arg *shipmentprice.CreateShipmentPriceArgs, out *shipmentprice.ShipmentPrice) *shipmentprice.ShipmentPrice {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentprice.ShipmentPrice{}
	}
	apply_shipmentprice_CreateShipmentPriceArgs_shipmentprice_ShipmentPrice(arg, out)
	return out
}

func apply_shipmentprice_CreateShipmentPriceArgs_shipmentprice_ShipmentPrice(arg *shipmentprice.CreateShipmentPriceArgs, out *shipmentprice.ShipmentPrice) {
	out.ID = 0                                        // zero value
	out.ShipmentPriceListID = arg.ShipmentPriceListID // simple assign
	out.ShipmentServiceID = arg.ShipmentServiceID     // simple assign
	out.Name = arg.Name                               // simple assign
	out.CustomRegionTypes = arg.CustomRegionTypes     // simple assign
	out.CustomRegionIDs = arg.CustomRegionIDs         // simple assign
	out.RegionTypes = arg.RegionTypes                 // simple assign
	out.ProvinceTypes = arg.ProvinceTypes             // simple assign
	out.UrbanTypes = arg.UrbanTypes                   // simple assign
	out.Details = arg.Details                         // simple assign
	out.AdditionalFees = arg.AdditionalFees           // simple assign
	out.PriorityPoint = arg.PriorityPoint             // simple assign
	out.CreatedAt = time.Time{}                       // zero value
	out.UpdatedAt = time.Time{}                       // zero value
	out.DeletedAt = time.Time{}                       // zero value
	out.WLPartnerID = 0                               // zero value
	out.Status = 0                                    // zero value
}

func Apply_shipmentprice_UpdateShipmentPriceArgs_shipmentprice_ShipmentPrice(arg *shipmentprice.UpdateShipmentPriceArgs, out *shipmentprice.ShipmentPrice) *shipmentprice.ShipmentPrice {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shipmentprice.ShipmentPrice{}
	}
	apply_shipmentprice_UpdateShipmentPriceArgs_shipmentprice_ShipmentPrice(arg, out)
	return out
}

func apply_shipmentprice_UpdateShipmentPriceArgs_shipmentprice_ShipmentPrice(arg *shipmentprice.UpdateShipmentPriceArgs, out *shipmentprice.ShipmentPrice) {
	out.ID = arg.ID                                   // simple assign
	out.ShipmentPriceListID = arg.ShipmentPriceListID // simple assign
	out.ShipmentServiceID = arg.ShipmentServiceID     // simple assign
	out.Name = arg.Name                               // simple assign
	out.CustomRegionTypes = arg.CustomRegionTypes     // simple assign
	out.CustomRegionIDs = arg.CustomRegionIDs         // simple assign
	out.RegionTypes = arg.RegionTypes                 // simple assign
	out.ProvinceTypes = arg.ProvinceTypes             // simple assign
	out.UrbanTypes = arg.UrbanTypes                   // simple assign
	out.Details = arg.Details                         // simple assign
	out.AdditionalFees = arg.AdditionalFees           // simple assign
	out.PriorityPoint = arg.PriorityPoint             // simple assign
	out.CreatedAt = out.CreatedAt                     // no change
	out.UpdatedAt = out.UpdatedAt                     // no change
	out.DeletedAt = out.DeletedAt                     // no change
	out.WLPartnerID = out.WLPartnerID                 // no change
	out.Status = arg.Status                           // simple assign
}
