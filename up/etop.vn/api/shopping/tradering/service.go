package tradering

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	DeleteTrader(ctx context.Context, ID dot.ID, shopID dot.ID) (deleted int, _ error)
}

type QueryService interface {
	GetTraderByID(ctx context.Context, _ *shopping.IDQueryShopArg) (*ShopTrader, error)
	GetTraderInfoByID(ctx context.Context, ID, ShopID dot.ID) (*ShopTrader, error)
	ListTradersByIDs(context.Context, *shopping.IDsQueryShopArgs) (*TradersResponse, error)
}

//-- queries --//

type TradersResponse struct {
	Traders []*ShopTrader
	Count   int
	Paging  meta.PageInfo
}
