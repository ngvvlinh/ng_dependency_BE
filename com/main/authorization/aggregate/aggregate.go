package aggregate

import (
	"context"

	"o.o/api/main/authorization"
	"o.o/backend/com/main/authorization/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

var _ authorization.Aggregate = &AuthorizationAggregate{}

type AuthorizationAggregate struct {
	Auth *auth.Authorizer

	AccountUserStore sqlstore.AccountUserStoreInterface
	ShopStore        sqlstore.ShopStoreInterface
}

func AuthorizationAggregateMessageBus(a *AuthorizationAggregate) authorization.CommandBus {
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
	if err := a.AccountUserStore.UpdateRole(ctx, updateRoleCmd); err != nil {
		return nil, err
	}

	relationship := &authorization.Relationship{
		UserID:    updateRoleCmd.Result.UserID,
		AccountID: updateRoleCmd.Result.AccountID,
		FullName:  updateRoleCmd.Result.FullName,
		ShortName: updateRoleCmd.Result.ShortName,
		Position:  updateRoleCmd.Result.Position,
		Roles:     convert.ConvertStringsToRoles(updateRoleCmd.Result.Permission.Roles),
		Actions:   convert.ConvertStringsToActions(a.Auth.ListActionsByRoles(updateRoleCmd.Result.Permission.Roles)),
	}

	return relationship, nil
}

func (a *AuthorizationAggregate) UpdateRelationship(
	ctx context.Context, args *authorization.UpdateRelationshipArgs,
) (*authorization.Relationship, error) {
	if args.UserID.Int64() == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "user_id ph???i kh??c 0")
	}
	updateRelationshipCmd := &identitymodelx.UpdateInfosCommand{
		AccountID: args.AccountID,
		UserID:    args.UserID,
		FullName:  args.FullName,
		ShortName: args.ShortName,
		Position:  args.Position,
	}
	if err := a.AccountUserStore.UpdateInfos(ctx, updateRelationshipCmd); err != nil {
		return nil, err
	}
	relationship := &authorization.Relationship{
		UserID:    updateRelationshipCmd.Result.UserID,
		AccountID: updateRelationshipCmd.Result.AccountID,
		FullName:  updateRelationshipCmd.Result.FullName,
		ShortName: updateRelationshipCmd.Result.ShortName,
		Position:  updateRelationshipCmd.Result.Position,
		Roles:     convert.ConvertStringsToRoles(updateRelationshipCmd.Result.Permission.Roles),
		Actions:   convert.ConvertStringsToActions(a.Auth.ListActionsByRoles(updateRelationshipCmd.Result.Permission.Roles)),
	}
	return relationship, nil
}

func (a *AuthorizationAggregate) validateAuthorization(ctx context.Context, args *authorization.UpdatePermissionArgs) (identitymodel.AccountUserExtended, error) {
	var currAccountUserHasOwner, currAccountUserHasStaffManagement, accountUserHasOwner, accountUserHasStaffManagement bool
	if len(args.Roles) == 0 {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.InvalidArgument, "roles kh??ng h???p l???", nil)
	}
	if args.CurrUserID == args.UserID {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "Kh??ng ???????c thay ?????i quy???n c???a t??i kho???n b???n ??ang s??? h???u", nil)
	}
	if authorization.IsContainsRole(args.Roles, authorization.RoleShopOwner) {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "Kh??ng ???????c g??n quy???n owner", nil)
	}
	if authorization.IsContainsRole(args.Roles, authorization.RoleAdmin) {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "Kh??ng ???????c g??n quy???n admin", nil)
	}
	getCurrAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		AccountID: args.AccountID,
		UserID:    args.CurrUserID,
	}
	if err := a.AccountUserStore.GetAccountUserExtended(ctx, getCurrAccountUserQuery); err != nil {
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
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "Ch??? c?? quy???n owner ho???c staff_management m???i ???????c th???c hi???n thao t??c n??y", nil)
	}
	getAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		AccountID: args.AccountID,
		UserID:    args.UserID,
	}
	if err := a.AccountUserStore.GetAccountUserExtended(ctx, getAccountUserQuery); err != nil {
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
			return identitymodel.AccountUserExtended{}, cm.Errorf(cm.InvalidArgument, nil, "role %v kh??ng h???p l???", role)
		}
	}
	if !currAccountUserHasOwner && authorization.IsContainsRole(args.Roles, authorization.RoleStaffManagement) {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "b???n kh??ng c?? quy???n owner ????? g??n quy???n staff_management", nil)
	}
	if accountUserHasOwner {
		return identitymodel.AccountUserExtended{}, cm.Error(cm.PermissionDenied, "kh??ng ???????c thay ?????i role c???a owner", nil)
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
		return 0, cm.Error(cm.InvalidArgument, "user_id ph???i kh??c 0", nil)
	}
	if accountID.Int64() == 0 {
		return 0, cm.Error(cm.InvalidArgument, "account_id ph???i kh??c 0", nil)
	}

	getAccountQuery := &identitymodelx.GetShopQuery{
		ShopID: accountID,
	}
	if err := a.ShopStore.GetShop(ctx, getAccountQuery); err != nil {
		return 0, err
	}
	if getAccountQuery.Result.OwnerID == userID {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Kh??ng th??? r???i kh???i shop b???n s??? h???u")
	}

	query := &identitymodelx.GetAccountUserQuery{
		UserID:    userID,
		AccountID: accountID,
	}
	if err := a.AccountUserStore.GetAccountUser(ctx, query); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "t??i kho???n b???n kh??ng thu???c shop n??y").
			Throw()
	}

	if !query.Result.DeletedAt.IsZero() {
		return 0, cm.Error(cm.FailedPrecondition, "t??i kho???n b???n kh??ng thu???c shop n??y	", nil)
	}

	cmd := &identitymodelx.DeleteAccountUserCommand{
		AccountID: accountID,
		UserID:    userID,
	}
	if err := a.AccountUserStore.DeleteAccountUser(ctx, cmd); err != nil {
		return 0, err
	}
	return cmd.Result.Updated, nil
}

func (a *AuthorizationAggregate) RemoveUser(
	ctx context.Context, args *authorization.RemoveUserArgs,
) (update int, _ error) {
	var accountUserHasRoleStaffManagement, currentAccountUserHasRoleOwner, currentAccountUserHasStaffManagement bool
	if args.CurrentUserID == args.UserID {
		return 0, cm.Error(cm.InvalidArgument, "Kh??ng th??? g??? b??? ch??nh m??nh kh???i shop", nil)
	}

	getAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		AccountID: args.AccountID,
		UserID:    args.UserID,
	}
	if err := a.AccountUserStore.GetAccountUserExtended(ctx, getAccountUserQuery); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "t??i kho???n b???n kh??ng thu???c shop n??y").
			Throw()
	}
	accountUser := getAccountUserQuery.Result.AccountUser
	rolesAccountUser := convert.ConvertStringsToRoles(accountUser.Permission.Roles)
	if authorization.IsContainsRole(rolesAccountUser, authorization.RoleShopOwner) {
		return 0, cm.Error(cm.InvalidArgument, "Kh??ng th??? g??? b??? user c?? quy???n owner", nil)
	}
	if authorization.IsContainsRole(rolesAccountUser, authorization.RoleStaffManagement) {
		accountUserHasRoleStaffManagement = true
	}

	getCurrentAccountUserQuery := &identitymodelx.GetAccountUserExtendedQuery{
		AccountID: args.AccountID,
		UserID:    args.CurrentUserID,
	}
	if err := a.AccountUserStore.GetAccountUserExtended(ctx, getCurrentAccountUserQuery); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "t??i kho???n b???n kh??ng thu???c shop n??y").
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
		return 0, cm.Error(cm.FailedPrecondition, "B???n kh??ng c?? quy???n th???c hi???n thao t??c n??y", nil)
	}
	if accountUserHasRoleStaffManagement && !currentAccountUserHasRoleOwner && currentAccountUserHasStaffManagement {
		return 0, cm.Error(cm.FailedPrecondition, "B???n kh??ng c?? quy???n ????? th???c hi???n thao t??c n??y", nil)
	}

	cmd := &identitymodelx.DeleteAccountUserCommand{
		AccountID: args.AccountID,
		UserID:    args.UserID,
	}
	if err := a.AccountUserStore.DeleteAccountUser(ctx, cmd); err != nil {
		return 0, err
	}
	return cmd.Result.Updated, nil
}
