package crm

import (
	"context"
	"fmt"
	"time"

	"o.o/api/supporting/crm/vht"
	"o.o/api/supporting/crm/vtiger"
	"o.o/api/top/services/crm"
	api "o.o/api/top/services/crm"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/shipping_provider"
	notimodel "o.o/backend/com/eventhandler/notifier/model"
	notistore "o.o/backend/com/eventhandler/notifier/sqlstore"
	"o.o/backend/com/main/shipping/carrier"
	shipmodel "o.o/backend/com/main/shipping/model"
	"o.o/backend/com/main/shipping/modelx"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

var (
	shipmentManager *carrier.ShipmentManager
	vtigerQS        vtiger.QueryBus
	vtigerAgg       vtiger.CommandBus
	vhtQS           vht.QueryBus
	vhtAgg          vht.CommandBus
)

func Init(shipmentManager *carrier.ShipmentManager,
	vtigerQuery vtiger.QueryBus,
	vtigerAggregate vtiger.CommandBus,
	vhtQuery vht.QueryBus,
	vhtAggregate vht.CommandBus) {
	shipmentManager = shipmentManager
	vtigerQS = vtigerQuery
	vtigerAgg = vtigerAggregate
	vhtQS = vhtQuery
	vhtAgg = vhtAggregate

}

type MiscService struct {
	session.Session
}

type CrmService struct {
	session.Session

	OrderStore sqlstore.OrderStoreInterface
}

type VtigerService struct {
	session.Session
}

type VhtService struct {
	session.Session
}

var miscService = &MiscService{}
var crmService = &CrmService{}
var vtigerService = &VtigerService{}
var vhtService = &VhtService{}

func (s *MiscService) Clone() api.MiscService {
	res := *s
	return &res
}

func (s *MiscService) VersionInfo(ctx context.Context, q *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	result := &pbcm.VersionInfoResponse{
		Service: "etop-crm",
		Version: "0.1",
	}
	return result, nil
}

func (s *CrmService) Clone() api.CrmService {
	res := *s
	return &res
}

func (s *CrmService) RefreshFulfillmentFromCarrier(ctx context.Context, r *api.RefreshFulfillmentFromCarrierRequest) (*pbcm.UpdatedResponse, error) {
	query := &shipmodelx.GetFulfillmentQuery{
		ShippingCode: r.ShippingCode,
	}
	if err := s.OrderStore.GetFulfillment(ctx, query); err != nil {
		return nil, err
	}
	ffm := query.Result
	var ffmUpdate *shipmodel.Fulfillment
	var err error
	switch ffm.ShippingProvider {
	case shipping_provider.GHN:
		ffmUpdate, err = shipmentManager.RefreshFulfillment(ctx, ffm)
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "This feature is not available for this carrier (%v)", ffm.ShippingProvider)
	}

	if err != nil {
		return nil, err
	}
	t0 := time.Now()
	ffmUpdate.LastSyncAt = t0
	update := &modelx.UpdateFulfillmentCommand{
		Fulfillment: ffmUpdate,
	}
	if err = s.OrderStore.UpdateFulfillment(ctx, update); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return result, nil
}

func (s *CrmService) SendNotification(ctx context.Context, r *api.SendNotificationRequest) (*pbcm.MessageResponse, error) {
	cmd := &notimodel.CreateNotificationsArgs{
		AccountIDs:       []dot.ID{r.AccountId},
		Title:            r.Title,
		Message:          r.Message,
		EntityID:         r.EntityId,
		Entity:           r.Entity,
		SendNotification: true,
		MetaData:         r.MetaData.Data,
	}
	_, _, err := notistore.NewNotificationStore(nil, nil).CreateNotifications(cmd)
	if err != nil {
		return nil, err
	}

	result := cmapi.Message("ok", fmt.Sprintf(
		"Create successful"))
	return result, nil
}

func (s *VtigerService) Clone() api.VtigerService {
	res := *s
	return &res
}

func (s *VtigerService) CreateTicket(ctx context.Context, r *api.CreateOrUpdateTicketRequest) (*api.Ticket, error) {
	cmd := &vtiger.CreateTicketCommand{
		FfmCode:     r.FfmCode,
		FfmID:       r.FfmId,
		ID:          r.Id,
		EtopUserID:  s.SS.Claim().UserID,
		Code:        r.Code,
		Title:       r.Title,
		Value:       r.Value,
		OldValue:    r.OldValue,
		Reason:      r.Reason,
		ShopID:      s.SS.Shop().ID,
		OrderID:     r.OrderId,
		OrderCode:   r.OrderCode,
		FfmUrl:      r.FfmUrl,
		Company:     r.Company,
		Provider:    r.Provider,
		Note:        r.Note,
		Environment: r.Environment,
		FromApp:     r.FromApp,
		Account: vtiger.Account{
			ID:        r.Account.Id,
			FullName:  r.Account.FullName,
			ShortName: r.Account.ShortName,
			Phone:     r.Account.Phone,
			Email:     r.Account.Email,
			Company:   r.Account.Company,
		},
	}
	if err := vtigerAgg.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.Ticket{
		TicketNo:         cmd.Result.TicketNo,
		AssignedUserId:   cmd.Result.AssignedUserId,
		ParentId:         cmd.Result.ParentID,
		Ticketpriorities: cmd.Result.Ticketpriorities,
		ProductId:        cmd.Result.ProductID,
		Ticketseverities: cmd.Result.Ticketseverities,
		Ticketstatus:     cmd.Result.Ticketstatus,
		Ticketcategories: cmd.Result.Ticketcategories,
		UpdateLog:        cmd.Result.UpdateLog,
		Hours:            cmd.Result.Hours,
		Days:             cmd.Result.Days,
		Createdtime:      cmapi.PbTime(cmd.Result.CreatedTime),
		Modifiedtime:     cmapi.PbTime(cmd.Result.ModifiedTime),
		FromPortal:       cmd.Result.FromPortal,
		Modifiedby:       cmd.Result.Modifiedby,
		TicketTitle:      cmd.Result.TicketTitle,
		Description:      cmd.Result.Description,
		Solution:         cmd.Result.Solution,
		ContactId:        cmd.Result.ContactId,
		Source:           cmd.Result.Source,
		Starred:          cmd.Result.Starred,
		Tags:             cmd.Result.Tags,
		Note:             cmd.Result.Note,
		FfmCode:          cmd.Result.FfmCode,
		FfmUrl:           cmd.Result.FfmUrl,
		FfmId:            cmd.Result.FfmId,
		EtopUserId:       cmd.Result.EtopUserID,
		OrderId:          cmd.Result.OrderId,
		OrderCode:        cmd.Result.OrderCode,
		Company:          cmd.Result.Company,
		Provider:         cmd.Result.Provider,
		FromApp:          cmd.Result.FromApp,
		Environment:      cmd.Result.Environment,
		Code:             cmd.Result.Code,
		OldValue:         cmd.Result.OldValue,
		NewValue:         cmd.Result.NewValue,
		Substatus:        cmd.Result.Substatus,
		EtopNote:         cmd.Result.EtopNote,
		Reason:           cmd.Result.Reason,
		Id:               cmd.Result.ID,
	}
	return result, nil
}

func (s *VtigerService) UpdateTicket(ctx context.Context, r *api.CreateOrUpdateTicketRequest) (*api.Ticket, error) {
	cmd := &vtiger.UpdateTicketCommand{
		FfmCode:     r.FfmCode,
		FfmID:       r.FfmId,
		ID:          r.Id,
		EtopUserID:  s.SS.Claim().UserID,
		Code:        r.Code,
		Title:       r.Title,
		Value:       r.Value,
		OldValue:    r.OldValue,
		Reason:      r.Reason,
		ShopID:      s.SS.Shop().ID,
		OrderID:     r.OrderId,
		OrderCode:   r.OrderCode,
		FfmUrl:      r.FfmUrl,
		Company:     r.Company,
		Provider:    r.Provider,
		Note:        r.Note,
		Environment: r.Environment,
		FromApp:     r.FromApp,
	}
	if err := vtigerAgg.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.Ticket{
		TicketNo:         cmd.Result.TicketNo,
		AssignedUserId:   cmd.Result.AssignedUserId,
		ParentId:         cmd.Result.ParentID,
		Ticketpriorities: cmd.Result.Ticketpriorities,
		ProductId:        cmd.Result.ProductID,
		Ticketseverities: cmd.Result.Ticketseverities,
		Ticketstatus:     cmd.Result.Ticketstatus,
		Ticketcategories: cmd.Result.Ticketcategories,
		UpdateLog:        cmd.Result.UpdateLog,
		Hours:            cmd.Result.Hours,
		Days:             cmd.Result.Days,
		Createdtime:      cmapi.PbTime(cmd.Result.CreatedTime),
		Modifiedtime:     cmapi.PbTime(cmd.Result.ModifiedTime),
		FromPortal:       cmd.Result.FromPortal,
		Modifiedby:       cmd.Result.Modifiedby,
		TicketTitle:      cmd.Result.TicketTitle,
		Description:      cmd.Result.Description,
		Solution:         cmd.Result.Solution,
		ContactId:        cmd.Result.ContactId,
		Source:           cmd.Result.Source,
		Starred:          cmd.Result.Starred,
		Tags:             cmd.Result.Tags,
		Note:             cmd.Result.Note,
		FfmCode:          cmd.Result.FfmCode,
		FfmUrl:           cmd.Result.FfmUrl,
		FfmId:            cmd.Result.FfmId,
		EtopUserId:       cmd.Result.EtopUserID,
		OrderId:          cmd.Result.OrderId,
		OrderCode:        cmd.Result.OrderCode,
		Company:          cmd.Result.Company,
		Provider:         cmd.Result.Provider,
		FromApp:          cmd.Result.FromApp,
		Environment:      cmd.Result.Environment,
		Code:             cmd.Result.Code,
		OldValue:         cmd.Result.OldValue,
		NewValue:         cmd.Result.NewValue,
		Substatus:        cmd.Result.Substatus,
		EtopNote:         cmd.Result.EtopNote,
		Reason:           cmd.Result.Reason,
		Id:               cmd.Result.ID,
	}
	return result, nil
}

func (s *VtigerService) CreateOrUpdateContact(ctx context.Context, r *api.ContactRequest) (*api.ContactResponse, error) {
	cmd := &vtiger.CreateOrUpdateContactCommand{
		ID:                   r.Id,
		EtopUserID:           s.SS.Claim().UserID,
		ContactNo:            r.ContactNo,
		Phone:                r.Phone,
		Lastname:             r.Lastname,
		Mobile:               r.Mobile,
		Leadsource:           r.Leadsource,
		Email:                r.Email,
		Description:          r.Description,
		Secondaryemail:       r.Secondaryemail,
		Modifiedby:           r.Modifiedby,
		Source:               r.Source,
		Company:              r.Company,
		Website:              r.Website,
		Lane:                 r.Lane,
		City:                 r.City,
		State:                r.State,
		Country:              r.Country,
		OrdersPerDay:         r.OrdersPerDay,
		UsedShippingProvider: r.UsedShippingProvider,
		Firstname:            r.Firstname,
		AssignedUserID:       r.AssignedUserId,
	}
	if err := vtigerAgg.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.ContactResponse{
		ContactNo:            cmd.Result.ContactNo,
		Phone:                cmd.Result.Phone,
		Lastname:             cmd.Result.Lastname,
		Mobile:               cmd.Result.Mobile,
		Leadsource:           cmd.Result.Leadsource,
		Email:                cmd.Result.Email,
		Description:          cmd.Result.Description,
		Secondaryemail:       cmd.Result.Secondaryemail,
		Modifiedby:           cmd.Result.Modifiedby,
		Source:               cmd.Result.Source,
		EtopUserId:           cmd.Result.EtopUserID,
		Company:              cmd.Result.Company,
		Website:              cmd.Result.Website,
		Lane:                 cmd.Result.Lane,
		City:                 cmd.Result.City,
		State:                cmd.Result.State,
		Country:              cmd.Result.Country,
		OrdersPerDay:         cmd.Result.OrdersPerDay,
		UsedShippingProvider: cmd.Result.UsedShippingProvider,
		Id:                   cmd.Result.ID,
		Firstname:            cmd.Result.Firstname,
		Createdtime:          cmapi.PbTime(cmd.Result.Createdtime),
		Modifiedtime:         cmapi.PbTime(cmd.Result.Createdtime),
		AssignedUserId:       cmd.Result.AssignedUserID,
	}
	return result, nil
}

func (s *VtigerService) CreateOrUpdateLead(ctx context.Context, r *api.LeadRequest) (*api.LeadResponse, error) {
	cmd := &vtiger.CreateOrUpdateLeadCommand{
		ID:                   r.Id,
		EtopUserID:           s.SS.Claim().UserID,
		ContactNo:            r.ContactNo,
		Phone:                r.Phone,
		Lastname:             r.Lastname,
		Mobile:               r.Mobile,
		Leadsource:           r.Leadsource,
		Email:                r.Email,
		Description:          r.Description,
		Secondaryemail:       r.Secondaryemail,
		Modifiedby:           r.Modifiedby,
		Source:               r.Source,
		Company:              r.Company,
		Website:              r.Website,
		Lane:                 r.Lane,
		City:                 r.City,
		State:                r.State,
		Country:              r.Country,
		OrdersPerDay:         r.OrdersPerDay,
		UsedShippingProvider: r.UsedShippingProvider,
		Firstname:            r.Firstname,
		AssignedUserID:       r.AssignedUserId,
	}
	if err := vtigerAgg.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.LeadResponse{
		ContactNo:            cmd.Result.ContactNo,
		Phone:                cmd.Result.Phone,
		Lastname:             cmd.Result.Lastname,
		Mobile:               cmd.Result.Mobile,
		Leadsource:           cmd.Result.Leadsource,
		Email:                cmd.Result.Email,
		Secondaryemail:       cmd.Result.Secondaryemail,
		AssignedUserId:       cmd.Result.AssignedUserID,
		Description:          cmd.Result.Description,
		Modifiedby:           cmd.Result.Modifiedby,
		Source:               cmd.Result.Source,
		EtopUserId:           cmd.Result.EtopUserID,
		Company:              cmd.Result.Company,
		Website:              cmd.Result.Website,
		Lane:                 cmd.Result.Lane,
		City:                 cmd.Result.City,
		State:                cmd.Result.State,
		Country:              cmd.Result.Country,
		OrdersPerDay:         cmd.Result.OrdersPerDay,
		UsedShippingProvider: cmd.Result.UsedShippingProvider,
		Id:                   cmd.Result.ID,
		Firstname:            cmd.Result.Firstname,
	}
	return result, nil
}

func (s *VhtService) Clone() api.VhtService {
	res := *s
	return &res
}

func (s *VhtService) CreateOrUpdateCallHistoryBySDKCallID(ctx context.Context, r *api.VHTCallLog) (*api.VHTCallLog, error) {
	cmd := &vht.CreateOrUpdateCallHistoryByCallIDCommand{
		Direction:       r.Direction,
		CdrID:           r.CdrId,
		CallID:          r.CallId,
		SipCallID:       r.SipCallId,
		SdkCallID:       r.SdkCallId,
		Cause:           r.Cause,
		Q850Cause:       r.Q850Cause,
		FromExtension:   r.FromExtension,
		ToExtension:     r.ToExtension,
		FromNumber:      r.FromNumber,
		ToNumber:        r.ToNumber,
		Duration:        r.Duration,
		TimeStarted:     cmapi.PbTimeToModel(r.TimeStarted),
		TimeConnected:   cmapi.PbTimeToModel(r.TimeConnected),
		TimeEnded:       cmapi.PbTimeToModel(r.TimeEnded),
		RecordingPath:   r.RecordingPath,
		RecordingUrl:    r.RecordingUrl,
		RecordFileSize:  r.RecordFileSize,
		EtopAccountID:   r.EtopAccountId,
		VtigerAccountID: r.VtigerAccountId,
	}
	if err := vhtAgg.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	// TODO: convert result
	return nil, nil
}

func (s *VhtService) CreateOrUpdateCallHistoryByCallID(ctx context.Context, r *api.VHTCallLog) (*api.VHTCallLog, error) {
	cmd := &vht.CreateOrUpdateCallHistoryByCallIDCommand{
		Direction:       r.Direction,
		CdrID:           r.CdrId,
		CallID:          r.CallId,
		SipCallID:       r.SipCallId,
		SdkCallID:       r.SdkCallId,
		Cause:           r.Cause,
		Q850Cause:       r.Q850Cause,
		FromExtension:   r.FromExtension,
		ToExtension:     r.ToExtension,
		FromNumber:      r.FromNumber,
		ToNumber:        r.ToNumber,
		Duration:        r.Duration,
		TimeStarted:     cmapi.PbTimeToModel(r.TimeStarted),
		TimeConnected:   cmapi.PbTimeToModel(r.TimeConnected),
		TimeEnded:       cmapi.PbTimeToModel(r.TimeEnded),
		RecordingPath:   r.RecordingPath,
		RecordingUrl:    r.RecordingUrl,
		RecordFileSize:  r.RecordFileSize,
		EtopAccountID:   r.EtopAccountId,
		VtigerAccountID: r.VtigerAccountId,
	}
	if err := vhtAgg.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	// TODO: convert result
	return nil, nil
}

func (s *VhtService) GetCallHistories(ctx context.Context, r *api.GetCallHistoriesRequest) (*api.GetCallHistoriesResponse, error) {
	query := &vht.GetCallHistoriesQuery{
		Paging:     cmapi.PagingToModel(r.Paging, 1, 100, 1000),
		TextSearch: r.TextSearch,
	}
	if err := vhtQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var vhtCallLog []*crm.VHTCallLog
	for _, value := range query.Result.VhtCallLog {
		vhtCallLog = append(vhtCallLog, &crm.VHTCallLog{
			CdrId:           value.CdrID,
			CallId:          value.CallID,
			SipCallId:       value.SipCallID,
			SdkCallId:       value.SdkCallID,
			Cause:           value.Cause,
			Q850Cause:       value.Q850Cause,
			FromExtension:   value.FromExtension,
			ToExtension:     value.ToExtension,
			FromNumber:      value.FromExtension,
			ToNumber:        value.ToNumber,
			Duration:        value.Duration,
			Direction:       value.Direction,
			TimeStarted:     cmapi.PbTime(value.TimeStarted),
			TimeConnected:   cmapi.PbTime(value.TimeConnected),
			TimeEnded:       cmapi.PbTime(value.TimeEnded),
			RecordingPath:   value.RecordingPath,
			RecordingUrl:    value.RecordingUrl,
			RecordFileSize:  value.RecordFileSize,
			EtopAccountId:   value.EtopAccountID,
			VtigerAccountId: value.VtigerAccountID,
		})
	}
	result := &api.GetCallHistoriesResponse{
		VhtCallLog: vhtCallLog,
	}
	return result, nil
}

func (s *VtigerService) GetContacts(ctx context.Context, r *api.GetContactsRequest) (*api.GetContactsResponse, error) {
	query := &vtiger.GetContactsQuery{
		Search: r.TextSearch,
		Paging: cmapi.PagingToModel(r.Paging, 1, 100, 1000),
	}
	if err := vtigerQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var arrContact []*crm.ContactResponse
	for _, value := range query.Result.Contacts {
		arrContact = append(arrContact, &crm.ContactResponse{
			ContactNo:            value.ContactNo,
			Phone:                value.Phone,
			Lastname:             value.Lastname,
			Mobile:               value.Mobile,
			Leadsource:           value.Leadsource,
			Email:                value.Email,
			Description:          value.Description,
			Secondaryemail:       value.Secondaryemail,
			Modifiedby:           value.Modifiedby,
			Source:               value.Source,
			EtopUserId:           value.EtopUserID,
			Company:              value.Company,
			Website:              value.Website,
			Lane:                 value.Lane,
			City:                 value.City,
			State:                value.State,
			Country:              value.Country,
			OrdersPerDay:         value.OrdersPerDay,
			UsedShippingProvider: value.UsedShippingProvider,
			Id:                   value.ID,
			Firstname:            value.Firstname,
			Createdtime:          cmapi.PbTime(value.Createdtime),
			Modifiedtime:         cmapi.PbTime(value.Createdtime),
			AssignedUserId:       value.AssignedUserID,
		})
	}
	result := &api.GetContactsResponse{
		Contacts: arrContact,
	}

	return result, nil
}

func (s *VtigerService) GetTickets(ctx context.Context, r *api.GetTicketsRequest) (*api.GetTicketsResponse, error) {
	paging := cmapi.PagingToModel(r.Paging, 1, 250, 250)
	query := &vtiger.GetTicketsQuery{
		Paging: paging,
		Ticket: vtiger.TicketArgs{
			ID:          r.Ticket.Id,
			EtopUserID:  r.Ticket.EtopUserId,
			Code:        r.Ticket.Code,
			Title:       r.Ticket.Title,
			Value:       r.Ticket.Value,
			OldValue:    r.Ticket.OldValue,
			Reason:      r.Ticket.Reason,
			OrderID:     r.Ticket.OrderId,
			OrderCode:   r.Ticket.OrderCode,
			FfmCode:     r.Ticket.FfmCode,
			FfmUrl:      r.Ticket.FfmUrl,
			FfmID:       r.Ticket.FfmId,
			Company:     r.Ticket.Company,
			Provider:    r.Ticket.Provider,
			Note:        r.Ticket.Note,
			Environment: r.Ticket.Environment,
			FromApp:     r.Ticket.FromApp,
		},
		Orderby: vtiger.OrderBy{
			Field: r.Orderby.Field,
			Sort:  r.Orderby.Sort,
		},
		Result: nil,
	}
	if err := vtigerQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var arrTicket []*crm.Ticket
	for _, value := range query.Result.Tickets {
		arrTicket = append(arrTicket, &crm.Ticket{
			TicketNo:         value.TicketNo,
			AssignedUserId:   value.AssignedUserId,
			ParentId:         value.ParentID,
			Ticketpriorities: value.Ticketpriorities,
			ProductId:        value.ProductID,
			Ticketseverities: value.Ticketseverities,
			Ticketstatus:     value.Ticketstatus,
			Ticketcategories: value.Ticketcategories,
			UpdateLog:        value.UpdateLog,
			Hours:            value.Hours,
			Days:             value.Days,
			Createdtime:      cmapi.PbTime(value.CreatedTime),
			Modifiedtime:     cmapi.PbTime(value.ModifiedTime),
			FromPortal:       value.FromPortal,
			Modifiedby:       value.Modifiedby,
			TicketTitle:      value.TicketTitle,
			Description:      value.Description,
			Solution:         value.Solution,
			ContactId:        value.ContactId,
			Source:           value.Source,
			Starred:          value.Starred,
			Tags:             value.Tags,
			Note:             value.Note,
			FfmCode:          value.FfmCode,
			FfmUrl:           value.FfmUrl,
			FfmId:            value.FfmId,
			EtopUserId:       value.EtopUserID,
			OrderId:          value.OrderId,
			OrderCode:        value.OrderCode,
			Company:          value.Company,
			Provider:         value.Provider,
			FromApp:          value.FromApp,
			Environment:      value.Environment,
			Code:             value.Code,
			OldValue:         value.OldValue,
			NewValue:         value.NewValue,
			Substatus:        value.Substatus,
			EtopNote:         value.EtopNote,
			Reason:           value.Reason,
			Id:               value.ID,
		})
	}
	result := &api.GetTicketsResponse{
		Tickets: arrTicket,
	}
	return result, nil
}

func (s *VtigerService) GetCategories(ctx context.Context, r *pbcm.Empty) (*api.GetCategoriesResponse, error) {
	query := &vtiger.GetCategoriesQuery{}
	if err := vtigerQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var arrCategory []*crm.Category
	for _, value := range query.Result.Categories {
		arrCategory = append(arrCategory, &crm.Category{
			Code:  value.Code,
			Label: value.Label,
		})
	}
	result := &api.GetCategoriesResponse{
		Categories: arrCategory,
	}
	return result, nil
}

func (s *VtigerService) GetTicketStatusCount(ctx context.Context, r *pbcm.Empty) (*api.GetTicketStatusCountResponse, error) {
	query := &vtiger.GetTicketStatusCountQuery{}
	if err := vtigerQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var arrStatus []*crm.CountTicketByStatusResponse
	for _, value := range query.Result.StatusCount {
		arrStatus = append(arrStatus, &crm.CountTicketByStatusResponse{
			Code:  value.Code,
			Count: value.Count,
		})
	}
	result := &api.GetTicketStatusCountResponse{
		StatusCount: arrStatus,
	}
	return result, nil
}
