// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	status3 "o.o/api/top/types/etc/status3"
	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
)

type FbExternalUserFilters struct{ prefix string }

func NewFbExternalUserFilters(prefix string) FbExternalUserFilters {
	return FbExternalUserFilters{prefix}
}

func (ft *FbExternalUserFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft FbExternalUserFilters) Prefix() string {
	return ft.prefix
}

func (ft *FbExternalUserFilters) ByExternalID(ExternalID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == "",
	}
}

func (ft *FbExternalUserFilters) ByExternalIDPtr(ExternalID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == nil,
		IsZero: ExternalID != nil && (*ExternalID) == "",
	}
}

func (ft *FbExternalUserFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *FbExternalUserFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *FbExternalUserFilters) ByExternalPageID(ExternalPageID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_page_id",
		Value:  ExternalPageID,
		IsNil:  ExternalPageID == "",
	}
}

func (ft *FbExternalUserFilters) ByExternalPageIDPtr(ExternalPageID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_page_id",
		Value:  ExternalPageID,
		IsNil:  ExternalPageID == nil,
		IsZero: ExternalPageID != nil && (*ExternalPageID) == "",
	}
}

func (ft *FbExternalUserFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *FbExternalUserFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *FbExternalUserFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *FbExternalUserFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type FbExternalUserInternalFilters struct{ prefix string }

func NewFbExternalUserInternalFilters(prefix string) FbExternalUserInternalFilters {
	return FbExternalUserInternalFilters{prefix}
}

func (ft *FbExternalUserInternalFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft FbExternalUserInternalFilters) Prefix() string {
	return ft.prefix
}

func (ft *FbExternalUserInternalFilters) ByExternalID(ExternalID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == "",
	}
}

func (ft *FbExternalUserInternalFilters) ByExternalIDPtr(ExternalID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == nil,
		IsZero: ExternalID != nil && (*ExternalID) == "",
	}
}

func (ft *FbExternalUserInternalFilters) ByToken(Token string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "token",
		Value:  Token,
		IsNil:  Token == "",
	}
}

func (ft *FbExternalUserInternalFilters) ByTokenPtr(Token *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "token",
		Value:  Token,
		IsNil:  Token == nil,
		IsZero: Token != nil && (*Token) == "",
	}
}

func (ft *FbExternalUserInternalFilters) ByExpiresIn(ExpiresIn int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "expires_in",
		Value:  ExpiresIn,
		IsNil:  ExpiresIn == 0,
	}
}

func (ft *FbExternalUserInternalFilters) ByExpiresInPtr(ExpiresIn *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "expires_in",
		Value:  ExpiresIn,
		IsNil:  ExpiresIn == nil,
		IsZero: ExpiresIn != nil && (*ExpiresIn) == 0,
	}
}

func (ft *FbExternalUserInternalFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *FbExternalUserInternalFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type FbExternalUserShopCustomerFilters struct{ prefix string }

func NewFbExternalUserShopCustomerFilters(prefix string) FbExternalUserShopCustomerFilters {
	return FbExternalUserShopCustomerFilters{prefix}
}

func (ft *FbExternalUserShopCustomerFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft FbExternalUserShopCustomerFilters) Prefix() string {
	return ft.prefix
}

func (ft *FbExternalUserShopCustomerFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByFbExternalUserID(FbExternalUserID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "fb_external_user_id",
		Value:  FbExternalUserID,
		IsNil:  FbExternalUserID == "",
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByFbExternalUserIDPtr(FbExternalUserID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "fb_external_user_id",
		Value:  FbExternalUserID,
		IsNil:  FbExternalUserID == nil,
		IsZero: FbExternalUserID != nil && (*FbExternalUserID) == "",
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByCustomerID(CustomerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "customer_id",
		Value:  CustomerID,
		IsNil:  CustomerID == 0,
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByCustomerIDPtr(CustomerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "customer_id",
		Value:  CustomerID,
		IsNil:  CustomerID == nil,
		IsZero: CustomerID != nil && (*CustomerID) == 0,
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *FbExternalUserShopCustomerFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}
