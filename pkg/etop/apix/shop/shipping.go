package xshop

import (
	"context"

	"etop.vn/backend/pkg/etop/apix/shipping"
	wrapxshop "etop.vn/backend/wrapper/external/shop"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("apix",
		GetLocationList,
		GetShippingServices,
		CreateAndConfirmOrder,
		CancelOrder,
		GetOrder,
		GetFulfillment,
	)
}

func GetLocationList(ctx context.Context, r *wrapxshop.GetLocationListEndpoint) error {
	resp, err := shipping.GetLocationList(ctx)
	r.Result = resp
	return err
}

func GetShippingServices(ctx context.Context, r *wrapxshop.GetShippingServicesEndpoint) error {
	resp, err := shipping.GetShippingServices(ctx, r.Context.Shop.ID, r.GetShippingServicesRequest)
	r.Result = resp
	return err
}

func CreateAndConfirmOrder(ctx context.Context, r *wrapxshop.CreateAndConfirmOrderEndpoint) error {
	resp, err := shipping.CreateAndConfirmOrder(ctx, r.Context.Shop.ID, &r.Context, r.CreateOrderRequest)
	r.Result = resp
	return err
}

func CancelOrder(ctx context.Context, r *wrapxshop.CancelOrderEndpoint) error {
	resp, err := shipping.CancelOrder(ctx, r.Context.Shop.ID, r.CancelOrderRequest)
	r.Result = resp
	return err
}

func GetOrder(ctx context.Context, r *wrapxshop.GetOrderEndpoint) error {
	resp, err := shipping.GetOrder(ctx, r.Context.Shop.ID, r.OrderIDRequest)
	r.Result = resp
	return err
}

func GetFulfillment(ctx context.Context, r *wrapxshop.GetFulfillmentEndpoint) error {
	resp, err := shipping.GetFulfillment(ctx, r.Context.Shop.ID, r.FulfillmentIDRequest)
	r.Result = resp
	return err
}
