package authorization

import "etop.vn/capi/dot"

type Role string
type Action string

const (
	RoleAdmin                Role = "admin"
	RoleInventoryManagement  Role = "inventory_management"
	RoleSalesMan             Role = "salesman"
	RoleShopOwner            Role = "owner"
	RoleAnalyst              Role = "analyst"
	RoleAccountant           Role = "accountant"
	RolePurchasingManagement Role = "purchasing_management"
	RoleStaffManagement      Role = "staff_management"
)

var Roles = [8]Role{
	RoleAdmin,
	RoleInventoryManagement,
	RoleSalesMan,
	RoleShopOwner,
	RoleAnalyst,
	RoleAccountant,
	RolePurchasingManagement,
	RoleStaffManagement,
}

type Authorization struct {
	UserID    dot.ID
	FullName  string
	ShortName string
	Position  string
	Email     string
	Roles     []Role
	Actions   []Action
}

type Relationship struct {
	UserID    dot.ID
	AccountID dot.ID
	FullName  string
	ShortName string
	Position  string
	Roles     []Role
	Actions   []Action
	Deleted   bool
}

func IsRole(arg Role) bool {
	for _, role := range Roles {
		if arg == role {
			return true
		}
	}
	return false
}

func IsContainsRole(roles []Role, arg Role) bool {
	for _, role := range roles {
		if role == arg {
			return true
		}
	}
	return false
}

func IsContainsActionString(actions []string, arg string) bool {
	for _, action := range actions {
		if action == arg {
			return true
		}
	}
	return false
}
