// +build !generator

// Code generated by generator api. DO NOT EDIT.

package pricelistpromotion

import (
	context "context"
	time "time"

	meta "o.o/api/meta"
	status3 "o.o/api/top/types/etc/status3"
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

type CreatePriceListPromotionCommand struct {
	PriceListID   dot.ID
	Name          string
	Description   string
	ConnectionID  dot.ID
	DateFrom      time.Time
	DateTo        time.Time
	AppliedRules  *AppliedRules
	PriorityPoint int

	Result *ShipmentPriceListPromotion `json:"-"`
}

func (h AggregateHandler) HandleCreatePriceListPromotion(ctx context.Context, msg *CreatePriceListPromotionCommand) (err error) {
	msg.Result, err = h.inner.CreatePriceListPromotion(msg.GetArgs(ctx))
	return err
}

type DeletePriceListPromotionCommand struct {
	ID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleDeletePriceListPromotion(ctx context.Context, msg *DeletePriceListPromotionCommand) (err error) {
	return h.inner.DeletePriceListPromotion(msg.GetArgs(ctx))
}

type UpdatePriceListPromotionCommand struct {
	ID            dot.ID
	Name          string
	Description   string
	DateFrom      time.Time
	DateTo        time.Time
	AppliedRules  *AppliedRules
	PriorityPoint int
	Status        status3.NullStatus
	ConnectionID  dot.ID
	PriceListID   dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdatePriceListPromotion(ctx context.Context, msg *UpdatePriceListPromotionCommand) (err error) {
	return h.inner.UpdatePriceListPromotion(msg.GetArgs(ctx))
}

type GetPriceListPromotionQuery struct {
	ID dot.ID

	Result *ShipmentPriceListPromotion `json:"-"`
}

func (h QueryServiceHandler) HandleGetPriceListPromotion(ctx context.Context, msg *GetPriceListPromotionQuery) (err error) {
	msg.Result, err = h.inner.GetPriceListPromotion(msg.GetArgs(ctx))
	return err
}

type GetValidPriceListPromotionQuery struct {
	ShopID           dot.ID
	FromProvinceCode string
	ConnectionID     dot.ID

	Result *ShipmentPriceListPromotion `json:"-"`
}

func (h QueryServiceHandler) HandleGetValidPriceListPromotion(ctx context.Context, msg *GetValidPriceListPromotionQuery) (err error) {
	msg.Result, err = h.inner.GetValidPriceListPromotion(msg.GetArgs(ctx))
	return err
}

type ListPriceListPromotionsQuery struct {
	ConnectionID dot.ID
	PriceListID  dot.ID
	Paging       meta.Paging

	Result []*ShipmentPriceListPromotion `json:"-"`
}

func (h QueryServiceHandler) HandleListPriceListPromotions(ctx context.Context, msg *ListPriceListPromotionsQuery) (err error) {
	msg.Result, err = h.inner.ListPriceListPromotions(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreatePriceListPromotionCommand) command() {}
func (q *DeletePriceListPromotionCommand) command() {}
func (q *UpdatePriceListPromotionCommand) command() {}

func (q *GetPriceListPromotionQuery) query()      {}
func (q *GetValidPriceListPromotionQuery) query() {}
func (q *ListPriceListPromotionsQuery) query()    {}

// implement conversion

func (q *CreatePriceListPromotionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreatePriceListPromotionArgs) {
	return ctx,
		&CreatePriceListPromotionArgs{
			PriceListID:   q.PriceListID,
			Name:          q.Name,
			Description:   q.Description,
			ConnectionID:  q.ConnectionID,
			DateFrom:      q.DateFrom,
			DateTo:        q.DateTo,
			AppliedRules:  q.AppliedRules,
			PriorityPoint: q.PriorityPoint,
		}
}

func (q *CreatePriceListPromotionCommand) SetCreatePriceListPromotionArgs(args *CreatePriceListPromotionArgs) {
	q.PriceListID = args.PriceListID
	q.Name = args.Name
	q.Description = args.Description
	q.ConnectionID = args.ConnectionID
	q.DateFrom = args.DateFrom
	q.DateTo = args.DateTo
	q.AppliedRules = args.AppliedRules
	q.PriorityPoint = args.PriorityPoint
}

func (q *DeletePriceListPromotionCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *UpdatePriceListPromotionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdatePriceListPromotionArgs) {
	return ctx,
		&UpdatePriceListPromotionArgs{
			ID:            q.ID,
			Name:          q.Name,
			Description:   q.Description,
			DateFrom:      q.DateFrom,
			DateTo:        q.DateTo,
			AppliedRules:  q.AppliedRules,
			PriorityPoint: q.PriorityPoint,
			Status:        q.Status,
			ConnectionID:  q.ConnectionID,
			PriceListID:   q.PriceListID,
		}
}

func (q *UpdatePriceListPromotionCommand) SetUpdatePriceListPromotionArgs(args *UpdatePriceListPromotionArgs) {
	q.ID = args.ID
	q.Name = args.Name
	q.Description = args.Description
	q.DateFrom = args.DateFrom
	q.DateTo = args.DateTo
	q.AppliedRules = args.AppliedRules
	q.PriorityPoint = args.PriorityPoint
	q.Status = args.Status
	q.ConnectionID = args.ConnectionID
	q.PriceListID = args.PriceListID
}

func (q *GetPriceListPromotionQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetValidPriceListPromotionQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetValidPriceListPromotionArgs) {
	return ctx,
		&GetValidPriceListPromotionArgs{
			ShopID:           q.ShopID,
			FromProvinceCode: q.FromProvinceCode,
			ConnectionID:     q.ConnectionID,
		}
}

func (q *GetValidPriceListPromotionQuery) SetGetValidPriceListPromotionArgs(args *GetValidPriceListPromotionArgs) {
	q.ShopID = args.ShopID
	q.FromProvinceCode = args.FromProvinceCode
	q.ConnectionID = args.ConnectionID
}

func (q *ListPriceListPromotionsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListPriceListPromotionArgs) {
	return ctx,
		&ListPriceListPromotionArgs{
			ConnectionID: q.ConnectionID,
			PriceListID:  q.PriceListID,
			Paging:       q.Paging,
		}
}

func (q *ListPriceListPromotionsQuery) SetListPriceListPromotionArgs(args *ListPriceListPromotionArgs) {
	q.ConnectionID = args.ConnectionID
	q.PriceListID = args.PriceListID
	q.Paging = args.Paging
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
	b.AddHandler(h.HandleCreatePriceListPromotion)
	b.AddHandler(h.HandleDeletePriceListPromotion)
	b.AddHandler(h.HandleUpdatePriceListPromotion)
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
	b.AddHandler(h.HandleGetPriceListPromotion)
	b.AddHandler(h.HandleGetValidPriceListPromotion)
	b.AddHandler(h.HandleListPriceListPromotions)
	return QueryBus{b}
}
