package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	catalogmodel "o.o/backend/com/main/catalog/model"
	"o.o/capi/dot"
)

// +sqlgen
type PurchaseRefund struct {
	ID              dot.ID
	ShopID          dot.ID
	PurchaseOrderID dot.ID
	Code            string
	CodeNorm        int
	Note            string
	Lines           []*PurchaseRefundLine
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CancelledAt     time.Time
	ConfirmedAt     time.Time
	CreatedBy       dot.ID
	UpdatedBy       dot.ID
	CancelReason    string
	Status          status3.Status `sql_type:"int4"`
	SupplierID      dot.ID
	TotalAmount     int
	BasketValue     int
	Rid             dot.ID
}

type PurchaseRefundLine struct {
	VariantID    dot.ID                           `json:"variant_id"`
	Quantity     int                              `json:"count"`
	Code         string                           `json:"code"`
	ImageURL     string                           `json:"image_url"`
	ProductName  string                           `json:"product_name"`
	PaymentPrice int                              `json:"payment_price"`
	ProductID    dot.ID                           `json:"product_id"`
	Attributes   []*catalogmodel.ProductAttribute `json:"attributes"`
	Adjustment   int                              `json:"adjustment"`
}
