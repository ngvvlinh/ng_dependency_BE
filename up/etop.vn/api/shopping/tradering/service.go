package tradering

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
)

// +gen:api

type QueryService interface {
	GetTraderByID(ctx context.Context, _ *shopping.IDQueryShopArg) (*ShopTrader, error)
	GetTraderInfoByID(ctx context.Context, ID, ShopID int64) (*ShopTrader, error)
	ListTradersByIDs(context.Context, *shopping.IDsQueryShopArgs) (*TradersResponse, error)
}

//-- queries --//

type TradersResponse struct {
	Traders []*ShopTrader
	Count   int32
	Paging  meta.PageInfo
}
