package carrying

import (
	"context"

	"o.o/api/meta"
	"o.o/api/shopping"
	. "o.o/capi/dot"
	dot "o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateCarrier(ctx context.Context, _ *CreateCarrierArgs) (*ShopCarrier, error)

	UpdateCarrier(ctx context.Context, _ *UpdateCarrierArgs) (*ShopCarrier, error)

	DeleteCarrier(ctx context.Context, ID dot.ID, shopID dot.ID) (deleted int, _ error)
}

type QueryService interface {
	GetCarrierByID(context.Context, *shopping.IDQueryShopArg) (*ShopCarrier, error)

	ListCarriers(context.Context, *shopping.ListQueryShopArgs) (*CarriersResponse, error)

	ListCarriersByIDs(context.Context, *shopping.IDsQueryShopArgs) (*CarriersResponse, error)
}

//-- queries --//

type CarriersResponse struct {
	Carriers []*ShopCarrier
	Count    int
	Paging   meta.PageInfo
}

//-- commands --//

// +convert:create=ShopCarrier
type CreateCarrierArgs struct {
	ShopID   dot.ID
	FullName string
	Note     string
}

// +convert:update=ShopCarrier(ID,ShopID)
type UpdateCarrierArgs struct {
	ID       dot.ID
	ShopID   dot.ID
	FullName NullString
	Note     NullString
}
