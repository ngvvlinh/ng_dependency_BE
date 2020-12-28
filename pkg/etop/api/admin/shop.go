package admin

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/main/moneytx"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

type ShopService struct {
	session.Session

	IdentityQuery identity.QueryBus
	IdentityAggr  identity.CommandBus
	ShopStore     sqlstore.ShopStoreInterface
	MoneyTxQuery  moneytx.QueryBus
}

func (s *ShopService) Clone() admin.ShopService {
	res := *s
	return &res
}

func (s *ShopService) GetShop(ctx context.Context, q *pbcm.IDRequest) (*etop.Shop, error) {
	query := &identitymodelx.GetShopExtendedQuery{
		ShopID: q.Id,
	}
	if err := s.ShopStore.GetShopExtended(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbShopExtended(query.Result)
	return result, nil
}

func (s *ShopService) GetShops(ctx context.Context, q *admin.GetShopsRequest) (*admin.GetShopsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	var fullTextSearch filter.FullTextSearch = ""
	var shopIDs []dot.ID
	if q.Filter != nil {
		fullTextSearch = q.Filter.Name
		shopIDs = q.Filter.ShopIDs
	}
	query := &identity.ListShopExtendedsQuery{
		Paging:               *paging,
		Filters:              cmapi.ToFilters(q.Filters),
		Name:                 fullTextSearch,
		ShopIDs:              shopIDs,
		IncludeWLPartnerShop: true,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	// count total money transaction shipping per shop
	shops := query.Result.Shops
	_shopIDs := []dot.ID{}
	for _, shop := range shops {
		_shopIDs = append(_shopIDs, shop.ID)
	}

	queryShopMoneyTxCount := &moneytx.CountMoneyTxShippingByShopIDsQuery{
		ShopIDs: _shopIDs,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, queryShopMoneyTxCount); err != nil {
		return nil, err
	}

	result := &admin.GetShopsResponse{
		Paging: cmapi.PbPageInfo(paging),
		Shops:  convertpball.Convert_core_ShopExtendeds_To_api_ShopExtendeds(shops, queryShopMoneyTxCount.Result),
	}
	return result, nil
}

func (s *ShopService) GetShopsByIDs(ctx context.Context, q *pbcm.IDsRequest) (*admin.GetShopsResponse, error) {
	query := &identity.ListShopsByIDsQuery{
		IDs:                  q.Ids,
		IncludeWLPartnerShop: true,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	// count total money transaction shipping per shop
	queryShopMoneyTxCount := &moneytx.CountMoneyTxShippingByShopIDsQuery{
		ShopIDs: q.Ids,
	}
	if err := s.MoneyTxQuery.Dispatch(ctx, queryShopMoneyTxCount); err != nil {
		return nil, err
	}
	result := &admin.GetShopsResponse{
		Shops: convertpball.Convert_core_Shops_To_api_Shops(query.Result, queryShopMoneyTxCount.Result),
	}
	return result, nil
}

func (s *ShopService) UpdateShopInfo(ctx context.Context, r *admin.UpdateShopInfoRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &identity.UpdateShopInfoCommand{
		ShopID:                  r.ID,
		MoneyTransactionRrule:   r.MoneyTransactionRrule,
		IsPriorMoneyTransaction: r.IsPriorMoneyTransaction,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}
