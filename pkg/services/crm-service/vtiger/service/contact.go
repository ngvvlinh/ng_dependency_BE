package vtigerservice

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/gorilla/schema"

	"etop.vn/api/meta"
	"etop.vn/backend/pb/services/crmservice"
	simpleSqlBuilder "etop.vn/backend/pkg/common/simple-sql-builder"
	"etop.vn/backend/pkg/services/crm-service/mapping"
	"etop.vn/backend/pkg/services/crm-service/model"
	"etop.vn/backend/pkg/services/crm-service/vtiger/client"

	cm "etop.vn/backend/pkg/common"
)

var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("json")
}

// CreateOrUpdateContact .
func (s *VtigerService) CreateOrUpdateContact(ctx context.Context, ct *crmservice.Contact) (*crmservice.Contact, error) {
	session, err := s.Client.GetSessionKey(s.Cfg.ServiceURL, s.Cfg.Username, s.Cfg.APIKey)
	if err != nil {
		return nil, err
	}

	// send value to vtiger service
	fileMapData := mapping.NewMappingConfigInfo(s.fieldMap)
	vtigerMap, err := fileMapData.MapingContactEtop2Vtiger(ct)
	if err != nil {
		return nil, err
	}
	contactResp, err := s.CreateOrUpdateVtiger(vtigerMap, session, fileMapData, "Contacts", Empty)

	contactReturn, err := fileMapData.MapingContactVtiger2Etop(contactResp)
	if err != nil {
		return nil, err
	}

	// save value to db
	contact := ConvertModelContact(contactReturn, session.UserID)
	if err = contact.BeforeInsertOrUpdate(); err != nil {
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
	Success bool
	Result  map[string]string
}

// CreateOrUpdateVtiger request create or to vtiger
func (s *VtigerService) CreateOrUpdateVtiger(
	etop2Vtiger map[string]string,
	session *client.VtigerSessionResult,
	fileMapData *mapping.Mapper,
	moduleName string,
	action string,
) (map[string]string, error) {
	etop2Vtiger["assigned_user_id"] = session.UserID
	bodyRequestVtiger := &BodyRequestVtiger{
		ElementType: moduleName,
		SessionName: session.SessionName,
	}

	//check exit id is already exit
	bodyRequestVtiger.Operation = "create"
	if moduleName == "HelpDesk" {
		bodyRequestVtiger.Operation = action
		if action == "update " && etop2Vtiger["id"] == "" {
			return nil, cm.Error(cm.InvalidArgument, "Missing id ", nil)
		}
		if action == "update " {
			var b simpleSqlBuilder.SimpleSQLBuilder

			b.Printf(`SELECT * from ? where ?=?;`, simpleSqlBuilder.Raw(moduleName), "id", etop2Vtiger["id"])
			sqlQuery, err := b.String()
			if err != nil {
				return nil, err
			}

			queryValues := make(url.Values)
			queryValues.Set("operation", "query")
			queryValues.Set("sessionName", session.SessionName)
			queryValues.Set("query", sqlQuery)

			resultVtiger, err := s.Client.RequestGet(queryValues)
			if err != nil {
				return nil, err
			}
			// create body request
			if len(resultVtiger.Result) == 0 {
				return nil, cm.Error(cm.InvalidArgument, "Missing id does not exist in vtiger", nil)
			}
		}
	} else
	//check Already exist in vtiger
	{
		fileMap := fileMapData.FieldMap[moduleName]

		if etop2Vtiger[fileMap["etop_id"]] == "" {
			return nil, cm.Error(cm.InvalidArgument, "Missing etop_id in request", nil)
		}

		var b simpleSqlBuilder.SimpleSQLBuilder

		b.Printf(`SELECT * from ? where ?=?;`, simpleSqlBuilder.Raw(moduleName), simpleSqlBuilder.Raw(fileMap["etop_id"]), etop2Vtiger[fileMap["etop_id"]])

		sqlQuery, err := b.String()
		if err != nil {
			return nil, err
		}

		queryValues := make(url.Values)
		queryValues.Set("operation", "query")
		queryValues.Set("sessionName", session.SessionName)
		queryValues.Set("query", sqlQuery)

		resultVtiger, err := s.Client.RequestGet(queryValues)
		if err != nil {
			return nil, err
		}
		// create body request
		if len(resultVtiger.Result) != 0 {
			bodyRequestVtiger.Operation = "update"
			resultVtigerArray := resultVtiger.Result
			etop2Vtiger["id"] = resultVtigerArray[0]["id"]
		}
	}

	etop2VtigerString, err := json.Marshal(etop2Vtiger)
	if err != nil {
		return nil, err
	}
	bodyRequestVtiger.Element = string(etop2VtigerString)

	// request update or create
	requestBody := url.Values{}
	if err = encoder.Encode(bodyRequestVtiger, requestBody); err != nil {
		return nil, err
	}

	var result BodyResponseVtiger
	err = s.Client.SendPost(requestBody, &result)

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

// sync data vtiger
func (s *VtigerService) SyncContac() error {

	fileMapData := mapping.NewMappingConfigInfo(s.fieldMap)
	ctx := context.Background()
	page := 0
	perPage := 50
	for true {
		var b simpleSqlBuilder.SimpleSQLBuilder

		b.Printf(`SELECT * FROM Contacts LIMIT ?, ? ;`, page*perPage, perPage)
		sqlQuery, err := b.String()
		if err != nil {
			return err
		}
		queryValues := make(url.Values)
		queryValues.Set("operation", "query")
		queryValues.Set("sessionName", s.Client.SessionInfo.VtigerSession.SessionName)
		queryValues.Set("query", sqlQuery)

		result, err := s.Client.RequestGet(queryValues)

		if err != nil {
			return err
		}
		if len(result.Result) == 0 {
			break
		}
		for _, value := range result.Result {
			var contact *crmservice.Contact

			contact, err = fileMapData.MapingContactVtiger2Etop(value)
			if err != nil {
				return err
			}
			modelContact := ConvertModelContact(contact, s.Client.SessionInfo.VtigerSession.UserID)
			err = s.CreateOrUpdateContactToDB(ctx, modelContact)
			if err != nil {
				return err
			}
		}
		page = page + 1
	}

	return nil
}

func (s *VtigerService) CreateOrUpdateContactToDB(ctx context.Context, contact *model.VtigerContact) error {
	_, err := s.vtigerContact(ctx).ByEtopID(contact.EtopID).GetContact()
	if err != nil && cm.ErrorCode(err) == cm.NotFound {
		return s.vtigerContact(ctx).ByEtopID(contact.EtopID).CreateVtigerContact(contact)
	} else if err != nil {
		return err
	} else {
		return s.vtigerContact(ctx).ByEtopID(contact.EtopID).UpdateVtigerContact(contact)
	}
}
