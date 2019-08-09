package vtigerservice

import (
	"strconv"
	"strings"

	"etop.vn/backend/pkg/services/crm-service/vtiger/client"

	"etop.vn/backend/pb/services/crmservice"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/services/crm-service/mapping"
	"etop.vn/backend/pkg/services/crm-service/model"
	"etop.vn/backend/pkg/services/crm-service/sqlstore"
)

// Config represents configuration for vtiger service
type Config struct {
	ServiceURL string `yaml:"service_url"`
	Username   string `yaml:"username"`
	APIKey     string `yaml:"api_key"`
}

func (c *Config) MustLoadEnv(prefix string) {
	p := prefix
	cc.EnvMap{
		p + "_SERVICE_URL": &c.ServiceURL,
		p + "_USERNAME":    &c.Username,
		p + "_API_KEY":     &c.APIKey,
	}.MustLoad()
}

var Categories = []*crmservice.Category{
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

// VtigerService controller vtiger
type VtigerService struct {
	vtigerContact sqlstore.VtigerContactStoreFactory
	cfg           Config
	fieldMap      mapping.ConfigMap
	client        *client.VtigerClient
}

// NewSVtigerService init Service
func NewSVtigerService(db cmsql.Database, vConfig Config, fieldMap mapping.ConfigMap) *VtigerService {
	s := &VtigerService{
		cfg:           vConfig,
		fieldMap:      fieldMap,
		vtigerContact: sqlstore.NewVtigerStore(db),
	}
	return s
}

// ConvertAccount convert Account to Contact
func ConvertAccount(a *crmservice.Account) *crmservice.Contact {
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
func MapTicketJSON(code string, categories []*crmservice.Category) (string, error) {
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

// GetCategories
func GetCategories() []*crmservice.Category {
	return Categories
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
