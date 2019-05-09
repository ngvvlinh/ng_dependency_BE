package relationship

import "etop.vn/backend/pkg/etop/model"

type InviterInfo struct {
	InviterUserID      int64
	InviterFullName    string
	InviterAccountName string
	InviterAccountType model.AccountType
}

type InvitationInfo struct {
	ShortName  string
	FullName   string
	Position   string
	Permission model.Permission
}

type InviteUserToAccountCommand struct {
	InviterInfo *InviterInfo
	Invitation  *InvitationInfo

	AccountID    int64
	EmailOrPhone string

	Result struct {
		AccountUser *model.AccountUser
	}
}

type GetOrCreateUserStubCommand struct {
	InviterInfo *InviterInfo
	*model.UserInner

	Result struct {
		User         *model.User
		UserInternal *model.UserInternal
	}
}

type GetOrCreateInvitationCommand struct {
	AccountID int64
	UserID    int64

	InviterInfo *InviterInfo
	Invitation  *InvitationInfo

	Result struct {
		AccountUser *model.AccountUser
	}
}
