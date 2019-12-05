package crm

import (
	common "etop.vn/api/top/types/common"
	notifier_entity "etop.vn/api/top/types/etc/notifier_entity"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type RefreshFulfillmentFromCarrierRequest struct {
	// @required
	ShippingCode string `json:"shipping_code"`
}

func (m *RefreshFulfillmentFromCarrierRequest) Reset()         { *m = RefreshFulfillmentFromCarrierRequest{} }
func (m *RefreshFulfillmentFromCarrierRequest) String() string { return jsonx.MustMarshalToString(m) }

type SendNotificationRequest struct {
	AccountId dot.ID                         `json:"account_id"`
	Title     string                         `json:"title"`
	Message   string                         `json:"message"`
	MetaData  common.RawJSONObject           `json:"meta_data"`
	Entity    notifier_entity.NotifierEntity `json:"entity"`
	EntityId  dot.ID                         `json:"entity_id"`
}

func (m *SendNotificationRequest) Reset()         { *m = SendNotificationRequest{} }
func (m *SendNotificationRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCallHistoriesRequest struct {
	Paging     *common.Paging `json:"paging"`
	TextSearch string         `json:"text_search"`
}

func (m *GetCallHistoriesRequest) Reset()         { *m = GetCallHistoriesRequest{} }
func (m *GetCallHistoriesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCallHistoriesResponse struct {
	VhtCallLog []*VHTCallLog `json:"vht_call_log"`
}

func (m *GetCallHistoriesResponse) Reset()         { *m = GetCallHistoriesResponse{} }
func (m *GetCallHistoriesResponse) String() string { return jsonx.MustMarshalToString(m) }

type VHTCallLog struct {
	CdrId           string   `json:"cdr_id"`
	CallId          string   `json:"call_id"`
	SipCallId       string   `json:"sip_call_id"`
	SdkCallId       string   `json:"sdk_call_id"`
	Cause           string   `json:"cause"`
	Q850Cause       string   `json:"q850_cause"`
	FromExtension   string   `json:"from_extension"`
	ToExtension     string   `json:"to_extension"`
	FromNumber      string   `json:"from_number"`
	ToNumber        string   `json:"to_number"`
	Duration        int      `json:"duration"`
	Direction       int      `json:"direction"`
	TimeStarted     dot.Time `json:"time_started"`
	TimeConnected   dot.Time `json:"time_connected"`
	TimeEnded       dot.Time `json:"time_ended"`
	RecordingPath   string   `json:"recording_path"`
	RecordingUrl    string   `json:"recording_url"`
	RecordFileSize  int      `json:"record_file_size"`
	EtopAccountId   dot.ID   `json:"etop_account_id"`
	VtigerAccountId string   `json:"vtiger_account_id"`
}

func (m *VHTCallLog) Reset()         { *m = VHTCallLog{} }
func (m *VHTCallLog) String() string { return jsonx.MustMarshalToString(m) }

type CountTicketByStatusRequest struct {
	Status string `json:"status"`
}

func (m *CountTicketByStatusRequest) Reset()         { *m = CountTicketByStatusRequest{} }
func (m *CountTicketByStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketStatusCountResponse struct {
	StatusCount []*CountTicketByStatusResponse `json:"status_count"`
}

func (m *GetTicketStatusCountResponse) Reset()         { *m = GetTicketStatusCountResponse{} }
func (m *GetTicketStatusCountResponse) String() string { return jsonx.MustMarshalToString(m) }

type CountTicketByStatusResponse struct {
	Code  string `json:"code"`
	Count int    `json:"count"`
}

func (m *CountTicketByStatusResponse) Reset()         { *m = CountTicketByStatusResponse{} }
func (m *CountTicketByStatusResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetContactsResponse struct {
	Contacts []*ContactResponse `json:"contacts"`
}

func (m *GetContactsResponse) Reset()         { *m = GetContactsResponse{} }
func (m *GetContactsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetContactsRequest struct {
	TextSearch string         `json:"text_search"`
	Paging     *common.Paging `json:"paging"`
}

func (m *GetContactsRequest) Reset()         { *m = GetContactsRequest{} }
func (m *GetContactsRequest) String() string { return jsonx.MustMarshalToString(m) }

type ContactRequest struct {
	ContactNo            string   `json:"contact_no"`
	Phone                string   `json:"phone"`
	Lastname             string   `json:"lastname"`
	Mobile               string   `json:"mobile"`
	Leadsource           string   `json:"leadsource"`
	Email                string   `json:"email"`
	Description          string   `json:"description"`
	Secondaryemail       string   `json:"secondaryemail"`
	Modifiedby           string   `json:"modifiedby"`
	Source               string   `json:"source"`
	EtopUserId           dot.ID   `json:"etop_user_id"`
	Company              string   `json:"company"`
	Website              string   `json:"website"`
	Lane                 string   `json:"lane"`
	City                 string   `json:"city"`
	State                string   `json:"state"`
	Country              string   `json:"country"`
	OrdersPerDay         string   `json:"orders_per_day"`
	UsedShippingProvider string   `json:"used_shipping_provider"`
	Id                   string   `json:"id"`
	Firstname            string   `json:"firstname"`
	Createdtime          dot.Time `json:"createdtime"`
	Modifiedtime         dot.Time `json:"modifiedtime"`
	AssignedUserId       string   `json:"assigned_user_id"`
}

func (m *ContactRequest) Reset()         { *m = ContactRequest{} }
func (m *ContactRequest) String() string { return jsonx.MustMarshalToString(m) }

type ContactResponse struct {
	ContactNo            string   `json:"contact_no"`
	Phone                string   `json:"phone"`
	Lastname             string   `json:"lastname"`
	Mobile               string   `json:"mobile"`
	Leadsource           string   `json:"leadsource"`
	Email                string   `json:"email"`
	Description          string   `json:"description"`
	Secondaryemail       string   `json:"secondaryemail"`
	Modifiedby           string   `json:"modifiedby"`
	Source               string   `json:"source"`
	EtopUserId           dot.ID   `json:"etop_user_id"`
	Company              string   `json:"company"`
	Website              string   `json:"website"`
	Lane                 string   `json:"lane"`
	City                 string   `json:"city"`
	State                string   `json:"state"`
	Country              string   `json:"country"`
	OrdersPerDay         string   `json:"orders_per_day"`
	UsedShippingProvider string   `json:"used_shipping_provider"`
	Id                   string   `json:"id"`
	Firstname            string   `json:"firstname"`
	Createdtime          dot.Time `json:"createdtime"`
	Modifiedtime         dot.Time `json:"modifiedtime"`
	AssignedUserId       string   `json:"assigned_user_id"`
}

func (m *ContactResponse) Reset()         { *m = ContactResponse{} }
func (m *ContactResponse) String() string { return jsonx.MustMarshalToString(m) }

type LeadRequest struct {
	ContactNo            string `json:"contact_no"`
	Phone                string `json:"phone"`
	Lastname             string `json:"lastname"`
	Mobile               string `json:"mobile"`
	Leadsource           string `json:"leadsource"`
	Email                string `json:"email"`
	Secondaryemail       string `json:"secondaryemail"`
	AssignedUserId       string `json:"assigned_user_id"`
	Description          string `json:"description"`
	Modifiedby           string `json:"modifiedby"`
	Source               string `json:"source"`
	EtopUserId           dot.ID `json:"etop_user_id"`
	Company              string `json:"company"`
	Website              string `json:"website"`
	Lane                 string `json:"lane"`
	City                 string `json:"city"`
	State                string `json:"state"`
	Country              string `json:"country"`
	OrdersPerDay         string `json:"orders_per_day"`
	UsedShippingProvider string `json:"used_shipping_provider"`
	Id                   string `json:"id"`
	Firstname            string `json:"firstname"`
}

func (m *LeadRequest) Reset()         { *m = LeadRequest{} }
func (m *LeadRequest) String() string { return jsonx.MustMarshalToString(m) }

type LeadResponse struct {
	ContactNo            string `json:"contact_no"`
	Phone                string `json:"phone"`
	Lastname             string `json:"lastname"`
	Mobile               string `json:"mobile"`
	Leadsource           string `json:"leadsource"`
	Email                string `json:"email"`
	Secondaryemail       string `json:"secondaryemail"`
	AssignedUserId       string `json:"assigned_user_id"`
	Description          string `json:"description"`
	Modifiedby           string `json:"modifiedby"`
	Source               string `json:"source"`
	EtopUserId           dot.ID `json:"etop_user_id"`
	Company              string `json:"company"`
	Website              string `json:"website"`
	Lane                 string `json:"lane"`
	City                 string `json:"city"`
	State                string `json:"state"`
	Country              string `json:"country"`
	OrdersPerDay         string `json:"orders_per_day"`
	UsedShippingProvider string `json:"used_shipping_provider"`
	Id                   string `json:"id"`
	Firstname            string `json:"firstname"`
}

func (m *LeadResponse) Reset()         { *m = LeadResponse{} }
func (m *LeadResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketsRequest struct {
	Paging  *common.Paging `json:"paging"`
	Ticket  TicketRequest  `json:"ticket"`
	Orderby OrderBy        `json:"orderby"`
}

func (m *GetTicketsRequest) Reset()         { *m = GetTicketsRequest{} }
func (m *GetTicketsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketsResponse struct {
	Tickets []*Ticket `json:"tickets"`
}

func (m *GetTicketsResponse) Reset()         { *m = GetTicketsResponse{} }
func (m *GetTicketsResponse) String() string { return jsonx.MustMarshalToString(m) }

type OrderBy struct {
	Field string `json:"field"`
	Sort  string `json:"sort"`
}

func (m *OrderBy) Reset()         { *m = OrderBy{} }
func (m *OrderBy) String() string { return jsonx.MustMarshalToString(m) }

type CreateOrUpdateTicketRequest struct {
	Id          string  `json:"id"`
	Code        string  `json:"code"`
	Title       string  `json:"title"`
	Value       string  `json:"value"`
	OldValue    string  `json:"old_value"`
	Reason      string  `json:"reason"`
	EtopUserId  dot.ID  `json:"etop_user_id"`
	OrderId     dot.ID  `json:"order_id"`
	OrderCode   string  `json:"order_code"`
	FfmCode     string  `json:"ffm_code"`
	FfmUrl      string  `json:"ffm_url"`
	FfmId       dot.ID  `json:"ffm_id"`
	Company     string  `json:"company"`
	Provider    string  `json:"provider"`
	Note        string  `json:"note"`
	Environment string  `json:"environment"`
	FromApp     string  `json:"from_app"`
	Account     Account `json:"account"`
}

func (m *CreateOrUpdateTicketRequest) Reset()         { *m = CreateOrUpdateTicketRequest{} }
func (m *CreateOrUpdateTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type TicketRequest struct {
	Id          string `json:"id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Value       string `json:"value"`
	OldValue    string `json:"old_value"`
	Reason      string `json:"reason"`
	EtopUserId  dot.ID `json:"etop_user_id"`
	OrderId     dot.ID `json:"order_id"`
	OrderCode   string `json:"order_code"`
	FfmCode     string `json:"ffm_code"`
	FfmUrl      string `json:"ffm_url"`
	FfmId       dot.ID `json:"ffm_id"`
	Company     string `json:"company"`
	Provider    string `json:"provider"`
	Note        string `json:"note"`
	Environment string `json:"environment"`
	FromApp     string `json:"from_app"`
}

func (m *TicketRequest) Reset()         { *m = TicketRequest{} }
func (m *TicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type Account struct {
	Id        dot.ID `json:"id"`
	FullName  string `json:"full_name"`
	ShortName string `json:"short_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Company   string `json:"company"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return jsonx.MustMarshalToString(m) }

type Ticket struct {
	TicketNo         string   `json:"ticket_no"`
	AssignedUserId   string   `json:"assigned_user_id"`
	ParentId         dot.ID   `json:"parent_id"`
	Ticketpriorities string   `json:"ticketpriorities"`
	ProductId        dot.ID   `json:"product_id"`
	Ticketseverities string   `json:"ticketseverities"`
	Ticketstatus     string   `json:"ticketstatus"`
	Ticketcategories string   `json:"ticketcategories"`
	UpdateLog        string   `json:"update_log"`
	Hours            string   `json:"hours"`
	Days             string   `json:"days"`
	Createdtime      dot.Time `json:"createdtime"`
	Modifiedtime     dot.Time `json:"modifiedtime"`
	FromPortal       string   `json:"from_portal"`
	Modifiedby       string   `json:"modifiedby"`
	TicketTitle      string   `json:"ticket_title"`
	Description      string   `json:"description"`
	Solution         string   `json:"solution"`
	ContactId        string   `json:"contact_id"`
	Source           string   `json:"source"`
	Starred          string   `json:"starred"`
	Tags             string   `json:"tags"`
	Note             string   `json:"note"`
	FfmCode          string   `json:"ffm_code"`
	FfmUrl           string   `json:"ffm_url"`
	FfmId            dot.ID   `json:"ffm_id"`
	EtopUserId       dot.ID   `json:"etop_user_id"`
	OrderId          dot.ID   `json:"order_id"`
	OrderCode        string   `json:"order_code"`
	Company          string   `json:"company"`
	Provider         string   `json:"provider"`
	FromApp          string   `json:"from_app"`
	Environment      string   `json:"environment"`
	Code             string   `json:"code"`
	OldValue         string   `json:"old_value"`
	NewValue         string   `json:"new_value"`
	Substatus        string   `json:"substatus"`
	EtopNote         string   `json:"etop_note"`
	Reason           string   `json:"reason"`
	Id               string   `json:"id"`
}

func (m *Ticket) Reset()         { *m = Ticket{} }
func (m *Ticket) String() string { return jsonx.MustMarshalToString(m) }

type GetCategoriesResponse struct {
	Categories []*Category `json:"categories"`
}

func (m *GetCategoriesResponse) Reset()         { *m = GetCategoriesResponse{} }
func (m *GetCategoriesResponse) String() string { return jsonx.MustMarshalToString(m) }

type Category struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

func (m *Category) Reset()         { *m = Category{} }
func (m *Category) String() string { return jsonx.MustMarshalToString(m) }

type GetStatusResponse struct {
	Status []*Status `json:"status"`
}

func (m *GetStatusResponse) Reset()         { *m = GetStatusResponse{} }
func (m *GetStatusResponse) String() string { return jsonx.MustMarshalToString(m) }

type Status struct {
	Code  string `json:"code"`
	Count string `json:"count"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return jsonx.MustMarshalToString(m) }
