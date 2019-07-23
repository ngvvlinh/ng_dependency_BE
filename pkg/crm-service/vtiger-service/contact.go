package vtigerservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"etop.vn/api/meta"
	"etop.vn/backend/pb/services/crmservice"
	"etop.vn/backend/pkg/crm-service/mapping"
	"etop.vn/backend/pkg/crm-service/model"
	"etop.vn/backend/pkg/crm-service/vtiger"
	"github.com/gorilla/schema"
	"github.com/k0kubun/pp"
)

var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("json")
}

// CreateOrUpdateContact .
func (v *VtigerService) CreateOrUpdateContact(ctx context.Context, ct *crmservice.Contact) (*crmservice.Contact, error) {

	// get session
	session, err := vtiger.GetSessionKey(v.cfg.VtigerService, v.cfg.VtigerUsername, v.cfg.VtigerAccesskey)
	if err != nil {
		return nil, err
	}

	// save value to db
	contact := ConvertModelContact(ct, session.UserID)

	if err := contact.BeforeInsertOrUpdate(); err != nil {
		return nil, err
	}
	_, err = v.vtigerContact(ctx).ByEtopID(contact.EtopID).GetContact()
	if err == nil {
		err = v.vtigerContact(ctx).UpdateVtigerContact(contact)
		if err != nil {
			return nil, err
		}
	} else {
		err = v.vtigerContact(ctx).CreateVtigerContact(contact)
		if err != nil {
			return nil, err
		}
	}
	// send value to vtiger service
	fileMapData := mapping.NewMappingConfigInfo(v.fieldMap)
	vtigerMap, err := fileMapData.MapingContactEtop2Vtiger(ct)
	if err != nil {
		return nil, err
	}
	contactResp, err := v.CreateOrUpdateVtiger(vtigerMap, session, fileMapData, "Contacts")
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

// BodyResponseViter define struct body response from vtiger
type BodyResponseViter struct {
	Success string
	Result  map[string]string
}

// CreateOrUpdateVtiger request create or to vtiger
func (v *VtigerService) CreateOrUpdateVtiger(etop2Vtiger map[string]string,
	session *vtiger.VtigerSessionResult,
	fileMapData *mapping.MappingConfigInfo,
	moduleName string) (map[string]string, error) {
	etop2Vtiger["assigned_user_id"] = session.UserID
	bodyRequestVtiger := &BodyRequestVtiger{
		ElementType: moduleName,
		SessionName: session.SessionName,
	}

	//check exit id is already exit
	bodyRequestVtiger.Operation = "create"

	pp.Println("ID :: ", etop2Vtiger["id"])

	if etop2Vtiger["id"] != "" {
		sq, err := singleQuote(etop2Vtiger["id"])
		if err != nil {
			return nil, err
		}
		sqlQuery := fmt.Sprintf(`SELECT * from %v where id=%v;`, moduleName, sq)

		u, err := url.Parse("")
		if err != nil {
			return nil, err
		}
		queryURL := u.Query()
		queryURL.Set("operation", "query")
		queryURL.Set("sessionName", session.SessionName)
		queryURL.Set("query", sqlQuery)
		u.RawQuery = queryURL.Encode()
		path := "?" + u.RawQuery

		vtigerService := vtiger.NewVigerClient(session.SessionName, v.cfg.VtigerService)
		resultVtiger, err := vtigerService.SendRequestVtigerValue(path, nil, "GET")
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
	data := url.Values{}
	if err := encoder.Encode(bodyRequestVtiger, data); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", v.cfg.VtigerService, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// get map result after update or create
	var result BodyResponseViter
	err = json.Unmarshal([]byte(bodyResp), &result)

	return result.Result, nil
}

// GetContacts get contact from db
func (v *VtigerService) GetContacts(ctx context.Context, getContactRequest *crmservice.GetContactsRequest) (*crmservice.GetContactsResponse, error) {
	// paging
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
		dbResult, err = v.vtigerContact(ctx).Paging(paging).SearchContact(textSearch)
	} else {
		dbResult, err = v.vtigerContact(ctx).Paging(paging).GetContacts()
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
