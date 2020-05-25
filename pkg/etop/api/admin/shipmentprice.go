package admin

import (
	"context"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	"o.o/api/main/shipmentpricing/subpricelist"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/com/main/shipping/carrier"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

type ShipmentPriceService struct {
	ShipmentManager            *carrier.ShipmentManager
	ShipmentPriceAggr          shipmentprice.CommandBus
	ShipmentPriceQuery         shipmentprice.QueryBus
	ShipmentServiceQuery       shipmentservice.QueryBus
	ShipmentServiceAggr        shipmentservice.CommandBus
	ShipmentPriceListAggr      pricelist.CommandBus
	ShipmentPriceListQuery     pricelist.QueryBus
	ShipmentSubPriceListQuery  subpricelist.QueryBus
	ShipmentSubPriceListAggr   subpricelist.CommandBus
	ShopShipmentPriceListQuery shopshipmentpricelist.QueryBus
	ShopShipmentPriceListAggr  shopshipmentpricelist.CommandBus
}

func (s *ShipmentPriceService) Clone() *ShipmentPriceService {
	res := *s
	return &res
}

func (s *ShipmentPriceService) GetShipmentService(ctx context.Context, r *GetShipmentServiceEndpoint) error {
	query := &shipmentservice.GetShipmentServiceQuery{
		ID: r.Id,
	}
	if err := s.ShipmentServiceQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentService(query.Result)
	return nil
}

func (s *ShipmentPriceService) GetShipmentServices(ctx context.Context, r *GetShipmentServicesEndpoint) error {
	query := &shipmentservice.ListShipmentServicesQuery{}
	if err := s.ShipmentServiceQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &admin.GetShipmentServicesResponse{
		ShipmentServices: convertpb.PbShipmentServices(query.Result),
	}
	return nil
}

func (s *ShipmentPriceService) CreateShipmentService(ctx context.Context, r *CreateShipmentServiceEndpoint) error {
	cmd := &shipmentservice.CreateShipmentServiceCommand{
		ConnectionID:       r.ConnectionID,
		Name:               r.Name,
		EdCode:             r.EdCode,
		ServiceIDs:         r.ServiceIDs,
		Description:        r.Description,
		ImageURL:           r.ImageURL,
		AvailableLocations: convertpb.AvailableLocations(r.AvailableLocations),
		BlacklistLocations: convertpb.BlacklistLocations(r.BlacklistLocations),
		OtherCondition:     convertpb.OtherCondition(r.OtherCondition),
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentService(cmd.Result)
	return nil
}

func (s *ShipmentPriceService) UpdateShipmentService(ctx context.Context, r *UpdateShipmentServiceEndpoint) error {
	cmd := &shipmentservice.UpdateShipmentServiceCommand{
		ID:             r.ID,
		ConnectionID:   r.ConnectionID,
		Name:           r.Name,
		EdCode:         r.EdCode,
		ServiceIDs:     r.ServiceIDs,
		Description:    r.Description,
		ImageURL:       r.ImageURL,
		Status:         r.Status,
		OtherCondition: convertpb.OtherCondition(r.OtherCondition),
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ShipmentPriceService) DeleteShipmentService(ctx context.Context, r *DeleteShipmentServiceEndpoint) error {
	cmd := &shipmentservice.DeleteShipmentServiceCommand{
		ID: r.Id,
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

func (s *ShipmentPriceService) UpdateShipmentServicesAvailableLocations(ctx context.Context, r *UpdateShipmentServicesAvailableLocationsEndpoint) error {
	cmd := &shipmentservice.UpdateShipmentServicesLocationConfigCommand{
		IDs:                r.IDs,
		AvailableLocations: convertpb.AvailableLocations(r.AvailableLocations),
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *ShipmentPriceService) UpdateShipmentServicesBlacklistLocations(ctx context.Context, r *UpdateShipmentServicesBlacklistLocationsEndpoint) error {
	cmd := &shipmentservice.UpdateShipmentServicesLocationConfigCommand{
		IDs:                r.IDs,
		BlacklistLocations: convertpb.BlacklistLocations(r.BlacklistLocations),
	}
	if err := s.ShipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

//-- End ShipmentService --//

//-- ShipmentSubPriceList --//

func (s *ShipmentPriceService) GetShipmentSubPriceLists(ctx context.Context, r *GetShipmentSubPriceListsEndpoint) error {
	query := &subpricelist.ListShipmentSubPriceListQuery{
		ConnectionID: r.ConnectionID,
	}
	if err := s.ShipmentSubPriceListQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &admin.GetShipmentSubPriceListsResponse{
		ShipmentSubPriceLists: convertpb.PbShipmentSubPriceLists(query.Result),
	}
	return nil
}

func (s *ShipmentPriceService) GetShipmentSubPriceList(ctx context.Context, r *GetShipmentSubPriceListEndpoint) error {
	query := &subpricelist.GetShipmentSubPriceListQuery{
		ID: r.Id,
	}
	if err := s.ShipmentSubPriceListQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentSubPriceList(query.Result)
	return nil
}

func (s *ShipmentPriceService) CreateShipmentSubPriceList(ctx context.Context, r *CreateShipmentSubPriceListEndpoint) error {
	cmd := &subpricelist.CreateShipmentSubPriceListCommand{
		Name:         r.Name,
		Description:  r.Description,
		ConnectionID: r.ConnectionID,
	}
	if err := s.ShipmentSubPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentSubPriceList(cmd.Result)
	return nil
}

func (s *ShipmentPriceService) UpdateShipmentSubPriceList(ctx context.Context, r *UpdateShipmentSubPriceListEndpoint) error {
	cmd := &subpricelist.UpdateShipmentSubPriceListCommand{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
	}
	if err := s.ShipmentSubPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ShipmentPriceService) DeleteShipmentSubPriceList(ctx context.Context, r *DeleteShipmentSubPriceListEndpoint) error {
	cmd := &subpricelist.DeleteShipmentSubPriceListCommand{
		ID: r.Id,
	}
	if err := s.ShipmentSubPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

//-- End ShipmentSubPriceList --//

//-- ShipmentPriceList --//

func (s *ShipmentPriceService) GetShipmentPriceList(ctx context.Context, r *GetShipmentPriceListEndpoint) error {
	query := &pricelist.GetShipmentPriceListQuery{
		ID: r.Id,
	}
	if err := s.ShipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPriceList(query.Result)
	return nil
}

func (s *ShipmentPriceService) GetShipmentPriceLists(ctx context.Context, r *GetShipmentPriceListsEndpoint) error {
	query := &pricelist.ListShipmentPriceListsQuery{}
	if err := s.ShipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &admin.GetShipmentPriceListsResponse{
		ShipmentPriceLists: convertpb.PbShipmentPriceLists(query.Result),
	}
	return nil
}

func (s *ShipmentPriceService) CreateShipmentPriceList(ctx context.Context, r *CreateShipmentPriceListEndpoint) error {
	cmd := &pricelist.CreateShipmentPriceListCommand{
		Name:                    r.Name,
		Description:             r.Description,
		IsActive:                r.IsActive,
		ShipmentSubPriceListIDs: r.ShipmentSubPriceListIDs,
	}
	if err := s.ShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPriceList(cmd.Result)
	return nil
}

func (s *ShipmentPriceService) UpdateShipmentPriceList(ctx context.Context, r *UpdateShipmentPriceListEndpoint) error {
	cmd := &pricelist.UpdateShipmentPriceListCommand{
		ID:                      r.ID,
		Name:                    r.Name,
		Description:             r.Description,
		ShipmentSubPriceListIDs: r.ShipmentSubPriceListIDs,
	}
	if err := s.ShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ShipmentPriceService) ActivateShipmentPriceList(ctx context.Context, r *ActivateShipmentPriceListEndpoint) error {
	cmd := &pricelist.ActivateShipmentPriceListCommand{
		ID: r.Id,
	}
	if err := s.ShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ShipmentPriceService) DeleteShipmentPriceList(ctx context.Context, r *DeleteShipmentPriceListEndpoint) error {
	cmd := &pricelist.DeleteShipmentPriceListCommand{
		ID: r.Id,
	}
	if err := s.ShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

//-- End ShipmentPriceList --//

//-- ShipmentPrice --//

func (s *ShipmentPriceService) GetShipmentPrice(ctx context.Context, r *GetShipmentPriceEndpoint) error {
	query := &shipmentprice.GetShipmentPriceQuery{
		ID: r.Id,
	}
	if err := s.ShipmentPriceQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPrice(query.Result)
	return nil
}

func (s *ShipmentPriceService) GetShipmentPrices(ctx context.Context, r *GetShipmentPricesEndpoint) error {
	query := &shipmentprice.ListShipmentPricesQuery{
		ShipmentSubPriceListID: r.ShipmentSubPriceListID,
		ShipmentServiceID:      r.ShipmentServiceID,
	}
	if err := s.ShipmentPriceQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &admin.GetShipmentPricesResponse{
		ShipmentPrices: convertpb.PbShipmentPrices(query.Result),
	}
	return nil
}

func (s *ShipmentPriceService) CreateShipmentPrice(ctx context.Context, r *CreateShipmentPriceEndpoint) error {
	cmd := &shipmentprice.CreateShipmentPriceCommand{
		Name:                   r.Name,
		ShipmentSubPriceListID: r.ShipmentSubPriceListID,
		ShipmentServiceID:      r.ShipmentServiceID,
		CustomRegionTypes:      r.CustomRegionTypes,
		CustomRegionIDs:        r.CustomRegionIDs,
		RegionTypes:            r.RegionTypes,
		ProvinceTypes:          r.ProvinceTypes,
		UrbanTypes:             r.UrbanTypes,
		PriorityPoint:          r.PriorityPoint,
		Details:                convertpb.PricingDetails(r.Details),
	}
	if err := s.ShipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPrice(cmd.Result)
	return nil
}

func (s *ShipmentPriceService) UpdateShipmentPrice(ctx context.Context, r *UpdateShipmentPriceEndpoint) error {
	cmd := &shipmentprice.UpdateShipmentPriceCommand{
		ID:                     r.ID,
		Name:                   r.Name,
		ShipmentSubPriceListID: r.ShipmentSubPriceListID,
		ShipmentServiceID:      r.ShipmentServiceID,
		CustomRegionTypes:      r.CustomRegionTypes,
		CustomRegionIDs:        r.CustomRegionIDs,
		RegionTypes:            r.RegionTypes,
		ProvinceTypes:          r.ProvinceTypes,
		UrbanTypes:             r.UrbanTypes,
		PriorityPoint:          r.PriorityPoint,
		Details:                convertpb.PricingDetails(r.Details),
		Status:                 r.Status,
	}
	if err := s.ShipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPrice(cmd.Result)
	return nil
}

func (s *ShipmentPriceService) DeleteShipmentPrice(ctx context.Context, r *DeleteShipmentPriceEndpoint) error {
	cmd := &shipmentprice.DeleteShipmentPriceCommand{
		ID: r.Id,
	}
	if err := s.ShipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

func (s *ShipmentPriceService) UpdateShipmentPricesPriorityPoint(ctx context.Context, r *UpdateShipmentPricesPriorityPointEndpoint) error {
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
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

//-- End ShipmentPrice --//

//-- Shop Shipment Price List --//

func (s *ShipmentPriceService) GetShopShipmentPriceLists(ctx context.Context, r *GetShopShipmentPriceListsEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	query := &shopshipmentpricelist.ListShopShipmentPriceListsQuery{
		ShipmentPriceListID: r.ShipmentPriceListID,
		Paging:              *paging,
	}
	if err := s.ShopShipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &admin.GetShopShipmentPriceListsResponse{
		PriceLists: convertpb.PbShopShipmentPriceLists(query.Result.ShopShipmentPriceLists),
		Paging:     cmapi.PbMetaPageInfo(query.Result.Paging),
	}
	return nil
}

func (s *ShipmentPriceService) GetShopShipmentPriceList(ctx context.Context, r *GetShopShipmentPriceListEndpoint) error {
	query := &shopshipmentpricelist.GetShopShipmentPriceListQuery{
		ShopID: r.ShopID,
	}
	if err := s.ShopShipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShopShipmentPriceList(query.Result)
	return nil
}

func (s *ShipmentPriceService) CreateShopShipmentPriceList(ctx context.Context, r *CreateShopShipmentPriceListEndpoint) error {
	cmd := &shopshipmentpricelist.CreateShopShipmentPriceListCommand{
		ShopID:              r.ShopID,
		ShipmentPriceListID: r.ShipmentPriceListID,
		Note:                r.Note,
		UpdatedBy:           r.Context.UserID,
	}
	if err := s.ShopShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShopShipmentPriceList(cmd.Result)
	return nil
}

func (s *ShipmentPriceService) UpdateShopShipmentPriceList(ctx context.Context, r *UpdateShopShipmentPriceListEndpoint) error {
	cmd := &shopshipmentpricelist.UpdateShopShipmentPriceListCommand{
		ShopID:              r.ShopID,
		ShipmentPriceListID: r.ShipmentPriceListID,
		Note:                r.Note,
		UpdatedBy:           r.Context.UserID,
	}
	if err := s.ShopShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ShipmentPriceService) DeleteShopShipmentPriceList(ctx context.Context, r *DeleteShopShipmentPriceListEndpoint) error {
	cmd := &shopshipmentpricelist.DeleteShopShipmentPriceListCommand{
		ShopID: r.ShopID,
	}
	if err := s.ShopShipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

//-- End Shop Shipment Price List --//

func (s *ShipmentPriceService) GetShippingServices(ctx context.Context, r *GetShippingServicesEndpoint) error {
	if r.ShipmentPriceListID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn bảng giá áp dụng")
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
		return err
	}
	r.Result = &types.GetShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return nil
}
