package root

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/api/main/invitation"
	api "o.o/api/top/int/etop"
	apietop "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

type UserRelationshipService struct {
	session.Session

	InvitationAggr         invitation.CommandBus
	InvitationQuery        invitation.QueryBus
	AuthorizationAggregate authorization.CommandBus
	ShopStore              sqlstore.ShopStoreInterface
	UserStore              sqlstore.UserStoreInterface
}

func (s *UserRelationshipService) Clone() api.UserRelationshipService {
	res := *s
	return &res
}

func (s *UserRelationshipService) AcceptInvitation(ctx context.Context, q *api.AcceptInvitationRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &invitation.AcceptInvitationCommand{
		UserID: s.SS.Claim().UserID,
		Token:  q.Token,
	}
	if err := s.InvitationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

func (s *UserRelationshipService) RejectInvitation(ctx context.Context, q *api.RejectInvitationRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &invitation.RejectInvitationCommand{
		UserID: s.SS.Claim().UserID,
		Token:  q.Token,
	}
	if err := s.InvitationAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

func (s *UserRelationshipService) GetInvitationByToken(ctx context.Context, q *api.GetInvitationByTokenRequest) (*api.Invitation, error) {
	query := &invitation.GetInvitationByTokenQuery{
		Token: q.Token,
	}
	if err := s.InvitationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbInvitation(query.Result)

	getAccountQuery := &identitymodelx.GetShopQuery{
		ShopID: query.Result.AccountID,
	}
	if err := s.ShopStore.GetShop(ctx, getAccountQuery); err != nil {
		return nil, err
	}
	result.ShopShort = &apietop.ShopShort{
		ID:       getAccountQuery.Result.ID,
		Name:     getAccountQuery.Result.Name,
		Code:     getAccountQuery.Result.Code,
		ImageUrl: getAccountQuery.Result.ImageURL,
	}

	getUserQuery := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Email: query.Result.Email,
		Phone: query.Result.Phone,
	}
	err := s.UserStore.GetUserByEmailOrPhone(ctx, getUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
	// no-op
	case cm.NoError:
		result.UserId = getUserQuery.Result.ID
	default:
		return nil, err
	}

	getInvitedByUserQuery := &identitymodelx.GetUserByIDQuery{
		UserID: query.Result.InvitedBy,
	}
	if err := s.UserStore.GetUserByID(ctx, getInvitedByUserQuery); err != nil {
		return nil, err
	}
	result.InvitedByUser = getInvitedByUserQuery.Result.FullName
	return result, nil
}

func (s *UserRelationshipService) GetInvitations(ctx context.Context, q *api.GetInvitationsRequest) (*api.InvitationsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &invitation.ListInvitationsByEmailAndPhoneQuery{
		Email:   s.SS.User().Email,
		Phone:   s.SS.User().Phone,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.InvitationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.InvitationsResponse{
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
	if err := s.ShopStore.GetShops(ctx, getAccountsQuery); err != nil {
		return nil, err
	}
	mapShop := make(map[dot.ID]*identitymodel.Shop)
	for _, shop := range getAccountsQuery.Result.Shops {
		mapShop[shop.ID] = shop
	}

	for _, invitationEl := range result.Invitations {
		invitationEl.ShopShort = &apietop.ShopShort{
			ID:       invitationEl.ShopId,
			Name:     mapShop[invitationEl.ShopId].Name,
			Code:     mapShop[invitationEl.ShopId].Code,
			ImageUrl: mapShop[invitationEl.ShopId].ImageURL,
		}
	}
	return result, nil
}

func (s *UserRelationshipService) LeaveAccount(ctx context.Context, q *api.UserRelationshipLeaveAccountRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &authorization.LeaveAccountCommand{
		UserID:    s.SS.Claim().UserID,
		AccountID: q.AccountID,
	}
	if err := s.AuthorizationAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}
