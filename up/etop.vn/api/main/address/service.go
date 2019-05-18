package address

import "context"

type QueryService interface {
	GetAddressByID(context.Context, *GetAddressByIDQueryArgs) (*Address, error)
}

type GetAddressByIDQueryArgs struct {
	ID int64
}
