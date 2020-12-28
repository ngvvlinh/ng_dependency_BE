package ticket

import (
	"context"

	"o.o/api/supporting/ticket"
	api "o.o/api/top/int/shop"
	shoptypes "o.o/api/top/int/shop/types"
	"o.o/api/top/types/etc/account_type"
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

func (s *TicketService) CreateTicketComment(ctx context.Context, request *api.CreateTicketCommentRequest) (*shoptypes.TicketComment, error) {
	var imageUrls []string
	if len(request.ImageUrls) > 0 {
		imageUrls = request.ImageUrls
	} else if request.ImageUrl != "" {
		imageUrls = []string{request.ImageUrl}
	}

	cmd := &ticket.CreateTicketCommentCommand{
		CreatedBy:     s.SS.User().ID,
		CreatedName:   s.SS.User().FullName,
		CreatedSource: account_type.Shop,
		TicketID:      request.TicketID,
		AccountID:     s.SS.Shop().ID,
		ParentID:      request.ParentID,
		Message:       request.Message,
		ImageUrls:     imageUrls,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpball.Convert_core_TicketComment_to_api_TicketComment(cmd.Result), nil
}

func (s *TicketService) UpdateTicketComment(ctx context.Context, request *api.UpdateTicketCommentRequest) (*shoptypes.TicketComment, error) {
	var imageUrls []string
	if len(request.ImageUrls) > 0 {
		imageUrls = request.ImageUrls
	} else if request.ImageUrl != "" {
		imageUrls = []string{request.ImageUrl}
	}

	cmd := &ticket.UpdateTicketCommentCommand{
		AccountID: s.SS.Shop().ID,
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
		DeletedBy: s.SS.User().ID,
	}
	if err := s.TicketAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &api.DeleteTicketCommentResponse{Count: cmd.Result}, nil
}

func (s *TicketService) GetTicketComments(ctx context.Context, request *api.GetTicketCommentsRequest) (*api.GetTicketCommentsResponse, error) {
	var filter = &ticket.FilterGetTicketComment{}
	filter.AccountID = s.SS.Shop().ID
	filter.TicketID = request.TicketID
	if request.Filter != nil {
		filter.CreatedBy = request.Filter.CreatedBy
		filter.ParentID = request.Filter.ParentID
		filter.Title = request.Filter.Title
		filter.IDs = request.Filter.IDs
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
		filter.Code = request.Filter.Code
		filter.RefID = request.Filter.RefID
		filter.RefType = request.Filter.RefType
		filter.State = request.Filter.State
		filter.RefCode = request.Filter.RefCode
	}
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
		AssignedUserIDs: nil,
		AccountID:       shopID,
		LabelIDs:        request.LabelIDs,
		Title:           request.Title,
		Description:     request.Description,
		Note:            request.Note,
		RefID:           request.RefID,
		RefType:         request.RefType,
		RefCode:         request.RefCode,
		Source:          request.Source,
		CreatedBy:       userID,
		CreatedName:     s.SS.User().FullName,
		CreatedSource:   account_type.Shop,
		Result:          nil,
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
