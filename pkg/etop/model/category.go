package model

type GetProductSourceCategoryQuery struct {
	SupplierID int64
	ShopID     int64
	CategoryID int64

	Result *ProductSourceCategoryExtended
}

type GetProductSourceCategoriesExtendedQuery struct {
	SupplierID        int64
	ShopID            int64
	IDs               []int64
	ProductSourceType string

	Result struct {
		Categories []*ProductSourceCategoryExtended
	}
}

type GetProductSourceCategoriesQuery struct {
	SupplierID        int64
	ShopID            int64
	IDs               []int64
	ProductSourceType string

	Result struct {
		Categories []*ProductSourceCategory
	}
}

type GetEtopCategoryQuery struct {
	CategoryID int64
	Status     *Status3

	Result *EtopCategory
}

type GetEtopCategoriesQuery struct {
	Status *Status3

	Result struct {
		Categories []*EtopCategory
	}
}

type CreateEtopCategoryCommand struct {
	Category *EtopCategory

	Result *EtopCategory
}

type UpdateShopProductSourceCategoryCommand struct {
	ID       int64
	ShopID   int64
	ParentID int64
	Name     string

	Result *ProductSourceCategoryExtended
}

type RemoveShopProductSourceCategoryCommand struct {
	ID     int64
	ShopID int64

	Result struct {
		Removed int
	}
}
