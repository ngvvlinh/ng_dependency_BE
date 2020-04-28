package suppliering

import (
	"context"

	"o.o/api/meta"
	"o.o/api/shopping"
	. "o.o/capi/dot"
	dot "o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateSupplier(ctx context.Context, _ *CreateSupplierArgs) (*ShopSupplier, error)

	UpdateSupplier(ctx context.Context, _ *UpdateSupplierArgs) (*ShopSupplier, error)

	DeleteSupplier(ctx context.Context, ID dot.ID, shopID dot.ID) (deleted int, _ error)
}

type QueryService interface {
	GetSupplierByID(context.Context, *shopping.IDQueryShopArg) (*ShopSupplier, error)

	ListSuppliers(context.Context, *shopping.ListQueryShopArgs) (*SuppliersResponse, error)

	ListSuppliersByIDs(context.Context, *shopping.IDsQueryShopArgs) (*SuppliersResponse, error)
}

//-- queries --//
type SuppliersResponse struct {
	Suppliers []*ShopSupplier
	Paging    meta.PageInfo
}

//-- commands --//

// +convert:create=ShopSupplier
type CreateSupplierArgs struct {
	ShopID            dot.ID
	FullName          string
	Phone             string
	Email             string
	CompanyName       string
	TaxNumber         string
	HeadquaterAddress string
	Note              string
}

// +convert:update=ShopSupplier(ID,ShopID)
type UpdateSupplierArgs struct {
	ID                dot.ID
	ShopID            dot.ID
	FullName          NullString
	Note              NullString
	Phone             NullString
	Email             NullString
	CompanyName       NullString
	TaxNumber         NullString
	HeadquaterAddress NullString
}
