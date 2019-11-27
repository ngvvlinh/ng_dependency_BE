// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/capi/dot"
)

type ShopTraderFilters struct{ prefix string }

func NewShopTraderFilters(prefix string) ShopTraderFilters {
	return ShopTraderFilters{prefix}
}

func (ft *ShopTraderFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopTraderFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopTraderFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShopTraderFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShopTraderFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopTraderFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopTraderFilters) ByType(Type string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == "",
	}
}

func (ft *ShopTraderFilters) ByTypePtr(Type *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == "",
	}
}

type ShopCustomerFilters struct{ prefix string }

func NewShopCustomerFilters(prefix string) ShopCustomerFilters {
	return ShopCustomerFilters{prefix}
}

func (ft *ShopCustomerFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopCustomerFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopCustomerFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShopCustomerFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShopCustomerFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopCustomerFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopCustomerFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *ShopCustomerFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *ShopCustomerFilters) ByCodeNorm(CodeNorm int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code_norm",
		Value:  CodeNorm,
		IsNil:  CodeNorm == 0,
	}
}

func (ft *ShopCustomerFilters) ByCodeNormPtr(CodeNorm *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code_norm",
		Value:  CodeNorm,
		IsNil:  CodeNorm == nil,
		IsZero: CodeNorm != nil && (*CodeNorm) == 0,
	}
}

func (ft *ShopCustomerFilters) ByFullName(FullName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "full_name",
		Value:  FullName,
		IsNil:  FullName == "",
	}
}

func (ft *ShopCustomerFilters) ByFullNamePtr(FullName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "full_name",
		Value:  FullName,
		IsNil:  FullName == nil,
		IsZero: FullName != nil && (*FullName) == "",
	}
}

func (ft *ShopCustomerFilters) ByGender(Gender string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "gender",
		Value:  Gender,
		IsNil:  Gender == "",
	}
}

func (ft *ShopCustomerFilters) ByGenderPtr(Gender *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "gender",
		Value:  Gender,
		IsNil:  Gender == nil,
		IsZero: Gender != nil && (*Gender) == "",
	}
}

func (ft *ShopCustomerFilters) ByType(Type customering.CustomerType) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == "",
	}
}

func (ft *ShopCustomerFilters) ByTypePtr(Type *customering.CustomerType) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == "",
	}
}

func (ft *ShopCustomerFilters) ByBirthday(Birthday string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "birthday",
		Value:  Birthday,
		IsNil:  Birthday == "",
	}
}

func (ft *ShopCustomerFilters) ByBirthdayPtr(Birthday *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "birthday",
		Value:  Birthday,
		IsNil:  Birthday == nil,
		IsZero: Birthday != nil && (*Birthday) == "",
	}
}

func (ft *ShopCustomerFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *ShopCustomerFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *ShopCustomerFilters) ByPhone(Phone string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == "",
	}
}

func (ft *ShopCustomerFilters) ByPhonePtr(Phone *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == nil,
		IsZero: Phone != nil && (*Phone) == "",
	}
}

func (ft *ShopCustomerFilters) ByEmail(Email string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "email",
		Value:  Email,
		IsNil:  Email == "",
	}
}

func (ft *ShopCustomerFilters) ByEmailPtr(Email *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "email",
		Value:  Email,
		IsNil:  Email == nil,
		IsZero: Email != nil && (*Email) == "",
	}
}

func (ft *ShopCustomerFilters) ByStatus(Status int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ShopCustomerFilters) ByStatusPtr(Status *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ShopCustomerFilters) ByFullNameNorm(FullNameNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "full_name_norm",
		Value:  FullNameNorm,
		IsNil:  FullNameNorm == "",
	}
}

func (ft *ShopCustomerFilters) ByFullNameNormPtr(FullNameNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "full_name_norm",
		Value:  FullNameNorm,
		IsNil:  FullNameNorm == nil,
		IsZero: FullNameNorm != nil && (*FullNameNorm) == "",
	}
}

func (ft *ShopCustomerFilters) ByPhoneNorm(PhoneNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "phone_norm",
		Value:  PhoneNorm,
		IsNil:  PhoneNorm == "",
	}
}

func (ft *ShopCustomerFilters) ByPhoneNormPtr(PhoneNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "phone_norm",
		Value:  PhoneNorm,
		IsNil:  PhoneNorm == nil,
		IsZero: PhoneNorm != nil && (*PhoneNorm) == "",
	}
}

func (ft *ShopCustomerFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopCustomerFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopCustomerFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopCustomerFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShopCustomerFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ShopCustomerFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

type ShopTraderAddressFilters struct{ prefix string }

func NewShopTraderAddressFilters(prefix string) ShopTraderAddressFilters {
	return ShopTraderAddressFilters{prefix}
}

func (ft *ShopTraderAddressFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopTraderAddressFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopTraderAddressFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShopTraderAddressFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShopTraderAddressFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopTraderAddressFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopTraderAddressFilters) ByTraderID(TraderID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "trader_id",
		Value:  TraderID,
		IsNil:  TraderID == 0,
	}
}

func (ft *ShopTraderAddressFilters) ByTraderIDPtr(TraderID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "trader_id",
		Value:  TraderID,
		IsNil:  TraderID == nil,
		IsZero: TraderID != nil && (*TraderID) == 0,
	}
}

func (ft *ShopTraderAddressFilters) ByFullName(FullName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "full_name",
		Value:  FullName,
		IsNil:  FullName == "",
	}
}

func (ft *ShopTraderAddressFilters) ByFullNamePtr(FullName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "full_name",
		Value:  FullName,
		IsNil:  FullName == nil,
		IsZero: FullName != nil && (*FullName) == "",
	}
}

func (ft *ShopTraderAddressFilters) ByPhone(Phone string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == "",
	}
}

func (ft *ShopTraderAddressFilters) ByPhonePtr(Phone *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == nil,
		IsZero: Phone != nil && (*Phone) == "",
	}
}

func (ft *ShopTraderAddressFilters) ByEmail(Email string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "email",
		Value:  Email,
		IsNil:  Email == "",
	}
}

func (ft *ShopTraderAddressFilters) ByEmailPtr(Email *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "email",
		Value:  Email,
		IsNil:  Email == nil,
		IsZero: Email != nil && (*Email) == "",
	}
}

func (ft *ShopTraderAddressFilters) ByCompany(Company string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "company",
		Value:  Company,
		IsNil:  Company == "",
	}
}

func (ft *ShopTraderAddressFilters) ByCompanyPtr(Company *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "company",
		Value:  Company,
		IsNil:  Company == nil,
		IsZero: Company != nil && (*Company) == "",
	}
}

func (ft *ShopTraderAddressFilters) ByAddress1(Address1 string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "address1",
		Value:  Address1,
		IsNil:  Address1 == "",
	}
}

func (ft *ShopTraderAddressFilters) ByAddress1Ptr(Address1 *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "address1",
		Value:  Address1,
		IsNil:  Address1 == nil,
		IsZero: Address1 != nil && (*Address1) == "",
	}
}

func (ft *ShopTraderAddressFilters) ByAddress2(Address2 string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "address2",
		Value:  Address2,
		IsNil:  Address2 == "",
	}
}

func (ft *ShopTraderAddressFilters) ByAddress2Ptr(Address2 *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "address2",
		Value:  Address2,
		IsNil:  Address2 == nil,
		IsZero: Address2 != nil && (*Address2) == "",
	}
}

func (ft *ShopTraderAddressFilters) ByDistrictCode(DistrictCode string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "district_code",
		Value:  DistrictCode,
		IsNil:  DistrictCode == "",
	}
}

func (ft *ShopTraderAddressFilters) ByDistrictCodePtr(DistrictCode *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "district_code",
		Value:  DistrictCode,
		IsNil:  DistrictCode == nil,
		IsZero: DistrictCode != nil && (*DistrictCode) == "",
	}
}

func (ft *ShopTraderAddressFilters) ByWardCode(WardCode string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ward_code",
		Value:  WardCode,
		IsNil:  WardCode == "",
	}
}

func (ft *ShopTraderAddressFilters) ByWardCodePtr(WardCode *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ward_code",
		Value:  WardCode,
		IsNil:  WardCode == nil,
		IsZero: WardCode != nil && (*WardCode) == "",
	}
}

func (ft *ShopTraderAddressFilters) ByIsDefault(IsDefault bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "is_default",
		Value:  IsDefault,
		IsNil:  bool(!IsDefault),
	}
}

func (ft *ShopTraderAddressFilters) ByIsDefaultPtr(IsDefault *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "is_default",
		Value:  IsDefault,
		IsNil:  IsDefault == nil,
		IsZero: IsDefault != nil && bool(!(*IsDefault)),
	}
}

func (ft *ShopTraderAddressFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopTraderAddressFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopTraderAddressFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopTraderAddressFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShopTraderAddressFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ShopTraderAddressFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *ShopTraderAddressFilters) ByStatus(Status etop.Status3) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ShopTraderAddressFilters) ByStatusPtr(Status *etop.Status3) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

type ShopCustomerGroupCustomerFilters struct{ prefix string }

func NewShopCustomerGroupCustomerFilters(prefix string) ShopCustomerGroupCustomerFilters {
	return ShopCustomerGroupCustomerFilters{prefix}
}

func (ft *ShopCustomerGroupCustomerFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopCustomerGroupCustomerFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopCustomerGroupCustomerFilters) ByGroupID(GroupID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "group_id",
		Value:  GroupID,
		IsNil:  GroupID == 0,
	}
}

func (ft *ShopCustomerGroupCustomerFilters) ByGroupIDPtr(GroupID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "group_id",
		Value:  GroupID,
		IsNil:  GroupID == nil,
		IsZero: GroupID != nil && (*GroupID) == 0,
	}
}

func (ft *ShopCustomerGroupCustomerFilters) ByCustomerID(CustomerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "customer_id",
		Value:  CustomerID,
		IsNil:  CustomerID == 0,
	}
}

func (ft *ShopCustomerGroupCustomerFilters) ByCustomerIDPtr(CustomerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "customer_id",
		Value:  CustomerID,
		IsNil:  CustomerID == nil,
		IsZero: CustomerID != nil && (*CustomerID) == 0,
	}
}

func (ft *ShopCustomerGroupCustomerFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopCustomerGroupCustomerFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopCustomerGroupCustomerFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopCustomerGroupCustomerFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type ShopCustomerGroupFilters struct{ prefix string }

func NewShopCustomerGroupFilters(prefix string) ShopCustomerGroupFilters {
	return ShopCustomerGroupFilters{prefix}
}

func (ft *ShopCustomerGroupFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopCustomerGroupFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopCustomerGroupFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShopCustomerGroupFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShopCustomerGroupFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *ShopCustomerGroupFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *ShopCustomerGroupFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopCustomerGroupFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopCustomerGroupFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopCustomerGroupFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}
