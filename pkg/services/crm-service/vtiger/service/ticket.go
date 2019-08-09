package vtigerservice

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"etop.vn/backend/pkg/services/crm-service/vtiger/client"

	"etop.vn/backend/pb/services/crmservice"
	cm "etop.vn/backend/pkg/common"
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
func (s *VtigerService) CreateOrUpdateTicket(ctx context.Context, req *crmservice.CreateOrUpdateTicketRequest) (*crmservice.Ticket, error) {

	// get session
	session, err := s.client.GetSessionKey(s.cfg.ServiceURL, s.cfg.Username, s.cfg.APIKey)
	if err != nil {
		return nil, err
	}
	accout2Contact := ConvertAccount(&req.Account)
	contactModel := ConvertModelContact(accout2Contact, session.UserID)
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

	ticket := ConvertTicket(&req.Ticket)
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

	ticketResp, err := s.CreateOrUpdateVtiger(vtigerMap, session, fileMapData, "HelpDesk")
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

	ticket := ConvertTicket(&getTicketsRequest.Ticket)

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
		mapTicket, err := fileMapData.MappingTicketVtiger2Etop(value)
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
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "SELECT * FROM %v ", module)
	arrCondition := make([]string, 0, len(condition))
	if page == 0 {
		page = 1
	}
	if perpage == 0 {
		perpage = 20
	}
	fieldMap := s.fieldMap[module]
	for key, value := range condition {
		etopField := key
		if fieldMap[key] != "" {
			etopField = fieldMap[etopField]
		}
		value, err := singleQuote(value)
		if err != nil {
			return "", err
		}
		arrCondition = append(arrCondition, fmt.Sprintf(" %v = %v ", etopField, value))
	}

	if len(arrCondition) > 0 {
		_, _ = fmt.Fprint(&b, "WHERE "+strings.Join(arrCondition, " AND "))
	}
	etopField := orderBy.Field
	if fieldMap[etopField] != "" {
		etopField = fieldMap[etopField]
	}
	if orderBy.Sort == "" {
		orderBy.Sort = "DESC"
	}
	if orderBy != nil {
		_, _ = fmt.Fprintf(&b, "ORDER BY %v %v ", etopField, orderBy.Sort)
	}
	_, _ = fmt.Fprintf(&b, "LIMIT %v, %v ;", (page-1)*perpage, perpage)
	return b.String(), nil
}

// CountTicketByStatus get count number of ticket follow status
func (s *VtigerService) CountTicketByStatus(ctx context.Context, countTicketByStatusRequest *crmservice.CountTicketByStatusRequest) (*crmservice.CountTicketByStatusResponse, error) {
	status := *countTicketByStatusRequest.Status
	statusSingleQuote, err := singleQuote(status)
	if err != nil {
		return nil, err
	}
	//make SQL query vtiger
	var b strings.Builder
	fmt.Fprint(&b, "SELECT COUNT(*) FROM HelpDesk WHERE ticketstatus = ")
	fmt.Fprint(&b, statusSingleQuote)
	fmt.Fprint(&b, " ;")

	resultVtiger, err := s.VtigerRawQuery(b.String())
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
	session, err := s.client.GetSessionKey(s.cfg.ServiceURL, s.cfg.Username, s.cfg.APIKey)
	if err != nil {
		return nil, err
	}

	queryURL := make(url.Values)
	queryURL.Set("operation", "query")
	queryURL.Set("sessionName", session.SessionName)
	queryURL.Set("query", query)

	vtigerService := client.NewVigerClient(session.SessionName, s.cfg.ServiceURL)
	return vtigerService.RequestGet(queryURL)
}

// GetTicketStatusCount get ticket status count
func (s *VtigerService) GetTicketStatusCount(ctx context.Context) (*crmservice.GetTicketStatusCountResponse, error) {
	session, err := s.client.GetSessionKey(s.cfg.ServiceURL, s.cfg.Username, s.cfg.APIKey)
	if err != nil {
		return nil, err
	}
	categories := GetCategories()

	var statusCounts []crmservice.CountTicketByStatusResponse
	for _, value := range categories {
		status := value.Label
		statusSingleQuote, err := singleQuote(status)
		if err != nil {
			return nil, err
		}

		//make SQL query vtiger
		var b strings.Builder
		_, _ = fmt.Fprint(&b, "SELECT COUNT(*) FROM HelpDesk WHERE ticketcategories = ")
		_, _ = fmt.Fprint(&b, statusSingleQuote)
		_, _ = fmt.Fprint(&b, " AND ticketstatus = 'Open';")

		queryValues := make(url.Values)
		queryValues.Set("operation", "query")
		queryValues.Set("sessionName", session.SessionName)
		queryValues.Set("query", b.String())

		vtigerService := client.NewVigerClient(session.SessionName, s.cfg.ServiceURL)
		response, err := vtigerService.RequestGet(queryValues)
		if err != nil {
			return nil, err
		}
		countAtoi, err := strconv.Atoi(response.Result[0]["count"])
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
