package relationship

import (
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

type InviterInfo struct {
	InviterUserID      dot.ID
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

	AccountID    dot.ID
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
	AccountID dot.ID
	UserID    dot.ID

	InviterInfo *InviterInfo
	Invitation  *InvitationInfo

	Result struct {
		AccountUser *model.AccountUser
	}
}
