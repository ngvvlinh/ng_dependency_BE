// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package admin_user_role

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumAdminUserRoleName = map[int]string{
	0: "admin",
	1: "ad_customerservice_lead",
	2: "ad_salelead",
	3: "ad_accountant",
	4: "ad_customerservice",
	5: "ad_sale",
	6: "ad_voip",
	7: "ad_debug_mode",
}

var enumAdminUserRoleValue = map[string]int{
	"admin":                   0,
	"ad_customerservice_lead": 1,
	"ad_salelead":             2,
	"ad_accountant":           3,
	"ad_customerservice":      4,
	"ad_sale":                 5,
	"ad_voip":                 6,
	"ad_debug_mode":           7,
}

func ParseAdminUserRole(s string) (AdminUserRole, bool) {
	val, ok := enumAdminUserRoleValue[s]
	return AdminUserRole(val), ok
}

func ParseAdminUserRoleWithDefault(s string, d AdminUserRole) AdminUserRole {
	val, ok := enumAdminUserRoleValue[s]
	if !ok {
		return d
	}
	return AdminUserRole(val)
}

func (e AdminUserRole) Apply(d AdminUserRole) AdminUserRole {
	if e == 0 {
		return d
	}
	return e
}

func (e AdminUserRole) Enum() int {
	return int(e)
}

func (e AdminUserRole) Name() string {
	return enumAdminUserRoleName[e.Enum()]
}

func (e AdminUserRole) String() string {
	s, ok := enumAdminUserRoleName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("AdminUserRole(%v)", e.Enum())
}

func (e AdminUserRole) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumAdminUserRoleName[e.Enum()] + "\""), nil
}

func (e *AdminUserRole) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumAdminUserRoleValue, data, "AdminUserRole")
	if err != nil {
		return err
	}
	*e = AdminUserRole(value)
	return nil
}

func (e AdminUserRole) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *AdminUserRole) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumAdminUserRoleValue, src, "AdminUserRole")
	*e = (AdminUserRole)(value)
	return err
}

func (e AdminUserRole) Wrap() NullAdminUserRole {
	return WrapAdminUserRole(e)
}

func ParseAdminUserRoleWithNull(s dot.NullString, d AdminUserRole) NullAdminUserRole {
	if !s.Valid {
		return NullAdminUserRole{}
	}
	val, ok := enumAdminUserRoleValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return AdminUserRole(val).Wrap()
}

func WrapAdminUserRole(enum AdminUserRole) NullAdminUserRole {
	return NullAdminUserRole{Enum: enum, Valid: true}
}

func (n NullAdminUserRole) Apply(s AdminUserRole) AdminUserRole {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullAdminUserRole) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullAdminUserRole) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullAdminUserRole) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullAdminUserRole) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
