package claims

import (
	"time"

	"etop.vn/backend/pkg/etop/model"
)

type ClaimInfo struct {
	Token         string `json:"-"`
	UserID        int64  `json:"-"`
	AdminID       int64  `json:"adm,omitempty"`
	AccountID     int64  `json:"acc,omitempty"`
	AuthPartnerID int64  `json:"auth_partner,omitempty"` // authenticated via partner

	SToken bool `json:"stoken,omitempty"`

	AccountIDs      map[int64]int `json:"acs,omitempty"`
	STokenExpiresAt *time.Time    `json:"stoken_expires_at,omitempty"`

	// check-and-set token for atomic updating
	CAS int64 `json:"cas,omitempty"`

	// extra value can be used for storing session data
	Extra map[string]string `json:"extra,omitempty"`
}

type ClaimInterface interface {
	GetClaim() *Claim
}

// Claim ...
type Claim struct {
	ClaimInfo

	LastLoginAt time.Time `json:"lla,omitempty"`
}

func (c *Claim) GetClaim() *Claim {
	return c
}

type CommonAccountClaim struct {
	IsOwner     bool
	Roles       []string
	Permissions []string
}

type EmptyClaim struct {
	*Claim

	IsSuperAdmin bool
}

type UserClaim struct {
	*Claim
	Admin *model.SignedInUser
	User  *model.SignedInUser
}

type PartnerClaim struct {
	UserClaim
	CommonAccountClaim
	User *model.SignedInUser

	Partner *model.Partner
}

type ShopClaim struct {
	UserClaim
	CommonAccountClaim

	Shop *model.Shop
}

type SupplierClaim struct {
	UserClaim
	CommonAccountClaim

	Supplier *model.Supplier
}

type AdminClaim struct {
	UserClaim
	CommonAccountClaim

	IsEtopAdmin bool
}
