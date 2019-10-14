package etop

import (
	"context"

	cm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
)

// +gen:apix

// +apix:path=/etop.User
type UserAPI interface {
	// Register
	//
	// Register a new user or after a user has login using generated password.
	Register(context.Context, *pbetop.CreateUserRequest) (*pbetop.RegisterResponse, error)

	// Login
	//
	// Log the user in and generate access token.
	Login(context.Context, *pbetop.LoginRequest) (*pbetop.LoginResponse, error)

	// SessionInfo
	//
	// Return current session info.
	SessionInfo(context.Context, *cm.Empty) (*pbetop.LoginResponse, error)

	// SwitchAccount
	//
	// Response error if the user does not have permission to the requested account.
	SwitchAccount(context.Context, *pbetop.SwitchAccountRequest) (*pbetop.AccessTokenResponse, error)

	SendSTokenEmail(context.Context, *pbetop.SendSTokenEmailRequest) (*cm.MessageResponse, error)

	UpgradeAccessToken(context.Context, *pbetop.UpgradeAccessTokenRequest) (*pbetop.AccessTokenResponse, error)

	// ResetPassword
	//
	// Send email or sms to allow the user reset their password.
	ResetPassword(context.Context, *pbetop.ResetPasswordRequest) (*cm.MessageResponse, error)

	// ChangePassword
	//
	// Change the user password
	ChangePassword(context.Context, *pbetop.ChangePasswordRequest) (*cm.Empty, error)

	// ChangePasswordUsingToken
	//
	// Reset password by providing the token sent to email or phone
	ChangePasswordUsingToken(context.Context, *pbetop.ChangePasswordUsingTokenRequest) (*cm.Empty, error)

	SendEmailVerification(context.Context, *pbetop.SendEmailVerificationRequest) (*cm.MessageResponse, error)

	SendPhoneVerification(context.Context, *pbetop.SendPhoneVerificationRequest) (*cm.MessageResponse, error)

	VerifyEmailUsingToken(context.Context, *pbetop.VerifyEmailUsingTokenRequest) (*cm.MessageResponse, error)

	VerifyPhoneUsingToken(context.Context, *pbetop.VerifyPhoneUsingTokenRequest) (*cm.MessageResponse, error)

	UpdatePermission(context.Context, *pbetop.UpdatePermissionRequest) (*pbetop.UpdatePermissionResponse, error)

	UpdateReferenceUser(context.Context, *pbetop.UpdateReferenceUserRequest) (*cm.UpdatedResponse, error)

	UpdateReferenceSale(context.Context, *pbetop.UpdateReferenceSaleRequest) (*cm.UpdatedResponse, error)
}
