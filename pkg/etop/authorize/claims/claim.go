package claims

import (
	"time"

	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/capi/dot"
)

type ClaimInfo struct {
	Token         string `json:"-"`
	UserID        dot.ID `json:"-"`
	AdminID       dot.ID `json:"adm,omitempty"`
	AccountID     dot.ID `json:"acc,omitempty"`
	AuthPartnerID dot.ID `json:"auth_partner,omitempty"` // authenticated via partner

	// Chỉ sử dụng cho 1 số trường hợp session gán sẵn wl_partner_id
	// Trường hợp này sẽ gán luôn vào ctx để sử dụng
	WLPartnerID dot.ID `json:"wl_partner_id"`

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
	Admin *identitymodelx.SignedInUser
	User  *identitymodelx.SignedInUser
}

type PartnerClaim struct {
	UserClaim
	CommonAccountClaim
	User *identitymodelx.SignedInUser

	Partner *identitymodel.Partner
}

type ShopClaim struct {
	UserClaim
	CommonAccountClaim
	Actions []string

	Shop *identitymodel.Shop
}

type AdminClaim struct {
	UserClaim
	CommonAccountClaim
	Actions []string

	IsEtopAdmin bool
}

type AffiliateClaim struct {
	UserClaim
	CommonAccountClaim

	Affiliate *identitymodel.Affiliate
}
