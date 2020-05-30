package moneytx

import (
	"time"

	"o.o/api/main/identity"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/main/ordering"
	"o.o/api/main/shipping"
	"o.o/api/meta"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:event:topic=event/moneytx

var ShippingAcceptStates = []string{
	shippingstate.Returned.String(), shippingstate.Returning.String(), shippingstate.Delivered.String(), shippingstate.Undeliverable.String(),
}

type MoneyTransactionShipping struct {
	ID                                 dot.ID
	ShopID                             dot.ID
	CreatedAt                          time.Time
	UpdatedAt                          time.Time
	ClosedAt                           time.Time
	Status                             status3.Status
	TotalCOD                           int
	TotalAmount                        int
	TotalOrders                        int
	Code                               string
	MoneyTransactionShippingExternalID dot.ID
	MoneyTransactionShippingEtopID     dot.ID
	Provider                           string
	ConfirmedAt                        time.Time
	EtopTransferedAt                   time.Time
	BankAccount                        *identitytypes.BankAccount
	Note                               string
	InvoiceNumber                      string
	Type                               string
}

type MoneyTransactionShippingExternal struct {
	ID             dot.ID
	Code           string
	TotalCOD       int
	TotalOrders    int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Status         status3.Status
	ExternalPaidAt time.Time
	Provider       shipping_provider.ShippingProvider
	BankAccount    *identitytypes.BankAccount
	Note           string
	InvoiceNumber  string
}

type MoneyTransactionShippingExternalLine struct {
	ID                                 dot.ID
	ExternalCode                       string
	ExternalCustomer                   string
	ExternalAddress                    string
	ExternalTotalCOD                   int
	ExternalCreatedAt                  time.Time
	ExternalClosedAt                   time.Time
	EtopFulfillmentIDRaw               string
	EtopFulfillmentID                  dot.ID
	Note                               string
	MoneyTransactionShippingExternalID dot.ID
	CreatedAt                          time.Time
	UpdatedAt                          time.Time
	ImportError                        *meta.Error
	ExternalTotalShippingFee           int
}

type MoneyTransactionShippingEtop struct {
	ID                    dot.ID
	Code                  string
	TotalCOD              int
	TotalOrders           int
	TotalAmount           int
	TotalFee              int
	TotalMoneyTransaction int
	CreatedAt             time.Time `sq:"create"`
	UpdatedAt             time.Time `sq:"update"`
	ConfirmedAt           time.Time
	Status                status3.Status
	BankAccount           *identitytypes.BankAccount
	Note                  string
	InvoiceNumber         string
}

type MoneyTransactionShippingExtended struct {
	*MoneyTransactionShipping
	Fulfillments []*shipping.FulfillmentExtended
}

type MoneyTransactionShippingEtopExtended struct {
	*MoneyTransactionShippingEtop
	MoneyTransactions []*MoneyTransactionShippingExtended
}

type MoneyTransactionShippingExternalExtended struct {
	*MoneyTransactionShippingExternal
	Lines []*MoneyTransactionShippingExternalLineExtended
}

type MoneyTransactionShippingExternalFtLine struct {
	*MoneyTransactionShippingExternal
	Lines []*MoneyTransactionShippingExternalLine
}

type MoneyTransactionShippingExternalLineExtended struct {
	*MoneyTransactionShippingExternalLine
	Fulfillment *shipping.Fulfillment
	Shop        *identity.Shop
	Order       *ordering.Order
}

// -- events -- //
type MoneyTxShippingConfirmedEvent struct {
	meta.EventMeta
	ShopID            dot.ID
	MoneyTxShippingID dot.ID
	ConfirmedAt       time.Time
}

type MoneyTxShippingCreatedEvent struct {
	meta.EventMeta
	MoneyTxShippingID dot.ID
	ShopID            dot.ID
	FulfillmentIDs    []dot.ID
}

type MoneyTxShippingDeletedEvent struct {
	meta.EventMeta
	MoneyTxShippingID dot.ID
	ShopID            dot.ID
}

type MoneyTxShippingRemovedFfmsEvent struct {
	meta.EventMeta
	MoneyTxShippingID dot.ID
	FulfillmentIDs    []dot.ID
}

type MoneyTxShippingExternalsConfirmingEvent struct {
	meta.EventMeta
	MoneyTxShippingExternalIDs []dot.ID
}

type MoneyTxShippingExternalLinesDeletedEvent struct {
	meta.EventMeta
	MoneyTxShippingExternalID      dot.ID
	MoneyTxShippingExternalLineIDs []dot.ID
	FulfillmentIDs                 []dot.ID
}

type MoneyTxShippingExternalCreatedEvent struct {
	meta.EventMeta
	MoneyTxShippingExternalID dot.ID
	FulfillementIDs           []dot.ID
}

type MoneyTxShippingExternalDeletedEvent struct {
	meta.EventMeta
	MoneyTxShippingExternalID dot.ID
}

type MoneyTxShippingEtopConfirmedEvent struct {
	meta.EventMeta
	MoneyTxShippingEtopID dot.ID
	MoneyTxShippingIDs    []dot.ID
	ConfirmedAt           time.Time
}
type FulfilmentStatistics struct {
	TotalCOD         int
	TotalAmount      int
	TotalOrders      int
	TotalShippingFee int
	FulfillmentIDs   []dot.ID
}
