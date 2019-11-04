package integration

import (
	"context"

	cm "etop.vn/backend/pb/common"
	inte "etop.vn/backend/pb/etop/integration"
)

// +gen:apix

// +apix:path=/integration.Misc
type MiscAPI interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/integration.Integration
type IntegrationAPI interface {
	Init(context.Context, *inte.InitRequest) (*inte.LoginResponse, error)
	// RequestLogin
	//
	// Check if the requested phone or email exists and take corresponding action.
	RequestLogin(context.Context, *inte.RequestLoginRequest) (*inte.RequestLoginResponse, error)

	LoginUsingToken(context.Context, *inte.LoginUsingTokenRequest) (*inte.LoginResponse, error)

	Register(context.Context, *inte.RegisterRequest) (*inte.RegisterResponse, error)

	GrantAccess(context.Context, *inte.GrantAccessRequest) (*inte.GrantAccessResponse, error)

	SessionInfo(context.Context, *cm.Empty) (*inte.LoginResponse, error)
}
