package service

import (
	"context"
	"log"

	"github.com/gorilla/schema"

	pbcm "etop.vn/backend/pb/common"
	"etop.vn/backend/pb/services/crmservice"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/crm-service/mapping"
	vs "etop.vn/backend/pkg/crm-service/vtiger-service"
	wrapcrm "etop.vn/backend/wrapper/services/crmservice"
	"etop.vn/common/bus"
)

var encoder = schema.NewEncoder()

// Service contain config for service
type Service struct {
	vtigerService *vs.VtigerService
}

// NewService init Service
func NewService(db cmsql.Database, vConfig vs.VtigerConfig, fieldMap *mapping.ConfigMap) *Service {
	s := &Service{
		vtigerService: vs.NewSVtigerService(db, vConfig, fieldMap),
	}
	return s
}

// Register to handler
func (s *Service) Register() {
	bus.AddHandler("", s.VersionInfo)
	bus.AddHandler("", s.GetContacts)
	bus.AddHandler("", s.CreateOrUpdateContact)
	bus.AddHandler("", s.GetCategories)
	bus.AddHandler("", s.CreateOrUpdateTicket)
	bus.AddHandler("", s.GetTickets)
	bus.AddHandler("", s.CountTicketByStatus)
	bus.AddHandler("", s.CreateOrUpdateLead)
	bus.AddHandler("", s.GetTicketStatusCount)
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

// ReadFileCategories read file json
func ReadFileCategories() ([]*crmservice.Categories, error) {
	return vs.ReadFileCategories()
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

// CreateOrUpdateTicket create or uodate ticket
func (s *Service) CreateOrUpdateTicket(ctx context.Context, q *wrapcrm.CreateOrUpdateTicketEndpoint) error {
	result, err := s.vtigerService.CreateOrUpdateTicket(ctx, q.CreateOrUpdateTicketRequest)
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
