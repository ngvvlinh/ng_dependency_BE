package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/inventory_policy"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/shipping"
	"o.o/backend/pkg/etop/authorize/session"
)

type OrderService struct {
	session.Session

	Shipping *shipping.Shipping
}

func (s *OrderService) Clone() api.OrderService { res := *s; return &res }

func (s *OrderService) GetOrder(ctx context.Context, r *externaltypes.OrderIDRequest) (*externaltypes.OrderAndFulfillments, error) {
	resp, err := s.Shipping.GetOrder(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *OrderService) ListOrders(ctx context.Context, r *externaltypes.ListOrdersRequest) (*externaltypes.OrdersResponse, error) {
	return nil, cm.ErrTODO
}

func (s *OrderService) CreateOrder(ctx context.Context, r *externaltypes.CreateOrderRequest) (*externaltypes.OrderWithoutShipping, error) {
	resp, err := s.Shipping.CreateOrder(ctx, s.SS.Shop(), s.SS.Partner(), r)
	return resp, err
}

func (s *OrderService) ConfirmOrder(ctx context.Context, r *externaltypes.ConfirmOrderRequest) (*pbcm.Empty, error) {
	autoInventoryVoucher := inventory_auto.Unknown
	if r.InventoryPolicy == inventory_policy.Obey {
		autoInventoryVoucher = r.AutoInventoryVoucher
	}
	err := s.Shipping.ConfirmOrder(ctx, s.SS.Claim().UserID, s.SS.Shop(), s.SS.Partner(), r.OrderId, autoInventoryVoucher)
	return &pbcm.Empty{}, err
}

func (s *OrderService) CancelOrder(ctx context.Context, r *externaltypes.CancelOrderRequest) (*pbcm.Empty, error) {
	_, err := s.Shipping.CancelOrder(ctx, s.SS.Claim().UserID, s.SS.Shop().ID, r)
	return &pbcm.Empty{}, err
}
