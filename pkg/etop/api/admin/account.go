package admin

import (
	"context"
	"fmt"

	"o.o/api/main/authorization"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	"o.o/api/top/types/etc/status3"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/admin/convert"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

const EtopAccountId = 101

type AccountService struct {
	session.Session

	AccountAuthStore sqlstore.AccountAuthStoreFactory
	UserStore        sqlstore.UserStoreFactory
	PartnerStore     sqlstore.PartnerStoreInterface
	AccountUserStore sqlstore.AccountUserStoreInterface
}

func (s *AccountService) Clone() admin.AccountService {
	res := *s
	return &res
}

func (s *AccountService) CreatePartner(ctx context.Context, q *admin.CreatePartnerRequest) (*etop.Partner, error) {
	cmd := &identitymodelx.CreatePartnerCommand{
		Partner: convert.CreatePartnerRequestToModel(q),
	}
	if err := s.PartnerStore.CreatePartner(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbPartner(cmd.Result.Partner)
	return result, nil
}

func (s *AccountService) GenerateAPIKey(ctx context.Context, q *admin.GenerateAPIKeyRequest) (*admin.GenerateAPIKeyResponse, error) {
	_, err := s.AccountAuthStore(ctx).AccountID(q.AccountId).Get()
	if cm.ErrorCode(err) != cm.NotFound {
		return nil, cm.MapError(err).
			Map(cm.OK, cm.AlreadyExists, "account already has an api_key").
			Throw()
	}

	aa := &identitymodel.AccountAuth{
		AccountID:   q.AccountId,
		Status:      status3.P,
		Roles:       nil,
		Permissions: nil,
	}
	err = s.AccountAuthStore(ctx).Create(aa)
	result := &admin.GenerateAPIKeyResponse{
		AccountId: q.AccountId,
		ApiKey:    aa.AuthKey,
	}
	return result, err
}

func (s *AccountService) CreateAdminUser(ctx context.Context, q *admin.CreateAdminUserRequest) (*admin.CreateAdminUserResponse, error) {
	for _, role := range q.Roles {
		if !authorization.IsContainsRole(authorization.InternalRoles, authorization.Role(role)) {
			return nil, cm.Errorf(cm.InvalidArgument, nil, fmt.Sprintf("Role không hợp lệ: %v", role))
		}
	}

	query := &identitymodelx.GetUserByEmailOrPhoneQuery{
		Email: q.Email,
	}
	if err := s.UserStore(ctx).GetUserByEmailOrPhone(ctx, query); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return nil, cm.Errorf(cm.NotFound, nil, "Email không tồn tại trong hệ thống.")
		}
		return nil, err
	}

	user := query.Result
	getAccountUserQuery := &identitymodelx.GetAccountUserQuery{
		AccountID:       EtopAccountId,
		UserID:          user.ID,
		FindByAccountID: false,
	}
	err := s.AccountUserStore.GetAccountUser(ctx, getAccountUserQuery)
	if err == nil {
		// this case mean `account_user` with `user_id` and `account_id` already exists.
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Tài khoản này đã là admin của etop")
	}

	// Oke let create
	accountUser := &identitymodel.AccountUser{
		UserID:    user.ID,
		AccountID: EtopAccountId,
		Status:    status3.P, // Enable
		Permission: identitymodel.Permission{
			Roles: q.Roles,
		},
	}
	createAccountUserCmd := &identitymodelx.CreateAccountUserCommand{
		AccountUser: accountUser,
	}
	if err := s.AccountUserStore.CreateAccountUser(ctx, createAccountUserCmd); err != nil {
		return nil, err
	}

	createdAccountUser := createAccountUserCmd.Result
	return &admin.CreateAdminUserResponse{
		UserId: createdAccountUser.UserID,
		Roles:  createdAccountUser.Roles,
		Status: createdAccountUser.Status,
	}, nil
}

func (s *AccountService) GetAPIKey(ctx context.Context, q *admin.GetAPIKeyRequest) (*admin.GetAPIKeyResponse, error) {
	accountAuth, err := s.AccountAuthStore(ctx).AccountID(q.AccountID).Get()
	if err != nil {
		return nil, err
	}
	return &admin.GetAPIKeyResponse{
		ApiKey: accountAuth.AuthKey,
	}, nil
}

func (s *AccountService) UpdateAdminUser(ctx context.Context, q *admin.UpdateAdminUserRequest) (*admin.UpdateAdminUserResponse, error) {
	for _, role := range q.Roles {
		if !authorization.IsContainsRole(authorization.InternalRoles, authorization.Role(role)) {
			return nil, cm.Errorf(cm.InvalidArgument, nil, fmt.Sprintf("invalid role %v", role))
		}
	}

	getAccountUserQuery := &identitymodelx.GetAccountUserQuery{
		AccountID:       EtopAccountId,
		UserID:          q.UserId,
		FindByAccountID: false,
	}
	if err := s.AccountUserStore.GetAccountUser(ctx, getAccountUserQuery); err != nil {
		return nil, err
	}

	accountUser := &identitymodel.AccountUser{
		UserID:    q.UserId,
		AccountID: EtopAccountId,
		Status:    q.Status,
	}

	if q.Status != status3.P && q.Status != status3.N {
		accountUser.Status = getAccountUserQuery.Result.Status
	}

	if len(q.Roles) == 0 {
		accountUser.Permission = identitymodel.Permission{
			Roles: getAccountUserQuery.Result.Permission.Roles,
		}
	} else {
		accountUser.Permission = identitymodel.Permission{
			Roles: q.Roles,
		}
	}

	updateInternalAccountCmd := &identitymodelx.UpdateAccountUserCommand{
		AccountUser: accountUser,
	}
	if err := s.AccountUserStore.UpdateAccountUser(ctx, updateInternalAccountCmd); err != nil {
		return nil, err
	}

	updatedUserAccount := updateInternalAccountCmd.Result
	return &admin.UpdateAdminUserResponse{
		UserId: updatedUserAccount.UserID,
		Roles:  updatedUserAccount.Roles,
		Status: updatedUserAccount.Status,
	}, nil
}

func (s *AccountService) GetAdminUsers(ctx context.Context, req *admin.GetAdminUsersRequest) (*admin.GetAdminUserResponse, error) {
	getAdminAccQuery := &identitymodelx.GetAccountUserExtendedsQuery{
		AccountIDs: []dot.ID{EtopAccountId},
		Roles:      req.Filter.Roles,
	}

	if err := s.AccountUserStore.GetAccountUserExtendeds(ctx, getAdminAccQuery); err != nil {
		return nil, err
	}

	res := &admin.GetAdminUserResponse{}
	for _, v := range getAdminAccQuery.Result.AccountUsers {
		res.Admins = append(res.Admins, &admin.AdminAccountResponse{
			UserId:    v.User.ID,
			FullName:  v.User.FullName,
			Email:     v.User.Email,
			Phone:     v.User.Phone,
			Roles:     v.AccountUser.Roles,
			CreatedAt: v.AccountUser.CreatedAt,
			UpdatedAt: v.AccountUser.UpdatedAt,
		})
	}
	return res, nil
}

func (s *AccountService) DeleteAdminUser(ctx context.Context, req *admin.DeleteAdminUserRequest) (*admin.DeleteAdminUserResponse, error) {
	deleteAccCmd := &identitymodelx.DeleteAccountUserCommand{
		AccountID: EtopAccountId,
		UserID:    req.UserID,
	}
	if err := s.AccountUserStore.DeleteAccountUser(ctx, deleteAccCmd); err != nil {
		return nil, err
	}
	return &admin.DeleteAdminUserResponse{Updated: deleteAccCmd.Result.Updated}, nil
}
