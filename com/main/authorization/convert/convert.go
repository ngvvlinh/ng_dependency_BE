package convert

import (
	"o.o/api/main/authorization"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/etop/authorize/auth"
)

func ConvertStringsToRoles(args []string) (roles []authorization.Role) {
	for _, arg := range args {
		roles = append(roles, authorization.Role(arg))
	}
	return
}

func ConvertRolesToStrings(roles []authorization.Role) (outs []string) {
	for _, role := range roles {
		outs = append(outs, role.String())
	}
	return
}

func ConvertStringsToActions(args []string) (outs []authorization.Action) {
	for _, action := range args {
		outs = append(outs, authorization.Action(action))
	}
	return
}

func ConvertAccountUserExtendedToAuthorization(aa *auth.Authorizer, accountUser *identitymodel.AccountUserExtended) *authorization.Authorization {
	user := accountUser.User
	roles := accountUser.AccountUser.Permission.Roles
	au := &authorization.Authorization{
		UserID:    user.ID,
		FullName:  accountUser.AccountUser.FullName,
		ShortName: accountUser.AccountUser.ShortName,
		Position:  accountUser.AccountUser.Position,
		Email:     user.Email,
		Roles:     ConvertStringsToRoles(roles),
		Actions:   ConvertStringsToActions(aa.ListActionsByRoles(roles)),
	}
	return au
}

func ConvertAccountUserToRelationship(aa *auth.Authorizer, accountUser *identitymodel.AccountUser) *authorization.Relationship {
	var isDeleted bool
	if !accountUser.DeletedAt.IsZero() {
		isDeleted = true
	}
	roles := ConvertStringsToRoles(accountUser.Permission.Roles)
	actions := ConvertStringsToActions(aa.ListActionsByRoles(accountUser.Permission.Roles))
	return &authorization.Relationship{
		UserID:    accountUser.UserID,
		AccountID: accountUser.AccountID,
		FullName:  accountUser.FullName,
		ShortName: accountUser.ShortName,
		Position:  accountUser.Position,
		Roles:     roles,
		Actions:   actions,
		Deleted:   isDeleted,
	}
}
