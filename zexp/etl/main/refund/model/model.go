package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/com/main/identity/sharemodel"
	"etop.vn/capi/dot"
)

// +sqlgen
type Refund struct {
	ID              dot.ID
	ShopID          dot.ID
	OrderID         dot.ID
	Code            string
	CodeNorm        int
	Note            string
	Lines           []*RefundLine
	AdjustmentLines []*sharemodel.AdjustmentLine
	TotalAdjustment int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CancelledAt     time.Time
	ConfirmedAt     time.Time
	CreatedBy       dot.ID
	UpdatedBy       dot.ID
	CancelReason    string
	Status          status3.Status `sql_type:"int4"`
	CustomerID      dot.ID
	TotalAmount     int
	BasketValue     int

	Rid dot.ID
}

type RefundLine struct {
	VariantID   dot.ID                           `json:"variant_id"`
	Quantity    int                              `json:"count"`
	Code        string                           `json:"code"`
	ImageURL    string                           `json:"image_url"`
	ProductName string                           `json:"product_name"`
	RetailPrice int                              `json:"retail_price"`
	ProductID   dot.ID                           `json:"product_id"`
	Attributes  []*catalogmodel.ProductAttribute `json:"attributes"`
	Adjustment  int                              `json:"adjustment"`
}
