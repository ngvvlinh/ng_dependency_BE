package ticket

import (
	"context"

	"o.o/api/supporting/ticket"
	api "o.o/api/top/int/shop"
	shoptypes "o.o/api/top/int/shop/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/ticket/ticket_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/authorize/session"
)

type TicketService struct {
	session.Session

	TicketQuery ticket.QueryBus
	TicketAggr  ticket.CommandBus
}

func (s *TicketService) Clone() api.TicketService { res := *s; return &res }

func (s *TicketService) GetTickets(ctx context.Context, request *api.GetTicketsRequest) (*api.GetTicketsResponse, error) {
	shopID := s.SS.Shop().ID
	var filter = &ticket.FilterGetTicket{
		AccountID: shopID,
	}
	paging := cmapi.CMPaging(request.Paging)
	if request.Filter != nil {
		filter.Title = request.Filter.Title
		filter.LabelIDs = request.Filter.LabelIDs
		filter.CreatedBy = request.Filter.CreatedBy
		filter.ConfirmedBy = request.Filter.ConfirmedBy
		filter.Code = request.Filter.Code
		filter.RefID = request.Filter.RefID
		filter.RefType = request.Filter.RefType
		filter.State = request.Filter.State
		filter.RefCode = request.Filter.RefCode
		filter.AssignedUserIDs = request.Filter.AssignedUserIDs
		if len(request.Filter.Types) != 0 {
			filter.Types = request.Filter.Types
		} else {
			filter.Types = []ticket_type.TicketType{ticket_type.System}
		}
	}
	filter.AccountID = s.SS.Shop().ID
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

func (s *TicketService) GetTicket(ctx context.Context, request *api.GetTicketRequest) (*shoptypes.Ticket, error) {
	query := &ticket.GetTicketByIDQuery{
		ID:        request.ID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.TicketQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	ticket := convertpball.Convert_core_Ticket_to_api_Ticket(query.Result)
	return ticket, nil
}

func (s *TicketService) CreateTicket(ctx context.Context, request *api.CreateTicketRequest) (*shoptypes.Ticket, error) {
	shopID := s.SS.Shop().ID
	userID := s.SS.User().ID
	cmd := &ticket.CreateTicketCommand{
		AccountID:     shopID,
		LabelIDs:      request.LabelIDs,
		Title:         request.Title,
		Description:   request.Description,
		Note:          request.Note,
		RefID:         request.RefID,
		RefType:       request.RefType,
		RefCode:       request.RefCode,
		Source:        request.Source,
		CreatedBy:     userID,
		CreatedName:   s.SS.User().FullName,
		CreatedSource: account_type.Shop,
	}
	if request.Type.Valid {
		cmd.Type = request.Type.Enum
	} else {
		cmd.Type = ticket_type.System // backward compatible
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) GetTicketsByRefTicketID(ctx context.Context, r *shoptypes.GetTicketsByRefTicketIDRequest) (*shoptypes.GetTicketsByRefTicketIDResponse, error) {
	query := &ticket.ListTicketsByRefTicketIDQuery{
		AccountID:   s.SS.Shop().ID,
		RefTicketID: r.RefTicketID,
	}
	if err := s.TicketQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &shoptypes.GetTicketsByRefTicketIDResponse{
		Tickets: convertpball.Convert_core_Tickets_to_api_Tickets(query.Result),
	}, nil
}

func (s *TicketService) AssignTicket(ctx context.Context, req *api.AssignTicketRequest) (*shoptypes.Ticket, error) {
	getTicketQuery := &ticket.GetTicketByIDQuery{
		ID:        req.TicketID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.TicketQuery.Dispatch(ctx, getTicketQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ticket %v not found", req.TicketID)
	}

	if getTicketQuery.Result.Type == ticket_type.System {
		return nil, cm.Error(cm.PermissionDenied, "", nil)
	}

	cmd := &ticket.AssignTicketCommand{
		ID:              req.TicketID,
		UpdatedBy:       s.SS.User().ID,
		AssignedUserIDs: req.AssignedUserIDs,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) UnassignTicket(ctx context.Context, req *api.AssignTicketRequest) (*shoptypes.Ticket, error) {
	getTicketQuery := &ticket.GetTicketByIDQuery{
		ID:        req.TicketID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.TicketQuery.Dispatch(ctx, getTicketQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ticket %v not found", req.TicketID)
	}

	if getTicketQuery.Result.Type == ticket_type.System {
		return nil, cm.Error(cm.PermissionDenied, "", nil)
	}

	cmd := &ticket.UnassignTicketCommand{
		ID:        req.TicketID,
		UpdatedBy: s.SS.User().ID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) ConfirmTicket(ctx context.Context, req *api.ConfirmTicketRequest) (*shoptypes.Ticket, error) {
	getTicketQuery := &ticket.GetTicketByIDQuery{
		ID:        req.TicketID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.TicketQuery.Dispatch(ctx, getTicketQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ticket %v not found", req.TicketID)
	}

	if getTicketQuery.Result.Type == ticket_type.System {
		return nil, cm.Error(cm.PermissionDenied, "", nil)
	}

	cmd := &ticket.ConfirmTicketCommand{
		ID:        req.TicketID,
		ConfirmBy: s.SS.User().ID,
		Note:      req.Note,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) CloseTicket(ctx context.Context, req *api.CloseTicketRequest) (*shoptypes.Ticket, error) {
	getTicketQuery := &ticket.GetTicketByIDQuery{
		ID:        req.TicketID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.TicketQuery.Dispatch(ctx, getTicketQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ticket %v not found", req.TicketID)
	}

	if getTicketQuery.Result.Type == ticket_type.System {
		return nil, cm.Error(cm.PermissionDenied, "", nil)
	}

	cmd := &ticket.CloseTicketCommand{
		ID:       req.TicketID,
		ClosedBy: s.SS.User().ID,
		Note:     req.Note,
		State:    req.State,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) ReopenTicket(ctx context.Context, req *api.ReopenTicketRequest) (*shoptypes.Ticket, error) {
	getTicketQuery := &ticket.GetTicketByIDQuery{
		ID:        req.TicketID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.TicketQuery.Dispatch(ctx, getTicketQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ticket %v not found", req.TicketID)
	}

	if getTicketQuery.Result.Type == ticket_type.System {
		return nil, cm.Error(cm.PermissionDenied, "", nil)
	}

	cmd := &ticket.ReopenTicketCommand{
		ID:   req.TicketID,
		Note: req.Note,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_Ticket_to_api_Ticket(cmd.Result), nil
}

func (s *TicketService) UpdateTicketRefTicketID(ctx context.Context, req *api.UpdateTicketRefTicketIDRequest) (*pbcm.UpdatedResponse, error) {
	getTicketQuery := &ticket.GetTicketByIDQuery{
		ID:        req.ID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.TicketQuery.Dispatch(ctx, getTicketQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ticket %v not found", req.ID)
	}

	if getTicketQuery.Result.Type == ticket_type.System {
		return nil, cm.Error(cm.PermissionDenied, "", nil)
	}

	cmd := &ticket.UpdateTicketRefTicketIDCommand{
		ID:          req.ID,
		RefTicketID: req.RefTicketID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return cmd.Result, nil
}

func (s *TicketService) CreateTicketComment(ctx context.Context, req *api.CreateTicketCommentRequest) (*shoptypes.TicketComment, error) {
	var imageUrls []string
	if len(req.ImageUrls) > 0 {
		imageUrls = req.ImageUrls
	} else if req.ImageUrl != "" {
		imageUrls = []string{req.ImageUrl}
	}

	cmd := &ticket.CreateTicketCommentCommand{
		CreatedBy:     s.SS.User().ID,
		CreatedName:   s.SS.User().FullName,
		CreatedSource: account_type.Shop,
		TicketID:      req.TicketID,
		AccountID:     s.SS.Shop().ID,
		ParentID:      req.ParentID,
		Message:       req.Message,
		ImageUrls:     imageUrls,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_TicketComment_to_api_TicketComment(cmd.Result), nil
}

func (s *TicketService) UpdateTicketComment(ctx context.Context, req *api.UpdateTicketCommentRequest) (*shoptypes.TicketComment, error) {
	var imageUrls []string
	if len(req.ImageUrls) > 0 {
		imageUrls = req.ImageUrls
	} else if req.ImageUrl != "" {
		imageUrls = []string{req.ImageUrl}
	}

	cmd := &ticket.UpdateTicketCommentCommand{
		AccountID: s.SS.Shop().ID,
		ID:        req.ID,
		UpdatedBy: s.SS.User().ID,
		Message:   req.Message,
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
		DeletedBy: s.SS.User().ID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &api.DeleteTicketCommentResponse{Count: cmd.Result}, nil
}

func (s *TicketService) GetTicketComments(ctx context.Context, req *api.GetTicketCommentsRequest) (*api.GetTicketCommentsResponse, error) {
	var filter = &ticket.FilterGetTicketComment{}
	filter.AccountID = s.SS.Shop().ID
	filter.TicketID = req.TicketID
	if req.Filter != nil {
		filter.CreatedBy = req.Filter.CreatedBy
		filter.ParentID = req.Filter.ParentID
		filter.Title = req.Filter.Title
		filter.IDs = req.Filter.IDs
	}
	paging := cmapi.CMPaging(req.Paging)
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

func (s *TicketService) CreateTicketLabel(ctx context.Context, req *api.CreateTicketLabelRequest) (*shoptypes.TicketLabel, error) {
	cmd := &ticket.CreateTicketLabelCommand{
		ShopID:   s.SS.Shop().ID,
		Type:     ticket_type.Internal,
		Name:     req.Name,
		Code:     req.Code,
		Color:    req.Color,
		ParentID: req.ParentID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_TicketLabel_to_api_TicketLabel(cmd.Result), nil
}

func (s *TicketService) UpdateTicketLabel(ctx context.Context, req *api.UpdateTicketLabelRequest) (*shoptypes.TicketLabel, error) {
	cmd := &ticket.UpdateTicketLabelCommand{
		ID:       req.ID,
		ShopID:   s.SS.Shop().ID,
		Type:     ticket_type.Internal,
		Name:     req.Name,
		Color:    req.Color,
		Code:     req.Code,
		ParentID: req.ParentID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_TicketLabel_to_api_TicketLabel(cmd.Result), nil
}

func (s *TicketService) DeleteTicketLabel(ctx context.Context, req *api.DeleteTicketLabelRequest) (*api.DeleteTicketLabelResponse, error) {
	cmd := &ticket.DeleteTicketLabelCommand{
		ID:          req.ID,
		ShopID:      s.SS.Shop().ID,
		Type:        ticket_type.Internal,
		DeleteChild: req.DeleteChild,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return &api.DeleteTicketLabelResponse{
		Count: 1,
	}, nil
}

func (s *TicketService) GetTicketLabels(ctx context.Context, req *api.GetTicketLabelsRequest) (*api.GetTicketLabelsResponse, error) {
	query := &ticket.ListTicketLabelsQuery{
		ShopID: s.SS.Shop().ID,
		Tree:   req.Tree,
	}
	if req.Filter != nil {
		query.Type = req.Filter.Type
	}

	err := s.TicketQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	return &api.GetTicketLabelsResponse{
		TicketLabels: convertpball.Convert_core_TicketLabels_to_api_TicketLabels(query.Result.TicketLabels),
	}, nil
}
