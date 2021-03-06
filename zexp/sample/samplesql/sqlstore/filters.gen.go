// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
)

type AccountFilters struct{ prefix string }

func NewAccountFilters(prefix string) AccountFilters {
	return AccountFilters{prefix}
}

func (ft *AccountFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft AccountFilters) Prefix() string {
	return ft.prefix
}

func (ft *AccountFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *AccountFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *AccountFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *AccountFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

type FooFilters struct{ prefix string }

func NewFooFilters(prefix string) FooFilters {
	return FooFilters{prefix}
}

func (ft *FooFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft FooFilters) Prefix() string {
	return ft.prefix
}

func (ft *FooFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *FooFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *FooFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *FooFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *FooFilters) ByABC(ABC string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "abc_2",
		Value:  ABC,
		IsNil:  ABC == "",
	}
}

func (ft *FooFilters) ByABCPtr(ABC *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "abc_2",
		Value:  ABC,
		IsNil:  ABC == nil,
		IsZero: ABC != nil && (*ABC) == "",
	}
}

func (ft *FooFilters) ByDef2(Def2 string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "def2",
		Value:  Def2,
		IsNil:  Def2 == "",
	}
}

func (ft *FooFilters) ByDef2Ptr(Def2 *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "def2",
		Value:  Def2,
		IsNil:  Def2 == nil,
		IsZero: Def2 != nil && (*Def2) == "",
	}
}

func (ft *FooFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *FooFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *FooFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *FooFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}
