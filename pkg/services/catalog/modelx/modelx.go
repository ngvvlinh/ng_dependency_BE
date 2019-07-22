package modelx

import (
	"etop.vn/api/main/catalog"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
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
type CreateVariantCommand struct {
	ShopID    int64
	ProductID int64
	// In `Dép Adidas Adilette Slides - Full Đỏ`, product_name is "Dép Adidas Adilette Slides"
	ProductName string
	// In `Dép Adidas Adilette Slides - Full Đỏ`, name is "Full Đỏ"
	Name              string
	Description       string
	ShortDesc         string
	ImageURLs         []string
	Tags              []string
	Status            model.Status3
	ListPrice         int32
	SKU               string
	Code              string
	QuantityAvailable int
	QuantityOnHand    int
	QuantityReserved  int
	CostPrice         int32
	Attributes        []catalogmodel.ProductAttribute
	DescHTML          string

	Result *catalog.ShopProductWithVariants
}

// deprecated
type CreateProductSourceCategoryCommand struct {
	ShopID   int64
	Name     string
	ParentID int64

	Result *catalogmodel.ProductSourceCategory
}

// deprecated
type UpdateProductsProductSourceCategoryCommand struct {
	CategoryID int64
	ProductIDs []int64
	ShopID     int64

	Result struct {
		Updated int
	}
}

// deprecated
type GetProductSourceCategoryQuery struct {
	ShopID     int64
	CategoryID int64

	Result *catalogmodel.ProductSourceCategory
}

// deprecated
type GetProductSourceCategoriesQuery struct {
	ShopID int64
	IDs    []int64

	Result struct {
		Categories []*catalogmodel.ProductSourceCategory
	}
}

// deprecated
type UpdateShopProductSourceCategoryCommand struct {
	ID       int64
	ShopID   int64
	ParentID int64
	Name     string

	Result *catalogmodel.ProductSourceCategory
}

// deprecated
type RemoveShopProductSourceCategoryCommand struct {
	ID     int64
	ShopID int64

	Result struct {
		Removed int
	}
}
