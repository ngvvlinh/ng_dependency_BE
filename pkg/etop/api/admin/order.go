package admin

import (
	"context"

	"etop.vn/backend/pkg/services/shipping/modelx"

	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"

	pbcm "etop.vn/backend/pb/common"
	pborder "etop.vn/backend/pb/etop/order"
	wrapadmin "etop.vn/backend/wrapper/etop/admin"
)

func init() {
	bus.AddHandler("api", GetOrder)
	bus.AddHandler("api", GetOrdersByIDs)
	bus.AddHandler("api", GetOrders)
	bus.AddHandler("api", GetFulfillment)
	bus.AddHandler("api", GetFulfillments)
}

func GetOrder(ctx context.Context, q *wrapadmin.GetOrderEndpoint) error {
	query := &model.GetOrderQuery{
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

func GetOrders(ctx context.Context, q *wrapadmin.GetOrdersEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &model.GetOrdersQuery{
		Paging:  paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.OrdersResponse{
		Paging: pbcm.PbPageInfo(paging, query.Result.Total),
		Orders: pborder.PbOrders(query.Result.Orders, model.TagEtop, query.Result.Shops),
	}
	return nil
}

func GetOrdersByIDs(ctx context.Context, q *wrapadmin.GetOrdersByIDsEndpoint) error {
	query := &model.GetOrdersQuery{
		IDs: q.Ids,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.OrdersResponse{
		Orders: pborder.PbOrders(query.Result.Orders, model.TagEtop, query.Result.Shops),
	}
	return nil
}

func GetFulfillment(ctx context.Context, q *wrapadmin.GetFulfillmentEndpoint) error {
	query := &modelx.GetFulfillmentExtendedQuery{
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbFulfillment(query.Result.Fulfillment, model.TagEtop, query.Result.Shop, query.Result.Order)
	return nil
}

func GetFulfillments(ctx context.Context, q *wrapadmin.GetFulfillmentsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &modelx.GetFulfillmentExtendedsQuery{
		SupplierID: q.SupplierId,
		OrderID:    q.OrderId,
		Status:     q.Status.ToModel(),
		Paging:     paging,
		Filters:    pbcm.ToFilters(q.Filters),
	}
	if q.ShopId != 0 {
		query.ShopIDs = []int64{q.ShopId}
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.FulfillmentsResponse{
		Fulfillments: pborder.PbFulfillmentExtendeds(query.Result.Fulfillments, model.TagEtop),
		Paging:       pbcm.PbPageInfo(paging, query.Result.Total),
	}
	return nil
}
