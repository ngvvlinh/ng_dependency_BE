package model

import (
	"time"

	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	identitymodel "o.o/backend/com/main/identity/model"
	identitysharemodel "o.o/backend/com/main/identity/sharemodel"
	ordermodel "o.o/backend/com/main/ordering/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

// +sqlgen
type MoneyTransactionShippingExternal struct {
	ID             dot.ID
	Code           string
	TotalCOD       int
	TotalOrders    int
	CreatedAt      time.Time `sq:"create"`
	UpdatedAt      time.Time `sq:"update"`
	Status         status3.Status
	ExternalPaidAt time.Time
	Provider       shipping_provider.ShippingProvider
	BankAccount    *identitysharemodel.BankAccount
	Note           string
	InvoiceNumber  string
	ConnectionID   dot.ID
	WLPartnerID    dot.ID
}

// +sqlgen
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
	CreatedAt                          time.Time `sq:"create"`
	UpdatedAt                          time.Time `sq:"update"`
	ImportError                        *model.Error
	ExternalTotalShippingFee           int
}

// +sqlgen: MoneyTransactionShippingExternalLine as m
// +sqlgen:left-join: Fulfillment as f on f.id = m.etop_fulfillment_id
// +sqlgen:left-join: Shop        as s on s.id = f.shop_id
// +sqlgen:left-join: Order       as o on o.id = f.order_id
type MoneyTransactionShippingExternalLineExtended struct {
	*MoneyTransactionShippingExternalLine
	Fulfillment *shipmodel.Fulfillment
	Shop        *identitymodel.Shop
	Order       *ordermodel.Order
}

type MoneyTransactionShippingExternalExtended struct {
	*MoneyTransactionShippingExternal
	Lines []*MoneyTransactionShippingExternalLineExtended
}

type MoneyTransactionShippingExternalFtLine struct {
	*MoneyTransactionShippingExternal
	Lines []*MoneyTransactionShippingExternalLine
}

// +sqlgen
type MoneyTransactionShipping struct {
	ID                                 dot.ID
	ShopID                             dot.ID
	CreatedAt                          time.Time `sq:"create"`
	UpdatedAt                          time.Time `sq:"update"`
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
	BankAccount                        *identitysharemodel.BankAccount
	Note                               string
	InvoiceNumber                      string
	Type                               string
	Rid                                dot.ID
	WLPartnerID                        dot.ID
}

// +sqlgen:           MoneyTransactionShipping as m
// +sqlgen:left-join: Shop as s on s.id = m.shop_id
type MoneyTransactionShippingFtShop struct {
	*MoneyTransactionShipping
	Shop *identitymodel.Shop
}

// +sqlgen
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
	BankAccount           *identitysharemodel.BankAccount
	Note                  string
	InvoiceNumber         string
	WLPartnerID           dot.ID
}

// +convert:type=moneytx.ShopFtMoneyTxShippingCount
// +sqlsel
type ShopFtMoneyTxShippingCount struct {
	MoneyTxShippingCount int    `sel:"count(*)"`
	ShopID               dot.ID `sel:"shop_id"`
}
