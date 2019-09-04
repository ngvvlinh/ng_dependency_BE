// Code generated by gen-cmd-query. DO NOT EDIT.

package vtiger

import (
	context "context"
	time "time"

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

type SyncContactCommand struct {
	SyncTime time.Time

	Result struct {
	} `json:"-"`
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

type CountTicketByStatusQuery struct {
	Status string

	Result *CountTicketByStatusResponse `json:"-"`
}

type GetCategoriesQuery struct {
	Result *GetCategoriesResponse `json:"-"`
}

type GetContactsQuery struct {
	Search string
	Paging *meta.Paging

	Result *ContactsResponse `json:"-"`
}

type GetLastTimeModifyQuery struct {
	Offset int32
	Limit  int32
	Sort   []string

	Result *Contact `json:"-"`
}

type GetTicketStatusCountQuery struct {
	Result *GetTicketStatusCountResponse `json:"-"`
}

type GetTicketsQuery struct {
	Paging  *meta.Paging
	Ticket  TicketArgs
	Orderby OrderBy

	Result *GetTicketsResponse `json:"-"`
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
func (q *GetLastTimeModifyQuery) query()         {}
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

func (q *SyncContactCommand) GetArgs(ctx context.Context) (_ context.Context, _ *SyncContactArgs) {
	return ctx,
		&SyncContactArgs{
			SyncTime: q.SyncTime,
		}
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

func (q *CountTicketByStatusQuery) GetArgs(ctx context.Context) (_ context.Context, _ *CountTicketByStatusArgs) {
	return ctx,
		&CountTicketByStatusArgs{
			Status: q.Status,
		}
}

func (q *GetCategoriesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *GetContactsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetContactsArgs) {
	return ctx,
		&GetContactsArgs{
			Search: q.Search,
			Paging: q.Paging,
		}
}

func (q *GetLastTimeModifyQuery) GetArgs(ctx context.Context) (_ context.Context, _ meta.Paging) {
	return ctx,
		meta.Paging{
			Offset: q.Offset,
			Limit:  q.Limit,
			Sort:   q.Sort,
		}
}

func (q *GetTicketStatusCountQuery) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *GetTicketsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTicketsArgs) {
	return ctx,
		&GetTicketsArgs{
			Paging:  q.Paging,
			Ticket:  q.Ticket,
			Orderby: q.Orderby,
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
	b.AddHandler(h.HandleCreateOrUpdateContact)
	b.AddHandler(h.HandleCreateOrUpdateLead)
	b.AddHandler(h.HandleCreateTicket)
	b.AddHandler(h.HandleSyncContact)
	b.AddHandler(h.HandleUpdateTicket)
	return CommandBus{b}
}

func (h AggregateHandler) HandleCreateOrUpdateContact(ctx context.Context, msg *CreateOrUpdateContactCommand) error {
	result, err := h.inner.CreateOrUpdateContact(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleCreateOrUpdateLead(ctx context.Context, msg *CreateOrUpdateLeadCommand) error {
	result, err := h.inner.CreateOrUpdateLead(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleCreateTicket(ctx context.Context, msg *CreateTicketCommand) error {
	result, err := h.inner.CreateTicket(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleSyncContact(ctx context.Context, msg *SyncContactCommand) error {
	return h.inner.SyncContact(msg.GetArgs(ctx))
}

func (h AggregateHandler) HandleUpdateTicket(ctx context.Context, msg *UpdateTicketCommand) error {
	result, err := h.inner.UpdateTicket(msg.GetArgs(ctx))
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
	b.AddHandler(h.HandleCountTicketByStatus)
	b.AddHandler(h.HandleGetCategories)
	b.AddHandler(h.HandleGetContacts)
	b.AddHandler(h.HandleGetLastTimeModify)
	b.AddHandler(h.HandleGetTicketStatusCount)
	b.AddHandler(h.HandleGetTickets)
	return QueryBus{b}
}

func (h QueryServiceHandler) HandleCountTicketByStatus(ctx context.Context, msg *CountTicketByStatusQuery) error {
	result, err := h.inner.CountTicketByStatus(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetCategories(ctx context.Context, msg *GetCategoriesQuery) error {
	result, err := h.inner.GetCategories(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetContacts(ctx context.Context, msg *GetContactsQuery) error {
	result, err := h.inner.GetContacts(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetLastTimeModify(ctx context.Context, msg *GetLastTimeModifyQuery) error {
	result, err := h.inner.GetRecordLastTimeModify(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetTicketStatusCount(ctx context.Context, msg *GetTicketStatusCountQuery) error {
	result, err := h.inner.GetTicketStatusCount(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetTickets(ctx context.Context, msg *GetTicketsQuery) error {
	result, err := h.inner.GetTickets(msg.GetArgs(ctx))
	msg.Result = result
	return err
}
