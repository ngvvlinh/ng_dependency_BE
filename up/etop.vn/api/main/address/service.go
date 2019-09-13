package address

import "context"

// +gen:api

type QueryService interface {
	GetAddressByID(context.Context, *GetAddressByIDQueryArgs) (*Address, error)
}

type GetAddressByIDQueryArgs struct {
	ID int64
}
