package identity

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/capi/dot"
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
	ShopID dot.ID
}

type GetExternalAccountHaravanByXShopIDQueryArgs struct {
	ExternalShopID int
}

type CreateExternalAccountHaravanArgs struct {
	ShopID      dot.ID
	Subdomain   string
	RedirectURI string
	Code        string
}

type UpdateExternalAccountHaravanTokenArgs struct {
	ShopID      dot.ID
	Subdomain   string
	RedirectURI string
	Code        string
}

type ConnectCarrierServiceExternalAccountHaravanArgs struct {
	ShopID dot.ID
}

type UpdateExternalShopIDAccountHaravanArgs struct {
	ShopID         dot.ID
	ExternalShopID int
}

type DeleteConnectedCarrierServiceExternalAccountHaravanArgs struct {
	ShopID dot.ID
}
