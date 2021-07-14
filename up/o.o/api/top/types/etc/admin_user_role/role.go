package admin_user_role

// +enum
// +enum:zero=null
type AdminUserRole int

type NullAdminUserRole struct {
	Enum  AdminUserRole
	Valid bool
}

const (
	// +enum=admin
	Admin AdminUserRole = 0

	// +enum=ad_customerservice_lead
	AdminCustomerServiceLead AdminUserRole = 1

	// +enum=ad_salelead
	AdminSaleLead AdminUserRole = 2

	// +enum=ad_accountant
	AdminAccountant AdminUserRole = 3

	// +enum=ad_customerservice
	AdminCustomerService AdminUserRole = 4

	// +enum=ad_sale
	AdminSale AdminUserRole = 5

	// +enum=ad_voip
	AdminVoip AdminUserRole = 6
)
