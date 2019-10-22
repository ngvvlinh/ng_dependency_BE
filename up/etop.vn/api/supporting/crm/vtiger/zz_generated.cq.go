// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package vtiger

import (
	context "context"
	time "time"

	meta "etop.vn/api/meta"
	capi "etop.vn/capi"
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

type CreateOrUpdateContactCommand struct {
	ID                   string    `json:"id"`
	EtopUserID           int64     `json:"etop_user_id"`
	ContactNo            string    `json:"contact_no"`
	Phone                string    `json:"phone"`
	Lastname             string    `json:"lastname"`
	Mobile               string    `json:"mobile"`
	Leadsource           string    `json:"leadsource"`
	Email                string    `json:"email"`
	Description          string    `json:"decription"`
	Secondaryemail       string    `json:"secondaryemail"`
	Modifiedby           string    `json:"modifiedby"`
	Source               string    `json:"source"`
	Company              string    `json:"company"`
	Website              string    `json:"website"`
	Lane                 string    `json:"lane"`
	City                 string    `json:"city"`
	State                string    `json:"state"`
	Country              string    `json:"country"`
	OrdersPerDay         string    `json:"orders_per_day"`
	UsedShippingProvider string    `json:"used_shipping_provider"`
	Firstname            string    `json:"firstname"`
	Createdtime          time.Time `json:"createdtime"`
	Modifiedtime         time.Time `json:"modifiedtime"`
	AssignedUserID       string    `json:"assigned_user_id"`

	Result *Contact `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateContact(ctx context.Context, msg *CreateOrUpdateContactCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateContact(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateLeadCommand struct {
	ID                   string    `json:"id"`
	EtopUserID           int64     `json:"etop_user_id"`
	ContactNo            string    `json:"contact_no"`
	Phone                string    `json:"phone"`
	Lastname             string    `json:"lastname"`
	Mobile               string    `json:"mobile"`
	Leadsource           string    `json:"leadsource"`
	Email                string    `json:"email"`
	Description          string    `json:"description"`
	Secondaryemail       string    `json:"secondaryemail"`
	Modifiedby           string    `json:"modifiedby"`
	Source               string    `json:"source"`
	Company              string    `json:"company"`
	Website              string    `json:"website"`
	Lane                 string    `json:"lane"`
	City                 string    `json:"city"`
	State                string    `json:"state"`
	Country              string    `json:"country"`
	OrdersPerDay         string    `json:"orders_per_day"`
	UsedShippingProvider string    `json:"used_shipping_provider"`
	Firstname            string    `json:"firstname"`
	AssignedUserID       string    `json:"assigned_user_id"`
	Createdtime          time.Time `json:"createdtime"`
	Modifiedtime         time.Time `json:"modifiedtime"`

	Result *Lead `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateLead(ctx context.Context, msg *CreateOrUpdateLeadCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateLead(msg.GetArgs(ctx))
	return err
}

type CreateTicketCommand struct {
	FfmCode     string
	FfmID       int64
	ID          string
	EtopUserID  int64 `json:"etop_user_id"`
	Code        string
	Title       string
	Value       string
	OldValue    string
	Reason      string
	ShopID      int64
	OrderID     int64
	OrderCode   string
	FfmUrl      string
	Company     string
	Provider    string
	Note        string
	Environment string
	FromApp     string
	Account     Account

	Result *Ticket `json:"-"`
}

func (h AggregateHandler) HandleCreateTicket(ctx context.Context, msg *CreateTicketCommand) (err error) {
	msg.Result, err = h.inner.CreateTicket(msg.GetArgs(ctx))
	return err
}

type SyncContactCommand struct {
	SyncTime time.Time

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleSyncContact(ctx context.Context, msg *SyncContactCommand) (err error) {
	return h.inner.SyncContact(msg.GetArgs(ctx))
}

type UpdateTicketCommand struct {
	FfmCode     string
	FfmID       int64
	ID          string
	EtopUserID  int64 `json:"etop_user_id"`
	Code        string
	Title       string
	Value       string
	OldValue    string
	Reason      string
	ShopID      int64
	OrderID     int64
	OrderCode   string
	FfmUrl      string
	Company     string
	Provider    string
	Note        string
	Environment string
	FromApp     string
	Account     Account

	Result *Ticket `json:"-"`
}

func (h AggregateHandler) HandleUpdateTicket(ctx context.Context, msg *UpdateTicketCommand) (err error) {
	msg.Result, err = h.inner.UpdateTicket(msg.GetArgs(ctx))
	return err
}

type CountTicketByStatusQuery struct {
	Status string

	Result *CountTicketByStatusResponse `json:"-"`
}

func (h QueryServiceHandler) HandleCountTicketByStatus(ctx context.Context, msg *CountTicketByStatusQuery) (err error) {
	msg.Result, err = h.inner.CountTicketByStatus(msg.GetArgs(ctx))
	return err
}

type GetCategoriesQuery struct {
	Result *GetCategoriesResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetCategories(ctx context.Context, msg *GetCategoriesQuery) (err error) {
	msg.Result, err = h.inner.GetCategories(msg.GetArgs(ctx))
	return err
}

type GetContactsQuery struct {
	Search string
	Paging *meta.Paging

	Result *ContactsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetContacts(ctx context.Context, msg *GetContactsQuery) (err error) {
	msg.Result, err = h.inner.GetContacts(msg.GetArgs(ctx))
	return err
}

type GetRecordLastTimeModifyQuery struct {
	Offset int32
	Limit  int32
	Sort   []string

	Result *Contact `json:"-"`
}

func (h QueryServiceHandler) HandleGetRecordLastTimeModify(ctx context.Context, msg *GetRecordLastTimeModifyQuery) (err error) {
	msg.Result, err = h.inner.GetRecordLastTimeModify(msg.GetArgs(ctx))
	return err
}

type GetTicketStatusCountQuery struct {
	Result *GetTicketStatusCountResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetTicketStatusCount(ctx context.Context, msg *GetTicketStatusCountQuery) (err error) {
	msg.Result, err = h.inner.GetTicketStatusCount(msg.GetArgs(ctx))
	return err
}

type GetTicketsQuery struct {
	Paging  *meta.Paging
	Ticket  TicketArgs
	Orderby OrderBy

	Result *GetTicketsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetTickets(ctx context.Context, msg *GetTicketsQuery) (err error) {
	msg.Result, err = h.inner.GetTickets(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateOrUpdateContactCommand) command() {}
func (q *CreateOrUpdateLeadCommand) command()    {}
func (q *CreateTicketCommand) command()          {}
func (q *SyncContactCommand) command()           {}
func (q *UpdateTicketCommand) command()          {}
func (q *CountTicketByStatusQuery) query()       {}
func (q *GetCategoriesQuery) query()             {}
func (q *GetContactsQuery) query()               {}
func (q *GetRecordLastTimeModifyQuery) query()   {}
func (q *GetTicketStatusCountQuery) query()      {}
func (q *GetTicketsQuery) query()                {}

// implement conversion

func (q *CreateOrUpdateContactCommand) GetArgs(ctx context.Context) (_ context.Context, _ *Contact) {
	return ctx,
		&Contact{
			ID:                   q.ID,
			EtopUserID:           q.EtopUserID,
			ContactNo:            q.ContactNo,
			Phone:                q.Phone,
			Lastname:             q.Lastname,
			Mobile:               q.Mobile,
			Leadsource:           q.Leadsource,
			Email:                q.Email,
			Description:          q.Description,
			Secondaryemail:       q.Secondaryemail,
			Modifiedby:           q.Modifiedby,
			Source:               q.Source,
			Company:              q.Company,
			Website:              q.Website,
			Lane:                 q.Lane,
			City:                 q.City,
			State:                q.State,
			Country:              q.Country,
			OrdersPerDay:         q.OrdersPerDay,
			UsedShippingProvider: q.UsedShippingProvider,
			Firstname:            q.Firstname,
			Createdtime:          q.Createdtime,
			Modifiedtime:         q.Modifiedtime,
			AssignedUserID:       q.AssignedUserID,
		}
}

func (q *CreateOrUpdateContactCommand) SetContact(args *Contact) {
	q.ID = args.ID
	q.EtopUserID = args.EtopUserID
	q.ContactNo = args.ContactNo
	q.Phone = args.Phone
	q.Lastname = args.Lastname
	q.Mobile = args.Mobile
	q.Leadsource = args.Leadsource
	q.Email = args.Email
	q.Description = args.Description
	q.Secondaryemail = args.Secondaryemail
	q.Modifiedby = args.Modifiedby
	q.Source = args.Source
	q.Company = args.Company
	q.Website = args.Website
	q.Lane = args.Lane
	q.City = args.City
	q.State = args.State
	q.Country = args.Country
	q.OrdersPerDay = args.OrdersPerDay
	q.UsedShippingProvider = args.UsedShippingProvider
	q.Firstname = args.Firstname
	q.Createdtime = args.Createdtime
	q.Modifiedtime = args.Modifiedtime
	q.AssignedUserID = args.AssignedUserID
}

func (q *CreateOrUpdateLeadCommand) GetArgs(ctx context.Context) (_ context.Context, _ *Lead) {
	return ctx,
		&Lead{
			ID:                   q.ID,
			EtopUserID:           q.EtopUserID,
			ContactNo:            q.ContactNo,
			Phone:                q.Phone,
			Lastname:             q.Lastname,
			Mobile:               q.Mobile,
			Leadsource:           q.Leadsource,
			Email:                q.Email,
			Description:          q.Description,
			Secondaryemail:       q.Secondaryemail,
			Modifiedby:           q.Modifiedby,
			Source:               q.Source,
			Company:              q.Company,
			Website:              q.Website,
			Lane:                 q.Lane,
			City:                 q.City,
			State:                q.State,
			Country:              q.Country,
			OrdersPerDay:         q.OrdersPerDay,
			UsedShippingProvider: q.UsedShippingProvider,
			Firstname:            q.Firstname,
			AssignedUserID:       q.AssignedUserID,
			Createdtime:          q.Createdtime,
			Modifiedtime:         q.Modifiedtime,
		}
}

func (q *CreateOrUpdateLeadCommand) SetLead(args *Lead) {
	q.ID = args.ID
	q.EtopUserID = args.EtopUserID
	q.ContactNo = args.ContactNo
	q.Phone = args.Phone
	q.Lastname = args.Lastname
	q.Mobile = args.Mobile
	q.Leadsource = args.Leadsource
	q.Email = args.Email
	q.Description = args.Description
	q.Secondaryemail = args.Secondaryemail
	q.Modifiedby = args.Modifiedby
	q.Source = args.Source
	q.Company = args.Company
	q.Website = args.Website
	q.Lane = args.Lane
	q.City = args.City
	q.State = args.State
	q.Country = args.Country
	q.OrdersPerDay = args.OrdersPerDay
	q.UsedShippingProvider = args.UsedShippingProvider
	q.Firstname = args.Firstname
	q.AssignedUserID = args.AssignedUserID
	q.Createdtime = args.Createdtime
	q.Modifiedtime = args.Modifiedtime
}

func (q *CreateTicketCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrUpdateTicketArgs) {
	return ctx,
		&CreateOrUpdateTicketArgs{
			FfmCode:     q.FfmCode,
			FfmID:       q.FfmID,
			ID:          q.ID,
			EtopUserID:  q.EtopUserID,
			Code:        q.Code,
			Title:       q.Title,
			Value:       q.Value,
			OldValue:    q.OldValue,
			Reason:      q.Reason,
			ShopID:      q.ShopID,
			OrderID:     q.OrderID,
			OrderCode:   q.OrderCode,
			FfmUrl:      q.FfmUrl,
			Company:     q.Company,
			Provider:    q.Provider,
			Note:        q.Note,
			Environment: q.Environment,
			FromApp:     q.FromApp,
			Account:     q.Account,
		}
}

func (q *CreateTicketCommand) SetCreateOrUpdateTicketArgs(args *CreateOrUpdateTicketArgs) {
	q.FfmCode = args.FfmCode
	q.FfmID = args.FfmID
	q.ID = args.ID
	q.EtopUserID = args.EtopUserID
	q.Code = args.Code
	q.Title = args.Title
	q.Value = args.Value
	q.OldValue = args.OldValue
	q.Reason = args.Reason
	q.ShopID = args.ShopID
	q.OrderID = args.OrderID
	q.OrderCode = args.OrderCode
	q.FfmUrl = args.FfmUrl
	q.Company = args.Company
	q.Provider = args.Provider
	q.Note = args.Note
	q.Environment = args.Environment
	q.FromApp = args.FromApp
	q.Account = args.Account
}

func (q *SyncContactCommand) GetArgs(ctx context.Context) (_ context.Context, _ *SyncContactArgs) {
	return ctx,
		&SyncContactArgs{
			SyncTime: q.SyncTime,
		}
}

func (q *SyncContactCommand) SetSyncContactArgs(args *SyncContactArgs) {
	q.SyncTime = args.SyncTime
}

func (q *UpdateTicketCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrUpdateTicketArgs) {
	return ctx,
		&CreateOrUpdateTicketArgs{
			FfmCode:     q.FfmCode,
			FfmID:       q.FfmID,
			ID:          q.ID,
			EtopUserID:  q.EtopUserID,
			Code:        q.Code,
			Title:       q.Title,
			Value:       q.Value,
			OldValue:    q.OldValue,
			Reason:      q.Reason,
			ShopID:      q.ShopID,
			OrderID:     q.OrderID,
			OrderCode:   q.OrderCode,
			FfmUrl:      q.FfmUrl,
			Company:     q.Company,
			Provider:    q.Provider,
			Note:        q.Note,
			Environment: q.Environment,
			FromApp:     q.FromApp,
			Account:     q.Account,
		}
}

func (q *UpdateTicketCommand) SetCreateOrUpdateTicketArgs(args *CreateOrUpdateTicketArgs) {
	q.FfmCode = args.FfmCode
	q.FfmID = args.FfmID
	q.ID = args.ID
	q.EtopUserID = args.EtopUserID
	q.Code = args.Code
	q.Title = args.Title
	q.Value = args.Value
	q.OldValue = args.OldValue
	q.Reason = args.Reason
	q.ShopID = args.ShopID
	q.OrderID = args.OrderID
	q.OrderCode = args.OrderCode
	q.FfmUrl = args.FfmUrl
	q.Company = args.Company
	q.Provider = args.Provider
	q.Note = args.Note
	q.Environment = args.Environment
	q.FromApp = args.FromApp
	q.Account = args.Account
}

func (q *CountTicketByStatusQuery) GetArgs(ctx context.Context) (_ context.Context, _ *CountTicketByStatusArgs) {
	return ctx,
		&CountTicketByStatusArgs{
			Status: q.Status,
		}
}

func (q *CountTicketByStatusQuery) SetCountTicketByStatusArgs(args *CountTicketByStatusArgs) {
	q.Status = args.Status
}

func (q *GetCategoriesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *GetCategoriesQuery) SetEmpty(args *meta.Empty) {
}

func (q *GetContactsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetContactsArgs) {
	return ctx,
		&GetContactsArgs{
			Search: q.Search,
			Paging: q.Paging,
		}
}

func (q *GetContactsQuery) SetGetContactsArgs(args *GetContactsArgs) {
	q.Search = args.Search
	q.Paging = args.Paging
}

func (q *GetRecordLastTimeModifyQuery) GetArgs(ctx context.Context) (_ context.Context, _ meta.Paging) {
	return ctx,
		meta.Paging{
			Offset: q.Offset,
			Limit:  q.Limit,
			Sort:   q.Sort,
		}
}

func (q *GetRecordLastTimeModifyQuery) SetPaging(args meta.Paging) {
	q.Offset = args.Offset
	q.Limit = args.Limit
	q.Sort = args.Sort
}

func (q *GetTicketStatusCountQuery) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *GetTicketStatusCountQuery) SetEmpty(args *meta.Empty) {
}

func (q *GetTicketsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTicketsArgs) {
	return ctx,
		&GetTicketsArgs{
			Paging:  q.Paging,
			Ticket:  q.Ticket,
			Orderby: q.Orderby,
		}
}

func (q *GetTicketsQuery) SetGetTicketsArgs(args *GetTicketsArgs) {
	q.Paging = args.Paging
	q.Ticket = args.Ticket
	q.Orderby = args.Orderby
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
	b.AddHandler(h.HandleCreateOrUpdateContact)
	b.AddHandler(h.HandleCreateOrUpdateLead)
	b.AddHandler(h.HandleCreateTicket)
	b.AddHandler(h.HandleSyncContact)
	b.AddHandler(h.HandleUpdateTicket)
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
	b.AddHandler(h.HandleCountTicketByStatus)
	b.AddHandler(h.HandleGetCategories)
	b.AddHandler(h.HandleGetContacts)
	b.AddHandler(h.HandleGetRecordLastTimeModify)
	b.AddHandler(h.HandleGetTicketStatusCount)
	b.AddHandler(h.HandleGetTickets)
	return QueryBus{b}
}
