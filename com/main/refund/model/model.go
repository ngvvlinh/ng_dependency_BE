package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/com/main/identity/sharemodel"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenRefund(&Refund{})

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
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
	CancelledAt     time.Time
	ConfirmedAt     time.Time
	CreatedBy       dot.ID
	UpdatedBy       dot.ID
	CancelReason    string
	Status          status3.Status
	CustomerID      dot.ID
	TotalAmount     int
	BasketValue     int
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
