package model

import (
	"time"

	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenMoneyTransactionShippingExternal(&MoneyTransactionShippingExternal{})

var _ = sqlgenMoneyTransactionShippingExternalLine(&MoneyTransactionShippingExternalLine{})

var _ = sqlgenMoneyTransactionShippingExternalLineExtended(
	&MoneyTransactionShippingExternalLineExtended{}, &MoneyTransactionShippingExternalLine{}, "m",
	sq.LEFT_JOIN, &shipmodel.Fulfillment{}, "f", "f.id = m.etop_fulfillment_id",
	sq.LEFT_JOIN, &identitymodel.Shop{}, "s", "s.id = f.shop_id",
	sq.LEFT_JOIN, &ordermodel.Order{}, "o", "o.id = f.order_id",
)

var _ = sqlgenMoneyTransactionShipping(&MoneyTransactionShipping{})

var _ = sqlgenMoneyTransactionShippingFtShop(
	&MoneyTransactionShippingFtShop{}, &MoneyTransactionShipping{}, "m",
	sq.LEFT_JOIN, &identitymodel.Shop{}, "s", "s.id = m.shop_id",
)

var _ = sqlgenMoneyTransactionShippingEtop(&MoneyTransactionShippingEtop{})

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
}

type MoneyTransactionShippingFtShop struct {
	*MoneyTransactionShipping
	Shop *identitymodel.Shop
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
	BankAccount           *identitysharemodel.BankAccount
	Note                  string
	InvoiceNumber         string
}
