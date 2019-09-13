// Code generated by generator cq. DO NOT EDIT.

// +build !generator

package shipping

import (
	context "context"

	types "etop.vn/api/main/shipping/types"
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

type CancelFulfillmentCommand struct {
	FulfillmentID int64
	CancelReason  string

	Result *meta.Empty `json:"-"`
}

func (h AggregateHandler) HandleCancelFulfillment(ctx context.Context, msg *CancelFulfillmentCommand) (err error) {
	msg.Result, err = h.inner.CancelFulfillment(msg.GetArgs(ctx))
	return err
}

type ConfirmFulfillmentCommand struct {
	FulfillmentID int64

	Result *meta.Empty `json:"-"`
}

func (h AggregateHandler) HandleConfirmFulfillment(ctx context.Context, msg *ConfirmFulfillmentCommand) (err error) {
	msg.Result, err = h.inner.ConfirmFulfillment(msg.GetArgs(ctx))
	return err
}

type CreateFulfillmentCommand struct {
	OrderID             int64
	PickupAddress       *Address
	ShippingAddress     *Address
	ReturnAddress       *Address
	Carrier             string
	ShippingServiceCode string
	ShippingServiceFee  string
	WeightInfo          WeightInfo
	ValueInfo           ValueInfo
	TryOn               types.TryOn
	ShippingNote        string

	Result *meta.Empty `json:"-"`
}

func (h AggregateHandler) HandleCreateFulfillment(ctx context.Context, msg *CreateFulfillmentCommand) (err error) {
	msg.Result, err = h.inner.CreateFulfillment(msg.GetArgs(ctx))
	return err
}

type GetFulfillmentByIDCommand struct {
	FulfillmentID int64

	Result *Fulfillment `json:"-"`
}

func (h AggregateHandler) HandleGetFulfillmentByID(ctx context.Context, msg *GetFulfillmentByIDCommand) (err error) {
	msg.Result, err = h.inner.GetFulfillmentByID(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CancelFulfillmentCommand) command()  {}
func (q *ConfirmFulfillmentCommand) command() {}
func (q *CreateFulfillmentCommand) command()  {}
func (q *GetFulfillmentByIDCommand) command() {}

// implement conversion

func (q *CancelFulfillmentCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CancelFulfillmentArgs) {
	return ctx,
		&CancelFulfillmentArgs{
			FulfillmentID: q.FulfillmentID,
			CancelReason:  q.CancelReason,
		}
}

func (q *ConfirmFulfillmentCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConfirmFulfillmentArgs) {
	return ctx,
		&ConfirmFulfillmentArgs{
			FulfillmentID: q.FulfillmentID,
		}
}

func (q *CreateFulfillmentCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateFulfillmentArgs) {
	return ctx,
		&CreateFulfillmentArgs{
			OrderID:             q.OrderID,
			PickupAddress:       q.PickupAddress,
			ShippingAddress:     q.ShippingAddress,
			ReturnAddress:       q.ReturnAddress,
			Carrier:             q.Carrier,
			ShippingServiceCode: q.ShippingServiceCode,
			ShippingServiceFee:  q.ShippingServiceFee,
			WeightInfo:          q.WeightInfo,
			ValueInfo:           q.ValueInfo,
			TryOn:               q.TryOn,
			ShippingNote:        q.ShippingNote,
		}
}

func (q *GetFulfillmentByIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *GetFulfillmentByIDQueryArgs) {
	return ctx,
		&GetFulfillmentByIDQueryArgs{
			FulfillmentID: q.FulfillmentID,
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
	b.AddHandler(h.HandleCancelFulfillment)
	b.AddHandler(h.HandleConfirmFulfillment)
	b.AddHandler(h.HandleCreateFulfillment)
	b.AddHandler(h.HandleGetFulfillmentByID)
	return CommandBus{b}
}
