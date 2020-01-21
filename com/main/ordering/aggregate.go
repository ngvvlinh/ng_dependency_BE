package ordering

import (
	"context"

	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/backend/com/main/ordering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

var _ ordering.Aggregate = &Aggregate{}

type Aggregate struct {
	store    sqlstore.OrderStoreFactory
	eventBus capi.EventBus
}

func NewAggregate(eventBus capi.EventBus, db *cmsql.Database) *Aggregate {
	return &Aggregate{
		store:    sqlstore.NewOrderStore(db),
		eventBus: eventBus,
	}
}

func (a *Aggregate) MessageBus() ordering.CommandBus {
	b := bus.New()
	return ordering.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) ValidateOrdersForShipping(ctx context.Context, args *ordering.ValidateOrdersForShippingArgs) (*ordering.ValidateOrdersForShippingResponse, error) {
	args1 := &ordering.GetOrdersArgs{
		IDs: args.OrderIDs,
	}
	orders, err := a.store(ctx).GetOrders(args1)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		if err := ValidateOrderStatus(order); err != nil {
			return nil, err
		}
	}
	return &ordering.ValidateOrdersForShippingResponse{}, nil
}

func ValidateOrderStatus(order *ordering.Order) error {
	switch order.Status {
	case status5.N:
		return cm.Error(cm.FailedPrecondition, "Đơn hàng đã huỷ.", nil).WithMetaID("orderID", order.ID)
	case status5.NS:
		return cm.Error(cm.FailedPrecondition, "Đơn hàng đã trả hàng.", nil).WithMetaID("orderID", order.ID)
	}

	if order.ConfirmStatus == status3.N {
		return cm.Error(cm.FailedPrecondition, "Đơn hàng đã huỷ.", nil).WithMetaID("orderID", order.ID)
	}
	return nil
}

func (a *Aggregate) ReserveOrdersForFfm(ctx context.Context, args *ordering.ReserveOrdersForFfmArgs) (*ordering.ReserveOrdersForFfmResponse, error) {
	orderIDs := args.OrderIDs
	orders, err := a.store(ctx).GetOrders(&ordering.GetOrdersArgs{
		IDs: orderIDs,
	})
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		if err := ValidateOrderForReserveFfm(order, args.Fulfill); err != nil {
			return nil, err
		}
	}

	update := sqlstore.UpdateOrdersForReserveOrdersFfmArgs{
		OrderIDs:   orderIDs,
		Fulfill:    args.Fulfill,
		FulfillIDs: args.FulfillIDs,
	}
	orders, err = a.store(ctx).UpdateOrdersForReserveOrdersFfm(update)
	if err != nil {
		return nil, err
	}
	return &ordering.ReserveOrdersForFfmResponse{
		Orders: orders,
	}, nil
}

func ValidateOrderForReserveFfm(order *ordering.Order, shippingType ordertypes.ShippingType) error {
	if err := ValidateOrderStatus(order); err != nil {
		return err
	}
	if shippingType != ordertypes.ShippingTypeShipnow {
		return nil
	}
	if len(order.FulfillmentIDs) != 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Order has been reserved").WithMetaID("order_id", order.ID)
	}
	return nil
}

func (a *Aggregate) ReleaseOrdersForFfm(ctx context.Context, args *ordering.ReleaseOrdersForFfmArgs) (*ordering.ReleaseOrdersForFfmResponse, error) {
	orderIDs := args.OrderIDs
	if len(orderIDs) == 0 {
		return &ordering.ReleaseOrdersForFfmResponse{
			Updated: 0,
		}, nil
	}
	update := sqlstore.UpdateOrdersForReleaseOrderFfmArgs{
		OrderIDs: orderIDs,
	}
	if err := a.store(ctx).UpdateOrdersForReleaseOrdersFfm(update); err != nil {
		return nil, err
	}
	return &ordering.ReleaseOrdersForFfmResponse{
		Updated: len(orderIDs),
	}, nil
}

func (a *Aggregate) UpdateOrderShippingStatus(ctx context.Context, args *ordering.UpdateOrderShippingStatusArgs) (*ordering.UpdateOrderShippingStatusResponse, error) {
	update := sqlstore.UpdateOrderShippingStatusArgs{
		ID:                         args.ID,
		FulfillmentShippingStates:  args.FulfillmentShippingStates,
		FulfillmentShippingStatus:  args.FulfillmentShippingStatus,
		FulfillmentPaymentStatuses: args.FulfillmentPaymentStatuses,
		EtopPaymentStatus:          args.EtopPaymentStatus,
	}
	err := a.store(ctx).UpdateOrderShippingStatus(update)
	return &ordering.UpdateOrderShippingStatusResponse{}, err
}

func (a *Aggregate) UpdateOrdersConfirmStatus(ctx context.Context, args *ordering.UpdateOrdersConfirmStatusArgs) (*ordering.UpdateOrdersConfirmStatusResponse, error) {
	update := sqlstore.UpdateOrdersConfirmStatusArgs{
		IDs:           args.IDs,
		ShopConfirm:   args.ShopConfirm,
		ConfirmStatus: args.ConfirmStatus,
	}
	err := a.store(ctx).UpdateOrdersConfirmStatus(update)
	return &ordering.UpdateOrdersConfirmStatusResponse{}, err
}

func (a *Aggregate) UpdateOrderPaymentInfo(ctx context.Context, args *ordering.UpdateOrderPaymentInfoArgs) error {
	update := sqlstore.UpdateOrderPaymentInfoArgs{
		ID:            args.ID,
		PaymentStatus: args.PaymentStatus,
		PaymentID:     args.PaymentID,
	}
	if err := a.store(ctx).UpdateOrderPaymentInfo(update); err != nil {
		return err
	}

	event := &ordering.OrderPaymentSuccessEvent{
		EventMeta: meta.NewEvent(),
		OrderID:   args.ID,
	}
	// ignore err
	_ = a.eventBus.Publish(ctx, event)
	return nil
}

func (a *Aggregate) CompleteOrder(ctx context.Context, orderID dot.ID, shopID dot.ID) error {
	update := &sqlstore.UpdateOrderStatus{
		ID:     orderID,
		ShopID: shopID,
		Status: status5.P,
	}
	return a.store(ctx).UpdateOrderStatus(update)
}

func (a *Aggregate) UpdateOrderStatus(ctx context.Context, args *ordering.UpdateOrderStatusArgs) error {
	update := &sqlstore.UpdateOrderStatus{
		ID:     args.OrderID,
		ShopID: args.ShopID,
		Status: args.Status,
	}
	return a.store(ctx).UpdateOrderStatus(update)
}
