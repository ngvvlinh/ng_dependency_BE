package admin

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/top/int/admin"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
)

type ShopService struct {
	IdentityQuery identity.QueryBus
}

func (s *ShopService) Clone() *ShopService {
	res := *s
	return &res
}

func (s *ShopService) GetShop(ctx context.Context, q *GetShopEndpoint) error {
	query := &identitymodelx.GetShopExtendedQuery{
		ShopID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbShopExtended(query.Result)
	return nil
}

func (s *ShopService) GetShops(ctx context.Context, q *GetShopsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &identity.ListShopExtendedsQuery{
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &admin.GetShopsResponse{
		Paging: cmapi.PbPageInfo(paging),
		Shops:  convertpb.Convert_core_ShopExtendeds_To_api_ShopExtendeds(query.Result.Shops),
	}
	return nil
}

func (s *ShopService) GetShopsByIDs(ctx context.Context, q *GetShopsByIDsEndpoint) error {
	query := &identity.ListShopsByIDsQuery{
		IDs: q.Ids,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &admin.GetShopsResponse{
		Shops: convertpb.Convert_core_Shops_To_api_Shops(query.Result),
	}
	return nil
}
