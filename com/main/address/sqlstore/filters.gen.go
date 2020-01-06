// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi/dot"
)

type AddressFilters struct{ prefix string }

func NewAddressFilters(prefix string) AddressFilters {
	return AddressFilters{prefix}
}

func (ft *AddressFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft AddressFilters) Prefix() string {
	return ft.prefix
}

func (ft *AddressFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *AddressFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *AddressFilters) ByFullName(FullName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "full_name",
		Value:  FullName,
		IsNil:  FullName == "",
	}
}

func (ft *AddressFilters) ByFullNamePtr(FullName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "full_name",
		Value:  FullName,
		IsNil:  FullName == nil,
		IsZero: FullName != nil && (*FullName) == "",
	}
}

func (ft *AddressFilters) ByFirstName(FirstName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "first_name",
		Value:  FirstName,
		IsNil:  FirstName == "",
	}
}

func (ft *AddressFilters) ByFirstNamePtr(FirstName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "first_name",
		Value:  FirstName,
		IsNil:  FirstName == nil,
		IsZero: FirstName != nil && (*FirstName) == "",
	}
}

func (ft *AddressFilters) ByLastName(LastName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "last_name",
		Value:  LastName,
		IsNil:  LastName == "",
	}
}

func (ft *AddressFilters) ByLastNamePtr(LastName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "last_name",
		Value:  LastName,
		IsNil:  LastName == nil,
		IsZero: LastName != nil && (*LastName) == "",
	}
}

func (ft *AddressFilters) ByPhone(Phone string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == "",
	}
}

func (ft *AddressFilters) ByPhonePtr(Phone *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == nil,
		IsZero: Phone != nil && (*Phone) == "",
	}
}

func (ft *AddressFilters) ByPosition(Position string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "position",
		Value:  Position,
		IsNil:  Position == "",
	}
}

func (ft *AddressFilters) ByPositionPtr(Position *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "position",
		Value:  Position,
		IsNil:  Position == nil,
		IsZero: Position != nil && (*Position) == "",
	}
}

func (ft *AddressFilters) ByEmail(Email string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "email",
		Value:  Email,
		IsNil:  Email == "",
	}
}

func (ft *AddressFilters) ByEmailPtr(Email *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "email",
		Value:  Email,
		IsNil:  Email == nil,
		IsZero: Email != nil && (*Email) == "",
	}
}

func (ft *AddressFilters) ByCountry(Country string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "country",
		Value:  Country,
		IsNil:  Country == "",
	}
}

func (ft *AddressFilters) ByCountryPtr(Country *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "country",
		Value:  Country,
		IsNil:  Country == nil,
		IsZero: Country != nil && (*Country) == "",
	}
}

func (ft *AddressFilters) ByCity(City string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "city",
		Value:  City,
		IsNil:  City == "",
	}
}

func (ft *AddressFilters) ByCityPtr(City *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "city",
		Value:  City,
		IsNil:  City == nil,
		IsZero: City != nil && (*City) == "",
	}
}

func (ft *AddressFilters) ByProvince(Province string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "province",
		Value:  Province,
		IsNil:  Province == "",
	}
}

func (ft *AddressFilters) ByProvincePtr(Province *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "province",
		Value:  Province,
		IsNil:  Province == nil,
		IsZero: Province != nil && (*Province) == "",
	}
}

func (ft *AddressFilters) ByDistrict(District string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "district",
		Value:  District,
		IsNil:  District == "",
	}
}

func (ft *AddressFilters) ByDistrictPtr(District *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "district",
		Value:  District,
		IsNil:  District == nil,
		IsZero: District != nil && (*District) == "",
	}
}

func (ft *AddressFilters) ByWard(Ward string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ward",
		Value:  Ward,
		IsNil:  Ward == "",
	}
}

func (ft *AddressFilters) ByWardPtr(Ward *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ward",
		Value:  Ward,
		IsNil:  Ward == nil,
		IsZero: Ward != nil && (*Ward) == "",
	}
}

func (ft *AddressFilters) ByZip(Zip string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "zip",
		Value:  Zip,
		IsNil:  Zip == "",
	}
}

func (ft *AddressFilters) ByZipPtr(Zip *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "zip",
		Value:  Zip,
		IsNil:  Zip == nil,
		IsZero: Zip != nil && (*Zip) == "",
	}
}

func (ft *AddressFilters) ByDistrictCode(DistrictCode string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "district_code",
		Value:  DistrictCode,
		IsNil:  DistrictCode == "",
	}
}

func (ft *AddressFilters) ByDistrictCodePtr(DistrictCode *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "district_code",
		Value:  DistrictCode,
		IsNil:  DistrictCode == nil,
		IsZero: DistrictCode != nil && (*DistrictCode) == "",
	}
}

func (ft *AddressFilters) ByProvinceCode(ProvinceCode string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "province_code",
		Value:  ProvinceCode,
		IsNil:  ProvinceCode == "",
	}
}

func (ft *AddressFilters) ByProvinceCodePtr(ProvinceCode *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "province_code",
		Value:  ProvinceCode,
		IsNil:  ProvinceCode == nil,
		IsZero: ProvinceCode != nil && (*ProvinceCode) == "",
	}
}

func (ft *AddressFilters) ByWardCode(WardCode string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ward_code",
		Value:  WardCode,
		IsNil:  WardCode == "",
	}
}

func (ft *AddressFilters) ByWardCodePtr(WardCode *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ward_code",
		Value:  WardCode,
		IsNil:  WardCode == nil,
		IsZero: WardCode != nil && (*WardCode) == "",
	}
}

func (ft *AddressFilters) ByCompany(Company string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "company",
		Value:  Company,
		IsNil:  Company == "",
	}
}

func (ft *AddressFilters) ByCompanyPtr(Company *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "company",
		Value:  Company,
		IsNil:  Company == nil,
		IsZero: Company != nil && (*Company) == "",
	}
}

func (ft *AddressFilters) ByAddress1(Address1 string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "address1",
		Value:  Address1,
		IsNil:  Address1 == "",
	}
}

func (ft *AddressFilters) ByAddress1Ptr(Address1 *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "address1",
		Value:  Address1,
		IsNil:  Address1 == nil,
		IsZero: Address1 != nil && (*Address1) == "",
	}
}

func (ft *AddressFilters) ByAddress2(Address2 string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "address2",
		Value:  Address2,
		IsNil:  Address2 == "",
	}
}

func (ft *AddressFilters) ByAddress2Ptr(Address2 *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "address2",
		Value:  Address2,
		IsNil:  Address2 == nil,
		IsZero: Address2 != nil && (*Address2) == "",
	}
}

func (ft *AddressFilters) ByType(Type string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == "",
	}
}

func (ft *AddressFilters) ByTypePtr(Type *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == "",
	}
}

func (ft *AddressFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *AddressFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *AddressFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *AddressFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *AddressFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *AddressFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}
