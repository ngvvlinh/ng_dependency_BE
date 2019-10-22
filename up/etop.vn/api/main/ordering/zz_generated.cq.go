// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package ordering

import (
	context "context"
	time "time"

	etop "etop.vn/api/main/etop"
	types "etop.vn/api/main/ordering/types"
	capi "etop.vn/capi"
)

type Command interface{ command() }
type Query interface{ query() }
type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus                          { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus                              { return QueryBus{bus} }
func (c CommandBus) Dispatch(ctx context.Context, msg Command) error { return c.bus.Dispatch(ctx, msg) }
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error     { return c.bus.Dispatch(ctx, msg) }
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

func (h AggregateHandler) HandleReleaseOrdersForFfm(ctx context.Context, msg *ReleaseOrdersForFfmCommand) (err error) {
	msg.Result, err = h.inner.ReleaseOrdersForFfm(msg.GetArgs(ctx))
	return err
}

type ReserveOrdersForFfmCommand struct {
	OrderIDs   []int64
	Fulfill    types.Fulfill
	FulfillIDs []int64

	Result *ReserveOrdersForFfmResponse `json:"-"`
}

func (h AggregateHandler) HandleReserveOrdersForFfm(ctx context.Context, msg *ReserveOrdersForFfmCommand) (err error) {
	msg.Result, err = h.inner.ReserveOrdersForFfm(msg.GetArgs(ctx))
	return err
}

type UpdateOrderPaymentInfoCommand struct {
	ID            int64
	PaymentStatus etop.Status4
	PaymentID     int64

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateOrderPaymentInfo(ctx context.Context, msg *UpdateOrderPaymentInfoCommand) (err error) {
	return h.inner.UpdateOrderPaymentInfo(msg.GetArgs(ctx))
}

type UpdateOrderShippingStatusCommand struct {
	ID                         int64
	FulfillmentShippingStates  []string
	FulfillmentShippingStatus  etop.Status5
	FulfillmentPaymentStatuses []int
	EtopPaymentStatus          etop.Status4
	CODEtopPaidAt              time.Time

	Result *UpdateOrderShippingStatusResponse `json:"-"`
}

func (h AggregateHandler) HandleUpdateOrderShippingStatus(ctx context.Context, msg *UpdateOrderShippingStatusCommand) (err error) {
	msg.Result, err = h.inner.UpdateOrderShippingStatus(msg.GetArgs(ctx))
	return err
}

type UpdateOrdersConfirmStatusCommand struct {
	IDs           []int64
	ShopConfirm   etop.Status3
	ConfirmStatus etop.Status3

	Result *UpdateOrdersConfirmStatusResponse `json:"-"`
}

func (h AggregateHandler) HandleUpdateOrdersConfirmStatus(ctx context.Context, msg *UpdateOrdersConfirmStatusCommand) (err error) {
	msg.Result, err = h.inner.UpdateOrdersConfirmStatus(msg.GetArgs(ctx))
	return err
}

type ValidateOrdersForShippingCommand struct {
	OrderIDs []int64

	Result *ValidateOrdersForShippingResponse `json:"-"`
}

func (h AggregateHandler) HandleValidateOrdersForShipping(ctx context.Context, msg *ValidateOrdersForShippingCommand) (err error) {
	msg.Result, err = h.inner.ValidateOrdersForShipping(msg.GetArgs(ctx))
	return err
}

type GetOrderByCodeQuery struct {
	Code string

	Result *Order `json:"-"`
}

func (h QueryServiceHandler) HandleGetOrderByCode(ctx context.Context, msg *GetOrderByCodeQuery) (err error) {
	msg.Result, err = h.inner.GetOrderByCode(msg.GetArgs(ctx))
	return err
}

type GetOrderByIDQuery struct {
	ID int64

	Result *Order `json:"-"`
}

func (h QueryServiceHandler) HandleGetOrderByID(ctx context.Context, msg *GetOrderByIDQuery) (err error) {
	msg.Result, err = h.inner.GetOrderByID(msg.GetArgs(ctx))
	return err
}

type GetOrdersQuery struct {
	ShopID int64
	IDs    []int64

	Result *OrdersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetOrders(ctx context.Context, msg *GetOrdersQuery) (err error) {
	msg.Result, err = h.inner.GetOrders(msg.GetArgs(ctx))
	return err
}

type GetOrdersByIDsAndCustomerIDQuery struct {
	ShopID     int64
	IDs        []int64
	CustomerID int64

	Result *OrdersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetOrdersByIDsAndCustomerID(ctx context.Context, msg *GetOrdersByIDsAndCustomerIDQuery) (err error) {
	msg.Result, err = h.inner.GetOrdersByIDsAndCustomerID(msg.GetArgs(ctx))
	return err
}

type ListOrdersByCustomerIDQuery struct {
	ShopID     int64
	CustomerID int64

	Result *OrdersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListOrdersByCustomerID(ctx context.Context, msg *ListOrdersByCustomerIDQuery) (err error) {
	msg.Result, err = h.inner.ListOrdersByCustomerID(msg.GetArgs(ctx))
	return err
}

type ListOrdersByCustomerIDsQuery struct {
	ShopID      int64
	CustomerIDs []int64

	Result *OrdersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListOrdersByCustomerIDs(ctx context.Context, msg *ListOrdersByCustomerIDsQuery) (err error) {
	msg.Result, err = h.inner.ListOrdersByCustomerIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *ReleaseOrdersForFfmCommand) command()       {}
func (q *ReserveOrdersForFfmCommand) command()       {}
func (q *UpdateOrderPaymentInfoCommand) command()    {}
func (q *UpdateOrderShippingStatusCommand) command() {}
func (q *UpdateOrdersConfirmStatusCommand) command() {}
func (q *ValidateOrdersForShippingCommand) command() {}
func (q *GetOrderByCodeQuery) query()                {}
func (q *GetOrderByIDQuery) query()                  {}
func (q *GetOrdersQuery) query()                     {}
func (q *GetOrdersByIDsAndCustomerIDQuery) query()   {}
func (q *ListOrdersByCustomerIDQuery) query()        {}
func (q *ListOrdersByCustomerIDsQuery) query()       {}

// implement conversion

func (q *ReleaseOrdersForFfmCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ReleaseOrdersForFfmArgs) {
	return ctx,
		&ReleaseOrdersForFfmArgs{
			OrderIDs: q.OrderIDs,
		}
}

func (q *ReleaseOrdersForFfmCommand) SetReleaseOrdersForFfmArgs(args *ReleaseOrdersForFfmArgs) {
	q.OrderIDs = args.OrderIDs
}

func (q *ReserveOrdersForFfmCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ReserveOrdersForFfmArgs) {
	return ctx,
		&ReserveOrdersForFfmArgs{
			OrderIDs:   q.OrderIDs,
			Fulfill:    q.Fulfill,
			FulfillIDs: q.FulfillIDs,
		}
}

func (q *ReserveOrdersForFfmCommand) SetReserveOrdersForFfmArgs(args *ReserveOrdersForFfmArgs) {
	q.OrderIDs = args.OrderIDs
	q.Fulfill = args.Fulfill
	q.FulfillIDs = args.FulfillIDs
}

func (q *UpdateOrderPaymentInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateOrderPaymentInfoArgs) {
	return ctx,
		&UpdateOrderPaymentInfoArgs{
			ID:            q.ID,
			PaymentStatus: q.PaymentStatus,
			PaymentID:     q.PaymentID,
		}
}

func (q *UpdateOrderPaymentInfoCommand) SetUpdateOrderPaymentInfoArgs(args *UpdateOrderPaymentInfoArgs) {
	q.ID = args.ID
	q.PaymentStatus = args.PaymentStatus
	q.PaymentID = args.PaymentID
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

func (q *UpdateOrderShippingStatusCommand) SetUpdateOrderShippingStatusArgs(args *UpdateOrderShippingStatusArgs) {
	q.ID = args.ID
	q.FulfillmentShippingStates = args.FulfillmentShippingStates
	q.FulfillmentShippingStatus = args.FulfillmentShippingStatus
	q.FulfillmentPaymentStatuses = args.FulfillmentPaymentStatuses
	q.EtopPaymentStatus = args.EtopPaymentStatus
	q.CODEtopPaidAt = args.CODEtopPaidAt
}

func (q *UpdateOrdersConfirmStatusCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateOrdersConfirmStatusArgs) {
	return ctx,
		&UpdateOrdersConfirmStatusArgs{
			IDs:           q.IDs,
			ShopConfirm:   q.ShopConfirm,
			ConfirmStatus: q.ConfirmStatus,
		}
}

func (q *UpdateOrdersConfirmStatusCommand) SetUpdateOrdersConfirmStatusArgs(args *UpdateOrdersConfirmStatusArgs) {
	q.IDs = args.IDs
	q.ShopConfirm = args.ShopConfirm
	q.ConfirmStatus = args.ConfirmStatus
}

func (q *ValidateOrdersForShippingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ValidateOrdersForShippingArgs) {
	return ctx,
		&ValidateOrdersForShippingArgs{
			OrderIDs: q.OrderIDs,
		}
}

func (q *ValidateOrdersForShippingCommand) SetValidateOrdersForShippingArgs(args *ValidateOrdersForShippingArgs) {
	q.OrderIDs = args.OrderIDs
}

func (q *GetOrderByCodeQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetOrderByCodeArgs) {
	return ctx,
		&GetOrderByCodeArgs{
			Code: q.Code,
		}
}

func (q *GetOrderByCodeQuery) SetGetOrderByCodeArgs(args *GetOrderByCodeArgs) {
	q.Code = args.Code
}

func (q *GetOrderByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetOrderByIDArgs) {
	return ctx,
		&GetOrderByIDArgs{
			ID: q.ID,
		}
}

func (q *GetOrderByIDQuery) SetGetOrderByIDArgs(args *GetOrderByIDArgs) {
	q.ID = args.ID
}

func (q *GetOrdersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetOrdersArgs) {
	return ctx,
		&GetOrdersArgs{
			ShopID: q.ShopID,
			IDs:    q.IDs,
		}
}

func (q *GetOrdersQuery) SetGetOrdersArgs(args *GetOrdersArgs) {
	q.ShopID = args.ShopID
	q.IDs = args.IDs
}

func (q *GetOrdersByIDsAndCustomerIDQuery) GetArgs(ctx context.Context) (_ context.Context, shopID int64, IDs []int64, customerID int64) {
	return ctx,
		q.ShopID,
		q.IDs,
		q.CustomerID
}

func (q *ListOrdersByCustomerIDQuery) GetArgs(ctx context.Context) (_ context.Context, shopID int64, customerID int64) {
	return ctx,
		q.ShopID,
		q.CustomerID
}

func (q *ListOrdersByCustomerIDsQuery) GetArgs(ctx context.Context) (_ context.Context, shopID int64, customerIDs []int64) {
	return ctx,
		q.ShopID,
		q.CustomerIDs
}

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleReleaseOrdersForFfm)
	b.AddHandler(h.HandleReserveOrdersForFfm)
	b.AddHandler(h.HandleUpdateOrderPaymentInfo)
	b.AddHandler(h.HandleUpdateOrderShippingStatus)
	b.AddHandler(h.HandleUpdateOrdersConfirmStatus)
	b.AddHandler(h.HandleValidateOrdersForShipping)
	return CommandBus{b}
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleGetOrderByCode)
	b.AddHandler(h.HandleGetOrderByID)
	b.AddHandler(h.HandleGetOrders)
	b.AddHandler(h.HandleGetOrdersByIDsAndCustomerID)
	b.AddHandler(h.HandleListOrdersByCustomerID)
	b.AddHandler(h.HandleListOrdersByCustomerIDs)
	return QueryBus{b}
}
