package account_type

// Indicates whether given account is **etop**, **shop**, **partner** or **sale**.
// +enum
type AccountType int

const (
	// +enum=unknown
	Unknown AccountType = 0

	// +enum=partner
	Partner AccountType = 21

	// +enum=shop
	Shop AccountType = 33

	// +enum=affiliate
	Affiliate AccountType = 35

	// +enum=etop
	Etop AccountType = 101
)
