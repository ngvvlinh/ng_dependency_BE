// +build !generator

// Code generated by generator api. DO NOT EDIT.

package moneytx

import (
	context "context"
	time "time"

	identity "o.o/api/main/identity"
	identitytypes "o.o/api/main/identity/types"
	shipping "o.o/api/main/shipping"
	meta "o.o/api/meta"
	shipping_provider "o.o/api/top/types/etc/shipping_provider"
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

type AddFulfillmentsMoneyTxShippingCommand struct {
	FulfillmentIDs    []dot.ID
	MoneyTxShippingID dot.ID
	ShopID            dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleAddFulfillmentsMoneyTxShipping(ctx context.Context, msg *AddFulfillmentsMoneyTxShippingCommand) (err error) {
	return h.inner.AddFulfillmentsMoneyTxShipping(msg.GetArgs(ctx))
}

type ConfirmMoneyTxShippingCommand struct {
	MoneyTxShippingID dot.ID
	ShopID            dot.ID
	TotalCOD          int
	TotalAmount       int
	TotalOrders       int

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleConfirmMoneyTxShipping(ctx context.Context, msg *ConfirmMoneyTxShippingCommand) (err error) {
	return h.inner.ConfirmMoneyTxShipping(msg.GetArgs(ctx))
}

type ConfirmMoneyTxShippingEtopCommand struct {
	MoneyTxShippingEtopID dot.ID
	TotalCOD              int
	TotalAmount           int
	TotalOrders           int

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleConfirmMoneyTxShippingEtop(ctx context.Context, msg *ConfirmMoneyTxShippingEtopCommand) (err error) {
	return h.inner.ConfirmMoneyTxShippingEtop(msg.GetArgs(ctx))
}

type ConfirmMoneyTxShippingExternalCommand struct {
	ID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleConfirmMoneyTxShippingExternal(ctx context.Context, msg *ConfirmMoneyTxShippingExternalCommand) (err error) {
	msg.Result, err = h.inner.ConfirmMoneyTxShippingExternal(msg.GetArgs(ctx))
	return err
}

type ConfirmMoneyTxShippingExternalsCommand struct {
	IDs []dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleConfirmMoneyTxShippingExternals(ctx context.Context, msg *ConfirmMoneyTxShippingExternalsCommand) (err error) {
	msg.Result, err = h.inner.ConfirmMoneyTxShippingExternals(msg.GetArgs(ctx))
	return err
}

type CreateMoneyTxShippingCommand struct {
	Shop           *identity.Shop
	FulfillmentIDs []dot.ID
	TotalCOD       int
	TotalAmount    int
	TotalOrders    int

	Result *MoneyTransactionShipping `json:"-"`
}

func (h AggregateHandler) HandleCreateMoneyTxShipping(ctx context.Context, msg *CreateMoneyTxShippingCommand) (err error) {
	msg.Result, err = h.inner.CreateMoneyTxShipping(msg.GetArgs(ctx))
	return err
}

type CreateMoneyTxShippingEtopCommand struct {
	MoneyTxShippingIDs []dot.ID
	BankAccount        *identitytypes.BankAccount
	Note               string
	InvoiceNumber      string

	Result *MoneyTransactionShippingEtop `json:"-"`
}

func (h AggregateHandler) HandleCreateMoneyTxShippingEtop(ctx context.Context, msg *CreateMoneyTxShippingEtopCommand) (err error) {
	msg.Result, err = h.inner.CreateMoneyTxShippingEtop(msg.GetArgs(ctx))
	return err
}

type CreateMoneyTxShippingExternalCommand struct {
	Provider       shipping_provider.ShippingProvider
	ConnectionID   dot.ID
	ExternalPaidAt time.Time
	Lines          []*MoneyTransactionShippingExternalLine
	BankAccount    *identitytypes.BankAccount
	Note           string
	InvoiceNumber  string

	Result *MoneyTransactionShippingExternalFtLine `json:"-"`
}

func (h AggregateHandler) HandleCreateMoneyTxShippingExternal(ctx context.Context, msg *CreateMoneyTxShippingExternalCommand) (err error) {
	msg.Result, err = h.inner.CreateMoneyTxShippingExternal(msg.GetArgs(ctx))
	return err
}

type CreateMoneyTxShippingExternalLineCommand struct {
	ExternalCode                       string
	ExternalTotalCOD                   int
	ExternalCreatedAt                  time.Time
	ExternalClosedAt                   time.Time
	EtopFulfillmentIDRaw               string
	ExternalCustomer                   string
	ExternalAddress                    string
	MoneyTransactionShippingExternalID dot.ID
	ExternalTotalShippingFee           int

	Result *MoneyTransactionShippingExternalLine `json:"-"`
}

func (h AggregateHandler) HandleCreateMoneyTxShippingExternalLine(ctx context.Context, msg *CreateMoneyTxShippingExternalLineCommand) (err error) {
	msg.Result, err = h.inner.CreateMoneyTxShippingExternalLine(msg.GetArgs(ctx))
	return err
}

type CreateMoneyTxShippingsCommand struct {
	ShopIDMapFfms map[dot.ID][]*shipping.Fulfillment

	Result int `json:"-"`
}

func (h AggregateHandler) HandleCreateMoneyTxShippings(ctx context.Context, msg *CreateMoneyTxShippingsCommand) (err error) {
	msg.Result, err = h.inner.CreateMoneyTxShippings(msg.GetArgs(ctx))
	return err
}

type DeleteMoneyTxShippingCommand struct {
	MoneyTxShippingID dot.ID
	ShopID            dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleDeleteMoneyTxShipping(ctx context.Context, msg *DeleteMoneyTxShippingCommand) (err error) {
	return h.inner.DeleteMoneyTxShipping(msg.GetArgs(ctx))
}

type DeleteMoneyTxShippingEtopCommand struct {
	ID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleDeleteMoneyTxShippingEtop(ctx context.Context, msg *DeleteMoneyTxShippingEtopCommand) (err error) {
	return h.inner.DeleteMoneyTxShippingEtop(msg.GetArgs(ctx))
}

type DeleteMoneyTxShippingExternalCommand struct {
	ID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteMoneyTxShippingExternal(ctx context.Context, msg *DeleteMoneyTxShippingExternalCommand) (err error) {
	msg.Result, err = h.inner.DeleteMoneyTxShippingExternal(msg.GetArgs(ctx))
	return err
}

type DeleteMoneyTxShippingExternalLinesCommand struct {
	MoneyTxShippingExternalID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleDeleteMoneyTxShippingExternalLines(ctx context.Context, msg *DeleteMoneyTxShippingExternalLinesCommand) (err error) {
	return h.inner.DeleteMoneyTxShippingExternalLines(msg.GetArgs(ctx))
}

type ReCalcMoneyTxShippingCommand struct {
	MoneyTxShippingID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleReCalcMoneyTxShipping(ctx context.Context, msg *ReCalcMoneyTxShippingCommand) (err error) {
	return h.inner.ReCalcMoneyTxShipping(msg.GetArgs(ctx))
}

type ReCalcMoneyTxShippingEtopCommand struct {
	MoneyTxShippingEtopID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleReCalcMoneyTxShippingEtop(ctx context.Context, msg *ReCalcMoneyTxShippingEtopCommand) (err error) {
	return h.inner.ReCalcMoneyTxShippingEtop(msg.GetArgs(ctx))
}

type RemoveFulfillmentsMoneyTxShippingCommand struct {
	FulfillmentIDs    []dot.ID
	MoneyTxShippingID dot.ID
	ShopID            dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleRemoveFulfillmentsMoneyTxShipping(ctx context.Context, msg *RemoveFulfillmentsMoneyTxShippingCommand) (err error) {
	return h.inner.RemoveFulfillmentsMoneyTxShipping(msg.GetArgs(ctx))
}

type RemoveMoneyTxShippingExternalLinesCommand struct {
	MoneyTxShippingExternalID dot.ID
	LineIDs                   []dot.ID

	Result *MoneyTransactionShippingExternalFtLine `json:"-"`
}

func (h AggregateHandler) HandleRemoveMoneyTxShippingExternalLines(ctx context.Context, msg *RemoveMoneyTxShippingExternalLinesCommand) (err error) {
	msg.Result, err = h.inner.RemoveMoneyTxShippingExternalLines(msg.GetArgs(ctx))
	return err
}

type UpdateMoneyTxShippingEtopCommand struct {
	MoneyTxShippingEtopID dot.ID
	BankAccount           *identitytypes.BankAccount
	Note                  string
	InvoiceNumber         string
	Adds                  []dot.ID
	Deletes               []dot.ID
	ReplaceAll            []dot.ID

	Result *MoneyTransactionShippingEtop `json:"-"`
}

func (h AggregateHandler) HandleUpdateMoneyTxShippingEtop(ctx context.Context, msg *UpdateMoneyTxShippingEtopCommand) (err error) {
	msg.Result, err = h.inner.UpdateMoneyTxShippingEtop(msg.GetArgs(ctx))
	return err
}

type UpdateMoneyTxShippingExternalInfoCommand struct {
	MoneyTxShippingExternalID dot.ID
	BankAccount               *identitytypes.BankAccount
	Note                      string
	InvoiceNumber             string

	Result *MoneyTransactionShippingExternalFtLine `json:"-"`
}

func (h AggregateHandler) HandleUpdateMoneyTxShippingExternalInfo(ctx context.Context, msg *UpdateMoneyTxShippingExternalInfoCommand) (err error) {
	msg.Result, err = h.inner.UpdateMoneyTxShippingExternalInfo(msg.GetArgs(ctx))
	return err
}

type UpdateMoneyTxShippingInfoCommand struct {
	MoneyTxShippingID dot.ID
	ShopID            dot.ID
	Note              string
	InvoiceNumber     string
	BankAccount       *identitytypes.BankAccount

	Result *MoneyTransactionShipping `json:"-"`
}

func (h AggregateHandler) HandleUpdateMoneyTxShippingInfo(ctx context.Context, msg *UpdateMoneyTxShippingInfoCommand) (err error) {
	msg.Result, err = h.inner.UpdateMoneyTxShippingInfo(msg.GetArgs(ctx))
	return err
}

type GetMoneyTxShippingByIDQuery struct {
	MoneyTxShippingID dot.ID
	ShopID            dot.ID

	Result *MoneyTransactionShipping `json:"-"`
}

func (h QueryServiceHandler) HandleGetMoneyTxShippingByID(ctx context.Context, msg *GetMoneyTxShippingByIDQuery) (err error) {
	msg.Result, err = h.inner.GetMoneyTxShippingByID(msg.GetArgs(ctx))
	return err
}

type GetMoneyTxShippingEtopQuery struct {
	ID dot.ID

	Result *MoneyTransactionShippingEtop `json:"-"`
}

func (h QueryServiceHandler) HandleGetMoneyTxShippingEtop(ctx context.Context, msg *GetMoneyTxShippingEtopQuery) (err error) {
	msg.Result, err = h.inner.GetMoneyTxShippingEtop(msg.GetArgs(ctx))
	return err
}

type GetMoneyTxShippingExternalQuery struct {
	ID dot.ID

	Result *MoneyTransactionShippingExternalFtLine `json:"-"`
}

func (h QueryServiceHandler) HandleGetMoneyTxShippingExternal(ctx context.Context, msg *GetMoneyTxShippingExternalQuery) (err error) {
	msg.Result, err = h.inner.GetMoneyTxShippingExternal(msg.GetArgs(ctx))
	return err
}

type ListMoneyTxShippingEtopsQuery struct {
	MoneyTxShippingEtopIDs []dot.ID
	Status                 status3.NullStatus
	Paging                 meta.Paging
	Filter                 meta.Filters

	Result *ListMoneyTxShippingEtopsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListMoneyTxShippingEtops(ctx context.Context, msg *ListMoneyTxShippingEtopsQuery) (err error) {
	msg.Result, err = h.inner.ListMoneyTxShippingEtops(msg.GetArgs(ctx))
	return err
}

type ListMoneyTxShippingExternalsQuery struct {
	MoneyTxShippingExternalIDs []dot.ID
	Paging                     meta.Paging
	Filters                    meta.Filters

	Result *ListMoneyTxShippingExternalsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListMoneyTxShippingExternals(ctx context.Context, msg *ListMoneyTxShippingExternalsQuery) (err error) {
	msg.Result, err = h.inner.ListMoneyTxShippingExternals(msg.GetArgs(ctx))
	return err
}

type ListMoneyTxShippingsQuery struct {
	MoneyTxShippingIDs    []dot.ID
	MoneyTxShippingEtopID dot.ID
	ShopID                dot.ID
	Paging                meta.Paging
	Filters               meta.Filters

	Result *ListMoneyTxShippingsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListMoneyTxShippings(ctx context.Context, msg *ListMoneyTxShippingsQuery) (err error) {
	msg.Result, err = h.inner.ListMoneyTxShippings(msg.GetArgs(ctx))
	return err
}

type ListMoneyTxShippingsByMoneyTxShippingExternalIDQuery struct {
	MoneyTxShippingExternalID dot.ID

	Result []*MoneyTransactionShipping `json:"-"`
}

func (h QueryServiceHandler) HandleListMoneyTxShippingsByMoneyTxShippingExternalID(ctx context.Context, msg *ListMoneyTxShippingsByMoneyTxShippingExternalIDQuery) (err error) {
	msg.Result, err = h.inner.ListMoneyTxShippingsByMoneyTxShippingExternalID(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *AddFulfillmentsMoneyTxShippingCommand) command()     {}
func (q *ConfirmMoneyTxShippingCommand) command()             {}
func (q *ConfirmMoneyTxShippingEtopCommand) command()         {}
func (q *ConfirmMoneyTxShippingExternalCommand) command()     {}
func (q *ConfirmMoneyTxShippingExternalsCommand) command()    {}
func (q *CreateMoneyTxShippingCommand) command()              {}
func (q *CreateMoneyTxShippingEtopCommand) command()          {}
func (q *CreateMoneyTxShippingExternalCommand) command()      {}
func (q *CreateMoneyTxShippingExternalLineCommand) command()  {}
func (q *CreateMoneyTxShippingsCommand) command()             {}
func (q *DeleteMoneyTxShippingCommand) command()              {}
func (q *DeleteMoneyTxShippingEtopCommand) command()          {}
func (q *DeleteMoneyTxShippingExternalCommand) command()      {}
func (q *DeleteMoneyTxShippingExternalLinesCommand) command() {}
func (q *ReCalcMoneyTxShippingCommand) command()              {}
func (q *ReCalcMoneyTxShippingEtopCommand) command()          {}
func (q *RemoveFulfillmentsMoneyTxShippingCommand) command()  {}
func (q *RemoveMoneyTxShippingExternalLinesCommand) command() {}
func (q *UpdateMoneyTxShippingEtopCommand) command()          {}
func (q *UpdateMoneyTxShippingExternalInfoCommand) command()  {}
func (q *UpdateMoneyTxShippingInfoCommand) command()          {}

func (q *GetMoneyTxShippingByIDQuery) query()                          {}
func (q *GetMoneyTxShippingEtopQuery) query()                          {}
func (q *GetMoneyTxShippingExternalQuery) query()                      {}
func (q *ListMoneyTxShippingEtopsQuery) query()                        {}
func (q *ListMoneyTxShippingExternalsQuery) query()                    {}
func (q *ListMoneyTxShippingsQuery) query()                            {}
func (q *ListMoneyTxShippingsByMoneyTxShippingExternalIDQuery) query() {}

// implement conversion

func (q *AddFulfillmentsMoneyTxShippingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *FfmsMoneyTxShippingArgs) {
	return ctx,
		&FfmsMoneyTxShippingArgs{
			FulfillmentIDs:    q.FulfillmentIDs,
			MoneyTxShippingID: q.MoneyTxShippingID,
			ShopID:            q.ShopID,
		}
}

func (q *AddFulfillmentsMoneyTxShippingCommand) SetFfmsMoneyTxShippingArgs(args *FfmsMoneyTxShippingArgs) {
	q.FulfillmentIDs = args.FulfillmentIDs
	q.MoneyTxShippingID = args.MoneyTxShippingID
	q.ShopID = args.ShopID
}

func (q *ConfirmMoneyTxShippingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConfirmMoneyTxShippingArgs) {
	return ctx,
		&ConfirmMoneyTxShippingArgs{
			MoneyTxShippingID: q.MoneyTxShippingID,
			ShopID:            q.ShopID,
			TotalCOD:          q.TotalCOD,
			TotalAmount:       q.TotalAmount,
			TotalOrders:       q.TotalOrders,
		}
}

func (q *ConfirmMoneyTxShippingCommand) SetConfirmMoneyTxShippingArgs(args *ConfirmMoneyTxShippingArgs) {
	q.MoneyTxShippingID = args.MoneyTxShippingID
	q.ShopID = args.ShopID
	q.TotalCOD = args.TotalCOD
	q.TotalAmount = args.TotalAmount
	q.TotalOrders = args.TotalOrders
}

func (q *ConfirmMoneyTxShippingEtopCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConfirmMoneyTxShippingEtopArgs) {
	return ctx,
		&ConfirmMoneyTxShippingEtopArgs{
			MoneyTxShippingEtopID: q.MoneyTxShippingEtopID,
			TotalCOD:              q.TotalCOD,
			TotalAmount:           q.TotalAmount,
			TotalOrders:           q.TotalOrders,
		}
}

func (q *ConfirmMoneyTxShippingEtopCommand) SetConfirmMoneyTxShippingEtopArgs(args *ConfirmMoneyTxShippingEtopArgs) {
	q.MoneyTxShippingEtopID = args.MoneyTxShippingEtopID
	q.TotalCOD = args.TotalCOD
	q.TotalAmount = args.TotalAmount
	q.TotalOrders = args.TotalOrders
}

func (q *ConfirmMoneyTxShippingExternalCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *ConfirmMoneyTxShippingExternalsCommand) GetArgs(ctx context.Context) (_ context.Context, IDs []dot.ID) {
	return ctx,
		q.IDs
}

func (q *CreateMoneyTxShippingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateMoneyTxShippingArgs) {
	return ctx,
		&CreateMoneyTxShippingArgs{
			Shop:           q.Shop,
			FulfillmentIDs: q.FulfillmentIDs,
			TotalCOD:       q.TotalCOD,
			TotalAmount:    q.TotalAmount,
			TotalOrders:    q.TotalOrders,
		}
}

func (q *CreateMoneyTxShippingCommand) SetCreateMoneyTxShippingArgs(args *CreateMoneyTxShippingArgs) {
	q.Shop = args.Shop
	q.FulfillmentIDs = args.FulfillmentIDs
	q.TotalCOD = args.TotalCOD
	q.TotalAmount = args.TotalAmount
	q.TotalOrders = args.TotalOrders
}

func (q *CreateMoneyTxShippingEtopCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateMoneyTxShippingEtopArgs) {
	return ctx,
		&CreateMoneyTxShippingEtopArgs{
			MoneyTxShippingIDs: q.MoneyTxShippingIDs,
			BankAccount:        q.BankAccount,
			Note:               q.Note,
			InvoiceNumber:      q.InvoiceNumber,
		}
}

func (q *CreateMoneyTxShippingEtopCommand) SetCreateMoneyTxShippingEtopArgs(args *CreateMoneyTxShippingEtopArgs) {
	q.MoneyTxShippingIDs = args.MoneyTxShippingIDs
	q.BankAccount = args.BankAccount
	q.Note = args.Note
	q.InvoiceNumber = args.InvoiceNumber
}

func (q *CreateMoneyTxShippingExternalCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateMoneyTxShippingExternalArgs) {
	return ctx,
		&CreateMoneyTxShippingExternalArgs{
			Provider:       q.Provider,
			ConnectionID:   q.ConnectionID,
			ExternalPaidAt: q.ExternalPaidAt,
			Lines:          q.Lines,
			BankAccount:    q.BankAccount,
			Note:           q.Note,
			InvoiceNumber:  q.InvoiceNumber,
		}
}

func (q *CreateMoneyTxShippingExternalCommand) SetCreateMoneyTxShippingExternalArgs(args *CreateMoneyTxShippingExternalArgs) {
	q.Provider = args.Provider
	q.ConnectionID = args.ConnectionID
	q.ExternalPaidAt = args.ExternalPaidAt
	q.Lines = args.Lines
	q.BankAccount = args.BankAccount
	q.Note = args.Note
	q.InvoiceNumber = args.InvoiceNumber
}

func (q *CreateMoneyTxShippingExternalLineCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateMoneyTxShippingExternalLineArgs) {
	return ctx,
		&CreateMoneyTxShippingExternalLineArgs{
			ExternalCode:                       q.ExternalCode,
			ExternalTotalCOD:                   q.ExternalTotalCOD,
			ExternalCreatedAt:                  q.ExternalCreatedAt,
			ExternalClosedAt:                   q.ExternalClosedAt,
			EtopFulfillmentIDRaw:               q.EtopFulfillmentIDRaw,
			ExternalCustomer:                   q.ExternalCustomer,
			ExternalAddress:                    q.ExternalAddress,
			MoneyTransactionShippingExternalID: q.MoneyTransactionShippingExternalID,
			ExternalTotalShippingFee:           q.ExternalTotalShippingFee,
		}
}

func (q *CreateMoneyTxShippingExternalLineCommand) SetCreateMoneyTxShippingExternalLineArgs(args *CreateMoneyTxShippingExternalLineArgs) {
	q.ExternalCode = args.ExternalCode
	q.ExternalTotalCOD = args.ExternalTotalCOD
	q.ExternalCreatedAt = args.ExternalCreatedAt
	q.ExternalClosedAt = args.ExternalClosedAt
	q.EtopFulfillmentIDRaw = args.EtopFulfillmentIDRaw
	q.ExternalCustomer = args.ExternalCustomer
	q.ExternalAddress = args.ExternalAddress
	q.MoneyTransactionShippingExternalID = args.MoneyTransactionShippingExternalID
	q.ExternalTotalShippingFee = args.ExternalTotalShippingFee
}

func (q *CreateMoneyTxShippingsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateMoneyTxShippingsArgs) {
	return ctx,
		&CreateMoneyTxShippingsArgs{
			ShopIDMapFfms: q.ShopIDMapFfms,
		}
}

func (q *CreateMoneyTxShippingsCommand) SetCreateMoneyTxShippingsArgs(args *CreateMoneyTxShippingsArgs) {
	q.ShopIDMapFfms = args.ShopIDMapFfms
}

func (q *DeleteMoneyTxShippingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteMoneyTxShippingArgs) {
	return ctx,
		&DeleteMoneyTxShippingArgs{
			MoneyTxShippingID: q.MoneyTxShippingID,
			ShopID:            q.ShopID,
		}
}

func (q *DeleteMoneyTxShippingCommand) SetDeleteMoneyTxShippingArgs(args *DeleteMoneyTxShippingArgs) {
	q.MoneyTxShippingID = args.MoneyTxShippingID
	q.ShopID = args.ShopID
}

func (q *DeleteMoneyTxShippingEtopCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *DeleteMoneyTxShippingExternalCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *DeleteMoneyTxShippingExternalLinesCommand) GetArgs(ctx context.Context) (_ context.Context, MoneyTxShippingExternalID dot.ID) {
	return ctx,
		q.MoneyTxShippingExternalID
}

func (q *ReCalcMoneyTxShippingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ReCalcMoneyTxShippingArgs) {
	return ctx,
		&ReCalcMoneyTxShippingArgs{
			MoneyTxShippingID: q.MoneyTxShippingID,
		}
}

func (q *ReCalcMoneyTxShippingCommand) SetReCalcMoneyTxShippingArgs(args *ReCalcMoneyTxShippingArgs) {
	q.MoneyTxShippingID = args.MoneyTxShippingID
}

func (q *ReCalcMoneyTxShippingEtopCommand) GetArgs(ctx context.Context) (_ context.Context, MoneyTxShippingEtopID dot.ID) {
	return ctx,
		q.MoneyTxShippingEtopID
}

func (q *RemoveFulfillmentsMoneyTxShippingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *FfmsMoneyTxShippingArgs) {
	return ctx,
		&FfmsMoneyTxShippingArgs{
			FulfillmentIDs:    q.FulfillmentIDs,
			MoneyTxShippingID: q.MoneyTxShippingID,
			ShopID:            q.ShopID,
		}
}

func (q *RemoveFulfillmentsMoneyTxShippingCommand) SetFfmsMoneyTxShippingArgs(args *FfmsMoneyTxShippingArgs) {
	q.FulfillmentIDs = args.FulfillmentIDs
	q.MoneyTxShippingID = args.MoneyTxShippingID
	q.ShopID = args.ShopID
}

func (q *RemoveMoneyTxShippingExternalLinesCommand) GetArgs(ctx context.Context) (_ context.Context, _ *RemoveMoneyTxShippingExternalLinesArgs) {
	return ctx,
		&RemoveMoneyTxShippingExternalLinesArgs{
			MoneyTxShippingExternalID: q.MoneyTxShippingExternalID,
			LineIDs:                   q.LineIDs,
		}
}

func (q *RemoveMoneyTxShippingExternalLinesCommand) SetRemoveMoneyTxShippingExternalLinesArgs(args *RemoveMoneyTxShippingExternalLinesArgs) {
	q.MoneyTxShippingExternalID = args.MoneyTxShippingExternalID
	q.LineIDs = args.LineIDs
}

func (q *UpdateMoneyTxShippingEtopCommand) GetArgs(ctx context.Context) (_ context.Context, _ UpdateMoneyTxShippingEtopArgs) {
	return ctx,
		UpdateMoneyTxShippingEtopArgs{
			MoneyTxShippingEtopID: q.MoneyTxShippingEtopID,
			BankAccount:           q.BankAccount,
			Note:                  q.Note,
			InvoiceNumber:         q.InvoiceNumber,
			Adds:                  q.Adds,
			Deletes:               q.Deletes,
			ReplaceAll:            q.ReplaceAll,
		}
}

func (q *UpdateMoneyTxShippingEtopCommand) SetUpdateMoneyTxShippingEtopArgs(args UpdateMoneyTxShippingEtopArgs) {
	q.MoneyTxShippingEtopID = args.MoneyTxShippingEtopID
	q.BankAccount = args.BankAccount
	q.Note = args.Note
	q.InvoiceNumber = args.InvoiceNumber
	q.Adds = args.Adds
	q.Deletes = args.Deletes
	q.ReplaceAll = args.ReplaceAll
}

func (q *UpdateMoneyTxShippingExternalInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateMoneyTxShippingExternalInfoArgs) {
	return ctx,
		&UpdateMoneyTxShippingExternalInfoArgs{
			MoneyTxShippingExternalID: q.MoneyTxShippingExternalID,
			BankAccount:               q.BankAccount,
			Note:                      q.Note,
			InvoiceNumber:             q.InvoiceNumber,
		}
}

func (q *UpdateMoneyTxShippingExternalInfoCommand) SetUpdateMoneyTxShippingExternalInfoArgs(args *UpdateMoneyTxShippingExternalInfoArgs) {
	q.MoneyTxShippingExternalID = args.MoneyTxShippingExternalID
	q.BankAccount = args.BankAccount
	q.Note = args.Note
	q.InvoiceNumber = args.InvoiceNumber
}

func (q *UpdateMoneyTxShippingInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateMoneyTxShippingInfoArgs) {
	return ctx,
		&UpdateMoneyTxShippingInfoArgs{
			MoneyTxShippingID: q.MoneyTxShippingID,
			ShopID:            q.ShopID,
			Note:              q.Note,
			InvoiceNumber:     q.InvoiceNumber,
			BankAccount:       q.BankAccount,
		}
}

func (q *UpdateMoneyTxShippingInfoCommand) SetUpdateMoneyTxShippingInfoArgs(args *UpdateMoneyTxShippingInfoArgs) {
	q.MoneyTxShippingID = args.MoneyTxShippingID
	q.ShopID = args.ShopID
	q.Note = args.Note
	q.InvoiceNumber = args.InvoiceNumber
	q.BankAccount = args.BankAccount
}

func (q *GetMoneyTxShippingByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetMoneyTxByIDQueryArgs) {
	return ctx,
		&GetMoneyTxByIDQueryArgs{
			MoneyTxShippingID: q.MoneyTxShippingID,
			ShopID:            q.ShopID,
		}
}

func (q *GetMoneyTxShippingByIDQuery) SetGetMoneyTxByIDQueryArgs(args *GetMoneyTxByIDQueryArgs) {
	q.MoneyTxShippingID = args.MoneyTxShippingID
	q.ShopID = args.ShopID
}

func (q *GetMoneyTxShippingEtopQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetMoneyTxShippingExternalQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *ListMoneyTxShippingEtopsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListMoneyTxShippingEtopsArgs) {
	return ctx,
		&ListMoneyTxShippingEtopsArgs{
			MoneyTxShippingEtopIDs: q.MoneyTxShippingEtopIDs,
			Status:                 q.Status,
			Paging:                 q.Paging,
			Filter:                 q.Filter,
		}
}

func (q *ListMoneyTxShippingEtopsQuery) SetListMoneyTxShippingEtopsArgs(args *ListMoneyTxShippingEtopsArgs) {
	q.MoneyTxShippingEtopIDs = args.MoneyTxShippingEtopIDs
	q.Status = args.Status
	q.Paging = args.Paging
	q.Filter = args.Filter
}

func (q *ListMoneyTxShippingExternalsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListMoneyTxShippingExternalsArgs) {
	return ctx,
		&ListMoneyTxShippingExternalsArgs{
			MoneyTxShippingExternalIDs: q.MoneyTxShippingExternalIDs,
			Paging:                     q.Paging,
			Filters:                    q.Filters,
		}
}

func (q *ListMoneyTxShippingExternalsQuery) SetListMoneyTxShippingExternalsArgs(args *ListMoneyTxShippingExternalsArgs) {
	q.MoneyTxShippingExternalIDs = args.MoneyTxShippingExternalIDs
	q.Paging = args.Paging
	q.Filters = args.Filters
}

func (q *ListMoneyTxShippingsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListMoneyTxShippingArgs) {
	return ctx,
		&ListMoneyTxShippingArgs{
			MoneyTxShippingIDs:    q.MoneyTxShippingIDs,
			MoneyTxShippingEtopID: q.MoneyTxShippingEtopID,
			ShopID:                q.ShopID,
			Paging:                q.Paging,
			Filters:               q.Filters,
		}
}

func (q *ListMoneyTxShippingsQuery) SetListMoneyTxShippingArgs(args *ListMoneyTxShippingArgs) {
	q.MoneyTxShippingIDs = args.MoneyTxShippingIDs
	q.MoneyTxShippingEtopID = args.MoneyTxShippingEtopID
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
}

func (q *ListMoneyTxShippingsByMoneyTxShippingExternalIDQuery) GetArgs(ctx context.Context) (_ context.Context, MoneyTxShippingExternalID dot.ID) {
	return ctx,
		q.MoneyTxShippingExternalID
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
	b.AddHandler(h.HandleAddFulfillmentsMoneyTxShipping)
	b.AddHandler(h.HandleConfirmMoneyTxShipping)
	b.AddHandler(h.HandleConfirmMoneyTxShippingEtop)
	b.AddHandler(h.HandleConfirmMoneyTxShippingExternal)
	b.AddHandler(h.HandleConfirmMoneyTxShippingExternals)
	b.AddHandler(h.HandleCreateMoneyTxShipping)
	b.AddHandler(h.HandleCreateMoneyTxShippingEtop)
	b.AddHandler(h.HandleCreateMoneyTxShippingExternal)
	b.AddHandler(h.HandleCreateMoneyTxShippingExternalLine)
	b.AddHandler(h.HandleCreateMoneyTxShippings)
	b.AddHandler(h.HandleDeleteMoneyTxShipping)
	b.AddHandler(h.HandleDeleteMoneyTxShippingEtop)
	b.AddHandler(h.HandleDeleteMoneyTxShippingExternal)
	b.AddHandler(h.HandleDeleteMoneyTxShippingExternalLines)
	b.AddHandler(h.HandleReCalcMoneyTxShipping)
	b.AddHandler(h.HandleReCalcMoneyTxShippingEtop)
	b.AddHandler(h.HandleRemoveFulfillmentsMoneyTxShipping)
	b.AddHandler(h.HandleRemoveMoneyTxShippingExternalLines)
	b.AddHandler(h.HandleUpdateMoneyTxShippingEtop)
	b.AddHandler(h.HandleUpdateMoneyTxShippingExternalInfo)
	b.AddHandler(h.HandleUpdateMoneyTxShippingInfo)
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
	b.AddHandler(h.HandleGetMoneyTxShippingByID)
	b.AddHandler(h.HandleGetMoneyTxShippingEtop)
	b.AddHandler(h.HandleGetMoneyTxShippingExternal)
	b.AddHandler(h.HandleListMoneyTxShippingEtops)
	b.AddHandler(h.HandleListMoneyTxShippingExternals)
	b.AddHandler(h.HandleListMoneyTxShippings)
	b.AddHandler(h.HandleListMoneyTxShippingsByMoneyTxShippingExternalID)
	return QueryBus{b}
}
