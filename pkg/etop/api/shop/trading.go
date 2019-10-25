package shop

import (
	"context"
	"fmt"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/ordering"
	ordertrading "etop.vn/api/main/ordering/trading"
	"etop.vn/api/meta"
	identityconvert "etop.vn/backend/com/main/identity/convert"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	pbcm "etop.vn/backend/pb/common"
	pborder "etop.vn/backend/pb/etop/order"
	pbsource "etop.vn/backend/pb/etop/order/source"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/authorize/claims"
	logicorder "etop.vn/backend/pkg/etop/logic/orders"
	"etop.vn/backend/pkg/etop/model"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
)

func init() {
	bus.AddHandlers("api",
		s.TradingGetProduct,
		s.TradingGetProducts,
		s.TradingCreateOrder,
		s.TradingGetOrder,
		s.TradingGetOrders,
	)
}

func (s *Service) TradingGetProduct(ctx context.Context, q *wrapshop.TradingGetProductEndpoint) error {
	query := &catalog.GetShopProductWithVariantsByIDQuery{
		ProductID: q.Id,
		ShopID:    model.EtopTradingAccountID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(query.Result)
	return nil

}

func (s *Service) TradingGetProducts(ctx context.Context, q *wrapshop.TradingGetProductsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  model.EtopTradingAccountID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.ShopProductsResponse{
		Paging:   pbcm.PbPageInfo(paging, query.Result.Count),
		Products: PbShopProductsWithVariants(query.Result.Products),
	}
	return nil
}

func (s *Service) TradingCreateOrder(ctx context.Context, r *wrapshop.TradingCreateOrderEndpoint) error {
	_, err := s.tradingCreateOrder(ctx, r)
	return err
}

func (s *Service) tradingCreateOrder(ctx context.Context, r *wrapshop.TradingCreateOrderEndpoint) (_orderID int64, _err error) {
	defer func() {
		if _err == nil {
			return
		}
		if _orderID != 0 {
			// cancel Order an inform error message
			if _, err := logicorder.CancelOrder(ctx, model.EtopTradingAccountID, 0, _orderID, fmt.Sprintf("Tạo đơn Trading thất bại (err = %v)", _err.Error())); err != nil {
				return
			}
		}
	}()
	req := &pborder.CreateOrderRequest{
		PaymentMethod:   r.PaymentMethod,
		Source:          pbsource.Source_etop_pos,
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
	resp, err := logicorder.CreateOrder(ctx, shopClaim, nil, req, &shopID)
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

func (s *Service) TradingGetOrder(ctx context.Context, q *wrapshop.TradingGetOrderEndpoint) error {
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

func (s *Service) TradingGetOrders(ctx context.Context, q *wrapshop.TradingGetOrdersEndpoint) error {
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
