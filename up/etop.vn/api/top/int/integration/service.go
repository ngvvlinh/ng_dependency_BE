package integration

import (
	"context"

	cm "etop.vn/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=etop/integration

// +apix:path=/integration.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/integration.Integration
type IntegrationService interface {
	Init(context.Context, *InitRequest) (*LoginResponse, error)
	// RequestLogin
	//
	// Check if the requested phone or email exists and take corresponding action.
	RequestLogin(context.Context, *RequestLoginRequest) (*RequestLoginResponse, error)

	LoginUsingToken(context.Context, *LoginUsingTokenRequest) (*LoginResponse, error)

	// LoginUsingTokenWL
	//
	// Do all stuff to grant access shop of whitelabel partner
	LoginUsingTokenWL(context.Context, *cm.Empty) (*LoginResponse, error)

	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)

	GrantAccess(context.Context, *GrantAccessRequest) (*GrantAccessResponse, error)

	SessionInfo(context.Context, *cm.Empty) (*LoginResponse, error)
}
