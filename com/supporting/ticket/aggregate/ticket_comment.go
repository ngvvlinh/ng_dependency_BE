package aggregate

import (
	"context"

	"o.o/api/supporting/ticket"
	"o.o/backend/com/supporting/ticket/model"
	cm "o.o/backend/pkg/common"
)

func (a TicketAggregate) CreateTicketComment(ctx context.Context, args *ticket.CreateTicketCommentArgs) (*ticket.TicketComment, error) {
	if args.TicketID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ticket_id")
	}
	ticketCore, err := a.TicketStore(ctx).ID(args.TicketID).AccountID(args.AccountID).GetTicket()
	if err != nil {
		return nil, err
	}
	// kiểm tra xem là admin, nếu là admin thì kiểm tra assign
	if !args.IsLeader && args.IsAdmin {
		isAssigned := false
		for _, v := range ticketCore.AssignedUserIDs {
			if v == args.CreatedBy {
				isAssigned = true
				break
			}
		}
		if !isAssigned {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Bạn chưa được assign vào ticket này.")
		}
	}
	if args.ParentID != 0 {
		_, err := a.TicketCommentStore(ctx).TicketID(args.TicketID).AccountID(args.AccountID).ID(args.ParentID).GetTicketComment()
		if err != nil {
			return nil, err
		}
	}
	var ticketComment = &ticket.TicketComment{}
	err = scheme.Convert(args, ticketComment)
	if err != nil {
		return nil, err
	}
	if err := a.TicketCommentStore(ctx).Create(ticketComment); err != nil {
		return nil, err
	}

	return a.TicketCommentStore(ctx).ID(ticketComment.ID).GetTicketComment()
}

func (a TicketAggregate) UpdateTicketComment(ctx context.Context, args *ticket.UpdateTicketCommentArgs) (*ticket.TicketComment, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing account_id")
	}
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing id")
	}
	ticketComment, err := a.TicketCommentStore(ctx).ID(args.ID).AccountID(args.AccountID).GetTicketComment()
	if err != nil {
		return nil, err
	}
	if ticketComment.CreatedBy != args.UpdatedBy {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Chỉ có người tạo comment mới được chỉnh sửa")
	}
	var ticketModel = &model.TicketComment{
		Message: args.Message,
	}
	if len(args.ImageUrls) > 0 {
		ticketModel.ImageUrls = args.ImageUrls
	}
	err = a.TicketCommentStore(ctx).ID(args.ID).UpdateTicketCommentDB(ticketModel)
	if err != nil {
		return nil, err
	}

	return a.TicketCommentStore(ctx).ID(args.ID).AccountID(args.AccountID).GetTicketComment()
}

func (a TicketAggregate) DeleteTicketComment(ctx context.Context, args *ticket.DeleteTicketCommentArgs) (int, error) {
	if args.AccountID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing account_id")
	}
	if args.ID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing id")
	}

	// admin can delete every comments into ticket
	// shop can delete only comments were created by shop
	q := a.TicketCommentStore(ctx).ID(args.ID).AccountID(args.AccountID)
	if !args.IsAdmin {
		q = q.CreatedBy(args.DeletedBy)
	}
	return q.SoftDelete(args.DeletedBy)
}
