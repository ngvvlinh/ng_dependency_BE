package crm

import (
	"context"

	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=services/crm

// +apix:path=/crm.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/crm.Crm
type CrmService interface {
	RefreshFulfillmentFromCarrier(context.Context, *RefreshFulfillmentFromCarrierRequest) (*cm.UpdatedResponse, error)

	SendNotification(context.Context, *SendNotificationRequest) (*cm.MessageResponse, error)
}

// +apix:path=/crm.Vtiger
type VtigerService interface {
	GetContacts(context.Context, *GetContactsRequest) (*GetContactsResponse, error)
	CreateOrUpdateContact(context.Context, *ContactRequest) (*ContactResponse, error)
	CreateOrUpdateLead(context.Context, *LeadRequest) (*LeadResponse, error)
	GetTickets(context.Context, *GetTicketsRequest) (*GetTicketsResponse, error)
	CreateTicket(context.Context, *CreateOrUpdateTicketRequest) (*Ticket, error)
	UpdateTicket(context.Context, *CreateOrUpdateTicketRequest) (*Ticket, error)
	GetCategories(context.Context, *cm.Empty) (*GetCategoriesResponse, error)
	// GetStatus(context.Context, *cm.Empty) (*crm.GetStatusResponse, error)
	// 	CountTicketByStatus(context.Context, *crm.CountTicketByStatusRequest) (*crm.CountTicketByStatusResponse, error)
	GetTicketStatusCount(context.Context, *cm.Empty) (*GetTicketStatusCountResponse, error)
}

// +apix:path=/crm.Vht
type VhtService interface {
	GetCallHistories(context.Context, *GetCallHistoriesRequest) (*GetCallHistoriesResponse, error)
	CreateOrUpdateCallHistoryBySDKCallID(context.Context, *VHTCallLog) (*VHTCallLog, error)
	CreateOrUpdateCallHistoryByCallID(context.Context, *VHTCallLog) (*VHTCallLog, error)
}
