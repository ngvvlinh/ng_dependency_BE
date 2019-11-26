package model

import (
	"time"

	"etop.vn/api/main/etop"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenPurchaseOrder(&PurchaseOrder{})

type PurchaseOrder struct {
	ID              dot.ID
	ShopID          dot.ID
	SupplierID      dot.ID
	Supplier        *PurchaseOrderSupplier
	BasketValue     int64
	TotalDiscount   int64
	TotalAmount     int64
	Code            string
	CodeNorm        int
	Note            string
	Status          etop.Status3
	VariantIDs      []dot.ID
	Lines           []*PurchaseOrderLine
	CreatedBy       dot.ID
	CancelledReason string
	ConfirmedAt     time.Time
	CancelledAt     time.Time
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
	DeletedAt       time.Time

	SupplierFullNameNorm string
	SupplierPhoneNorm    string
}

type PurchaseOrderLine struct {
	ProductName string `json:"product_name"`
	ProductID   dot.ID `json:"product_id"`

	VariantID    dot.ID `json:"variant_id"`
	Quantity     int64  `json:"quantity"`
	PaymentPrice int64  `json:"payment_price"`
	Code         string `json:"code"`

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
