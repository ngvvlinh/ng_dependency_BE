package model

import (
	"time"

	"etop.vn/api/main/etop"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenInventoryVariant(&InventoryVariant{})

type InventoryVariant struct {
	ShopID         int64
	VariantID      int64
	QuantityOnHand int32
	QuantityPicked int32
	CostPrice      int32

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenInventoryVoucher(&InventoryVoucher{})

type InventoryVoucher struct {
	ShopID       int64
	ID           int64
	CreatedBy    int64
	UpdatedBy    int64
	Code         string
	CodeNorm     int32
	Status       etop.Status3
	Note         string
	TraderID     int64
	Trader       *Trader
	TotalAmount  int32
	Type         string
	Lines        []*InventoryVoucherItem
	VariantIDs   []int64
	RefID        int64
	RefCode      string
	RefType      string
	RefName      string
	Title        string
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	ConfirmedAt  time.Time
	CancelledAt  time.Time
	CancelReason string
}

type InventoryVoucherItem struct {
	ProductName string `json:"product_name"`
	ProductID   int64  `json:"product_id"`
	VariantID   int64  `json:"variant_id"`
	VariantName string `json:"variant_name"`

	Price    int32 `json:"price"`
	Quantity int32 `json:"quantity"`

	Code       string       `json:"code"`
	ImageURL   string       `json:"image_url"`
	Attributes []*Attribute `json:"attributes"`
}

type Trader struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
