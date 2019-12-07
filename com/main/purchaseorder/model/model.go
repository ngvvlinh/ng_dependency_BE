package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"

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
	BasketValue     int
	TotalDiscount   int
	TotalAmount     int
	Code            string
	CodeNorm        int
	Note            string
	Status          status3.Status
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
	Quantity     int    `json:"quantity"`
	PaymentPrice int    `json:"payment_price"`
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
