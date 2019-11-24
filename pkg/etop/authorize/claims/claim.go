package claims

import (
	"time"

	identitymodel "etop.vn/backend/com/main/identity/model"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

type ClaimInfo struct {
	Token         string `json:"-"`
	UserID        dot.ID `json:"-"`
	AdminID       dot.ID `json:"adm,omitempty"`
	AccountID     dot.ID `json:"acc,omitempty"`
	AuthPartnerID dot.ID `json:"auth_partner,omitempty"` // authenticated via partner

	SToken bool `json:"stoken,omitempty"`

	AccountIDs      map[dot.ID]int `json:"acs,omitempty"`
	STokenExpiresAt *time.Time     `json:"stoken_expires_at,omitempty"`

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

type AdminClaim struct {
	UserClaim
	CommonAccountClaim

	IsEtopAdmin bool
}

type AffiliateClaim struct {
	UserClaim
	CommonAccountClaim

	Affiliate *identitymodel.Affiliate
}
