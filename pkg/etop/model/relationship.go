package model

import (
	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

type GetAccountUserQuery struct {
	UserID          dot.ID
	AccountID       dot.ID
	FindByAccountID bool

	Result *AccountUser
}

type GetAccountUserExtendedQuery struct {
	UserID    dot.ID
	AccountID dot.ID

	Result AccountUserExtended
}

type GetAccountUserExtendedsQuery struct {
	AccountIDs []dot.ID

	Paging  *cm.Paging
	Filters []cm.Filter
	Status  dot.NullInt

	Result struct {
		AccountUsers []*AccountUserExtended

		Total int
	}
}

type GetAccountRolesQuery = GetAccountUserExtendedQuery

type CreateAccountUserCommand struct {
	AccountUser *AccountUser

	Result *AccountUser
}

type UpdateAccountUserCommand struct {
	AccountUser *AccountUser
}

type AccountPermission struct {
	Account    `sq:"inline"`
	Permission `sq:"inline"`
}

type RemoveUserFromAccount struct {
	AccountID dot.ID
	UserID    dot.ID

	Result bool
}

type UpdateRoleCommand struct {
	AccountID dot.ID
	UserID    dot.ID
	Permission

	Result *AccountUser
}

type GetAllAccountRolesQuery struct {
	UserID dot.ID
	Type   AccountType

	Result []*AccountUserExtended
}

type GetAllAccountUsersQuery struct {
	UserIDs []dot.ID
	Type    AccountType

	Result []*AccountUser
}
