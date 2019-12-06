// +build !generator

// Code generated by generator api. DO NOT EDIT.

package gateway

import (
	context "context"

	haravan "etop.vn/api/external/haravan"
	capi "etop.vn/capi"
	dot "etop.vn/capi/dot"
)

type CommandBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus { return CommandBus{bus} }

func (b CommandBus) Dispatch(ctx context.Context, msg interface{ command() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type CancelOrderCommand struct {
	EtopShopID     dot.ID
	TrackingNumber string `json:"tracking_number"`

	Result *GetOrderResponse `json:"-"`
}

func (h AggregateHandler) HandleCancelOrder(ctx context.Context, msg *CancelOrderCommand) (err error) {
	msg.Result, err = h.inner.CancelOrder(msg.GetArgs(ctx))
	return err
}

type CreateOrderCommand struct {
	EtopShopID            dot.ID
	Origin                *haravan.Address `json:"origin"`
	Destination           *haravan.Address `json:"destination"`
	Items                 []*haravan.Item  `json:"items"`
	CodAmount             float32          `json:"cod_amount"`
	TotalGrams            float32          `json:"total_grams"`
	ExternalStoreID       int              `json:"external_store_id"`
	ExternalOrderID       int              `json:"external_order_id"`
	ExternalFulfillmentID int              `json:"external_fulfillment_id"`
	ExternalCode          string           `json:"external_code"`
	Note                  string           `json:"note"`
	ShippingRateID        int              `json:"shipping_rate_id"`

	Result *CreateOrderResponse `json:"-"`
}

func (h AggregateHandler) HandleCreateOrder(ctx context.Context, msg *CreateOrderCommand) (err error) {
	msg.Result, err = h.inner.CreateOrder(msg.GetArgs(ctx))
	return err
}

type GetOrderCommand struct {
	EtopShopID     dot.ID
	TrackingNumber string `json:"tracking_number"`

	Result *GetOrderResponse `json:"-"`
}

func (h AggregateHandler) HandleGetOrder(ctx context.Context, msg *GetOrderCommand) (err error) {
	msg.Result, err = h.inner.GetOrder(msg.GetArgs(ctx))
	return err
}

type GetShippingRateCommand struct {
	EtopShopID  dot.ID
	Origin      *haravan.Address `json:"origin"`
	Destination *haravan.Address `json:"destination"`
	CodAmount   float32          `json:"cod_amount"`
	TotalGrams  float32          `json:"total_grams"`

	Result *GetShippingRateResponse `json:"-"`
}

func (h AggregateHandler) HandleGetShippingRate(ctx context.Context, msg *GetShippingRateCommand) (err error) {
	msg.Result, err = h.inner.GetShippingRate(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CancelOrderCommand) command()     {}
func (q *CreateOrderCommand) command()     {}
func (q *GetOrderCommand) command()        {}
func (q *GetShippingRateCommand) command() {}

// implement conversion

func (q *CancelOrderCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CancelOrderRequestArgs) {
	return ctx,
		&CancelOrderRequestArgs{
			EtopShopID:     q.EtopShopID,
			TrackingNumber: q.TrackingNumber,
		}
}

func (q *CancelOrderCommand) SetCancelOrderRequestArgs(args *CancelOrderRequestArgs) {
	q.EtopShopID = args.EtopShopID
	q.TrackingNumber = args.TrackingNumber
}

func (q *CreateOrderCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrderRequestArgs) {
	return ctx,
		&CreateOrderRequestArgs{
			EtopShopID:            q.EtopShopID,
			Origin:                q.Origin,
			Destination:           q.Destination,
			Items:                 q.Items,
			CodAmount:             q.CodAmount,
			TotalGrams:            q.TotalGrams,
			ExternalStoreID:       q.ExternalStoreID,
			ExternalOrderID:       q.ExternalOrderID,
			ExternalFulfillmentID: q.ExternalFulfillmentID,
			ExternalCode:          q.ExternalCode,
			Note:                  q.Note,
			ShippingRateID:        q.ShippingRateID,
		}
}

func (q *CreateOrderCommand) SetCreateOrderRequestArgs(args *CreateOrderRequestArgs) {
	q.EtopShopID = args.EtopShopID
	q.Origin = args.Origin
	q.Destination = args.Destination
	q.Items = args.Items
	q.CodAmount = args.CodAmount
	q.TotalGrams = args.TotalGrams
	q.ExternalStoreID = args.ExternalStoreID
	q.ExternalOrderID = args.ExternalOrderID
	q.ExternalFulfillmentID = args.ExternalFulfillmentID
	q.ExternalCode = args.ExternalCode
	q.Note = args.Note
	q.ShippingRateID = args.ShippingRateID
}

func (q *GetOrderCommand) GetArgs(ctx context.Context) (_ context.Context, _ *GetOrderRequestArgs) {
	return ctx,
		&GetOrderRequestArgs{
			EtopShopID:     q.EtopShopID,
			TrackingNumber: q.TrackingNumber,
		}
}

func (q *GetOrderCommand) SetGetOrderRequestArgs(args *GetOrderRequestArgs) {
	q.EtopShopID = args.EtopShopID
	q.TrackingNumber = args.TrackingNumber
}

func (q *GetShippingRateCommand) GetArgs(ctx context.Context) (_ context.Context, _ *GetShippingRateRequestArgs) {
	return ctx,
		&GetShippingRateRequestArgs{
			EtopShopID:  q.EtopShopID,
			Origin:      q.Origin,
			Destination: q.Destination,
			CodAmount:   q.CodAmount,
			TotalGrams:  q.TotalGrams,
		}
}

func (q *GetShippingRateCommand) SetGetShippingRateRequestArgs(args *GetShippingRateRequestArgs) {
	q.EtopShopID = args.EtopShopID
	q.Origin = args.Origin
	q.Destination = args.Destination
	q.CodAmount = args.CodAmount
	q.TotalGrams = args.TotalGrams
}

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleCancelOrder)
	b.AddHandler(h.HandleCreateOrder)
	b.AddHandler(h.HandleGetOrder)
	b.AddHandler(h.HandleGetShippingRate)
	return CommandBus{b}
}
