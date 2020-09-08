package admin

import (
	"context"

	"o.o/api/top/int/admin"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_tag"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
)

type OrderService struct {
	session.Session

	OrderStore sqlstore.OrderStoreInterface
}

func (s *OrderService) Clone() admin.OrderService {
	res := *s
	return &res
}

func (s *OrderService) GetOrder(ctx context.Context, q *pbcm.IDRequest) (*types.Order, error) {
	query := &ordermodelx.GetOrderQuery{
		OrderID:            q.Id,
		IncludeFulfillment: true,
	}
	if err := s.OrderStore.GetOrder(ctx, query); err != nil {
		return nil, err
	}

	result := convertpb.PbOrder(
		query.Result.Order,
		query.Result.Fulfillments,
		account_tag.TagEtop,
	)
	return result, nil
}

func (s *OrderService) GetOrders(ctx context.Context, q *admin.GetOrdersRequest) (*types.OrdersResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &ordermodelx.GetOrdersQuery{
		Paging:  paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.OrderStore.GetOrders(ctx, query); err != nil {
		return nil, err
	}
	result := &types.OrdersResponse{
		Paging: cmapi.PbPageInfo(paging),
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, account_tag.TagEtop, query.Result.Shops),
	}
	return result, nil
}

func (s *OrderService) GetOrdersByIDs(ctx context.Context, q *pbcm.IDsRequest) (*types.OrdersResponse, error) {
	query := &ordermodelx.GetOrdersQuery{
		IDs: q.Ids,
	}
	if err := s.OrderStore.GetOrders(ctx, query); err != nil {
		return nil, err
	}
	result := &types.OrdersResponse{
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, account_tag.TagEtop, query.Result.Shops),
	}
	return result, nil
}
