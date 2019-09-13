// Code generated by generator cq. DO NOT EDIT.

// +build !generator

package address

import (
	context "context"

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

type GetAddressByIDQuery struct {
	ID int64

	Result *Address `json:"-"`
}

func (h QueryServiceHandler) HandleGetAddressByID(ctx context.Context, msg *GetAddressByIDQuery) (err error) {
	msg.Result, err = h.inner.GetAddressByID(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *GetAddressByIDQuery) query() {}

// implement conversion

func (q *GetAddressByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetAddressByIDQueryArgs) {
	return ctx,
		&GetAddressByIDQueryArgs{
			ID: q.ID,
		}
}

// implement dispatching

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
	b.AddHandler(h.HandleGetAddressByID)
	return QueryBus{b}
}
