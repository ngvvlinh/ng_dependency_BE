package order

import "context"

type Aggregate interface {
	GetOrderByID(ctx context.Context, args GetOrderByIDArgs) (*Order, error)
}

type GetOrderByIDArgs struct {
	ID int64
}

type Order struct {
	ID int64
}
