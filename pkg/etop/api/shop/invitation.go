package shop

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/invitation"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
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

	q.Result = pbshop.PbInvitation(cmd.Result)
	return nil
}

func (s *InvitationService) GetInvitations(ctx context.Context, q *GetInvitationsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &invitation.ListInvitationsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := invitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.InvitationsResponse{
		Invitations: pbshop.PbInvitations(query.Result.Invitations),
		Paging:      pbcm.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}
