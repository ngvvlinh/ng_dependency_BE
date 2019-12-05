package api

import (
	"context"
	"fmt"
	"time"

	"etop.vn/api/top/int/etop"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/logic/relationship"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/l"
)

func init() {
	bus.AddHandlers("api",
		relationshipService.AnswerInvitation,
		relationshipService.GetUsersInCurrentAccounts,
		relationshipService.InviteUserToAccount,
		relationshipService.LeaveAccount,
		relationshipService.RemoveUserFromCurrentAccount,
	)
}

type RelationshipService struct{}

var relationshipService = &RelationshipService{}

func (s *RelationshipService) AnswerInvitation(ctx context.Context, r *AnswerInvitationEndpoint) error {
	resp, err := s.answerInvitation(ctx, r)
	if err != nil {
		return err
	}
	r.Result = resp.Result
	return nil
}

func (s *RelationshipService) answerInvitation(ctx context.Context, r *AnswerInvitationEndpoint) (*AnswerInvitationEndpoint, error) {
	if r.AccountId == 0 {
		return r, cm.Error(cm.InvalidArgument, "Missing Name", nil)
	}
	if r.Response == nil {
		return r, cm.Error(cm.InvalidArgument, "Invalid response", nil)
	}
	response := *convertpb.Status3ToModel(r.Response)

	userID := r.Context.UserID
	accountID := r.AccountId
	accUserQuery := &model.GetAccountUserExtendedQuery{
		UserID:    userID,
		AccountID: accountID,
	}
	if err := bus.Dispatch(ctx, accUserQuery); err != nil {
		return r, err
	}

	updateAccUser := &model.AccountUser{
		UserID:    userID,
		AccountID: accountID,
	}

	accUser := accUserQuery.Result.AccountUser
	switch accUser.Status {
	case model.S3Zero:
		switch response {
		case model.S3Positive, model.S3Negative:
			updateAccUser.Status = response
			updateAccUser.ResponseStatus = response
		default:
			return r, cm.Error(cm.InvalidArgument, "Invalid response", nil)
		}

	case model.S3Positive, model.S3Negative:
		// If the response is the same as the status, just respond it
		if response == accUser.Status {
			r.Result = convertpb.PbUserAccount(&accUserQuery.Result)
			return r, nil
		}

		// Positive response for negative status, can not accept
		if response > accUser.Status {
			return r, cm.Error(cm.FailedPrecondition, "Bạn không thể tham gia vào tài khoản này.", nil).
				Log("positive response for negative status", l.Int("response", int(response)), l.Int("status", int(accUser.Status)))
		}

		// Negative response for positive status, set both to negative
		updateAccUser.Status = response
		updateAccUser.ResponseStatus = response

	default:
		return r, cm.Error(cm.FailedPrecondition, "Bạn không thể tham gia vào tài khoản này.", nil).
			Log("unexpected status")
	}

	updateCmd := &model.UpdateAccountUserCommand{
		AccountUser: updateAccUser,
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return r, err
	}

	// Get it again
	if err := bus.Dispatch(ctx, accUserQuery); err != nil {
		return r, cm.Error(cm.Internal, "", err).
			Log("unexpected")
	}
	r.Result = convertpb.PbUserAccount(&accUserQuery.Result)
	return r, nil
}

func (s *RelationshipService) GetUsersInCurrentAccounts(ctx context.Context, r *GetUsersInCurrentAccountsEndpoint) error {
	accountIDs, err := MixAccount(r.Context.Claim, r.Mixed)
	if err != nil {
		return err
	}

	paging := cmapi.CMPaging(r.Paging)
	query := &model.GetAccountUserExtendedsQuery{
		AccountIDs: accountIDs,
		Paging:     paging,
		Filters:    cmapi.ToFilters(r.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = &etop.ProtectedUsersResponse{
		Paging: cmapi.PbPageInfo(paging, query.Result.Total),
		Users:  convertpb.PbUserAccounts(query.Result.AccountUsers),
	}
	return nil
}

func (s *RelationshipService) InviteUserToAccount(ctx context.Context, r *InviteUserToAccountEndpoint) error {
	key := fmt.Sprintf("InviteUserToAccount %v-%v", r.Context.User.ID, r.InviteeIdentifier)
	resp, err := idempgroup.DoAndWrap(key, 10*time.Second, func() (interface{}, error) {
		return s.inviteUserToAccount(ctx, r)
	}, "thêm người dùng")

	if err != nil {
		return err
	}
	r.Result = resp.(*InviteUserToAccountEndpoint).Result
	return nil
}

func (s *RelationshipService) inviteUserToAccount(ctx context.Context, r *InviteUserToAccountEndpoint) (*InviteUserToAccountEndpoint, error) {

	inviter := r.Context.User.User
	accountQuery := &model.GetAccountRolesQuery{
		UserID:    inviter.ID,
		AccountID: r.AccountId,
	}
	if err := bus.Dispatch(ctx, accountQuery); err != nil {
		return r, err
	}
	account := accountQuery.Result.Account

	// The user must be the owner of the given account
	if account.OwnerID != inviter.ID {
		return r, cm.Errorf(cm.PermissionDenied, nil, "Bạn không có quyền thực hiện thao tác này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.")
	}

	inviteCmd := &relationship.InviteUserToAccountCommand{
		InviterInfo: &relationship.InviterInfo{
			InviterUserID:      inviter.ID,
			InviterFullName:    inviter.FullName,
			InviterAccountName: account.Name,
			InviterAccountType: account.Type,
		},
		Invitation: &relationship.InvitationInfo{
			ShortName:  "",
			FullName:   "",
			Position:   "",
			Permission: model.Permission{},
		},
		AccountID:    r.AccountId,
		EmailOrPhone: r.InviteeIdentifier,
	}
	if err := bus.Dispatch(ctx, inviteCmd); err != nil {
		return r, err
	}
	accUser := inviteCmd.Result.AccountUser
	r.Result = convertpb.PbUserAccountIncomplete(accUser, account)
	return r, nil
}

func (s *RelationshipService) LeaveAccount(ctx context.Context, r *LeaveAccountEndpoint) error {
	return nil
}

func (s *RelationshipService) RemoveUserFromCurrentAccount(ctx context.Context, r *RemoveUserFromCurrentAccountEndpoint) error {
	return nil
}
