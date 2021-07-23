package shop_user_role

// +enum
// +enum:zero=null
type UserRole int

type NullUserRole struct {
	Enum  UserRole
	Valid bool
}

const (
	// +enum=unknown
	Unknown UserRole = 0

	// +enum=owner
	Owner UserRole = 1

	// +enum=staff_management
	StaffManagement UserRole = 2

	// +enum=telecom_customerservice
	TelecomCustomerService UserRole = 3

	// +enum=inventory_management
	InventoryManagement UserRole = 4

	// +enum=purchasing_management
	PurchasingManagement UserRole = 5

	// +enum=accountant
	Accountant UserRole = 6

	// +enum=analyst
	Analyst UserRole = 7

	// +enum=salesman
	SalesMan UserRole = 8
)

func ContainsUserRoles(roles []UserRole, role UserRole) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func ToRolesString(roles []UserRole) (res []string) {
	for _, r := range roles {
		res = append(res, r.String())
	}
	return
}
