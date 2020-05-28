package admin

import (
	"context"

	"o.o/api/top/int/types"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/model"
)

type OrderService struct{}

func (s *OrderService) Clone() *OrderService {
	res := *s
	return &res
}

func (s *OrderService) GetOrder(ctx context.Context, q *GetOrderEndpoint) error {
	query := &ordermodelx.GetOrderQuery{
		OrderID:            q.Id,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = convertpb.PbOrder(
		query.Result.Order,
		query.Result.Fulfillments,
		model.TagEtop,
	)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, q *GetOrdersEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &ordermodelx.GetOrdersQuery{
		Paging:  paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.OrdersResponse{
		Paging: cmapi.PbPageInfo(paging),
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, model.TagEtop, query.Result.Shops),
	}
	return nil
}

func (s *OrderService) GetOrdersByIDs(ctx context.Context, q *GetOrdersByIDsEndpoint) error {
	query := &ordermodelx.GetOrdersQuery{
		IDs: q.Ids,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.OrdersResponse{
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, model.TagEtop, query.Result.Shops),
	}
	return nil
}
