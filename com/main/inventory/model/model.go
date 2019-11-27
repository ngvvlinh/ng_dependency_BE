package model

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenInventoryVariant(&InventoryVariant{})

type InventoryVariant struct {
	ShopID         dot.ID
	VariantID      dot.ID
	QuantityOnHand int
	QuantityPicked int
	CostPrice      int

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenInventoryVoucher(&InventoryVoucher{})

type InventoryVoucher struct {
	ShopID       dot.ID
	ID           dot.ID
	CreatedBy    dot.ID
	UpdatedBy    dot.ID
	Code         string
	CodeNorm     int
	Status       etop.Status3
	Note         string
	TraderID     dot.ID
	Trader       *Trader
	TotalAmount  int
	Type         string
	Lines        []*InventoryVoucherItem
	VariantIDs   []dot.ID
	RefID        dot.ID
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
	ProductID   dot.ID `json:"product_id"`
	VariantID   dot.ID `json:"variant_id"`
	VariantName string `json:"variant_name"`

	Price    int `json:"price"`
	Quantity int `json:"quantity"`

	Code       string       `json:"code"`
	ImageURL   string       `json:"image_url"`
	Attributes []*Attribute `json:"attributes"`
}

type Trader struct {
	ID       dot.ID `json:"id"`
	Type     string `json:"type"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
