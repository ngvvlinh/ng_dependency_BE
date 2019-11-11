// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/backend/pkg/common/sq"
)

type PurchaseOrderFilters struct{ prefix string }

func NewPurchaseOrderFilters(prefix string) PurchaseOrderFilters {
	return PurchaseOrderFilters{prefix}
}

func (ft *PurchaseOrderFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft PurchaseOrderFilters) Prefix() string {
	return ft.prefix
}

func (ft *PurchaseOrderFilters) ByID(ID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *PurchaseOrderFilters) ByIDPtr(ID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *PurchaseOrderFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *PurchaseOrderFilters) BySupplierID(SupplierID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "supplier_id",
		Value:  SupplierID,
		IsNil:  SupplierID == 0,
	}
}

func (ft *PurchaseOrderFilters) BySupplierIDPtr(SupplierID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "supplier_id",
		Value:  SupplierID,
		IsNil:  SupplierID == nil,
		IsZero: SupplierID != nil && (*SupplierID) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByBasketValue(BasketValue int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "basket_value",
		Value:  BasketValue,
		IsNil:  BasketValue == 0,
	}
}

func (ft *PurchaseOrderFilters) ByBasketValuePtr(BasketValue *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "basket_value",
		Value:  BasketValue,
		IsNil:  BasketValue == nil,
		IsZero: BasketValue != nil && (*BasketValue) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalDiscount(TotalDiscount int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_discount",
		Value:  TotalDiscount,
		IsNil:  TotalDiscount == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalDiscountPtr(TotalDiscount *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_discount",
		Value:  TotalDiscount,
		IsNil:  TotalDiscount == nil,
		IsZero: TotalDiscount != nil && (*TotalDiscount) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalAmount(TotalAmount int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalAmountPtr(TotalAmount *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == nil,
		IsZero: TotalAmount != nil && (*TotalAmount) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *PurchaseOrderFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *PurchaseOrderFilters) ByCodeNorm(CodeNorm int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code_norm",
		Value:  CodeNorm,
		IsNil:  CodeNorm == 0,
	}
}

func (ft *PurchaseOrderFilters) ByCodeNormPtr(CodeNorm *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code_norm",
		Value:  CodeNorm,
		IsNil:  CodeNorm == nil,
		IsZero: CodeNorm != nil && (*CodeNorm) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *PurchaseOrderFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *PurchaseOrderFilters) ByStatus(Status etop.Status3) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *PurchaseOrderFilters) ByStatusPtr(Status *etop.Status3) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByCreatedBy(CreatedBy int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == 0,
	}
}

func (ft *PurchaseOrderFilters) ByCreatedByPtr(CreatedBy *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == nil,
		IsZero: CreatedBy != nil && (*CreatedBy) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByCancelledReason(CancelledReason string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancelled_reason",
		Value:  CancelledReason,
		IsNil:  CancelledReason == "",
	}
}

func (ft *PurchaseOrderFilters) ByCancelledReasonPtr(CancelledReason *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancelled_reason",
		Value:  CancelledReason,
		IsNil:  CancelledReason == nil,
		IsZero: CancelledReason != nil && (*CancelledReason) == "",
	}
}

func (ft *PurchaseOrderFilters) ByConfirmedAt(ConfirmedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt.IsZero(),
	}
}

func (ft *PurchaseOrderFilters) ByConfirmedAtPtr(ConfirmedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt == nil,
		IsZero: ConfirmedAt != nil && (*ConfirmedAt).IsZero(),
	}
}

func (ft *PurchaseOrderFilters) ByCancelledAt(CancelledAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancelled_at",
		Value:  CancelledAt,
		IsNil:  CancelledAt.IsZero(),
	}
}

func (ft *PurchaseOrderFilters) ByCancelledAtPtr(CancelledAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancelled_at",
		Value:  CancelledAt,
		IsNil:  CancelledAt == nil,
		IsZero: CancelledAt != nil && (*CancelledAt).IsZero(),
	}
}

func (ft *PurchaseOrderFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *PurchaseOrderFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *PurchaseOrderFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *PurchaseOrderFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *PurchaseOrderFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *PurchaseOrderFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}
