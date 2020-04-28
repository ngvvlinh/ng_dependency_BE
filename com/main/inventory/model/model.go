package model

import (
	"time"

	"o.o/api/top/types/etc/inventory_type"
	"o.o/api/top/types/etc/inventory_voucher_ref"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/catalog/model"
	"o.o/capi/dot"
)

// +sqlgen
type InventoryVariant struct {
	ShopID         dot.ID
	VariantID      dot.ID
	QuantityOnHand int
	QuantityPicked int
	CostPrice      int

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	Rid dot.ID
}

// +sqlgen
type InventoryVoucher struct {
	ShopID       dot.ID
	ID           dot.ID
	CreatedBy    dot.ID
	UpdatedBy    dot.ID
	Code         string
	CodeNorm     int
	Status       status3.Status
	TraderID     dot.ID
	Trader       *Trader
	TotalAmount  int
	Type         inventory_type.InventoryVoucherType
	Lines        []*InventoryVoucherItem
	VariantIDs   []dot.ID
	RefID        dot.ID
	RefCode      string
	RefType      inventory_voucher_ref.InventoryVoucherRef
	Title        string
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	ConfirmedAt  time.Time
	CancelledAt  time.Time
	CancelReason string
	ProductIDs   []dot.ID

	Rid dot.ID
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
