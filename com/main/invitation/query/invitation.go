package query

import (
	"context"

	"o.o/api/main/invitation"
	"o.o/api/shopping"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/invitation/aggregate"
	"o.o/backend/com/main/invitation/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

var _ invitation.QueryService = &InvitationQuery{}

type InvitationQuery struct {
	db          *cmsql.Database
	store       sqlstore.InvitationStoreFactory
	flagNewLink aggregate.FlagEnableNewLinkInvitation
}

func NewInvitationQuery(
	db com.MainDB,
	flagNewLink aggregate.FlagEnableNewLinkInvitation,
) *InvitationQuery {
	return &InvitationQuery{
		db:          db,
		store:       sqlstore.NewInvitationStore(db),
		flagNewLink: flagNewLink,
	}
}

func InvitationQueryMessageBus(q *InvitationQuery) invitation.QueryBus {
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
	URL, err := aggregate.GetInvitationURL(ctx, invitation, q.flagNewLink)
	if err != nil {
		return nil, err
	}
	invitation.InvitationURL = URL.String()
	return invitation, nil
}

func (q *InvitationQuery) GetInvitationByToken(
	ctx context.Context, token string,
) (*invitation.Invitation, error) {
	invitation, err := q.store(ctx).Token(token).NotExpires().GetInvitation()
	if err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy lời mời").
			Throw()
	}
	URL, err := aggregate.GetInvitationURL(ctx, invitation, q.flagNewLink)
	if err != nil {
		return nil, err
	}
	invitation.InvitationURL = URL.String()
	return invitation, nil
}

func (q *InvitationQuery) ListInvitationsByEmailAndPhone(
	ctx context.Context, args *invitation.ListInvitationsByEmailAndPhoneArgs,
) (*invitation.InvitationsResponse, error) {
	query := q.store(ctx).PhoneOrEmail(args.Phone, args.Email).Filters(args.Filters)
	invitations, err := query.WithPaging(args.Paging).ListInvitations()
	if err != nil {
		return nil, err
	}
	invitations, err = PopulateInvitationURL(ctx, invitations, q.flagNewLink)
	if err != nil {
		return nil, err
	}
	return &invitation.InvitationsResponse{
		Invitations: invitations,
	}, nil
}

func (q *InvitationQuery) ListInvitations(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*invitation.InvitationsResponse, error) {
	query := q.store(ctx).AccountID(args.ShopID).Filters(args.Filters)
	invitations, err := query.WithPaging(args.Paging).ListInvitations()
	if err != nil {
		return nil, err
	}
	invitations, err = PopulateInvitationURL(ctx, invitations, q.flagNewLink)
	if err != nil {
		return nil, err
	}
	return &invitation.InvitationsResponse{
		Invitations: invitations,
	}, nil
}

func (q *InvitationQuery) ListInvitationsAcceptedByEmail(
	ctx context.Context, email string,
) (*invitation.InvitationsResponse, error) {
	emailNorm, ok := validate.NormalizeEmail(email)
	if !ok {
		return nil, cm.Error(cm.InvalidArgument, "Email không hợp lệ", nil)
	}
	query := q.store(ctx).Email(emailNorm.String()).Accepted()
	invitations, err := query.ListInvitations()
	if err != nil {
		return nil, err
	}
	invitations, err = PopulateInvitationURL(ctx, invitations, q.flagNewLink)
	if err != nil {
		return nil, err
	}
	return &invitation.InvitationsResponse{
		Invitations: invitations,
	}, nil
}

func PopulateInvitationURL(ctx context.Context, invitations []*invitation.Invitation, flag aggregate.FlagEnableNewLinkInvitation) ([]*invitation.Invitation, error) {
	for k, v := range invitations {
		URL, err := aggregate.GetInvitationURL(ctx, v, flag)
		invitations[k].InvitationURL = URL.String()
		if err != nil {
			return nil, err
		}
	}
	return invitations, nil
}
