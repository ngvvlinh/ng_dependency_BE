package account_user

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	etopapi "o.o/api/top/int/etop"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/shop_user_role"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type AccountUserService struct {
	session.Session

	IdentityAggr  identity.CommandBus
	IdentityQuery identity.QueryBus
}

func (s *AccountUserService) Clone() api.AccountUserService {
	res := *s
	return &res
}

func (s *AccountUserService) CreateAccountUser(ctx context.Context, r *api.CreateAccountUserRequest) (*etopapi.User, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if !hasPermision(s.SS.GetRoles(), r.Roles) {
		return nil, cm.ErrPermissionDenied
	}

	phoneNorm, ok := validate.NormalizePhone(r.Phone)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}

	phone := phoneNorm.String()
	// register user if needed
	cmd := &identity.RegisterSimplifyCommand{
		Phone:            phone,
		FullName:         r.FullName,
		Password:         r.Password,
		IsUpdatePassword: true,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	// get user
	query := &identity.GetUserByPhoneOrEmailQuery{
		Phone: phone,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	userID := query.Result.ID

	if err := s.createAccountUser(ctx, userID, r.GetAccountUserRoles()); err != nil {
		return nil, err
	}

	queryAccountUser := &identity.GetAccountUserQuery{
		UserID:    userID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.IdentityQuery.Dispatch(ctx, queryAccountUser); err != nil {
		return nil, err
	}

	res := convertpb.Convert_core_User_To_api_User(query.Result)
	return res, nil
}

func (s *AccountUserService) GetAccountUsers(ctx context.Context, r *api.GetAccountUsersRequest) (*api.GetAccountUsersResponse, error) {
	// Parse Paging
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}

	// Get Extended Account Users
	listExtendedAccountUsersQuery := &identity.ListExtendedAccountUsersQuery{
		AccountID: s.SS.Shop().ID,
		Paging:    *paging,
	}
	if r.Filter != nil {
		listExtendedAccountUsersQuery.FullNameNorm = r.Filter.Name
		listExtendedAccountUsersQuery.PhoneNorm = r.Filter.Phone
		listExtendedAccountUsersQuery.ExtensionNumberNorm = r.Filter.ExtensionNumber
		listExtendedAccountUsersQuery.Roles = r.Filter.Roles
		listExtendedAccountUsersQuery.UserIDs = r.Filter.UserIDs
		listExtendedAccountUsersQuery.HasExtension = r.Filter.HasExtension
	}
	if err = s.IdentityQuery.Dispatch(ctx, listExtendedAccountUsersQuery); err != nil {
		return nil, err
	}
	extendedAccountUsers := listExtendedAccountUsersQuery.Result.AccountUsers

	result := &api.GetAccountUsersResponse{
		AccountUsers: convertpb.Convert_core_ExtendedAccountUsers_To_api_ExtendedAccountUsers(extendedAccountUsers),
		Paging:       cmapi.PbCursorPageInfo(paging, &listExtendedAccountUsersQuery.Result.Paging),
	}
	return result, nil
}

func (s *AccountUserService) createAccountUser(ctx context.Context, userID dot.ID, roles []string) error {
	accountID := s.SS.Shop().ID
	query := &identity.GetAccountUserQuery{
		UserID:    userID,
		AccountID: accountID,
	}
	err := s.IdentityQuery.Dispatch(ctx, query)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		// create new
		cmd := &identity.CreateAccountUserCommand{
			AccountID: accountID,
			UserID:    userID,
			Status:    status3.P,
			Permission: identity.Permission{
				Roles: roles,
			},
		}
		return s.IdentityAggr.Dispatch(ctx, cmd)

	case cm.NoError:
		// user is a staff in shop already
		return cm.Errorf(cm.FailedPrecondition, nil, "Người dùng đã tồn tại.")

	default:
		return err
	}
}

func (s *AccountUserService) UpdateAccountUser(ctx context.Context, r *api.UpdateAccountUserRequest) (*pbcm.UpdatedResponse, error) {
	if r.UserID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "user_id không được để trống")
	}
	query := &identity.GetAccountUserQuery{
		UserID:    r.UserID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	accountUser := query.Result
	roles := accountUser.Permission.GetShopUserRoles()
	if len(r.Roles) > 0 {
		for _, _role := range r.Roles {
			if !shop_user_role.ContainsUserRoles(roles, _role) {
				roles = append(roles, _role)
			}
		}
	}
	if !hasPermision(s.SS.GetRoles(), roles) {
		return nil, cm.ErrPermissionDenied
	}
	updated := 0
	if r.Password != "" {
		updateUserPass := &identity.UpdateUserPasswordCommand{
			UserID:   r.UserID,
			Password: r.Password,
		}
		if err := s.IdentityAggr.Dispatch(ctx, updateUserPass); err != nil {
			return nil, err
		}
		updated++
	}
	if len(r.Roles) > 0 {
		// update roles
		update := &identity.UpdateAccountUserPermissionCommand{
			AccountID: s.SS.Shop().ID,
			UserID:    r.UserID,
			Permission: identity.Permission{
				Roles: shop_user_role.ToRolesString(r.Roles),
			},
		}
		if err := s.IdentityAggr.Dispatch(ctx, update); err != nil {
			return nil, err
		}
		updated++
	}

	// update other info
	if r.FullName != "" {
		updateUserInfo := &identity.UpdateUserInfoCommand{
			AccountID: s.SS.Shop().ID,
			UserID:    r.UserID,
			FullName:  r.FullName,
		}
		if err := s.IdentityAggr.Dispatch(ctx, updateUserInfo); err != nil {
			return nil, err
		}
		updated++
	}

	if updated > 1 {
		updated = 1
	}
	return &pbcm.UpdatedResponse{
		Updated: updated,
	}, nil
}

func (s *AccountUserService) DeleteAccountUser(ctx context.Context, r *api.DeleteAccountUserRequest) (*pbcm.DeletedResponse, error) {
	if r.UserID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "user_id không được để trống")
	}
	query := &identity.GetAccountUserQuery{
		UserID:    r.UserID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	accountUser := query.Result
	roles := accountUser.Permission.GetShopUserRoles()
	if !hasPermision(s.SS.GetRoles(), roles) {
		return nil, cm.ErrPermissionDenied
	}

	cmd := &identity.DeleteAccountUsersCommand{
		AccountID: s.SS.Shop().ID,
		UserID:    r.UserID,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.DeletedResponse{Deleted: 1}, nil
}

//Kiểm tra xem người dùng hiện tại (currentRole) có đủ quyền
//để tương tác với các người dùng có quyền khác không (roles)
func hasPermision(currentUserRoles []string, roles []shop_user_role.UserRole) bool {
	// Không người dùng nào được quyền tương tác với người dùng có quyền owner
	if shop_user_role.ContainsUserRoles(roles, shop_user_role.Owner) {
		return false
	}

	// Người dùng có quyền owner được tương tác với tất cả người dùng có quyền khác
	if cm.StringsContain(currentUserRoles, authorization.RoleShopOwner.String()) {
		return true
	}

	// Chỉ người dùng có quyền owner và staff management mới có quyền tương tác với các người dùng có quyền khác
	if !cm.StringsContain(currentUserRoles, authorization.RoleStaffManagement.String()) {
		return false
	}

	var currentRole shop_user_role.UserRole
	// lấy quyền lớn nhất của user hiện tại
	// user_role number càng nhỏ => quyền càng lớn (trừ số 0)
	for _, crole := range currentUserRoles {
		_crole, ok := shop_user_role.ParseUserRole(crole)
		if !ok || _crole == shop_user_role.Unknown {
			continue
		}

		if currentRole == shop_user_role.Unknown {
			currentRole = _crole
			continue
		}
		if currentRole >= _crole {
			currentRole = _crole
		}
	}

	if currentRole == shop_user_role.Unknown {
		return false
	}
	for _, role := range roles {
		if role <= currentRole {
			return false
		}
	}
	return true
}
