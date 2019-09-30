package vendoring

import (
	"context"

	"etop.vn/api/meta"

	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateVendor(ctx context.Context, _ *CreateVendorArgs) (*ShopVendor, error)

	UpdateVendor(ctx context.Context, _ *UpdateVendorArgs) (*ShopVendor, error)

	DeleteVendor(ctx context.Context, ID int64, shopID int64) (deleted int, _ error)
}

type QueryService interface {
	GetVendorByID(context.Context, *shopping.IDQueryShopArg) (*ShopVendor, error)

	ListVendors(context.Context, *shopping.ListQueryShopArgs) (*VendorsResponse, error)

	ListVendorsByIDs(context.Context, *shopping.IDsQueryShopArgs) (*VendorsResponse, error)
}

//-- queries --//
type VendorsResponse struct {
	Vendors []*ShopVendor
	Count   int32
	Paging  meta.PageInfo
}

//-- commands --//

// +convert:create=ShopVendor
type CreateVendorArgs struct {
	ShopID   int64
	FullName string
	Note     string
}

// +convert:update=ShopVendor(ID,ShopID)
type UpdateVendorArgs struct {
	ID       int64
	ShopID   int64
	FullName NullString
	Note     NullString
}
