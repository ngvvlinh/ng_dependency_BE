// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package identity

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

type ConnectCarrierServiceExternalAccountHaravanCommand struct {
	ShopID int64

	Result *meta.Empty `json:"-"`
}

func (h AggregateHandler) HandleConnectCarrierServiceExternalAccountHaravan(ctx context.Context, msg *ConnectCarrierServiceExternalAccountHaravanCommand) (err error) {
	msg.Result, err = h.inner.ConnectCarrierServiceExternalAccountHaravan(msg.GetArgs(ctx))
	return err
}

type CreateExternalAccountHaravanCommand struct {
	ShopID      int64
	Subdomain   string
	RedirectURI string
	Code        string

	Result *ExternalAccountHaravan `json:"-"`
}

func (h AggregateHandler) HandleCreateExternalAccountHaravan(ctx context.Context, msg *CreateExternalAccountHaravanCommand) (err error) {
	msg.Result, err = h.inner.CreateExternalAccountHaravan(msg.GetArgs(ctx))
	return err
}

type DeleteConnectedCarrierServiceExternalAccountHaravanCommand struct {
	ShopID int64

	Result *meta.Empty `json:"-"`
}

func (h AggregateHandler) HandleDeleteConnectedCarrierServiceExternalAccountHaravan(ctx context.Context, msg *DeleteConnectedCarrierServiceExternalAccountHaravanCommand) (err error) {
	msg.Result, err = h.inner.DeleteConnectedCarrierServiceExternalAccountHaravan(msg.GetArgs(ctx))
	return err
}

type UpdateExternalAccountHaravanTokenCommand struct {
	ShopID      int64
	Subdomain   string
	RedirectURI string
	Code        string

	Result *ExternalAccountHaravan `json:"-"`
}

func (h AggregateHandler) HandleUpdateExternalAccountHaravanToken(ctx context.Context, msg *UpdateExternalAccountHaravanTokenCommand) (err error) {
	msg.Result, err = h.inner.UpdateExternalAccountHaravanToken(msg.GetArgs(ctx))
	return err
}

type GetExternalAccountHaravanByShopIDQuery struct {
	ShopID int64

	Result *ExternalAccountHaravan `json:"-"`
}

func (h QueryServiceHandler) HandleGetExternalAccountHaravanByShopID(ctx context.Context, msg *GetExternalAccountHaravanByShopIDQuery) (err error) {
	msg.Result, err = h.inner.GetExternalAccountHaravanByShopID(msg.GetArgs(ctx))
	return err
}

type GetExternalAccountHaravanByXShopIDQuery struct {
	ExternalShopID int

	Result *ExternalAccountHaravan `json:"-"`
}

func (h QueryServiceHandler) HandleGetExternalAccountHaravanByXShopID(ctx context.Context, msg *GetExternalAccountHaravanByXShopIDQuery) (err error) {
	msg.Result, err = h.inner.GetExternalAccountHaravanByXShopID(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *ConnectCarrierServiceExternalAccountHaravanCommand) command()         {}
func (q *CreateExternalAccountHaravanCommand) command()                        {}
func (q *DeleteConnectedCarrierServiceExternalAccountHaravanCommand) command() {}
func (q *UpdateExternalAccountHaravanTokenCommand) command()                   {}
func (q *GetExternalAccountHaravanByShopIDQuery) query()                       {}
func (q *GetExternalAccountHaravanByXShopIDQuery) query()                      {}

// implement conversion

func (q *ConnectCarrierServiceExternalAccountHaravanCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConnectCarrierServiceExternalAccountHaravanArgs) {
	return ctx,
		&ConnectCarrierServiceExternalAccountHaravanArgs{
			ShopID: q.ShopID,
		}
}

func (q *CreateExternalAccountHaravanCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateExternalAccountHaravanArgs) {
	return ctx,
		&CreateExternalAccountHaravanArgs{
			ShopID:      q.ShopID,
			Subdomain:   q.Subdomain,
			RedirectURI: q.RedirectURI,
			Code:        q.Code,
		}
}

func (q *DeleteConnectedCarrierServiceExternalAccountHaravanCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteConnectedCarrierServiceExternalAccountHaravanArgs) {
	return ctx,
		&DeleteConnectedCarrierServiceExternalAccountHaravanArgs{
			ShopID: q.ShopID,
		}
}

func (q *UpdateExternalAccountHaravanTokenCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateExternalAccountHaravanTokenArgs) {
	return ctx,
		&UpdateExternalAccountHaravanTokenArgs{
			ShopID:      q.ShopID,
			Subdomain:   q.Subdomain,
			RedirectURI: q.RedirectURI,
			Code:        q.Code,
		}
}

func (q *GetExternalAccountHaravanByShopIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetExternalAccountHaravanByShopIDQueryArgs) {
	return ctx,
		&GetExternalAccountHaravanByShopIDQueryArgs{
			ShopID: q.ShopID,
		}
}

func (q *GetExternalAccountHaravanByXShopIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetExternalAccountHaravanByXShopIDQueryArgs) {
	return ctx,
		&GetExternalAccountHaravanByXShopIDQueryArgs{
			ExternalShopID: q.ExternalShopID,
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
	b.AddHandler(h.HandleConnectCarrierServiceExternalAccountHaravan)
	b.AddHandler(h.HandleCreateExternalAccountHaravan)
	b.AddHandler(h.HandleDeleteConnectedCarrierServiceExternalAccountHaravan)
	b.AddHandler(h.HandleUpdateExternalAccountHaravanToken)
	return CommandBus{b}
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
	b.AddHandler(h.HandleGetExternalAccountHaravanByShopID)
	b.AddHandler(h.HandleGetExternalAccountHaravanByXShopID)
	return QueryBus{b}
}
