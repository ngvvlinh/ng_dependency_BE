package subscriptionproduct

import (
	"context"

	"o.o/api/top/types/etc/subscription_product_type"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateSubrProduct(context.Context, *CreateSubrProductArgs) (*SubscriptionProduct, error)
	UpdateSubrProduct(context.Context, *UpdateSubrProductArgs) error
	DeleteSubrProduct(ctx context.Context, ID dot.ID) error
}

type QueryService interface {
	GetSubrProductByID(ctx context.Context, ID dot.ID) (*SubscriptionProduct, error)
	ListSubrProducts(context.Context, *ListSubrProductsArgs) ([]*SubscriptionProduct, error)
}

// +convert:create=SubscriptionProduct
type CreateSubrProductArgs struct {
	Name        string
	Description string
	ImageURL    string
	Type        subscription_product_type.ProductSubscriptionType
}

// +convert:update=SubscriptionProduct(ID)
type UpdateSubrProductArgs struct {
	ID          dot.ID
	Name        string
	Description string
	ImageURL    string
}

type ListSubrProductsArgs struct {
	Type subscription_product_type.ProductSubscriptionType
}
