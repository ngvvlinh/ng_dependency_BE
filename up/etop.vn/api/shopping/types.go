package shopping

import (
	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

type IDQueryShopArg struct {
	ID             dot.ID
	ShopID         dot.ID
	IncludeDeleted bool
}

type IDsQueryShopArgs struct {
	IDs    []dot.ID
	ShopID dot.ID
	Paging meta.Paging
}

type ListQueryShopArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}
