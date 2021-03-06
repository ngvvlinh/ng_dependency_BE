// +build !generator

// Code generated by generator api. DO NOT EDIT.

package subscription

import (
	context "context"
	time "time"

	meta "o.o/api/meta"
	subscriptingtypes "o.o/api/subscripting/types"
	status3 "o.o/api/top/types/etc/status3"
	subscription_product_type "o.o/api/top/types/etc/subscription_product_type"
	capi "o.o/capi"
	dot "o.o/capi/dot"
)

type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus     { return QueryBus{bus} }

func (b CommandBus) Dispatch(ctx context.Context, msg interface{ command() }) error {
	return b.bus.Dispatch(ctx, msg)
}
func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type ActivateSubscriptionCommand struct {
	ID        dot.ID
	AccountID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleActivateSubscription(ctx context.Context, msg *ActivateSubscriptionCommand) (err error) {
	return h.inner.ActivateSubscription(msg.GetArgs(ctx))
}

type CancelSubscriptionCommand struct {
	ID        dot.ID
	AccountID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleCancelSubscription(ctx context.Context, msg *CancelSubscriptionCommand) (err error) {
	return h.inner.CancelSubscription(msg.GetArgs(ctx))
}

type CreateSubscriptionCommand struct {
	AccountID            dot.ID
	CancelAtPeriodEnd    bool
	CurrentPeriodEndAt   time.Time
	CurrentPeriodStartAt time.Time
	Lines                []*SubscriptionLine
	BillingCycleAnchorAt time.Time
	Customer             *subscriptingtypes.CustomerInfo

	Result *SubscriptionFtLine `json:"-"`
}

func (h AggregateHandler) HandleCreateSubscription(ctx context.Context, msg *CreateSubscriptionCommand) (err error) {
	msg.Result, err = h.inner.CreateSubscription(msg.GetArgs(ctx))
	return err
}

type DeleteSubscriptionCommand struct {
	ID        dot.ID
	AccountID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleDeleteSubscription(ctx context.Context, msg *DeleteSubscriptionCommand) (err error) {
	return h.inner.DeleteSubscription(msg.GetArgs(ctx))
}

type UpdateSubscripionStatusCommand struct {
	ID        dot.ID
	AccountID dot.ID
	Status    status3.NullStatus

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateSubscripionStatus(ctx context.Context, msg *UpdateSubscripionStatusCommand) (err error) {
	return h.inner.UpdateSubscripionStatus(msg.GetArgs(ctx))
}

type UpdateSubscriptionInfoCommand struct {
	ID                   dot.ID
	AccountID            dot.ID
	CancelAtPeriodEnd    dot.NullBool
	BillingCycleAnchorAt time.Time
	Customer             *subscriptingtypes.CustomerInfo
	Lines                []*SubscriptionLine

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateSubscriptionInfo(ctx context.Context, msg *UpdateSubscriptionInfoCommand) (err error) {
	return h.inner.UpdateSubscriptionInfo(msg.GetArgs(ctx))
}

type UpdateSubscriptionPeriodCommand struct {
	ID                   dot.ID
	AccountID            dot.ID
	CancelAtPeriodEnd    bool
	CurrentPeriodStartAt time.Time
	CurrentPeriodEndAt   time.Time
	BillingCycleAnchorAt time.Time
	StartedAt            time.Time

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateSubscriptionPeriod(ctx context.Context, msg *UpdateSubscriptionPeriodCommand) (err error) {
	return h.inner.UpdateSubscriptionPeriod(msg.GetArgs(ctx))
}

type GetLastestSubscriptionByProductTypeQuery struct {
	AccountID   dot.ID
	ProductType subscription_product_type.ProductSubscriptionType

	Result *SubscriptionFtLine `json:"-"`
}

func (h QueryServiceHandler) HandleGetLastestSubscriptionByProductType(ctx context.Context, msg *GetLastestSubscriptionByProductTypeQuery) (err error) {
	msg.Result, err = h.inner.GetLastestSubscriptionByProductType(msg.GetArgs(ctx))
	return err
}

type GetSubscriptionByIDQuery struct {
	ID        dot.ID
	AccountID dot.ID

	Result *SubscriptionFtLine `json:"-"`
}

func (h QueryServiceHandler) HandleGetSubscriptionByID(ctx context.Context, msg *GetSubscriptionByIDQuery) (err error) {
	msg.Result, err = h.inner.GetSubscriptionByID(msg.GetArgs(ctx))
	return err
}

type ListSubscriptionsQuery struct {
	AccountID dot.ID
	Paging    meta.Paging
	Filters   meta.Filters

	Result *ListSubscriptionsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListSubscriptions(ctx context.Context, msg *ListSubscriptionsQuery) (err error) {
	msg.Result, err = h.inner.ListSubscriptions(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *ActivateSubscriptionCommand) command()     {}
func (q *CancelSubscriptionCommand) command()       {}
func (q *CreateSubscriptionCommand) command()       {}
func (q *DeleteSubscriptionCommand) command()       {}
func (q *UpdateSubscripionStatusCommand) command()  {}
func (q *UpdateSubscriptionInfoCommand) command()   {}
func (q *UpdateSubscriptionPeriodCommand) command() {}

func (q *GetLastestSubscriptionByProductTypeQuery) query() {}
func (q *GetSubscriptionByIDQuery) query()                 {}
func (q *ListSubscriptionsQuery) query()                   {}

// implement conversion

func (q *ActivateSubscriptionCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, AccountID dot.ID) {
	return ctx,
		q.ID,
		q.AccountID
}

func (q *CancelSubscriptionCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, AccountID dot.ID) {
	return ctx,
		q.ID,
		q.AccountID
}

func (q *CreateSubscriptionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateSubscriptionArgs) {
	return ctx,
		&CreateSubscriptionArgs{
			AccountID:            q.AccountID,
			CancelAtPeriodEnd:    q.CancelAtPeriodEnd,
			CurrentPeriodEndAt:   q.CurrentPeriodEndAt,
			CurrentPeriodStartAt: q.CurrentPeriodStartAt,
			Lines:                q.Lines,
			BillingCycleAnchorAt: q.BillingCycleAnchorAt,
			Customer:             q.Customer,
		}
}

func (q *CreateSubscriptionCommand) SetCreateSubscriptionArgs(args *CreateSubscriptionArgs) {
	q.AccountID = args.AccountID
	q.CancelAtPeriodEnd = args.CancelAtPeriodEnd
	q.CurrentPeriodEndAt = args.CurrentPeriodEndAt
	q.CurrentPeriodStartAt = args.CurrentPeriodStartAt
	q.Lines = args.Lines
	q.BillingCycleAnchorAt = args.BillingCycleAnchorAt
	q.Customer = args.Customer
}

func (q *DeleteSubscriptionCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, AccountID dot.ID) {
	return ctx,
		q.ID,
		q.AccountID
}

func (q *UpdateSubscripionStatusCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateSubscriptionStatusArgs) {
	return ctx,
		&UpdateSubscriptionStatusArgs{
			ID:        q.ID,
			AccountID: q.AccountID,
			Status:    q.Status,
		}
}

func (q *UpdateSubscripionStatusCommand) SetUpdateSubscriptionStatusArgs(args *UpdateSubscriptionStatusArgs) {
	q.ID = args.ID
	q.AccountID = args.AccountID
	q.Status = args.Status
}

func (q *UpdateSubscriptionInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateSubscriptionInfoArgs) {
	return ctx,
		&UpdateSubscriptionInfoArgs{
			ID:                   q.ID,
			AccountID:            q.AccountID,
			CancelAtPeriodEnd:    q.CancelAtPeriodEnd,
			BillingCycleAnchorAt: q.BillingCycleAnchorAt,
			Customer:             q.Customer,
			Lines:                q.Lines,
		}
}

func (q *UpdateSubscriptionInfoCommand) SetUpdateSubscriptionInfoArgs(args *UpdateSubscriptionInfoArgs) {
	q.ID = args.ID
	q.AccountID = args.AccountID
	q.CancelAtPeriodEnd = args.CancelAtPeriodEnd
	q.BillingCycleAnchorAt = args.BillingCycleAnchorAt
	q.Customer = args.Customer
	q.Lines = args.Lines
}

func (q *UpdateSubscriptionPeriodCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateSubscriptionPeriodArgs) {
	return ctx,
		&UpdateSubscriptionPeriodArgs{
			ID:                   q.ID,
			AccountID:            q.AccountID,
			CancelAtPeriodEnd:    q.CancelAtPeriodEnd,
			CurrentPeriodStartAt: q.CurrentPeriodStartAt,
			CurrentPeriodEndAt:   q.CurrentPeriodEndAt,
			BillingCycleAnchorAt: q.BillingCycleAnchorAt,
			StartedAt:            q.StartedAt,
		}
}

func (q *UpdateSubscriptionPeriodCommand) SetUpdateSubscriptionPeriodArgs(args *UpdateSubscriptionPeriodArgs) {
	q.ID = args.ID
	q.AccountID = args.AccountID
	q.CancelAtPeriodEnd = args.CancelAtPeriodEnd
	q.CurrentPeriodStartAt = args.CurrentPeriodStartAt
	q.CurrentPeriodEndAt = args.CurrentPeriodEndAt
	q.BillingCycleAnchorAt = args.BillingCycleAnchorAt
	q.StartedAt = args.StartedAt
}

func (q *GetLastestSubscriptionByProductTypeQuery) GetArgs(ctx context.Context) (_ context.Context, AccountID dot.ID, ProductType subscription_product_type.ProductSubscriptionType) {
	return ctx,
		q.AccountID,
		q.ProductType
}

func (q *GetSubscriptionByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, AccountID dot.ID) {
	return ctx,
		q.ID,
		q.AccountID
}

func (q *ListSubscriptionsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListSubscriptionsArgs) {
	return ctx,
		&ListSubscriptionsArgs{
			AccountID: q.AccountID,
			Paging:    q.Paging,
			Filters:   q.Filters,
		}
}

func (q *ListSubscriptionsQuery) SetListSubscriptionsArgs(args *ListSubscriptionsArgs) {
	q.AccountID = args.AccountID
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
	b.AddHandler(h.HandleActivateSubscription)
	b.AddHandler(h.HandleCancelSubscription)
	b.AddHandler(h.HandleCreateSubscription)
	b.AddHandler(h.HandleDeleteSubscription)
	b.AddHandler(h.HandleUpdateSubscripionStatus)
	b.AddHandler(h.HandleUpdateSubscriptionInfo)
	b.AddHandler(h.HandleUpdateSubscriptionPeriod)
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
	b.AddHandler(h.HandleGetLastestSubscriptionByProductType)
	b.AddHandler(h.HandleGetSubscriptionByID)
	b.AddHandler(h.HandleListSubscriptions)
	return QueryBus{b}
}
