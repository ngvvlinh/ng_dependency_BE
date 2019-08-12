// Code generated by gen-cmd-query. DO NOT EDIT.

package ordering

import (
	context "context"
	time "time"

	etopv1 "etop.vn/api/main/etop/v1"
	orderingv1types "etop.vn/api/main/ordering/v1/types"
	meta "etop.vn/api/meta"
)

type Command interface{ command() }
type Query interface{ query() }
type CommandBus struct{ bus meta.Bus }
type QueryBus struct{ bus meta.Bus }

func (c CommandBus) Dispatch(ctx context.Context, msg Command) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c CommandBus) DispatchAll(ctx context.Context, msgs ...Command) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
func (c QueryBus) DispatchAll(ctx context.Context, msgs ...Query) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

type ReleaseOrdersForFfmCommand struct {
	OrderIDs []int64

	Result *ReleaseOrdersForFfmResponse `json:"-"`
}

type ReserveOrdersForFfmCommand struct {
	OrderIDs   []int64
	Fulfill    orderingv1types.Fulfill
	FulfillIDs []int64

	Result *ReserveOrdersForFfmResponse `json:"-"`
}

type UpdateOrderShippingStatusCommand struct {
	ID                         int64
	FulfillmentShippingStates  []string
	FulfillmentShippingStatus  etopv1.Status5
	FulfillmentPaymentStatuses []int
	EtopPaymentStatus          etopv1.Status4
	CODEtopPaidAt              time.Time

	Result *UpdateOrderShippingStatusResponse `json:"-"`
}

type UpdateOrdersConfirmStatusCommand struct {
	IDs           []int64
	ShopConfirm   etopv1.Status3
	ConfirmStatus etopv1.Status3

	Result *UpdateOrdersConfirmStatusResponse `json:"-"`
}

type ValidateOrdersForShippingCommand struct {
	OrderIDs []int64

	Result *ValidateOrdersForShippingResponse `json:"-"`
}

type GetOrderByCodeQuery struct {
	Code string

	Result *Order `json:"-"`
}

type GetOrderByIDQuery struct {
	ID int64

	Result *Order `json:"-"`
}

type GetOrdersQuery struct {
	ShopID int64
	IDs    []int64

	Result *OrdersResponse `json:"-"`
}

// implement interfaces

func (q *ReleaseOrdersForFfmCommand) command()       {}
func (q *ReserveOrdersForFfmCommand) command()       {}
func (q *UpdateOrderShippingStatusCommand) command() {}
func (q *UpdateOrdersConfirmStatusCommand) command() {}
func (q *ValidateOrdersForShippingCommand) command() {}
func (q *GetOrderByCodeQuery) query()                {}
func (q *GetOrderByIDQuery) query()                  {}
func (q *GetOrdersQuery) query()                     {}

// implement conversion

func (q *ReleaseOrdersForFfmCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ReleaseOrdersForFfmArgs) {
	return ctx,
		&ReleaseOrdersForFfmArgs{
			OrderIDs: q.OrderIDs,
		}
}

func (q *ReserveOrdersForFfmCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ReserveOrdersForFfmArgs) {
	return ctx,
		&ReserveOrdersForFfmArgs{
			OrderIDs:   q.OrderIDs,
			Fulfill:    q.Fulfill,
			FulfillIDs: q.FulfillIDs,
		}
}

func (q *UpdateOrderShippingStatusCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateOrderShippingStatusArgs) {
	return ctx,
		&UpdateOrderShippingStatusArgs{
			ID:                         q.ID,
			FulfillmentShippingStates:  q.FulfillmentShippingStates,
			FulfillmentShippingStatus:  q.FulfillmentShippingStatus,
			FulfillmentPaymentStatuses: q.FulfillmentPaymentStatuses,
			EtopPaymentStatus:          q.EtopPaymentStatus,
			CODEtopPaidAt:              q.CODEtopPaidAt,
		}
}

func (q *UpdateOrdersConfirmStatusCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateOrdersConfirmStatusArgs) {
	return ctx,
		&UpdateOrdersConfirmStatusArgs{
			IDs:           q.IDs,
			ShopConfirm:   q.ShopConfirm,
			ConfirmStatus: q.ConfirmStatus,
		}
}

func (q *ValidateOrdersForShippingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ValidateOrdersForShippingArgs) {
	return ctx,
		&ValidateOrdersForShippingArgs{
			OrderIDs: q.OrderIDs,
		}
}

func (q *GetOrderByCodeQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetOrderByCodeArgs) {
	return ctx,
		&GetOrderByCodeArgs{
			Code: q.Code,
		}
}

func (q *GetOrderByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetOrderByIDArgs) {
	return ctx,
		&GetOrderByIDArgs{
			ID: q.ID,
		}
}

func (q *GetOrdersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetOrdersArgs) {
	return ctx,
		&GetOrdersArgs{
			ShopID: q.ShopID,
			IDs:    q.IDs,
		}
}

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	meta.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleReleaseOrdersForFfm)
	b.AddHandler(h.HandleReserveOrdersForFfm)
	b.AddHandler(h.HandleUpdateOrderShippingStatus)
	b.AddHandler(h.HandleUpdateOrdersConfirmStatus)
	b.AddHandler(h.HandleValidateOrdersForShipping)
	return CommandBus{b}
}

func (h AggregateHandler) HandleReleaseOrdersForFfm(ctx context.Context, msg *ReleaseOrdersForFfmCommand) error {
	result, err := h.inner.ReleaseOrdersForFfm(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleReserveOrdersForFfm(ctx context.Context, msg *ReserveOrdersForFfmCommand) error {
	result, err := h.inner.ReserveOrdersForFfm(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleUpdateOrderShippingStatus(ctx context.Context, msg *UpdateOrderShippingStatusCommand) error {
	result, err := h.inner.UpdateOrderShippingStatus(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleUpdateOrdersConfirmStatus(ctx context.Context, msg *UpdateOrdersConfirmStatusCommand) error {
	result, err := h.inner.UpdateOrdersConfirmStatus(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleValidateOrdersForShipping(ctx context.Context, msg *ValidateOrdersForShippingCommand) error {
	result, err := h.inner.ValidateOrdersForShipping(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	meta.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleGetOrderByCode)
	b.AddHandler(h.HandleGetOrderByID)
	b.AddHandler(h.HandleGetOrders)
	return QueryBus{b}
}

func (h QueryServiceHandler) HandleGetOrderByCode(ctx context.Context, msg *GetOrderByCodeQuery) error {
	result, err := h.inner.GetOrderByCode(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetOrderByID(ctx context.Context, msg *GetOrderByIDQuery) error {
	result, err := h.inner.GetOrderByID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetOrders(ctx context.Context, msg *GetOrdersQuery) error {
	result, err := h.inner.GetOrders(msg.GetArgs(ctx))
	msg.Result = result
	return err
}
