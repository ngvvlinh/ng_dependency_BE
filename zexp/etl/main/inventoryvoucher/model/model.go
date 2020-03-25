package model

import (
	"time"

	"etop.vn/api/top/types/etc/inventory_type"
	"etop.vn/api/top/types/etc/inventory_voucher_ref"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/catalog/model"
	"etop.vn/capi/dot"
)

// +sqlgen
type InventoryVoucher struct {
	ShopID       dot.ID
	ID           dot.ID
	CreatedBy    dot.ID
	UpdatedBy    dot.ID
	Code         string
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
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
