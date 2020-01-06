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
	CreateMoneyTxShipping(context.Context, CreateMoneyTxShippingArgs) (*MoneyTransactionShippingExtended, error)
	CreateMoneyTxShippings(context.Context, *CreateMoneyTxShippingsArgs) (created int, _ error)
	UpdateMoneyTxShippingInfo(context.Context, *UpdateMoneyTxShippingInfoArgs) (*MoneyTransactionShippingExtended, error)
	ConfirmMoneyTxShipping(context.Context, *ConfirmMoneyTxShippingArgs) (updated int, _ error)
	DeleteMoneyTxShipping(context.Context, *DeleteMoneyTxShippingArgs) (deleted int, _ error)
	AddFulfillmentMoneyTxShipping(context.Context, *FfmMoneyTxShippingArgs) (updated int, _ error)
	RemoveFulfillmentMoneyTxShipping(context.Context, *FfmMoneyTxShippingArgs) (removed int, _ error)
	ReCalcMoneyTxShipping(ctx context.Context, MoneyTxShippingID dot.ID) error

	// -- Money transaction shipping external -- //
	CreateMoneyTxShippingExternal(context.Context, *CreateMoneyTxShippingExternalArgs) (*MoneyTransactionShippingExternalExtended, error)
	CreateMoneyTxShippingExternalLine(context.Context, *CreateMoneyTxShippingExternalLineArgs) (*MoneyTransactionShippingExternalLine, error)
	UpdateMoneyTxShippingExternalInfo(context.Context, *UpdateMoneyTxShippingExternalInfoArgs) (*MoneyTransactionShippingExternalExtended, error)
	ConfirmMoneyTxShippingExternals(ctx context.Context, IDs []dot.ID) (updated int, _ error)
	RemoveMoneyTxShippingExternalLines(context.Context, *RemoveMoneyTxShippingExternalLinesArgs) (*MoneyTransactionShippingExternalExtended, error)
	DeleteMoneyTxShippingExternal(ctx context.Context, ID dot.ID) (deleted int, _ error)

	// -- Money transaction shipping etop -- //
	CreateMoneyTxShippingEtop(context.Context, *CreateMoneyTxShippingEtopArgs) (*MoneyTransactionShippingEtopExtended, error)
	UpdateMoneyTxShippingEtop(context.Context, UpdateMoneyTxShippingEtopArgs) (*MoneyTransactionShippingEtopExtended, error)
	ConfirmMoneyTxShippingEtop(context.Context, *ConfirmMoneyTxShippingEtopArgs) (updated int, _ error)
	DeleteMoneyTxShippingEtop(ctx context.Context, ID dot.ID) (deleted int, _ error)
}

type QueryService interface {
	// -- Money transaction shipping -- //
	GetMoneyTxShippingByID(context.Context, *GetMoneyTxByIDQueryArgs) (*MoneyTransactionShippingExtended, error)
	ListMoneyTxShippings(context.Context, *ListMoneyTxArgs) (*ListMoneyTxShippingsResponse, error)
	ListMoneyTxShippingsByMoneyTxShippingExternalID(ctx context.Context, MoneyTxShippingExternalID dot.ID) ([]*MoneyTransactionShippingExtended, error)

	// -- Money transaction shipping external -- //
	GetMoneyTxShippingExternal(ctx context.Context, ID dot.ID) (*MoneyTransactionShippingExternalExtended, error)
	ListMoneyTxShippingExternals(context.Context, *ListMoneyTxShippingExternalsArgs) (*ListMoneyTxShippingExternalsResponse, error)

	// -- Money transaction shipping etop -- //
	GetMoneyTxShippingEtop(ctx context.Context, ID dot.ID) (*MoneyTransactionShippingEtopExtended, error)
	ListMoneyTxShippingEtops(context.Context, *ListMoneyTxShippingEtopsArgs) (*ListMoneyTxShippingEtopsResponse, error)
}

type FfmMoneyTxShippingArgs struct {
	FulfillmentID     dot.ID
	MoneyTxShippingID dot.ID
}

type GetMoneyTxByIDQueryArgs struct {
	ID     dot.ID
	ShopID dot.ID
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
	ShopIDMap     map[dot.ID]*identity.Shop
}

type UpdateMoneyTxShippingInfoArgs struct {
	ID            dot.ID
	ShopID        dot.ID
	Note          string
	InvoiceNumber string
	BankAccount   *identitytypes.BankAccount
}

type ConfirmMoneyTxShippingArgs struct {
	ID          dot.ID
	ShopID      dot.ID
	TotalCOD    int
	TotalAmount int
	TotalOrders int
}

type DeleteMoneyTxShippingArgs struct {
	ID     dot.ID
	ShopID dot.ID
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
	ID            dot.ID
	BankAccount   *identitytypes.BankAccount
	Note          string
	InvoiceNumber string
}

type RemoveMoneyTxShippingExternalLinesArgs struct {
	ID      dot.ID
	LineIDs []dot.ID
}

type CreateMoneyTxShippingEtopArgs struct {
	MoneyTransactionShippingIDs []dot.ID
	BankAccount                 *identitytypes.BankAccount
	Note                        string
	InvoiceNumber               string
}

type UpdateMoneyTxShippingEtopArgs struct {
	ID            dot.ID
	BankAccount   *identitytypes.BankAccount
	Note          string
	InvoiceNumber string

	// MoneyTransactionShipping IDs
	Adds       []dot.ID
	Deletes    []dot.ID
	ReplaceAll []dot.ID
}

type ConfirmMoneyTxShippingEtopArgs struct {
	ID          dot.ID
	TotalCOD    int
	TotalAmount int
	TotalOrders int
}

type ListMoneyTxArgs struct {
	IDs                []dot.ID
	ShopID             dot.ID
	IncludeFulfillment bool
	Paging             meta.Paging
	Filters            meta.Filters
}

type ListMoneyTxShippingsResponse struct {
	MoneyTxShippings []*MoneyTransactionShippingExtended
	Paging           meta.PageInfo
}

type ListMoneyTxShippingExternalsArgs struct {
	IDs     []dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type ListMoneyTxShippingExternalsResponse struct {
	MoneyTxShippingExternals []*MoneyTransactionShippingExternalExtended
	Paging                   meta.PageInfo
}

type ListMoneyTxShippingEtopsArgs struct {
	IDs    []dot.ID
	Status status3.NullStatus
	Paging meta.Paging
	Filter meta.Filters
}

type ListMoneyTxShippingEtopsResponse struct {
	MoneyTxShippingEtops []*MoneyTransactionShippingEtopExtended
	Paging               meta.PageInfo
}
