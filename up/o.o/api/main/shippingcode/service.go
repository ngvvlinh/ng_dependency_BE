package shippingcode

import (
	"context"

	cm "o.o/api/top/types/common"
)

// +gen:api

type Aggregate interface{}

type QueryService interface {
	GenerateShippingCode(context.Context, *cm.Empty) (string, error)
}

type GenerateShippingCodeArgs struct{}
