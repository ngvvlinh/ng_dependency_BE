// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/top/types/etc/receipt_ref"
	"etop.vn/api/top/types/etc/receipt_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi/dot"
)

type ReceiptFilters struct{ prefix string }

func NewReceiptFilters(prefix string) ReceiptFilters {
	return ReceiptFilters{prefix}
}

func (ft *ReceiptFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ReceiptFilters) Prefix() string {
	return ft.prefix
}

func (ft *ReceiptFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ReceiptFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ReceiptFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ReceiptFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ReceiptFilters) ByTraderID(TraderID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "trader_id",
		Value:  TraderID,
		IsNil:  TraderID == 0,
	}
}

func (ft *ReceiptFilters) ByTraderIDPtr(TraderID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "trader_id",
		Value:  TraderID,
		IsNil:  TraderID == nil,
		IsZero: TraderID != nil && (*TraderID) == 0,
	}
}

func (ft *ReceiptFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *ReceiptFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *ReceiptFilters) ByCodeNorm(CodeNorm int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code_norm",
		Value:  CodeNorm,
		IsNil:  CodeNorm == 0,
	}
}

func (ft *ReceiptFilters) ByCodeNormPtr(CodeNorm *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code_norm",
		Value:  CodeNorm,
		IsNil:  CodeNorm == nil,
		IsZero: CodeNorm != nil && (*CodeNorm) == 0,
	}
}

func (ft *ReceiptFilters) ByTitle(Title string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "title",
		Value:  Title,
		IsNil:  Title == "",
	}
}

func (ft *ReceiptFilters) ByTitlePtr(Title *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "title",
		Value:  Title,
		IsNil:  Title == nil,
		IsZero: Title != nil && (*Title) == "",
	}
}

func (ft *ReceiptFilters) ByType(Type receipt_type.ReceiptType) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == 0,
	}
}

func (ft *ReceiptFilters) ByTypePtr(Type *receipt_type.ReceiptType) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == 0,
	}
}

func (ft *ReceiptFilters) ByDescription(Description string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == "",
	}
}

func (ft *ReceiptFilters) ByDescriptionPtr(Description *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == nil,
		IsZero: Description != nil && (*Description) == "",
	}
}

func (ft *ReceiptFilters) ByTraderFullNameNorm(TraderFullNameNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "trader_full_name_norm",
		Value:  TraderFullNameNorm,
		IsNil:  TraderFullNameNorm == "",
	}
}

func (ft *ReceiptFilters) ByTraderFullNameNormPtr(TraderFullNameNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "trader_full_name_norm",
		Value:  TraderFullNameNorm,
		IsNil:  TraderFullNameNorm == nil,
		IsZero: TraderFullNameNorm != nil && (*TraderFullNameNorm) == "",
	}
}

func (ft *ReceiptFilters) ByTraderPhoneNorm(TraderPhoneNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "trader_phone_norm",
		Value:  TraderPhoneNorm,
		IsNil:  TraderPhoneNorm == "",
	}
}

func (ft *ReceiptFilters) ByTraderPhoneNormPtr(TraderPhoneNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "trader_phone_norm",
		Value:  TraderPhoneNorm,
		IsNil:  TraderPhoneNorm == nil,
		IsZero: TraderPhoneNorm != nil && (*TraderPhoneNorm) == "",
	}
}

func (ft *ReceiptFilters) ByTraderType(TraderType string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "trader_type",
		Value:  TraderType,
		IsNil:  TraderType == "",
	}
}

func (ft *ReceiptFilters) ByTraderTypePtr(TraderType *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "trader_type",
		Value:  TraderType,
		IsNil:  TraderType == nil,
		IsZero: TraderType != nil && (*TraderType) == "",
	}
}

func (ft *ReceiptFilters) ByAmount(Amount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "amount",
		Value:  Amount,
		IsNil:  Amount == 0,
	}
}

func (ft *ReceiptFilters) ByAmountPtr(Amount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "amount",
		Value:  Amount,
		IsNil:  Amount == nil,
		IsZero: Amount != nil && (*Amount) == 0,
	}
}

func (ft *ReceiptFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ReceiptFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ReceiptFilters) ByRefType(RefType receipt_ref.ReceiptRef) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ref_type",
		Value:  RefType,
		IsNil:  RefType == 0,
	}
}

func (ft *ReceiptFilters) ByRefTypePtr(RefType *receipt_ref.ReceiptRef) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ref_type",
		Value:  RefType,
		IsNil:  RefType == nil,
		IsZero: RefType != nil && (*RefType) == 0,
	}
}

func (ft *ReceiptFilters) ByLedgerID(LedgerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ledger_id",
		Value:  LedgerID,
		IsNil:  LedgerID == 0,
	}
}

func (ft *ReceiptFilters) ByLedgerIDPtr(LedgerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ledger_id",
		Value:  LedgerID,
		IsNil:  LedgerID == nil,
		IsZero: LedgerID != nil && (*LedgerID) == 0,
	}
}

func (ft *ReceiptFilters) ByCancelledReason(CancelledReason string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancelled_reason",
		Value:  CancelledReason,
		IsNil:  CancelledReason == "",
	}
}

func (ft *ReceiptFilters) ByCancelledReasonPtr(CancelledReason *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancelled_reason",
		Value:  CancelledReason,
		IsNil:  CancelledReason == nil,
		IsZero: CancelledReason != nil && (*CancelledReason) == "",
	}
}

func (ft *ReceiptFilters) ByCreatedType(CreatedType string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_type",
		Value:  CreatedType,
		IsNil:  CreatedType == "",
	}
}

func (ft *ReceiptFilters) ByCreatedTypePtr(CreatedType *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_type",
		Value:  CreatedType,
		IsNil:  CreatedType == nil,
		IsZero: CreatedType != nil && (*CreatedType) == "",
	}
}

func (ft *ReceiptFilters) ByCreatedBy(CreatedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == 0,
	}
}

func (ft *ReceiptFilters) ByCreatedByPtr(CreatedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == nil,
		IsZero: CreatedBy != nil && (*CreatedBy) == 0,
	}
}

func (ft *ReceiptFilters) ByPaidAt(PaidAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "paid_at",
		Value:  PaidAt,
		IsNil:  PaidAt.IsZero(),
	}
}

func (ft *ReceiptFilters) ByPaidAtPtr(PaidAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "paid_at",
		Value:  PaidAt,
		IsNil:  PaidAt == nil,
		IsZero: PaidAt != nil && (*PaidAt).IsZero(),
	}
}

func (ft *ReceiptFilters) ByConfirmedAt(ConfirmedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt.IsZero(),
	}
}

func (ft *ReceiptFilters) ByConfirmedAtPtr(ConfirmedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt == nil,
		IsZero: ConfirmedAt != nil && (*ConfirmedAt).IsZero(),
	}
}

func (ft *ReceiptFilters) ByCancelledAt(CancelledAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancelled_at",
		Value:  CancelledAt,
		IsNil:  CancelledAt.IsZero(),
	}
}

func (ft *ReceiptFilters) ByCancelledAtPtr(CancelledAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancelled_at",
		Value:  CancelledAt,
		IsNil:  CancelledAt == nil,
		IsZero: CancelledAt != nil && (*CancelledAt).IsZero(),
	}
}

func (ft *ReceiptFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ReceiptFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ReceiptFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ReceiptFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ReceiptFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ReceiptFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}
