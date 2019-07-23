package vtigerservice

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"etop.vn/backend/pb/services/crmservice"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/crm-service/mapping"
	"etop.vn/backend/pkg/crm-service/model"
	"etop.vn/backend/pkg/crm-service/sqlstore"
)

// VtigerService controller vtiger
type VtigerService struct {
	vtigerContact sqlstore.VtigerContactStoreFactory
	cfg           VtigerConfig
	fieldMap      *mapping.ConfigMap
}

// VtigerConfig information vtiger's config
type VtigerConfig struct {
	VtigerService   string `yaml:"vtiger_service"`
	VtigerUsername  string `yaml:"vtiger_username"`
	VtigerAccesskey string `yaml:"vtiger_accesskey"`
}

// NewSVtigerService init Service
func NewSVtigerService(db cmsql.Database, vConfig VtigerConfig, fieldMap *mapping.ConfigMap) *VtigerService {
	s := &VtigerService{
		cfg:           vConfig,
		fieldMap:      fieldMap,
		vtigerContact: sqlstore.NewVtigerStore(db),
	}
	return s
}

// ConvertAccout convert Account to Contact
func ConvertAccout(a *crmservice.Account) *crmservice.Contact {
	return &crmservice.Contact{
		EtopId:   a.Id,
		Lastname: a.FullName,
		Phone:    a.Phone,
		Email:    a.Email,
		Company:  a.Company,
	}
}

// ConvertTicket convert TicketRequest to Ticket protobuf. Ticket protobuf is used like DTO
func ConvertTicket(t *crmservice.TicketRequest) *crmservice.Ticket {
	ticket := &crmservice.Ticket{
		Code:        t.Code,
		TicketTitle: t.Title,
		NewValue:    t.Value,
		OldValue:    t.OldValue,
		Reason:      t.Reason,
		EtopId:      t.EtopId,
		OrderId:     t.OrderId,
		FfmCode:     t.FfmCode,
		FfmUrl:      t.FfmUrl,
		FfmId:       t.FfmId,
		Company:     t.Company,
		Provider:    t.Provider,
		Note:        t.Note,
		Environment: t.Environment,
		FromApp:     t.FromApp,
		Id:          t.Id,
	}
	return ticket
}

// MapTicketJSON Get label of reason follow code
func MapTicketJSON(code string, categories []*crmservice.Categories) (string, error) {
	for _, value := range categories {
		if value.Code == code {
			return value.Label, nil
		}
	}
	return "", cm.Errorf(cm.InvalidArgument, nil, "Code categories not existed")
}

func singleQuote(value string) (string, error) {
	if strings.Contains(value, `"`) {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Value string contain Double-Qoute")
	}
	return strings.ReplaceAll(strconv.Quote(value), `"`, `'`), nil
}

// ReadFileCategories read file reason
func ReadFileCategories() ([]*crmservice.Categories, error) {
	configFile, err := os.Open("../pkg/crm-service/reason_mapping.json")
	if err != nil {
		configFile, err = os.Open("pkg/crm-service/reason_mapping.json")
		if err != nil {
			return nil, err
		}
	}

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var categoriesResponse []*crmservice.Categories
	err = json.Unmarshal(byteValue, &categoriesResponse)
	if err != nil {
		return nil, err
	}
	return categoriesResponse, nil
}

// ConvertModelContact covert protobuf to model Contact
func ConvertModelContact(c *crmservice.Contact, AssignedUserID string) *model.VtigerContact {
	contact := &model.VtigerContact{
		ID:                   c.Id,
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
		EtopID:               c.EtopId,
		Source:               c.Source,
		UsedShippingProvider: c.UsedShippingProvider,
		OrdersPerDay:         c.OrdersPerDay,
		Company:              c.Company,
		City:                 c.City,
		State:                c.State,
		Website:              c.Website,
		Lane:                 c.Lane,
		Country:              c.Country,
	}
	return contact
}

// ReadFileConfig read json file which is use for map vtiger and etop
func ReadFileConfig() (string, error) {

	configFile, err := os.Open("../pkg/crm-service/field_mapping.json")
	if err != nil {
		configFile, err = os.Open("pkg/crm-service/field_mapping.json")
		if err != nil {
			return "", err
		}
	}

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return "", err
	}
	defer configFile.Close()

	return string(byteValue), nil
}
