package suppliering

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateSupplier(ctx context.Context, _ *CreateSupplierArgs) (*ShopSupplier, error)

	UpdateSupplier(ctx context.Context, _ *UpdateSupplierArgs) (*ShopSupplier, error)

	DeleteSupplier(ctx context.Context, ID int64, shopID int64) (deleted int, _ error)
}

type QueryService interface {
	GetSupplierByID(context.Context, *shopping.IDQueryShopArg) (*ShopSupplier, error)

	ListSuppliers(context.Context, *shopping.ListQueryShopArgs) (*SuppliersResponse, error)

	ListSuppliersByIDs(context.Context, *shopping.IDsQueryShopArgs) (*SuppliersResponse, error)
}

//-- queries --//
type SuppliersResponse struct {
	Suppliers []*ShopSupplier
	Count     int32
	Paging    meta.PageInfo
}

//-- commands --//

// +convert:create=ShopSupplier
type CreateSupplierArgs struct {
	ShopID            int64
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
	ID                int64
	ShopID            int64
	FullName          NullString
	Note              NullString
	Phone             NullString
	Email             NullString
	CompanyName       NullString
	TaxNumber         NullString
	HeadquaterAddress NullString
}
