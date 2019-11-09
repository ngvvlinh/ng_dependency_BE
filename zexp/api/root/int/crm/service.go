package crm

import (
	"context"

	cm "etop.vn/backend/pb/common"
	crm "etop.vn/backend/pb/services/crm"
)

// +gen:apix

// +apix:path=/crm.Misc
type MiscAPI interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/crm.Crm
type CrmAPI interface {
	RefreshFulfillmentFromCarrier(context.Context, *crm.RefreshFulfillmentFromCarrierRequest) (*cm.UpdatedResponse, error)

	SendNotification(context.Context, *crm.SendNotificationRequest) (*cm.MessageResponse, error)
}

// +apix:path=/crm.Vtiger
type VtigerAPI interface {
	GetContacts(context.Context, *crm.GetContactsRequest) (*crm.GetContactsResponse, error)
	CreateOrUpdateContact(context.Context, *crm.ContactRequest) (*crm.ContactResponse, error)
	CreateOrUpdateLead(context.Context, *crm.LeadRequest) (*crm.LeadResponse, error)
	GetTickets(context.Context, *crm.GetTicketsRequest) (*crm.GetTicketsResponse, error)
	CreateTicket(context.Context, *crm.CreateOrUpdateTicketRequest) (*crm.Ticket, error)
	UpdateTicket(context.Context, *crm.CreateOrUpdateTicketRequest) (*crm.Ticket, error)
	GetCategories(context.Context, *cm.Empty) (*crm.GetCategoriesResponse, error)
	// GetStatus(context.Context, *cm.Empty) (*crm.GetStatusResponse, error)
	// 	CountTicketByStatus(context.Context, *crm.CountTicketByStatusRequest) (*crm.CountTicketByStatusResponse, error)
	GetTicketStatusCount(context.Context, *cm.Empty) (*crm.GetTicketStatusCountResponse, error)
}

// +apix:path=/crm.Vht
type VhtAPI interface {
	GetCallHistories(context.Context, *crm.GetCallHistoriesRequest) (*crm.GetCallHistoriesResponse, error)
	CreateOrUpdateCallHistoryBySDKCallID(context.Context, *crm.VHTCallLog) (*crm.VHTCallLog, error)
	CreateOrUpdateCallHistoryByCallID(context.Context, *crm.VHTCallLog) (*crm.VHTCallLog, error)
}
