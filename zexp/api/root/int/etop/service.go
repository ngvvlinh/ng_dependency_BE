package etop

import (
	"context"

	cm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
)

// +gen:apix
// +gen:apix:doc-path=etop

// +apix:path=/etop.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/etop.User
type UserService interface {
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

// +apix:path=/etop.Account
type AccountService interface {
	UpdateURLSlug(context.Context, *pbetop.UpdateURLSlugRequest) (*cm.Empty, error)
	GetPublicPartnerInfo(context.Context, *cm.IDRequest) (*pbetop.PublicAccountInfo, error)

	// leave ids empty to get all connected partners
	GetPublicPartners(context.Context, *cm.IDsRequest) (*pbetop.PublicAccountsResponse, error)
}

// +apix:path=/etop.Relationship
type RelationshipService interface {
	InviteUserToAccount(context.Context, *pbetop.InviteUserToAccountRequest) (*pbetop.UserAccountInfo, error)
	AnswerInvitation(context.Context, *pbetop.AnswerInvitationRequest) (*pbetop.UserAccountInfo, error)

	GetUsersInCurrentAccounts(context.Context, *pbetop.GetUsersInCurrentAccountsRequest) (*pbetop.ProtectedUsersResponse, error)

	LeaveAccount(context.Context, *pbetop.LeaveAccountRequest) (*cm.Empty, error)

	RemoveUserFromCurrentAccount(context.Context, *pbetop.RemoveUserFromCurrentAccountRequest) (*cm.Empty, error)
}

// +apix:path=/etop.Location
type LocationService interface {
	GetProvinces(context.Context, *cm.Empty) (*pbetop.GetProvincesResponse, error)

	GetDistricts(context.Context, *cm.Empty) (*pbetop.GetDistrictsResponse, error)

	GetDistrictsByProvince(context.Context, *pbetop.GetDistrictsByProvinceRequest) (*pbetop.GetDistrictsResponse, error)

	GetWards(context.Context, *cm.Empty) (*pbetop.GetWardsResponse, error)

	GetWardsByDistrict(context.Context, *pbetop.GetWardsByDistrictRequest) (*pbetop.GetWardsResponse, error)

	ParseLocation(context.Context, *pbetop.ParseLocationRequest) (*pbetop.ParseLocationResponse, error)
}

// +apix:path=/etop.Bank
type BankService interface {
	GetBanks(context.Context, *cm.Empty) (*pbetop.GetBanksResponse, error)

	GetProvincesByBank(context.Context, *pbetop.GetProvincesByBankResquest) (*pbetop.GetBankProvincesResponse, error)

	GetBranchesByBankProvince(context.Context, *pbetop.GetBranchesByBankProvinceResquest) (*pbetop.GetBranchesByBankProvinceResponse, error)
}

// +apix:path=/etop.Address
type AddressService interface {
	CreateAddress(context.Context, *pbetop.CreateAddressRequest) (*pbetop.Address, error)

	GetAddresses(context.Context, *cm.Empty) (*pbetop.GetAddressResponse, error)

	UpdateAddress(context.Context, *pbetop.UpdateAddressRequest) (*pbetop.Address, error)

	RemoveAddress(context.Context, *cm.IDRequest) (*cm.Empty, error)
}

// +apix:path=/etop.Invitation
type InvitationService interface {
	AcceptInvitation(context.Context, *pbetop.AcceptInvitationRequest) (*cm.UpdatedResponse, error)
	RejectInvitation(context.Context, *pbetop.RejectInvitationRequest) (*cm.UpdatedResponse, error)
	GetInvitationByToken(context.Context, *pbetop.GetInvitationByTokenRequest) (*pbetop.Invitation, error)
	GetInvitations(context.Context, *pbetop.GetInvitationsRequest) (*pbetop.InvitationsResponse, error)
}
