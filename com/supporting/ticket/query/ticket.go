package query

import (
	"context"

	"o.o/api/supporting/ticket"
	com "o.o/backend/com/main"
	"o.o/backend/com/supporting/ticket/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/capi"
)

type TicketQuery struct {
	TicketStore        sqlstore.TicketStoreFactory
	TicketLabelStore   sqlstore.TicketLabelStoreFactory
	TicketCommentStore sqlstore.TicketCommentStoreFactory
	EventBus           capi.EventBus
	RedisStore         redis.Store
}

var _ ticket.QueryService = &TicketQuery{}

func NewTicketQuery(redisStore redis.Store, eventBus capi.EventBus, db com.MainDB) *TicketQuery {
	return &TicketQuery{
		TicketStore:        sqlstore.NewTicketStore(db),
		TicketLabelStore:   sqlstore.NewTicketLabelStore(db),
		TicketCommentStore: sqlstore.NewTicketCommentStore(db),
		EventBus:           eventBus,
		RedisStore:         redisStore,
	}
}

func TicketQueryMessageBus(q *TicketQuery) ticket.QueryBus {
	b := bus.New()
	return ticket.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q TicketQuery) GetTicketByID(ctx context.Context, args *ticket.GetTicketByIDArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	ticketCore, err := q.TicketStore(ctx).ID(args.ID).AccountID(args.AccountID).GetTicket()
	if err != nil {
		return nil, err
	}
	return ticketCore, nil
}

func (q TicketQuery) ListTickets(ctx context.Context, args *ticket.GetTicketsArgs) (*ticket.GetTicketsResponse, error) {
	query := q.TicketStore(ctx)
	if args.Filter != nil {
		// optional
		if args.Filter.Code != "" {
			query = query.Code(args.Filter.Code)
		}
		if args.Filter.AccountID != 0 {
			query = query.AccountID(args.Filter.AccountID)
		}
		if args.Filter.CreatedBy != 0 {
			query = query.CreatedBy(args.Filter.CreatedBy)
		}
		if args.Filter.IDs != nil {
			query = query.IDs(args.Filter.IDs...)
		}
		if args.Filter.LabelIDs != nil {
			query = query.LabelIDs(args.Filter.LabelIDs)
		}
		if args.Filter.Title != "" {
			query = query.TitleFullTextSearch(args.Filter.Title)
		}
		if args.Filter.RefType != 0 {
			query = query.RefType(args.Filter.RefType)
		}
		if args.Filter.RefID != 0 {
			query = query.RefID(args.Filter.RefID)
		}
		if args.Filter.AssignedUserIDs != nil && len(args.Filter.AssignedUserIDs) > 0 {
			query = query.AssignedUserIDs(args.Filter.AssignedUserIDs)
		}
		if args.Filter.State != 0 {
			query = query.State(args.Filter.State)
		}
		if args.Filter.RefCode != "" {
			query = query.State(args.Filter.State)
		}
	}
	result, err := query.WithPaging(args.Paging).ListTickets()
	if err != nil {
		return nil, err
	}
	return &ticket.GetTicketsResponse{
		Tickets: result,
		Paging:  query.GetPaging(),
	}, nil
}