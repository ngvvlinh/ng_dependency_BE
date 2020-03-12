// +build !generator

// Code generated by generator api. DO NOT EDIT.

package shipmentservice

import (
	context "context"

	meta "etop.vn/api/meta"
	status3 "etop.vn/api/top/types/etc/status3"
	capi "etop.vn/capi"
	dot "etop.vn/capi/dot"
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

type CreateShipmentServiceCommand struct {
	ConnectionID dot.ID
	Name         string
	EdCode       string
	ServiceIDs   []string
	Description  string
	ImageURL     string

	Result *ShipmentService `json:"-"`
}

func (h AggregateHandler) HandleCreateShipmentService(ctx context.Context, msg *CreateShipmentServiceCommand) (err error) {
	msg.Result, err = h.inner.CreateShipmentService(msg.GetArgs(ctx))
	return err
}

type DeleteShipmentServiceCommand struct {
	ID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleDeleteShipmentService(ctx context.Context, msg *DeleteShipmentServiceCommand) (err error) {
	return h.inner.DeleteShipmentService(msg.GetArgs(ctx))
}

type UpdateShipmentServiceCommand struct {
	ID           dot.ID
	ConnectionID dot.ID
	Name         string
	EdCode       string
	ServiceIDs   []string
	Description  string
	ImageURL     string
	Status       status3.NullStatus

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateShipmentService(ctx context.Context, msg *UpdateShipmentServiceCommand) (err error) {
	return h.inner.UpdateShipmentService(msg.GetArgs(ctx))
}

type GetShipmentServiceQuery struct {
	ID dot.ID

	Result *ShipmentService `json:"-"`
}

func (h QueryServiceHandler) HandleGetShipmentService(ctx context.Context, msg *GetShipmentServiceQuery) (err error) {
	msg.Result, err = h.inner.GetShipmentService(msg.GetArgs(ctx))
	return err
}

type GetShipmentServiceByServiceIDQuery struct {
	ServiceID string
	ConnID    dot.ID

	Result *ShipmentService `json:"-"`
}

func (h QueryServiceHandler) HandleGetShipmentServiceByServiceID(ctx context.Context, msg *GetShipmentServiceByServiceIDQuery) (err error) {
	msg.Result, err = h.inner.GetShipmentServiceByServiceID(msg.GetArgs(ctx))
	return err
}

type ListShipmentServicesQuery struct {
	Result []*ShipmentService `json:"-"`
}

func (h QueryServiceHandler) HandleListShipmentServices(ctx context.Context, msg *ListShipmentServicesQuery) (err error) {
	msg.Result, err = h.inner.ListShipmentServices(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateShipmentServiceCommand) command() {}
func (q *DeleteShipmentServiceCommand) command() {}
func (q *UpdateShipmentServiceCommand) command() {}

func (q *GetShipmentServiceQuery) query()            {}
func (q *GetShipmentServiceByServiceIDQuery) query() {}
func (q *ListShipmentServicesQuery) query()          {}

// implement conversion

func (q *CreateShipmentServiceCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateShipmentServiceArgs) {
	return ctx,
		&CreateShipmentServiceArgs{
			ConnectionID: q.ConnectionID,
			Name:         q.Name,
			EdCode:       q.EdCode,
			ServiceIDs:   q.ServiceIDs,
			Description:  q.Description,
			ImageURL:     q.ImageURL,
		}
}

func (q *CreateShipmentServiceCommand) SetCreateShipmentServiceArgs(args *CreateShipmentServiceArgs) {
	q.ConnectionID = args.ConnectionID
	q.Name = args.Name
	q.EdCode = args.EdCode
	q.ServiceIDs = args.ServiceIDs
	q.Description = args.Description
	q.ImageURL = args.ImageURL
}

func (q *DeleteShipmentServiceCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *UpdateShipmentServiceCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateShipmentServiceArgs) {
	return ctx,
		&UpdateShipmentServiceArgs{
			ID:           q.ID,
			ConnectionID: q.ConnectionID,
			Name:         q.Name,
			EdCode:       q.EdCode,
			ServiceIDs:   q.ServiceIDs,
			Description:  q.Description,
			ImageURL:     q.ImageURL,
			Status:       q.Status,
		}
}

func (q *UpdateShipmentServiceCommand) SetUpdateShipmentServiceArgs(args *UpdateShipmentServiceArgs) {
	q.ID = args.ID
	q.ConnectionID = args.ConnectionID
	q.Name = args.Name
	q.EdCode = args.EdCode
	q.ServiceIDs = args.ServiceIDs
	q.Description = args.Description
	q.ImageURL = args.ImageURL
	q.Status = args.Status
}

func (q *GetShipmentServiceQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetShipmentServiceByServiceIDQuery) GetArgs(ctx context.Context) (_ context.Context, ServiceID string, ConnID dot.ID) {
	return ctx,
		q.ServiceID,
		q.ConnID
}

func (q *ListShipmentServicesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *ListShipmentServicesQuery) SetEmpty(args *meta.Empty) {
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
	b.AddHandler(h.HandleCreateShipmentService)
	b.AddHandler(h.HandleDeleteShipmentService)
	b.AddHandler(h.HandleUpdateShipmentService)
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
	b.AddHandler(h.HandleGetShipmentService)
	b.AddHandler(h.HandleGetShipmentServiceByServiceID)
	b.AddHandler(h.HandleListShipmentServices)
	return QueryBus{b}
}
