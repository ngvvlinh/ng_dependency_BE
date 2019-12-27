package model

import (
	"time"

	"etop.vn/api/top/types/etc/inventory_type"
	"etop.vn/api/top/types/etc/inventory_voucher_ref"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/catalog/model"
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
	Status       status3.Status
	Note         string
	TraderID     dot.ID
	Trader       *Trader
	TotalAmount  int
	Type         inventory_type.InventoryVoucherType
	Lines        []*InventoryVoucherItem
	VariantIDs   []dot.ID
	RefID        dot.ID
	RefCode      string
	RefType      inventory_voucher_ref.InventoryVoucherRef
	RefName      string
	Title        string
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	ConfirmedAt  time.Time
	CancelledAt  time.Time
	CancelReason string
	ProductIDs   []dot.ID
	Rollback     bool
}

type InventoryVoucherItem struct {
	ProductName string `json:"product_name"`
	ProductID   dot.ID `json:"product_id"`
	VariantID   dot.ID `json:"variant_id"`
	VariantName string `json:"variant_name"`

	Price    int `json:"price"`
	Quantity int `json:"quantity"`

	Code       string                    `json:"code"`
	ImageURL   string                    `json:"image_url"`
	Attributes []*model.ProductAttribute `json:"attributes"`
}

type Trader struct {
	ID       dot.ID `json:"id"`
	Type     string `json:"type"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}
