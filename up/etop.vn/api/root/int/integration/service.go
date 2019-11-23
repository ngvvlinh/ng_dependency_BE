package integration

import (
	"context"

	cm "etop.vn/api/pb/common"
	inte "etop.vn/api/pb/etop/integration"
)

// +gen:apix
// +gen:apix:doc-path=etop/integration

// +apix:path=/integration.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/integration.Integration
type IntegrationService interface {
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