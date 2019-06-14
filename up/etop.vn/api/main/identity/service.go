package identity

import (
	"context"
)

type Aggregate interface {
	// External Account Ahamove

	CreateExternalAccountAhamove(ctx context.Context, args *CreateExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	RequestVerifyExternalAccountAhamove(ctx context.Context, args *RequestVerifyExternalAccountAhamoveArgs) (*RequestVerifyExternalAccountAhamoveResult, error)

	UpdateVerifiedExternalAccountAhamove(ctx context.Context, args *UpdateVerifiedExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	UpdateExternalAccountAhamoveVerificationImages(ctx context.Context, args *UpdateExternalAccountAhamoveVerificationImagesArgs) (*ExternalAccountAhamove, error)
}

type QueryService interface {
	GetShopByID(context.Context, *GetShopByIDQueryArgs) (*GetShopByIDQueryResult, error)

	GetUserByID(context.Context, *GetUserByIDQueryArgs) (*User, error)

	GetExternalAccountAhamove(context.Context, *GetExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)
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

//-- commands --//
type CreateExternalAccountAhamoveArgs struct {
	OwnerID int64 // user id
	Phone   string
	Name    string
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

type UpdateExternalAccountAhamoveVerificationImagesArgs struct {
	UserID         int64
	IDCardFrontImg string
	IDCardBackImg  string
	PortraitImg    string
}
