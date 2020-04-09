package partner

import (
	"context"

	"etop.vn/api/main/connectioning"
	shippingcore "etop.vn/api/main/shipping"
	pbcm "etop.vn/api/top/types/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/apix/shipping"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("apix",
		miscService.GetLocationList,
		shippingService.GetShippingServices,
		shippingService.CreateAndConfirmOrder,
		shippingService.CancelOrder,
		shippingService.GetOrder,
		shippingService.GetFulfillment,
	)
}

func (s *MiscService) GetLocationList(ctx context.Context, r *GetLocationListEndpoint) error {
	resp, err := shipping.GetLocationList(ctx)
	r.Result = resp
	return err
}

func (s *ShippingService) GetShippingServices(ctx context.Context, r *GetShippingServicesEndpoint) error {
	resp, err := shipping.GetShippingServices(ctx, r.Context.Shop.ID, r.GetShippingServicesRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) CreateAndConfirmOrder(ctx context.Context, r *CreateAndConfirmOrderEndpoint) error {
	userID := cm.CoalesceID(r.Context.UserID, r.Context.Shop.OwnerID)
	resp, err := shipping.CreateAndConfirmOrder(ctx, userID, r.Context.Shop.ID, &r.Context, r.CreateAndConfirmOrderRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) CancelOrder(ctx context.Context, r *CancelOrderEndpoint) error {
	userID := cm.CoalesceID(r.Context.UserID, r.Context.Shop.OwnerID)
	resp, err := shipping.CancelOrder(ctx, userID, r.Context.Shop.ID, r.CancelOrderRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) GetOrder(ctx context.Context, r *GetOrderEndpoint) error {
	resp, err := shipping.GetOrder(ctx, r.Context.Shop.ID, r.OrderIDRequest)
	r.Result = resp
	return err
}

func (s *ShippingService) GetFulfillment(ctx context.Context, r *GetFulfillmentEndpoint) error {
	resp, err := shipping.GetFulfillment(ctx, r.Context.Shop.ID, r.FulfillmentIDRequest)
	r.Result = resp
	return err
}

/*
	UpdateFulfillment

	Api này chỉ sử dụng cho partner là nhà vận chuyển
	Chỉ update fulfillment thuộc connection của NVC
*/

func (s *ShipmentService) UpdateFulfillment(ctx context.Context, r *UpdateFulfillmentEndpoint) error {
	if r.ShippingCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing shipping_code")
	}
	if !r.ShippingState.Valid {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing shipping_state")
	}

	query := &connectioning.ListConnectionsQuery{
		PartnerID: r.Context.Partner.ID,
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	connIDs := []dot.ID{}
	for _, conn := range query.Result {
		connIDs = append(connIDs, conn.ID)
	}
	if len(connIDs) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Không thể sử dụng api này. Vui lòng liên hệ %v để biết thêm chi tiết.", wl.X(ctx).CSEmail)
	}

	cmd := &shippingcore.UpdateFulfillmentShippingStateCommand{
		ShippingCode:  r.ShippingCode,
		ShippingState: r.ShippingState.Enum,
		ConnectionIDs: connIDs,
	}
	if err := shippingAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}
