package model

import (
	"time"

	sq "etop.vn/backend/pkg/common/sql"
)

var _ = sqlgenProductBrand(&ProductBrand{})

type ProductBrand struct {
	ID          int64
	Name        string
	Description string
	Policy      string
	ImageURLs   []string `sq:"'image_urls'"`
	SupplierID  int64
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
}

var _ = sqlgenProductBrandExtended(
	&ProductBrandExtended{}, &ProductBrand{}, sq.AS("pb"),
	sq.LEFT_JOIN, &Supplier{}, sq.AS("s"), "pb.supplier_id = s.id",
)

type ProductBrandExtended struct {
	*ProductBrand
	Supplier *Supplier
}

type GetProductBrandQuery struct {
	SupplierID int64
	ID         int64

	Result *ProductBrandExtended
}

type GetProductBrandsQuery struct {
	SupplierID int64
	Ids        []int64

	Result struct {
		Brands []*ProductBrandExtended
	}
}

type CreateProductBrandCommand struct {
	Brand *ProductBrand

	Result *ProductBrand
}

type UpdateProductBrandCommand struct {
	Brand *ProductBrand

	Result *ProductBrand
}

type DeleteProductBrandCommand struct {
	ID         int64
	SupplierID int64
}
