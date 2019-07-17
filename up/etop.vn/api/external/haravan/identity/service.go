package identity

import (
	"context"

	"etop.vn/api/meta"
)

type Aggregate interface {
	CreateExternalAccountHaravan(ctx context.Context, args *CreateExternalAccountHaravanArgs) (*ExternalAccountHaravan, error)

	UpdateExternalAccountHaravanToken(ctx context.Context, args *UpdateExternalAccountHaravanTokenArgs) (*ExternalAccountHaravan, error)

	ConnectCarrierServiceExternalAccountHaravan(ctx context.Context, args *ConnectCarrierServiceExternalAccountHaravanArgs) (*meta.Empty, error)

	DeleteConnectedCarrierServiceExternalAccountHaravan(ctx context.Context, args *DeleteConnectedCarrierServiceExternalAccountHaravanArgs) (*meta.Empty, error)
}

type QueryService interface {
	GetExternalAccountHaravanByShopID(context.Context, *GetExternalAccountHaravanByShopIDQueryArgs) (*ExternalAccountHaravan, error)

	GetExternalAccountHaravanByXShopID(context.Context, *GetExternalAccountHaravanByXShopIDQueryArgs) (*ExternalAccountHaravan, error)
}

type GetExternalAccountHaravanByShopIDQueryArgs struct {
	ShopID int64
}

type GetExternalAccountHaravanByXShopIDQueryArgs struct {
	ExternalShopID int
}

type CreateExternalAccountHaravanArgs struct {
	ShopID      int64
	Subdomain   string
	RedirectURI string
	Code        string
}

type UpdateExternalAccountHaravanTokenArgs struct {
	ShopID      int64
	Subdomain   string
	RedirectURI string
	Code        string
}

type ConnectCarrierServiceExternalAccountHaravanArgs struct {
	ShopID int64
}

type UpdateExternalShopIDAccountHaravanArgs struct {
	ShopID         int64
	ExternalShopID int
}

type DeleteConnectedCarrierServiceExternalAccountHaravanArgs struct {
	ShopID int64
}
