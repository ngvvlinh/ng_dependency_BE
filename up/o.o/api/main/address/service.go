package address

import (
	"context"

	"o.o/capi/dot"
)

// +gen:api

type QueryService interface {
	GetAddressByID(context.Context, *GetAddressByIDQueryArgs) (*Address, error)
}

type GetAddressByIDQueryArgs struct {
	ID dot.ID
}
