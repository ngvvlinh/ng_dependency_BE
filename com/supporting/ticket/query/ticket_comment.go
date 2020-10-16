package query

import (
	"context"

	"o.o/api/supporting/ticket"
	cm "o.o/backend/pkg/common"
)

func (q TicketQuery) GetTicketCommentByID(ctx context.Context, args *ticket.GetTicketCommentByIDArgs) (*ticket.TicketComment, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing account_id")
	}
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing id")
	}
	return q.TicketCommentStore(ctx).ID(args.ID).AccountID(args.AccountID).GetTicketComment()
}

func (q TicketQuery) ListTicketComments(ctx context.Context, args *ticket.GetTicketCommentsArgs) (*ticket.GetTicketCommentsResponse, error) {
	var query = q.TicketCommentStore(ctx).WithPaging(args.Paging)
	if args.Filter != nil {
		if args.Filter.ParentID != 0 {
			query = query.ParentID(args.Filter.ParentID)
		}
		if args.Filter.IDs != nil && len(args.Filter.IDs) > 0 {
			query = query.IDs(args.Filter.IDs...)
		}
		if args.Filter.CreatedBy != 0 {
			query = query.CreatedBy(args.Filter.CreatedBy)
		}
		if args.Filter.AccountID != 0 {
			query = query.AccountID(args.Filter.AccountID)
		}
		if args.Filter.TicketID != 0 {
			query = query.TicketID(args.Filter.TicketID)
		}
	}
	result, err := query.ListTicketComments()
	if err != nil {
		return nil, err
	}
	return &ticket.GetTicketCommentsResponse{
		TicketComments: result,
		Paging:         query.GetPaging(),
	}, nil

}
