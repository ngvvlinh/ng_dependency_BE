package vtigerstore

import (
	"net/url"
	"strconv"

	"github.com/gorilla/schema"

	"etop.vn/api/supporting/crm/vtiger"
	"etop.vn/backend/com/supporting/crm/vtiger/mapping"
	cm "etop.vn/backend/pkg/common"
	sqlbuilder "etop.vn/backend/pkg/common/simple-sql-builder"
	vtigerclient "etop.vn/backend/pkg/integration/vtiger/client"
	"etop.vn/common/jsonx"
)

var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("url")
}

type VtigerStore struct {
	Client   *vtigerclient.VtigerClient
	FieldMap *mapping.Mapper
}

func NewVtigerStore(client *vtigerclient.VtigerClient, mapField mapping.ConfigMap) *VtigerStore {
	return &VtigerStore{
		Client:   client,
		FieldMap: mapping.NewMappingConfigInfo(mapField),
	}
}

// BodyRequestVtiger define struct body request to vtiger
type BodyRequestVtiger struct {
	SessionName string `url:"sessionName"`
	Operation   string `url:"operation"`
	Element     string `url:"element"`
	ElementType string `url:"elementType"`
}

// BodyResponseVtiger define struct body response from vtiger
type BodyResponseVtiger struct {
	Success bool              `json:"suscess"`
	Result  map[string]string `json:"result"`
}

func (v *VtigerStore) CheckExist(module string, field string, value string, sessionName string) (string, error) {
	var b sqlbuilder.SimpleSQLBuilder

	b.Printf(`SELECT * from ? where ?=?;`, sqlbuilder.Raw(module), sqlbuilder.Raw(field), value)

	sqlQuery, err := b.String()
	if err != nil {
		return "", err
	}

	queryValues := make(url.Values)
	queryValues.Set("operation", "query")
	queryValues.Set("sessionName", sessionName)
	queryValues.Set("query", sqlQuery)

	resultVtiger, err := v.Client.RequestGet(queryValues)
	if err != nil {
		return "", err
	}
	if len(resultVtiger.Result) != 0 {
		resultVtigerArray := resultVtiger.Result
		return resultVtigerArray[0]["id"], nil
	}
	return "", nil
}

// CreateOrUpdateVtiger request create or to vtiger
func (v *VtigerStore) CreateOrUpdateVtiger(
	etop2Vtiger map[string]string,
	session *vtigerclient.VtigerSessionResult,
	fileMapData *mapping.Mapper,
	moduleName string,
	action string,
) (map[string]string, error) {
	etop2Vtiger["assigned_user_id"] = session.UserID
	bodyRequestVtiger := &BodyRequestVtiger{
		ElementType: moduleName,
		SessionName: session.SessionName,
	}
	vs := VtigerStore{Client: v.Client}
	bodyRequestVtiger.Operation = "create"
	if moduleName == "HelpDesk" {
		bodyRequestVtiger.Operation = action
		if action == "update" && etop2Vtiger["id"] == "" {
			return nil, cm.Error(cm.InvalidArgument, "Missing id ", nil)
		}
		if action == "update" {
			id, err := vs.CheckExist(moduleName, "id", etop2Vtiger["id"], session.SessionName)
			if err != nil {
				return nil, err
			}
			if id == "" {
				return nil, cm.Error(cm.InvalidArgument, "Missing id does not exist in vtiger", nil)
			}
		}
	} else {
		fileMap := fileMapData.FieldMap[moduleName]

		if etop2Vtiger[fileMap["etop_user_id"]] == "" {
			return nil, cm.Error(cm.InvalidArgument, "Missing etop_user_id in request", nil)
		}
		id, err := vs.CheckExist(moduleName, fileMap["etop_user_id"], etop2Vtiger[fileMap["etop_user_id"]], session.SessionName)
		if err != nil {
			return nil, err
		}
		if id != "" {
			bodyRequestVtiger.Operation = "update"
			etop2Vtiger["id"] = id
		}
	}

	etop2VtigerString, err := jsonx.Marshal(etop2Vtiger)
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
	err = v.Client.SendPost(requestBody, &result)

	if err != nil {
		return nil, err
	}

	return result.Result, nil
}

func (v *VtigerStore) GetTickeyStatusByCategories(status string) (*vtiger.CountTicketByStatusResponse, error) {
	_, err := v.Client.GetSessionKey()
	if err != nil {
		return nil, err
	}
	//make SQL query vtiger
	var b sqlbuilder.SimpleSQLBuilder
	b.Printf("SELECT COUNT(*) FROM HelpDesk WHERE ticketcategories = ? AND ticketstatus = 'Open';", status)
	var sql string
	sql, err = b.String()
	if err != nil {
		return nil, err
	}

	queryValues := make(url.Values)
	queryValues.Set("operation", "query")
	queryValues.Set("sessionName", v.Client.SessionInfo.VtigerSession.SessionName)
	queryValues.Set("query", sql)

	var response *vtigerclient.VtigerResponse
	response, err = v.Client.RequestGet(queryValues)
	if err != nil {
		return nil, err
	}
	var countAtoi int
	countAtoi, err = strconv.Atoi(response.Result[0]["count"])
	if err != nil {
		return nil, err
	}
	return &vtiger.CountTicketByStatusResponse{
		Code:  status,
		Count: int32(countAtoi),
	}, nil
}
