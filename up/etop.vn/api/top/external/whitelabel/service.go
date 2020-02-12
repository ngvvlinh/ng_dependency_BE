package whitelabel

import (
	"context"
)

// +gen:apix
// +gen:swagger:doc-path=external/whitelabel

// +apix:path=/partner.Import
type ImportService interface {
	Products(context.Context, *ImportProductsRequest) (*ImportProductsResponse, error)
	Brands(context.Context, *ImportBrandsRequest) (*ImportBrandsResponse, error)
	Categories(context.Context, *ImportCategoriesRequest) (*ImportCategoriesResponse, error)
	Customers(context.Context, *ImportCustomersRequest) (*ImportCustomersResponse, error)
	Variants(context.Context, *ImportShopVariantsRequest) (*ImportShopVariantsResponse, error)
	Collections(context.Context, *ImportCollectionsRequest) (*ImportCollectionsResponse, error)
	ProductCollections(context.Context, *ImportProductCollectionsRequest) (*ImportProductCollectionsResponse, error)
}
