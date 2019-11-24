package model

import (
	"time"

	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenStocktake(&ShopStocktake{})

type ShopStocktake struct {
	ID            dot.ID
	ShopID        dot.ID
	TotalQuantity int32
	CreatedBy     dot.ID
	UpdatedBy     dot.ID
	CancelReason  string
	Code          string
	CodeNorm      int32
	Status        model.Status3
	CreatedAt     time.Time `sq:"create"`
	UpdatedAt     time.Time `sq:"update"`
	ConfirmedAt   time.Time
	CancelledAt   time.Time
	Lines         []*StocktakeLine
	Note          string
}

type StocktakeLine struct {
	ProductName string       `json:"product_name"`
	ProductID   dot.ID       `json:"product_id"`
	VariantID   dot.ID       `json:"variant_id"`
	OldQuantity int32        `json:"old_quantity"`
	NewQuantity int32        `json:"new_quantity"`
	VariantName string       `json:"variant_name"`
	Name        string       `json:"name"`
	Code        string       `json:"code"`
	ImageURL    string       `json:"image_url"`
	Attributes  []*Attribute `json:"attributes"`
	CostPrice   int32        `json:"cost_price"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
