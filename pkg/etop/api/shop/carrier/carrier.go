package carrier

import (
	"context"

	"o.o/api/shopping/carrying"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/authorize/session"
)

type CarrierService struct {
	session.Session

	CarrierAggr  carrying.CommandBus
	CarrierQuery carrying.QueryBus
}

func (s *CarrierService) Clone() api.CarrierService { res := *s; return &res }

func (s *CarrierService) GetCarrier(ctx context.Context, r *pbcm.IDRequest) (*api.Carrier, error) {
	query := &carrying.GetCarrierByIDQuery{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CarrierQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpball.PbCarrier(query.Result)
	return result, nil
}

func (s *CarrierService) GetCarriers(ctx context.Context, r *api.GetCarriersRequest) (*api.CarriersResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	query := &carrying.ListCarriersQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := s.CarrierQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.CarriersResponse{
		Carriers: convertpball.PbCarriers(query.Result.Carriers),
		Paging:   cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *CarrierService) GetCarriersByIDs(ctx context.Context, r *pbcm.IDsRequest) (*api.CarriersResponse, error) {
	query := &carrying.ListCarriersByIDsQuery{
		IDs:    r.Ids,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CarrierQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.CarriersResponse{
		Carriers: convertpball.PbCarriers(query.Result.Carriers),
	}
	return result, nil
}

func (s *CarrierService) CreateCarrier(ctx context.Context, r *api.CreateCarrierRequest) (*api.Carrier, error) {
	cmd := &carrying.CreateCarrierCommand{
		ShopID:   s.SS.Shop().ID,
		FullName: r.FullName,
		Note:     r.Note,
	}
	if err := s.CarrierAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpball.PbCarrier(cmd.Result)
	return result, nil
}

func (s *CarrierService) UpdateCarrier(ctx context.Context, r *api.UpdateCarrierRequest) (*api.Carrier, error) {
	cmd := &carrying.UpdateCarrierCommand{
		ID:       r.Id,
		ShopID:   s.SS.Shop().ID,
		FullName: r.FullName,
		Note:     r.Note,
	}
	if err := s.CarrierAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpball.PbCarrier(cmd.Result)
	return result, nil
}

func (s *CarrierService) DeleteCarrier(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &carrying.DeleteCarrierCommand{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CarrierAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: cmd.Result}
	return result, nil
}
