// +build !generator

// Code generated by generator api. DO NOT EDIT.

package shipmentprice

import (
	context "context"

	route_type "o.o/api/top/types/etc/route_type"
	status3 "o.o/api/top/types/etc/status3"
	capi "o.o/capi"
	dot "o.o/capi/dot"
)

type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus     { return QueryBus{bus} }

func (b CommandBus) Dispatch(ctx context.Context, msg interface{ command() }) error {
	return b.bus.Dispatch(ctx, msg)
}
func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type CreateShipmentPriceCommand struct {
	Name                string
	ShipmentPriceListID dot.ID
	ShipmentServiceID   dot.ID
	CustomRegionTypes   []route_type.CustomRegionRouteType
	CustomRegionIDs     []dot.ID
	RegionTypes         []route_type.RegionRouteType
	ProvinceTypes       []route_type.ProvinceRouteType
	UrbanTypes          []route_type.UrbanType
	PriorityPoint       int
	Details             []*PricingDetail
	AdditionalFees      []*AdditionalFee

	Result *ShipmentPrice `json:"-"`
}

func (h AggregateHandler) HandleCreateShipmentPrice(ctx context.Context, msg *CreateShipmentPriceCommand) (err error) {
	msg.Result, err = h.inner.CreateShipmentPrice(msg.GetArgs(ctx))
	return err
}

type DeleteShipmentPriceCommand struct {
	ID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleDeleteShipmentPrice(ctx context.Context, msg *DeleteShipmentPriceCommand) (err error) {
	return h.inner.DeleteShipmentPrice(msg.GetArgs(ctx))
}

type UpdateShipmentPriceCommand struct {
	ID                  dot.ID
	Name                string
	ShipmentPriceListID dot.ID
	ShipmentServiceID   dot.ID
	CustomRegionTypes   []route_type.CustomRegionRouteType
	CustomRegionIDs     []dot.ID
	RegionTypes         []route_type.RegionRouteType
	ProvinceTypes       []route_type.ProvinceRouteType
	UrbanTypes          []route_type.UrbanType
	PriorityPoint       int
	Details             []*PricingDetail
	AdditionalFees      []*AdditionalFee
	Status              status3.Status

	Result *ShipmentPrice `json:"-"`
}

func (h AggregateHandler) HandleUpdateShipmentPrice(ctx context.Context, msg *UpdateShipmentPriceCommand) (err error) {
	msg.Result, err = h.inner.UpdateShipmentPrice(msg.GetArgs(ctx))
	return err
}

type UpdateShipmentPricesPriorityPointCommand struct {
	ShipmentPrices []*UpdateShipmentPricePriorityPointArgs

	Result int `json:"-"`
}

func (h AggregateHandler) HandleUpdateShipmentPricesPriorityPoint(ctx context.Context, msg *UpdateShipmentPricesPriorityPointCommand) (err error) {
	msg.Result, err = h.inner.UpdateShipmentPricesPriorityPoint(msg.GetArgs(ctx))
	return err
}

type CalculateShippingFeesQuery struct {
	AccountID           dot.ID
	FromProvince        string
	FromProvinceCode    string
	FromDistrict        string
	FromDistrictCode    string
	ToProvince          string
	ToProvinceCode      string
	ToDistrict          string
	ToDistrictCode      string
	ShipmentServiceID   dot.ID
	ConnectionID        dot.ID
	ShipmentPriceListID dot.ID
	Weight              int
	BasketValue         int
	CODAmount           int
	IncludeInsurance    bool

	Result *CalculateShippingFeesResponse `json:"-"`
}

func (h QueryServiceHandler) HandleCalculateShippingFees(ctx context.Context, msg *CalculateShippingFeesQuery) (err error) {
	msg.Result, err = h.inner.CalculateShippingFees(msg.GetArgs(ctx))
	return err
}

type GetShipmentPriceQuery struct {
	ID dot.ID

	Result *ShipmentPrice `json:"-"`
}

func (h QueryServiceHandler) HandleGetShipmentPrice(ctx context.Context, msg *GetShipmentPriceQuery) (err error) {
	msg.Result, err = h.inner.GetShipmentPrice(msg.GetArgs(ctx))
	return err
}

type ListShipmentPricesQuery struct {
	ShipmentPriceListID dot.ID
	ShipmentServiceID   dot.ID

	Result []*ShipmentPrice `json:"-"`
}

func (h QueryServiceHandler) HandleListShipmentPrices(ctx context.Context, msg *ListShipmentPricesQuery) (err error) {
	msg.Result, err = h.inner.ListShipmentPrices(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateShipmentPriceCommand) command()               {}
func (q *DeleteShipmentPriceCommand) command()               {}
func (q *UpdateShipmentPriceCommand) command()               {}
func (q *UpdateShipmentPricesPriorityPointCommand) command() {}

func (q *CalculateShippingFeesQuery) query() {}
func (q *GetShipmentPriceQuery) query()      {}
func (q *ListShipmentPricesQuery) query()    {}

// implement conversion

func (q *CreateShipmentPriceCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateShipmentPriceArgs) {
	return ctx,
		&CreateShipmentPriceArgs{
			Name:                q.Name,
			ShipmentPriceListID: q.ShipmentPriceListID,
			ShipmentServiceID:   q.ShipmentServiceID,
			CustomRegionTypes:   q.CustomRegionTypes,
			CustomRegionIDs:     q.CustomRegionIDs,
			RegionTypes:         q.RegionTypes,
			ProvinceTypes:       q.ProvinceTypes,
			UrbanTypes:          q.UrbanTypes,
			PriorityPoint:       q.PriorityPoint,
			Details:             q.Details,
			AdditionalFees:      q.AdditionalFees,
		}
}

func (q *CreateShipmentPriceCommand) SetCreateShipmentPriceArgs(args *CreateShipmentPriceArgs) {
	q.Name = args.Name
	q.ShipmentPriceListID = args.ShipmentPriceListID
	q.ShipmentServiceID = args.ShipmentServiceID
	q.CustomRegionTypes = args.CustomRegionTypes
	q.CustomRegionIDs = args.CustomRegionIDs
	q.RegionTypes = args.RegionTypes
	q.ProvinceTypes = args.ProvinceTypes
	q.UrbanTypes = args.UrbanTypes
	q.PriorityPoint = args.PriorityPoint
	q.Details = args.Details
	q.AdditionalFees = args.AdditionalFees
}

func (q *DeleteShipmentPriceCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *UpdateShipmentPriceCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateShipmentPriceArgs) {
	return ctx,
		&UpdateShipmentPriceArgs{
			ID:                  q.ID,
			Name:                q.Name,
			ShipmentPriceListID: q.ShipmentPriceListID,
			ShipmentServiceID:   q.ShipmentServiceID,
			CustomRegionTypes:   q.CustomRegionTypes,
			CustomRegionIDs:     q.CustomRegionIDs,
			RegionTypes:         q.RegionTypes,
			ProvinceTypes:       q.ProvinceTypes,
			UrbanTypes:          q.UrbanTypes,
			PriorityPoint:       q.PriorityPoint,
			Details:             q.Details,
			AdditionalFees:      q.AdditionalFees,
			Status:              q.Status,
		}
}

func (q *UpdateShipmentPriceCommand) SetUpdateShipmentPriceArgs(args *UpdateShipmentPriceArgs) {
	q.ID = args.ID
	q.Name = args.Name
	q.ShipmentPriceListID = args.ShipmentPriceListID
	q.ShipmentServiceID = args.ShipmentServiceID
	q.CustomRegionTypes = args.CustomRegionTypes
	q.CustomRegionIDs = args.CustomRegionIDs
	q.RegionTypes = args.RegionTypes
	q.ProvinceTypes = args.ProvinceTypes
	q.UrbanTypes = args.UrbanTypes
	q.PriorityPoint = args.PriorityPoint
	q.Details = args.Details
	q.AdditionalFees = args.AdditionalFees
	q.Status = args.Status
}

func (q *UpdateShipmentPricesPriorityPointCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateShipmentPricesPriorityPointArgs) {
	return ctx,
		&UpdateShipmentPricesPriorityPointArgs{
			ShipmentPrices: q.ShipmentPrices,
		}
}

func (q *UpdateShipmentPricesPriorityPointCommand) SetUpdateShipmentPricesPriorityPointArgs(args *UpdateShipmentPricesPriorityPointArgs) {
	q.ShipmentPrices = args.ShipmentPrices
}

func (q *CalculateShippingFeesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *CalculateShippingFeeArgs) {
	return ctx,
		&CalculateShippingFeeArgs{
			AccountID:           q.AccountID,
			FromProvince:        q.FromProvince,
			FromProvinceCode:    q.FromProvinceCode,
			FromDistrict:        q.FromDistrict,
			FromDistrictCode:    q.FromDistrictCode,
			ToProvince:          q.ToProvince,
			ToProvinceCode:      q.ToProvinceCode,
			ToDistrict:          q.ToDistrict,
			ToDistrictCode:      q.ToDistrictCode,
			ShipmentServiceID:   q.ShipmentServiceID,
			ConnectionID:        q.ConnectionID,
			ShipmentPriceListID: q.ShipmentPriceListID,
			Weight:              q.Weight,
			BasketValue:         q.BasketValue,
			CODAmount:           q.CODAmount,
			IncludeInsurance:    q.IncludeInsurance,
		}
}

func (q *CalculateShippingFeesQuery) SetCalculateShippingFeeArgs(args *CalculateShippingFeeArgs) {
	q.AccountID = args.AccountID
	q.FromProvince = args.FromProvince
	q.FromProvinceCode = args.FromProvinceCode
	q.FromDistrict = args.FromDistrict
	q.FromDistrictCode = args.FromDistrictCode
	q.ToProvince = args.ToProvince
	q.ToProvinceCode = args.ToProvinceCode
	q.ToDistrict = args.ToDistrict
	q.ToDistrictCode = args.ToDistrictCode
	q.ShipmentServiceID = args.ShipmentServiceID
	q.ConnectionID = args.ConnectionID
	q.ShipmentPriceListID = args.ShipmentPriceListID
	q.Weight = args.Weight
	q.BasketValue = args.BasketValue
	q.CODAmount = args.CODAmount
	q.IncludeInsurance = args.IncludeInsurance
}

func (q *GetShipmentPriceQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *ListShipmentPricesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListShipmentPricesArgs) {
	return ctx,
		&ListShipmentPricesArgs{
			ShipmentPriceListID: q.ShipmentPriceListID,
			ShipmentServiceID:   q.ShipmentServiceID,
		}
}

func (q *ListShipmentPricesQuery) SetListShipmentPricesArgs(args *ListShipmentPricesArgs) {
	q.ShipmentPriceListID = args.ShipmentPriceListID
	q.ShipmentServiceID = args.ShipmentServiceID
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
	b.AddHandler(h.HandleCreateShipmentPrice)
	b.AddHandler(h.HandleDeleteShipmentPrice)
	b.AddHandler(h.HandleUpdateShipmentPrice)
	b.AddHandler(h.HandleUpdateShipmentPricesPriorityPoint)
	return CommandBus{b}
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleCalculateShippingFees)
	b.AddHandler(h.HandleGetShipmentPrice)
	b.AddHandler(h.HandleListShipmentPrices)
	return QueryBus{b}
}
