// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/backend/pkg/common/sq"
	m "etop.vn/backend/tools/pkg/sqlgen/test"
	"etop.vn/capi/dot"
)

type UserFilters struct{ prefix string }

func NewUserFilters(prefix string) UserFilters {
	return UserFilters{prefix}
}

func (ft *UserFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft UserFilters) Prefix() string {
	return ft.prefix
}

func (ft *UserFilters) ByID(ID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == "",
	}
}

func (ft *UserFilters) ByIDPtr(ID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == "",
	}
}

func (ft *UserFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *UserFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *UserFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *UserFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *UserFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && false,
	}
}

func (ft *UserFilters) ByBool(Bool bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "bool",
		Value:  Bool,
		IsNil:  bool(!Bool),
	}
}

func (ft *UserFilters) ByBoolPtr(Bool *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "bool",
		Value:  Bool,
		IsNil:  Bool == nil,
		IsZero: Bool != nil && bool(!(*Bool)),
	}
}

func (ft *UserFilters) ByFloat64(Float64 float64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "float64",
		Value:  Float64,
		IsNil:  Float64 == 0,
	}
}

func (ft *UserFilters) ByFloat64Ptr(Float64 *float64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "float64",
		Value:  Float64,
		IsNil:  Float64 == nil,
		IsZero: Float64 != nil && (*Float64) == 0,
	}
}

func (ft *UserFilters) ByInt(Int int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "int",
		Value:  Int,
		IsNil:  Int == 0,
	}
}

func (ft *UserFilters) ByIntPtr(Int *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "int",
		Value:  Int,
		IsNil:  Int == nil,
		IsZero: Int != nil && (*Int) == 0,
	}
}

func (ft *UserFilters) ByInt64(Int64 int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "int64",
		Value:  Int64,
		IsNil:  Int64 == 0,
	}
}

func (ft *UserFilters) ByInt64Ptr(Int64 *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "int64",
		Value:  Int64,
		IsNil:  Int64 == nil,
		IsZero: Int64 != nil && (*Int64) == 0,
	}
}

func (ft *UserFilters) ByString(String string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "string",
		Value:  String,
		IsNil:  String == "",
	}
}

func (ft *UserFilters) ByStringPtr(String *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "string",
		Value:  String,
		IsNil:  String == nil,
		IsZero: String != nil && (*String) == "",
	}
}

func (ft *UserFilters) ByPBoolPtr(PBool *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_bool",
		Value:  PBool,
		IsNil:  PBool == nil,
		IsZero: PBool != nil && bool(!(*PBool)),
	}
}

func (ft *UserFilters) ByPFloat64Ptr(PFloat64 *float64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_float64",
		Value:  PFloat64,
		IsNil:  PFloat64 == nil,
		IsZero: PFloat64 != nil && (*PFloat64) == 0,
	}
}

func (ft *UserFilters) ByPIntPtr(PInt *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_int",
		Value:  PInt,
		IsNil:  PInt == nil,
		IsZero: PInt != nil && (*PInt) == 0,
	}
}

func (ft *UserFilters) ByPInt64Ptr(PInt64 *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_int64",
		Value:  PInt64,
		IsNil:  PInt64 == nil,
		IsZero: PInt64 != nil && (*PInt64) == 0,
	}
}

func (ft *UserFilters) ByPStringPtr(PString *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_string",
		Value:  PString,
		IsNil:  PString == nil,
		IsZero: PString != nil && (*PString) == "",
	}
}

type UserSubsetFilters struct{ prefix string }

func NewUserSubsetFilters(prefix string) UserSubsetFilters {
	return UserSubsetFilters{prefix}
}

func (ft *UserSubsetFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft UserSubsetFilters) Prefix() string {
	return ft.prefix
}

func (ft *UserSubsetFilters) ByID(ID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == "",
	}
}

func (ft *UserSubsetFilters) ByIDPtr(ID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == "",
	}
}

func (ft *UserSubsetFilters) ByBool(Bool bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "bool",
		Value:  Bool,
		IsNil:  bool(!Bool),
	}
}

func (ft *UserSubsetFilters) ByBoolPtr(Bool *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "bool",
		Value:  Bool,
		IsNil:  Bool == nil,
		IsZero: Bool != nil && bool(!(*Bool)),
	}
}

func (ft *UserSubsetFilters) ByFloat64(Float64 float64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "float64",
		Value:  Float64,
		IsNil:  Float64 == 0,
	}
}

func (ft *UserSubsetFilters) ByFloat64Ptr(Float64 *float64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "float64",
		Value:  Float64,
		IsNil:  Float64 == nil,
		IsZero: Float64 != nil && (*Float64) == 0,
	}
}

func (ft *UserSubsetFilters) ByInt(Int int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "int",
		Value:  Int,
		IsNil:  Int == 0,
	}
}

func (ft *UserSubsetFilters) ByIntPtr(Int *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "int",
		Value:  Int,
		IsNil:  Int == nil,
		IsZero: Int != nil && (*Int) == 0,
	}
}

func (ft *UserSubsetFilters) ByInt64(Int64 int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "int64",
		Value:  Int64,
		IsNil:  Int64 == 0,
	}
}

func (ft *UserSubsetFilters) ByInt64Ptr(Int64 *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "int64",
		Value:  Int64,
		IsNil:  Int64 == nil,
		IsZero: Int64 != nil && (*Int64) == 0,
	}
}

func (ft *UserSubsetFilters) ByString(String string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "string",
		Value:  String,
		IsNil:  String == "",
	}
}

func (ft *UserSubsetFilters) ByStringPtr(String *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "string",
		Value:  String,
		IsNil:  String == nil,
		IsZero: String != nil && (*String) == "",
	}
}

func (ft *UserSubsetFilters) ByPBoolPtr(PBool *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_bool",
		Value:  PBool,
		IsNil:  PBool == nil,
		IsZero: PBool != nil && bool(!(*PBool)),
	}
}

func (ft *UserSubsetFilters) ByPFloat64Ptr(PFloat64 *float64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_float64",
		Value:  PFloat64,
		IsNil:  PFloat64 == nil,
		IsZero: PFloat64 != nil && (*PFloat64) == 0,
	}
}

func (ft *UserSubsetFilters) ByPIntPtr(PInt *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_int",
		Value:  PInt,
		IsNil:  PInt == nil,
		IsZero: PInt != nil && (*PInt) == 0,
	}
}

func (ft *UserSubsetFilters) ByPInt64Ptr(PInt64 *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_int64",
		Value:  PInt64,
		IsNil:  PInt64 == nil,
		IsZero: PInt64 != nil && (*PInt64) == 0,
	}
}

func (ft *UserSubsetFilters) ByPStringPtr(PString *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_string",
		Value:  PString,
		IsNil:  PString == nil,
		IsZero: PString != nil && (*PString) == "",
	}
}

type UserInfoFilters struct{ prefix string }

func NewUserInfoFilters(prefix string) UserInfoFilters {
	return UserInfoFilters{prefix}
}

func (ft *UserInfoFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft UserInfoFilters) Prefix() string {
	return ft.prefix
}

func (ft *UserInfoFilters) ByUserID(UserID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == 0,
	}
}

func (ft *UserInfoFilters) ByUserIDPtr(UserID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == nil,
		IsZero: UserID != nil && (*UserID) == 0,
	}
}

func (ft *UserInfoFilters) ByMetadata(Metadata string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "metadata",
		Value:  Metadata,
		IsNil:  Metadata == "",
	}
}

func (ft *UserInfoFilters) ByMetadataPtr(Metadata *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "metadata",
		Value:  Metadata,
		IsNil:  Metadata == nil,
		IsZero: Metadata != nil && (*Metadata) == "",
	}
}

func (ft *UserInfoFilters) ByBool(Bool bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "bool",
		Value:  Bool,
		IsNil:  bool(!Bool),
	}
}

func (ft *UserInfoFilters) ByBoolPtr(Bool *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "bool",
		Value:  Bool,
		IsNil:  Bool == nil,
		IsZero: Bool != nil && bool(!(*Bool)),
	}
}

func (ft *UserInfoFilters) ByFloat64(Float64 float64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "float64",
		Value:  Float64,
		IsNil:  Float64 == 0,
	}
}

func (ft *UserInfoFilters) ByFloat64Ptr(Float64 *float64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "float64",
		Value:  Float64,
		IsNil:  Float64 == nil,
		IsZero: Float64 != nil && (*Float64) == 0,
	}
}

func (ft *UserInfoFilters) ByInt(Int int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "int",
		Value:  Int,
		IsNil:  Int == 0,
	}
}

func (ft *UserInfoFilters) ByIntPtr(Int *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "int",
		Value:  Int,
		IsNil:  Int == nil,
		IsZero: Int != nil && (*Int) == 0,
	}
}

func (ft *UserInfoFilters) ByInt64(Int64 int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "int64",
		Value:  Int64,
		IsNil:  Int64 == 0,
	}
}

func (ft *UserInfoFilters) ByInt64Ptr(Int64 *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "int64",
		Value:  Int64,
		IsNil:  Int64 == nil,
		IsZero: Int64 != nil && (*Int64) == 0,
	}
}

func (ft *UserInfoFilters) ByString(String string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "string",
		Value:  String,
		IsNil:  String == "",
	}
}

func (ft *UserInfoFilters) ByStringPtr(String *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "string",
		Value:  String,
		IsNil:  String == nil,
		IsZero: String != nil && (*String) == "",
	}
}

func (ft *UserInfoFilters) ByPBoolPtr(PBool *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_bool",
		Value:  PBool,
		IsNil:  PBool == nil,
		IsZero: PBool != nil && bool(!(*PBool)),
	}
}

func (ft *UserInfoFilters) ByPFloat64Ptr(PFloat64 *float64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_float64",
		Value:  PFloat64,
		IsNil:  PFloat64 == nil,
		IsZero: PFloat64 != nil && (*PFloat64) == 0,
	}
}

func (ft *UserInfoFilters) ByPIntPtr(PInt *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_int",
		Value:  PInt,
		IsNil:  PInt == nil,
		IsZero: PInt != nil && (*PInt) == 0,
	}
}

func (ft *UserInfoFilters) ByPInt64Ptr(PInt64 *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_int64",
		Value:  PInt64,
		IsNil:  PInt64 == nil,
		IsZero: PInt64 != nil && (*PInt64) == 0,
	}
}

func (ft *UserInfoFilters) ByPStringPtr(PString *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "p_string",
		Value:  PString,
		IsNil:  PString == nil,
		IsZero: PString != nil && (*PString) == "",
	}
}

type ComplexInfoFilters struct{ prefix string }

func NewComplexInfoFilters(prefix string) ComplexInfoFilters {
	return ComplexInfoFilters{prefix}
}

func (ft *ComplexInfoFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ComplexInfoFilters) Prefix() string {
	return ft.prefix
}

func (ft *ComplexInfoFilters) ByID(ID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == "",
	}
}

func (ft *ComplexInfoFilters) ByIDPtr(ID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == "",
	}
}

func (ft *ComplexInfoFilters) ByAliasString(AliasString m.AliasString) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "alias_string",
		Value:  AliasString,
		IsNil:  AliasString == "",
	}
}

func (ft *ComplexInfoFilters) ByAliasStringPtr(AliasString *m.AliasString) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_string",
		Value:  AliasString,
		IsNil:  AliasString == nil,
		IsZero: AliasString != nil && (*AliasString) == "",
	}
}

func (ft *ComplexInfoFilters) ByAliasInt64(AliasInt64 m.AliasInt64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "alias_int64",
		Value:  AliasInt64,
		IsNil:  AliasInt64 == 0,
	}
}

func (ft *ComplexInfoFilters) ByAliasInt64Ptr(AliasInt64 *m.AliasInt64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_int64",
		Value:  AliasInt64,
		IsNil:  AliasInt64 == nil,
		IsZero: AliasInt64 != nil && (*AliasInt64) == 0,
	}
}

func (ft *ComplexInfoFilters) ByAliasInt(AliasInt m.AliasInt) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "alias_int",
		Value:  AliasInt,
		IsNil:  AliasInt == 0,
	}
}

func (ft *ComplexInfoFilters) ByAliasIntPtr(AliasInt *m.AliasInt) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_int",
		Value:  AliasInt,
		IsNil:  AliasInt == nil,
		IsZero: AliasInt != nil && (*AliasInt) == 0,
	}
}

func (ft *ComplexInfoFilters) ByAliasBool(AliasBool m.AliasBool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "alias_bool",
		Value:  AliasBool,
		IsNil:  bool(!AliasBool),
	}
}

func (ft *ComplexInfoFilters) ByAliasBoolPtr(AliasBool *m.AliasBool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_bool",
		Value:  AliasBool,
		IsNil:  AliasBool == nil,
		IsZero: AliasBool != nil && bool(!(*AliasBool)),
	}
}

func (ft *ComplexInfoFilters) ByAliasFloat64(AliasFloat64 m.AliasFloat64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "alias_float64",
		Value:  AliasFloat64,
		IsNil:  AliasFloat64 == 0,
	}
}

func (ft *ComplexInfoFilters) ByAliasFloat64Ptr(AliasFloat64 *m.AliasFloat64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_float64",
		Value:  AliasFloat64,
		IsNil:  AliasFloat64 == nil,
		IsZero: AliasFloat64 != nil && (*AliasFloat64) == 0,
	}
}

func (ft *ComplexInfoFilters) ByAliasPStringPtr(AliasPString m.AliasPString) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_p_string",
		Value:  AliasPString,
		IsNil:  AliasPString == nil,
		IsZero: AliasPString != nil && (*AliasPString) == "",
	}
}

func (ft *ComplexInfoFilters) ByAliasPInt64Ptr(AliasPInt64 m.AliasPInt64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_p_int64",
		Value:  AliasPInt64,
		IsNil:  AliasPInt64 == nil,
		IsZero: AliasPInt64 != nil && (*AliasPInt64) == 0,
	}
}

func (ft *ComplexInfoFilters) ByAliasPIntPtr(AliasPInt m.AliasPInt) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_p_int",
		Value:  AliasPInt,
		IsNil:  AliasPInt == nil,
		IsZero: AliasPInt != nil && (*AliasPInt) == 0,
	}
}

func (ft *ComplexInfoFilters) ByAliasPBoolPtr(AliasPBool m.AliasPBool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_p_bool",
		Value:  AliasPBool,
		IsNil:  AliasPBool == nil,
		IsZero: AliasPBool != nil && bool(!(*AliasPBool)),
	}
}

func (ft *ComplexInfoFilters) ByAliasPFloat64Ptr(AliasPFloat64 m.AliasPFloat64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "alias_p_float64",
		Value:  AliasPFloat64,
		IsNil:  AliasPFloat64 == nil,
		IsZero: AliasPFloat64 != nil && (*AliasPFloat64) == 0,
	}
}

type UserTagFilters struct{ prefix string }

func NewUserTagFilters(prefix string) UserTagFilters {
	return UserTagFilters{prefix}
}

func (ft *UserTagFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft UserTagFilters) Prefix() string {
	return ft.prefix
}

func (ft *UserTagFilters) ByProvince(Province string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "province",
		Value:  Province,
		IsNil:  Province == "",
	}
}

func (ft *UserTagFilters) ByProvincePtr(Province *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "province",
		Value:  Province,
		IsNil:  Province == nil,
		IsZero: Province != nil && (*Province) == "",
	}
}

func (ft *UserTagFilters) ByRename(Rename string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "new_name",
		Value:  Rename,
		IsNil:  Rename == "",
	}
}

func (ft *UserTagFilters) ByRenamePtr(Rename *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "new_name",
		Value:  Rename,
		IsNil:  Rename == nil,
		IsZero: Rename != nil && (*Rename) == "",
	}
}

type UserInlineFilters struct{ prefix string }

func NewUserInlineFilters(prefix string) UserInlineFilters {
	return UserInlineFilters{prefix}
}

func (ft *UserInlineFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft UserInlineFilters) Prefix() string {
	return ft.prefix
}

func (ft *UserInlineFilters) ByProvince(Province string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "province",
		Value:  Province,
		IsNil:  Province == "",
	}
}

func (ft *UserInlineFilters) ByProvincePtr(Province *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "province",
		Value:  Province,
		IsNil:  Province == nil,
		IsZero: Province != nil && (*Province) == "",
	}
}

func (ft *UserInlineFilters) ByDistrict(District string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "district",
		Value:  District,
		IsNil:  District == "",
	}
}

func (ft *UserInlineFilters) ByDistrictPtr(District *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "district",
		Value:  District,
		IsNil:  District == nil,
		IsZero: District != nil && (*District) == "",
	}
}

type ProfileFilters struct{ prefix string }

func NewProfileFilters(prefix string) ProfileFilters {
	return ProfileFilters{prefix}
}

func (ft *ProfileFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ProfileFilters) Prefix() string {
	return ft.prefix
}

func (ft *ProfileFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ProfileFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ProfileFilters) ByStyle(Style string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "style",
		Value:  Style,
		IsNil:  Style == "",
	}
}

func (ft *ProfileFilters) ByStylePtr(Style *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "style",
		Value:  Style,
		IsNil:  Style == nil,
		IsZero: Style != nil && (*Style) == "",
	}
}
