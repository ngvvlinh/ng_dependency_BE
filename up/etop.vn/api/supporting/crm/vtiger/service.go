package vtiger

import (
	"context"
	"time"

	"etop.vn/api/meta"
)

// +gen:api

type Aggregate interface {
	CreateOrUpdateContact(context.Context, *Contact) (*Contact, error)
	CreateOrUpdateLead(context.Context, *Lead) (*Lead, error)
	CreateTicket(context.Context, *CreateOrUpdateTicketArgs) (*Ticket, error)
	UpdateTicket(context.Context, *CreateOrUpdateTicketArgs) (*Ticket, error)
	SyncContact(context.Context, *SyncContactArgs) error
}

type QueryService interface {
	GetContacts(context.Context, *GetContactsArgs) (*ContactsResponse, error)
	GetTickets(context.Context, *GetTicketsArgs) (*GetTicketsResponse, error)
	GetCategories(context.Context, *meta.Empty) (*GetCategoriesResponse, error)
	CountTicketByStatus(context.Context, *CountTicketByStatusArgs) (*CountTicketByStatusResponse, error)
	GetTicketStatusCount(context.Context, *meta.Empty) (*GetTicketStatusCountResponse, error)
	GetRecordLastTimeModify(context.Context, meta.Paging) (*Contact, error)
}

type GetTicketsResponse struct {
	Tickets []*Ticket
}

type GetContactsArgs struct {
	Search string
	Paging *meta.Paging
}

type GetTicketsArgs struct {
	Paging  *meta.Paging
	Ticket  TicketArgs
	Orderby OrderBy
}

type OrderBy struct {
	Field string
	Sort  string
}

type TicketArgs struct {
	ID          string
	EtopUserID  int64 `json:"etop_user_id"`
	shopID      int64
	Code        string
	Title       string
	Value       string
	OldValue    string
	Reason      string
	ShopID      int64
	OrderID     int64
	OrderCode   string
	FfmCode     string
	FfmUrl      string
	FfmID       int64
	Company     string
	Provider    string
	Note        string
	Environment string
	FromApp     string
}

type ContactsResponse struct {
	Contacts []*Contact
}

type Contact struct {
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
}

type Lead struct {
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
}

type Ticket struct {
	ID               string    `json:"id"`
	EtopUserID       int64     `json:"etop_user_id"`
	ShopID           int64     `json:"shop_id"`
	TicketNo         string    `json:"ticket_no"`
	ParentID         int64     `json:"parent_id"`
	Ticketpriorities string    `json:"ticketpriorities"`
	ProductID        int64     `json:"productID"`
	Ticketseverities string    `json:"ticketseverities"`
	Ticketstatus     string    `json:"ticketstatus"`
	Ticketcategories string    `json:"ticketcategories"`
	UpdateLog        string    `json:"updateLog"`
	Hours            string    `json:"hours"`
	Days             string    `json:"days"`
	CreatedTime      time.Time `json:"created_time"`
	ModifiedTime     time.Time `json:"modified_time"`
	FromPortal       string    `json:"from_portal"`
	Modifiedby       string    `json:"modifiedby"`
	TicketTitle      string    `json:"ticket_title"`
	Description      string    `json:"description"`
	Solution         string    `json:"solution"`
	ContactId        string    `json:"contactId"`
	Source           string    `json:"source"`
	Starred          string    `json:"starred"`
	Tags             string    `json:"tags"`
	Note             string    `json:"note"`
	FfmCode          string    `json:"ffm_code"`
	FfmUrl           string    `json:"ffm_url"`
	FfmId            int64     `json:"ffm_id"`
	OrderId          int64     `json:"order_id"`
	OrderCode        string    `json:"order_code"`
	Company          string    `json:"company"`
	Provider         string    `json:"provider"`
	FromApp          string    `json:"from_app"`
	Environment      string    `json:"environment"`
	Code             string    `json:"code"`
	OldValue         string    `json:"old_value"`
	NewValue         string    `json:"new_value"`
	Substatus        string    `json:"substatus"`
	EtopNote         string    `json:"etopNote"`
	Reason           string    `json:"reason"`
	AssignedUserId   string    `json:"assigned_user_id"`
}

type CreateOrUpdateTicketArgs struct {
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
}

type Account struct {
	ID        int64
	FullName  string
	ShortName string
	Phone     string
	Email     string
	Company   string
}

type GetCategoriesResponse struct {
	Categories []*Category
}
type Category struct {
	Code  string
	Label string
}

type GetStatusResponse struct {
	Status []*Status
}

type Status struct {
	Code  string
	Count int64
}

type CountTicketByStatusArgs struct {
	Status string
}

type CountTicketByStatusResponse struct {
	Code  string
	Count int32
}

type GetTicketStatusCountResponse struct {
	StatusCount []CountTicketByStatusResponse
}

type SyncContactArgs struct {
	SyncTime time.Time
}
