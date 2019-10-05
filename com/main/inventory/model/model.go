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
	PurchasePrice  int32

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenInventoryVoucher(&InventoryVoucher{})

type InventoryVoucher struct {
	ShopID      int64
	ID          int64
	CreatedBy   int64
	UpdatedBy   int64
	Status      etop.Status3
	Note        string
	TraderID    int64
	TotalAmount int32
	Type        string
	Lines       []*InventoryVoucherItem
	Title       string
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	ConfirmedAt time.Time
	CancelledAt time.Time

	CancelledReason string
}

type InventoryVoucherItem struct {
	VariantID int64 `json:"variant_id"`
	Price     int32 `json:"price"`
	Quantity  int32 `json:"quantity"`
	Discount  int32 `json:"discount"`
}
