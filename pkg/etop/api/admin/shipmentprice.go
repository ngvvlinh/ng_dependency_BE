package admin

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/pricelistpromotion"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/com/main/shipping/carrier"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/admin/convert"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type ShipmentPriceService struct {
	session.Session

	ShipmentManager            *carrier.ShipmentManager
	ShipmentPriceAggr          shipmentprice.CommandBus
	ShipmentPriceQuery         shipmentprice.QueryBus
	ShipmentServiceQuery       shipmentservice.QueryBus
	ShipmentServiceAggr        shipmentservice.CommandBus
	ShipmentPriceListAggr      pricelist.CommandBus
	ShipmentPriceListQuery     pricelist.QueryBus
	ShopShipmentPriceListQuery shopshipmentpricelist.QueryBus
	ShopShipmentPriceListAggr  shopshipmentpricelist.CommandBus
	PriceListPromotionQuery    pricelistpromotion.QueryBus
	PriceListPromotionAggr     pricelistpromotion.CommandBus
}

func (s *ShipmentPriceService) Clone() admin.ShipmentPriceService {
	res := *s
	return &res
}

func (s *ShipmentPriceService) GetShipmentService(ctx context.Context, r *pbcm.IDRequest) (*admin.ShipmentService, error) {
	query := &shipmentservice.GetShipmentServiceQuery{
		ID: r.Id,
	}
	if err := s.ShipmentServiceQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convert.PbShipmentService(query.Result)
	return result, nil
}

func (s *ShipmentPriceService) GetShipmentServices(ctx context.Context, r *admin.GetShipmentServicesRequest) (*admin.GetShipmentServicesResponse, error) {
	query := &shipmentservice.ListShipmentServicesQuery{
		ConnectionID: r.ConnectionID,
	}
	if err := s.ShipmentServiceQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &admin.GetShipmentServicesResponse{
		ShipmentServices: convert.PbShipmentServices(query.Result),
	}
	return result, nil
}

func (s *ShipmentPriceService) CreateShipmentService(ctx context.Context, r *admin.CreateShipmentServiceRequest) (*admin.ShipmentService, error) {
	cmd := &shipmentservice.CreateShipmentServiceCommand{
		ConnectionID:       r.ConnectionID,
		Name:               r.Name,
		EdCode:             r.EdCode,
		ServiceIDs:         r.ServiceIDs,
		Description:        r.Description,
		ImageURL:           r.ImageURL,
		AvailableLocations: convert.AvailableLocations(r.AvailableLocations),
		BlacklistLocations: convert.BlacklistLocations(r.BlacklistLocations),
		OtherCondition:     convert.OtherCondition(r.OtherCondition),
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convert.PbShipmentService(cmd.Result)
	return result, nil
}

func (s *ShipmentPriceService) UpdateShipmentService(ctx context.Context, r *admin.UpdateShipmentServiceRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipmentservice.UpdateShipmentServiceCommand{
		ID:             r.ID,
		ConnectionID:   r.ConnectionID,
		Name:           r.Name,
		EdCode:         r.EdCode,
		ServiceIDs:     r.ServiceIDs,
		Description:    r.Description,
		ImageURL:       r.ImageURL,
		Status:         r.Status,
		OtherCondition: convert.OtherCondition(r.OtherCondition),
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *ShipmentPriceService) DeleteShipmentService(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &shipmentservice.DeleteShipmentServiceCommand{
		ID: r.Id,
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}

func (s *ShipmentPriceService) UpdateShipmentServicesAvailableLocations(ctx context.Context, r *admin.UpdateShipmentServicesAvailableLocationsRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipmentservice.UpdateShipmentServicesLocationConfigCommand{
		IDs:                r.IDs,
		AvailableLocations: convert.AvailableLocations(r.AvailableLocations),
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

func (s *ShipmentPriceService) UpdateShipmentServicesBlacklistLocations(ctx context.Context, r *admin.UpdateShipmentServicesBlacklistLocationsRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shipmentservice.UpdateShipmentServicesLocationConfigCommand{
		IDs:                r.IDs,
		BlacklistLocations: convert.BlacklistLocations(r.BlacklistLocations),
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

//-- End ShipmentService --//

//-- ShipmentPriceList --//

func (s *ShipmentPriceService) GetShipmentPriceList(ctx context.Context, r *pbcm.IDRequest) (*admin.ShipmentPriceList, error) {
	query := &pricelist.GetShipmentPriceListQuery{
		ID: r.Id,
	}
	if err := s.ShipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convert.PbShipmentPriceList(query.Result)
	return result, nil
}

func (s *ShipmentPriceService) GetShipmentPriceLists(ctx context.Context, r *admin.GetShipmentPriceListsRequest) (*admin.GetShipmentPriceListsResponse, error) {
	query := &pricelist.ListShipmentPriceListsQuery{
		ConnectionID: r.ConnectionID,
		IsDefault:    r.IsDefault,
	}
	if err := s.ShipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &admin.GetShipmentPriceListsResponse{
		ShipmentPriceLists: convert.PbShipmentPriceLists(query.Result),
	}
	return result, nil
}

func (s *ShipmentPriceService) CreateShipmentPriceList(ctx context.Context, r *admin.CreateShipmentPriceListRequest) (*admin.ShipmentPriceList, error) {
	cmd := &pricelist.CreateShipmentPriceListCommand{
		Name:         r.Name,
		Description:  r.Description,
		ConnectionID: r.ConnectionID,
		IsDefault:    r.IsDefault,
	}
	if err := s.ShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convert.PbShipmentPriceList(cmd.Result)
	return result, nil
}

func (s *ShipmentPriceService) UpdateShipmentPriceList(ctx context.Context, r *admin.UpdateShipmentPriceListRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &pricelist.UpdateShipmentPriceListCommand{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
	if err := s.ShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *ShipmentPriceService) SetDefaultShipmentPriceList(ctx context.Context, r *admin.ActiveShipmentPriceListRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &pricelist.SetDefaultShipmentPriceListCommand{
		ID:           r.ID,
		ConnectionID: r.ConnectionID,
	}
	if err := s.ShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *ShipmentPriceService) DeleteShipmentPriceList(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &pricelist.DeleteShipmentPriceListCommand{
		ID: r.Id,
	}
	if err := s.ShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}

//-- End ShipmentPriceList --//

//-- ShipmentPrice --//

func (s *ShipmentPriceService) GetShipmentPrice(ctx context.Context, r *pbcm.IDRequest) (*admin.ShipmentPrice, error) {
	query := &shipmentprice.GetShipmentPriceQuery{
		ID: r.Id,
	}
	if err := s.ShipmentPriceQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convert.PbShipmentPrice(query.Result)
	return result, nil
}

func (s *ShipmentPriceService) GetShipmentPrices(ctx context.Context, r *admin.GetShipmentPricesRequest) (*admin.GetShipmentPricesResponse, error) {
	query := &shipmentprice.ListShipmentPricesQuery{
		ShipmentPriceListID: r.ShipmentPriceListID,
		ShipmentServiceID:   r.ShipmentServiceID,
	}
	if err := s.ShipmentPriceQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &admin.GetShipmentPricesResponse{
		ShipmentPrices: convert.PbShipmentPrices(query.Result),
	}
	return result, nil
}

func (s *ShipmentPriceService) CreateShipmentPrice(ctx context.Context, r *admin.CreateShipmentPriceRequest) (*admin.ShipmentPrice, error) {
	cmd := &shipmentprice.CreateShipmentPriceCommand{
		Name:                r.Name,
		ShipmentPriceListID: r.ShipmentPriceListID,
		ShipmentServiceID:   r.ShipmentServiceID,
		CustomRegionTypes:   r.CustomRegionTypes,
		CustomRegionIDs:     r.CustomRegionIDs,
		RegionTypes:         r.RegionTypes,
		ProvinceTypes:       r.ProvinceTypes,
		UrbanTypes:          r.UrbanTypes,
		PriorityPoint:       r.PriorityPoint,
		Details:             convert.PricingDetails(r.Details),
		AdditionalFees:      convert.Convert_api_AdditionalFees_To_core_AdditionalFees(r.AdditionalFees),
	}
	if err := s.ShipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convert.PbShipmentPrice(cmd.Result)
	return result, nil
}

func (s *ShipmentPriceService) UpdateShipmentPrice(ctx context.Context, r *admin.UpdateShipmentPriceRequest) (*admin.ShipmentPrice, error) {
	cmd := &shipmentprice.UpdateShipmentPriceCommand{
		ID:                  r.ID,
		Name:                r.Name,
		ShipmentPriceListID: r.ShipmentPriceListID,
		ShipmentServiceID:   r.ShipmentServiceID,
		CustomRegionTypes:   r.CustomRegionTypes,
		CustomRegionIDs:     r.CustomRegionIDs,
		RegionTypes:         r.RegionTypes,
		ProvinceTypes:       r.ProvinceTypes,
		UrbanTypes:          r.UrbanTypes,
		PriorityPoint:       r.PriorityPoint,
		Details:             convert.PricingDetails(r.Details),
		AdditionalFees:      convert.Convert_api_AdditionalFees_To_core_AdditionalFees(r.AdditionalFees),
		Status:              r.Status,
	}
	if err := s.ShipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convert.PbShipmentPrice(cmd.Result)
	return result, nil
}

func (s *ShipmentPriceService) DeleteShipmentPrice(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &shipmentprice.DeleteShipmentPriceCommand{
		ID: r.Id,
	}
	if err := s.ShipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}

func (s *ShipmentPriceService) UpdateShipmentPricesPriorityPoint(ctx context.Context, r *admin.UpdateShipmentPricesPriorityPointRequest) (*pbcm.UpdatedResponse, error) {
	updates := []*shipmentprice.UpdateShipmentPricePriorityPointArgs{}
	for _, sp := range r.ShipmentPrices {
		updates = append(updates, &shipmentprice.UpdateShipmentPricePriorityPointArgs{
			ID:            sp.ID,
			PriorityPoint: sp.PriorityPoint,
		})
	}
	cmd := &shipmentprice.UpdateShipmentPricesPriorityPointCommand{
		ShipmentPrices: updates,
	}
	if err := s.ShipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

//-- End ShipmentPrice --//

//-- Shop Shipment Price List --//

func (s *ShipmentPriceService) GetShopShipmentPriceLists(ctx context.Context, r *admin.GetShopShipmentPriceListsRequest) (*admin.GetShopShipmentPriceListsResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	query := &shopshipmentpricelist.ListShopShipmentPriceListsQuery{
		ShipmentPriceListID: r.ShipmentPriceListID,
		ConnectionID:        r.ConnectionID,
		ShopID:              r.ShopID,
		Paging:              *paging,
	}
	if err := s.ShopShipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &admin.GetShopShipmentPriceListsResponse{
		PriceLists: convert.PbShopShipmentPriceLists(query.Result.ShopShipmentPriceLists),
		Paging:     cmapi.PbMetaPageInfo(query.Result.Paging),
	}
	return result, nil
}

func (s *ShipmentPriceService) GetShopShipmentPriceList(ctx context.Context, r *admin.GetShopShipmentPriceListRequest) (*admin.ShopShipmentPriceList, error) {
	query := &shopshipmentpricelist.GetShopShipmentPriceListQuery{
		ShopID: r.ShopID,
	}
	if err := s.ShopShipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convert.PbShopShipmentPriceList(query.Result)
	return result, nil
}

func (s *ShipmentPriceService) CreateShopShipmentPriceList(ctx context.Context, r *admin.CreateShopShipmentPriceList) (*admin.ShopShipmentPriceList, error) {
	cmd := &shopshipmentpricelist.CreateShopShipmentPriceListCommand{
		ShopID:              r.ShopID,
		ShipmentPriceListID: r.ShipmentPriceListID,
		ConnectionID:        r.ConnectionID,
		Note:                r.Note,
		UpdatedBy:           s.SS.Claim().UserID,
	}
	if err := s.ShopShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convert.PbShopShipmentPriceList(cmd.Result)
	return result, nil
}

func (s *ShipmentPriceService) UpdateShopShipmentPriceList(ctx context.Context, r *admin.UpdateShopShipmentPriceListRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &shopshipmentpricelist.UpdateShopShipmentPriceListCommand{
		ShopID:              r.ShopID,
		ShipmentPriceListID: r.ShipmentPriceListID,
		ConnectionID:        r.ConnectionID,
		Note:                r.Note,
		UpdatedBy:           s.SS.Claim().UserID,
	}
	if err := s.ShopShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *ShipmentPriceService) DeleteShopShipmentPriceList(ctx context.Context, r *admin.GetShopShipmentPriceListRequest) (*pbcm.DeletedResponse, error) {
	cmd := &shopshipmentpricelist.DeleteShopShipmentPriceListCommand{
		ShopID:       r.ShopID,
		ConnectionID: r.ConnectionID,
	}
	if err := s.ShopShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}

//-- End Shop Shipment Price List --//

func (s *ShipmentPriceService) GetShippingServices(ctx context.Context, r *admin.GetShippingServicesRequest) (*types.GetShippingServicesResponse, error) {
	if r.ShipmentPriceListID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui l??ng ch???n b???ng gi?? ??p d???ng")
	}
	args := &carrier.GetShippingServicesArgs{
		AccountID:           r.AccountID,
		ShipmentPriceListID: r.ShipmentPriceListID,
		ConnectionIDs:       r.ConnectionIDs,
		FromDistrictCode:    r.FromDistrictCode,
		FromProvinceCode:    r.FromProvinceCode,
		ToDistrictCode:      r.ToDistrictCode,
		ToProvinceCode:      r.ToProvinceCode,
		ChargeableWeight:    r.GrossWeight,
		IncludeInsurance:    r.IncludeInsurance.Apply(false),
		BasketValue:         r.BasketValue,
		CODAmount:           r.TotalCodAmount,
	}
	resp, err := s.ShipmentManager.GetShippingServices(ctx, args)
	if err != nil {
		return nil, err
	}
	result := &types.GetShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return result, nil
}

//-- Shipment Price List Promotion --//

func (s *ShipmentPriceService) GetShipmentPriceListPromotions(ctx context.Context, r *admin.GetShipmentPriceListPromotionsRequest) (*admin.GetShipmentPriceListPromotionsResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	query := &pricelistpromotion.ListPriceListPromotionsQuery{
		ConnectionID: r.ConnectionID,
		PriceListID:  r.ShipmentPriceListID,
		Paging:       *paging,
	}
	if err := s.PriceListPromotionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var res = &admin.GetShipmentPriceListPromotionsResponse{
		ShipmentPriceListPromotions: convert.Convert_core_PriceListPromotions_To_api_PriceListPromotions(query.Result),
	}
	return res, nil
}

func (s *ShipmentPriceService) GetShipmentPriceListPromotion(ctx context.Context, r *pbcm.IDRequest) (*admin.ShipmentPriceListPromotion, error) {
	query := &pricelistpromotion.GetPriceListPromotionQuery{
		ID: r.Id,
	}
	if err := s.PriceListPromotionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convert.Convert_core_PriceListPromotion_To_api_PriceListPromotion(query.Result), nil
}

func (s *ShipmentPriceService) CreateShipmentPriceListPromotion(ctx context.Context, r *admin.CreateShipmentPriceListPromotionRequest) (*admin.ShipmentPriceListPromotion, error) {
	cmd := &pricelistpromotion.CreatePriceListPromotionCommand{
		PriceListID:   r.ShipmentPriceListID,
		Name:          r.Name,
		Description:   r.Description,
		ConnectionID:  r.ConnectionID,
		DateFrom:      r.DateFrom,
		DateTo:        r.DateTo,
		AppliedRules:  convert.Convert_api_PriceListPromotionAppliedRules_To_core_PriceListPromotionAppliedRules(r.AppliedRules),
		PriorityPoint: r.PriorityPoint,
	}
	if err := s.PriceListPromotionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convert.Convert_core_PriceListPromotion_To_api_PriceListPromotion(cmd.Result), nil
}

func (s *ShipmentPriceService) UpdateShipmentPriceListPromotion(ctx context.Context, r *admin.UpdateShipmentPriceListPromotionRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &pricelistpromotion.UpdatePriceListPromotionCommand{
		ID:            r.ID,
		Name:          r.Name,
		Description:   r.Description,
		DateFrom:      r.DateFrom,
		DateTo:        r.DateTo,
		AppliedRules:  convert.Convert_api_PriceListPromotionAppliedRules_To_core_PriceListPromotionAppliedRules(r.AppliedRules),
		PriorityPoint: r.PriorityPoint,
		Status:        r.Status,
		ConnectionID:  r.ConnectionID,
		PriceListID:   r.ShipmentPriceListID,
	}
	if err := s.PriceListPromotionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func (s *ShipmentPriceService) DeleteShipmentPriceListPromotion(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &pricelistpromotion.DeletePriceListPromotionCommand{
		ID: r.Id,
	}
	if err := s.PriceListPromotionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.DeletedResponse{Deleted: 1}, nil
}

// -- End Shipment Price List Promotion --//
