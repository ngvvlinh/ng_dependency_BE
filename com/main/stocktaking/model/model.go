package model

import (
	"time"

	"etop.vn/backend/pkg/etop/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenStocktake(&ShopStocktake{})

type ShopStocktake struct {
	ID            int64
	ShopID        int64
	TotalQuantity int32
	CreatedBy     int64
	UpdatedBy     int64
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
	ProductID   int64        `json:"product_id"`
	VariantID   int64        `json:"variant_id"`
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
