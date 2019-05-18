package identity

import "context"

type QueryService interface {
	GetShopByID(context.Context, *GetShopByIDQueryArgs) (*Shop, error)
}

type GetShopByIDQueryArgs struct {
	ID int64
}
