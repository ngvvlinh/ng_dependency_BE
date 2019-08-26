// Code generated by gen-cmd-query. DO NOT EDIT.

package addressing

import (
	context "context"

	types "etop.vn/api/main/ordering/types"
	meta "etop.vn/api/meta"
	dot "etop.vn/capi/dot"
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

type CreateAddressCommand struct {
	ShopID       int64
	TraderID     int64
	FullName     string
	Phone        string
	Email        string
	Company      string
	Address1     string
	Address2     string
	DistrictCode string
	WardCode     string
	Coordinates  *types.Coordinates

	Result *ShopTraderAddress `json:"-"`
}

type DeleteAddressCommand struct {
	ID     int64
	ShopID int64

	Result int `json:"-"`
}

type UpdateAddressCommand struct {
	ID           int64
	ShopID       int64
	FullName     dot.NullString
	Phone        dot.NullString
	Email        dot.NullString
	Company      dot.NullString
	Address1     dot.NullString
	Address2     dot.NullString
	DistrictCode dot.NullString
	WardCode     dot.NullString
	Coordinates  *types.Coordinates

	Result *ShopTraderAddress `json:"-"`
}

type GetAddressByIDQuery struct {
	ID     int64
	ShopID int64

	Result *ShopTraderAddress `json:"-"`
}

type ListAddressesByTraderIDQuery struct {
	ShopID   int64
	TraderID int64

	Result []*ShopTraderAddress `json:"-"`
}

// implement interfaces

func (q *CreateAddressCommand) command()       {}
func (q *DeleteAddressCommand) command()       {}
func (q *UpdateAddressCommand) command()       {}
func (q *GetAddressByIDQuery) query()          {}
func (q *ListAddressesByTraderIDQuery) query() {}

// implement conversion

func (q *CreateAddressCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateAddressArgs) {
	return ctx,
		&CreateAddressArgs{
			ShopID:       q.ShopID,
			TraderID:     q.TraderID,
			FullName:     q.FullName,
			Phone:        q.Phone,
			Email:        q.Email,
			Company:      q.Company,
			Address1:     q.Address1,
			Address2:     q.Address2,
			DistrictCode: q.DistrictCode,
			WardCode:     q.WardCode,
			Coordinates:  q.Coordinates,
		}
}

func (q *DeleteAddressCommand) GetArgs(ctx context.Context) (_ context.Context, ID int64, ShopID int64) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *UpdateAddressCommand) GetArgs(ctx context.Context) (_ context.Context, ID int64, ShopID int64, _ *UpdateAddressArgs) {
	return ctx,
		q.ID,
		q.ShopID,
		&UpdateAddressArgs{
			FullName:     q.FullName,
			Phone:        q.Phone,
			Email:        q.Email,
			Company:      q.Company,
			Address1:     q.Address1,
			Address2:     q.Address2,
			DistrictCode: q.DistrictCode,
			WardCode:     q.WardCode,
			Coordinates:  q.Coordinates,
		}
}

func (q *GetAddressByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID int64, ShopID int64) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *ListAddressesByTraderIDQuery) GetArgs(ctx context.Context) (_ context.Context, ShopID int64, TraderID int64) {
	return ctx,
		q.ShopID,
		q.TraderID
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
	b.AddHandler(h.HandleCreateAddress)
	b.AddHandler(h.HandleDeleteAddress)
	b.AddHandler(h.HandleUpdateAddress)
	return CommandBus{b}
}

func (h AggregateHandler) HandleCreateAddress(ctx context.Context, msg *CreateAddressCommand) error {
	result, err := h.inner.CreateAddress(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleDeleteAddress(ctx context.Context, msg *DeleteAddressCommand) error {
	result, err := h.inner.DeleteAddress(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleUpdateAddress(ctx context.Context, msg *UpdateAddressCommand) error {
	result, err := h.inner.UpdateAddress(msg.GetArgs(ctx))
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
	b.AddHandler(h.HandleGetAddressByID)
	b.AddHandler(h.HandleListAddressesByTraderID)
	return QueryBus{b}
}

func (h QueryServiceHandler) HandleGetAddressByID(ctx context.Context, msg *GetAddressByIDQuery) error {
	result, err := h.inner.GetAddressByID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleListAddressesByTraderID(ctx context.Context, msg *ListAddressesByTraderIDQuery) error {
	result, err := h.inner.ListAddressesByTraderID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}
