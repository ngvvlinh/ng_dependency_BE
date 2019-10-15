// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package haravan

import (
	context "context"

	meta "etop.vn/api/meta"
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

type SendUpdateExternalFulfillmentStateCommand struct {
	FulfillmentID int64

	Result *meta.Empty `json:"-"`
}

func (h AggregateHandler) HandleSendUpdateExternalFulfillmentState(ctx context.Context, msg *SendUpdateExternalFulfillmentStateCommand) (err error) {
	msg.Result, err = h.inner.SendUpdateExternalFulfillmentState(msg.GetArgs(ctx))
	return err
}

type SendUpdateExternalPaymentStatusCommand struct {
	FulfillmentID int64

	Result *meta.Empty `json:"-"`
}

func (h AggregateHandler) HandleSendUpdateExternalPaymentStatus(ctx context.Context, msg *SendUpdateExternalPaymentStatusCommand) (err error) {
	msg.Result, err = h.inner.SendUpdateExternalPaymentStatus(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *SendUpdateExternalFulfillmentStateCommand) command() {}
func (q *SendUpdateExternalPaymentStatusCommand) command()    {}

// implement conversion

func (q *SendUpdateExternalFulfillmentStateCommand) GetArgs(ctx context.Context) (_ context.Context, _ *SendUpdateExternalFulfillmentStateArgs) {
	return ctx,
		&SendUpdateExternalFulfillmentStateArgs{
			FulfillmentID: q.FulfillmentID,
		}
}

func (q *SendUpdateExternalPaymentStatusCommand) GetArgs(ctx context.Context) (_ context.Context, _ *SendUpdateExternalPaymentStatusArgs) {
	return ctx,
		&SendUpdateExternalPaymentStatusArgs{
			FulfillmentID: q.FulfillmentID,
		}
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
	b.AddHandler(h.HandleSendUpdateExternalFulfillmentState)
	b.AddHandler(h.HandleSendUpdateExternalPaymentStatus)
	return CommandBus{b}
}
