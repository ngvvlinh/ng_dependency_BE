// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	status3 "o.o/api/top/types/etc/status3"
	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
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

func (ft *PurchaseOrderFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *PurchaseOrderFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *PurchaseOrderFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *PurchaseOrderFilters) BySupplierID(SupplierID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "supplier_id",
		Value:  SupplierID,
		IsNil:  SupplierID == 0,
	}
}

func (ft *PurchaseOrderFilters) BySupplierIDPtr(SupplierID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "supplier_id",
		Value:  SupplierID,
		IsNil:  SupplierID == nil,
		IsZero: SupplierID != nil && (*SupplierID) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByBasketValue(BasketValue int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "basket_value",
		Value:  BasketValue,
		IsNil:  BasketValue == 0,
	}
}

func (ft *PurchaseOrderFilters) ByBasketValuePtr(BasketValue *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "basket_value",
		Value:  BasketValue,
		IsNil:  BasketValue == nil,
		IsZero: BasketValue != nil && (*BasketValue) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalDiscount(TotalDiscount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_discount",
		Value:  TotalDiscount,
		IsNil:  TotalDiscount == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalDiscountPtr(TotalDiscount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_discount",
		Value:  TotalDiscount,
		IsNil:  TotalDiscount == nil,
		IsZero: TotalDiscount != nil && (*TotalDiscount) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalFee(TotalFee int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_fee",
		Value:  TotalFee,
		IsNil:  TotalFee == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalFeePtr(TotalFee *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_fee",
		Value:  TotalFee,
		IsNil:  TotalFee == nil,
		IsZero: TotalFee != nil && (*TotalFee) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalAmount(TotalAmount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == 0,
	}
}

func (ft *PurchaseOrderFilters) ByTotalAmountPtr(TotalAmount *int) *sq.ColumnFilterPtr {
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

func (ft *PurchaseOrderFilters) ByCodeNorm(CodeNorm int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code_norm",
		Value:  CodeNorm,
		IsNil:  CodeNorm == 0,
	}
}

func (ft *PurchaseOrderFilters) ByCodeNormPtr(CodeNorm *int) *sq.ColumnFilterPtr {
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

func (ft *PurchaseOrderFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *PurchaseOrderFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *PurchaseOrderFilters) ByCreatedBy(CreatedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == 0,
	}
}

func (ft *PurchaseOrderFilters) ByCreatedByPtr(CreatedBy *dot.ID) *sq.ColumnFilterPtr {
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

func (ft *PurchaseOrderFilters) BySupplierFullNameNorm(SupplierFullNameNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "supplier_full_name_norm",
		Value:  SupplierFullNameNorm,
		IsNil:  SupplierFullNameNorm == "",
	}
}

func (ft *PurchaseOrderFilters) BySupplierFullNameNormPtr(SupplierFullNameNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "supplier_full_name_norm",
		Value:  SupplierFullNameNorm,
		IsNil:  SupplierFullNameNorm == nil,
		IsZero: SupplierFullNameNorm != nil && (*SupplierFullNameNorm) == "",
	}
}

func (ft *PurchaseOrderFilters) BySupplierPhoneNorm(SupplierPhoneNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "supplier_phone_norm",
		Value:  SupplierPhoneNorm,
		IsNil:  SupplierPhoneNorm == "",
	}
}

func (ft *PurchaseOrderFilters) BySupplierPhoneNormPtr(SupplierPhoneNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "supplier_phone_norm",
		Value:  SupplierPhoneNorm,
		IsNil:  SupplierPhoneNorm == nil,
		IsZero: SupplierPhoneNorm != nil && (*SupplierPhoneNorm) == "",
	}
}

func (ft *PurchaseOrderFilters) ByRid(Rid dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "rid",
		Value:  Rid,
		IsNil:  Rid == 0,
	}
}

func (ft *PurchaseOrderFilters) ByRidPtr(Rid *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "rid",
		Value:  Rid,
		IsNil:  Rid == nil,
		IsZero: Rid != nil && (*Rid) == 0,
	}
}
