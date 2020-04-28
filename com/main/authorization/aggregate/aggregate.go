package aggregate

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/backend/com/main/authorization/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	authservice "o.o/backend/pkg/etop/authorize/auth"
	"o.o/capi/dot"
)

var _ authorization.Aggregate = &AuthorizationAggregate{}

type AuthorizationAggregate struct{}

func NewAuthorizationAggregate() *AuthorizationAggregate {
	return &AuthorizationAggregate{}
}

func (a *AuthorizationAggregate) MessageBus() authorization.CommandBus {
	b := bus.New()
	return authorization.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *AuthorizationAggregate) UpdatePermission(
	ctx context.Context, args *authorization.UpdatePermissionArgs,
) (*authorization.Relationship, error) {
	accountUser, err := a.validateAuthorization(ctx, args)
	if err != nil {
		return nil, err
	}

	accountUser.AccountUser.Permission.Roles = convert.ConvertRolesToStrings(args.Roles)

	updateRoleCmd := &identitymodelx.UpdateRoleCommand{
		AccountID: args.AccountID,
		UserID:    args.UserID,
		Permission: identitymodel.Permission{
			Roles:       accountUser.AccountUser.Roles,
			Permissions: accountUser.AccountUser.Permissions,
		},
	}
	if err := bus.Dispatch(ctx, updateRoleCmd); err != nil {
		return nil, err
	}

	relationship := &authorization.Relationship{
		UserID:    updateRoleCmd.Result.UserID,
		AccountID: updateRoleCmd.Result.AccountID,
		FullName:  updateRoleCmd.Result.FullName,
		ShortName: updateRoleCmd.Result.ShortName,
		Position:  updateRoleCmd.Result.Position,
		Roles:     convert.ConvertStringsToRoles(updateRoleCmd.Result.Permission.Roles),
		Actions:   convert.ConvertStringsToActions(authservice.ListActionsByRoles(updateRoleCmd.Result.Permission.Roles)),
	}

	return relationship, nil
}

func (a *AuthorizationAggregate) UpdateRelationship(
	ctx context.Context, args *authorization.UpdateRelationshipArgs,
) (*authorization.Relationship, error) {
	if args.UserID.Int64() == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "user_id phải khác 0")
	}
	updateRelationshipCmd := &identitymodelx.UpdateInfosCommand{
		AccountID: args.AccountID,
		UserID:    args.UserID,
		FullName:  args.FullName,
		ShortName: args.ShortName,
		Position:  args.Position,
	}
	if err := bus.Dispatch(ctx, updateRelationshipCmd); err != nil {
		return nil, err
	}
	relationship := &authorization.Relationship{
		UserID:    updateRelationshipCmd.Result.UserID,
		AccountID: updateRelationshipCmd.Result.AccountID,
		FullName:  updateRelationshipCmd.Result.FullName,
		ShortName: updateRelationshipCmd.Result.ShortName,
		Position:  updateRelationshipCmd.Result.Position,
		Roles:     convert.ConvertStringsToRoles(updateRelationshipCmd.Result.Permission.Roles),
		Actions:   convert.ConvertStringsToActions(authservice.ListActionsByRoles(updateRelationshipCmd.Result.Permission.Roles)),
	}
	return relationship, nil
}

func (a *AuthorizationAggregate) validateAuthorization(ctx context.Context, args *authorization.UpdatePermissionArgs) (identitymodel.AccountUserExtended, error) {
	var currAccountUserHasOwner, currAccountUserHasStaffManagement, accountUserHasOwner, accountUserHasStaffManagement bool
	if len(args.Roles) == 0 {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.InvalidArgument, "roles không hợp lệ", nil)
	}
	if args.CurrUserID == args.UserID {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "Không được thay đổi quyền của tài khoản bạn đang sở hữu", nil)
	}
	if authorization.IsContainsRole(args.Roles, authorization.RoleShopOwner) {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "Không được gán quyền owner", nil)
	}
	if authorization.IsContainsRole(args.Roles, authorization.RoleAdmin) {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "Không được gán quyền admin", nil)
	}
	getCurrAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		AccountID: args.AccountID,
		UserID:    args.CurrUserID,
	}
	if err := bus.Dispatch(ctx, getCurrAccountUserQuery); err != nil {
		return identitymodel.AccountUserExtended{}, err
	}
	currAccountUser := getCurrAccountUserQuery.Result
	rolesCurrAccountUser := convert.ConvertStringsToRoles(currAccountUser.AccountUser.Permission.Roles)
	if authorization.IsContainsRole(rolesCurrAccountUser, authorization.RoleShopOwner) {
		currAccountUserHasOwner = true
	}
	if authorization.IsContainsRole(rolesCurrAccountUser, authorization.RoleStaffManagement) {
		currAccountUserHasStaffManagement = true
	}
	if !currAccountUserHasOwner && !currAccountUserHasStaffManagement {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "Chỉ có quyền owner hoặc staff_management mới được thực hiện thao tác này", nil)
	}
	getAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		AccountID: args.AccountID,
		UserID:    args.UserID,
	}
	if err := bus.Dispatch(ctx, getAccountUserQuery); err != nil {
		return identitymodel.AccountUserExtended{}, err
	}
	accountUser := getAccountUserQuery.Result
	rolesAccountUser := convert.ConvertStringsToRoles(accountUser.AccountUser.Permission.Roles)
	if authorization.IsContainsRole(rolesAccountUser, authorization.RoleShopOwner) {
		accountUserHasOwner = true
	}
	if authorization.IsContainsRole(rolesAccountUser, authorization.RoleStaffManagement) {
		accountUserHasStaffManagement = true
	}
	for _, role := range rolesAccountUser {
		if !authorization.IsRole(role) {
			return identitymodel.AccountUserExtended{}, cm.Errorf(cm.InvalidArgument, nil, "role %v không hợp lệ", role)
		}
	}
	if !currAccountUserHasOwner && authorization.IsContainsRole(args.Roles, authorization.RoleStaffManagement) {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "bạn không có quyền owner để gán quyền staff_management", nil)
	}
	if accountUserHasOwner {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "không được thay đổi role của owner", nil)
	}
	if !currAccountUserHasOwner {
		if accountUserHasStaffManagement {
			return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "", nil)
		} else {
			if !currAccountUserHasStaffManagement {
				return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "", nil)
			}
		}
	}
	return accountUser, nil
}

func (a *AuthorizationAggregate) LeaveAccount(
	ctx context.Context, userID, accountID dot.ID,
) (updated int, _ error) {
	if userID.Int64() == 0 {
		return 0, cm.Error(cm.InvalidArgument, "user_id phải khác 0", nil)
	}
	if accountID.Int64() == 0 {
		return 0, cm.Error(cm.InvalidArgument, "account_id phải khác 0", nil)
	}

	getAccountQuery := &identitymodelx.GetShopQuery{
		ShopID: accountID,
	}
	if err := bus.Dispatch(ctx, getAccountQuery); err != nil {
		return 0, err
	}
	if getAccountQuery.Result.OwnerID == userID {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Không thể rời khỏi shop bạn sở hữu")
	}

	query := &identitymodelx.GetAccountUserQuery{
		UserID:    userID,
		AccountID: accountID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "tài khoản bạn không thuộc shop này").
			Throw()
	}

	if !query.Result.DeletedAt.IsZero() {
		return 0, cm.Error(cm.FailedPrecondition, "tài khoản bạn không thuộc shop này	", nil)
	}

	cmd := &identitymodelx.DeleteAccountUserCommand{
		AccountID: accountID,
		UserID:    userID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return 0, err
	}
	return cmd.Result.Updated, nil
}

func (a *AuthorizationAggregate) RemoveUser(
	ctx context.Context, args *authorization.RemoveUserArgs,
) (update int, _ error) {
	var accountUserHasRoleStaffManagement, currentAccountUserHasRoleOwner, currentAccountUserHasStaffManagement bool
	if args.CurrentUserID == args.UserID {
		return 0, cm.Error(cm.InvalidArgument, "Không thể gỡ bỏ chính mình khỏi shop", nil)
	}

	getAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		AccountID: args.AccountID,
		UserID:    args.UserID,
	}
	if err := bus.Dispatch(ctx, getAccountUserQuery); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "tài khoản bạn không thuộc shop này").
			Throw()
	}
	accountUser := getAccountUserQuery.Result.AccountUser
	rolesAccountUser := convert.ConvertStringsToRoles(accountUser.Permission.Roles)
	if authorization.IsContainsRole(rolesAccountUser, authorization.RoleShopOwner) {
		return 0, cm.Error(cm.InvalidArgument, "Không thể gỡ bở user có quyền owner", nil)
	}
	if authorization.IsContainsRole(rolesAccountUser, authorization.RoleStaffManagement) {
		accountUserHasRoleStaffManagement = true
	}

	getCurrentAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		AccountID: args.AccountID,
		UserID:    args.CurrentUserID,
	}
	if err := bus.Dispatch(ctx, getCurrentAccountUserQuery); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "tài khoản bạn không thuộc shop này").
			Throw()
	}
	currentAccountUser := getCurrentAccountUserQuery.Result.AccountUser
	rolesCurrentAccountUser := convert.ConvertStringsToRoles(currentAccountUser.Permission.Roles)
	if authorization.IsContainsRole(rolesCurrentAccountUser, authorization.RoleShopOwner) {
		currentAccountUserHasRoleOwner = true
	}
	if authorization.IsContainsRole(rolesCurrentAccountUser, authorization.RoleStaffManagement) {
		currentAccountUserHasStaffManagement = true
	}

	if !currentAccountUserHasRoleOwner && !currentAccountUserHasStaffManagement {
		return 0, cm.Error(cm.FailedPrecondition, "Bạn không có quyền thực hiện thao tác này", nil)
	}
	if accountUserHasRoleStaffManagement && !currentAccountUserHasRoleOwner && currentAccountUserHasStaffManagement {
		return 0, cm.Error(cm.FailedPrecondition, "Bạn không có quyền để thực hiện thao tác này", nil)
	}

	cmd := &identitymodelx.DeleteAccountUserCommand{
		AccountID: args.AccountID,
		UserID:    args.UserID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return 0, err
	}
	return cmd.Result.Updated, nil
}
