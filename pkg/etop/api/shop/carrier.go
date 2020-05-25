package shop

import (
	"context"

	"o.o/api/shopping/carrying"
	"o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

type CarrierService struct {
	CarrierAggr  carrying.CommandBus
	CarrierQuery carrying.QueryBus
}

func (s *CarrierService) Clone() *CarrierService { res := *s; return &res }

func (s *CarrierService) GetCarrier(ctx context.Context, r *GetCarrierEndpoint) error {
	query := &carrying.GetCarrierByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := s.CarrierQuery.Dispatch(ctx, query); err != nil {
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
	if err := s.CarrierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &shop.CarriersResponse{
		Carriers: convertpb.PbCarriers(query.Result.Carriers),
		Paging:   cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *CarrierService) GetCarriersByIDs(ctx context.Context, r *GetCarriersByIDsEndpoint) error {
	query := &carrying.ListCarriersByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := s.CarrierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &shop.CarriersResponse{
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
	if err := s.CarrierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCarrier(cmd.Result)
	return nil
}

func (s *CarrierService) UpdateCarrier(ctx context.Context, r *UpdateCarrierEndpoint) error {
	cmd := &carrying.UpdateCarrierCommand{
		ID:       r.Id,
		ShopID:   r.Context.Shop.ID,
		FullName: r.FullName,
		Note:     r.Note,
	}
	if err := s.CarrierAggr.Dispatch(ctx, cmd); err != nil {
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
	if err := s.CarrierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: cmd.Result}
	return nil
}
