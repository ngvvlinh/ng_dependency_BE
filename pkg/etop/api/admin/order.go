package admin

import (
	"context"

	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	shipmodelx "etop.vn/backend/com/main/shipping/modelx"
	pbcm "etop.vn/backend/pb/common"
	pborder "etop.vn/backend/pb/etop/order"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
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

	q.Result = pborder.PbOrder(
		query.Result.Order,
		query.Result.Fulfillments,
		model.TagEtop,
	)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, q *GetOrdersEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &ordermodelx.GetOrdersQuery{
		Paging:  paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.OrdersResponse{
		Paging: pbcm.PbPageInfo(paging, int32(query.Result.Total)),
		Orders: pborder.PbOrdersWithFulfillments(query.Result.Orders, model.TagEtop, query.Result.Shops),
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
	q.Result = &pborder.OrdersResponse{
		Orders: pborder.PbOrdersWithFulfillments(query.Result.Orders, model.TagEtop, query.Result.Shops),
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
	q.Result = pborder.PbFulfillment(query.Result.Fulfillment, model.TagEtop, query.Result.Shop, query.Result.Order)
	return nil
}

func (s *FulfillmentService) GetFulfillments(ctx context.Context, q *GetFulfillmentsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &shipmodelx.GetFulfillmentExtendedsQuery{
		OrderID: q.OrderId,
		Status:  q.Status.ToModel(),
		Paging:  paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if q.ShopId != 0 {
		query.ShopIDs = []int64{q.ShopId}
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.FulfillmentsResponse{
		Fulfillments: pborder.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagEtop),
		Paging:       pbcm.PbPageInfo(paging, int32(query.Result.Total)),
	}
	return nil
}
