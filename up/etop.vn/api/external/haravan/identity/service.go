package identity

import (
	"context"

	"etop.vn/api/meta"
)

// +gen:api

type Aggregate interface {
	CreateExternalAccountHaravan(context.Context, *CreateExternalAccountHaravanArgs) (*ExternalAccountHaravan, error)

	UpdateExternalAccountHaravanToken(context.Context, *UpdateExternalAccountHaravanTokenArgs) (*ExternalAccountHaravan, error)

	ConnectCarrierServiceExternalAccountHaravan(context.Context, *ConnectCarrierServiceExternalAccountHaravanArgs) (*meta.Empty, error)

	DeleteConnectedCarrierServiceExternalAccountHaravan(context.Context, *DeleteConnectedCarrierServiceExternalAccountHaravanArgs) (*meta.Empty, error)
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
