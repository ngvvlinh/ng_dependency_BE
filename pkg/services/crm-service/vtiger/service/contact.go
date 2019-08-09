package vtigerservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/schema"

	"etop.vn/api/meta"
	"etop.vn/backend/pb/services/crmservice"
	"etop.vn/backend/pkg/services/crm-service/mapping"
	"etop.vn/backend/pkg/services/crm-service/model"
	"etop.vn/backend/pkg/services/crm-service/vtiger/client"
)

var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("json")
}

// CreateOrUpdateContact .
func (s *VtigerService) CreateOrUpdateContact(ctx context.Context, ct *crmservice.Contact) (*crmservice.Contact, error) {

	session, err := s.client.GetSessionKey(s.cfg.ServiceURL, s.cfg.Username, s.cfg.APIKey)
	if err != nil {
		return nil, err
	}

	// save value to db
	contact := ConvertModelContact(ct, session.UserID)
	if err := contact.BeforeInsertOrUpdate(); err != nil {
		return nil, err
	}
	_, err = s.vtigerContact(ctx).ByEtopID(contact.EtopID).GetContact()
	if err == nil {
		err = s.vtigerContact(ctx).UpdateVtigerContact(contact)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.vtigerContact(ctx).CreateVtigerContact(contact)
		if err != nil {
			return nil, err
		}
	}

	// send value to vtiger service
	fileMapData := mapping.NewMappingConfigInfo(s.fieldMap)
	vtigerMap, err := fileMapData.MapingContactEtop2Vtiger(ct)
	if err != nil {
		return nil, err
	}
	contactResp, err := s.CreateOrUpdateVtiger(vtigerMap, session, fileMapData, "Contacts")
	//convert
	contactReturn, err := fileMapData.MapingContactVtiger2Etop(contactResp)
	if err != nil {
		return nil, err
	}
	return contactReturn, nil
}

// BodyRequestVtiger define struct body request to vtiger
type BodyRequestVtiger struct {
	SessionName string `json:"sessionName"`
	Operation   string `json:"operation"`
	Element     string `json:"element"`
	ElementType string `json:"elementType"`
}

// BodyResponseVtiger define struct body response from vtiger
type BodyResponseVtiger struct {
	Success string
	Result  map[string]string
}

// CreateOrUpdateVtiger request create or to vtiger
func (s *VtigerService) CreateOrUpdateVtiger(
	etop2Vtiger map[string]string,
	session *client.VtigerSessionResult,
	fileMapData *mapping.Mapper,
	moduleName string,
) (map[string]string, error) {
	etop2Vtiger["assigned_user_id"] = session.UserID
	bodyRequestVtiger := &BodyRequestVtiger{
		ElementType: moduleName,
		SessionName: session.SessionName,
	}

	//check exit id is already exit
	bodyRequestVtiger.Operation = "create"

	if etop2Vtiger["id"] != "" {
		sq, err := singleQuote(etop2Vtiger["id"])
		if err != nil {
			return nil, err
		}
		sqlQuery := fmt.Sprintf(`SELECT * from %v where id=%v;`, moduleName, sq)

		queryValues := make(url.Values)
		queryValues.Set("operation", "query")
		queryValues.Set("sessionName", session.SessionName)
		queryValues.Set("query", sqlQuery)

		vtigerService := client.NewVigerClient(session.SessionName, s.cfg.ServiceURL)
		resultVtiger, err := vtigerService.RequestGet(queryValues)
		if err != nil {
			return nil, err
		}
		// create body request
		if len(resultVtiger.Result) != 0 {
			bodyRequestVtiger.Operation = "update"
		} else {
			delete(etop2Vtiger, "id")
		}
	}

	etop2VtigerString, err := json.Marshal(etop2Vtiger)
	if err != nil {
		return nil, err
	}
	bodyRequestVtiger.Element = string(etop2VtigerString)

	// request update or create
	requestBody := url.Values{}
	if err := encoder.Encode(bodyRequestVtiger, requestBody); err != nil {
		return nil, err
	}

	var result BodyResponseVtiger
	err = s.client.SendPost(requestBody, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}

// GetContacts get contact from db
func (s *VtigerService) GetContacts(ctx context.Context, getContactRequest *crmservice.GetContactsRequest) (*crmservice.GetContactsResponse, error) {
	page := getContactRequest.Page
	perpage := getContactRequest.Perpage
	if perpage == 0 {
		perpage = 100
	}
	if perpage > 1000 {
		perpage = 1000
	}
	var paging meta.Paging
	paging.Offset = page
	paging.Limit = perpage

	//search in db
	textSearch := getContactRequest.TextSearch
	var dbResult []*model.VtigerContact
	var err error
	if textSearch != "" {
		dbResult, err = s.vtigerContact(ctx).Paging(paging).SearchContact(textSearch)
	} else {
		dbResult, err = s.vtigerContact(ctx).Paging(paging).GetContacts()
	}
	if err != nil {
		return nil, err
	}
	var contactResult []*crmservice.Contact

	for i := 0; i < len(dbResult); i++ {
		dataRow := dbResult[i]
		contactResponseRow := &crmservice.Contact{
			Id:                   dataRow.ID,
			Firstname:            dataRow.Firstname,
			ContactNo:            dataRow.ContactNo,
			Phone:                dataRow.Phone,
			Lastname:             dataRow.Lastname,
			Mobile:               dataRow.Mobile,
			Email:                dataRow.Email,
			Leadsource:           dataRow.Leadsource,
			Secondaryemail:       dataRow.Secondaryemail,
			AssignedUserId:       dataRow.AssignedUserID,
			EtopId:               dataRow.EtopID,
			Source:               dataRow.Source,
			UsedShippingProvider: dataRow.UsedShippingProvider,
			OrdersPerDay:         dataRow.OrdersPerDay,
			Company:              dataRow.Company,
			City:                 dataRow.City,
			State:                dataRow.State,
			Website:              dataRow.Website,
			Lane:                 dataRow.Lane,
			Country:              dataRow.Country,
		}

		contactResult = append(contactResult, contactResponseRow)
	}
	return &crmservice.GetContactsResponse{
		Contacts: contactResult,
	}, nil
}
