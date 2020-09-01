package authorization

import "o.o/capi/dot"

type Role string
type Action string

func (r Role) String() string   { return string(r) }
func (a Action) String() string { return string(a) }

var InternalRoles = []Role{RoleAdmin, RoleAdminSaleLead, RoleAdminSale, RoleAdminCustomerService, RoleAdminAccountant, RoleAdminCustomerServiceLead}

const (
	RoleAdmin                    Role = "admin"
	RoleAdminSaleLead            Role = "ad_salelead"
	RoleAdminSale                Role = "ad_sale"
	RoleAdminCustomerService     Role = "ad_customerservice"
	RoleAdminAccountant          Role = "ad_accountant"
	RoleAdminCustomerServiceLead Role = "ad_customerservice_lead"

	RoleInventoryManagement  Role = "inventory_management"
	RoleSalesMan             Role = "salesman"
	RoleShopOwner            Role = "owner"
	RoleAnalyst              Role = "analyst"
	RoleAccountant           Role = "accountant"
	RolePurchasingManagement Role = "purchasing_management"
	RoleStaffManagement      Role = "staff_management"
)

var Roles = [13]Role{
	RoleAdmin,
	RoleAdminSaleLead,
	RoleAdminSale,
	RoleAdminAccountant,
	RoleAdminCustomerService,
	RoleAdminCustomerServiceLead,
	RoleInventoryManagement,
	RoleSalesMan,
	RoleShopOwner,
	RoleAnalyst,
	RoleAccountant,
	RolePurchasingManagement,
	RoleStaffManagement,
}

var roleLabels = map[Role]string{
	RoleAdmin:                    "Quản trị viên",
	RoleAdminSaleLead:            "Trưởng Sale",
	RoleAdminSale:                "Sale",
	RoleAdminAccountant:          "Kế Toán",
	RoleAdminCustomerService:     "Chăm Sóc Khách Hàng",
	RoleAdminCustomerServiceLead: "Chăm Sóc Khách Hàng - Trưởng",
	RoleShopOwner:                "Chủ sở hữu",
	RoleStaffManagement:          "Quản lý nhân viên",
	RoleAnalyst:                  "Phân tích",
	RoleAccountant:               "Kế toán",
	RoleSalesMan:                 "Bán hàng",
	RoleInventoryManagement:      "Quản lý kho",
	RolePurchasingManagement:     "Thu mua",
}

func ParseRoleLabels(roles []Role) (result []string) {
	for _, role := range roles {
		result = append(result, roleLabels[role])
	}
	return
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
