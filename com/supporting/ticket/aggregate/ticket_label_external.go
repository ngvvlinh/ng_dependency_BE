package aggregate

import (
	"context"

	"o.o/api/supporting/ticket"
	"o.o/backend/com/supporting/ticket/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
)

func (a TicketAggregate) CreateTicketLabelExternal(
	ctx context.Context, args *ticket.CreateTicketLabelExternalArgs,
) (*ticket.TicketLabelExternal, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}

	ticketLabelExternal := &ticket.TicketLabelExternal{}
	if err := scheme.Convert(args, ticketLabelExternal); err != nil {
		return nil, err
	}

	ticketLabelExternal.ID = cm.NewID()
	if err := a.TicketLabelExternalStore(ctx).Create(ticketLabelExternal); err != nil {
		return nil, err
	}

	return ticketLabelExternal, nil
}

func (a TicketAggregate) UpdateTicketLabelExternal(
	ctx context.Context, args *ticket.UpdateTicketLabelExternalArgs,
) (*ticket.TicketLabelExternal, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}

	ticketLabelExternalModel := &model.TicketLabelExternal{
		ExternalName: args.ExternalName,
	}
	if err := a.TicketLabelExternalStore(ctx).ID(args.ID).UpdateTicketLabelExternalDB(ticketLabelExternalModel); err != nil {
		return nil, err
	}

	return a.TicketLabelExternalStore(ctx).ID(args.ID).GetTicketLabelExternal()
}

func (a TicketAggregate) DeleteTicketLabelExternal(
	ctx context.Context, args *ticket.DeleteTicketLabelExternalArgs,
) (int, error) {
	if args.ID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}

	deleted := 0
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var _err error
		deleted, _err = a.TicketLabelExternalStore(ctx).ID(args.ID).SoftDelete()
		if _err != nil {
			return _err
		}

		// delete all ticket_label_ticket_label_external by ticket_label_external_id
		if _, err := a.TicketLabelTicketLabelExternalStore(ctx).TicketLabelExternalID(args.ID).SoftDelete(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return 0, err
	}
	return deleted, nil
}
