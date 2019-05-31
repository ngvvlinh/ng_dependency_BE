package identity

import (
	"context"
)

type Aggregate interface {
	CreateExternalAccountAhamove(ctx context.Context, cmd *CreateExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)
}

type QueryService interface {
	GetShopByID(context.Context, *GetShopByIDQueryArgs) (*GetShopByIDQueryResult, error)

	GetExternalAccountAhamoveByPhone(context.Context, *GetExternalAccountAhamoveByPhoneArgs) (*ExternalAccountAhamove, error)
}

//-- queries --//
type GetShopByIDQueryArgs struct {
	ID int64
}

type GetShopByIDQueryResult struct {
	Shop *Shop
}

type GetExternalAccountAhamoveByPhoneArgs struct {
	Phone   string
	OwnerID int64
}

//-- commands --//
type CreateExternalAccountAhamoveArgs struct {
	OwnerID int64 // user id
	Phone   string
	Name    string
}
