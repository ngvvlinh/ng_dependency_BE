package model

import (
	"time"

	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/model"
	ordermodel "etop.vn/backend/pkg/services/ordering/model"
	shipmodel "etop.vn/backend/pkg/services/shipping/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenMoneyTransactionShippingExternal(&MoneyTransactionShippingExternal{})

var _ = sqlgenMoneyTransactionShippingExternalLine(&MoneyTransactionShippingExternalLine{})

var _ = sqlgenMoneyTransactionShippingExternalLineExtended(
	&MoneyTransactionShippingExternalLineExtended{}, &MoneyTransactionShippingExternalLine{}, sq.AS("m"),
	sq.LEFT_JOIN, &shipmodel.Fulfillment{}, sq.AS("f"), "f.id = m.etop_fulfillment_id",
	sq.LEFT_JOIN, &model.Shop{}, sq.AS("s"), "s.id = f.shop_id",
	sq.LEFT_JOIN, &ordermodel.Order{}, sq.AS("o"), "o.id = f.order_id",
)

var _ = sqlgenMoneyTransactionShipping(&MoneyTransactionShipping{})

var _ = sqlgenMoneyTransactionShippingFtShop(
	&MoneyTransactionShippingFtShop{}, &MoneyTransactionShipping{}, sq.AS("m"),
	sq.LEFT_JOIN, &model.Shop{}, sq.AS("s"), "s.id = m.shop_id",
)

var _ = sqlgenMoneyTransactionShippingEtop(&MoneyTransactionShippingEtop{})

type MoneyTransactionShippingExternal struct {
	ID             int64
	Code           string
	TotalCOD       int
	TotalOrders    int
	CreatedAt      time.Time `sq:"create"`
	UpdatedAt      time.Time `sq:"update"`
	Status         model.Status3
	ExternalPaidAt time.Time
	Provider       string
	BankAccount    *model.BankAccount
	Note           string
	InvoiceNumber  string
}

type MoneyTransactionShippingExternalLine struct {
	ID                                 int64
	ExternalCode                       string
	ExternalCustomer                   string
	ExternalAddress                    string
	ExternalTotalCOD                   int
	ExternalCreatedAt                  time.Time
	ExternalClosedAt                   time.Time
	EtopFulfillmentIdRaw               string
	EtopFulfillmentID                  int64
	Note                               string
	MoneyTransactionShippingExternalID int64
	CreatedAt                          time.Time `sq:"create"`
	UpdatedAt                          time.Time `sq:"update"`
	ImportError                        *model.Error
	ExternalTotalShippingFee           int
}

type MoneyTransactionShippingExternalLineExtended struct {
	*MoneyTransactionShippingExternalLine
	Fulfillment *shipmodel.Fulfillment
	Shop        *model.Shop
	Order       *ordermodel.Order
}

type MoneyTransactionShippingExternalExtended struct {
	*MoneyTransactionShippingExternal
	Lines []*MoneyTransactionShippingExternalLineExtended
}

type MoneyTransactionShipping struct {
	ID                                 int64
	ShopID                             int64
	CreatedAt                          time.Time `sq:"create"`
	UpdatedAt                          time.Time `sq:"update"`
	ClosedAt                           time.Time
	Status                             model.Status3
	TotalCOD                           int
	TotalAmount                        int
	TotalOrders                        int
	Code                               string
	MoneyTransactionShippingExternalID int64
	MoneyTransactionShippingEtopID     int64
	Provider                           string
	ConfirmedAt                        time.Time
	EtopTransferedAt                   time.Time
	BankAccount                        *model.BankAccount
	Note                               string
	InvoiceNumber                      string
	Type                               string
}

type MoneyTransactionShippingFtShop struct {
	*MoneyTransactionShipping
	Shop *model.Shop
}

type MoneyTransactionShippingEtop struct {
	ID                    int64
	Code                  string
	TotalCOD              int
	TotalOrders           int
	TotalAmount           int
	TotalFee              int
	TotalMoneyTransaction int
	CreatedAt             time.Time `sq:"create"`
	UpdatedAt             time.Time `sq:"update"`
	ConfirmedAt           time.Time
	Status                model.Status3
	BankAccount           *model.BankAccount
	Note                  string
	InvoiceNumber         string
}
