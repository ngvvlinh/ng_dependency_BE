package invitation

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"etop.vn/api/main/etop"
	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/invitation

type Role string

const (
	RoleInventoryManagement  Role = "inventory_management"
	RoleSalesMan             Role = "salesman"
	RoleShopOwner            Role = "owner"
	RoleAnalyst              Role = "analyst"
	RolePurchasingManagement Role = "purchasing_management"
	RoleStaffManagement      Role = "staff_management"
)

var Roles = [6]Role{RoleInventoryManagement, RoleSalesMan, RoleShopOwner, RoleAnalyst, RolePurchasingManagement, RoleStaffManagement}

type Config struct {
	Secret string `yaml:"secret"`
}

type Invitation struct {
	ID         dot.ID
	AccountID  dot.ID
	Email      string
	Roles      []Role
	Token      string
	Status     etop.Status3
	InvitedBy  dot.ID
	AcceptedAt time.Time
	RejectedAt time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Claims struct {
	Email          string `json:"email"`
	AccountID      dot.ID `json:"account_id"`
	Roles          []Role `json:"roles"`
	StandardClaims jwt.StandardClaims
}

func (c Claims) Valid() error {
	panic("implement me")
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

type InvitationAcceptedEvent struct {
	meta.EventMeta
	ID dot.ID
}
