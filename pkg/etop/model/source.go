package model

//-- Query & Command

type GetProductSourceProps struct {
	ID          int64
	Type        string
	ExternalKey string
}

type GetProductSourceQuery struct {
	GetProductSourceProps
	Result *ProductSource
}

type GetProductSourceExtendedQuery struct {
	GetProductSourceProps
	Result *ProductSourceExtended
}

type GetAllProductSourcesQuery struct {
	External *bool
	Supplier *bool

	Result struct {
		Sources []*ProductSource
	}
}

type GetOrderSourceProps struct {
	ID          int64
	ShopID      int64
	Type        string
	ExternalKey string
}

type GetOrderSourceQuery struct {
	GetOrderSourceProps
	Result *OrderSource
}

type GetOrderSourceExtendedQuery struct {
	GetOrderSourceProps
	Result *OrderSourceExtended
}
