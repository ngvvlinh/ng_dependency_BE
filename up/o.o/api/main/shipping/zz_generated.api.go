// +build !generator

// Code generated by generator api. DO NOT EDIT.

package shipping

import (
	context "context"
	json "encoding/json"
	time "time"

	orderingtypes "o.o/api/main/ordering/types"
	shippingtypes "o.o/api/main/shipping/types"
	shipping "o.o/api/top/types/etc/shipping"
	shipping_fee_type "o.o/api/top/types/etc/shipping_fee_type"
	shipping_provider "o.o/api/top/types/etc/shipping_provider"
	status3 "o.o/api/top/types/etc/status3"
	status4 "o.o/api/top/types/etc/status4"
	status5 "o.o/api/top/types/etc/status5"
	try_on "o.o/api/top/types/etc/try_on"
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

type AddFulfillmentShippingFeeCommand struct {
	FulfillmentID   dot.ID
	ShippingCode    string
	ShippingFeeType shipping_fee_type.ShippingFeeType
	UpdatedBy       dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleAddFulfillmentShippingFee(ctx context.Context, msg *AddFulfillmentShippingFeeCommand) (err error) {
	return h.inner.AddFulfillmentShippingFee(msg.GetArgs(ctx))
}

type CancelFulfillmentCommand struct {
	FulfillmentID dot.ID
	CancelReason  string

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleCancelFulfillment(ctx context.Context, msg *CancelFulfillmentCommand) (err error) {
	return h.inner.CancelFulfillment(msg.GetArgs(ctx))
}

type CreateFulfillmentsCommand struct {
	ShopID              dot.ID
	OrderID             dot.ID
	PickupAddress       *orderingtypes.Address
	ShippingAddress     *orderingtypes.Address
	ReturnAddress       *orderingtypes.Address
	ShippingType        orderingtypes.ShippingType
	ShippingServiceCode string
	ShippingServiceFee  int
	ShippingServiceName string
	WeightInfo          shippingtypes.WeightInfo
	ValueInfo           shippingtypes.ValueInfo
	TryOn               try_on.TryOnCode
	ShippingNote        string
	ConnectionID        dot.ID
	ShopCarrierID       dot.ID
	Coupon              string

	Result []dot.ID `json:"-"`
}

func (h AggregateHandler) HandleCreateFulfillments(ctx context.Context, msg *CreateFulfillmentsCommand) (err error) {
	msg.Result, err = h.inner.CreateFulfillments(msg.GetArgs(ctx))
	return err
}

type RemoveFulfillmentsMoneyTxIDCommand struct {
	FulfillmentIDs            []dot.ID
	MoneyTxShippingID         dot.ID
	MoneyTxShippingExternalID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleRemoveFulfillmentsMoneyTxID(ctx context.Context, msg *RemoveFulfillmentsMoneyTxIDCommand) (err error) {
	msg.Result, err = h.inner.RemoveFulfillmentsMoneyTxID(msg.GetArgs(ctx))
	return err
}

type ShopUpdateFulfillmentCODCommand struct {
	FulfillmentID  dot.ID
	ShippingCode   string
	TotalCODAmount dot.NullInt
	UpdatedBy      dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleShopUpdateFulfillmentCOD(ctx context.Context, msg *ShopUpdateFulfillmentCODCommand) (err error) {
	msg.Result, err = h.inner.ShopUpdateFulfillmentCOD(msg.GetArgs(ctx))
	return err
}

type ShopUpdateFulfillmentInfoCommand struct {
	FulfillmentID    dot.ID
	AddressTo        *orderingtypes.Address
	AddressFrom      *orderingtypes.Address
	IncludeInsurance dot.NullBool
	InsuranceValue   dot.NullInt
	GrossWeight      dot.NullInt
	TryOn            try_on.TryOnCode
	ShippingNote     dot.NullString

	Result int `json:"-"`
}

func (h AggregateHandler) HandleShopUpdateFulfillmentInfo(ctx context.Context, msg *ShopUpdateFulfillmentInfoCommand) (err error) {
	msg.Result, err = h.inner.ShopUpdateFulfillmentInfo(msg.GetArgs(ctx))
	return err
}

type UpdateFulfillmentExternalShippingInfoCommand struct {
	FulfillmentID             dot.ID
	ShippingState             shipping.State
	ShippingStatus            status5.Status
	ExternalShippingData      json.RawMessage
	ExternalShippingState     string
	ExternalShippingSubState  dot.NullString
	ExternalShippingStatus    status5.Status
	ExternalShippingNote      dot.NullString
	ExternalShippingUpdatedAt time.Time
	ExternalShippingLogs      []*ExternalShippingLog
	ExternalShippingStateCode string
	Weight                    int
	ClosedAt                  time.Time
	LastSyncAt                time.Time
	ShippingCreatedAt         time.Time
	ShippingPickingAt         time.Time
	ShippingDeliveringAt      time.Time
	ShippingDeliveredAt       time.Time
	ShippingReturningAt       time.Time
	ShippingReturnedAt        time.Time
	ShippingCancelledAt       time.Time

	Result int `json:"-"`
}

func (h AggregateHandler) HandleUpdateFulfillmentExternalShippingInfo(ctx context.Context, msg *UpdateFulfillmentExternalShippingInfoCommand) (err error) {
	msg.Result, err = h.inner.UpdateFulfillmentExternalShippingInfo(msg.GetArgs(ctx))
	return err
}

type UpdateFulfillmentInfoCommand struct {
	FulfillmentID dot.ID
	ShippingCode  string
	FullName      dot.NullString
	Phone         dot.NullString
	AdminNote     string

	Result int `json:"-"`
}

func (h AggregateHandler) HandleUpdateFulfillmentInfo(ctx context.Context, msg *UpdateFulfillmentInfoCommand) (err error) {
	msg.Result, err = h.inner.UpdateFulfillmentInfo(msg.GetArgs(ctx))
	return err
}

type UpdateFulfillmentShippingFeesCommand struct {
	FulfillmentID            dot.ID
	ShippingCode             string
	ProviderShippingFeeLines []*ShippingFeeLine
	ShippingFeeLines         []*ShippingFeeLine
	TotalCODAmount           dot.NullInt
	UpdatedBy                dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleUpdateFulfillmentShippingFees(ctx context.Context, msg *UpdateFulfillmentShippingFeesCommand) (err error) {
	msg.Result, err = h.inner.UpdateFulfillmentShippingFees(msg.GetArgs(ctx))
	return err
}

type UpdateFulfillmentShippingFeesFromWebhookCommand struct {
	FulfillmentID    dot.ID
	NewWeight        int
	NewState         shipping.State
	ProviderFeeLines []*ShippingFeeLine

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateFulfillmentShippingFeesFromWebhook(ctx context.Context, msg *UpdateFulfillmentShippingFeesFromWebhookCommand) (err error) {
	return h.inner.UpdateFulfillmentShippingFeesFromWebhook(msg.GetArgs(ctx))
}

type UpdateFulfillmentShippingStateCommand struct {
	PartnerID                dot.ID
	FulfillmentID            dot.ID
	ShippingCode             string
	ShippingState            shipping.State
	ActualCompensationAmount dot.NullInt
	UpdatedBy                dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleUpdateFulfillmentShippingState(ctx context.Context, msg *UpdateFulfillmentShippingStateCommand) (err error) {
	msg.Result, err = h.inner.UpdateFulfillmentShippingState(msg.GetArgs(ctx))
	return err
}

type UpdateFulfillmentsCODTransferedAtCommand struct {
	FulfillmentIDs     []dot.ID
	MoneyTxShippingIDs []dot.ID
	CODTransferedAt    time.Time

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateFulfillmentsCODTransferedAt(ctx context.Context, msg *UpdateFulfillmentsCODTransferedAtCommand) (err error) {
	return h.inner.UpdateFulfillmentsCODTransferedAt(msg.GetArgs(ctx))
}

type UpdateFulfillmentsMoneyTxIDCommand struct {
	FulfillmentIDs            []dot.ID
	MoneyTxShippingExternalID dot.ID
	MoneyTxShippingID         dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleUpdateFulfillmentsMoneyTxID(ctx context.Context, msg *UpdateFulfillmentsMoneyTxIDCommand) (err error) {
	msg.Result, err = h.inner.UpdateFulfillmentsMoneyTxID(msg.GetArgs(ctx))
	return err
}

type UpdateFulfillmentsStatusCommand struct {
	FulfillmentIDs []dot.ID
	Status         status4.NullStatus
	ShopConfirm    status3.NullStatus
	SyncStatus     status4.NullStatus

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateFulfillmentsStatus(ctx context.Context, msg *UpdateFulfillmentsStatusCommand) (err error) {
	return h.inner.UpdateFulfillmentsStatus(msg.GetArgs(ctx))
}

type GetFulfillmentByIDOrShippingCodeQuery struct {
	ID            dot.ID
	ShippingCode  string
	ConnectionIDs []dot.ID

	Result *Fulfillment `json:"-"`
}

func (h QueryServiceHandler) HandleGetFulfillmentByIDOrShippingCode(ctx context.Context, msg *GetFulfillmentByIDOrShippingCodeQuery) (err error) {
	msg.Result, err = h.inner.GetFulfillmentByIDOrShippingCode(msg.GetArgs(ctx))
	return err
}

type GetFulfillmentExtendedQuery struct {
	ID           dot.ID
	ShippingCode string

	Result *FulfillmentExtended `json:"-"`
}

func (h QueryServiceHandler) HandleGetFulfillmentExtended(ctx context.Context, msg *GetFulfillmentExtendedQuery) (err error) {
	msg.Result, err = h.inner.GetFulfillmentExtended(msg.GetArgs(ctx))
	return err
}

type ListFulfillmentExtendedsByIDsQuery struct {
	IDs    []dot.ID
	ShopID dot.ID

	Result []*FulfillmentExtended `json:"-"`
}

func (h QueryServiceHandler) HandleListFulfillmentExtendedsByIDs(ctx context.Context, msg *ListFulfillmentExtendedsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListFulfillmentExtendedsByIDs(msg.GetArgs(ctx))
	return err
}

type ListFulfillmentExtendedsByMoneyTxShippingIDQuery struct {
	ShopID            dot.ID
	MoneyTxShippingID dot.ID

	Result []*FulfillmentExtended `json:"-"`
}

func (h QueryServiceHandler) HandleListFulfillmentExtendedsByMoneyTxShippingID(ctx context.Context, msg *ListFulfillmentExtendedsByMoneyTxShippingIDQuery) (err error) {
	msg.Result, err = h.inner.ListFulfillmentExtendedsByMoneyTxShippingID(msg.GetArgs(ctx))
	return err
}

type ListFulfillmentsByIDsQuery struct {
	IDs    []dot.ID
	ShopID dot.ID

	Result []*Fulfillment `json:"-"`
}

func (h QueryServiceHandler) HandleListFulfillmentsByIDs(ctx context.Context, msg *ListFulfillmentsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListFulfillmentsByIDs(msg.GetArgs(ctx))
	return err
}

type ListFulfillmentsByMoneyTxQuery struct {
	MoneyTxShippingIDs        []dot.ID
	MoneyTxShippingExternalID dot.ID

	Result []*Fulfillment `json:"-"`
}

func (h QueryServiceHandler) HandleListFulfillmentsByMoneyTx(ctx context.Context, msg *ListFulfillmentsByMoneyTxQuery) (err error) {
	msg.Result, err = h.inner.ListFulfillmentsByMoneyTx(msg.GetArgs(ctx))
	return err
}

type ListFulfillmentsByShippingCodesQuery struct {
	Codes []string

	Result []*Fulfillment `json:"-"`
}

func (h QueryServiceHandler) HandleListFulfillmentsByShippingCodes(ctx context.Context, msg *ListFulfillmentsByShippingCodesQuery) (err error) {
	msg.Result, err = h.inner.ListFulfillmentsByShippingCodes(msg.GetArgs(ctx))
	return err
}

type ListFulfillmentsForMoneyTxQuery struct {
	ShippingProvider shipping_provider.ShippingProvider
	ShippingStates   []shipping.State
	IsNoneCOD        dot.NullBool

	Result []*Fulfillment `json:"-"`
}

func (h QueryServiceHandler) HandleListFulfillmentsForMoneyTx(ctx context.Context, msg *ListFulfillmentsForMoneyTxQuery) (err error) {
	msg.Result, err = h.inner.ListFulfillmentsForMoneyTx(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *AddFulfillmentShippingFeeCommand) command()                {}
func (q *CancelFulfillmentCommand) command()                        {}
func (q *CreateFulfillmentsCommand) command()                       {}
func (q *RemoveFulfillmentsMoneyTxIDCommand) command()              {}
func (q *ShopUpdateFulfillmentCODCommand) command()                 {}
func (q *ShopUpdateFulfillmentInfoCommand) command()                {}
func (q *UpdateFulfillmentExternalShippingInfoCommand) command()    {}
func (q *UpdateFulfillmentInfoCommand) command()                    {}
func (q *UpdateFulfillmentShippingFeesCommand) command()            {}
func (q *UpdateFulfillmentShippingFeesFromWebhookCommand) command() {}
func (q *UpdateFulfillmentShippingStateCommand) command()           {}
func (q *UpdateFulfillmentsCODTransferedAtCommand) command()        {}
func (q *UpdateFulfillmentsMoneyTxIDCommand) command()              {}
func (q *UpdateFulfillmentsStatusCommand) command()                 {}

func (q *GetFulfillmentByIDOrShippingCodeQuery) query()            {}
func (q *GetFulfillmentExtendedQuery) query()                      {}
func (q *ListFulfillmentExtendedsByIDsQuery) query()               {}
func (q *ListFulfillmentExtendedsByMoneyTxShippingIDQuery) query() {}
func (q *ListFulfillmentsByIDsQuery) query()                       {}
func (q *ListFulfillmentsByMoneyTxQuery) query()                   {}
func (q *ListFulfillmentsByShippingCodesQuery) query()             {}
func (q *ListFulfillmentsForMoneyTxQuery) query()                  {}

// implement conversion

func (q *AddFulfillmentShippingFeeCommand) GetArgs(ctx context.Context) (_ context.Context, _ *AddFulfillmentShippingFeeArgs) {
	return ctx,
		&AddFulfillmentShippingFeeArgs{
			FulfillmentID:   q.FulfillmentID,
			ShippingCode:    q.ShippingCode,
			ShippingFeeType: q.ShippingFeeType,
			UpdatedBy:       q.UpdatedBy,
		}
}

func (q *AddFulfillmentShippingFeeCommand) SetAddFulfillmentShippingFeeArgs(args *AddFulfillmentShippingFeeArgs) {
	q.FulfillmentID = args.FulfillmentID
	q.ShippingCode = args.ShippingCode
	q.ShippingFeeType = args.ShippingFeeType
	q.UpdatedBy = args.UpdatedBy
}

func (q *CancelFulfillmentCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CancelFulfillmentArgs) {
	return ctx,
		&CancelFulfillmentArgs{
			FulfillmentID: q.FulfillmentID,
			CancelReason:  q.CancelReason,
		}
}

func (q *CancelFulfillmentCommand) SetCancelFulfillmentArgs(args *CancelFulfillmentArgs) {
	q.FulfillmentID = args.FulfillmentID
	q.CancelReason = args.CancelReason
}

func (q *CreateFulfillmentsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateFulfillmentsArgs) {
	return ctx,
		&CreateFulfillmentsArgs{
			ShopID:              q.ShopID,
			OrderID:             q.OrderID,
			PickupAddress:       q.PickupAddress,
			ShippingAddress:     q.ShippingAddress,
			ReturnAddress:       q.ReturnAddress,
			ShippingType:        q.ShippingType,
			ShippingServiceCode: q.ShippingServiceCode,
			ShippingServiceFee:  q.ShippingServiceFee,
			ShippingServiceName: q.ShippingServiceName,
			WeightInfo:          q.WeightInfo,
			ValueInfo:           q.ValueInfo,
			TryOn:               q.TryOn,
			ShippingNote:        q.ShippingNote,
			ConnectionID:        q.ConnectionID,
			ShopCarrierID:       q.ShopCarrierID,
			Coupon:              q.Coupon,
		}
}

func (q *CreateFulfillmentsCommand) SetCreateFulfillmentsArgs(args *CreateFulfillmentsArgs) {
	q.ShopID = args.ShopID
	q.OrderID = args.OrderID
	q.PickupAddress = args.PickupAddress
	q.ShippingAddress = args.ShippingAddress
	q.ReturnAddress = args.ReturnAddress
	q.ShippingType = args.ShippingType
	q.ShippingServiceCode = args.ShippingServiceCode
	q.ShippingServiceFee = args.ShippingServiceFee
	q.ShippingServiceName = args.ShippingServiceName
	q.WeightInfo = args.WeightInfo
	q.ValueInfo = args.ValueInfo
	q.TryOn = args.TryOn
	q.ShippingNote = args.ShippingNote
	q.ConnectionID = args.ConnectionID
	q.ShopCarrierID = args.ShopCarrierID
	q.Coupon = args.Coupon
}

func (q *RemoveFulfillmentsMoneyTxIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *RemoveFulfillmentsMoneyTxIDArgs) {
	return ctx,
		&RemoveFulfillmentsMoneyTxIDArgs{
			FulfillmentIDs:            q.FulfillmentIDs,
			MoneyTxShippingID:         q.MoneyTxShippingID,
			MoneyTxShippingExternalID: q.MoneyTxShippingExternalID,
		}
}

func (q *RemoveFulfillmentsMoneyTxIDCommand) SetRemoveFulfillmentsMoneyTxIDArgs(args *RemoveFulfillmentsMoneyTxIDArgs) {
	q.FulfillmentIDs = args.FulfillmentIDs
	q.MoneyTxShippingID = args.MoneyTxShippingID
	q.MoneyTxShippingExternalID = args.MoneyTxShippingExternalID
}

func (q *ShopUpdateFulfillmentCODCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ShopUpdateFulfillmentCODArgs) {
	return ctx,
		&ShopUpdateFulfillmentCODArgs{
			FulfillmentID:  q.FulfillmentID,
			ShippingCode:   q.ShippingCode,
			TotalCODAmount: q.TotalCODAmount,
			UpdatedBy:      q.UpdatedBy,
		}
}

func (q *ShopUpdateFulfillmentCODCommand) SetShopUpdateFulfillmentCODArgs(args *ShopUpdateFulfillmentCODArgs) {
	q.FulfillmentID = args.FulfillmentID
	q.ShippingCode = args.ShippingCode
	q.TotalCODAmount = args.TotalCODAmount
	q.UpdatedBy = args.UpdatedBy
}

func (q *ShopUpdateFulfillmentInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateFulfillmentInfoArgs) {
	return ctx,
		&UpdateFulfillmentInfoArgs{
			FulfillmentID:    q.FulfillmentID,
			AddressTo:        q.AddressTo,
			AddressFrom:      q.AddressFrom,
			IncludeInsurance: q.IncludeInsurance,
			InsuranceValue:   q.InsuranceValue,
			GrossWeight:      q.GrossWeight,
			TryOn:            q.TryOn,
			ShippingNote:     q.ShippingNote,
		}
}

func (q *ShopUpdateFulfillmentInfoCommand) SetUpdateFulfillmentInfoArgs(args *UpdateFulfillmentInfoArgs) {
	q.FulfillmentID = args.FulfillmentID
	q.AddressTo = args.AddressTo
	q.AddressFrom = args.AddressFrom
	q.IncludeInsurance = args.IncludeInsurance
	q.InsuranceValue = args.InsuranceValue
	q.GrossWeight = args.GrossWeight
	q.TryOn = args.TryOn
	q.ShippingNote = args.ShippingNote
}

func (q *UpdateFulfillmentExternalShippingInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateFfmExternalShippingInfoArgs) {
	return ctx,
		&UpdateFfmExternalShippingInfoArgs{
			FulfillmentID:             q.FulfillmentID,
			ShippingState:             q.ShippingState,
			ShippingStatus:            q.ShippingStatus,
			ExternalShippingData:      q.ExternalShippingData,
			ExternalShippingState:     q.ExternalShippingState,
			ExternalShippingSubState:  q.ExternalShippingSubState,
			ExternalShippingStatus:    q.ExternalShippingStatus,
			ExternalShippingNote:      q.ExternalShippingNote,
			ExternalShippingUpdatedAt: q.ExternalShippingUpdatedAt,
			ExternalShippingLogs:      q.ExternalShippingLogs,
			ExternalShippingStateCode: q.ExternalShippingStateCode,
			Weight:                    q.Weight,
			ClosedAt:                  q.ClosedAt,
			LastSyncAt:                q.LastSyncAt,
			ShippingCreatedAt:         q.ShippingCreatedAt,
			ShippingPickingAt:         q.ShippingPickingAt,
			ShippingDeliveringAt:      q.ShippingDeliveringAt,
			ShippingDeliveredAt:       q.ShippingDeliveredAt,
			ShippingReturningAt:       q.ShippingReturningAt,
			ShippingReturnedAt:        q.ShippingReturnedAt,
			ShippingCancelledAt:       q.ShippingCancelledAt,
		}
}

func (q *UpdateFulfillmentExternalShippingInfoCommand) SetUpdateFfmExternalShippingInfoArgs(args *UpdateFfmExternalShippingInfoArgs) {
	q.FulfillmentID = args.FulfillmentID
	q.ShippingState = args.ShippingState
	q.ShippingStatus = args.ShippingStatus
	q.ExternalShippingData = args.ExternalShippingData
	q.ExternalShippingState = args.ExternalShippingState
	q.ExternalShippingSubState = args.ExternalShippingSubState
	q.ExternalShippingStatus = args.ExternalShippingStatus
	q.ExternalShippingNote = args.ExternalShippingNote
	q.ExternalShippingUpdatedAt = args.ExternalShippingUpdatedAt
	q.ExternalShippingLogs = args.ExternalShippingLogs
	q.ExternalShippingStateCode = args.ExternalShippingStateCode
	q.Weight = args.Weight
	q.ClosedAt = args.ClosedAt
	q.LastSyncAt = args.LastSyncAt
	q.ShippingCreatedAt = args.ShippingCreatedAt
	q.ShippingPickingAt = args.ShippingPickingAt
	q.ShippingDeliveringAt = args.ShippingDeliveringAt
	q.ShippingDeliveredAt = args.ShippingDeliveredAt
	q.ShippingReturningAt = args.ShippingReturningAt
	q.ShippingReturnedAt = args.ShippingReturnedAt
	q.ShippingCancelledAt = args.ShippingCancelledAt
}

func (q *UpdateFulfillmentInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateFulfillmentInfoByAdminArgs) {
	return ctx,
		&UpdateFulfillmentInfoByAdminArgs{
			FulfillmentID: q.FulfillmentID,
			ShippingCode:  q.ShippingCode,
			FullName:      q.FullName,
			Phone:         q.Phone,
			AdminNote:     q.AdminNote,
		}
}

func (q *UpdateFulfillmentInfoCommand) SetUpdateFulfillmentInfoByAdminArgs(args *UpdateFulfillmentInfoByAdminArgs) {
	q.FulfillmentID = args.FulfillmentID
	q.ShippingCode = args.ShippingCode
	q.FullName = args.FullName
	q.Phone = args.Phone
	q.AdminNote = args.AdminNote
}

func (q *UpdateFulfillmentShippingFeesCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateFulfillmentShippingFeesArgs) {
	return ctx,
		&UpdateFulfillmentShippingFeesArgs{
			FulfillmentID:            q.FulfillmentID,
			ShippingCode:             q.ShippingCode,
			ProviderShippingFeeLines: q.ProviderShippingFeeLines,
			ShippingFeeLines:         q.ShippingFeeLines,
			TotalCODAmount:           q.TotalCODAmount,
			UpdatedBy:                q.UpdatedBy,
		}
}

func (q *UpdateFulfillmentShippingFeesCommand) SetUpdateFulfillmentShippingFeesArgs(args *UpdateFulfillmentShippingFeesArgs) {
	q.FulfillmentID = args.FulfillmentID
	q.ShippingCode = args.ShippingCode
	q.ProviderShippingFeeLines = args.ProviderShippingFeeLines
	q.ShippingFeeLines = args.ShippingFeeLines
	q.TotalCODAmount = args.TotalCODAmount
	q.UpdatedBy = args.UpdatedBy
}

func (q *UpdateFulfillmentShippingFeesFromWebhookCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateFulfillmentShippingFeesFromWebhookArgs) {
	return ctx,
		&UpdateFulfillmentShippingFeesFromWebhookArgs{
			FulfillmentID:    q.FulfillmentID,
			NewWeight:        q.NewWeight,
			NewState:         q.NewState,
			ProviderFeeLines: q.ProviderFeeLines,
		}
}

func (q *UpdateFulfillmentShippingFeesFromWebhookCommand) SetUpdateFulfillmentShippingFeesFromWebhookArgs(args *UpdateFulfillmentShippingFeesFromWebhookArgs) {
	q.FulfillmentID = args.FulfillmentID
	q.NewWeight = args.NewWeight
	q.NewState = args.NewState
	q.ProviderFeeLines = args.ProviderFeeLines
}

func (q *UpdateFulfillmentShippingStateCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateFulfillmentShippingStateArgs) {
	return ctx,
		&UpdateFulfillmentShippingStateArgs{
			PartnerID:                q.PartnerID,
			FulfillmentID:            q.FulfillmentID,
			ShippingCode:             q.ShippingCode,
			ShippingState:            q.ShippingState,
			ActualCompensationAmount: q.ActualCompensationAmount,
			UpdatedBy:                q.UpdatedBy,
		}
}

func (q *UpdateFulfillmentShippingStateCommand) SetUpdateFulfillmentShippingStateArgs(args *UpdateFulfillmentShippingStateArgs) {
	q.PartnerID = args.PartnerID
	q.FulfillmentID = args.FulfillmentID
	q.ShippingCode = args.ShippingCode
	q.ShippingState = args.ShippingState
	q.ActualCompensationAmount = args.ActualCompensationAmount
	q.UpdatedBy = args.UpdatedBy
}

func (q *UpdateFulfillmentsCODTransferedAtCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateFulfillmentsCODTransferedAtArgs) {
	return ctx,
		&UpdateFulfillmentsCODTransferedAtArgs{
			FulfillmentIDs:     q.FulfillmentIDs,
			MoneyTxShippingIDs: q.MoneyTxShippingIDs,
			CODTransferedAt:    q.CODTransferedAt,
		}
}

func (q *UpdateFulfillmentsCODTransferedAtCommand) SetUpdateFulfillmentsCODTransferedAtArgs(args *UpdateFulfillmentsCODTransferedAtArgs) {
	q.FulfillmentIDs = args.FulfillmentIDs
	q.MoneyTxShippingIDs = args.MoneyTxShippingIDs
	q.CODTransferedAt = args.CODTransferedAt
}

func (q *UpdateFulfillmentsMoneyTxIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateFulfillmentsMoneyTxIDArgs) {
	return ctx,
		&UpdateFulfillmentsMoneyTxIDArgs{
			FulfillmentIDs:            q.FulfillmentIDs,
			MoneyTxShippingExternalID: q.MoneyTxShippingExternalID,
			MoneyTxShippingID:         q.MoneyTxShippingID,
		}
}

func (q *UpdateFulfillmentsMoneyTxIDCommand) SetUpdateFulfillmentsMoneyTxIDArgs(args *UpdateFulfillmentsMoneyTxIDArgs) {
	q.FulfillmentIDs = args.FulfillmentIDs
	q.MoneyTxShippingExternalID = args.MoneyTxShippingExternalID
	q.MoneyTxShippingID = args.MoneyTxShippingID
}

func (q *UpdateFulfillmentsStatusCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateFulfillmentsStatusArgs) {
	return ctx,
		&UpdateFulfillmentsStatusArgs{
			FulfillmentIDs: q.FulfillmentIDs,
			Status:         q.Status,
			ShopConfirm:    q.ShopConfirm,
			SyncStatus:     q.SyncStatus,
		}
}

func (q *UpdateFulfillmentsStatusCommand) SetUpdateFulfillmentsStatusArgs(args *UpdateFulfillmentsStatusArgs) {
	q.FulfillmentIDs = args.FulfillmentIDs
	q.Status = args.Status
	q.ShopConfirm = args.ShopConfirm
	q.SyncStatus = args.SyncStatus
}

func (q *GetFulfillmentByIDOrShippingCodeQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetFulfillmentByIDOrShippingCodeArgs) {
	return ctx,
		&GetFulfillmentByIDOrShippingCodeArgs{
			ID:            q.ID,
			ShippingCode:  q.ShippingCode,
			ConnectionIDs: q.ConnectionIDs,
		}
}

func (q *GetFulfillmentByIDOrShippingCodeQuery) SetGetFulfillmentByIDOrShippingCodeArgs(args *GetFulfillmentByIDOrShippingCodeArgs) {
	q.ID = args.ID
	q.ShippingCode = args.ShippingCode
	q.ConnectionIDs = args.ConnectionIDs
}

func (q *GetFulfillmentExtendedQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, ShippingCode string) {
	return ctx,
		q.ID,
		q.ShippingCode
}

func (q *ListFulfillmentExtendedsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, IDs []dot.ID, ShopID dot.ID) {
	return ctx,
		q.IDs,
		q.ShopID
}

func (q *ListFulfillmentExtendedsByMoneyTxShippingIDQuery) GetArgs(ctx context.Context) (_ context.Context, shopID dot.ID, moneyTxShippingID dot.ID) {
	return ctx,
		q.ShopID,
		q.MoneyTxShippingID
}

func (q *ListFulfillmentsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, IDs []dot.ID, shopID dot.ID) {
	return ctx,
		q.IDs,
		q.ShopID
}

func (q *ListFulfillmentsByMoneyTxQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListFullfillmentsByMoneyTxArgs) {
	return ctx,
		&ListFullfillmentsByMoneyTxArgs{
			MoneyTxShippingIDs:        q.MoneyTxShippingIDs,
			MoneyTxShippingExternalID: q.MoneyTxShippingExternalID,
		}
}

func (q *ListFulfillmentsByMoneyTxQuery) SetListFullfillmentsByMoneyTxArgs(args *ListFullfillmentsByMoneyTxArgs) {
	q.MoneyTxShippingIDs = args.MoneyTxShippingIDs
	q.MoneyTxShippingExternalID = args.MoneyTxShippingExternalID
}

func (q *ListFulfillmentsByShippingCodesQuery) GetArgs(ctx context.Context) (_ context.Context, Codes []string) {
	return ctx,
		q.Codes
}

func (q *ListFulfillmentsForMoneyTxQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListFulfillmentForMoneyTxArgs) {
	return ctx,
		&ListFulfillmentForMoneyTxArgs{
			ShippingProvider: q.ShippingProvider,
			ShippingStates:   q.ShippingStates,
			IsNoneCOD:        q.IsNoneCOD,
		}
}

func (q *ListFulfillmentsForMoneyTxQuery) SetListFulfillmentForMoneyTxArgs(args *ListFulfillmentForMoneyTxArgs) {
	q.ShippingProvider = args.ShippingProvider
	q.ShippingStates = args.ShippingStates
	q.IsNoneCOD = args.IsNoneCOD
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
	b.AddHandler(h.HandleAddFulfillmentShippingFee)
	b.AddHandler(h.HandleCancelFulfillment)
	b.AddHandler(h.HandleCreateFulfillments)
	b.AddHandler(h.HandleRemoveFulfillmentsMoneyTxID)
	b.AddHandler(h.HandleShopUpdateFulfillmentCOD)
	b.AddHandler(h.HandleShopUpdateFulfillmentInfo)
	b.AddHandler(h.HandleUpdateFulfillmentExternalShippingInfo)
	b.AddHandler(h.HandleUpdateFulfillmentInfo)
	b.AddHandler(h.HandleUpdateFulfillmentShippingFees)
	b.AddHandler(h.HandleUpdateFulfillmentShippingFeesFromWebhook)
	b.AddHandler(h.HandleUpdateFulfillmentShippingState)
	b.AddHandler(h.HandleUpdateFulfillmentsCODTransferedAt)
	b.AddHandler(h.HandleUpdateFulfillmentsMoneyTxID)
	b.AddHandler(h.HandleUpdateFulfillmentsStatus)
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
	b.AddHandler(h.HandleGetFulfillmentByIDOrShippingCode)
	b.AddHandler(h.HandleGetFulfillmentExtended)
	b.AddHandler(h.HandleListFulfillmentExtendedsByIDs)
	b.AddHandler(h.HandleListFulfillmentExtendedsByMoneyTxShippingID)
	b.AddHandler(h.HandleListFulfillmentsByIDs)
	b.AddHandler(h.HandleListFulfillmentsByMoneyTx)
	b.AddHandler(h.HandleListFulfillmentsByShippingCodes)
	b.AddHandler(h.HandleListFulfillmentsForMoneyTx)
	return QueryBus{b}
}
