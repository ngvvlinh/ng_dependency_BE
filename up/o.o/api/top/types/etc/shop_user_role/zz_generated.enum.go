// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package shop_user_role

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumUserRoleName = map[int]string{
	0:  "unknown",
	1:  "owner",
	2:  "staff_management",
	3:  "telecom_customerservice",
	4:  "inventory_management",
	5:  "purchasing_management",
	6:  "accountant",
	7:  "analyst",
	8:  "salesman",
	9:  "telecom_customerservice_management",
	10: "m_admin",
}

var enumUserRoleValue = map[string]int{
	"unknown":                            0,
	"owner":                              1,
	"staff_management":                   2,
	"telecom_customerservice":            3,
	"inventory_management":               4,
	"purchasing_management":              5,
	"accountant":                         6,
	"analyst":                            7,
	"salesman":                           8,
	"telecom_customerservice_management": 9,
	"m_admin":                            10,
}

func ParseUserRole(s string) (UserRole, bool) {
	val, ok := enumUserRoleValue[s]
	return UserRole(val), ok
}

func ParseUserRoleWithDefault(s string, d UserRole) UserRole {
	val, ok := enumUserRoleValue[s]
	if !ok {
		return d
	}
	return UserRole(val)
}

func (e UserRole) Apply(d UserRole) UserRole {
	if e == 0 {
		return d
	}
	return e
}

func (e UserRole) Enum() int {
	return int(e)
}

func (e UserRole) Name() string {
	return enumUserRoleName[e.Enum()]
}

func (e UserRole) String() string {
	s, ok := enumUserRoleName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("UserRole(%v)", e.Enum())
}

func (e UserRole) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumUserRoleName[e.Enum()] + "\""), nil
}

func (e *UserRole) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumUserRoleValue, data, "UserRole")
	if err != nil {
		return err
	}
	*e = UserRole(value)
	return nil
}

func (e UserRole) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *UserRole) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumUserRoleValue, src, "UserRole")
	*e = (UserRole)(value)
	return err
}

func (e UserRole) Wrap() NullUserRole {
	return WrapUserRole(e)
}

func ParseUserRoleWithNull(s dot.NullString, d UserRole) NullUserRole {
	if !s.Valid {
		return NullUserRole{}
	}
	val, ok := enumUserRoleValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return UserRole(val).Wrap()
}

func WrapUserRole(enum UserRole) NullUserRole {
	return NullUserRole{Enum: enum, Valid: true}
}

func (n NullUserRole) Apply(s UserRole) UserRole {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullUserRole) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullUserRole) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullUserRole) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullUserRole) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
