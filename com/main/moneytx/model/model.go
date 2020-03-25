package model

import (
	"time"

	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
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
}
