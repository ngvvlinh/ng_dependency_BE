package identity

import (
	"context"
)

type Aggregate interface {
	// External Account Ahamove

	CreateExternalAccountAhamove(ctx context.Context, args *CreateExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	RequestVerifyExternalAccountAhamove(ctx context.Context, args *RequestVerifyExternalAccountAhamoveArgs) (*RequestVerifyExternalAccountAhamoveResult, error)

	UpdateVerifiedExternalAccountAhamove(ctx context.Context, args *UpdateVerifiedExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	UpdateExternalAccountAhamoveVerification(ctx context.Context, args *UpdateExternalAccountAhamoveVerificationArgs) (*ExternalAccountAhamove, error)
}

type QueryService interface {
	GetShopByID(context.Context, *GetShopByIDQueryArgs) (*GetShopByIDQueryResult, error)

	GetUserByID(context.Context, *GetUserByIDQueryArgs) (*User, error)

	GetExternalAccountAhamove(context.Context, *GetExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	GetExternalAccountAhamoveByExternalID(context.Context, *GetExternalAccountAhamoveByExternalIDQueryArgs) (*ExternalAccountAhamove, error)
}

//-- queries --//
type GetShopByIDQueryArgs struct {
	ID int64
}

type GetShopByIDQueryResult struct {
	Shop *Shop
}

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
