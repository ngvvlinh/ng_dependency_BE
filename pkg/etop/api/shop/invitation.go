package shop

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/invitation"
	pbshop "etop.vn/api/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/api/convertpb"
)

func init() {
	bus.AddHandlers("api",
		invitationService.CreateInvitation,
		invitationService.GetInvitations)
}

func (s *InvitationService) CreateInvitation(ctx context.Context, q *CreateInvitationEndpoint) error {
	var roles []invitation.Role
	for _, role := range q.Roles {
		roles = append(roles, invitation.Role(role))
	}
	cmd := &invitation.CreateInvitationCommand{
		AccountID: q.Context.Shop.ID,
		Email:     q.Email,
		Roles:     roles,
		Status:    etop.S3Zero,
		InvitedBy: q.Context.UserID,
	}
	if err := invitationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = convertpb.PbInvitation(cmd.Result)
	return nil
}

func (s *InvitationService) GetInvitations(ctx context.Context, q *GetInvitationsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &invitation.ListInvitationsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := invitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.InvitationsResponse{
		Invitations: convertpb.PbInvitations(query.Result.Invitations),
		Paging:      cmapi.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}
