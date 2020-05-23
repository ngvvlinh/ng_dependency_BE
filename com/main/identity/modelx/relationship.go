package modelx

import (
	"o.o/api/main/authorization"
	"o.o/api/top/types/etc/account_type"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

type GetAccountUserQuery struct {
	UserID          dot.ID
	AccountID       dot.ID
	FindByAccountID bool

	Result *identitymodel.AccountUser
}

type GetAccountUserExtendedQuery struct {
	UserID    dot.ID
	AccountID dot.ID

	Result identitymodel.AccountUserExtended
}

type GetAccountUserExtendedsQuery struct {
	AccountIDs []dot.ID

	Paging  *cm.Paging
	Filters []cm.Filter
	Status  dot.NullInt

	IncludeDeleted bool

	Result struct {
		AccountUsers []*identitymodel.AccountUserExtended
	}
}

type GetAccountRolesQuery = GetAccountUserExtendedQuery

type CreateAccountUserCommand struct {
	AccountUser *identitymodel.AccountUser

	Result *identitymodel.AccountUser
}

type UpdateAccountUserCommand struct {
	AccountUser *identitymodel.AccountUser
}

type AccountPermission struct {
	identitymodel.Account    `sq:"inline"`
	identitymodel.Permission `sq:"inline"`
}

type RemoveUserFromAccount struct {
	AccountID dot.ID
	UserID    dot.ID

	Result bool
}

type UpdateRoleCommand struct {
	AccountID dot.ID
	UserID    dot.ID
	identitymodel.Permission

	Result *identitymodel.AccountUser
}

type UpdateInfosCommand struct {
	AccountID dot.ID
	UserID    dot.ID
	FullName  dot.NullString
	ShortName dot.NullString
	Position  dot.NullString

	Result *identitymodel.AccountUser
}

type DeleteAccountUserCommand struct {
	AccountID dot.ID
	UserID    dot.ID

	Result struct {
		Updated int
	}
}

type GetAllAccountRolesQuery struct {
	UserID dot.ID
	Type   account_type.NullAccountType

	Result []*identitymodel.AccountUserExtended
}

type GetAllAccountUsersQuery struct {
	UserIDs []dot.ID
	Type    account_type.NullAccountType
	Role    authorization.Role

	Result []*identitymodel.AccountUser
}
