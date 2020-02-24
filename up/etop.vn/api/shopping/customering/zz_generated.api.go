// +build !generator

// Code generated by generator api. DO NOT EDIT.

package customering

import (
	context "context"

	meta "etop.vn/api/meta"
	shopping "etop.vn/api/shopping"
	customer_type "etop.vn/api/shopping/customering/customer_type"
	gender "etop.vn/api/top/types/etc/gender"
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

type AddCustomersToGroupCommand struct {
	GroupID     dot.ID
	CustomerIDs []dot.ID
	ShopID      dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleAddCustomersToGroup(ctx context.Context, msg *AddCustomersToGroupCommand) (err error) {
	msg.Result, err = h.inner.AddCustomersToGroup(msg.GetArgs(ctx))
	return err
}

type BatchSetCustomersStatusCommand struct {
	IDs    []dot.ID
	ShopID dot.ID
	Status int

	Result *meta.UpdatedResponse `json:"-"`
}

func (h AggregateHandler) HandleBatchSetCustomersStatus(ctx context.Context, msg *BatchSetCustomersStatusCommand) (err error) {
	msg.Result, err = h.inner.BatchSetCustomersStatus(msg.GetArgs(ctx))
	return err
}

type CreateCustomerCommand struct {
	ExternalID   string
	ExternalCode string
	PartnerID    dot.ID
	ShopID       dot.ID
	FullName     string
	Gender       gender.Gender
	Type         customer_type.CustomerType
	Birthday     string
	Note         string
	Phone        string
	Email        string

	Result *ShopCustomer `json:"-"`
}

func (h AggregateHandler) HandleCreateCustomer(ctx context.Context, msg *CreateCustomerCommand) (err error) {
	msg.Result, err = h.inner.CreateCustomer(msg.GetArgs(ctx))
	return err
}

type CreateCustomerGroupCommand struct {
	Name   string
	ShopID dot.ID

	Result *ShopCustomerGroup `json:"-"`
}

func (h AggregateHandler) HandleCreateCustomerGroup(ctx context.Context, msg *CreateCustomerGroupCommand) (err error) {
	msg.Result, err = h.inner.CreateCustomerGroup(msg.GetArgs(ctx))
	return err
}

type DeleteCustomerCommand struct {
	ID         dot.ID
	ShopID     dot.ID
	ExternalID string
	Code       string

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteCustomer(ctx context.Context, msg *DeleteCustomerCommand) (err error) {
	msg.Result, err = h.inner.DeleteCustomer(msg.GetArgs(ctx))
	return err
}

type DeleteGroupCommand struct {
	ShopID  dot.ID
	GroupID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteGroup(ctx context.Context, msg *DeleteGroupCommand) (err error) {
	msg.Result, err = h.inner.DeleteGroup(msg.GetArgs(ctx))
	return err
}

type RemoveCustomersFromGroupCommand struct {
	GroupID     dot.ID
	CustomerIDs []dot.ID
	ShopID      dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleRemoveCustomersFromGroup(ctx context.Context, msg *RemoveCustomersFromGroupCommand) (err error) {
	msg.Result, err = h.inner.RemoveCustomersFromGroup(msg.GetArgs(ctx))
	return err
}

type UpdateCustomerCommand struct {
	ID       dot.ID
	ShopID   dot.ID
	FullName dot.NullString
	Gender   gender.NullGender
	Type     customer_type.CustomerType
	Birthday dot.NullString
	Note     dot.NullString
	Phone    dot.NullString
	Email    dot.NullString

	Result *ShopCustomer `json:"-"`
}

func (h AggregateHandler) HandleUpdateCustomer(ctx context.Context, msg *UpdateCustomerCommand) (err error) {
	msg.Result, err = h.inner.UpdateCustomer(msg.GetArgs(ctx))
	return err
}

type UpdateCustomerGroupCommand struct {
	ID   dot.ID
	Name string

	Result *ShopCustomerGroup `json:"-"`
}

func (h AggregateHandler) HandleUpdateCustomerGroup(ctx context.Context, msg *UpdateCustomerGroupCommand) (err error) {
	msg.Result, err = h.inner.UpdateCustomerGroup(msg.GetArgs(ctx))
	return err
}

type GetCustomerQuery struct {
	ID         dot.ID
	ShopID     dot.ID
	Code       string
	ExternalID string

	Result *ShopCustomer `json:"-"`
}

func (h QueryServiceHandler) HandleGetCustomer(ctx context.Context, msg *GetCustomerQuery) (err error) {
	msg.Result, err = h.inner.GetCustomer(msg.GetArgs(ctx))
	return err
}

type GetCustomerByCodeQuery struct {
	Code   string
	ShopID dot.ID

	Result *ShopCustomer `json:"-"`
}

func (h QueryServiceHandler) HandleGetCustomerByCode(ctx context.Context, msg *GetCustomerByCodeQuery) (err error) {
	msg.Result, err = h.inner.GetCustomerByCode(msg.GetArgs(ctx))
	return err
}

type GetCustomerByEmailQuery struct {
	Email  string
	ShopID dot.ID

	Result *ShopCustomer `json:"-"`
}

func (h QueryServiceHandler) HandleGetCustomerByEmail(ctx context.Context, msg *GetCustomerByEmailQuery) (err error) {
	msg.Result, err = h.inner.GetCustomerByEmail(msg.GetArgs(ctx))
	return err
}

type GetCustomerByIDQuery struct {
	ID             dot.ID
	ShopID         dot.ID
	IncludeDeleted bool

	Result *ShopCustomer `json:"-"`
}

func (h QueryServiceHandler) HandleGetCustomerByID(ctx context.Context, msg *GetCustomerByIDQuery) (err error) {
	msg.Result, err = h.inner.GetCustomerByID(msg.GetArgs(ctx))
	return err
}

type GetCustomerByPhoneQuery struct {
	Phone  string
	ShopID dot.ID

	Result *ShopCustomer `json:"-"`
}

func (h QueryServiceHandler) HandleGetCustomerByPhone(ctx context.Context, msg *GetCustomerByPhoneQuery) (err error) {
	msg.Result, err = h.inner.GetCustomerByPhone(msg.GetArgs(ctx))
	return err
}

type GetCustomerGroupQuery struct {
	ID dot.ID

	Result *ShopCustomerGroup `json:"-"`
}

func (h QueryServiceHandler) HandleGetCustomerGroup(ctx context.Context, msg *GetCustomerGroupQuery) (err error) {
	msg.Result, err = h.inner.GetCustomerGroup(msg.GetArgs(ctx))
	return err
}

type GetCustomerIndependentQuery struct {
	Result *ShopCustomer `json:"-"`
}

func (h QueryServiceHandler) HandleGetCustomerIndependent(ctx context.Context, msg *GetCustomerIndependentQuery) (err error) {
	msg.Result, err = h.inner.GetCustomerIndependent(msg.GetArgs(ctx))
	return err
}

type ListCustomerGroupsQuery struct {
	Paging         meta.Paging
	Filters        meta.Filters
	IncludeDeleted bool

	Result *CustomerGroupsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListCustomerGroups(ctx context.Context, msg *ListCustomerGroupsQuery) (err error) {
	msg.Result, err = h.inner.ListCustomerGroups(msg.GetArgs(ctx))
	return err
}

type ListCustomerGroupsCustomersQuery struct {
	CustomerIDs    []dot.ID
	GroupIDs       []dot.ID
	Paging         meta.Paging
	IncludeDeleted bool

	Result *CustomerGroupsCustomersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListCustomerGroupsCustomers(ctx context.Context, msg *ListCustomerGroupsCustomersQuery) (err error) {
	msg.Result, err = h.inner.ListCustomerGroupsCustomers(msg.GetArgs(ctx))
	return err
}

type ListCustomersQuery struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters

	Result *CustomersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListCustomers(ctx context.Context, msg *ListCustomersQuery) (err error) {
	msg.Result, err = h.inner.ListCustomers(msg.GetArgs(ctx))
	return err
}

type ListCustomersByIDsQuery struct {
	IDs            []dot.ID
	ShopIDs        []dot.ID
	ShopID         dot.ID
	Paging         meta.Paging
	IncludeDeleted bool

	Result *CustomersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListCustomersByIDs(ctx context.Context, msg *ListCustomersByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListCustomersByIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *AddCustomersToGroupCommand) command()      {}
func (q *BatchSetCustomersStatusCommand) command()  {}
func (q *CreateCustomerCommand) command()           {}
func (q *CreateCustomerGroupCommand) command()      {}
func (q *DeleteCustomerCommand) command()           {}
func (q *DeleteGroupCommand) command()              {}
func (q *RemoveCustomersFromGroupCommand) command() {}
func (q *UpdateCustomerCommand) command()           {}
func (q *UpdateCustomerGroupCommand) command()      {}

func (q *GetCustomerQuery) query()                 {}
func (q *GetCustomerByCodeQuery) query()           {}
func (q *GetCustomerByEmailQuery) query()          {}
func (q *GetCustomerByIDQuery) query()             {}
func (q *GetCustomerByPhoneQuery) query()          {}
func (q *GetCustomerGroupQuery) query()            {}
func (q *GetCustomerIndependentQuery) query()      {}
func (q *ListCustomerGroupsQuery) query()          {}
func (q *ListCustomerGroupsCustomersQuery) query() {}
func (q *ListCustomersQuery) query()               {}
func (q *ListCustomersByIDsQuery) query()          {}

// implement conversion

func (q *AddCustomersToGroupCommand) GetArgs(ctx context.Context) (_ context.Context, _ *AddCustomerToGroupArgs) {
	return ctx,
		&AddCustomerToGroupArgs{
			GroupID:     q.GroupID,
			CustomerIDs: q.CustomerIDs,
			ShopID:      q.ShopID,
		}
}

func (q *AddCustomersToGroupCommand) SetAddCustomerToGroupArgs(args *AddCustomerToGroupArgs) {
	q.GroupID = args.GroupID
	q.CustomerIDs = args.CustomerIDs
	q.ShopID = args.ShopID
}

func (q *BatchSetCustomersStatusCommand) GetArgs(ctx context.Context) (_ context.Context, IDs []dot.ID, shopID dot.ID, status int) {
	return ctx,
		q.IDs,
		q.ShopID,
		q.Status
}

func (q *CreateCustomerCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateCustomerArgs) {
	return ctx,
		&CreateCustomerArgs{
			ExternalID:   q.ExternalID,
			ExternalCode: q.ExternalCode,
			PartnerID:    q.PartnerID,
			ShopID:       q.ShopID,
			FullName:     q.FullName,
			Gender:       q.Gender,
			Type:         q.Type,
			Birthday:     q.Birthday,
			Note:         q.Note,
			Phone:        q.Phone,
			Email:        q.Email,
		}
}

func (q *CreateCustomerCommand) SetCreateCustomerArgs(args *CreateCustomerArgs) {
	q.ExternalID = args.ExternalID
	q.ExternalCode = args.ExternalCode
	q.PartnerID = args.PartnerID
	q.ShopID = args.ShopID
	q.FullName = args.FullName
	q.Gender = args.Gender
	q.Type = args.Type
	q.Birthday = args.Birthday
	q.Note = args.Note
	q.Phone = args.Phone
	q.Email = args.Email
}

func (q *CreateCustomerGroupCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateCustomerGroupArgs) {
	return ctx,
		&CreateCustomerGroupArgs{
			Name:   q.Name,
			ShopID: q.ShopID,
		}
}

func (q *CreateCustomerGroupCommand) SetCreateCustomerGroupArgs(args *CreateCustomerGroupArgs) {
	q.Name = args.Name
	q.ShopID = args.ShopID
}

func (q *DeleteCustomerCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteCustomerArgs) {
	return ctx,
		&DeleteCustomerArgs{
			ID:         q.ID,
			ShopID:     q.ShopID,
			ExternalID: q.ExternalID,
			Code:       q.Code,
		}
}

func (q *DeleteCustomerCommand) SetDeleteCustomerArgs(args *DeleteCustomerArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.ExternalID = args.ExternalID
	q.Code = args.Code
}

func (q *DeleteGroupCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteGroupArgs) {
	return ctx,
		&DeleteGroupArgs{
			ShopID:  q.ShopID,
			GroupID: q.GroupID,
		}
}

func (q *DeleteGroupCommand) SetDeleteGroupArgs(args *DeleteGroupArgs) {
	q.ShopID = args.ShopID
	q.GroupID = args.GroupID
}

func (q *RemoveCustomersFromGroupCommand) GetArgs(ctx context.Context) (_ context.Context, _ *RemoveCustomerOutOfGroupArgs) {
	return ctx,
		&RemoveCustomerOutOfGroupArgs{
			GroupID:     q.GroupID,
			CustomerIDs: q.CustomerIDs,
			ShopID:      q.ShopID,
		}
}

func (q *RemoveCustomersFromGroupCommand) SetRemoveCustomerOutOfGroupArgs(args *RemoveCustomerOutOfGroupArgs) {
	q.GroupID = args.GroupID
	q.CustomerIDs = args.CustomerIDs
	q.ShopID = args.ShopID
}

func (q *UpdateCustomerCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateCustomerArgs) {
	return ctx,
		&UpdateCustomerArgs{
			ID:       q.ID,
			ShopID:   q.ShopID,
			FullName: q.FullName,
			Gender:   q.Gender,
			Type:     q.Type,
			Birthday: q.Birthday,
			Note:     q.Note,
			Phone:    q.Phone,
			Email:    q.Email,
		}
}

func (q *UpdateCustomerCommand) SetUpdateCustomerArgs(args *UpdateCustomerArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.FullName = args.FullName
	q.Gender = args.Gender
	q.Type = args.Type
	q.Birthday = args.Birthday
	q.Note = args.Note
	q.Phone = args.Phone
	q.Email = args.Email
}

func (q *UpdateCustomerGroupCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateCustomerGroupArgs) {
	return ctx,
		&UpdateCustomerGroupArgs{
			ID:   q.ID,
			Name: q.Name,
		}
}

func (q *UpdateCustomerGroupCommand) SetUpdateCustomerGroupArgs(args *UpdateCustomerGroupArgs) {
	q.ID = args.ID
	q.Name = args.Name
}

func (q *GetCustomerQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetCustomerArgs) {
	return ctx,
		&GetCustomerArgs{
			ID:         q.ID,
			ShopID:     q.ShopID,
			Code:       q.Code,
			ExternalID: q.ExternalID,
		}
}

func (q *GetCustomerQuery) SetGetCustomerArgs(args *GetCustomerArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.Code = args.Code
	q.ExternalID = args.ExternalID
}

func (q *GetCustomerByCodeQuery) GetArgs(ctx context.Context) (_ context.Context, code string, shopID dot.ID) {
	return ctx,
		q.Code,
		q.ShopID
}

func (q *GetCustomerByEmailQuery) GetArgs(ctx context.Context) (_ context.Context, email string, shopID dot.ID) {
	return ctx,
		q.Email,
		q.ShopID
}

func (q *GetCustomerByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDQueryShopArg) {
	return ctx,
		&shopping.IDQueryShopArg{
			ID:             q.ID,
			ShopID:         q.ShopID,
			IncludeDeleted: q.IncludeDeleted,
		}
}

func (q *GetCustomerByIDQuery) SetIDQueryShopArg(args *shopping.IDQueryShopArg) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.IncludeDeleted = args.IncludeDeleted
}

func (q *GetCustomerByPhoneQuery) GetArgs(ctx context.Context) (_ context.Context, phone string, shopID dot.ID) {
	return ctx,
		q.Phone,
		q.ShopID
}

func (q *GetCustomerGroupQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetCustomerGroupArgs) {
	return ctx,
		&GetCustomerGroupArgs{
			ID: q.ID,
		}
}

func (q *GetCustomerGroupQuery) SetGetCustomerGroupArgs(args *GetCustomerGroupArgs) {
	q.ID = args.ID
}

func (q *GetCustomerIndependentQuery) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *GetCustomerIndependentQuery) SetEmpty(args *meta.Empty) {
}

func (q *ListCustomerGroupsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListCustomerGroupArgs) {
	return ctx,
		&ListCustomerGroupArgs{
			Paging:         q.Paging,
			Filters:        q.Filters,
			IncludeDeleted: q.IncludeDeleted,
		}
}

func (q *ListCustomerGroupsQuery) SetListCustomerGroupArgs(args *ListCustomerGroupArgs) {
	q.Paging = args.Paging
	q.Filters = args.Filters
	q.IncludeDeleted = args.IncludeDeleted
}

func (q *ListCustomerGroupsCustomersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListCustomerGroupsCustomersArgs) {
	return ctx,
		&ListCustomerGroupsCustomersArgs{
			CustomerIDs:    q.CustomerIDs,
			GroupIDs:       q.GroupIDs,
			Paging:         q.Paging,
			IncludeDeleted: q.IncludeDeleted,
		}
}

func (q *ListCustomerGroupsCustomersQuery) SetListCustomerGroupsCustomersArgs(args *ListCustomerGroupsCustomersArgs) {
	q.CustomerIDs = args.CustomerIDs
	q.GroupIDs = args.GroupIDs
	q.Paging = args.Paging
	q.IncludeDeleted = args.IncludeDeleted
}

func (q *ListCustomersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListCustomersQuery) SetListQueryShopArgs(args *shopping.ListQueryShopArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
}

func (q *ListCustomersByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListCustomerByIDsArgs) {
	return ctx,
		&ListCustomerByIDsArgs{
			IDs:            q.IDs,
			ShopIDs:        q.ShopIDs,
			ShopID:         q.ShopID,
			Paging:         q.Paging,
			IncludeDeleted: q.IncludeDeleted,
		}
}

func (q *ListCustomersByIDsQuery) SetListCustomerByIDsArgs(args *ListCustomerByIDsArgs) {
	q.IDs = args.IDs
	q.ShopIDs = args.ShopIDs
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.IncludeDeleted = args.IncludeDeleted
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
	b.AddHandler(h.HandleAddCustomersToGroup)
	b.AddHandler(h.HandleBatchSetCustomersStatus)
	b.AddHandler(h.HandleCreateCustomer)
	b.AddHandler(h.HandleCreateCustomerGroup)
	b.AddHandler(h.HandleDeleteCustomer)
	b.AddHandler(h.HandleDeleteGroup)
	b.AddHandler(h.HandleRemoveCustomersFromGroup)
	b.AddHandler(h.HandleUpdateCustomer)
	b.AddHandler(h.HandleUpdateCustomerGroup)
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
	b.AddHandler(h.HandleGetCustomer)
	b.AddHandler(h.HandleGetCustomerByCode)
	b.AddHandler(h.HandleGetCustomerByEmail)
	b.AddHandler(h.HandleGetCustomerByID)
	b.AddHandler(h.HandleGetCustomerByPhone)
	b.AddHandler(h.HandleGetCustomerGroup)
	b.AddHandler(h.HandleGetCustomerIndependent)
	b.AddHandler(h.HandleListCustomerGroups)
	b.AddHandler(h.HandleListCustomerGroupsCustomers)
	b.AddHandler(h.HandleListCustomers)
	b.AddHandler(h.HandleListCustomersByIDs)
	return QueryBus{b}
}
