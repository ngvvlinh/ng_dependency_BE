package sadmin

import (
	"context"

	cm "etop.vn/backend/pb/common"
	"etop.vn/backend/pb/etop"
	sadmin "etop.vn/backend/pb/etop/sadmin"
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
