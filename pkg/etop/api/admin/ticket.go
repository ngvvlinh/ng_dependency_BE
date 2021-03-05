package admin

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	"o.o/api/supporting/ticket"
	api "o.o/api/top/int/admin"
	shoptypes "o.o/api/top/int/shop/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/ticket/ticket_type"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
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
	if err := s.IndentityQuery.Dispatch(ctx, queryAccount); err != nil {
		return nil, err
	}

	query := &ticket.GetTicketByIDQuery{
		ID:        request.ID,
		AccountID: request.AccountID,
	}
	if err := s.TicketQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	ticket := convertpball.Convert_core_Ticket_to_api_Ticket(query.Result)
	return ticket, nil
}

func (s *TicketService) ReopenTicket(ctx context.Context, request *api.ReopenTicketRequest) (*shoptypes.Ticket, error) {
	cmd := &ticket.ReopenTicketCommand{
		ID:   request.TicketID,
		Note: request.Note,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) CreateTicketComment(ctx context.Context, request *api.CreateTicketCommentRequest) (*shoptypes.TicketComment, error) {
	isLeader := false
	for _, role := range s.SS.Permission().Roles {
		if role == string(authorization.RoleAdminCustomerServiceLead) || role == string(authorization.RoleAdmin) || role == string(authorization.RoleAdminSaleLead) {
			isLeader = true
			break
		}
	}

	var imageUrls []string
	if len(request.ImageUrls) > 0 {
		imageUrls = append(imageUrls, request.ImageUrls...)
	} else if request.ImageUrl != "" {
		imageUrls = append(imageUrls, request.ImageUrl)
	}
	cmd := &ticket.CreateTicketCommentCommand{
		CreatedBy:     s.SS.User().ID,
		CreatedName:   s.SS.User().FullName,
		CreatedSource: account_type.Etop,
		TicketID:      request.TicketID,
		AccountID:     request.AccountID,
		ParentID:      request.ParentID,
		Message:       request.Message,
		ImageUrls:     imageUrls,
		IsLeader:      isLeader,
		IsAdmin:       true,
	}
	err := s.TicketAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return convertpball.Convert_core_TicketComment_to_api_TicketComment(cmd.Result), nil
}

func (s *TicketService) UpdateTicketComment(ctx context.Context, request *api.UpdateTicketCommentRequest) (*shoptypes.TicketComment, error) {
	var imageUrls []string
	if len(request.ImageUrls) > 0 {
		imageUrls = append(imageUrls, request.ImageUrls...)
	} else if request.ImageUrl != "" {
		imageUrls = append(imageUrls, request.ImageUrl)
	}
	cmd := &ticket.UpdateTicketCommentCommand{
		AccountID: request.AccountID,
		ID:        request.ID,
		UpdatedBy: s.SS.User().ID,
		Message:   request.Message,
		ImageUrls: imageUrls,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_TicketComment_to_api_TicketComment(cmd.Result), nil
}

func (s *TicketService) DeleteTicketComment(ctx context.Context, req *api.DeleteTicketCommentRequest) (*api.DeleteTicketCommentResponse, error) {
	cmd := &ticket.DeleteTicketCommentCommand{
		AccountID: s.SS.Shop().ID,
		ID:        req.ID,
		IsAdmin:   true,
		DeletedBy: s.SS.User().ID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &api.DeleteTicketCommentResponse{Count: cmd.Result}, nil
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
	if err := s.TicketQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	return &api.GetTicketCommentsResponse{
		TicketComments: convertpball.Convert_core_TicketComments_to_api_TicketComments(query.Result.TicketComments),
		Paging:         cmapi.PbPaging(query.Paging),
	}, nil
}

func (s *TicketService) CreateTicket(ctx context.Context, request *api.CreateTicketRequest) (*shoptypes.Ticket, error) {
	userID := s.SS.User().ID
	query := &identity.GetAccountByIDQuery{
		ID: request.AccountID,
	}
	if err := s.IndentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	cmd := &ticket.CreateTicketCommand{
		AccountID:     request.AccountID,
		LabelIDs:      request.LabelIDs,
		Title:         request.Title,
		Description:   request.Description,
		AdminNote:     request.Note,
		RefID:         request.RefID,
		RefType:       request.RefType,
		RefCode:       request.RefCode,
		Source:        request.Source,
		CreatedBy:     userID,
		CreatedName:   s.SS.User().FullName,
		CreatedSource: account_type.Etop,
		RefTicketID:   request.RefTicketID,
		Type:          ticket_type.System,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
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
			if err := s.IndentityQuery.Dispatch(ctx, query); err != nil {
				return nil, err
			}
		}

	}
	filter.Types = []ticket_type.TicketType{ticket_type.System}
	paging := cmapi.CMPaging(request.Paging)
	query := &ticket.ListTicketsQuery{
		Filter: filter,
		Paging: *paging,
	}
	if err := s.TicketQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	tickets := convertpball.Convert_core_Tickets_to_api_Tickets(query.Result.Tickets)
	return &api.GetTicketsResponse{
		Paging:  cmapi.PbPaging(query.Paging),
		Tickets: tickets,
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
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) UnassignTicket(ctx context.Context, request *api.AssignTicketRequest) (*shoptypes.Ticket, error) {
	cmd := &ticket.UnassignTicketCommand{
		ID:        request.TicketID,
		UpdatedBy: s.SS.User().ID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
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
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
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
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) GetTicketLabels(ctx context.Context, req *api.GetTicketLabelsRequest) (*api.GetTicketLabelsResponse, error) {
	query := &ticket.ListTicketLabelsQuery{
		Tree: req.Tree,
		Type: ticket_type.System.Wrap(),
	}

	err := s.TicketQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	return &api.GetTicketLabelsResponse{
		TicketLabels: convertpball.Convert_core_TicketLabels_to_api_TicketLabels(query.Result.TicketLabels),
	}, nil
}

func (s *TicketService) CreateTicketLabel(ctx context.Context, request *api.CreateTicketLabelRequest) (*shoptypes.TicketLabel, error) {
	cmd := &ticket.CreateTicketLabelCommand{
		Type:     ticket_type.System,
		Name:     request.Name,
		Code:     request.Code,
		Color:    request.Color,
		ParentID: request.ParentID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_TicketLabel_to_api_TicketLabel(cmd.Result), nil
}

func (s *TicketService) UpdateTicketLabel(ctx context.Context, request *api.UpdateTicketLabelRequest) (*shoptypes.TicketLabel, error) {
	cmd := &ticket.UpdateTicketLabelCommand{
		ID:       request.ID,
		Type:     ticket_type.System,
		Name:     request.Name,
		Color:    request.Color,
		Code:     request.Code,
		ParentID: request.ParentID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_TicketLabel_to_api_TicketLabel(cmd.Result), nil
}

func (s *TicketService) DeleteTicketLabel(ctx context.Context, request *api.DeleteTicketLabelRequest) (*api.DeleteTicketLabelResponse, error) {
	cmd := &ticket.DeleteTicketLabelCommand{
		ID:          request.ID,
		Type:        ticket_type.System,
		DeleteChild: request.DeleteChild,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return &api.DeleteTicketLabelResponse{
		Count: 1,
	}, nil
}

func (s *TicketService) GetTicketsByRefTicketID(ctx context.Context, r *shoptypes.GetTicketsByRefTicketIDRequest) (*shoptypes.GetTicketsByRefTicketIDResponse, error) {
	query := &ticket.ListTicketsByRefTicketIDQuery{
		RefTicketID: r.RefTicketID,
	}
	if err := s.TicketQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &shoptypes.GetTicketsByRefTicketIDResponse{
		Tickets: convertpball.Convert_core_Tickets_to_api_Tickets(query.Result),
	}, nil
}

func (s *TicketService) UpdateTicketRefTicketID(ctx context.Context, r *api.UpdateTicketRefTicketIDRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &ticket.UpdateTicketRefTicketIDCommand{
		ID:          r.ID,
		RefTicketID: r.RefTicketID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return cmd.Result, nil
}
