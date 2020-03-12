package convertpb

import (
	"etop.vn/api/main/shipmentpricing/pricelist"
	"etop.vn/api/main/shipmentpricing/shipmentprice"
	"etop.vn/api/main/shipmentpricing/shipmentservice"
	"etop.vn/api/top/int/admin"
)

func PbShipmentService(in *shipmentservice.ShipmentService) *admin.ShipmentService {
	if in == nil {
		return nil
	}
	return &admin.ShipmentService{
		ID:           in.ID,
		ConnectionID: in.ConnectionID,
		Name:         in.Name,
		EdCode:       in.EdCode,
		ServiceIDs:   in.ServiceIDs,
		Description:  in.Description,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		Status:       in.Status,
		ImageURL:     in.ImageURL,
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
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		IsActive:    in.IsActive,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
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

func PricingDetails(ins []*admin.PricingDetail) []*shipmentprice.PricingDetail {
	result := make([]*shipmentprice.PricingDetail, len(ins))
	for i, in := range ins {
		result[i] = &shipmentprice.PricingDetail{
			Weight:     in.Weight,
			Price:      in.Price,
			Overweight: PricingDetailOverweights(in.Overweight),
		}
	}
	return result
}

func PricingDetailOverweights(items []*admin.PricingDetailOverweight) []*shipmentprice.PricingDetailOverweight {
	result := make([]*shipmentprice.PricingDetailOverweight, len(items))
	for i, item := range items {
		result[i] = &shipmentprice.PricingDetailOverweight{
			MinWeight:  item.MinWeight,
			MaxWeight:  item.MaxWeight,
			WeightStep: item.WeightStep,
			PriceStep:  item.PriceStep,
		}
	}
	return result
}
