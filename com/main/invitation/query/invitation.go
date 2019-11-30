package query

import (
	"context"

	"etop.vn/api/main/invitation"
	"etop.vn/api/shopping"
	"etop.vn/backend/com/main/invitation/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/capi/dot"
)

var _ invitation.QueryService = &InvitationQuery{}

type InvitationQuery struct {
	db    *cmsql.Database
	store sqlstore.InvitationStoreFactory
}

func NewInvitationQuery(
	database *cmsql.Database,
) *InvitationQuery {
	return &InvitationQuery{
		db:    database,
		store: sqlstore.NewInvitationStore(database),
	}
}

func (q *InvitationQuery) MessageBus() invitation.QueryBus {
	b := bus.New()
	return invitation.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *InvitationQuery) GetInvitation(
	ctx context.Context, ID dot.ID,
) (*invitation.Invitation, error) {
	invitation, err := q.store(ctx).ID(ID).GetInvitation()
	if err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy lời mời").
			Throw()
	}

	return invitation, nil
}

func (q *InvitationQuery) GetInvitationByToken(
	ctx context.Context, token string,
) (*invitation.Invitation, error) {
	invitation, err := q.store(ctx).Token(token).GetInvitation()
	if err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy lời mời").
			Throw()
	}

	return invitation, nil
}

func (q *InvitationQuery) ListInvitationsByEmail(
	ctx context.Context, args *invitation.ListInvitationsByEmailArgs,
) (*invitation.InvitationsResponse, error) {
	query := q.store(ctx).Email(args.Email).Filters(args.Filters)
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	invitations, err := query.Paging(args.Paging).ListInvitations()
	if err != nil {
		return nil, err
	}
	return &invitation.InvitationsResponse{
		Invitations: invitations,
		Count:       count,
	}, nil
}

func (q *InvitationQuery) ListInvitations(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*invitation.InvitationsResponse, error) {
	query := q.store(ctx).AccountID(args.ShopID).Filters(args.Filters)
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	invitations, err := query.Paging(args.Paging).ListInvitations()
	if err != nil {
		return nil, err
	}
	return &invitation.InvitationsResponse{
		Invitations: invitations,
		Count:       count,
	}, nil
}

func (q *InvitationQuery) ListInvitationsAcceptedByEmail(
	ctx context.Context, email string,
) (*invitation.InvitationsResponse, error) {
	query := q.store(ctx).Email(email).Accepted()
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	invitations, err := query.ListInvitations()
	if err != nil {
		return nil, err
	}
	return &invitation.InvitationsResponse{
		Invitations: invitations,
		Count:       count,
	}, nil
}
