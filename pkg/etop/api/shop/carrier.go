package shop

import (
	"context"

	pbcm "etop.vn/api/pb/common"
	pbshop "etop.vn/api/pb/etop/shop"
	"etop.vn/api/shopping/carrying"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/api/convertpb"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		carrierService.GetCarrier,
		carrierService.GetCarriers,
		carrierService.GetCarriersByIDs,
		carrierService.CreateCarrier,
		carrierService.UpdateCarrier,
		carrierService.DeleteCarrier)
}

func (s *CarrierService) GetCarrier(ctx context.Context, r *GetCarrierEndpoint) error {
	query := &carrying.GetCarrierByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := carrierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbCarrier(query.Result)
	return nil
}

func (s *CarrierService) GetCarriers(ctx context.Context, r *GetCarriersEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	query := &carrying.ListCarriersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := carrierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.CarriersResponse{
		Carriers: convertpb.PbCarriers(query.Result.Carriers),
		Paging:   cmapi.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}

func (s *CarrierService) GetCarriersByIDs(ctx context.Context, r *GetCarriersByIDsEndpoint) error {
	query := &carrying.ListCarriersByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := carrierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.CarriersResponse{
		Carriers: convertpb.PbCarriers(query.Result.Carriers),
	}
	return nil
}

func (s *CarrierService) CreateCarrier(ctx context.Context, r *CreateCarrierEndpoint) error {
	cmd := &carrying.CreateCarrierCommand{
		ShopID:   r.Context.Shop.ID,
		FullName: r.FullName,
		Note:     r.Note,
	}
	if err := carrierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCarrier(cmd.Result)
	return nil
}

func (s *CarrierService) UpdateCarrier(ctx context.Context, r *UpdateCarrierEndpoint) error {
	cmd := &carrying.UpdateCarrierCommand{
		ID:       r.Id,
		ShopID:   r.Context.Shop.ID,
		FullName: PString(r.FullName),
		Note:     PString(r.Note),
	}
	if err := carrierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCarrier(cmd.Result)
	return nil
}

func (s *CarrierService) DeleteCarrier(ctx context.Context, r *DeleteCarrierEndpoint) error {
	cmd := &carrying.DeleteCarrierCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := carrierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
	return nil
}
