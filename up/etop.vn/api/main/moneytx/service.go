package moneytx

import (
	"context"
	"time"

	"etop.vn/api/main/identity"
	identitytypes "etop.vn/api/main/identity/types"
	"etop.vn/api/main/shipping"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	// -- Money transaction shipping -- //
	CreateMoneyTxShipping(context.Context, *CreateMoneyTxShippingArgs) (*MoneyTransactionShipping, error)
	CreateMoneyTxShippings(context.Context, *CreateMoneyTxShippingsArgs) (created int, _ error)
	UpdateMoneyTxShippingInfo(context.Context, *UpdateMoneyTxShippingInfoArgs) (*MoneyTransactionShipping, error)
	ConfirmMoneyTxShipping(context.Context, *ConfirmMoneyTxShippingArgs) error
	DeleteMoneyTxShipping(context.Context, *DeleteMoneyTxShippingArgs) error
	AddFulfillmentsMoneyTxShipping(context.Context, *FfmsMoneyTxShippingArgs) error
	RemoveFulfillmentsMoneyTxShipping(context.Context, *FfmsMoneyTxShippingArgs) error
	ReCalcMoneyTxShipping(ctx context.Context, MoneyTxShippingID dot.ID) error

	// -- Money transaction shipping external -- //
	CreateMoneyTxShippingExternal(context.Context, *CreateMoneyTxShippingExternalArgs) (*MoneyTransactionShippingExternalFtLine, error)
	CreateMoneyTxShippingExternalLine(context.Context, *CreateMoneyTxShippingExternalLineArgs) (*MoneyTransactionShippingExternalLine, error)
	UpdateMoneyTxShippingExternalInfo(context.Context, *UpdateMoneyTxShippingExternalInfoArgs) (*MoneyTransactionShippingExternalFtLine, error)
	ConfirmMoneyTxShippingExternal(ctx context.Context, ID dot.ID) (updated int, _ error)
	ConfirmMoneyTxShippingExternals(ctx context.Context, IDs []dot.ID) (updated int, _ error)
	RemoveMoneyTxShippingExternalLines(context.Context, *RemoveMoneyTxShippingExternalLinesArgs) (*MoneyTransactionShippingExternalFtLine, error)
	DeleteMoneyTxShippingExternal(ctx context.Context, ID dot.ID) (deleted int, _ error)
	DeleteMoneyTxShippingExternalLines(ctx context.Context, MoneyTxShippingExternalID dot.ID) error

	// -- Money transaction shipping etop -- //
	CreateMoneyTxShippingEtop(context.Context, *CreateMoneyTxShippingEtopArgs) (*MoneyTransactionShippingEtop, error)
	UpdateMoneyTxShippingEtop(context.Context, UpdateMoneyTxShippingEtopArgs) (*MoneyTransactionShippingEtop, error)
	ConfirmMoneyTxShippingEtop(context.Context, *ConfirmMoneyTxShippingEtopArgs) error
	DeleteMoneyTxShippingEtop(ctx context.Context, ID dot.ID) error
}

type QueryService interface {
	// -- Money transaction shipping -- //
	GetMoneyTxShippingByID(context.Context, *GetMoneyTxByIDQueryArgs) (*MoneyTransactionShipping, error)
	ListMoneyTxShippings(context.Context, *ListMoneyTxShippingArgs) (*ListMoneyTxShippingsResponse, error)
	ListMoneyTxShippingsByMoneyTxShippingExternalID(ctx context.Context, MoneyTxShippingExternalID dot.ID) ([]*MoneyTransactionShipping, error)

	// -- Money transaction shipping external -- //
	GetMoneyTxShippingExternal(ctx context.Context, ID dot.ID) (*MoneyTransactionShippingExternalFtLine, error)
	ListMoneyTxShippingExternals(context.Context, *ListMoneyTxShippingExternalsArgs) (*ListMoneyTxShippingExternalsResponse, error)

	// -- Money transaction shipping etop -- //
	GetMoneyTxShippingEtop(ctx context.Context, ID dot.ID) (*MoneyTransactionShippingEtop, error)
	ListMoneyTxShippingEtops(context.Context, *ListMoneyTxShippingEtopsArgs) (*ListMoneyTxShippingEtopsResponse, error)
}

type FfmsMoneyTxShippingArgs struct {
	FulfillmentIDs    []dot.ID
	MoneyTxShippingID dot.ID
	ShopID            dot.ID
}

type GetMoneyTxByIDQueryArgs struct {
	MoneyTxShippingID dot.ID
	ShopID            dot.ID
}

type CreateMoneyTxShippingArgs struct {
	Shop           *identity.Shop
	FulfillmentIDs []dot.ID
	TotalCOD       int
	TotalAmount    int
	TotalOrders    int
}

type CreateMoneyTxShippingsArgs struct {
	ShopIDMapFfms map[dot.ID][]*shipping.Fulfillment
}

type UpdateMoneyTxShippingInfoArgs struct {
	MoneyTxShippingID dot.ID
	ShopID            dot.ID
	Note              string
	InvoiceNumber     string
	BankAccount       *identitytypes.BankAccount
}

type ConfirmMoneyTxShippingArgs struct {
	MoneyTxShippingID dot.ID
	ShopID            dot.ID
	TotalCOD          int
	TotalAmount       int
	TotalOrders       int
}

type DeleteMoneyTxShippingArgs struct {
	MoneyTxShippingID dot.ID
	ShopID            dot.ID
}

type CreateMoneyTxShippingExternalArgs struct {
	Provider       shipping_provider.ShippingProvider
	ExternalPaidAt time.Time
	Lines          []*MoneyTransactionShippingExternalLine
	BankAccount    *identitytypes.BankAccount
	Note           string
	InvoiceNumber  string
}

type CreateMoneyTxShippingExternalLineArgs struct {
	ExternalCode                       string
	ExternalTotalCOD                   int
	ExternalCreatedAt                  time.Time
	ExternalClosedAt                   time.Time
	EtopFulfillmentIDRaw               string
	ExternalCustomer                   string
	ExternalAddress                    string
	MoneyTransactionShippingExternalID dot.ID
	ExternalTotalShippingFee           int
}

type UpdateMoneyTxShippingExternalInfoArgs struct {
	MoneyTxShippingExternalID dot.ID
	BankAccount               *identitytypes.BankAccount
	Note                      string
	InvoiceNumber             string
}

type RemoveMoneyTxShippingExternalLinesArgs struct {
	MoneyTxShippingExternalID dot.ID
	LineIDs                   []dot.ID
}

type CreateMoneyTxShippingEtopArgs struct {
	MoneyTxShippingIDs []dot.ID
	BankAccount        *identitytypes.BankAccount
	Note               string
	InvoiceNumber      string
}

type UpdateMoneyTxShippingEtopArgs struct {
	MoneyTxShippingEtopID dot.ID
	BankAccount           *identitytypes.BankAccount
	Note                  string
	InvoiceNumber         string

	// MoneyTransactionShipping IDs
	Adds       []dot.ID
	Deletes    []dot.ID
	ReplaceAll []dot.ID
}

type ConfirmMoneyTxShippingEtopArgs struct {
	MoneyTxShippingEtopID dot.ID
	TotalCOD              int
	TotalAmount           int
	TotalOrders           int
}

type ListMoneyTxShippingArgs struct {
	MoneyTxShippingIDs    []dot.ID
	MoneyTxShippingEtopID dot.ID
	ShopID                dot.ID
	Paging                meta.Paging
	Filters               meta.Filters
}

type ListMoneyTxShippingsResponse struct {
	MoneyTxShippings []*MoneyTransactionShipping
	Paging           meta.PageInfo
}

type ListMoneyTxShippingExternalsArgs struct {
	MoneyTxShippingExternalIDs []dot.ID
	Paging                     meta.Paging
	Filters                    meta.Filters
}

type ListMoneyTxShippingExternalsResponse struct {
	MoneyTxShippingExternals []*MoneyTransactionShippingExternalFtLine
	Paging                   meta.PageInfo
}

type ListMoneyTxShippingEtopsArgs struct {
	MoneyTxShippingEtopIDs []dot.ID
	Status                 status3.NullStatus
	Paging                 meta.Paging
	Filter                 meta.Filters
}

type ListMoneyTxShippingEtopsResponse struct {
	MoneyTxShippingEtops []*MoneyTransactionShippingEtop
	Paging               meta.PageInfo
}
