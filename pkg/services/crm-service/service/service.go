package service

import (
	"context"
	"log"

	pbcm "etop.vn/backend/pb/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/crm-service/mapping"
	"etop.vn/backend/pkg/services/crm-service/vht/service"
	vtigerservice "etop.vn/backend/pkg/services/crm-service/vtiger/service"
	wrapcrm "etop.vn/backend/wrapper/services/crmservice"
	"etop.vn/common/bus"
)

// Service contain config for service
type Service struct {
	vtigerService *vtigerservice.VtigerService
	vhtservice    *service.VhtService
}

// NewService init Service
func NewService(db cmsql.Database, vConfig vtigerservice.Config, fieldMap mapping.ConfigMap) *Service {
	s := &Service{
		vtigerService: vtigerservice.NewSVtigerService(db, vConfig, fieldMap),
		vhtservice:    service.NewVhtService(db),
	}
	return s
}

// Register to handler
func (s *Service) Register() {
	bus.AddHandler("", s.VersionInfo)
	bus.AddHandler("", s.GetContacts)
	bus.AddHandler("", s.CreateOrUpdateContact)
	bus.AddHandler("", s.GetCategories)
	bus.AddHandler("", s.CreateTicket)
	bus.AddHandler("", s.UpdateTicket)
	bus.AddHandler("", s.GetTickets)
	bus.AddHandler("", s.CountTicketByStatus)
	bus.AddHandler("", s.CreateOrUpdateLead)
	bus.AddHandler("", s.GetTicketStatusCount)
	bus.AddHandler("", s.GetCallHistories)
	bus.AddHandler("", s.CreateOrUpdateCallHistoryBySDKCallID)
}

// VersionInfo get version info
func (s *Service) VersionInfo(ctx context.Context, q *wrapcrm.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Shop",
		Version: "0.1",
	}
	log.Println("change")
	return nil
}

//CreateOrUpdateCallHistoryBySDKCallID Create Or Update Call History By SDKCallID
func (s *Service) CreateOrUpdateCallHistoryByCallID(ctx context.Context, q *wrapcrm.CreateOrUpdateCallHistoryByCallIDEndpoint) error {
	result, err := s.vhtservice.CreateOrUpdateCallHistoryByCallID(ctx, q.VHTCallLog)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

//CreateOrUpdateCallHistoryBySDKCallID Create Or Update Call History By SDKCallID
func (s *Service) CreateOrUpdateCallHistoryBySDKCallID(ctx context.Context, q *wrapcrm.CreateOrUpdateCallHistoryBySDKCallIDEndpoint) error {
	result, err := s.vhtservice.CreateOrUpdateBySDKCallID(ctx, q.VHTCallLog)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

// GetCallHistories get call histories
func (s *Service) GetCallHistories(ctx context.Context, q *wrapcrm.GetCallHistoriesEndpoint) error {
	result, err := s.vhtservice.GetCallHistories(ctx, q.GetCallHistoriesRequest)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

// GetCategories get categories
func (s *Service) GetTicketStatusCount(ctx context.Context, q *wrapcrm.GetTicketStatusCountEndpoint) error {
	result, err := s.vtigerService.GetTicketStatusCount(ctx)
	if err != nil {
		return err
	}
	q.Result = result
	return nil

}

// GetCategories get categories
func (s *Service) GetCategories(ctx context.Context, q *wrapcrm.GetCategoriesEndpoint) error {
	result, err := s.vtigerService.GetCategories(ctx)
	if err != nil {
		return err
	}
	q.Result = result
	return nil

}

// CountTicketByStatus  count ticket by status which is code from reason map
func (s *Service) CountTicketByStatus(ctx context.Context, q *wrapcrm.CountTicketByStatusEndpoint) error {
	result, err := s.vtigerService.CountTicketByStatus(ctx, q.CountTicketByStatusRequest)
	if err != nil {
		return nil
	}
	q.Result = result
	return nil
}

// CreateOrUpdateLead create or update lead
func (s *Service) CreateOrUpdateLead(ctx context.Context, q *wrapcrm.CreateOrUpdateLeadEndpoint) error {
	result, err := s.vtigerService.CreateOrUpdateLead(ctx, q.Lead)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

// GetTickets get ticket from vtiger
func (s *Service) GetTickets(ctx context.Context, q *wrapcrm.GetTicketsEndpoint) error {
	result, err := s.vtigerService.GetTickets(ctx, q.GetTicketsRequest)
	if err != nil {
		return nil
	}
	q.Result = result
	return nil
}

// CreateTicket create a ticket
func (s *Service) CreateTicket(ctx context.Context, q *wrapcrm.CreateTicketEndpoint) error {
	result, err := s.vtigerService.CreateOrUpdateTicket(ctx, q.CreateOrUpdateTicketRequest, "create")
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

// UpdateTicket update a ticket
func (s *Service) UpdateTicket(ctx context.Context, q *wrapcrm.UpdateTicketEndpoint) error {
	result, err := s.vtigerService.CreateOrUpdateTicket(ctx, q.CreateOrUpdateTicketRequest, "update")
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

// GetContacts get contact from db
func (s *Service) GetContacts(ctx context.Context, q *wrapcrm.GetContactsEndpoint) error {
	result, err := s.vtigerService.GetContacts(ctx, q.GetContactsRequest)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}

// CreateOrUpdateContact .
func (s *Service) CreateOrUpdateContact(ctx context.Context, q *wrapcrm.CreateOrUpdateContactEndpoint) error {
	result, err := s.vtigerService.CreateOrUpdateContact(ctx, q.Contact)
	if err != nil {
		return err
	}
	q.Result = result
	return nil
}
