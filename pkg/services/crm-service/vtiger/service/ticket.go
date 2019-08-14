package vtigerservice

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"

	"etop.vn/backend/pkg/services/crm-service/vtiger/client"

	"etop.vn/backend/pb/services/crmservice"
	cm "etop.vn/backend/pkg/common"
	simpleSqlBuilder "etop.vn/backend/pkg/common/simple-sql-builder"
	"etop.vn/backend/pkg/services/crm-service/mapping"
)

// GetCategories get categories
func (s *VtigerService) GetCategories(ctx context.Context) (*crmservice.GetCategoriesResponse, error) {
	categories := GetCategories()
	return &crmservice.GetCategoriesResponse{
		Categories: categories,
	}, nil
}

// CreateOrUpdateTicket create or uodate ticket
func (s *VtigerService) CreateOrUpdateTicket(ctx context.Context, req *crmservice.CreateOrUpdateTicketRequest, action string) (*crmservice.Ticket, error) {

	// get session
	session, err := s.Client.GetSessionKey(s.Cfg.ServiceURL, s.Cfg.Username, s.Cfg.APIKey)
	if err != nil {
		return nil, err
	}
	accout2Contact := ConvertAccount(&req.Account)
	contactModel := ConvertModelContact(accout2Contact, session.UserID)
	contactModel.EtopID = req.EtopId
	pp.Print("contactModel ::", contactModel)
	if contactModel.EtopID == 0 || contactModel.Phone == "" || contactModel.Email == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing argument Email, Phone or EtopID in body request")
	}
	query := s.vtigerContact(ctx).ByEtopID(contactModel.EtopID).ByEmail(contactModel.Email).ByPhone(contactModel.Phone)
	result, err := query.GetContact()

	// if not exist, create new Contact
	if err != nil {
		if err = contactModel.BeforeInsertOrUpdate(); err != nil {
			return nil, err
		}
		err = s.vtigerContact(ctx).CreateVtigerContact(contactModel)
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide exactly one argument")
		}
	} else {
		contactModel = result
	}

	ticket := ConvertTicket(req)
	// Process Ticket
	ticketTitle := ticket.TicketTitle

	if ticket.Code != "others" {
		ticketTitle = fmt.Sprintf("%v đơn giao hàng %v", ticketTitle, ticket.FfmCode)
	}

	if ticket.Code == "change-shop-cod" && ticket.NewValue != "" {
		ticket.NewValue = strings.ReplaceAll(ticket.NewValue, ".", "")
	}

	ticket.ContactId = contactModel.ID
	ticket.TicketTitle = ticketTitle
	ticket.Ticketpriorities = "Normal"
	ticket.Ticketstatus = "Open"

	categories := GetCategories()
	ticket.Ticketcategories, err = MapTicketJSON(ticket.Code, categories)
	if err != nil {
		return nil, err
	}

	ticket.Description = ""
	if ticket.Code == "change-shop-cod" || ticket.Code == "change-phone" {
		ticket.Description = fmt.Sprintf("%v → %v", ticket.OldValue, ticket.NewValue)
	}

	fileMapData := mapping.NewMappingConfigInfo(s.fieldMap)
	vtigerMap, err := fileMapData.MappingTicketEtop2Vtiger(ticket)
	if err != nil {
		return nil, err
	}

	ticketResp, err := s.CreateOrUpdateVtiger(vtigerMap, session, fileMapData, "HelpDesk", action)
	//convert
	ticketReturn, err := fileMapData.MappingTicketVtiger2Etop(ticketResp)
	if err != nil {
		return nil, err
	}
	return ticketReturn, nil
}

// GetTickets get ticket from vtiger
func (s *VtigerService) GetTickets(ctx context.Context, getTicketsRequest *crmservice.GetTicketsRequest) (*crmservice.GetTicketsResponse, error) {
	page := getTicketsRequest.Page
	perPager := getTicketsRequest.Perpage

	ticket := ConvertTicketGetReq(&getTicketsRequest.Ticket)

	fileMapData := mapping.NewMappingConfigInfo(s.fieldMap)
	vtigerMap, err := fileMapData.MappingTicketEtop2Vtiger(ticket)
	if err != nil {

	}

	// create query
	query, err := s.BuildVtigerQuery("HelpDesk", vtigerMap, &getTicketsRequest.Orderby, page, perPager)
	if err != nil {
		return nil, err
	}

	resultVtiger, err := s.VtigerRawQuery(query)
	if err != nil {
		return nil, err
	}

	response := make([]*crmservice.Ticket, 0)
	for _, value := range resultVtiger.Result {
		var mapTicket *crmservice.Ticket
		mapTicket, err = fileMapData.MappingTicketVtiger2Etop(value)
		if err != nil {
			return nil, err
		}
		response = append(response, mapTicket)
	}
	return &crmservice.GetTicketsResponse{
		Tickets: response,
	}, nil
}

// BuildVtigerQuery build query type sql vtiger
func (s *VtigerService) BuildVtigerQuery(module string, condition map[string]string, orderBy *crmservice.OrderBy, page int32, perpage int32) (string, error) {
	var b simpleSqlBuilder.SimpleSQLBuilder
	b.Printf("SELECT * FROM ? ", simpleSqlBuilder.Raw(module))
	if page == 0 {
		page = 1
	}
	if perpage == 0 {
		perpage = 20
	}
	fieldMap := s.fieldMap[module]

	if len(condition) > 0 {
		b.Printf(" WHERE ")
	}
	i := 1
	for key, value := range condition {
		etopField := key
		if fieldMap[key] != "" {
			etopField = fieldMap[etopField]
		}
		if i == 1 {
			b.Printf(" ? = ? ", simpleSqlBuilder.Raw(etopField), value)
			continue
		}
		b.Printf(" AND ? = ? ", simpleSqlBuilder.Raw(etopField), value)
		i++
	}

	etopField := orderBy.Field
	if fieldMap[etopField] != "" {
		etopField = fieldMap[etopField]
	}
	if orderBy.Sort == "" {
		orderBy.Sort = "DESC"
	}
	if orderBy != nil {
		b.Printf("ORDER BY ? ? ", simpleSqlBuilder.Raw(etopField), simpleSqlBuilder.Raw(orderBy.Sort))
	}
	b.Printf("LIMIT ?, ? ;", (page-1)*perpage, perpage)
	returnValue, err := b.String()
	if err != nil {
		return "", err
	}
	return returnValue, nil
}

// CountTicketByStatus get count number of ticket follow status
func (s *VtigerService) CountTicketByStatus(ctx context.Context, countTicketByStatusRequest *crmservice.CountTicketByStatusRequest) (*crmservice.CountTicketByStatusResponse, error) {
	status := *countTicketByStatusRequest.Status
	//make SQL query vtiger
	var b simpleSqlBuilder.SimpleSQLBuilder
	b.Printf("SELECT COUNT(*) FROM HelpDesk WHERE ticketstatus = ? ;", status)
	sqlVtiger, err := b.String()
	if err != nil {
		return nil, err
	}
	resultVtiger, err := s.VtigerRawQuery(sqlVtiger)
	if err != nil {
		return nil, err
	}
	countAtoi, err := strconv.Atoi(resultVtiger.Result[0]["count"])
	if err != nil {
		return nil, err
	}
	return &crmservice.CountTicketByStatusResponse{
		Code:  status,
		Count: int32(countAtoi),
	}, nil
}

// VtigerRawQuery request select
func (s *VtigerService) VtigerRawQuery(query string) (*client.VtigerResponse, error) {
	session, err := s.Client.GetSessionKey(s.Cfg.ServiceURL, s.Cfg.Username, s.Cfg.APIKey)
	if err != nil {
		return nil, err
	}

	queryURL := make(url.Values)
	queryURL.Set("operation", "query")
	queryURL.Set("sessionName", session.SessionName)
	queryURL.Set("query", query)

	return s.Client.RequestGet(queryURL)
}

// GetTicketStatusCount get ticket status count
func (s *VtigerService) GetTicketStatusCount(ctx context.Context) (*crmservice.GetTicketStatusCountResponse, error) {
	session, err := s.Client.GetSessionKey(s.Cfg.ServiceURL, s.Cfg.Username, s.Cfg.APIKey)
	if err != nil {
		return nil, err
	}
	categories := GetCategories()

	var statusCounts []crmservice.CountTicketByStatusResponse
	for _, value := range categories {
		status := value.Label

		//make SQL query vtiger
		var b simpleSqlBuilder.SimpleSQLBuilder
		b.Printf("SELECT COUNT(*) FROM HelpDesk WHERE ticketcategories = ? AND ticketstatus = 'Open';", status)
		var sql string
		sql, err = b.String()
		if err != nil {
			return nil, err
		}

		queryValues := make(url.Values)
		queryValues.Set("operation", "query")
		queryValues.Set("sessionName", session.SessionName)
		queryValues.Set("query", sql)

		var response *client.VtigerResponse
		response, err = s.Client.RequestGet(queryValues)
		if err != nil {
			return nil, err
		}
		var countAtoi int
		countAtoi, err = strconv.Atoi(response.Result[0]["count"])
		if err != nil {
			return nil, err
		}
		statusCounts = append(statusCounts, crmservice.CountTicketByStatusResponse{
			Code:  status,
			Count: int32(countAtoi),
		})
	}
	return &crmservice.GetTicketStatusCountResponse{
		StatusCount: statusCounts,
	}, nil
}
