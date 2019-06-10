package identity

import (
	"context"
)

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
