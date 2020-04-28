package modelx

import (
	"o.o/api/main/catalog"
	"o.o/api/main/catalog/types"
	"o.o/api/top/types/etc/status3"
	catalogmodel "o.o/backend/com/main/catalog/model"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

// deprecated
type GetShopVariantQuery struct {
	ShopID    dot.ID
	VariantID dot.ID

	ShopVariantStatus dot.NullInt

	Result *catalog.ShopVariant
}

// deprecated
type RemoveShopVariantsCommand struct {
	ShopID dot.ID
	IDs    []dot.ID

	Result struct {
		Removed int
	}
}

// deprecated
type UpdateShopVariantCommand struct {
	ShopID     dot.ID
	Variant    *catalogmodel.ShopVariant
	CostPrice  int
	Code       string
	Attributes []catalogmodel.ProductAttribute

	Result *catalog.ShopVariant
}

// deprecated
type AddShopProductsCommand struct {
	ShopID dot.ID
	IDs    []dot.ID

	Result struct {
		Products []*catalogmodel.ShopProduct
		Errors   []error
	}
}

// deprecated
type RemoveShopProductsCommand struct {
	ShopID dot.ID
	IDs    []dot.ID

	Result struct {
		Removed int
	}
}

// deprecated
type UpdateShopProductCommand struct {
	ShopID  dot.ID
	Product *catalogmodel.ShopProduct
	Code    string

	Result *catalog.ShopProductWithVariants
}

// deprecated
type UpdateShopProductsTagsCommand struct {
	ShopID     dot.ID
	ProductIDs []dot.ID

	Update *model.UpdateListRequest

	Result struct {
		Updated int
	}
}

// deprecated
type DeprecatedCreateVariantCommand struct {
	ShopID    dot.ID
	ProductID dot.ID
	// In `Dép Adidas Adilette Slides - Full Đỏ`, product_name is "Dép Adidas Adilette Slides"
	ProductName string
	// In `Dép Adidas Adilette Slides - Full Đỏ`, name is "Full Đỏ"
	Name        string
	Description string
	ShortDesc   string
	ImageURLs   []string
	Tags        []string
	Status      status3.Status
	ProductCode string
	VariantCode string

	QuantityAvailable int
	QuantityOnHand    int
	QuantityReserved  int

	ListPrice   int
	RetailPrice int
	CostPrice   int

	Attributes []*types.Attribute
	DescHTML   string

	Result *catalog.ShopProductWithVariants
}

// deprecated
type CreateShopCategoryCommand struct {
	ShopID   dot.ID
	Name     string
	ParentID dot.ID

	Result *catalogmodel.ShopCategory
}

// deprecated
type UpdateProductsShopCategoryCommand struct {
	CategoryID dot.ID
	ProductIDs []dot.ID
	ShopID     dot.ID

	Result struct {
		Updated int
	}
}

// deprecated
type GetShopCategoryQuery struct {
	ShopID     dot.ID
	CategoryID dot.ID

	Result *catalogmodel.ShopCategory
}

// deprecated
type GetProductSourceCategoriesQuery struct {
	ShopID dot.ID
	IDs    []dot.ID

	Result struct {
		Categories []*catalogmodel.ShopCategory
	}
}

// deprecated
type UpdateShopCategoryCommand struct {
	ID       dot.ID
	ShopID   dot.ID
	ParentID dot.ID
	Name     string

	Result *catalogmodel.ShopCategory
}

// deprecated
type RemoveShopCategoryCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result struct {
		Removed int
	}
}
