package shop

import (
	"context"

	"etop.vn/api/shopping/carrying"

	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		GetCarrier,
		GetCarriers,
		GetCarriersByIDs,
		CreateCarrier,
		UpdateCarrier,
		DeleteCarrier)
}

func GetCarrier(ctx context.Context, r *wrapshop.GetCarrierEndpoint) error {
	query := &carrying.GetCarrierByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := carrierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = pbshop.PbCarrier(query.Result)
	return nil
}

func GetCarriers(ctx context.Context, r *wrapshop.GetCarriersEndpoint) error {
	paging := r.Paging.CMPaging()
	query := &carrying.ListCarriersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(r.Filters),
	}
	if err := carrierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.CarriersResponse{
		Carriers: pbshop.PbCarriers(query.Result.Carriers),
		Paging:   pbcm.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}

func GetCarriersByIDs(ctx context.Context, r *wrapshop.GetCarriersByIDsEndpoint) error {
	query := &carrying.ListCarriersByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := carrierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.CarriersResponse{
		Carriers: pbshop.PbCarriers(query.Result.Carriers),
	}
	return nil
}

func CreateCarrier(ctx context.Context, r *wrapshop.CreateCarrierEndpoint) error {
	cmd := &carrying.CreateCarrierCommand{
		ShopID:   r.Context.Shop.ID,
		FullName: r.FullName,
		Note:     r.Note,
	}
	if err := carrierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbCarrier(cmd.Result)
	return nil
}

func UpdateCarrier(ctx context.Context, r *wrapshop.UpdateCarrierEndpoint) error {
	cmd := &carrying.UpdateCarrierCommand{
		ID:       r.Id,
		ShopID:   r.Context.Shop.ID,
		FullName: PString(r.FullName),
		Note:     PString(r.Note),
	}
	if err := carrierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbCarrier(cmd.Result)
	return nil
}

func DeleteCarrier(ctx context.Context, r *wrapshop.DeleteCarrierEndpoint) error {
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
