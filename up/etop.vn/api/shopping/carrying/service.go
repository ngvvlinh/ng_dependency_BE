package carrying

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateCarrier(ctx context.Context, _ *CreateCarrierArgs) (*ShopCarrier, error)

	UpdateCarrier(ctx context.Context, _ *UpdateCarrierArgs) (*ShopCarrier, error)

	DeleteCarrier(ctx context.Context, ID int64, shopID int64) (deleted int, _ error)
}

type QueryService interface {
	GetCarrierByID(context.Context, *shopping.IDQueryShopArg) (*ShopCarrier, error)

	ListCarriers(context.Context, *shopping.ListQueryShopArgs) (*CarriersResponse, error)

	ListCarriersByIDs(context.Context, *shopping.IDsQueryShopArgs) (*CarriersResponse, error)
}

//-- queries --//

type CarriersResponse struct {
	Carriers []*ShopCarrier
	Count    int32
	Paging   meta.PageInfo
}

//-- commands --//

// +convert:create=ShopCarrier
type CreateCarrierArgs struct {
	ShopID   int64
	FullName string
	Note     string
}

// +convert:update=ShopCarrier(ID,ShopID)
type UpdateCarrierArgs struct {
	ID       int64
	ShopID   int64
	FullName NullString
	Note     NullString
}
