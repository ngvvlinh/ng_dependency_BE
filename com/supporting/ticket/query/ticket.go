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
	"o.o/common/l"
)

var ll = l.New()

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

func (q *TicketQuery) GetTicketByID(ctx context.Context, args *ticket.GetTicketByIDArgs) (*ticket.Ticket, error) {
	query := q.TicketStore(ctx)
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	if args.CreatedBy != 0 || len(args.AssignedUserIDs) > 0 {
		query = query.AssignedUserIDsOrCreatedBy(args.CreatedBy, args.AssignedUserIDs)
	}
	ticketCore, err := query.ID(args.ID).AccountID(args.AccountID).GetTicket()
	if err != nil {
		return nil, err
	}
	return ticketCore, nil
}

func (q *TicketQuery) GetTicketByExternalID(ctx context.Context, args *ticket.GetTicketByExternalIDArgs) (*ticket.Ticket, error) {
	if args.ExternalID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ExternalID")
	}
	ticketCore, err := q.TicketStore(ctx).ExternalID(args.ExternalID).GetTicket()
	if err != nil {
		return nil, err
	}
	return ticketCore, nil
}

func (q *TicketQuery) ListTickets(ctx context.Context, args *ticket.GetTicketsArgs) (*ticket.ListTicketsResponse, error) {
	query := q.TicketStore(ctx)
	if args.Filter != nil {
		// optional
		if args.Filter.Code != "" {
			query = query.Code(args.Filter.Code)
		}
		if args.Filter.AccountID != 0 {
			query = query.AccountID(args.Filter.AccountID)
		}
		if args.Filter.ConfirmedBy != 0 {
			query = query.ConfirmedBy(args.Filter.ConfirmedBy)
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
		ll.Info("args.Filter.State", l.Object("args.Filter.State", args.Filter.State))
		if args.Filter.State != 0 {
			query = query.State(args.Filter.State)
		}
		if args.Filter.RefCode != "" {
			query = query.State(args.Filter.State)
		}
		if len(args.Filter.Types) != 0 {
			query = query.Types(args.Filter.Types)
		}

		if args.IsLeader {
			// Nếu là chủ shop
			if len(args.Filter.AssignedUserIDs) > 0 {
				query = query.AssignedUserIDs(args.Filter.AssignedUserIDs)
			}
			if args.Filter.CreatedBy != 0 {
				query = query.CreatedBy(args.Filter.CreatedBy)
			}
		} else {
			// Nếu là nhân viên:
			//	+ Không filter theo field assigned_user_ids:
			//		* Ticket: created_by OR assigned_user_ids
			//  + Có filter theo field assigned_user_ids:
			//		* Ticket: created_by AND assigned_user_ids
			if !args.HasFilter {
				if len(args.Filter.AssignedUserIDs) > 0 || args.Filter.CreatedBy != 0 {
					query = query.AssignedUserIDsOrCreatedBy(args.Filter.CreatedBy, args.Filter.AssignedUserIDs)
				}
			} else {
				if len(args.Filter.AssignedUserIDs) > 0 || args.Filter.CreatedBy != 0 {
					query = query.AssignedUserIDsAndCreatedBy(args.Filter.CreatedBy, args.Filter.AssignedUserIDs)
				}
			}
		}

	}
	tickets, err := query.WithPaging(args.Paging).ListTickets()
	if err != nil {
		return nil, err
	}

	return &ticket.ListTicketsResponse{
		Tickets: tickets,
		Paging:  query.GetPaging(),
	}, nil
}

func (q *TicketQuery) ListTicketsByRefTicketID(ctx context.Context, args *ticket.ListTicketsByRefTicketIDArgs) ([]*ticket.Ticket, error) {
	query := q.TicketStore(ctx)
	if args.RefTicketID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin ref ticket ID")
	}
	if len(args.AssignedUserIDs) > 0 || args.CreatedBy != 0 {
		query = query.AssignedUserIDsOrCreatedBy(args.CreatedBy, args.AssignedUserIDs)
	}
	return query.RefTicketID(args.RefTicketID).OptionalAccountID(args.AccountID).ListTickets()
}
