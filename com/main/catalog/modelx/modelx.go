package modelx

import (
	"etop.vn/api/main/catalog"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/pkg/etop/model"
)

// deprecated
type GetShopVariantQuery struct {
	ShopID    int64
	VariantID int64

	ShopVariantStatus *int

	Result *catalog.ShopVariant
}

// deprecated
type RemoveShopVariantsCommand struct {
	ShopID int64
	IDs    []int64

	Result struct {
		Removed int
	}
}

// deprecated
type UpdateShopVariantCommand struct {
	ShopID     int64
	Variant    *catalogmodel.ShopVariant
	CostPrice  int32
	Code       string
	Attributes []catalogmodel.ProductAttribute

	Result *catalog.ShopVariant
}

// deprecated
type AddShopProductsCommand struct {
	ShopID int64
	IDs    []int64

	Result struct {
		Products []*catalogmodel.ShopProduct
		Errors   []error
	}
}

// deprecated
type RemoveShopProductsCommand struct {
	ShopID int64
	IDs    []int64

	Result struct {
		Removed int
	}
}

// deprecated
type UpdateShopProductCommand struct {
	ShopID  int64
	Product *catalogmodel.ShopProduct
	Code    string

	Result *catalog.ShopProductWithVariants
}

// deprecated
type UpdateShopProductsTagsCommand struct {
	ShopID     int64
	ProductIDs []int64

	Update *model.UpdateListRequest

	Result struct {
		Updated int
	}
}

// deprecated
type DeprecatedCreateVariantCommand struct {
	ShopID    int64
	ProductID int64
	// In `Dép Adidas Adilette Slides - Full Đỏ`, product_name is "Dép Adidas Adilette Slides"
	ProductName string
	// In `Dép Adidas Adilette Slides - Full Đỏ`, name is "Full Đỏ"
	Name        string
	Description string
	ShortDesc   string
	ImageURLs   []string
	Tags        []string
	Status      model.Status3
	ProductCode string
	VariantCode string

	QuantityAvailable int
	QuantityOnHand    int
	QuantityReserved  int

	ListPrice   int32
	RetailPrice int32
	CostPrice   int32

	Attributes []*catalogmodel.ProductAttribute
	DescHTML   string

	Result *catalog.ShopProductWithVariants
}

// deprecated
type CreateShopCategoryCommand struct {
	ShopID   int64
	Name     string
	ParentID int64

	Result *catalogmodel.ShopCategory
}

// deprecated
type UpdateProductsShopCategoryCommand struct {
	CategoryID int64
	ProductIDs []int64
	ShopID     int64

	Result struct {
		Updated int
	}
}

// deprecated
type GetShopCategoryQuery struct {
	ShopID     int64
	CategoryID int64

	Result *catalogmodel.ShopCategory
}

// deprecated
type GetProductSourceCategoriesQuery struct {
	ShopID int64
	IDs    []int64

	Result struct {
		Categories []*catalogmodel.ShopCategory
	}
}

// deprecated
type UpdateShopCategoryCommand struct {
	ID       int64
	ShopID   int64
	ParentID int64
	Name     string

	Result *catalogmodel.ShopCategory
}

// deprecated
type RemoveShopCategoryCommand struct {
	ID     int64
	ShopID int64

	Result struct {
		Removed int
	}
}
