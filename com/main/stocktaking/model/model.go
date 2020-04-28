package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/stocktake_type"
	catalogmodel "o.o/backend/com/main/catalog/model"
	"o.o/capi/dot"
)

// +sqlgen
type ShopStocktake struct {
	ID            dot.ID
	ShopID        dot.ID
	TotalQuantity int
	CreatedBy     dot.ID
	UpdatedBy     dot.ID
	CancelReason  string
	Type          stocktake_type.StocktakeType
	Code          string
	CodeNorm      int
	Status        status3.Status
	CreatedAt     time.Time `sq:"create"`
	UpdatedAt     time.Time `sq:"update"`
	ConfirmedAt   time.Time
	CancelledAt   time.Time
	Lines         []*StocktakeLine
	Note          string
	ProductIDs    []dot.ID

	Rid dot.ID
}

type StocktakeLine struct {
	ProductName string                           `json:"product_name"`
	ProductID   dot.ID                           `json:"product_id"`
	VariantID   dot.ID                           `json:"variant_id"`
	OldQuantity int                              `json:"old_quantity"`
	NewQuantity int                              `json:"new_quantity"`
	VariantName string                           `json:"variant_name"`
	Name        string                           `json:"name"`
	Code        string                           `json:"code"`
	ImageURL    string                           `json:"image_url"`
	Attributes  []*catalogmodel.ProductAttribute `json:"attributes"`
	CostPrice   int                              `json:"cost_price"`
}
