package modelx

import (
	"etop.vn/api/main/catalog"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

type DeprecatedGetProductsExtendedQuery struct {
	Paging  *cm.Paging
	Filters []cm.Filter
	IDs     []int64
	model.StatusQuery
	ProductSourceType string

	Result struct {
		Products []*catalogmodel.ProductFtVariant
		Total    int
	}
}

type GetVariantsExtendedQuery struct {
	ProductSourceID int64

	Paging  *cm.Paging
	Filters []cm.Filter
	IDs     []int64
	Codes   []string
	EdCodes []string
	model.StatusQuery

	SkipPaging bool

	Result struct {
		Variants []*catalogmodel.VariantExtended
		Total    int
	}
}

func (g *GetVariantsExtendedQuery) IsPaging() bool {
	return !g.SkipPaging
}

type GetShopVariantQuery struct {
	ShopID    int64
	VariantID int64

	ShopVariantStatus *int

	Result *catalog.ShopVariantExtended
}

type AddProductsToShopCollectionCommand struct {
	ShopID       int64
	ProductIDs   []int64
	CollectionID int64

	Result struct {
		Updated int
		Errors  []error
	}
}

type RemoveProductsFromShopCollectionCommand struct {
	ShopID       int64
	ProductIDs   []int64
	CollectionID int64

	Result struct {
		Updated int
	}
}

type AddShopVariantsCommand struct {
	ShopID int64
	IDs    []int64

	Result struct {
		Variants []*catalog.ShopVariantExtended
		Errors   []error
	}
}

type RemoveShopVariantsCommand struct {
	ShopID          int64
	IDs             []int64
	ProductSourceID int64

	Result struct {
		Removed int
	}
}

type UpdateShopVariantCommand struct {
	ShopID          int64
	Variant         *catalogmodel.ShopVariant
	CostPrice       int32
	Code            string
	Attributes      []catalogmodel.ProductAttribute
	ProductSourceID int64

	Result *catalog.ShopVariantExtended
}

type UpdateShopVariantsCommand struct {
	ShopID          int64
	Products        []*catalogmodel.ShopVariant
	ProductSourceID int64

	Result struct {
		Products []*catalogmodel.ShopVariantExtended
		Errors   []error
	}
}

type AddShopProductsCommand struct {
	ShopID int64
	IDs    []int64

	Result struct {
		Products []*catalogmodel.ShopProduct
		Errors   []error
	}
}

type RemoveShopProductsCommand struct {
	ShopID          int64
	IDs             []int64
	ProductSourceID int64

	Result struct {
		Removed int
	}
}

type UpdateShopProductCommand struct {
	ShopID          int64
	Product         *catalogmodel.ShopProduct
	Code            string
	ProductSourceID int64

	Result *catalog.ShopProductWithVariants
}

type UpdateShopProductsTagsCommand struct {
	ShopID          int64
	ProductIDs      []int64
	ProductSourceID int64
	Update          *model.UpdateListRequest

	Result struct {
		Updated int
	}
}

type UpdateShopProductImagesCommand struct {
	ShopID    int64
	ProductID int64
	ImageURLs []string

	Result *catalogmodel.ShopProductFtVariant
}

type GetShopCollectionQuery struct {
	ShopID       int64
	CollectionID int64

	Result *catalogmodel.ShopCollection
}

type GetShopCollectionsQuery struct {
	ShopID        int64
	CollectionIDs []int64

	Result struct {
		Collections []*catalogmodel.ShopCollection
	}
}

type CreateShopCollectionCommand struct {
	Collection *catalogmodel.ShopCollection

	Result *catalogmodel.ShopCollection
}

type UpdateShopCollectionCommand struct {
	Collection *catalogmodel.ShopCollection

	Result *catalogmodel.ShopCollection
}

type RemoveShopCollectionCommand struct {
	ShopID       int64
	CollectionID int64

	Result struct {
		Deleted int
	}
}

type CreateVariantCommand struct {
	ShopID          int64
	ProductSourceID int64
	ProductID       int64
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

type CreateProductSourceCommand struct {
	ShopID int64
	Name   string
	Type   string

	Result *catalogmodel.ProductSource
}

type GetShopProductSourcesCommand struct {
	UserID int64
	ShopID int64

	Result []*catalogmodel.ProductSource
}

type CreateProductSourceCategoryCommand struct {
	ShopID            int64
	Name              string
	ProductSourceID   int64
	ProductSourceType string
	ParentID          int64

	Result *catalogmodel.ProductSourceCategory
}

type UpdateProductsProductSourceCategoryCommand struct {
	CategoryID      int64
	ProductIDs      []int64
	ShopID          int64
	ProductSourceID int64

	Result struct {
		Updated int
	}
}

type GetProductSourceCategoryQuery struct {
	ShopID     int64
	CategoryID int64

	Result *catalogmodel.ProductSourceCategory
}

type GetProductSourceCategoriesQuery struct {
	ShopID            int64
	IDs               []int64
	ProductSourceType string

	Result struct {
		Categories []*catalogmodel.ProductSourceCategory
	}
}

type UpdateShopProductSourceCategoryCommand struct {
	ID       int64
	ShopID   int64
	ParentID int64
	Name     string

	Result *catalogmodel.ProductSourceCategory
}

type RemoveShopProductSourceCategoryCommand struct {
	ID     int64
	ShopID int64

	Result struct {
		Removed int
	}
}
