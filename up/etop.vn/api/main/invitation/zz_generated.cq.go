// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package invitation

import (
	context "context"
	time "time"

	etop "etop.vn/api/main/etop"
	meta "etop.vn/api/meta"
	shopping "etop.vn/api/shopping"
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

type AcceptInvitationCommand struct {
	UserID dot.ID
	Token  string

	Result int `json:"-"`
}

func (h AggregateHandler) HandleAcceptInvitation(ctx context.Context, msg *AcceptInvitationCommand) (err error) {
	msg.Result, err = h.inner.AcceptInvitation(msg.GetArgs(ctx))
	return err
}

type CreateInvitationCommand struct {
	AccountID dot.ID
	Email     string
	Roles     []Role
	Status    etop.Status3
	InvitedBy dot.ID
	CreatedBy time.Time

	Result *Invitation `json:"-"`
}

func (h AggregateHandler) HandleCreateInvitation(ctx context.Context, msg *CreateInvitationCommand) (err error) {
	msg.Result, err = h.inner.CreateInvitation(msg.GetArgs(ctx))
	return err
}

type RejectInvitationCommand struct {
	UserID dot.ID
	Token  string

	Result int `json:"-"`
}

func (h AggregateHandler) HandleRejectInvitation(ctx context.Context, msg *RejectInvitationCommand) (err error) {
	msg.Result, err = h.inner.RejectInvitation(msg.GetArgs(ctx))
	return err
}

type GetInvitationQuery struct {
	ID dot.ID

	Result *Invitation `json:"-"`
}

func (h QueryServiceHandler) HandleGetInvitation(ctx context.Context, msg *GetInvitationQuery) (err error) {
	msg.Result, err = h.inner.GetInvitation(msg.GetArgs(ctx))
	return err
}

type GetInvitationByTokenQuery struct {
	Token string

	Result *Invitation `json:"-"`
}

func (h QueryServiceHandler) HandleGetInvitationByToken(ctx context.Context, msg *GetInvitationByTokenQuery) (err error) {
	msg.Result, err = h.inner.GetInvitationByToken(msg.GetArgs(ctx))
	return err
}

type ListInvitationsQuery struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters

	Result *InvitationsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListInvitations(ctx context.Context, msg *ListInvitationsQuery) (err error) {
	msg.Result, err = h.inner.ListInvitations(msg.GetArgs(ctx))
	return err
}

type ListInvitationsAcceptedByEmailQuery struct {
	Email string

	Result *InvitationsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListInvitationsAcceptedByEmail(ctx context.Context, msg *ListInvitationsAcceptedByEmailQuery) (err error) {
	msg.Result, err = h.inner.ListInvitationsAcceptedByEmail(msg.GetArgs(ctx))
	return err
}

type ListInvitationsByEmailQuery struct {
	Email   string
	Paging  meta.Paging
	Filters meta.Filters

	Result *InvitationsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListInvitationsByEmail(ctx context.Context, msg *ListInvitationsByEmailQuery) (err error) {
	msg.Result, err = h.inner.ListInvitationsByEmail(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *AcceptInvitationCommand) command()           {}
func (q *CreateInvitationCommand) command()           {}
func (q *RejectInvitationCommand) command()           {}
func (q *GetInvitationQuery) query()                  {}
func (q *GetInvitationByTokenQuery) query()           {}
func (q *ListInvitationsQuery) query()                {}
func (q *ListInvitationsAcceptedByEmailQuery) query() {}
func (q *ListInvitationsByEmailQuery) query()         {}

// implement conversion

func (q *AcceptInvitationCommand) GetArgs(ctx context.Context) (_ context.Context, userID dot.ID, token string) {
	return ctx,
		q.UserID,
		q.Token
}

func (q *CreateInvitationCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateInvitationArgs) {
	return ctx,
		&CreateInvitationArgs{
			AccountID: q.AccountID,
			Email:     q.Email,
			Roles:     q.Roles,
			Status:    q.Status,
			InvitedBy: q.InvitedBy,
			CreatedBy: q.CreatedBy,
		}
}

func (q *CreateInvitationCommand) SetCreateInvitationArgs(args *CreateInvitationArgs) {
	q.AccountID = args.AccountID
	q.Email = args.Email
	q.Roles = args.Roles
	q.Status = args.Status
	q.InvitedBy = args.InvitedBy
	q.CreatedBy = args.CreatedBy
}

func (q *RejectInvitationCommand) GetArgs(ctx context.Context) (_ context.Context, userID dot.ID, token string) {
	return ctx,
		q.UserID,
		q.Token
}

func (q *GetInvitationQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetInvitationByTokenQuery) GetArgs(ctx context.Context) (_ context.Context, token string) {
	return ctx,
		q.Token
}

func (q *ListInvitationsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListInvitationsQuery) SetListQueryShopArgs(args *shopping.ListQueryShopArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
}

func (q *ListInvitationsAcceptedByEmailQuery) GetArgs(ctx context.Context) (_ context.Context, email string) {
	return ctx,
		q.Email
}

func (q *ListInvitationsByEmailQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListInvitationsByEmailArgs) {
	return ctx,
		&ListInvitationsByEmailArgs{
			Email:   q.Email,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListInvitationsByEmailQuery) SetListInvitationsByEmailArgs(args *ListInvitationsByEmailArgs) {
	q.Email = args.Email
	q.Paging = args.Paging
	q.Filters = args.Filters
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
	b.AddHandler(h.HandleAcceptInvitation)
	b.AddHandler(h.HandleCreateInvitation)
	b.AddHandler(h.HandleRejectInvitation)
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
	b.AddHandler(h.HandleGetInvitation)
	b.AddHandler(h.HandleGetInvitationByToken)
	b.AddHandler(h.HandleListInvitations)
	b.AddHandler(h.HandleListInvitationsAcceptedByEmail)
	b.AddHandler(h.HandleListInvitationsByEmail)
	return QueryBus{b}
}
