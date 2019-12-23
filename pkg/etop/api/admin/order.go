package admin

import (
	"context"

	"etop.vn/api/top/int/types"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	shipmodelx "etop.vn/backend/com/main/shipping/modelx"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandler("api", orderService.GetOrder)
	bus.AddHandler("api", orderService.GetOrdersByIDs)
	bus.AddHandler("api", orderService.GetOrders)
	bus.AddHandler("api", fulfillmentService.GetFulfillment)
	bus.AddHandler("api", fulfillmentService.GetFulfillments)
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

func (s *FulfillmentService) GetFulfillment(ctx context.Context, q *GetFulfillmentEndpoint) error {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbFulfillment(query.Result.Fulfillment, model.TagEtop, query.Result.Shop, query.Result.Order)
	return nil
}

func (s *FulfillmentService) GetFulfillments(ctx context.Context, q *GetFulfillmentsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		OrderID: q.OrderId,
		Status:  q.Status,
		Paging:  paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if q.ShopId != 0 {
		query.ShopIDs = []dot.ID{q.ShopId}
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.FulfillmentsResponse{
		Fulfillments: convertpb.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagEtop),
		Paging:       cmapi.PbPageInfo(paging),
	}
	return nil
}
