// Code generated by gen-cmd-query. DO NOT EDIT.

package customering

import (
	context "context"

	meta "etop.vn/api/meta"
	metav1 "etop.vn/api/meta/v1"
	shopping "etop.vn/api/shopping"
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

type BatchSetCustomersStatusCommand struct {
	IDs    []int64
	ShopID int64
	Status int32

	Result *meta.UpdatedResponse `json:"-"`
}

type CreateCustomerCommand struct {
	ShopID   int64
	Code     string
	FullName string
	Gender   string
	Type     string
	Birthday string
	Note     string
	Phone    string
	Email    string

	Result *ShopCustomer `json:"-"`
}

type DeleteCustomerCommand struct {
	ID     int64
	ShopID int64

	Result struct {
	} `json:"-"`
}

type UpdateCustomerCommand struct {
	ID       int64
	ShopID   int64
	Code     dot.NullString
	FullName dot.NullString
	Gender   dot.NullString
	Type     dot.NullString
	Birthday dot.NullString
	Note     dot.NullString
	Phone    dot.NullString
	Email    dot.NullString

	Result *ShopCustomer `json:"-"`
}

type GetCustomerByIDQuery struct {
	ID     int64
	ShopID int64

	Result *ShopCustomer `json:"-"`
}

type ListCustomersQuery struct {
	ShopID  int64
	Paging  metav1.Paging
	Filters meta.Filters

	Result *CustomersResponse `json:"-"`
}

type ListCustomersByIDsQuery struct {
	IDs    []int64
	ShopID int64

	Result *CustomersResponse `json:"-"`
}

// implement interfaces

func (q *BatchSetCustomersStatusCommand) command() {}
func (q *CreateCustomerCommand) command()          {}
func (q *DeleteCustomerCommand) command()          {}
func (q *UpdateCustomerCommand) command()          {}
func (q *GetCustomerByIDQuery) query()             {}
func (q *ListCustomersQuery) query()               {}
func (q *ListCustomersByIDsQuery) query()          {}

// implement conversion

func (q *BatchSetCustomersStatusCommand) GetArgs(ctx context.Context) (_ context.Context, IDs []int64, shopID int64, status int32) {
	return ctx,
		q.IDs,
		q.ShopID,
		q.Status
}

func (q *CreateCustomerCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateCustomerArgs) {
	return ctx,
		&CreateCustomerArgs{
			ShopID:   q.ShopID,
			Code:     q.Code,
			FullName: q.FullName,
			Gender:   q.Gender,
			Type:     q.Type,
			Birthday: q.Birthday,
			Note:     q.Note,
			Phone:    q.Phone,
			Email:    q.Email,
		}
}

func (q *DeleteCustomerCommand) GetArgs(ctx context.Context) (_ context.Context, ID int64, shopID int64) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *UpdateCustomerCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateCustomerArgs) {
	return ctx,
		&UpdateCustomerArgs{
			ID:       q.ID,
			ShopID:   q.ShopID,
			Code:     q.Code,
			FullName: q.FullName,
			Gender:   q.Gender,
			Type:     q.Type,
			Birthday: q.Birthday,
			Note:     q.Note,
			Phone:    q.Phone,
			Email:    q.Email,
		}
}

func (q *GetCustomerByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDQueryShopArg) {
	return ctx,
		&shopping.IDQueryShopArg{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *ListCustomersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListCustomersByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
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
	b.AddHandler(h.HandleBatchSetCustomersStatus)
	b.AddHandler(h.HandleCreateCustomer)
	b.AddHandler(h.HandleDeleteCustomer)
	b.AddHandler(h.HandleUpdateCustomer)
	return CommandBus{b}
}

func (h AggregateHandler) HandleBatchSetCustomersStatus(ctx context.Context, msg *BatchSetCustomersStatusCommand) error {
	result, err := h.inner.BatchSetCustomersStatus(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleCreateCustomer(ctx context.Context, msg *CreateCustomerCommand) error {
	result, err := h.inner.CreateCustomer(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleDeleteCustomer(ctx context.Context, msg *DeleteCustomerCommand) error {
	return h.inner.DeleteCustomer(msg.GetArgs(ctx))
}

func (h AggregateHandler) HandleUpdateCustomer(ctx context.Context, msg *UpdateCustomerCommand) error {
	result, err := h.inner.UpdateCustomer(msg.GetArgs(ctx))
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
	b.AddHandler(h.HandleGetCustomerByID)
	b.AddHandler(h.HandleListCustomers)
	b.AddHandler(h.HandleListCustomersByIDs)
	return QueryBus{b}
}

func (h QueryServiceHandler) HandleGetCustomerByID(ctx context.Context, msg *GetCustomerByIDQuery) error {
	result, err := h.inner.GetCustomerByID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleListCustomers(ctx context.Context, msg *ListCustomersQuery) error {
	result, err := h.inner.ListCustomers(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleListCustomersByIDs(ctx context.Context, msg *ListCustomersByIDsQuery) error {
	result, err := h.inner.ListCustomersByIDs(msg.GetArgs(ctx))
	msg.Result = result
	return err
}
