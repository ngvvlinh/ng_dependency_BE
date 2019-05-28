package identity

import (
	"context"

	"etop.vn/api/meta"
)

type QueryBus struct{ meta.Bus }

type QueryService interface {
	GetShopByID(context.Context, *GetShopByIDQueryArgs) (*GetShopByIDQueryResult, error)
}

//-- queries --//
type GetShopByIDQueryArgs struct {
	ID int64
}

type GetShopByIDQueryResult struct {
	Shop *Shop
}
