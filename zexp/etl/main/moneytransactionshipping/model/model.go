package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	"etop.vn/capi/dot"
)

// +sqlgen
type MoneyTransactionShipping struct {
	ID            dot.ID
	ShopID        dot.ID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ClosedAt      time.Time
	Status        status3.Status `sql_gen:"int2"`
	TotalCOD      int
	TotalAmount   int
	TotalOrders   int
	Code          string
	Provider      string `sql_gen:"enum(shipping_provider)"`
	ConfirmedAt   time.Time
	BankAccount   *identitysharemodel.BankAccount `sql_gen:"json"`
	Note          string
	InvoiceNumber string
	Type          string
	Rid           dot.ID
}
