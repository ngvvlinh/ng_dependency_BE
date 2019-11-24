package shopping

import (
	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

type IDQueryShopArg struct {
	ID     dot.ID
	ShopID dot.ID
}

type IDsQueryShopArgs struct {
	IDs    []dot.ID
	ShopID dot.ID
}

type ListQueryShopArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}
