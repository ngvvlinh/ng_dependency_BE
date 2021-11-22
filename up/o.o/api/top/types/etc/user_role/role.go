package user_role

// +enum
// +enum:zero=null
type UserRole int

type NullUserRole struct {
	Enum  UserRole
	Valid bool
}

const (
	// +enum=unknown
	MerchantUnknown UserRole = 0

	// +enum=owner
	MerchantSupperAdmin UserRole = 1

	// +enum=m_admin
	MerchantAdmin UserRole = 2
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
