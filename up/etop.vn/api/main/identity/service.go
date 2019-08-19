package identity

import (
	"context"
)

type Aggregate interface {
	// -- User -- //
	UpdateUserReferenceUserID(context.Context, *UpdateUserReferenceUserIDArgs) error

	UpdateUserReferenceSaleID(context.Context, *UpdateUserReferenceSaleIDArgs) error

	// -- External Account Ahamove -- //
	// TODO: move External Account Ahamove to its own module
	CreateExternalAccountAhamove(context.Context, *CreateExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	RequestVerifyExternalAccountAhamove(context.Context, *RequestVerifyExternalAccountAhamoveArgs) (*RequestVerifyExternalAccountAhamoveResult, error)

	UpdateVerifiedExternalAccountAhamove(context.Context, *UpdateVerifiedExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	UpdateExternalAccountAhamoveVerification(context.Context, *UpdateExternalAccountAhamoveVerificationArgs) (*ExternalAccountAhamove, error)

	// -- Affiliate -- //
	CreateAffiliate(context.Context, *CreateAffiliateArgs) (*Affiliate, error)

	UpdateAffiliate(context.Context, *UpdateAffiliateArgs) (*Affiliate, error)

	DeleteAffiliate(context.Context, *DeleteAffiliateArgs) error
}

type QueryService interface {
	// -- Shop -- //
	GetShopByID(ctx context.Context, ID int64) (*Shop, error)

	// -- User -- //
	GetUserByID(context.Context, *GetUserByIDQueryArgs) (*User, error)

	GetUserByPhone(ctx context.Context, phone string) (*User, error)

	// -- ExternalAccountAhamove -- //
	GetExternalAccountAhamove(context.Context, *GetExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	GetExternalAccountAhamoveByExternalID(context.Context, *GetExternalAccountAhamoveByExternalIDQueryArgs) (*ExternalAccountAhamove, error)

	// -- Affiliate -- //
	GetAffiliateByID(ctx context.Context, ID int64) (*Affiliate, error)

	GetAffiliateWithPermission(ctx context.Context, AffiliateID int64, UserID int64) (*GetAffiliateWithPermissionResult, error)
}

//-- queries --//
type GetUserByIDQueryArgs struct {
	UserID int64
}

type GetExternalAccountAhamoveArgs struct {
	OwnerID int64
	Phone   string
}

type GetExternalAccountAhamoveByExternalIDQueryArgs struct {
	ExternalID string
}

//-- commands --//
type CreateExternalAccountAhamoveArgs struct {
	OwnerID int64 // user id
	Phone   string
	Name    string
	Address string
}

type RequestVerifyExternalAccountAhamoveArgs struct {
	OwnerID int64
	Phone   string
}

type RequestVerifyExternalAccountAhamoveResult struct {
}

type UpdateVerifiedExternalAccountAhamoveArgs struct {
	OwnerID int64
	Phone   string
}

type UpdateExternalAccountAhamoveVerificationArgs struct {
	OwnerID             int64
	Phone               string
	IDCardFrontImg      string
	IDCardBackImg       string
	PortraitImg         string
	WebsiteURL          string
	FanpageURL          string
	CompanyImgs         []string
	BusinessLicenseImgs []string
}

type UpdateUserReferenceUserIDArgs struct {
	UserID       int64
	RefUserPhone string
}

type UpdateUserReferenceSaleIDArgs struct {
	UserID       int64
	RefSalePhone string
}

type GetAffiliateWithPermissionResult struct {
	Affiliate  *Affiliate
	Permission Permission
}

type CreateAffiliateArgs struct {
	Name    string
	OwnerID int64
	Phone   string
	Email   string
	IsTest  bool
}

type UpdateAffiliateArgs struct {
	ID      int64
	OwnerID int64
	Phone   string
	Email   string
	Name    string
}

type DeleteAffiliateArgs struct {
	ID      int64
	OwnerID int64
}
