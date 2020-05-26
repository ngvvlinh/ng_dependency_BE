package partner

import (
	"context"

	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/inventory_policy"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/shipping"
)

type OrderService struct {
	Shipping *shipping.Shipping
}

func (s *OrderService) Clone() *OrderService { res := *s; return &res }

func (s *OrderService) GetOrder(ctx context.Context, r *OrderGetOrderEndpoint) error {
	resp, err := s.Shipping.GetOrder(ctx, r.Context.Shop.ID, r.OrderIDRequest)
	r.Result = resp
	return err
}

func (s *OrderService) ListOrders(ctx context.Context, r *OrderListOrdersEndpoint) error {
	return cm.ErrTODO
}

func (s *OrderService) CreateOrder(ctx context.Context, r *OrderCreateOrderEndpoint) error {
	resp, err := s.Shipping.CreateOrder(ctx, &r.Context, r.CreateOrderRequest)
	r.Result = resp
	return err
}

func (s *OrderService) ConfirmOrder(ctx context.Context, r *OrderConfirmOrderEndpoint) error {
	autoInventoryVoucher := inventory_auto.Unknown
	if r.InventoryPolicy == inventory_policy.Obey {
		autoInventoryVoucher = r.AutoInventoryVoucher
	}
	err := s.Shipping.ConfirmOrder(ctx, r.Context.UserID, &r.Context, r.OrderId, autoInventoryVoucher)
	r.Result = &common.Empty{}
	return err
}

func (s *OrderService) CancelOrder(ctx context.Context, r *OrderCancelOrderEndpoint) error {
	_, err := s.Shipping.CancelOrder(ctx, r.Context.UserID, r.Context.Shop.ID, r.CancelOrderRequest)
	r.Result = &common.Empty{}
	return err
}
