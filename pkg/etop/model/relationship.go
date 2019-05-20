package model

import cm "etop.vn/backend/pkg/common"

type GetAccountUserQuery struct {
	UserID          int64
	AccountID       int64
	FindByAccountID bool

	Result *AccountUser
}

type GetAccountUserExtendedQuery struct {
	UserID    int64
	AccountID int64

	Result AccountUserExtended
}

type GetAccountUserExtendedsQuery struct {
	AccountIDs []int64

	Paging  *cm.Paging
	Filters []cm.Filter
	Status  *int

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
	AccountID int64
	UserID    int64

	Result bool
}

type UpdateRoleCommand struct {
	AccountID int64
	UserID    int64
	Permission

	Result *AccountUser
}

type GetAllAccountRolesQuery struct {
	UserID int64
	Type   AccountType

	Result []*AccountUserExtended
}

type GetAllAccountUsersQuery struct {
	UserIDs []int64
	Type    AccountType

	Result []*AccountUser
}
