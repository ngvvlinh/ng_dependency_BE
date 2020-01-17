package permission

// PermType declares permission type
type PermType int
type AuthType int
type AuthOpt int
type ActionType string

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
	Permissions string
	Validate    string
	Captcha     string
	AuthPartner AuthOpt
	Actions     []ActionType

	Rename string
}
