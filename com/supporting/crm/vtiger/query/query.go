package query

import (
	"context"
	"net/url"
	"sort"
	"strconv"

	"etop.vn/api/meta"
	"etop.vn/api/supporting/crm/vtiger"
	"etop.vn/backend/com/supporting/crm/vtiger/convert"
	"etop.vn/backend/com/supporting/crm/vtiger/mapping"
	"etop.vn/backend/com/supporting/crm/vtiger/model"
	"etop.vn/backend/com/supporting/crm/vtiger/sqlstore"
	"etop.vn/backend/com/supporting/crm/vtiger/vtigerstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	sqlbuilder "etop.vn/backend/pkg/common/simple-sql-builder"
	"etop.vn/backend/pkg/integration/vtiger/client"
)

var (
	Categories = []*vtiger.Category{
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
var _ vtiger.QueryService = &QueryService{}

type QueryService struct {
	vcsf     sqlstore.VtigerContactStoreFactory
	fieldMap mapping.ConfigMap
	vs       *vtigerstore.VtigerStore
}

func New(db *cmsql.Database, fieldMap mapping.ConfigMap, client *client.VtigerClient) *QueryService {
	return &QueryService{
		fieldMap: fieldMap,
		vcsf:     sqlstore.NewVtigerStore(db),
		vs:       vtigerstore.NewVtigerStore(client, fieldMap),
	}
}

func (q *QueryService) MessageBus() vtiger.QueryBus {
	b := bus.New()
	return vtiger.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetCategories(context.Context, *meta.Empty) (*vtiger.GetCategoriesResponse, error) {
	return &vtiger.GetCategoriesResponse{
		Categories: Categories,
	}, nil
}

func (q *QueryService) GetTicketStatusCount(context.Context, *meta.Empty) (*vtiger.GetTicketStatusCountResponse, error) {
	var statusCounts []vtiger.CountTicketByStatusResponse
	for _, value := range Categories {
		status := value.Label
		categoriesStatus, err := q.vs.GetTickeyStatusByCategories(status)
		if err != nil {
			return nil, err
		}
		statusCounts = append(statusCounts, *categoriesStatus)
	}
	return &vtiger.GetTicketStatusCountResponse{
		StatusCount: statusCounts,
	}, nil
}

func (q *QueryService) GetTickets(ctx context.Context, req *vtiger.GetTicketsArgs) (*vtiger.GetTicketsResponse, error) {
	ticket := convert.ConvertTicketGetReq(&req.Ticket)

	fileMapData := q.vs.FieldMap
	vtigerMap, err := fileMapData.MappingTicketEtop2Vtiger(ticket)
	if err != nil {
		return nil, err
	}

	// create query
	query, err := q.BuildVtigerQuery("HelpDesk", vtigerMap, &req.Orderby, req.Paging)
	if err != nil {
		return nil, err
	}

	resultVtiger, err := q.VtigerRawQuery(query)
	if err != nil {
		return nil, err
	}

	response := make([]*vtiger.Ticket, 0)
	for _, value := range resultVtiger.Result {
		var mapTicket *vtiger.Ticket
		mapTicket, err = fileMapData.MappingTicketVtiger2Etop(value)
		if err != nil {
			return nil, err
		}
		response = append(response, mapTicket)
	}
	return &vtiger.GetTicketsResponse{
		Tickets: response,
	}, nil
}

func (q *QueryService) CountTicketByStatus(ctx context.Context, req *vtiger.CountTicketByStatusArgs) (*vtiger.CountTicketByStatusResponse, error) {
	status := req.Status

	var b sqlbuilder.SimpleSQLBuilder
	b.Printf("SELECT COUNT(*) FROM HelpDesk WHERE ticketstatus = ? ;", status)
	sqlVtiger, err := b.String()
	if err != nil {
		return nil, err
	}
	resultVtiger, err := q.VtigerRawQuery(sqlVtiger)
	if err != nil {
		return nil, err
	}
	countAtoi, err := strconv.Atoi(resultVtiger.Result[0]["count"])
	if err != nil {
		return nil, err
	}
	return &vtiger.CountTicketByStatusResponse{
		Code:  status,
		Count: int32(countAtoi),
	}, nil
}

func (q *QueryService) GetContacts(ctx context.Context, req *vtiger.GetContactsArgs) (*vtiger.ContactsResponse, error) {

	// search in db
	textSearch := req.Search
	var dbResult []*model.VtigerContact
	var err error
	if textSearch != "" {
		dbResult, err = q.vcsf(ctx).Paging(*req.Paging).SearchContact(textSearch)
	} else {
		dbResult, err = q.vcsf(ctx).Paging(*req.Paging).GetContacts()
	}
	if err != nil {
		return nil, err
	}
	var contactResult []*vtiger.Contact

	for i := 0; i < len(dbResult); i++ {
		dataRow := dbResult[i]
		contactResponseRow := convert.ConvertContactFromModel(dataRow)

		contactResult = append(contactResult, contactResponseRow)
	}
	var response vtiger.ContactsResponse
	response.Contacts = contactResult
	return &response, nil
}

// VtigerRawQuery request select
func (q *QueryService) VtigerRawQuery(query string) (*client.VtigerResponse, error) {
	session, err := q.vs.Client.GetSessionKey()
	if err != nil {
		return nil, err
	}

	queryURL := make(url.Values)
	queryURL.Set("operation", "query")
	queryURL.Set("sessionName", session.SessionName)
	queryURL.Set("query", query)

	return q.vs.Client.RequestGet(queryURL)
}

// BuildVtigerQuery build query type sql vtiger
func (q *QueryService) BuildVtigerQuery(module string, condition map[string]string, orderBy *vtiger.OrderBy, paging *meta.Paging) (string, error) {
	return buildVtigerQuery(q.fieldMap[module], module, condition, orderBy, paging)
}

func buildVtigerQuery(fieldMap mapping.ConfigGroup, module string, condition map[string]string, orderBy *vtiger.OrderBy, paging *meta.Paging) (string, error) {
	var b sqlbuilder.SimpleSQLBuilder
	b.Printf("SELECT * FROM ?", sqlbuilder.Raw(module))
	if len(condition) > 0 {
		b.Printf(" WHERE")
	}

	keys := make([]string, 0, len(condition))
	for key := range condition {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, key := range keys {
		value := condition[key]
		etopField := key
		if fieldMap[key] != "" {
			etopField = fieldMap[key]
		}
		if i != 0 {
			b.Printf(" AND")
		}
		b.Printf(" ? = ?", sqlbuilder.Raw(etopField), value)
	}

	if orderBy != nil {
		etopField := orderBy.Field
		if fieldMap[etopField] != "" {
			etopField = fieldMap[etopField]
		}
		if orderBy.Sort == "" {
			orderBy.Sort = " DESC"
		}
		b.Printf(" ORDER BY ? ?", sqlbuilder.Raw(etopField), sqlbuilder.Raw(orderBy.Sort))
	}
	if paging != nil {
		b.Printf(" LIMIT ?, ?", paging.Offset, paging.Limit)
	}
	b.Printf(";")
	returnValue, err := b.String()
	if err != nil {
		return "", err
	}
	return returnValue, nil
}

func (q *QueryService) GetRecordLastTimeModify(ctx context.Context, paging meta.Paging) (*vtiger.Contact, error) {
	result, err := q.vcsf(ctx).Paging(paging).SortBy("vtiger_updated_at desc").GetContacts()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return convert.ConvertContactFromModel(result[0]), nil
}
