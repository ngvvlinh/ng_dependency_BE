// Code generated by gen-cmd-query. DO NOT EDIT.

package gateway

import (
	context "context"

	haravan "etop.vn/api/external/haravan"
	meta "etop.vn/api/meta"
)

type Command interface{ command() }
type Query interface{ query() }
type CommandBus struct{ bus meta.Bus }
type QueryBus struct{ bus meta.Bus }

func (c CommandBus) Dispatch(ctx context.Context, msg Command) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c CommandBus) DispatchAll(ctx context.Context, msgs ...Command) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
func (c QueryBus) DispatchAll(ctx context.Context, msgs ...Query) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

type CancelOrderCommand struct {
	EtopShopID     int64
	TrackingNumber string `json:"tracking_number"`

	Result *GetOrderResponse `json:"-"`
}

type CreateOrderCommand struct {
	EtopShopID            int64
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
	ShippingRateID        int32            `json:"shipping_rate_id"`

	Result *CreateOrderResponse `json:"-"`
}

type GetOrderCommand struct {
	EtopShopID     int64
	TrackingNumber string `json:"tracking_number"`

	Result *GetOrderResponse `json:"-"`
}

type GetShippingRateCommand struct {
	EtopShopID  int64
	Origin      *haravan.Address `json:"origin"`
	Destination *haravan.Address `json:"destination"`
	CodAmount   float32          `json:"cod_amount"`
	TotalGrams  float32          `json:"total_grams"`

	Result *GetShippingRateResponse `json:"-"`
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

func (q *GetOrderCommand) GetArgs(ctx context.Context) (_ context.Context, _ *GetOrderRequestArgs) {
	return ctx,
		&GetOrderRequestArgs{
			EtopShopID:     q.EtopShopID,
			TrackingNumber: q.TrackingNumber,
		}
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

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	meta.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleCancelOrder)
	b.AddHandler(h.HandleCreateOrder)
	b.AddHandler(h.HandleGetOrder)
	b.AddHandler(h.HandleGetShippingRate)
	return CommandBus{b}
}

func (h AggregateHandler) HandleCancelOrder(ctx context.Context, msg *CancelOrderCommand) error {
	result, err := h.inner.CancelOrder(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleCreateOrder(ctx context.Context, msg *CreateOrderCommand) error {
	result, err := h.inner.CreateOrder(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleGetOrder(ctx context.Context, msg *GetOrderCommand) error {
	result, err := h.inner.GetOrder(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleGetShippingRate(ctx context.Context, msg *GetShippingRateCommand) error {
	result, err := h.inner.GetShippingRate(msg.GetArgs(ctx))
	msg.Result = result
	return err
}
