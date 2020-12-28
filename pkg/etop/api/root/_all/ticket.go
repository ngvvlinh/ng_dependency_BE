package _all

import (
	"context"

	"o.o/api/supporting/ticket"
	api "o.o/api/top/int/etop"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/authorize/session"
)

type TicketService struct {
	session.Session

	TicketQuery ticket.QueryBus
}

func (s *TicketService) Clone() api.TicketService {
	res := *s
	return &res
}

func (s *TicketService) GetTicketLabels(ctx context.Context, request *api.GetTicketLabelsRequest) (*api.GetTicketLabelsResponse, error) {
	query := &ticket.ListTicketLabelsQuery{
		Tree: request.Tree,
	}
	err := s.TicketQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	return &api.GetTicketLabelsResponse{
		TicketLabels: convertpball.Convert_core_TicketLabels_to_api_TicketLabels(query.Result.TicketLabels),
	}, nil
}
