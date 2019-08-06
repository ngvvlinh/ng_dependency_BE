package shopping

import "etop.vn/api/meta"

type IDQueryShopArg struct {
	ID     int64
	ShopID int64
}

type IDsQueryShopArgs struct {
	IDs    []int64
	ShopID int64
}

type ListQueryShopArgs struct {
	ShopID  int64
	Paging  meta.Paging
	Filters meta.Filters
}
