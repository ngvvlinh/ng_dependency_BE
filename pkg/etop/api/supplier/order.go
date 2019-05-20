package supplier

import (
	"context"

	modelx2 "etop.vn/backend/pkg/services/selling/modelx"

	"etop.vn/backend/pkg/services/shipping/modelx"

	cmP "etop.vn/backend/pb/common"
	orderP "etop.vn/backend/pb/etop/order"
	supplierP "etop.vn/backend/pb/etop/supplier"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	supplierW "etop.vn/backend/wrapper/etop/supplier"
)

func init() {
	bus.AddHandler("api", GetOrder)
	bus.AddHandler("api", GetOrdersByIDs)
	bus.AddHandler("api", GetOrders)
	bus.AddHandler("api", UpdateOrdersStatus)
	bus.AddHandler("api", GetFulfillment)
	bus.AddHandler("api", GetFulfillments)
}

func GetOrder(ctx context.Context, q *supplierW.GetOrderEndpoint) error {
	query := &modelx2.GetOrderQuery{
		OrderID:    q.Id,
		SupplierID: q.Context.Supplier.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	order := query.Result.Order

	shopQuery := &model.GetShopExtendedQuery{
		ShopID: order.ShopID,
	}
	if err := bus.Dispatch(ctx, shopQuery); err != nil {
		return err
	}

	// TODO: Fulfillment
	q.Result = supplierP.PbOrder(order, shopQuery.Result)
	return nil
}

func GetOrdersByIDs(ctx context.Context, q *supplierW.GetOrdersByIDsEndpoint) error {
	query := &modelx2.GetOrdersQuery{
		SupplierID: q.Context.Supplier.ID,
		IDs:        q.Ids,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &supplierP.OrdersResponse{
		// Orders: supplierP.PbOrders(query.Result.Orders),
	}
	return nil
}

func GetOrders(ctx context.Context, q *supplierW.GetOrdersEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &modelx2.GetOrdersQuery{
		SupplierID: q.Context.Supplier.ID,
		Paging:     paging,
		Filters:    cmP.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &supplierP.OrdersResponse{
		Paging: cmP.PbPageInfo(paging, query.Result.Total),
		// Orders: supplierP.PbOrders(query.Result.Orders),
	}
	return nil
}

func UpdateOrdersStatus(ctx context.Context, q *supplierW.UpdateOrdersStatusEndpoint) error {
	if len(q.Updates) == 0 {
		return cm.Error(cm.InvalidArgument, "Nothing to update", nil)
	}

	updates := make([]model.UpdateOrderLinesStatus, len(q.Updates))
	for i, update := range q.Updates {
		if update == nil {
			return cm.Error(cm.InvalidArgument, "Empty update", nil)
		}
		updates[i] = model.UpdateOrderLinesStatus{
			OrderID:         update.OrderId,
			ProductIDs:      update.ProductIds,
			SupplierConfirm: update.SConfirm.ToModel(),
			CancelReason:    update.CancelReason,
		}
	}

	cmd := &model.UpdateOrderLinesStatusCommand{
		SupplierID: q.Context.Supplier.ID,
		Updates:    updates,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.UpdatedResponse{Updated: int32(cmd.Result.Updated)}
	return nil
}

func GetFulfillment(ctx context.Context, q *supplierW.GetFulfillmentEndpoint) error {
	query := &modelx.GetFulfillmentQuery{
		SupplierID:    q.Context.Supplier.ID,
		FulfillmentID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = orderP.PbFulfillment(query.Result, model.TagShop, nil, nil)
	return nil
}

func GetFulfillments(ctx context.Context, q *supplierW.GetFulfillmentsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &modelx.GetFulfillmentsQuery{
		SupplierID: q.Context.Supplier.ID,
		OrderID:    q.OrderId,
		Status:     q.Status.ToModel(),
		Paging:     paging,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &orderP.FulfillmentsResponse{
		Fulfillments: orderP.PbFulfillments(query.Result.Fulfillments, model.TagShop),
		Paging:       cmP.PbPageInfo(paging, query.Result.Total),
	}
	return nil
}
