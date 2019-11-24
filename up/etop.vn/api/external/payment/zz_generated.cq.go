// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package payment

import (
	context "context"
	json "encoding/json"
	time "time"

	etop "etop.vn/api/main/etop"
	capi "etop.vn/capi"
	dot "etop.vn/capi/dot"
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

type CreateOrUpdatePaymentCommand struct {
	Amount          int
	Status          etop.Status4
	State           PaymentState
	PaymentProvider PaymentProvider
	ExternalTransID string
	ExternalData    json.RawMessage
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`

	Result *Payment `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdatePayment(ctx context.Context, msg *CreateOrUpdatePaymentCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdatePayment(msg.GetArgs(ctx))
	return err
}

type UpdateExternalPaymentInfoCommand struct {
	ID              dot.ID
	Amount          int
	Status          etop.Status4
	State           PaymentState
	ExternalData    json.RawMessage
	ExternalTransID string

	Result *Payment `json:"-"`
}

func (h AggregateHandler) HandleUpdateExternalPaymentInfo(ctx context.Context, msg *UpdateExternalPaymentInfoCommand) (err error) {
	msg.Result, err = h.inner.UpdateExternalPaymentInfo(msg.GetArgs(ctx))
	return err
}

type GetPaymentByExternalTransIDQuery struct {
	TransactionID string

	Result *Payment `json:"-"`
}

func (h QueryServiceHandler) HandleGetPaymentByExternalTransID(ctx context.Context, msg *GetPaymentByExternalTransIDQuery) (err error) {
	msg.Result, err = h.inner.GetPaymentByExternalTransID(msg.GetArgs(ctx))
	return err
}

type GetPaymentByIDQuery struct {
	ID dot.ID

	Result *Payment `json:"-"`
}

func (h QueryServiceHandler) HandleGetPaymentByID(ctx context.Context, msg *GetPaymentByIDQuery) (err error) {
	msg.Result, err = h.inner.GetPaymentByID(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateOrUpdatePaymentCommand) command()     {}
func (q *UpdateExternalPaymentInfoCommand) command() {}
func (q *GetPaymentByExternalTransIDQuery) query()   {}
func (q *GetPaymentByIDQuery) query()                {}

// implement conversion

func (q *CreateOrUpdatePaymentCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreatePaymentArgs) {
	return ctx,
		&CreatePaymentArgs{
			Amount:          q.Amount,
			Status:          q.Status,
			State:           q.State,
			PaymentProvider: q.PaymentProvider,
			ExternalTransID: q.ExternalTransID,
			ExternalData:    q.ExternalData,
			CreatedAt:       q.CreatedAt,
			UpdatedAt:       q.UpdatedAt,
		}
}

func (q *CreateOrUpdatePaymentCommand) SetCreatePaymentArgs(args *CreatePaymentArgs) {
	q.Amount = args.Amount
	q.Status = args.Status
	q.State = args.State
	q.PaymentProvider = args.PaymentProvider
	q.ExternalTransID = args.ExternalTransID
	q.ExternalData = args.ExternalData
	q.CreatedAt = args.CreatedAt
	q.UpdatedAt = args.UpdatedAt
}

func (q *UpdateExternalPaymentInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateExternalPaymentInfoArgs) {
	return ctx,
		&UpdateExternalPaymentInfoArgs{
			ID:              q.ID,
			Amount:          q.Amount,
			Status:          q.Status,
			State:           q.State,
			ExternalData:    q.ExternalData,
			ExternalTransID: q.ExternalTransID,
		}
}

func (q *UpdateExternalPaymentInfoCommand) SetUpdateExternalPaymentInfoArgs(args *UpdateExternalPaymentInfoArgs) {
	q.ID = args.ID
	q.Amount = args.Amount
	q.Status = args.Status
	q.State = args.State
	q.ExternalData = args.ExternalData
	q.ExternalTransID = args.ExternalTransID
}

func (q *GetPaymentByExternalTransIDQuery) GetArgs(ctx context.Context) (_ context.Context, TransactionID string) {
	return ctx,
		q.TransactionID
}

func (q *GetPaymentByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
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
	b.AddHandler(h.HandleCreateOrUpdatePayment)
	b.AddHandler(h.HandleUpdateExternalPaymentInfo)
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
	b.AddHandler(h.HandleGetPaymentByExternalTransID)
	b.AddHandler(h.HandleGetPaymentByID)
	return QueryBus{b}
}
