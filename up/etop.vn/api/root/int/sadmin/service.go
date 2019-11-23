package sadmin

import (
	"context"

	cm "etop.vn/api/pb/common"
	"etop.vn/api/pb/etop"
	sadmin "etop.vn/api/pb/etop/sadmin"
)

// +gen:apix
// +gen:apix:doc-path=etop/sadmin

// +apix:path=/sadmin.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/sadmin.User
type UserService interface {
	CreateUser(context.Context, *sadmin.SAdminCreateUserRequest) (*etop.RegisterResponse, error)
	ResetPassword(context.Context, *sadmin.SAdminResetPasswordRequest) (*cm.Empty, error)
	LoginAsAccount(context.Context, *sadmin.LoginAsAccountRequest) (*etop.LoginResponse, error)
}
