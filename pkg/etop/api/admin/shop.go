package admin

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/filter"
)

type ShopService struct {
	session.Session

	IdentityQuery identity.QueryBus
}

func (s *ShopService) Clone() admin.ShopService {
	res := *s
	return &res
}

func (s *ShopService) GetShop(ctx context.Context, q *pbcm.IDRequest) (*etop.Shop, error) {
	query := &identitymodelx.GetShopExtendedQuery{
		ShopID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbShopExtended(query.Result)
	return result, nil
}

func (s *ShopService) GetShops(ctx context.Context, q *admin.GetShopsRequest) (*admin.GetShopsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	var fullTextSearch filter.FullTextSearch = ""
	if q.Filter != nil {
		fullTextSearch = q.Filter.Name
	}
	query := &identity.ListShopExtendedsQuery{
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
		Name:    fullTextSearch,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := &admin.GetShopsResponse{
		Paging: cmapi.PbPageInfo(paging),
		Shops:  convertpb.Convert_core_ShopExtendeds_To_api_ShopExtendeds(query.Result.Shops),
	}
	return result, nil
}

func (s *ShopService) GetShopsByIDs(ctx context.Context, q *pbcm.IDsRequest) (*admin.GetShopsResponse, error) {
	query := &identity.ListShopsByIDsQuery{
		IDs: q.Ids,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &admin.GetShopsResponse{
		Shops: convertpb.Convert_core_Shops_To_api_Shops(query.Result),
	}
	return result, nil
}
