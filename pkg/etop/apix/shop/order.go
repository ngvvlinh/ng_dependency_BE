package xshop

import (
	"context"

	"etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/inventory_auto"
	"etop.vn/api/top/types/etc/inventory_policy"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/backend/pkg/etop/apix/shipping"
)

func (s *OrderService) CancelOrder(ctx context.Context, r *OrderCancelOrderEndpoint) error {
	resp, err := shipping.CancelOrder(ctx, r.Context.User.ID, r.Context.Shop.ID, r.CancelOrderRequest)
	r.Result = resp
	return err
}

func (s *OrderService) GetOrder(ctx context.Context, r *OrderGetOrderEndpoint) error {
	resp, err := shipping.GetOrder(ctx, r.Context.Shop.ID, r.OrderIDRequest)
	r.Result = resp
	return err
}

func (s *OrderService) ListOrders(ctx context.Context, r *OrderListOrdersEndpoint) error {
	panic("implement me")
}

func (s *OrderService) CreateOrder(ctx context.Context, r *OrderCreateOrderEndpoint) error {
	resp, err := shipping.CreateOrder(ctx, &r.Context, r.CreateOrderRequest)
	r.Result = convertpb.PbOrderToOrderWithoutShipping(resp)
	return err
}

func (s *OrderService) ConfirmOrder(ctx context.Context, r *OrderConfirmOrderEndpoint) error {
	autoInventoryVoucher := inventory_auto.Unknown
	if r.InventoryPolicy == inventory_policy.Obey {
		autoInventoryVoucher = r.AutoInventoryVoucher
	}
	err := shipping.ConfirmOrder(ctx, r.Context.UserID, &r.Context, r.OrderId, autoInventoryVoucher)
	r.Result = &common.Empty{}
	return err
}
