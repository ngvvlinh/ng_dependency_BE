package query

import (
	"context"

	"o.o/api/supporting/ticket"
	"o.o/api/top/types/etc/ticket/ticket_type"
	"o.o/backend/com/supporting/ticket/aggregate"
	cm "o.o/backend/pkg/common"
)

func (q TicketQuery) GetTicketLabelByID(ctx context.Context, args *ticket.GetTicketLabelByIDArgs) (*ticket.TicketLabel, error) {
	return q.TicketLabelStore(ctx).ID(args.ID).GetTicketLabel()
}

func (q TicketQuery) ListTicketLabels(ctx context.Context, args *ticket.GetTicketLabelsArgs) (*ticket.GetTicketLabelsResponse, error) {
	if args.Type.Valid && args.Type.Enum == ticket_type.Internal && args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "internal: missing shop_id")
	}
	if !args.Type.Valid && args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing shop_id")
	}

	query := q.TicketLabelStore(ctx)
	if args.Type.Valid {
		if args.Type.Enum == ticket_type.Internal {
			query = query.Type(ticket_type.Internal).ShopID(args.ShopID)
		}
		if args.Type.Enum == ticket_type.System {
			query = query.Type(ticket_type.System)
		}
	} else {
		// get all internal && system
		query = query.InternalAndSystem(args.ShopID)
	}
	ticketLabels, err := query.ListTicketLabels()
	if err != nil {
		return nil, err
	}

	if args.Tree {
		ticketLabels = aggregate.MakeTreeLabel(ticketLabels)
	}
	return &ticket.GetTicketLabelsResponse{TicketLabels: ticketLabels}, nil
}
