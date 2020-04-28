package aggregate

import (
	"context"
	"fmt"
	"strings"

	crmvtiger "o.o/api/supporting/crm/vtiger"
	"o.o/backend/com/supporting/crm/vtiger/convert"
	"o.o/backend/com/supporting/crm/vtiger/mapping"
	vtigermapping "o.o/backend/com/supporting/crm/vtiger/mapping"
	"o.o/backend/com/supporting/crm/vtiger/sqlstore"
	"o.o/backend/com/supporting/crm/vtiger/sync"
	"o.o/backend/com/supporting/crm/vtiger/vtigerstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	vtigerclient "o.o/backend/pkg/integration/vtiger/client"
)

var (
	Empty      = ""
	Categories = []*crmvtiger.Category{
		{
			Code:  "force-picking",
			Label: "Giục lấy hàng",
		},
		{
			Code:  "force-delivering",
			Label: "Giục giao hàng",
		},
		{
			Code:  "change-shop-cod",
			Label: "Thay đổi COD",
		},
		{
			Code:  "change-phone",
			Label: "Thay đổi SDT",
		},
		{
			Code:  "change-customer-name",
			Label: "Thay đổi Tên KH",
		},
		{
			Code:  "change-shipping-address",
			Label: "Thay đổi địa chỉ giao",
		},
		{
			Code:  "request-redelivering",
			Label: "Yêu cầu giao lại",
		},

		{
			Code:  "service-rating",
			Label: "Đánh giá dịch vụ",
		},
		{
			Code:  "request-contact",
			Label: "Liên hệ",
		},
		{
			Code:  "others",
			Label: "Yêu cầu khác",
		},
	}
)
var _ crmvtiger.Aggregate = &Aggregate{}

type Aggregate struct {
	vcsf        sqlstore.VtigerContactStoreFactory
	fieldMap    vtigermapping.ConfigMap
	vs          *vtigerstore.VtigerStore
	syncContact *sync.SyncVtiger
}

type FieldMap map[string]ItemFieldMap
type ItemFieldMap map[string]string

func New(db *cmsql.Database, fieldMap vtigermapping.ConfigMap, client *vtigerclient.VtigerClient) *Aggregate {
	return &Aggregate{
		fieldMap: fieldMap,
		vcsf:     sqlstore.NewVtigerStore(db),
		vs:       vtigerstore.NewVtigerStore(client, fieldMap),
	}
}

func (q *Aggregate) MessageBus() crmvtiger.CommandBus {
	b := bus.New()
	return crmvtiger.NewAggregateHandler(q).RegisterHandlers(b)
}

func (q *Aggregate) SyncContact(ctx context.Context, req *crmvtiger.SyncContactArgs) error {
	syncTime := req.SyncTime
	if q.syncContact == nil {
		q.syncContact = sync.NewSyncVtiger(q.vs, q.vcsf)
	}
	err := q.syncContact.SyncContact(syncTime)
	return err
}

func (q *Aggregate) CreateTicket(ctx context.Context, req *crmvtiger.CreateOrUpdateTicketArgs) (*crmvtiger.Ticket, error) {
	result, err := q.CreateOrUpdateTicket(ctx, req, "create")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (q *Aggregate) UpdateTicket(ctx context.Context, req *crmvtiger.CreateOrUpdateTicketArgs) (*crmvtiger.Ticket, error) {
	result, err := q.CreateOrUpdateTicket(ctx, req, "update")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (q *Aggregate) CreateOrUpdateLead(ctx context.Context, lead *crmvtiger.Lead) (*crmvtiger.Lead, error) {
	session, err := q.vs.Client.GetSessionKey()
	if err != nil {
		return nil, err
	}

	// send value to vtiger service
	fileMapData := q.vs.FieldMap
	vtigerMap, err := fileMapData.MappingLeadEtop2Vtiger(lead)
	if err != nil {
		return nil, err
	}
	leadResp, err := q.vs.CreateOrUpdateVtiger(vtigerMap, session, fileMapData, "Leads", Empty)
	if err != nil {
		return nil, err
	}
	leadReturn, err := fileMapData.MappingLeadVtiger2Etop(leadResp)

	// save to database
	contact := convert.ConvertLeadtoModelContact(leadReturn)
	query := q.vcsf(ctx).ByEtopUserID(contact.EtopUserID)
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
	return leadReturn, nil
}

func (q *Aggregate) CreateOrUpdateContact(ctx context.Context, ct *crmvtiger.Contact) (*crmvtiger.Contact, error) {
	session, err := q.vs.Client.GetSessionKey()
	if err != nil {
		return nil, err
	}
	// send value to vtiger service
	fileMapData := mapping.NewMappingConfigInfo(q.fieldMap)
	vtigerMap, err := fileMapData.MapingContactEtop2Vtiger(ct)
	if err != nil {
		return nil, err
	}
	contactResp, err := q.vs.CreateOrUpdateVtiger(vtigerMap, session, fileMapData, "Contacts", Empty)
	if err != nil {
		return nil, err
	}
	contactReturn, err := fileMapData.MapingContactVtiger2Etop(contactResp)
	if err != nil {
		return nil, err
	}
	// save value to db
	contact := convert.ConvertModelContact(contactReturn, session.UserID)
	if err = contact.BeforeInsertOrUpdate(); err != nil {
		return nil, err
	}
	_, err = q.vcsf(ctx).ByEtopUserID(contact.EtopUserID).GetContact()
	if err == nil {
		err = q.vcsf(ctx).UpdateVtigerContact(contact)
		if err != nil {
			return nil, err
		}
	} else {
		err = q.vcsf(ctx).CreateVtigerContact(contact)
		if err != nil {
			return nil, err
		}
	}
	return contactReturn, nil
}

// CreateOrUpdateTicket create or uodate ticket
func (q *Aggregate) CreateOrUpdateTicket(ctx context.Context, req *crmvtiger.CreateOrUpdateTicketArgs, action string) (*crmvtiger.Ticket, error) {

	// get session
	session, err := q.vs.Client.GetSessionKey()

	if err != nil {
		return nil, err
	}
	accout2Contact := convert.ConvertAccount(&req.Account)
	contactModel := convert.ConvertModelContact(accout2Contact, session.UserID)
	contactModel.EtopUserID = req.EtopUserID

	if contactModel.EtopUserID == 0 || contactModel.Phone == "" || contactModel.Email == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing argument Email, Phone or EtopUserID in body request")
	}
	query := q.vcsf(ctx).ByEtopUserID(contactModel.EtopUserID).ByEmail(contactModel.Email).ByPhone(contactModel.Phone)
	result, err := query.GetContact()

	// if not exist, create new Contact
	if err != nil {
		if err = contactModel.BeforeInsertOrUpdate(); err != nil {
			return nil, err
		}
		err = q.vcsf(ctx).CreateVtigerContact(contactModel)
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide exactly one argument")
		}
	} else {
		contactModel = result
	}

	ticket := convert.ConvertTicket(req)
	// Process Ticket
	ticketTitle := ticket.TicketTitle

	if ticket.Code != "others" {
		ticketTitle = ticketTitle + " đơn giao hàng " + ticket.FfmCode
	}

	if ticket.Code == "change-shop-cod" && ticket.NewValue != "" {
		ticket.NewValue = strings.ReplaceAll(ticket.NewValue, ".", "")
	}

	ticket.ContactId = contactModel.ID
	ticket.TicketTitle = ticketTitle
	ticket.Ticketpriorities = "Normal"
	ticket.Ticketstatus = "Open"

	categories := Categories
	ticket.Ticketcategories, err = MapTicketJSON(ticket.Code, categories)
	if err != nil {
		return nil, err
	}

	ticket.Description = ""
	if ticket.Code == "change-shop-cod" || ticket.Code == "change-phone" {
		ticket.Description = fmt.Sprintf("%v → %v", ticket.OldValue, ticket.NewValue)
	}

	fileMapData := q.vs.FieldMap
	vtigerMap, err := fileMapData.MappingTicketEtop2Vtiger(ticket)
	if err != nil {
		return nil, err
	}

	ticketResp, err := q.vs.CreateOrUpdateVtiger(vtigerMap, session, fileMapData, "HelpDesk", action)

	//convert
	ticketReturn, err := fileMapData.MappingTicketVtiger2Etop(ticketResp)

	if err != nil {
		return nil, err
	}
	return ticketReturn, nil
}

// MapTicketJSON Get label of reason follow code
func MapTicketJSON(code string, categories []*crmvtiger.Category) (string, error) {
	for _, value := range categories {
		if value.Code == code {
			return value.Label, nil
		}
	}
	return "", cm.Errorf(cm.InvalidArgument, nil, "Code categories not existed")
}
