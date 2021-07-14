package shop_user_role

// +enum
// +enum:zero=null
type UserRole int

type NullUserRole struct {
	Enum  UserRole
	Valid bool
}

const (
	// +enum=owner
	Owner UserRole = 0

	// +enum=staff_management
	StaffManagement UserRole = 1

	// +enum=telecom_customerservice
	TelecomCustomerService UserRole = 2

	// +enum=inventory_management
	InventoryManagement UserRole = 3

	// +enum=purchasing_management
	PurchasingManagement UserRole = 4

	// +enum=accountant
	Accountant UserRole = 5

	// +enum=analyst
	Analyst UserRole = 6

	// +enum=salesman
	SalesMan UserRole = 7
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
