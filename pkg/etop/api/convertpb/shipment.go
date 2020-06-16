package convertpb

import (
	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	"o.o/api/top/int/admin"
)

func PbShipmentService(in *shipmentservice.ShipmentService) *admin.ShipmentService {
	if in == nil {
		return nil
	}
	return &admin.ShipmentService{
		ID:                 in.ID,
		ConnectionID:       in.ConnectionID,
		Name:               in.Name,
		EdCode:             in.EdCode,
		ServiceIDs:         in.ServiceIDs,
		Description:        in.Description,
		CreatedAt:          in.CreatedAt,
		UpdatedAt:          in.UpdatedAt,
		Status:             in.Status,
		ImageURL:           in.ImageURL,
		AvailableLocations: PbAvailableLocations(in.AvailableLocations),
		BlacklistLocations: PbBlacklistLocations(in.BlacklistLocations),
		OtherCondition:     PbOtherCondition(in.OtherCondition),
	}
}

func PbAvailableLocation(in *shipmentservice.AvailableLocation) *admin.AvailableLocation {
	if in == nil {
		return nil
	}
	return &admin.AvailableLocation{
		FilterType:           in.FilterType,
		ShippingLocationType: in.ShippingLocationType,
		RegionTypes:          in.RegionTypes,
		CustomRegionIDs:      in.CustomRegionIDs,
		ProvinceCodes:        in.ProvinceCodes,
	}
}

func PbAvailableLocations(items []*shipmentservice.AvailableLocation) []*admin.AvailableLocation {
	result := make([]*admin.AvailableLocation, len(items))
	for i, item := range items {
		result[i] = PbAvailableLocation(item)
	}
	return result
}

func AvailableLocation(in *admin.AvailableLocation) *shipmentservice.AvailableLocation {
	if in == nil {
		return nil
	}
	return &shipmentservice.AvailableLocation{
		FilterType:           in.FilterType,
		ShippingLocationType: in.ShippingLocationType,
		RegionTypes:          in.RegionTypes,
		CustomRegionIDs:      in.CustomRegionIDs,
		ProvinceCodes:        in.ProvinceCodes,
	}
}

func AvailableLocations(items []*admin.AvailableLocation) []*shipmentservice.AvailableLocation {
	result := make([]*shipmentservice.AvailableLocation, len(items))
	for i, item := range items {
		result[i] = AvailableLocation(item)
	}
	return result
}

func PbBlacklistLocation(in *shipmentservice.BlacklistLocation) *admin.BlacklistLocation {
	if in == nil {
		return nil
	}
	return &admin.BlacklistLocation{
		ShippingLocationType: in.ShippingLocationType,
		ProvinceCodes:        in.ProvinceCodes,
		DistrictCodes:        in.DistrictCodes,
		WardCodes:            in.WardCodes,
		Reason:               in.Reason,
	}
}

func PbBlacklistLocations(items []*shipmentservice.BlacklistLocation) []*admin.BlacklistLocation {
	result := make([]*admin.BlacklistLocation, len(items))
	for i, item := range items {
		result[i] = PbBlacklistLocation(item)
	}
	return result
}

func BlacklistLocation(in *admin.BlacklistLocation) *shipmentservice.BlacklistLocation {
	if in == nil {
		return nil
	}
	return &shipmentservice.BlacklistLocation{
		ShippingLocationType: in.ShippingLocationType,
		ProvinceCodes:        in.ProvinceCodes,
		DistrictCodes:        in.DistrictCodes,
		WardCodes:            in.WardCodes,
		Reason:               in.Reason,
	}
}

func BlacklistLocations(items []*admin.BlacklistLocation) []*shipmentservice.BlacklistLocation {
	result := make([]*shipmentservice.BlacklistLocation, len(items))
	for i, item := range items {
		result[i] = BlacklistLocation(item)
	}
	return result
}

func PbOtherCondition(in *shipmentservice.OtherCondition) *admin.OtherCondition {
	if in == nil {
		return nil
	}
	return &admin.OtherCondition{
		MinWeight: in.MinWeight,
		MaxWeight: in.MaxWeight,
	}
}

func OtherCondition(in *admin.OtherCondition) *shipmentservice.OtherCondition {
	if in == nil {
		return nil
	}
	return &shipmentservice.OtherCondition{
		MinWeight: in.MinWeight,
		MaxWeight: in.MaxWeight,
	}
}

func PbShipmentServices(items []*shipmentservice.ShipmentService) []*admin.ShipmentService {
	result := make([]*admin.ShipmentService, len(items))
	for i, item := range items {
		result[i] = PbShipmentService(item)
	}
	return result
}

func PbShipmentPriceList(in *pricelist.ShipmentPriceList) *admin.ShipmentPriceList {
	if in == nil {
		return nil
	}
	return &admin.ShipmentPriceList{
		ID:           in.ID,
		Name:         in.Name,
		Description:  in.Description,
		IsDefault:    in.IsDefault,
		ConnectionID: in.ConnectionID,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
	}
}

func PbShipmentPriceLists(items []*pricelist.ShipmentPriceList) []*admin.ShipmentPriceList {
	result := make([]*admin.ShipmentPriceList, len(items))
	for i, item := range items {
		result[i] = PbShipmentPriceList(item)
	}
	return result
}

func PbShipmentPrice(in *shipmentprice.ShipmentPrice) *admin.ShipmentPrice {
	if in == nil {
		return nil
	}
	return &admin.ShipmentPrice{
		ID:                  in.ID,
		ShipmentPriceListID: in.ShipmentPriceListID,
		ShipmentServiceID:   in.ShipmentServiceID,
		Name:                in.Name,
		CustomRegionTypes:   in.CustomRegionTypes,
		CustomRegionIDs:     in.CustomRegionIDs,
		RegionTypes:         in.RegionTypes,
		ProvinceTypes:       in.ProvinceTypes,
		UrbanTypes:          in.UrbanTypes,
		PriorityPoint:       in.PriorityPoint,
		Details:             PbPricingDetails(in.Details),
		CreatedAt:           in.CreatedAt,
		UpdatedAt:           in.UpdatedAt,
		Status:              in.Status,
		AdditionalFees:      Convert_core_AdditionalFees_To_api_AdditionalFees(in.AdditionalFees),
	}
}

func PbPricingDetail(in *shipmentprice.PricingDetail) *admin.PricingDetail {
	if in == nil {
		return nil
	}
	return &admin.PricingDetail{
		Weight:     in.Weight,
		Price:      in.Price,
		Overweight: PbPricingDetailOverweight(in.Overweight),
	}
}

func PbPricingDetailOverweight(ins []*shipmentprice.PricingDetailOverweight) (res []*admin.PricingDetailOverweight) {
	for _, in := range ins {
		res = append(res, &admin.PricingDetailOverweight{
			MinWeight:  in.MinWeight,
			MaxWeight:  in.MaxWeight,
			WeightStep: in.WeightStep,
			PriceStep:  in.PriceStep,
		})
	}
	return
}

func PbPricingDetails(items []*shipmentprice.PricingDetail) []*admin.PricingDetail {
	result := make([]*admin.PricingDetail, len(items))
	for i, item := range items {
		result[i] = PbPricingDetail(item)
	}
	return result
}

func PbShipmentPrices(items []*shipmentprice.ShipmentPrice) []*admin.ShipmentPrice {
	result := make([]*admin.ShipmentPrice, len(items))
	for i, item := range items {
		result[i] = PbShipmentPrice(item)
	}
	return result
}

func PricingDetails(ins []*admin.PricingDetail) (res []*shipmentprice.PricingDetail) {
	for _, in := range ins {
		res = append(res, &shipmentprice.PricingDetail{
			Weight:     in.Weight,
			Price:      in.Price,
			Overweight: PricingDetailOverweights(in.Overweight),
		})
	}
	return
}

func PricingDetailOverweights(items []*admin.PricingDetailOverweight) (res []*shipmentprice.PricingDetailOverweight) {
	for _, item := range items {
		res = append(res, &shipmentprice.PricingDetailOverweight{
			MinWeight:  item.MinWeight,
			MaxWeight:  item.MaxWeight,
			WeightStep: item.WeightStep,
			PriceStep:  item.PriceStep,
		})
	}
	return
}

func Convert_api_AdditionalFee_To_core_AdditionalFee(in *admin.AdditionalFee) *shipmentprice.AdditionalFee {
	if in == nil {
		return nil
	}
	return &shipmentprice.AdditionalFee{
		FeeType: in.FeeType,
		Rules:   Convert_api_AdditionalFeeRules_To_core_AdditionalFeeRules(in.Rules),
	}
}

func Convert_api_AdditionalFees_To_core_AdditionalFees(items []*admin.AdditionalFee) (res []*shipmentprice.AdditionalFee) {
	for _, item := range items {
		res = append(res, Convert_api_AdditionalFee_To_core_AdditionalFee(item))
	}
	return
}

func Convert_api_AdditionalFeeRule_To_core_AdditionalFeeRule(in *admin.AdditionalFeeRule) *shipmentprice.AdditionalFeeRule {
	if in == nil {
		return nil
	}
	return &shipmentprice.AdditionalFeeRule{
		MinValue:          in.MinValue,
		MaxValue:          in.MaxValue,
		PriceModifierType: in.PriceModifierType,
		Amount:            in.Amount,
		MinPrice:          in.MinPrice,
	}
}

func Convert_api_AdditionalFeeRules_To_core_AdditionalFeeRules(items []*admin.AdditionalFeeRule) []*shipmentprice.AdditionalFeeRule {
	result := make([]*shipmentprice.AdditionalFeeRule, len(items))
	for i, item := range items {
		result[i] = Convert_api_AdditionalFeeRule_To_core_AdditionalFeeRule(item)
	}
	return result
}

func Convert_core_AdditionalFee_To_api_AdditionalFee(in *shipmentprice.AdditionalFee) *admin.AdditionalFee {
	if in == nil {
		return nil
	}
	return &admin.AdditionalFee{
		FeeType: in.FeeType,
		Rules:   Convert_core_AdditionalFeeRules_To_api_AdditionalFeeRules(in.Rules),
	}
}

func Convert_core_AdditionalFees_To_api_AdditionalFees(items []*shipmentprice.AdditionalFee) []*admin.AdditionalFee {
	result := make([]*admin.AdditionalFee, len(items))
	for i, item := range items {
		result[i] = Convert_core_AdditionalFee_To_api_AdditionalFee(item)
	}
	return result
}

func Convert_core_AdditionalFeeRule_To_api_AdditionalFeeRule(in *shipmentprice.AdditionalFeeRule) *admin.AdditionalFeeRule {
	if in == nil {
		return nil
	}
	return &admin.AdditionalFeeRule{
		MinValue:          in.MinValue,
		MaxValue:          in.MaxValue,
		PriceModifierType: in.PriceModifierType,
		Amount:            in.Amount,
		MinPrice:          in.MinPrice,
	}
}

func Convert_core_AdditionalFeeRules_To_api_AdditionalFeeRules(items []*shipmentprice.AdditionalFeeRule) []*admin.AdditionalFeeRule {
	result := make([]*admin.AdditionalFeeRule, len(items))
	for i, item := range items {
		result[i] = Convert_core_AdditionalFeeRule_To_api_AdditionalFeeRule(item)
	}
	return result
}

func PbShopShipmentPriceList(in *shopshipmentpricelist.ShopShipmentPriceList) *admin.ShopShipmentPriceList {
	if in == nil {
		return nil
	}
	return &admin.ShopShipmentPriceList{
		ShopID:              in.ShopID,
		ShipmentPriceListID: in.ShipmentPriceListID,
		ConnectionID:        in.ConnectionID,
		Note:                in.Note,
		CreatedAt:           in.CreatedAt,
		UpdatedAt:           in.UpdatedAt,
		UpdatedBy:           in.UpdatedBy,
	}
}

func PbShopShipmentPriceLists(items []*shopshipmentpricelist.ShopShipmentPriceList) []*admin.ShopShipmentPriceList {
	result := make([]*admin.ShopShipmentPriceList, len(items))
	for i, item := range items {
		result[i] = PbShopShipmentPriceList(item)
	}
	return result
}
