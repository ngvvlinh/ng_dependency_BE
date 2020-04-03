package admin

import (
	"context"

	"etop.vn/api/main/shipmentpricing/pricelist"
	"etop.vn/api/main/shipmentpricing/shipmentprice"
	"etop.vn/api/main/shipmentpricing/shipmentservice"
	"etop.vn/api/top/int/admin"
	"etop.vn/api/top/int/types"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/backend/com/main/shipping/carrier"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/api/convertpb"
)

//-- ShipmentService --//

func (s *ShipmentPriceService) GetShipmentService(ctx context.Context, r *GetShipmentServiceEndpoint) error {
	query := &shipmentservice.GetShipmentServiceQuery{
		ID: r.Id,
	}
	if err := shipmentServiceQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentService(query.Result)
	return nil
}

func (s *ShipmentPriceService) GetShipmentServices(ctx context.Context, r *GetShipmentServicesEndpoint) error {
	query := &shipmentservice.ListShipmentServicesQuery{}
	if err := shipmentServiceQuery.Dispatch(ctx, query); err != nil {
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
	if err := shipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := shipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ShipmentPriceService) DeleteShipmentService(ctx context.Context, r *DeleteShipmentServiceEndpoint) error {
	cmd := &shipmentservice.DeleteShipmentServiceCommand{
		ID: r.Id,
	}
	if err := shipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := shipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := shipmentServiceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

//-- End ShipmentService --//

//-- ShipmentPriceList --//

func (s *ShipmentPriceService) GetShipmentPriceList(ctx context.Context, r *GetShipmentPriceListEndpoint) error {
	query := &pricelist.GetShipmentPriceListQuery{
		ID: r.Id,
	}
	if err := shipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPriceList(query.Result)
	return nil
}

func (s *ShipmentPriceService) GetShipmentPriceLists(ctx context.Context, r *GetShipmentPriceListsEndpoint) error {
	query := &pricelist.ListShipmentPriceListQuery{}
	if err := shipmentPriceListQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &admin.GetShipmentPriceListsResponse{
		ShipmentPriceLists: convertpb.PbShipmentPriceLists(query.Result),
	}
	return nil
}

func (s *ShipmentPriceService) CreateShipmentPriceList(ctx context.Context, r *CreateShipmentPriceListEndpoint) error {
	cmd := &pricelist.CreateShipmentPriceListCommand{
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
	if err := shipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPriceList(cmd.Result)
	return nil
}

func (s *ShipmentPriceService) UpdateShipmentPriceList(ctx context.Context, r *UpdateShipmentPriceListEndpoint) error {
	cmd := &pricelist.UpdateShipmentPriceListCommand{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
	if err := shipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ShipmentPriceService) ActivateShipmentPriceList(ctx context.Context, r *ActivateShipmentPriceListEndpoint) error {
	cmd := &pricelist.ActivateShipmentPriceListCommand{
		ID: r.Id,
	}
	if err := shipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ShipmentPriceService) DeleteShipmentPriceList(ctx context.Context, r *DeleteShipmentPriceListEndpoint) error {
	cmd := &pricelist.DeleteShipmentPriceListCommand{
		ID: r.Id,
	}
	if err := shipmentPriceListAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := shipmentPriceQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPrice(query.Result)
	return nil
}

func (s *ShipmentPriceService) GetShipmentPrices(ctx context.Context, r *GetShipmentPricesEndpoint) error {
	query := &shipmentprice.ListShipmentPricesQuery{
		ShipmentPriceListID: r.ShipmentPriceListID,
		ShipmentServiceID:   r.ShipmentServiceID,
	}
	if err := shipmentPriceQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &admin.GetShipmentPricesResponse{
		ShipmentPrices: convertpb.PbShipmentPrices(query.Result),
	}
	return nil
}

func (s *ShipmentPriceService) CreateShipmentPrice(ctx context.Context, r *CreateShipmentPriceEndpoint) error {
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
		Details:             convertpb.PricingDetails(r.Details),
	}
	if err := shipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPrice(cmd.Result)
	return nil
}

func (s *ShipmentPriceService) UpdateShipmentPrice(ctx context.Context, r *UpdateShipmentPriceEndpoint) error {
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
		Details:             convertpb.PricingDetails(r.Details),
		Status:              r.Status,
	}
	if err := shipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentPrice(cmd.Result)
	return nil
}

func (s *ShipmentPriceService) DeleteShipmentPrice(ctx context.Context, r *DeleteShipmentPriceEndpoint) error {
	cmd := &shipmentprice.DeleteShipmentPriceCommand{
		ID: r.Id,
	}
	if err := shipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := shipmentPriceAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

//-- End ShipmentPrice --//

func (s *ShipmentPriceService) GetShippingServices(ctx context.Context, r *GetShippingServicesEndpoint) error {
	if r.ShipmentPriceListID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn bảng giá áp dụng")
	}
	args := &carrier.GetShippingServicesArgs{
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
	resp, err := shipmentManager.GetShippingServices(ctx, 0, args)
	if err != nil {
		return err
	}
	r.Result = &types.GetShippingServicesResponse{
		Services: convertpb.PbAvailableShippingServices(resp),
	}
	return nil
}
