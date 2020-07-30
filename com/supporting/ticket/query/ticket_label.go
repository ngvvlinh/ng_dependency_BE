package query

import (
	"context"

	"o.o/api/supporting/ticket"
	"o.o/backend/com/supporting/ticket/aggregate"
)

func (q TicketQuery) GetTicketLabelByID(ctx context.Context, args *ticket.GetTicketLabelByIDArgs) (*ticket.TicketLabel, error) {
	return q.TicketLabelStore(ctx).ID(args.ID).GetTicketLabel()
}

func (q TicketQuery) ListTicketLabels(ctx context.Context, args *ticket.GetTicketLabelsArgs) (*ticket.GetTicketLabelsResponse, error) {
	ticketLabels, err := q.TicketLabelStore(ctx).ListTicketLabels()
	if err != nil {
		return nil, err
	}
	if args.Tree {
		ticketLabels = aggregate.MakeTreeLabel(ticketLabels)
	}
	return &ticket.GetTicketLabelsResponse{TicketLabels: ticketLabels}, nil
}
