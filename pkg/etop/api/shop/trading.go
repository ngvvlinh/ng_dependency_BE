package shop

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/identity"
	identityconvert "etop.vn/backend/com/main/identity/convert"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	pbcm "etop.vn/backend/pb/common"
	pborder "etop.vn/backend/pb/etop/order"
	pbsource "etop.vn/backend/pb/etop/order/source"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/etop/authorize/claims"
	logicorder "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/backend/pkg/etop/model"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("api",
		TradingGetProduct,
		TradingGetProducts,
		TradingCreateOrder,
		TradingGetOrder,
		TradingGetOrders,
	)
}

func TradingGetProduct(ctx context.Context, q *wrapshop.TradingGetProductEndpoint) error {
	query := &catalog.GetShopProductByIDQuery{
		ProductID: q.Id,
		ShopID:    model.EtopTradingAccountID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopProduct(query.Result)
	return nil

}

func TradingGetProducts(ctx context.Context, q *wrapshop.TradingGetProductsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &catalog.ListShopProductsQuery{
		ShopID:  model.EtopTradingAccountID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.ShopProductsResponse{
		Paging:   pbcm.PbPageInfo(paging, query.Result.Count),
		Products: PbShopProducts(query.Result.Products),
	}
	return nil
}

func TradingCreateOrder(ctx context.Context, r *wrapshop.TradingCreateOrderEndpoint) error {
	req := &pborder.CreateOrderRequest{
		PaymentMethod:   model.PaymentMethodOther,
		Source:          pbsource.Source_etop_pos,
		Customer:        r.Customer,
		CustomerAddress: r.CustomerAddress,
		BillingAddress:  r.BillingAddress,
		Lines:           r.Lines,
		Discounts:       r.Discounts,
		TotalItems:      r.TotalItems,
		BasketValue:     r.BasketValue,
		OrderDiscount:   r.OrderDiscount,
		TotalFee:        r.TotalFee,
		FeeLines:        r.FeeLines,
		TotalDiscount:   r.TotalDiscount,
		TotalAmount:     r.TotalAmount,
		OrderNote:       r.OrderNote,
	}

	query := &identity.GetShopByIDQuery{
		ID: model.EtopTradingAccountID,
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	eTopTrading := identityconvert.ShopDB(query.Result)
	shopClaim := &claims.ShopClaim{Shop: eTopTrading}
	shopID := r.Context.Shop.ID
	resp, err := logicorder.CreateOrder(ctx, shopClaim, nil, req, &shopID)
	if err != nil {
		return err
	}

	r.Result = resp
	return nil
}

func TradingGetOrder(ctx context.Context, q *wrapshop.TradingGetOrderEndpoint) error {
	query := &ordermodelx.GetOrderQuery{
		OrderID:            q.Id,
		ShopID:             model.EtopTradingAccountID,
		TradingShopID:      q.Context.Shop.ID,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbOrder(query.Result.Order, nil, model.TagShop)
	q.Result.Fulfillments = pborder.XPbFulfillments(query.Result.XFulfillments, model.TagShop)
	return nil
}

func TradingGetOrders(ctx context.Context, q *wrapshop.TradingGetOrdersEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:       []int64{model.EtopTradingAccountID},
		TradingShopID: q.Context.Shop.ID,
		Paging:        paging,
		Filters:       pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.OrdersResponse{
		Paging: pbcm.PbPageInfo(paging, int32(query.Result.Total)),
		Orders: pborder.PbOrdersWithFulfillments(query.Result.Orders, model.TagShop, query.Result.Shops),
	}
	return nil
}
