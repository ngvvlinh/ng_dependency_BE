package etop

import (
	"context"

	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=etop
// +gen:swagger:description=description.md

// +apix:path=/etop.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/etop.User
type UserService interface {
	// Register
	//
	// Register a new user or after a user has login using generated password.
	Register(context.Context, *CreateUserRequest) (*RegisterResponse, error)

	RegisterUsingToken(context.Context, *CreateUserRequest) (*RegisterResponse, error)

	// Login
	//
	// Log the user in and generate access token.
	Login(context.Context, *LoginRequest) (*LoginResponse, error)

	// SessionInfo
	//
	// Return current session info.
	SessionInfo(context.Context, *cm.Empty) (*LoginResponse, error)

	InitSession(context.Context, *cm.Empty) (*LoginResponse, error)

	// SwitchAccount
	//
	// Response error if the user does not have permission to the requested account.
	SwitchAccount(context.Context, *SwitchAccountRequest) (*AccessTokenResponse, error)

	SendSTokenEmail(context.Context, *SendSTokenEmailRequest) (*cm.MessageResponse, error)

	UpgradeAccessToken(context.Context, *UpgradeAccessTokenRequest) (*AccessTokenResponse, error)

	UpdateUserEmail(context.Context, *UpdateUserEmailRequest) (*UpdateUserEmailResponse, error)

	UpdateUserPhone(context.Context, *UpdateUserPhoneRequest) (*UpdateUserPhoneResponse, error)

	// ResetPassword
	//
	// Send email or sms to allow the user reset their password.
	ResetPassword(context.Context, *ResetPasswordRequest) (*ResetPasswordResponse, error)
	// ChangePassword
	//
	// Change the user password
	ChangePassword(context.Context, *ChangePasswordRequest) (*cm.Empty, error)

	ChangeRefAff(ctx context.Context, request *ChangeUserRefAffRequest) (*cm.Empty, error)
	// ChangePasswordUsingToken
	//
	// Reset password by providing the token sent to email or phone
	ChangePasswordUsingToken(context.Context, *ChangePasswordUsingTokenRequest) (*cm.Empty, error)

	SendEmailVerification(context.Context, *SendEmailVerificationRequest) (*cm.MessageResponse, error)

	SendEmailVerificationUsingOTP(context.Context, *SendEmailVerificationUsingOTPRequest) (*cm.MessageResponse, error)

	SendPhoneVerification(context.Context, *SendPhoneVerificationRequest) (*cm.MessageResponse, error)

	VerifyEmailUsingToken(context.Context, *VerifyEmailUsingTokenRequest) (*cm.MessageResponse, error)

	VerifyEmailUsingOTP(context.Context, *VerifyEmailUsingOTPRequest) (*cm.MessageResponse, error)

	VerifyPhoneUsingToken(context.Context, *VerifyPhoneUsingTokenRequest) (*cm.MessageResponse, error)

	VerifyPhoneResetPasswordUsingToken(context.Context, *VerifyPhoneResetPasswordUsingTokenRequest) (*VerifyPhoneResetPasswordUsingTokenResponse, error)

	UpdatePermission(context.Context, *UpdatePermissionRequest) (*UpdatePermissionResponse, error)

	UpdateReferenceUser(context.Context, *UpdateReferenceUserRequest) (*cm.UpdatedResponse, error)

	UpdateReferenceSale(context.Context, *UpdateReferenceSaleRequest) (*cm.UpdatedResponse, error)

	CheckUserRegistration(context.Context, *GetUserByPhoneRequest) (*GetUserByPhoneResponse, error)
}

// +apix:path=/etop.Account
type AccountService interface {
	UpdateURLSlug(context.Context, *UpdateURLSlugRequest) (*cm.Empty, error)
	GetPublicPartnerInfo(context.Context, *cm.IDRequest) (*PublicAccountInfo, error)

	// leave ids empty to get all connected partners
	GetPublicPartners(context.Context, *cm.IDsRequest) (*PublicAccountsResponse, error)
}

// +apix:path=/etop.Location
type LocationService interface {
	GetProvinces(context.Context, *cm.Empty) (*GetProvincesResponse, error)

	GetDistricts(context.Context, *cm.Empty) (*GetDistrictsResponse, error)

	GetDistrictsByProvince(context.Context, *GetDistrictsByProvinceRequest) (*GetDistrictsResponse, error)

	GetWards(context.Context, *cm.Empty) (*GetWardsResponse, error)

	GetWardsByDistrict(context.Context, *GetWardsByDistrictRequest) (*GetWardsResponse, error)

	ParseLocation(context.Context, *ParseLocationRequest) (*ParseLocationResponse, error)
}

// +apix:path=/etop.Bank
type BankService interface {
	GetBanks(context.Context, *cm.Empty) (*GetBanksResponse, error)

	GetProvincesByBank(context.Context, *GetProvincesByBankResquest) (*GetBankProvincesResponse, error)

	GetBranchesByBankProvince(context.Context, *GetBranchesByBankProvinceResquest) (*GetBranchesByBankProvinceResponse, error)

	GetBankProvinces(context.Context, *GetBankProvincesRequest) (*GetBankProvinceResponse, error)

	GetBankBranches(context.Context, *GetBankBranchesRequest) (*GetBankBranchesResponse, error)
}

// +apix:path=/etop.Address
type AddressService interface {
	CreateAddress(context.Context, *CreateAddressRequest) (*Address, error)

	GetAddresses(context.Context, *cm.Empty) (*GetAddressResponse, error)

	UpdateAddress(context.Context, *UpdateAddressRequest) (*Address, error)

	RemoveAddress(context.Context, *cm.IDRequest) (*cm.Empty, error)
}

// +apix:path=/etop.UserRelationship
// +wrapper:endpoint-prefix=UserRelationship
type UserRelationshipService interface {
	AcceptInvitation(context.Context, *AcceptInvitationRequest) (*cm.UpdatedResponse, error)
	RejectInvitation(context.Context, *RejectInvitationRequest) (*cm.UpdatedResponse, error)
	GetInvitationByToken(context.Context, *GetInvitationByTokenRequest) (*Invitation, error)
	GetInvitations(context.Context, *GetInvitationsRequest) (*InvitationsResponse, error)
	LeaveAccount(context.Context, *UserRelationshipLeaveAccountRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/etop.AccountRelationship
// +wrapper:endpoint-prefix=AccountRelationship
type AccountRelationshipService interface {
	CreateInvitation(context.Context, *CreateInvitationRequest) (*Invitation, error)
	ResendInvitation(context.Context, *ResendInvitationRequest) (*Invitation, error)
	GetInvitations(context.Context, *GetInvitationsRequest) (*InvitationsResponse, error)
	DeleteInvitation(context.Context, *DeleteInvitationRequest) (*cm.UpdatedResponse, error)
	UpdatePermission(context.Context, *UpdateAccountUserPermissionRequest) (*Relationship, error)
	UpdateRelationship(context.Context, *UpdateRelationshipRequest) (*Relationship, error)
	GetRelationships(context.Context, *GetRelationshipsRequest) (*RelationshipsResponse, error)
	RemoveUser(context.Context, *RemoveUserRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/etop.Ecom
// +wrapper:endpoint-prefix=Ecom
type EcomService interface {
	SessionInfo(context.Context, *cm.Empty) (*EcomSessionInfoResponse, error)
}

// +apix:path=/etop.Ticket
type TicketService interface {
	GetTicketLabels(context.Context, *GetTicketLabelsRequest) (*GetTicketLabelsResponse, error)
}
