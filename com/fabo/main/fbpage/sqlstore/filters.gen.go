// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	status3 "etop.vn/api/top/types/etc/status3"
	sq "etop.vn/backend/pkg/common/sql/sq"
	dot "etop.vn/capi/dot"
)

type FbPageFilters struct{ prefix string }

func NewFbPageFilters(prefix string) FbPageFilters {
	return FbPageFilters{prefix}
}

func (ft *FbPageFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft FbPageFilters) Prefix() string {
	return ft.prefix
}

func (ft *FbPageFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *FbPageFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *FbPageFilters) ByExternalID(ExternalID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == "",
	}
}

func (ft *FbPageFilters) ByExternalIDPtr(ExternalID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == nil,
		IsZero: ExternalID != nil && (*ExternalID) == "",
	}
}

func (ft *FbPageFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *FbPageFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *FbPageFilters) ByUserID(UserID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == 0,
	}
}

func (ft *FbPageFilters) ByUserIDPtr(UserID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == nil,
		IsZero: UserID != nil && (*UserID) == 0,
	}
}

func (ft *FbPageFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *FbPageFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *FbPageFilters) ByCategory(Category string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "category",
		Value:  Category,
		IsNil:  Category == "",
	}
}

func (ft *FbPageFilters) ByCategoryPtr(Category *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "category",
		Value:  Category,
		IsNil:  Category == nil,
		IsZero: Category != nil && (*Category) == "",
	}
}

func (ft *FbPageFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *FbPageFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *FbPageFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *FbPageFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *FbPageFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *FbPageFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *FbPageFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *FbPageFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

type FbPageInternalFilters struct{ prefix string }

func NewFbPageInternalFilters(prefix string) FbPageInternalFilters {
	return FbPageInternalFilters{prefix}
}

func (ft *FbPageInternalFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft FbPageInternalFilters) Prefix() string {
	return ft.prefix
}

func (ft *FbPageInternalFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *FbPageInternalFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *FbPageInternalFilters) ByToken(Token string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "token",
		Value:  Token,
		IsNil:  Token == "",
	}
}

func (ft *FbPageInternalFilters) ByTokenPtr(Token *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "token",
		Value:  Token,
		IsNil:  Token == nil,
		IsZero: Token != nil && (*Token) == "",
	}
}

func (ft *FbPageInternalFilters) ByExpiresIn(ExpiresIn int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "expires_in",
		Value:  ExpiresIn,
		IsNil:  ExpiresIn == 0,
	}
}

func (ft *FbPageInternalFilters) ByExpiresInPtr(ExpiresIn *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "expires_in",
		Value:  ExpiresIn,
		IsNil:  ExpiresIn == nil,
		IsZero: ExpiresIn != nil && (*ExpiresIn) == 0,
	}
}

func (ft *FbPageInternalFilters) ByUpdateAt(UpdateAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "update_at",
		Value:  UpdateAt,
		IsNil:  UpdateAt.IsZero(),
	}
}

func (ft *FbPageInternalFilters) ByUpdateAtPtr(UpdateAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "update_at",
		Value:  UpdateAt,
		IsNil:  UpdateAt == nil,
		IsZero: UpdateAt != nil && (*UpdateAt).IsZero(),
	}
}
