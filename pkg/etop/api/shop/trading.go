package shop

import (
	"context"
	"fmt"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/ordering"
	ordertrading "etop.vn/api/main/ordering/trading"
	"etop.vn/api/meta"
	"etop.vn/api/top/int/shop"
	"etop.vn/api/top/int/types"
	"etop.vn/api/top/types/etc/inventory_auto"
	pbsource "etop.vn/api/top/types/etc/order_source"
	identityconvert "etop.vn/backend/com/main/identity/convert"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/authorize/claims"
	logicorder "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		tradingService.TradingGetProduct,
		tradingService.TradingGetProducts,
		tradingService.TradingCreateOrder,
		tradingService.TradingGetOrder,
		tradingService.TradingGetOrders,
	)
}

func (s *TradingService) TradingGetProduct(ctx context.Context, q *TradingGetProductEndpoint) error {
	query := &catalog.GetShopProductWithVariantsByIDQuery{
		ProductID: q.Id,
		ShopID:    model.EtopTradingAccountID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	result := PbShopProductWithVariants(query.Result)
	result, err := PopulateTradingProductWithInventoryCount(ctx, result)
	if err != nil {
		return err
	}
	q.Result = result
	return nil

}

func (s *TradingService) TradingGetProducts(ctx context.Context, q *TradingGetProductsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  model.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	result := PbShopProductsWithVariants(query.Result.Products)
	result, err := PopulateTradingProductsWithInventoryCount(ctx, result)
	if err != nil {
		return err
	}
	q.Result = &shop.ShopProductsResponse{
		Paging:   cmapi.PbPageInfo(paging),
		Products: result,
	}
	return nil
}

func (s *TradingService) TradingCreateOrder(ctx context.Context, r *TradingCreateOrderEndpoint) error {
	_, err := s.tradingCreateOrder(ctx, r)
	return err
}

func (s *TradingService) tradingCreateOrder(ctx context.Context, r *TradingCreateOrderEndpoint) (_orderID dot.ID, _err error) {
	defer func() {
		if _err == nil {
			return
		}
		if _orderID != 0 {
			// cancel Order an inform error message
			if _, err := logicorder.CancelOrder(ctx, r.Context.UserID, model.EtopTradingAccountID, 0, _orderID, fmt.Sprintf("Tạo đơn Trading thất bại (err = %v)", _err.Error()), inventory_auto.Unknown); err != nil {
				return
			}
		}
	}()
	req := &types.CreateOrderRequest{
		PaymentMethod:   r.PaymentMethod,
		Source:          pbsource.EtopPOS,
		Customer:        r.Customer,
		CustomerAddress: r.CustomerAddress,
		BillingAddress:  r.BillingAddress,
		ShippingAddress: r.ShippingAddress,
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
		ReferralMeta:    r.ReferralMeta,
	}

	query := &identity.GetShopByIDQuery{
		ID: model.EtopTradingAccountID,
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return 0, err
	}
	{
		referralCode := r.ReferralMeta["referral_code"]
		if referralCode != "" {
			tradingOrderCreating := &ordertrading.TradingOrderCreatingEvent{
				EventMeta:    meta.NewEvent(),
				ReferralCode: referralCode,
				UserID:       r.Context.Shop.OwnerID,
			}
			if err := eventBus.Publish(ctx, tradingOrderCreating); err != nil {
				return 0, err
			}
		}
	}
	eTopTrading := identityconvert.ShopDB(query.Result)
	shopClaim := &claims.ShopClaim{Shop: eTopTrading}
	shopID := r.Context.Shop.ID
	resp, err := logicorder.CreateOrder(ctx, shopClaim, nil, req, &shopID, 0)
	if err != nil {
		return 0, err
	}

	{
		_query := &ordering.GetOrderByIDQuery{
			ID: resp.Id,
		}
		if err := orderQuery.Dispatch(ctx, _query); err != nil {
			return resp.Id, err
		}
		tradingOrderCreatedEvent := &ordertrading.TradingOrderCreatedEvent{
			EventMeta:    meta.NewEvent(),
			OrderID:      _query.Result.ID,
			ReferralCode: _query.Result.ReferralMeta.ReferralCode,
		}
		if err := eventBus.Publish(ctx, tradingOrderCreatedEvent); err != nil {
			return resp.Id, err
		}
	}

	r.Result = resp
	return resp.Id, nil
}

func (s *TradingService) TradingGetOrder(ctx context.Context, q *TradingGetOrderEndpoint) error {
	query := &ordermodelx.GetOrderQuery{
		OrderID:            q.Id,
		ShopID:             model.EtopTradingAccountID,
		TradingShopID:      q.Context.Shop.ID,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbOrder(query.Result.Order, nil, model.TagShop)
	q.Result.Fulfillments = convertpb.XPbFulfillments(query.Result.XFulfillments, model.TagShop)
	return nil
}

func (s *TradingService) TradingGetOrders(ctx context.Context, q *TradingGetOrdersEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:       []dot.ID{model.EtopTradingAccountID},
		TradingShopID: q.Context.Shop.ID,
		Paging:        paging,
		Filters:       cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.OrdersResponse{
		Paging: cmapi.PbPageInfo(paging),
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, model.TagShop, query.Result.Shops),
	}
	return nil
}

func PopulateTradingProductsWithInventoryCount(ctx context.Context, args []*shop.ShopProduct) ([]*shop.ShopProduct, error) {
	var variantIDs []dot.ID
	for _, p := range args {
		for _, v := range p.Variants {
			variantIDs = append(variantIDs, v.Id)
		}
	}
	query := &inventory.GetInventoryVariantsByVariantIDsQuery{
		ShopID:     model.EtopTradingAccountID,
		VariantIDs: variantIDs,
	}
	err := inventoryQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	var mapInventoryVariants = make(map[dot.ID]*inventory.InventoryVariant)
	for _, v := range query.Result.InventoryVariants {
		mapInventoryVariants[v.VariantID] = v
	}
	for k1, p := range args {
		for k2, v := range p.Variants {
			if mapInventoryVariants[v.Id] != nil {
				args[k1].Variants[k2].QuantityOnHand = mapInventoryVariants[v.Id].QuantityOnHand
				args[k1].Variants[k2].QuantityPicked = mapInventoryVariants[v.Id].QuantityPicked
				args[k1].Variants[k2].Quantity = mapInventoryVariants[v.Id].QuantitySummary
			}
		}
	}
	return args, nil
}

func PopulateTradingProductWithInventoryCount(ctx context.Context, args *shop.ShopProduct) (*shop.ShopProduct, error) {
	if args == nil {
		return nil, nil
	}
	var variantIDs []dot.ID
	for _, v := range args.Variants {
		variantIDs = append(variantIDs, v.Id)
	}
	query := &inventory.GetInventoryVariantsByVariantIDsQuery{
		ShopID:     model.EtopTradingAccountID,
		VariantIDs: variantIDs,
	}
	err := inventoryQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	var mapInventoryVariants = make(map[dot.ID]*inventory.InventoryVariant)
	for _, v := range query.Result.InventoryVariants {
		mapInventoryVariants[v.VariantID] = v
	}
	for k2, v := range args.Variants {
		if mapInventoryVariants[v.Id] != nil {
			args.Variants[k2].QuantityOnHand = mapInventoryVariants[v.Id].QuantityOnHand
			args.Variants[k2].QuantityPicked = mapInventoryVariants[v.Id].QuantityPicked
			args.Variants[k2].Quantity = mapInventoryVariants[v.Id].QuantitySummary
		}
	}
	return args, nil
}
