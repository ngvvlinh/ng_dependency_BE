package permission

import "etop.vn/backend/pkg/etop/model"

// PermType declares permission type
type PermType int
type AuthType int
type AuthOpt int

// PermType constants
const (
	Public PermType = iota + 1
	Protected
	CurUsr
	Partner
	Shop
	EtopAdmin

	SuperAdmin PermType = 100
	Custom     PermType = -1
	Secret     PermType = -2
)

const (
	User AuthType = iota
	APIKey
	APIPartnerShopKey
)

const (
	Default AuthOpt = iota // Reject when auth is not Public, Protected
	Optional
	Required
)

// PermissionDecl ...
type PermissionDecl struct {
	Type        PermType
	Auth        AuthType
	Role        RoleType
	Permissions string
	Validate    string
	Captcha     string
	AuthPartner AuthOpt

	Rename string
}

type RoleType = model.RoleType

const (
	RoleOwner = model.RoleOwner
	RoleAdmin = model.RoleAdmin
	RoleStaff = model.RoleStaff
	Role3rd   = model.Role3rd
)

func RoleLevel(r RoleType) int {
	switch r {
	case model.RoleOwner:
		return 8
	case model.RoleAdmin:
		return 4
	case model.RoleStaff:
		return 2
	}
	return 0
}

func MaxRoleLevel(roles []string) int {
	if len(roles) == 0 {
		// backward-compatible for old tokens without roles TODO: remove
		return RoleLevel(RoleStaff)
	}
	if len(roles) == 1 {
		return RoleLevel(RoleType(roles[0]))
	}

	max := 0
	for _, role := range roles {
		if lvl := RoleLevel(RoleType(role)); lvl > max {
			max = lvl
		}
	}
	return max
}
