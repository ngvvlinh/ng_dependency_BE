package vtigerservice

import (
	"context"

	"etop.vn/backend/pb/services/crmservice"
	"etop.vn/backend/pkg/crm-service/mapping"
	"etop.vn/backend/pkg/crm-service/model"
	"etop.vn/backend/pkg/crm-service/vtiger"
)

// CreateOrUpdateLead create or update lead
func (s *VtigerService) CreateOrUpdateLead(ctx context.Context, req *crmservice.Lead) (*crmservice.Lead, error) {
	// get session
	session, err := vtiger.GetSessionKey(s.cfg.VtigerService, s.cfg.VtigerUsername, s.cfg.VtigerAccesskey)
	if err != nil {
		return nil, err
	}

	// save value to database.
	contact := &model.VtigerContact{
		ID:                   req.Id,
		Firstname:            req.Firstname,
		ContactNo:            req.ContactNo,
		Phone:                req.Phone,
		Description:          req.Description,
		Lastname:             req.Lastname,
		Mobile:               req.Mobile,
		Email:                req.Email,
		Leadsource:           req.Leadsource,
		Secondaryemail:       req.Secondaryemail,
		AssignedUserID:       req.AssignedUserId,
		EtopID:               req.EtopId,
		Source:               req.Source,
		UsedShippingProvider: req.UsedShippingProvider,
		OrdersPerDay:         req.OrdersPerDay,
		Company:              req.Company,
		City:                 req.City,
		State:                req.State,
		Website:              req.Website,
		Lane:                 req.Lane,
		Country:              req.Country,
	}
	query := s.vtigerContact(ctx).ByEtopID(contact.EtopID)
	_, err = query.GetContact()
	if err == nil {
		err = query.UpdateVtigerContact(contact)
		if err != nil {
			return nil, err
		}
	} else {
		err = query.CreateVtigerContact(contact)
		if err != nil {
			return nil, err
		}
	}
	//s.UpdateOrCreateContactToVtiger(contact)
	// send value to vtiger service
	fileMapData := mapping.NewMappingConfigInfo(s.fieldMap)
	vtigerMap, err := fileMapData.MappingLeadEtop2Vtiger(req)
	if err != nil {
		return nil, err
	}
	leadResp, err := s.CreateOrUpdateVtiger(vtigerMap, session, fileMapData, "Leads")
	if err != nil {
		return nil, err
	}
	leadReturn, err := fileMapData.MappingLeadVtiger2Etop(leadResp)
	return leadReturn, nil
}
