// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/backend/pkg/common/sq"
)

type InventoryVariantFilters struct{ prefix string }

func NewInventoryVariantFilters(prefix string) InventoryVariantFilters {
	return InventoryVariantFilters{prefix}
}

func (ft *InventoryVariantFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft InventoryVariantFilters) Prefix() string {
	return ft.prefix
}

func (ft *InventoryVariantFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *InventoryVariantFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *InventoryVariantFilters) ByVariantID(VariantID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "variant_id",
		Value:  VariantID,
		IsNil:  VariantID == 0,
	}
}

func (ft *InventoryVariantFilters) ByVariantIDPtr(VariantID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "variant_id",
		Value:  VariantID,
		IsNil:  VariantID == nil,
		IsZero: VariantID != nil && (*VariantID) == 0,
	}
}

func (ft *InventoryVariantFilters) ByQuantityOnHand(QuantityOnHand int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "quantity_on_hand",
		Value:  QuantityOnHand,
		IsNil:  QuantityOnHand == 0,
	}
}

func (ft *InventoryVariantFilters) ByQuantityOnHandPtr(QuantityOnHand *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "quantity_on_hand",
		Value:  QuantityOnHand,
		IsNil:  QuantityOnHand == nil,
		IsZero: QuantityOnHand != nil && (*QuantityOnHand) == 0,
	}
}

func (ft *InventoryVariantFilters) ByQuantityPicked(QuantityPicked int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "quantity_picked",
		Value:  QuantityPicked,
		IsNil:  QuantityPicked == 0,
	}
}

func (ft *InventoryVariantFilters) ByQuantityPickedPtr(QuantityPicked *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "quantity_picked",
		Value:  QuantityPicked,
		IsNil:  QuantityPicked == nil,
		IsZero: QuantityPicked != nil && (*QuantityPicked) == 0,
	}
}

func (ft *InventoryVariantFilters) ByPurchasePrice(PurchasePrice int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "purchase_price",
		Value:  PurchasePrice,
		IsNil:  PurchasePrice == 0,
	}
}

func (ft *InventoryVariantFilters) ByPurchasePricePtr(PurchasePrice *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "purchase_price",
		Value:  PurchasePrice,
		IsNil:  PurchasePrice == nil,
		IsZero: PurchasePrice != nil && (*PurchasePrice) == 0,
	}
}

func (ft *InventoryVariantFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *InventoryVariantFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *InventoryVariantFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *InventoryVariantFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type InventoryVoucherFilters struct{ prefix string }

func NewInventoryVoucherFilters(prefix string) InventoryVoucherFilters {
	return InventoryVoucherFilters{prefix}
}

func (ft *InventoryVoucherFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft InventoryVoucherFilters) Prefix() string {
	return ft.prefix
}

func (ft *InventoryVoucherFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *InventoryVoucherFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByID(ID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *InventoryVoucherFilters) ByIDPtr(ID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByCreatedBy(CreatedBy int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == 0,
	}
}

func (ft *InventoryVoucherFilters) ByCreatedByPtr(CreatedBy *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == nil,
		IsZero: CreatedBy != nil && (*CreatedBy) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByUpdatedBy(UpdatedBy int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_by",
		Value:  UpdatedBy,
		IsNil:  UpdatedBy == 0,
	}
}

func (ft *InventoryVoucherFilters) ByUpdatedByPtr(UpdatedBy *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_by",
		Value:  UpdatedBy,
		IsNil:  UpdatedBy == nil,
		IsZero: UpdatedBy != nil && (*UpdatedBy) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByStatus(Status etop.Status3) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *InventoryVoucherFilters) ByStatusPtr(Status *etop.Status3) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *InventoryVoucherFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *InventoryVoucherFilters) ByTraderID(TraderID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "trader_id",
		Value:  TraderID,
		IsNil:  TraderID == 0,
	}
}

func (ft *InventoryVoucherFilters) ByTraderIDPtr(TraderID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "trader_id",
		Value:  TraderID,
		IsNil:  TraderID == nil,
		IsZero: TraderID != nil && (*TraderID) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByTotalAmount(TotalAmount int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == 0,
	}
}

func (ft *InventoryVoucherFilters) ByTotalAmountPtr(TotalAmount *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == nil,
		IsZero: TotalAmount != nil && (*TotalAmount) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByType(Type string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == "",
	}
}

func (ft *InventoryVoucherFilters) ByTypePtr(Type *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == "",
	}
}

func (ft *InventoryVoucherFilters) ByTitle(Title string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "title",
		Value:  Title,
		IsNil:  Title == "",
	}
}

func (ft *InventoryVoucherFilters) ByTitlePtr(Title *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "title",
		Value:  Title,
		IsNil:  Title == nil,
		IsZero: Title != nil && (*Title) == "",
	}
}

func (ft *InventoryVoucherFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *InventoryVoucherFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *InventoryVoucherFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *InventoryVoucherFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *InventoryVoucherFilters) ByConfirmedAt(ConfirmedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt.IsZero(),
	}
}

func (ft *InventoryVoucherFilters) ByConfirmedAtPtr(ConfirmedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt == nil,
		IsZero: ConfirmedAt != nil && (*ConfirmedAt).IsZero(),
	}
}

func (ft *InventoryVoucherFilters) ByCancelledAt(CancelledAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancelled_at",
		Value:  CancelledAt,
		IsNil:  CancelledAt.IsZero(),
	}
}

func (ft *InventoryVoucherFilters) ByCancelledAtPtr(CancelledAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancelled_at",
		Value:  CancelledAt,
		IsNil:  CancelledAt == nil,
		IsZero: CancelledAt != nil && (*CancelledAt).IsZero(),
	}
}

func (ft *InventoryVoucherFilters) ByCancelledReason(CancelledReason string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancelled_reason",
		Value:  CancelledReason,
		IsNil:  CancelledReason == "",
	}
}

func (ft *InventoryVoucherFilters) ByCancelledReasonPtr(CancelledReason *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancelled_reason",
		Value:  CancelledReason,
		IsNil:  CancelledReason == nil,
		IsZero: CancelledReason != nil && (*CancelledReason) == "",
	}
}