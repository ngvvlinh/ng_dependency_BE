package api

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/api/main/invitation"
	apietop "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/capi/dot"
)

type UserRelationshipService struct {
	InvitationAggr         invitation.CommandBus
	InvitationQuery        invitation.QueryBus
	AuthorizationAggregate authorization.CommandBus
}

func (s *UserRelationshipService) Clone() *UserRelationshipService {
	res := *s
	return &res
}

func (s *UserRelationshipService) AcceptInvitation(ctx context.Context, q *UserRelationshipAcceptInvitationEndpoint) error {
	cmd := &invitation.AcceptInvitationCommand{
		UserID: q.Context.UserID,
		Token:  q.Token,
	}
	if err := s.InvitationAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *UserRelationshipService) RejectInvitation(ctx context.Context, q *UserRelationshipRejectInvitationEndpoint) error {
	cmd := &invitation.RejectInvitationCommand{
		UserID: q.Context.UserID,
		Token:  q.Token,
	}
	if err := s.InvitationAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *UserRelationshipService) GetInvitationByToken(ctx context.Context, q *UserRelationshipGetInvitationByTokenEndpoint) error {
	query := &invitation.GetInvitationByTokenQuery{
		Token: q.Token,
	}
	if err := s.InvitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbInvitation(query.Result)

	getAccountQuery := &identitymodelx.GetShopQuery{
		ShopID: query.Result.AccountID,
	}
	if err := bus.Dispatch(ctx, getAccountQuery); err != nil {
		return err
	}
	q.Result.ShopShort = &apietop.ShopShort{
		ID:       getAccountQuery.Result.ID,
		Name:     getAccountQuery.Result.Name,
		Code:     getAccountQuery.Result.Code,
		ImageUrl: getAccountQuery.Result.ImageURL,
	}

	getUserQuery := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Email: query.Result.Email,
		Phone: query.Result.Phone,
	}
	err := bus.Dispatch(ctx, getUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
	// no-op
	case cm.NoError:
		q.Result.UserId = getUserQuery.Result.ID
	default:
		return err
	}

	getInvitedByUserQuery := &identitymodelx.GetUserByIDQuery{
		UserID: query.Result.InvitedBy,
	}
	if err := bus.Dispatch(ctx, getInvitedByUserQuery); err != nil {
		return err
	}
	q.Result.InvitedByUser = getInvitedByUserQuery.Result.FullName

	return nil
}

func (s *UserRelationshipService) GetInvitations(ctx context.Context, q *UserRelationshipGetInvitationsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &invitation.ListInvitationsByEmailAndPhoneQuery{
		Email:   q.Context.User.Email,
		Phone:   q.Context.User.Phone,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.InvitationQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apietop.InvitationsResponse{
		Invitations: convertpb.PbInvitations(query.Result.Invitations),
		Paging:      cmapi.PbPageInfo(paging),
	}

	var accountIDs []dot.ID
	var hasAccountID bool
	for _, invitationEl := range query.Result.Invitations {
		hasAccountID = false
		for _, accountID := range accountIDs {
			if accountID == invitationEl.AccountID {
				hasAccountID = true
			}
		}
		if !hasAccountID {
			accountIDs = append(accountIDs, invitationEl.AccountID)
		}
	}

	getAccountsQuery := &identitymodelx.GetShopsQuery{
		ShopIDs: accountIDs,
	}
	if err := bus.Dispatch(ctx, getAccountsQuery); err != nil {
		return err
	}
	mapShop := make(map[dot.ID]*identitymodel.Shop)
	for _, shop := range getAccountsQuery.Result.Shops {
		mapShop[shop.ID] = shop
	}

	for _, invitationEl := range q.Result.Invitations {
		invitationEl.ShopShort = &apietop.ShopShort{
			ID:       invitationEl.ShopId,
			Name:     mapShop[invitationEl.ShopId].Name,
			Code:     mapShop[invitationEl.ShopId].Code,
			ImageUrl: mapShop[invitationEl.ShopId].ImageURL,
		}
	}

	return nil
}

func (s *UserRelationshipService) LeaveAccount(ctx context.Context, q *UserRelationshipLeaveAccountEndpoint) error {
	cmd := &authorization.LeaveAccountCommand{
		UserID:    q.Context.UserID,
		AccountID: q.AccountID,
	}
	if err := s.AuthorizationAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}
