// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	status3 "o.o/api/top/types/etc/status3"
	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
)

type SubscriptionFilters struct{ prefix string }

func NewSubscriptionFilters(prefix string) SubscriptionFilters {
	return SubscriptionFilters{prefix}
}

func (ft *SubscriptionFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft SubscriptionFilters) Prefix() string {
	return ft.prefix
}

func (ft *SubscriptionFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *SubscriptionFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *SubscriptionFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *SubscriptionFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *SubscriptionFilters) ByCancelAtPeriodEnd(CancelAtPeriodEnd bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancel_at_period_end",
		Value:  CancelAtPeriodEnd,
		IsNil:  bool(!CancelAtPeriodEnd),
	}
}

func (ft *SubscriptionFilters) ByCancelAtPeriodEndPtr(CancelAtPeriodEnd *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancel_at_period_end",
		Value:  CancelAtPeriodEnd,
		IsNil:  CancelAtPeriodEnd == nil,
		IsZero: CancelAtPeriodEnd != nil && bool(!(*CancelAtPeriodEnd)),
	}
}

func (ft *SubscriptionFilters) ByCurrentPeriodEndAt(CurrentPeriodEndAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "current_period_end_at",
		Value:  CurrentPeriodEndAt,
		IsNil:  CurrentPeriodEndAt.IsZero(),
	}
}

func (ft *SubscriptionFilters) ByCurrentPeriodEndAtPtr(CurrentPeriodEndAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "current_period_end_at",
		Value:  CurrentPeriodEndAt,
		IsNil:  CurrentPeriodEndAt == nil,
		IsZero: CurrentPeriodEndAt != nil && (*CurrentPeriodEndAt).IsZero(),
	}
}

func (ft *SubscriptionFilters) ByCurrentPeriodStartAt(CurrentPeriodStartAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "current_period_start_at",
		Value:  CurrentPeriodStartAt,
		IsNil:  CurrentPeriodStartAt.IsZero(),
	}
}

func (ft *SubscriptionFilters) ByCurrentPeriodStartAtPtr(CurrentPeriodStartAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "current_period_start_at",
		Value:  CurrentPeriodStartAt,
		IsNil:  CurrentPeriodStartAt == nil,
		IsZero: CurrentPeriodStartAt != nil && (*CurrentPeriodStartAt).IsZero(),
	}
}

func (ft *SubscriptionFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *SubscriptionFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *SubscriptionFilters) ByBillingCycleAnchorAt(BillingCycleAnchorAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "billing_cycle_anchor_at",
		Value:  BillingCycleAnchorAt,
		IsNil:  BillingCycleAnchorAt.IsZero(),
	}
}

func (ft *SubscriptionFilters) ByBillingCycleAnchorAtPtr(BillingCycleAnchorAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "billing_cycle_anchor_at",
		Value:  BillingCycleAnchorAt,
		IsNil:  BillingCycleAnchorAt == nil,
		IsZero: BillingCycleAnchorAt != nil && (*BillingCycleAnchorAt).IsZero(),
	}
}

func (ft *SubscriptionFilters) ByStartedAt(StartedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "started_at",
		Value:  StartedAt,
		IsNil:  StartedAt.IsZero(),
	}
}

func (ft *SubscriptionFilters) ByStartedAtPtr(StartedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "started_at",
		Value:  StartedAt,
		IsNil:  StartedAt == nil,
		IsZero: StartedAt != nil && (*StartedAt).IsZero(),
	}
}

func (ft *SubscriptionFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *SubscriptionFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *SubscriptionFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *SubscriptionFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *SubscriptionFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *SubscriptionFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *SubscriptionFilters) ByWLPartnerID(WLPartnerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "wl_partner_id",
		Value:  WLPartnerID,
		IsNil:  WLPartnerID == 0,
	}
}

func (ft *SubscriptionFilters) ByWLPartnerIDPtr(WLPartnerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "wl_partner_id",
		Value:  WLPartnerID,
		IsNil:  WLPartnerID == nil,
		IsZero: WLPartnerID != nil && (*WLPartnerID) == 0,
	}
}

type SubscriptionLineFilters struct{ prefix string }

func NewSubscriptionLineFilters(prefix string) SubscriptionLineFilters {
	return SubscriptionLineFilters{prefix}
}

func (ft *SubscriptionLineFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft SubscriptionLineFilters) Prefix() string {
	return ft.prefix
}

func (ft *SubscriptionLineFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *SubscriptionLineFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *SubscriptionLineFilters) ByPlanID(PlanID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "plan_id",
		Value:  PlanID,
		IsNil:  PlanID == 0,
	}
}

func (ft *SubscriptionLineFilters) ByPlanIDPtr(PlanID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "plan_id",
		Value:  PlanID,
		IsNil:  PlanID == nil,
		IsZero: PlanID != nil && (*PlanID) == 0,
	}
}

func (ft *SubscriptionLineFilters) BySubscriptionID(SubscriptionID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "subscription_id",
		Value:  SubscriptionID,
		IsNil:  SubscriptionID == 0,
	}
}

func (ft *SubscriptionLineFilters) BySubscriptionIDPtr(SubscriptionID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "subscription_id",
		Value:  SubscriptionID,
		IsNil:  SubscriptionID == nil,
		IsZero: SubscriptionID != nil && (*SubscriptionID) == 0,
	}
}

func (ft *SubscriptionLineFilters) ByQuantity(Quantity int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "quantity",
		Value:  Quantity,
		IsNil:  Quantity == 0,
	}
}

func (ft *SubscriptionLineFilters) ByQuantityPtr(Quantity *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "quantity",
		Value:  Quantity,
		IsNil:  Quantity == nil,
		IsZero: Quantity != nil && (*Quantity) == 0,
	}
}

func (ft *SubscriptionLineFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *SubscriptionLineFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *SubscriptionLineFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *SubscriptionLineFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}