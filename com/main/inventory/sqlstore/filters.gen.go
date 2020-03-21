// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/top/types/etc/inventory_type"
	"etop.vn/api/top/types/etc/inventory_voucher_ref"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi/dot"
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

func (ft *InventoryVariantFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *InventoryVariantFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *InventoryVariantFilters) ByVariantID(VariantID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "variant_id",
		Value:  VariantID,
		IsNil:  VariantID == 0,
	}
}

func (ft *InventoryVariantFilters) ByVariantIDPtr(VariantID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "variant_id",
		Value:  VariantID,
		IsNil:  VariantID == nil,
		IsZero: VariantID != nil && (*VariantID) == 0,
	}
}

func (ft *InventoryVariantFilters) ByQuantityOnHand(QuantityOnHand int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "quantity_on_hand",
		Value:  QuantityOnHand,
		IsNil:  QuantityOnHand == 0,
	}
}

func (ft *InventoryVariantFilters) ByQuantityOnHandPtr(QuantityOnHand *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "quantity_on_hand",
		Value:  QuantityOnHand,
		IsNil:  QuantityOnHand == nil,
		IsZero: QuantityOnHand != nil && (*QuantityOnHand) == 0,
	}
}

func (ft *InventoryVariantFilters) ByQuantityPicked(QuantityPicked int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "quantity_picked",
		Value:  QuantityPicked,
		IsNil:  QuantityPicked == 0,
	}
}

func (ft *InventoryVariantFilters) ByQuantityPickedPtr(QuantityPicked *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "quantity_picked",
		Value:  QuantityPicked,
		IsNil:  QuantityPicked == nil,
		IsZero: QuantityPicked != nil && (*QuantityPicked) == 0,
	}
}

func (ft *InventoryVariantFilters) ByCostPrice(CostPrice int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cost_price",
		Value:  CostPrice,
		IsNil:  CostPrice == 0,
	}
}

func (ft *InventoryVariantFilters) ByCostPricePtr(CostPrice *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cost_price",
		Value:  CostPrice,
		IsNil:  CostPrice == nil,
		IsZero: CostPrice != nil && (*CostPrice) == 0,
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

func (ft *InventoryVoucherFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *InventoryVoucherFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *InventoryVoucherFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByCreatedBy(CreatedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == 0,
	}
}

func (ft *InventoryVoucherFilters) ByCreatedByPtr(CreatedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == nil,
		IsZero: CreatedBy != nil && (*CreatedBy) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByUpdatedBy(UpdatedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_by",
		Value:  UpdatedBy,
		IsNil:  UpdatedBy == 0,
	}
}

func (ft *InventoryVoucherFilters) ByUpdatedByPtr(UpdatedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_by",
		Value:  UpdatedBy,
		IsNil:  UpdatedBy == nil,
		IsZero: UpdatedBy != nil && (*UpdatedBy) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *InventoryVoucherFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *InventoryVoucherFilters) ByCodeNorm(CodeNorm int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code_norm",
		Value:  CodeNorm,
		IsNil:  CodeNorm == 0,
	}
}

func (ft *InventoryVoucherFilters) ByCodeNormPtr(CodeNorm *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code_norm",
		Value:  CodeNorm,
		IsNil:  CodeNorm == nil,
		IsZero: CodeNorm != nil && (*CodeNorm) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *InventoryVoucherFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByTraderID(TraderID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "trader_id",
		Value:  TraderID,
		IsNil:  TraderID == 0,
	}
}

func (ft *InventoryVoucherFilters) ByTraderIDPtr(TraderID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "trader_id",
		Value:  TraderID,
		IsNil:  TraderID == nil,
		IsZero: TraderID != nil && (*TraderID) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByTotalAmount(TotalAmount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == 0,
	}
}

func (ft *InventoryVoucherFilters) ByTotalAmountPtr(TotalAmount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == nil,
		IsZero: TotalAmount != nil && (*TotalAmount) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByType(Type inventory_type.InventoryVoucherType) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == 0,
	}
}

func (ft *InventoryVoucherFilters) ByTypePtr(Type *inventory_type.InventoryVoucherType) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByRefID(RefID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ref_id",
		Value:  RefID,
		IsNil:  RefID == 0,
	}
}

func (ft *InventoryVoucherFilters) ByRefIDPtr(RefID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ref_id",
		Value:  RefID,
		IsNil:  RefID == nil,
		IsZero: RefID != nil && (*RefID) == 0,
	}
}

func (ft *InventoryVoucherFilters) ByRefCode(RefCode string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ref_code",
		Value:  RefCode,
		IsNil:  RefCode == "",
	}
}

func (ft *InventoryVoucherFilters) ByRefCodePtr(RefCode *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ref_code",
		Value:  RefCode,
		IsNil:  RefCode == nil,
		IsZero: RefCode != nil && (*RefCode) == "",
	}
}

func (ft *InventoryVoucherFilters) ByRefType(RefType inventory_voucher_ref.InventoryVoucherRef) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ref_type",
		Value:  RefType,
		IsNil:  RefType == 0,
	}
}

func (ft *InventoryVoucherFilters) ByRefTypePtr(RefType *inventory_voucher_ref.InventoryVoucherRef) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ref_type",
		Value:  RefType,
		IsNil:  RefType == nil,
		IsZero: RefType != nil && (*RefType) == 0,
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

func (ft *InventoryVoucherFilters) ByCancelReason(CancelReason string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancel_reason",
		Value:  CancelReason,
		IsNil:  CancelReason == "",
	}
}

func (ft *InventoryVoucherFilters) ByCancelReasonPtr(CancelReason *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancel_reason",
		Value:  CancelReason,
		IsNil:  CancelReason == nil,
		IsZero: CancelReason != nil && (*CancelReason) == "",
	}
}
