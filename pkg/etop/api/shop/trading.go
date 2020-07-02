package shop

import (
	"context"
	"fmt"

	"o.o/api/main/catalog"
	"o.o/api/main/identity"
	"o.o/api/main/inventory"
	"o.o/api/main/ordering"
	ordertrading "o.o/api/main/ordering/trading"
	"o.o/api/meta"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	inttypes "o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/inventory_auto"
	pbsource "o.o/api/top/types/etc/order_source"
	identityconvert "o.o/backend/com/main/identity/convert"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	"o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
)

type TradingService struct {
	session.Session

	EventBus       capi.EventBus
	IdentityQuery  identity.QueryBus
	CatalogQuery   catalog.QueryBus
	OrderQuery     ordering.QueryBus
	InventoryQuery inventory.QueryBus
	OrderLogic     *logicorder.OrderLogic
}

func (s *TradingService) Clone() api.TradingService { res := *s; return &res }

func (s *TradingService) TradingGetProduct(ctx context.Context, q *pbcm.IDRequest) (*api.ShopProduct, error) {
	query := &catalog.GetShopProductWithVariantsByIDQuery{
		ProductID: q.Id,
		ShopID:    model.EtopTradingAccountID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := PbShopProductWithVariants(query.Result)
	result, err := PopulateTradingProductWithInventoryCount(ctx, s.InventoryQuery, result)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (s *TradingService) TradingGetProducts(ctx context.Context, q *pbcm.CommonListRequest) (*api.ShopProductsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  model.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	resp := PbShopProductsWithVariants(query.Result.Products)
	resp, err := s.populateTradingProductsWithInventoryCount(ctx, resp)
	if err != nil {
		return nil, err
	}
	result := &api.ShopProductsResponse{
		Paging:   cmapi.PbPageInfo(paging),
		Products: resp,
	}
	return result, nil
}

func (s *TradingService) TradingCreateOrder(ctx context.Context, r *inttypes.TradingCreateOrderRequest) (*inttypes.Order, error) {
	order, _, err := s.tradingCreateOrder(ctx, r)
	return order, err
}

func (s *TradingService) tradingCreateOrder(ctx context.Context, r *inttypes.TradingCreateOrderRequest) (_ *inttypes.Order, _orderID dot.ID, _err error) {
	defer func() {
		if _err == nil {
			return
		}
		if _orderID != 0 {
			// cancel Order an inform error message
			if _, err := s.OrderLogic.CancelOrder(ctx, s.SS.Claim().UserID, model.EtopTradingAccountID, 0, _orderID, fmt.Sprintf("Tạo đơn Trading thất bại (err = %v)", _err.Error()), inventory_auto.Unknown); err != nil {
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
	if err := s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, 0, err
	}
	{
		referralCode := r.ReferralMeta["referral_code"]
		if referralCode != "" {
			tradingOrderCreating := &ordertrading.TradingOrderCreatingEvent{
				EventMeta:    meta.NewEvent(),
				ReferralCode: referralCode,
				UserID:       s.SS.Shop().OwnerID,
			}
			if err := s.EventBus.Publish(ctx, tradingOrderCreating); err != nil {
				return nil, 0, err
			}
		}
	}
	eTopTrading := identityconvert.ShopDB(query.Result)
	shopID := s.SS.Shop().ID
	resp, err := s.OrderLogic.CreateOrder(ctx, eTopTrading, nil, req, &shopID, 0)
	if err != nil {
		return nil, 0, err
	}

	{
		_query := &ordering.GetOrderByIDQuery{
			ID: resp.Id,
		}
		if err := s.OrderQuery.Dispatch(ctx, _query); err != nil {
			return nil, resp.Id, err
		}
		tradingOrderCreatedEvent := &ordertrading.TradingOrderCreatedEvent{
			EventMeta:    meta.NewEvent(),
			OrderID:      _query.Result.ID,
			ReferralCode: _query.Result.ReferralMeta.ReferralCode,
		}
		if err := s.EventBus.Publish(ctx, tradingOrderCreatedEvent); err != nil {
			return nil, resp.Id, err
		}
	}
	return resp, resp.Id, nil
}

func (s *TradingService) TradingGetOrder(ctx context.Context, q *pbcm.IDRequest) (*inttypes.Order, error) {
	query := &ordermodelx.GetOrderQuery{
		OrderID:            q.Id,
		ShopID:             model.EtopTradingAccountID,
		TradingShopID:      s.SS.Shop().ID,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbOrder(query.Result.Order, nil, model.TagShop)
	result.Fulfillments = convertpb.XPbFulfillments(query.Result.XFulfillments, model.TagShop)
	return result, nil
}

func (s *TradingService) TradingGetOrders(ctx context.Context, q *api.GetOrdersRequest) (*inttypes.OrdersResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &ordermodelx.GetOrdersQuery{
		ShopIDs:       []dot.ID{model.EtopTradingAccountID},
		TradingShopID: s.SS.Shop().ID,
		Paging:        paging,
		Filters:       cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &inttypes.OrdersResponse{
		Paging: cmapi.PbPageInfo(paging),
		Orders: convertpb.PbOrdersWithFulfillments(query.Result.Orders, model.TagShop, query.Result.Shops),
	}
	return result, nil
}

func (s *TradingService) populateTradingProductsWithInventoryCount(ctx context.Context, args []*shop.ShopProduct) ([]*shop.ShopProduct, error) {
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
	err := s.InventoryQuery.Dispatch(ctx, query)
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

func PopulateTradingProductWithInventoryCount(ctx context.Context, inventoryQuery inventory.QueryBus, args *shop.ShopProduct) (*shop.ShopProduct, error) {
	if args == nil {
		return nil, nil
	}
	var variantIDs []dot.ID
	for _, v := range args.Variants {
		variantIDs = append(variantIDs, v.Id)
	}
	if len(variantIDs) == 0 {
		return args, nil
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
