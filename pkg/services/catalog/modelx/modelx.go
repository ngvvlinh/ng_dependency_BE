package modelx

import (
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

type GetProductQuery struct {
	SupplierID int64
	ProductID  int64
	model.StatusQuery

	Result *catalogmodel.ProductFtVariant
}

type GetVariantQuery struct {
	SupplierID int64
	VariantID  int64
	model.StatusQuery

	Result *catalogmodel.VariantExtended
}

type GetProductsQuery struct {
	ProductSourceID int64
	IncludeDeleted  bool
	ExcludeEdCode   bool

	EdCodes []string

	// must be normalized names
	NameNormUas []string

	Result struct {
		Products []*catalogmodel.Product
	}
}

type GetVariantsQuery struct {
	ProductSourceID int64
	IncludeDeleted  bool
	Inclusive       bool // Include both ed_code and attr_norm_kv

	EdCodes []string

	// must be group of (product_id, attr_norm_kv)
	// the default variant's attr_norm_kv is '_'
	AttrNorms []interface{}

	Result struct {
		Variants []*catalogmodel.Variant
	}
}

type GetProductsExtendedQuery struct {
	SupplierID int64

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

type GetVariantExternalsQuery struct {
	SupplierID int64

	Paging  *cm.Paging
	Filters []cm.Filter
	IDs     []int64
	model.StatusQuery

	Result struct {
		Variants []*catalogmodel.VariantExternalExtended
		Total    int
	}
}

type GetVariantExternalsFromIDQuery struct {
	FromID int64
	Limit  int

	Result ScanVariantExternalsResult
}

type ScanVariantExternalsQuery struct {
	FromID   int64
	Limit    int
	PageSize int

	Result chan ScanVariantExternalsResult
}

type ScanVariantExternalsResult struct {
	MaxID    int64
	Variants []*catalogmodel.VariantExternalExtended
}

type UpdateProductCommand struct {
	SupplierID int64

	Product *catalogmodel.Product

	Result *catalogmodel.ProductFtVariant
}

type UpdateVariantCommand struct {
	SupplierID int64

	Variant *catalogmodel.Variant

	Result *catalogmodel.VariantExtended
}

type UpdateProductImagesCommand struct {
	SupplierID int64
	ProductID  int64
	ImageURLs  []string

	Result *catalogmodel.ProductFtVariant
}

type UpdateVariantImagesCommand struct {
	SupplierID int64
	VariantID  int64
	ImageURLs  []string

	Result *catalogmodel.VariantExtended
}

type UpdateVariantPriceCommand struct {
	VariantID int64
	PriceDef  *catalogmodel.PriceDef
}

type UpdateProductsCommand struct {
	SupplierID int64
	Products   []*catalogmodel.Product

	Result struct {
		Products []*catalogmodel.Product
		Errors   []error
	}
}

type UpdateVariantsCommand struct {
	SupplierID int64
	Variants   []*catalogmodel.Variant

	Result struct {
		Variants []*catalogmodel.VariantExtended
		Errors   []error
	}
}

type UpdateVariantsStatusCommand struct {
	SupplierID  int64
	IDs         []int64
	StatusQuery model.StatusQuery
	Update      model.ProductStatusUpdate

	Result struct {
		Updated int
	}
}

type UpdateProductsStatusCommand struct {
	SupplierID  int64
	IDs         []int64
	StatusQuery model.StatusQuery
	Update      model.ProductStatusUpdate

	Result struct {
		Updated int
	}
}

type UpdateProductsEtopCategoryCommand struct {
	ProductIDs     []int64
	EtopCategoryID int64

	Result struct {
		Updated int
	}
}

type RemoveProductsEtopCategoryCommand struct {
	ProductIDs []int64

	Result struct {
		Updated int
	}
}

type GetShopVariantQuery struct {
	ShopID    int64
	VariantID int64

	ShopVariantStatus *int

	Result *catalogmodel.ShopVariantExtended
}

type GetShopVariantsQuery struct {
	ShopID     int64
	Paging     *cm.Paging
	Filters    []cm.Filter
	VariantIDs []int64

	ShopVariantStatus *int

	Result struct {
		Total    int
		Variants []*catalogmodel.ShopVariantExtended
	}
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
		Variants []*catalogmodel.ShopVariantExtended
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
	CostPrice       int
	Inventory       int
	EdCode          string
	Attributes      []catalogmodel.ProductAttribute
	ProductSourceID int64

	Result *catalogmodel.ShopVariantExtended
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

type UpdateShopVariantsStatusCommand struct {
	ShopID          int64
	VariantIDs      []int64
	ProductSourceID int64
	Update          struct {
		Status *model.Status3
	}

	Result struct {
		Updated int
	}
}

type UpdateShopVariantsTagsCommand struct {
	ShopID          int64
	VariantIDs      []int64
	Update          *model.UpdateListRequest
	ProductSourceID int64

	Result struct {
		Updated int
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

type GetShopProductQuery struct {
	ShopID          int64
	ProductID       int64
	ProductSourceID int64

	ShopProductStatus *model.Status3

	Result *catalogmodel.ShopProductFtVariant
}

type GetShopProductsQuery struct {
	ShopID          int64
	Paging          *cm.Paging
	Filters         []cm.Filter
	ProductIDs      []int64
	ProductSourceID int64

	ShopProductStatus *model.Status3

	Result struct {
		Total    int
		Products []*catalogmodel.ShopProductFtVariant
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

	Result *catalogmodel.ShopProductFtVariant
}

type UpdateShopProductsStatusCommand struct {
	ShopID          int64
	ProductIDs      []int64
	ProductSourceID int64
	Update          struct {
		Status *model.Status3
	}

	Result struct {
		Updated int
	}
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

type GetAllShopVariantsQuery struct {
	ShopID          int64
	VariantIDs      []int64
	ProductSourceID int64

	Result struct {
		Variants []*catalogmodel.ShopVariantExtended
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
	ListPrice         int
	SKU               string
	Code              string
	QuantityAvailable int
	QuantityOnHand    int
	QuantityReserved  int
	CostPrice         int
	Attributes        []catalogmodel.ProductAttribute
	DescHTML          string

	Result *catalogmodel.ShopProductFtVariant
}
