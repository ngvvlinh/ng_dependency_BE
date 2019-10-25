package admin

import (
	"context"

	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	shipmodelx "etop.vn/backend/com/main/shipping/modelx"
	pbcm "etop.vn/backend/pb/common"
	pborder "etop.vn/backend/pb/etop/order"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	wrapadmin "etop.vn/backend/wrapper/etop/admin"
)

func init() {
	bus.AddHandler("api", s.GetOrder)
	bus.AddHandler("api", s.GetOrdersByIDs)
	bus.AddHandler("api", s.GetOrders)
	bus.AddHandler("api", s.GetFulfillment)
	bus.AddHandler("api", s.GetFulfillments)
}

func (s *Service) GetOrder(ctx context.Context, q *wrapadmin.GetOrderEndpoint) error {
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

func (s *Service) GetOrders(ctx context.Context, q *wrapadmin.GetOrdersEndpoint) error {
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

func (s *Service) GetOrdersByIDs(ctx context.Context, q *wrapadmin.GetOrdersByIDsEndpoint) error {
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

func (s *Service) GetFulfillment(ctx context.Context, q *wrapadmin.GetFulfillmentEndpoint) error {
	query := &shipmodelx.GetFulfillmentExtendedQuery{
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbFulfillment(query.Result.Fulfillment, model.TagEtop, query.Result.Shop, query.Result.Order)
	return nil
}

func (s *Service) GetFulfillments(ctx context.Context, q *wrapadmin.GetFulfillmentsEndpoint) error {
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
