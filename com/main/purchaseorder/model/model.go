package model

import (
	"time"

	"etop.vn/api/main/etop"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenPurchaseOrder(&PurchaseOrder{})

type PurchaseOrder struct {
	ID              int64
	ShopID          int64
	SupplierID      int64
	Supplier        *PurchaseOrderSupplier
	BasketValue     int64
	TotalDiscount   int64
	TotalAmount     int64
	Code            string
	CodeNorm        int32
	Note            string
	Status          etop.Status3
	VariantIDs      []int64
	Lines           []*PurchaseOrderLine
	CreatedBy       int64
	CancelledReason string
	ConfirmedAt     time.Time
	CancelledAt     time.Time
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
	DeletedAt       time.Time
}

type PurchaseOrderLine struct {
	ProductName string `json:"product_name"`
	ProductID   int64  `json:"product_id"`

	VariantID int64  `json:"variant_id"`
	Quantity  int64  `json:"quantity"`
	Price     int64  `json:"price"`
	Code      string `json:"code"`

	ImageUrl   string                           `json:"image_url"`
	Attributes []*catalogmodel.ProductAttribute `json:"attributes"`
}

type PurchaseOrderSupplier struct {
	FullName           string `json:"full_name"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	CompanyName        string `json:"company_name"`
	TaxNumber          string `json:"tax_number"`
	HeadquarterAddress string `json:"headquarter_address"`
}
