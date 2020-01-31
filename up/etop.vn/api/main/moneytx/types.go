package moneytx

import (
	"time"

	"etop.vn/api/main/identity"
	identitytypes "etop.vn/api/main/identity/types"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/shipping"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/moneytx

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
	Fulfillments []*shipping.Fulfillment
}

type MoneyTransactionShippingEtopExtended struct {
	*MoneyTransactionShippingEtop
	MoneyTransactions []*MoneyTransactionShippingExtended
}

type MoneyTransactionShippingExternalExtended struct {
	*MoneyTransactionShippingExternal
	Lines []*MoneyTransactionShippingExternalLineExtended
}

type MoneyTransactionShippingExternalLineExtended struct {
	*MoneyTransactionShippingExternalLine
	Fulfillment *shipping.Fulfillment
	Shop        *identity.Shop
	Order       *ordering.Order
}

// -- events -- //
type MoneyTransactionShippingExternalCreatedEvent struct {
	MoneyTxShippingExternalID dot.ID
	FulfillementIDs           []dot.ID
}

type MoneyTransactionConfirmedEvent struct {
	meta.EventMeta
	ShopID             dot.ID
	MoneyTransactionID dot.ID
}

type MoneyTxShippingEtopConfirmedEvent struct {
	meta.EventMeta
	MoneyTxShippingEtopID dot.ID
}
