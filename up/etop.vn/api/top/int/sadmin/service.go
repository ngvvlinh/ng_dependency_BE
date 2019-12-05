package sadmin

import (
	"context"

	etop "etop.vn/api/top/int/etop"
	cm "etop.vn/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=etop/sadmin

// +apix:path=/sadmin.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/sadmin.User
type UserService interface {
	CreateUser(context.Context, *SAdminCreateUserRequest) (*etop.RegisterResponse, error)
	ResetPassword(context.Context, *SAdminResetPasswordRequest) (*cm.Empty, error)
	LoginAsAccount(context.Context, *LoginAsAccountRequest) (*etop.LoginResponse, error)
}
