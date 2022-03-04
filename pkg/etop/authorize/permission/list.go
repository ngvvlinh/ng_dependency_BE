package permission

// PermType declares permission type
type PermType int
type AuthType int
type AuthOpt int
type ActionType string
type Actions []ActionType

func (as Actions) Contains(action ActionType) bool {
	for _, a := range as {
		if a == action {
			return true
		}
	}
	return false
}

// PermType constants
const (
	Public PermType = iota + 1
	Protected
	CurUsr
	Partner
	Shop

	Affiliate
	EtopAdmin

	SuperAdmin PermType = 100
	Secret     PermType = -2
)

const (
	User AuthType = iota + 1
	APIKey
	APIPartnerShopKey
	APIPartnerCarrierKey
)

const (
	Default AuthOpt = iota // Reject when auth is not Public, Protected
	Optional
	Required
)

func (a AuthOpt) AuthPartner() bool {
	return a != Default
}

type Decl struct {
	Type                 PermType
	Auth                 AuthType
	Permissions          string
	Validate             string
	Captcha              string
	AuthPartner          AuthOpt
	Actions              Actions
	IncludeFaboInfo      bool
	RequiredWhitelistIPs bool

	Rename string
}
