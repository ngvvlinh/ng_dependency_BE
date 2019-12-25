// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/top/types/etc/payment_provider"
	"etop.vn/api/top/types/etc/payment_state"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi/dot"
)

type PaymentFilters struct{ prefix string }

func NewPaymentFilters(prefix string) PaymentFilters {
	return PaymentFilters{prefix}
}

func (ft *PaymentFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft PaymentFilters) Prefix() string {
	return ft.prefix
}

func (ft *PaymentFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *PaymentFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *PaymentFilters) ByAmount(Amount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "amount",
		Value:  Amount,
		IsNil:  Amount == 0,
	}
}

func (ft *PaymentFilters) ByAmountPtr(Amount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "amount",
		Value:  Amount,
		IsNil:  Amount == nil,
		IsZero: Amount != nil && (*Amount) == 0,
	}
}

func (ft *PaymentFilters) ByStatus(Status status4.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *PaymentFilters) ByStatusPtr(Status *status4.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *PaymentFilters) ByState(State payment_state.PaymentState) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "state",
		Value:  State,
		IsNil:  State == 0,
	}
}

func (ft *PaymentFilters) ByStatePtr(State *payment_state.PaymentState) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "state",
		Value:  State,
		IsNil:  State == nil,
		IsZero: State != nil && (*State) == 0,
	}
}

func (ft *PaymentFilters) ByPaymentProvider(PaymentProvider payment_provider.PaymentProvider) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "payment_provider",
		Value:  PaymentProvider,
		IsNil:  PaymentProvider == 0,
	}
}

func (ft *PaymentFilters) ByPaymentProviderPtr(PaymentProvider *payment_provider.PaymentProvider) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "payment_provider",
		Value:  PaymentProvider,
		IsNil:  PaymentProvider == nil,
		IsZero: PaymentProvider != nil && (*PaymentProvider) == 0,
	}
}

func (ft *PaymentFilters) ByExternalTransID(ExternalTransID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_trans_id",
		Value:  ExternalTransID,
		IsNil:  ExternalTransID == "",
	}
}

func (ft *PaymentFilters) ByExternalTransIDPtr(ExternalTransID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_trans_id",
		Value:  ExternalTransID,
		IsNil:  ExternalTransID == nil,
		IsZero: ExternalTransID != nil && (*ExternalTransID) == "",
	}
}

func (ft *PaymentFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *PaymentFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *PaymentFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *PaymentFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}
