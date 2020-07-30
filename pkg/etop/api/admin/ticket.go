package admin

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	"o.o/api/supporting/ticket"
	api "o.o/api/top/int/admin"
	shoptypes "o.o/api/top/int/shop/types"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type TicketService struct {
	session.Session

	TicketQuery    ticket.QueryBus
	TicketAggr     ticket.CommandBus
	IndentityQuery identity.QueryBus
}

func (s *TicketService) Clone() api.TicketService { res := *s; return &res }

func (s *TicketService) GetTicket(ctx context.Context, request *api.GetTicketRequest) (*shoptypes.Ticket, error) {
	queryAccount := &identity.GetAccountByIDQuery{
		ID: request.AccountID,
	}
	err := s.IndentityQuery.Dispatch(ctx, queryAccount)
	if err != nil {
		return nil, err
	}
	query := &ticket.GetTicketByIDQuery{
		ID:        request.ID,
		AccountID: request.AccountID,
	}
	err = s.TicketQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_Ticket_to_api_Ticket(query.Result), nil
}

func (s *TicketService) ReopenTicket(ctx context.Context, request *api.ReopenTicketRequest) (*shoptypes.Ticket, error) {
	cmd := &ticket.ReopenTicketCommand{
		ID:   request.TicketID,
		Note: request.Note,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) CreateTicketComment(ctx context.Context, request *api.CreateTicketCommentRequest) (*shoptypes.TicketComment, error) {
	isLeader := false
	for _, role := range s.SS.Permission().Roles {
		if role == string(authorization.RoleAdminCustomerServiceLead) || role == string(authorization.RoleAdmin) || role == string(authorization.RoleAdminSaleLead) {
			isLeader = true
			break
		}
	}
	cmd := &ticket.CreateTicketCommentCommand{
		CreatedBy: s.SS.User().ID,
		TicketID:  request.TicketID,
		AccountID: request.AccountID,
		ParentID:  request.ParentID,
		Message:   request.Message,
		ImageUrl:  request.ImageUrl,
		IsLeader:  isLeader,
		IsAdmin:   true,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_TicketComment_to_api_TicketComment(cmd.Result), nil
}

func (s *TicketService) UpdateTicketComment(ctx context.Context, request *api.UpdateTicketCommentRequest) (*shoptypes.TicketComment, error) {
	cmd := &ticket.UpdateTicketCommentCommand{
		AccountID: request.AccountID,
		ID:        request.ID,
		UpdatedBy: s.SS.User().ID,
		Message:   request.Message,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_TicketComment_to_api_TicketComment(cmd.Result), nil
}

func (s *TicketService) GetTicketComments(ctx context.Context, request *api.GetTicketCommentsRequest) (*api.GetTicketCommentsResponse, error) {
	var filter = &ticket.FilterGetTicketComment{}
	if request.Filter != nil {
		filter.CreatedBy = request.Filter.CreatedBy
		filter.ParentID = request.Filter.ParentID
		filter.Title = request.Filter.Title
		filter.IDs = request.Filter.IDs
		filter.AccountID = request.Filter.AccountID
		filter.TicketID = request.Filter.TicketID
	}
	paging := cmapi.CMPaging(request.Paging)
	query := &ticket.ListTicketCommentsQuery{
		Filter: filter,
		Paging: *paging,
	}
	err := s.TicketQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	return &api.GetTicketCommentsResponse{
		TicketComments: convertpb.Convert_core_TicketComments_to_api_TicketComments(query.Result.TicketComments),
		Paging:         cmapi.PbPaging(query.Paging),
	}, nil
}

func (s *TicketService) CreateTicket(ctx context.Context, request *api.CreateTicketRequest) (*shoptypes.Ticket, error) {
	userID := s.SS.User().ID
	query := &identity.GetAccountByIDQuery{
		ID: request.AccountID,
	}
	err := s.IndentityQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	cmd := &ticket.CreateTicketCommand{
		AccountID:   request.AccountID,
		LabelIDs:    request.LabelIDs,
		Title:       request.Title,
		Description: request.Description,
		AdminNote:   request.Note,
		RefID:       request.RefID,
		RefType:     request.RefType,
		RefCode:     request.RefCode,
		Source:      request.Source,
		CreatedBy:   userID,
	}
	err = s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) GetTickets(ctx context.Context, request *api.GetTicketsRequest) (*api.GetTicketsResponse, error) {
	var filter = &ticket.FilterGetTicket{}
	if request.Filter != nil {
		filter.Title = request.Filter.Title
		filter.LabelIDs = request.Filter.LabelIDs
		filter.CreatedBy = request.Filter.CreatedBy
		filter.ClosedBy = request.Filter.ClosedBy
		filter.IDs = request.Filter.IDs
		filter.AccountID = request.Filter.AccountID
		filter.AssignedUserIDs = request.Filter.AssignedUserID
		filter.Code = request.Filter.Code
		filter.RefID = request.Filter.RefID
		filter.State = request.Filter.State
		filter.RefType = request.Filter.RefType
		filter.RefCode = request.Filter.RefCode
		if filter.AccountID != 0 {
			query := &identity.GetAccountByIDQuery{
				ID: request.Filter.AccountID,
			}
			err := s.IndentityQuery.Dispatch(ctx, query)
			if err != nil {
				return nil, err
			}
		}

	}
	paging := cmapi.CMPaging(request.Paging)
	query := &ticket.ListTicketsQuery{
		Filter: filter,
		Paging: *paging,
	}
	err := s.TicketQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	result := convertpb.Convert_core_Tickets_to_api_Tickets(query.Result.Tickets)
	return &api.GetTicketsResponse{
		Paging:  cmapi.PbPaging(query.Paging),
		Tickets: result,
	}, nil
}

func (s *TicketService) AssignTicket(ctx context.Context, request *api.AssignTicketRequest) (*shoptypes.Ticket, error) {
	isLeader := false
	for _, role := range s.SS.Permission().Roles {
		if role == string(authorization.RoleAdminCustomerServiceLead) || role == string(authorization.RoleAdmin) || role == string(authorization.RoleAdminSaleLead) {
			isLeader = true
			break
		}
	}
	cmd := &ticket.AssignTicketCommand{
		ID:              request.TicketID,
		UpdatedBy:       s.SS.User().ID,
		IsLeader:        isLeader,
		AssignedUserIDs: request.AssignedUserIDs,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) UnassignTicket(ctx context.Context, request *api.AssignTicketRequest) (*shoptypes.Ticket, error) {
	cmd := &ticket.UnassignTicketCommand{
		ID:        request.TicketID,
		UpdatedBy: s.SS.User().ID,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) ConfirmTicket(ctx context.Context, request *api.ConfirmTicketRequest) (*shoptypes.Ticket, error) {
	isLeader := false
	for _, role := range s.SS.Permission().Roles {
		if role == string(authorization.RoleAdminCustomerServiceLead) || role == string(authorization.RoleAdmin) || role == string(authorization.RoleAdminSaleLead) {
			isLeader = true
			break
		}
	}
	cmd := &ticket.ConfirmTicketCommand{
		IsLeader:  isLeader,
		ID:        request.TicketID,
		ConfirmBy: s.SS.User().ID,
		Note:      request.Note,
		Result:    nil,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) CloseTicket(ctx context.Context, request *api.CloseTicketRequest) (*shoptypes.Ticket, error) {
	isLeader := false
	for _, role := range s.SS.Permission().Roles {
		if role == string(authorization.RoleAdminCustomerServiceLead) || role == string(authorization.RoleAdmin) || role == string(authorization.RoleAdminSaleLead) {
			isLeader = true
			break
		}
	}
	cmd := &ticket.CloseTicketCommand{
		IsLeader: isLeader,
		ID:       request.TicketID,
		ClosedBy: s.SS.User().ID,
		Note:     request.Note,
		State:    request.State,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) CreateTicketLabel(ctx context.Context, request *api.CreateTicketLabelRequest) (*shoptypes.TicketLabel, error) {
	cmd := &ticket.CreateTicketLabelCommand{
		Name:     request.Name,
		Code:     request.Code,
		Color:    request.Color,
		ParentID: request.ParentID,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_TicketLabel_to_api_TicketLabel(cmd.Result), nil
}

func (s *TicketService) UpdateTicketLabel(ctx context.Context, request *api.UpdateTicketLabelRequest) (*shoptypes.TicketLabel, error) {
	cmd := &ticket.UpdateTicketLabelCommand{
		ID:       request.ID,
		Name:     request.Name,
		Color:    request.Color,
		Code:     request.Code,
		ParentID: request.ParentID,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpb.Convert_core_TicketLabel_to_api_TicketLabel(cmd.Result), nil
}

func (s *TicketService) DeleteTicketLabel(ctx context.Context, request *api.DeleteTicketLabelRequest) (*api.DeleteTicketLabelResponse, error) {
	cmd := &ticket.DeleteTicketLabelCommand{
		ID:          request.ID,
		DeleteChild: request.DeleteChild,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &api.DeleteTicketLabelResponse{
		Count: 1,
	}, nil
}
