package model

import (
	cm "etop.vn/backend/pkg/common"
)

//-- Query

type GetProductQuery struct {
	SupplierID int64
	ProductID  int64
	StatusQuery

	Result *ProductFtVariant
}

type GetVariantQuery struct {
	SupplierID int64
	VariantID  int64
	StatusQuery

	Result *VariantExtended
}

type GetProductsQuery struct {
	ProductSourceID int64
	IncludeDeleted  bool
	ExcludeEdCode   bool

	EdCodes []string

	// must be normalized names
	NameNormUas []string

	Result struct {
		Products []*Product
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
		Variants []*Variant
	}
}

type GetProductsExtendedQuery struct {
	SupplierID int64

	Paging  *cm.Paging
	Filters []cm.Filter
	IDs     []int64
	StatusQuery
	ProductSourceType string

	Result struct {
		Products []*ProductFtVariant
		Total    int
	}
}

type GetVariantsExtendedQuery struct {
	SupplierID      int64
	ProductSourceID int64

	Paging  *cm.Paging
	Filters []cm.Filter
	IDs     []int64
	Codes   []string
	EdCodes []string
	StatusQuery

	SkipPaging bool

	Result struct {
		Variants []*VariantExtended
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
	StatusQuery

	Result struct {
		Variants []*VariantExternalExtended
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
	Variants []*VariantExternalExtended
}

type UpdateProductCommand struct {
	SupplierID int64

	Product *Product

	Result *ProductFtVariant
}

type UpdateVariantCommand struct {
	SupplierID int64

	Variant *Variant

	Result *VariantExtended
}

type UpdateProductImagesCommand struct {
	SupplierID int64
	ProductID  int64
	ImageURLs  []string

	Result *ProductFtVariant
}

type UpdateVariantImagesCommand struct {
	SupplierID int64
	VariantID  int64
	ImageURLs  []string

	Result *VariantExtended
}

type UpdateVariantPriceCommand struct {
	VariantID int64
	PriceDef  *PriceDef
}

type UpdateProductsCommand struct {
	SupplierID int64
	Products   []*Product

	Result struct {
		Products []*Product
		Errors   []error
	}
}

type UpdateVariantsCommand struct {
	SupplierID int64
	Variants   []*Variant

	Result struct {
		Variants []*VariantExtended
		Errors   []error
	}
}

type UpdateVariantsStatusCommand struct {
	SupplierID  int64
	IDs         []int64
	StatusQuery StatusQuery
	Update      ProductStatusUpdate

	Result struct {
		Updated int
	}
}

type UpdateProductsStatusCommand struct {
	SupplierID  int64
	IDs         []int64
	StatusQuery StatusQuery
	Update      ProductStatusUpdate

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

	Result *ShopVariantExtended
}

type GetShopVariantsQuery struct {
	ShopID     int64
	Paging     *cm.Paging
	Filters    []cm.Filter
	VariantIDs []int64

	ShopVariantStatus *int

	Result struct {
		Total    int
		Variants []*ShopVariantExtended
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
		Variants []*ShopVariantExtended
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
	Variant         *ShopVariant
	CostPrice       int
	Inventory       int
	EdCode          string
	Attributes      []ProductAttribute
	ProductSourceID int64

	Result *ShopVariantExtended
}

type UpdateShopVariantsCommand struct {
	ShopID          int64
	Products        []*ShopVariant
	ProductSourceID int64

	Result struct {
		Products []*ShopVariantExtended
		Errors   []error
	}
}

type UpdateShopVariantsStatusCommand struct {
	ShopID          int64
	VariantIDs      []int64
	ProductSourceID int64
	Update          struct {
		Status *Status3
	}

	Result struct {
		Updated int
	}
}

type UpdateShopVariantsTagsCommand struct {
	ShopID          int64
	VariantIDs      []int64
	Update          *UpdateListRequest
	ProductSourceID int64

	Result struct {
		Updated int
	}
}

type AddShopProductsCommand struct {
	ShopID int64
	IDs    []int64

	Result struct {
		Products []*ShopProduct
		Errors   []error
	}
}

type GetShopProductQuery struct {
	ShopID          int64
	ProductID       int64
	ProductSourceID int64

	ShopProductStatus *Status3

	Result *ShopProductFtVariant
}

type GetShopProductsQuery struct {
	ShopID          int64
	Paging          *cm.Paging
	Filters         []cm.Filter
	ProductIDs      []int64
	ProductSourceID int64

	ShopProductStatus *Status3

	Result struct {
		Total    int
		Products []*ShopProductFtVariant
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
	Product         *ShopProduct
	Code            string
	ProductSourceID int64

	Result *ShopProductFtVariant
}

type UpdateShopProductsStatusCommand struct {
	ShopID          int64
	ProductIDs      []int64
	ProductSourceID int64
	Update          struct {
		Status *Status3
	}

	Result struct {
		Updated int
	}
}

type UpdateShopProductsTagsCommand struct {
	ShopID          int64
	ProductIDs      []int64
	ProductSourceID int64
	Update          *UpdateListRequest

	Result struct {
		Updated int
	}
}

type GetAllShopVariantsQuery struct {
	ShopID          int64
	VariantIDs      []int64
	ProductSourceID int64

	Result struct {
		Variants []*ShopVariantExtended
	}
}

type UpdateShopProductImagesCommand struct {
	ShopID    int64
	ProductID int64
	ImageURLs []string

	Result *ShopProductFtVariant
}
