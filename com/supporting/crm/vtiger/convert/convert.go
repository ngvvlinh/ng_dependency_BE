package convert

import (
	"o.o/api/supporting/crm/vtiger"
	"o.o/backend/com/supporting/crm/vtiger/model"
)

// ConvertModelContact covert protobuf to model Contact
func ConvertModelContact(c *vtiger.Contact, AssignedUserID string) *model.VtigerContact {
	contact := &model.VtigerContact{
		ID:                   c.ID,
		Firstname:            c.Firstname,
		ContactNo:            c.ContactNo,
		Phone:                c.Phone,
		Description:          c.Description,
		Lastname:             c.Lastname,
		Mobile:               c.Mobile,
		Email:                c.Email,
		Leadsource:           c.Leadsource,
		Secondaryemail:       c.Secondaryemail,
		AssignedUserID:       AssignedUserID,
		EtopUserID:           c.EtopUserID,
		Source:               c.Source,
		UsedShippingProvider: c.UsedShippingProvider,
		OrdersPerDay:         c.OrdersPerDay,
		Company:              c.Company,
		City:                 c.City,
		State:                c.State,
		Website:              c.Website,
		Lane:                 c.Lane,
		Country:              c.Country,
		VtigerCreatedAt:      c.Createdtime,
		VtigerUpdatedAt:      c.Modifiedtime,
	}
	return contact
}

func ConvertContactFromModel(contact *model.VtigerContact) *vtiger.Contact {
	return &vtiger.Contact{
		ID:                   contact.ID,
		Firstname:            contact.Firstname,
		ContactNo:            contact.ContactNo,
		Phone:                contact.Phone,
		Lastname:             contact.Lastname,
		Mobile:               contact.Mobile,
		Email:                contact.Email,
		Leadsource:           contact.Leadsource,
		Secondaryemail:       contact.Secondaryemail,
		AssignedUserID:       contact.AssignedUserID,
		EtopUserID:           contact.EtopUserID,
		Source:               contact.Source,
		UsedShippingProvider: contact.UsedShippingProvider,
		OrdersPerDay:         contact.OrdersPerDay,
		Company:              contact.Company,
		City:                 contact.City,
		State:                contact.State,
		Website:              contact.Website,
		Lane:                 contact.Lane,
		Country:              contact.Country,
		Modifiedtime:         contact.VtigerUpdatedAt,
		Createdtime:          contact.VtigerCreatedAt,
	}
}

func ConvertLeadtoModelContact(lead *vtiger.Lead) *model.VtigerContact {
	contact := &model.VtigerContact{
		ID:                   lead.ID,
		Firstname:            lead.Firstname,
		ContactNo:            lead.ContactNo,
		Phone:                lead.Phone,
		Description:          lead.Description,
		Lastname:             lead.Lastname,
		Mobile:               lead.Mobile,
		Email:                lead.Email,
		Leadsource:           lead.Leadsource,
		Secondaryemail:       lead.Secondaryemail,
		AssignedUserID:       lead.AssignedUserID,
		EtopUserID:           lead.EtopUserID,
		Source:               lead.Source,
		UsedShippingProvider: lead.UsedShippingProvider,
		OrdersPerDay:         lead.OrdersPerDay,
		Company:              lead.Company,
		City:                 lead.City,
		State:                lead.State,
		Website:              lead.Website,
		Lane:                 lead.Lane,
		Country:              lead.Country,
		VtigerCreatedAt:      lead.Createdtime,
		VtigerUpdatedAt:      lead.Modifiedtime,
	}
	return contact
}

// ConvertAccount convert Account to Contact
func ConvertAccount(a *vtiger.Account) *vtiger.Contact {
	return &vtiger.Contact{
		EtopUserID: a.ID,
		Lastname:   a.FullName,
		Phone:      a.Phone,
		Email:      a.Email,
		Company:    a.Company,
	}
}

// ConvertTicket convert TicketRequest to Ticket protobuf. Ticket protobuf is used like DTO
func ConvertTicketGetReq(t *vtiger.TicketArgs) *vtiger.Ticket {
	ticket := &vtiger.Ticket{
		Code:        t.Code,
		TicketTitle: t.Title,
		NewValue:    t.Value,
		OldValue:    t.OldValue,
		Reason:      t.Reason,
		EtopUserID:  t.EtopUserID,
		OrderId:     t.OrderID,
		FfmCode:     t.FfmCode,
		FfmUrl:      t.FfmUrl,
		FfmId:       t.FfmID,
		Company:     t.Company,
		Provider:    t.Provider,
		Note:        t.Note,
		Environment: t.Environment,
		FromApp:     t.FromApp,
		ID:          t.ID,
	}
	return ticket
}

// ConvertTicket convert TicketRequest to Ticket protobuf. Ticket protobuf is used like DTO
func ConvertTicket(t *vtiger.CreateOrUpdateTicketArgs) *vtiger.Ticket {
	ticket := &vtiger.Ticket{
		Code:        t.Code,
		TicketTitle: t.Title,
		NewValue:    t.Value,
		OldValue:    t.OldValue,
		Reason:      t.Reason,
		EtopUserID:  t.EtopUserID,
		OrderId:     t.OrderID,
		FfmCode:     t.FfmCode,
		FfmUrl:      t.FfmUrl,
		FfmId:       t.FfmID,
		Company:     t.Company,
		Provider:    t.Provider,
		Note:        t.Note,
		Environment: t.Environment,
		FromApp:     t.FromApp,
		ID:          t.ID,
	}
	return ticket
}
