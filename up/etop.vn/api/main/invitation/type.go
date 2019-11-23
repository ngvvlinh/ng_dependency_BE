package invitation

import (
	"time"

	"etop.vn/api/main/authorization"

	"github.com/dgrijalva/jwt-go"

	"etop.vn/api/main/etop"
	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/invitation

type Config struct {
	Secret string `yaml:"secret"`
}

type Invitation struct {
	ID         dot.ID
	AccountID  dot.ID
	Email      string
	FullName   string
	ShortName  string
	Position   string
	Roles      []authorization.Role
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
	Email          string               `json:"email"`
	AccountID      dot.ID               `json:"account_id"`
	Roles          []authorization.Role `json:"roles"`
	StandardClaims jwt.StandardClaims
}

func (c Claims) Valid() error {
	panic("implement me")
}

type InvitationAcceptedEvent struct {
	meta.EventMeta
	ID dot.ID
}
