package ordering

import (
	"context"

	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/ordering/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ ordering.Aggregate = &Aggregate{}

type Aggregate struct {
	store    sqlstore.OrderStoreFactory
	eventBus capi.EventBus
}

func NewAggregate(eventBus capi.EventBus, db com.MainDB) *Aggregate {
	return &Aggregate{
		store:    sqlstore.NewOrderStore(db),
		eventBus: eventBus,
	}
}

func AggregateMessageBus(a *Aggregate) ordering.CommandBus {
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
		FulfillmentStatuses:        args.FulfillmentStatuses,
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
	// Kiểm tra lại hàm này
	// Chỉ bắn event khi đây là đơn trading
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

func (a *Aggregate) UpdateOrderPaymentStatus(ctx context.Context, args *ordering.UpdateOrderPaymentStatusArgs) error {
	if args.OrderID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing OrderID", nil)
	}
	if !args.PaymentStatus.Valid {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing payment status")
	}
	if _, err := a.store(ctx).ID(args.OrderID).OptionalShopID(args.ShopID).GetOrder(); err != nil {
		return err
	}

	return a.store(ctx).UpdateOrderPaymentStatus(args)
}
